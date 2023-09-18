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

	"github.com/hadithopen-io/back/internal/config"
	"github.com/hadithopen-io/back/internal/story"
	"github.com/hadithopen-io/back/internal/story/dhttp"
	"github.com/hadithopen-io/back/internal/story/postgre"
	"github.com/hadithopen-io/back/pkg/empty"
	"github.com/hadithopen-io/back/pkg/errors"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				nil,
			),
		),
	)
}

const baseConfigPath = "./configs/main.yaml"
const configPathKey = "CONFIG_PATH"

func run() (
	err error,
) {
	slog.Info("init app context")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
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

	slog.Info("init hadith store")
	hadithStore := postgre.NewHadith()

	slog.Info("init story service")
	storyService := story.NewStory(
		hadithStore,
		nil,
	)

	slog.Info("init graph store")
	graphStore := postgre.NewGraph()

	slog.Info("init transmitters service")
	transmittersService := story.NewTransmitters(
		graphStore,
	)

	slog.Info("init story handler")
	handler, err := dhttp.NewStoryHandler(storyService, transmittersService).Handler()
	if err != nil {
		return errors.Wrap(err, "after init story handler")
	}

	http.Handle("/", handler)

	eg, ctx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr: conf.API.Host,
	}

	eg.Go(func() error {
		slog.Info("init http server")

		return errors.Wrap(
			server.ListenAndServe(),
			"after init http listening",
		)
	})

	<-ctx.Done()
	slog.Info("shutdown app")

	if err := server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server")
	}

	time.Sleep(time.Second * 1)

	return eg.Wait()
}
