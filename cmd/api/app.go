package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/hadithopen-io/back/internal/config"
	"github.com/hadithopen-io/back/internal/story"
	"github.com/hadithopen-io/back/internal/story/dhttp"
	"github.com/hadithopen-io/back/internal/story/postgres"
	"github.com/hadithopen-io/back/internal/story/translate"
	"github.com/hadithopen-io/back/pkg/db/conn"
	"github.com/hadithopen-io/back/pkg/empty"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/http/middleware"
	"github.com/hadithopen-io/back/pkg/tx"
)

const (
	baseConfigPath = "./configs/main.yaml"

	configPathKey = "CONFIG_PATH"
)

func run() (
	err error,
) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				nil,
			),
		),
	)

	slog.Info("init app context")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		os.Interrupt,
	)
	defer cancel()

	slog.Info("init app config")
	conf, err := config.NewConfig(
		empty.Coalesce(
			os.Getenv(
				configPathKey,
			),
			baseConfigPath,
		),
	)
	if err != nil {
		return errors.Wrap(err, "after init app config")
	}

	slog.Info("init pgx connection")
	dbconn, err := conn.NewConn(ctx, conf.DB.Conn)
	if err != nil {
		return err
	}
	defer func() { err = errors.JoinCloser(err, dbconn) }()

	slog.Info("init hadith store")
	hadithStore := postgres.NewHadith(
		dbconn,
	)

	slog.Info("init story service")
	storyService := story.NewStory(
		hadithStore,
		nil,
	)

	slog.Info("init translate repo")
	translateStore := postgres.NewTranslate(
		dbconn,
	)

	slog.Info("init comment repo")
	comment := postgres.NewComment(
		dbconn,
	)

	slog.Info("init brought repo")
	brought := postgres.NewBrought(
		dbconn,
	)

	slog.Info("init version repo")
	version := postgres.NewVersion(
		dbconn,
	)

	slog.Info("init story object repo")
	hadithObject := postgres.NewStoryObject(
		dbconn,
	)

	slog.Info("init object tx")
	objectWrapper := tx.NewWrapper(
		postgres.NewObjectTX(
			dbconn,
		),
	)

	slog.Info("init translate service")
	translateService := translate.NewTranslate(
		translateStore,
	)

	storyObjectService := story.NewObject(
		translateService,
		comment,
		brought,
		version,
		hadithObject,
		objectWrapper,
	)

	slog.Info("init graph store")
	graphStore := postgres.NewGraph(
		dbconn,
	)

	slog.Info("init transmitters service")
	transmittersService := story.NewTransmitters(
		graphStore,
	)

	slog.Info("init story delivery")
	storyDelivery := dhttp.NewStoryHandler(
		storyService,
		transmittersService,
		storyObjectService,
	)

	slog.Info("init story delivery handler")
	storyHandler, err := storyDelivery.Handler(
		middleware.CookieAuth,
		middleware.QueryLang,
	)
	if err != nil {
		return errors.Wrap(err, "after init story handler")
	}

	server := &http.Server{
		Addr:              conf.API.Host,
		ReadHeaderTimeout: conf.HTTP.ReadHeaderTimeout,
	}

	http.Handle(
		storyDelivery.Path()+"/",
		storyHandler,
	)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer slog.Info("stop listen http server")

		slog.Info("init http server", "host", conf.API.Host)
		serr := server.ListenAndServe()
		if errors.Is(serr, http.ErrServerClosed) {
			return nil
		}
		return errors.Wrap(
			serr,
			"group listen http server",
		)
	})

	eg.Go(func() error {
		defer slog.Info("shutdown http server")

		<-ctx.Done()
		return errors.Wrap(
			server.Shutdown(
				ctx,
			),
			"group shutdown http server",
		)
	})

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "waited group")
	}

	slog.Info("shutdown app")
	time.Sleep(time.Second * 1)

	return nil
}
