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

	"github.com/hadithopen-io/back/internal/auth"
	authhttp "github.com/hadithopen-io/back/internal/auth/dhttp"
	authpostgres "github.com/hadithopen-io/back/internal/auth/postgres"
	"github.com/hadithopen-io/back/internal/auth/sha256"
	"github.com/hadithopen-io/back/internal/config"
	"github.com/hadithopen-io/back/internal/story"
	"github.com/hadithopen-io/back/internal/story/dhttp"
	"github.com/hadithopen-io/back/internal/story/postgres"
	"github.com/hadithopen-io/back/internal/story/translate"
	"github.com/hadithopen-io/back/pkg/db/conn"
	"github.com/hadithopen-io/back/pkg/empty"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/http/middleware"
	"github.com/hadithopen-io/back/pkg/jwt"
	"github.com/hadithopen-io/back/pkg/jwt/hs256"
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

	slog.Info("init story object service")
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

	slog.Info("init tokenize/parser")
	tokenizeParser := hs256.NewHS256(
		conf.JWT.ExpiresTime,
		conf.JWT.RefreshTime,
		[]byte(
			// The mutable nature of secret may affect attempts to decode and encode data in the past.
			// You need to keep in mind that by changing the secret, previous users will not be able to log in.
			// Need migration when changing.
			conf.JWT.Secret,
		),
	)

	slog.Info("init authentication wrapper")
	authWrapper := jwt.NewWrapper(
		tokenizeParser,
	)

	slog.Info("init story delivery")
	storyDelivery := dhttp.NewStoryHandler(
		storyService,
		transmittersService,
		storyObjectService,
		authWrapper,
	)

	slog.Info("init story delivery handler")
	storyHandler, err := storyDelivery.Handler(
		middleware.UserLang,
	)
	if err != nil {
		return errors.Wrap(err, "after init story delivery handler")
	}

	slog.Info("init auth postgres repo")
	authUser := authpostgres.NewUser(
		dbconn,
	)

	slog.Info("init pwd encoder")
	pwdEncoder := sha256.NewEncoder(
		conf.Auth.Secret,
	)

	slog.Info("init auth service")
	authService := auth.NewAuth(
		authUser,
		tokenizeParser,
		pwdEncoder,
	)

	slog.Info("init auth delivery")
	authDelivery := authhttp.NewAuthHandler(
		authService,
		authWrapper,
	)

	slog.Info("init auth delivery handler")
	authHandler, err := authDelivery.Handler(
		middleware.UserLang,
	)
	if err != nil {
		return errors.Wrap(err, "after init auth delivery handler")
	}

	server := &http.Server{
		Addr:              conf.HTTP.Host,
		ReadHeaderTimeout: conf.HTTP.ReadHeaderTimeout,
	}

	http.Handle(
		storyDelivery.Path()+"/",
		storyHandler,
	)

	http.Handle(
		authDelivery.Path()+"/",
		authHandler,
	)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer slog.Info("stop listen http server")

		slog.Info("init http server", "host", conf.HTTP.Host)
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
