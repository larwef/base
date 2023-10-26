package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/kelseyhightower/envconfig"
	"github.com/larwef/base/internal/server"
)

// Variables injected at compile time.
var (
	appName = "app"
	version = "No version provided"
)

type Config struct {
	Addr string `envconfig:"ADDRESS" default:":8080"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})).With(slog.Group("application", "name", appName, "version", version))
	slog.SetDefault(logger)

	logger.Info("starting application")

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx, logger)
	done()
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("program finished with error", "error", err)
	} else {
		logger.Info("program finished gracefully")
	}
}

func realMain(ctx context.Context, logger *slog.Logger) error {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		return err
	}

	srv := server.New(conf.Addr, server.WithLogger(logger))

	return srv.ListenAndServeContext(ctx)
}
