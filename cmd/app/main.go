package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/kelseyhightower/envconfig"
	"github.com/larwef/base/internal/handler"
	"github.com/larwef/base/internal/server"
)

// Variables injected at compile time.
var (
	appName = "No name provided"
	version = "No version provided"
)

type Config struct {
	LogLvl    slog.Level `envconfig:"LOG_LEVEL" default:"info"`
	LogSource bool       `envconfig:"LOG_SOURCE"`
	LogJSON   bool       `envconfig:"LOG_JSON" default:"true"`

	Addr string `envconfig:"ADDRESS" default:":8080"`
}

func main() {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Set the default slog logger in stead of passing it around.
	logHandlerOpts := &slog.HandlerOptions{
		Level:     conf.LogLvl,
		AddSource: conf.LogSource,
	}

	var logHandler slog.Handler
	if conf.LogJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOpts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOpts)
	}
	logger := slog.
		New(logHandler).
		With(slog.Group("application", "name", appName, "version", version))
	slog.SetDefault(logger)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx, conf)
	done()
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("program finished with error", "error", err)
	} else {
		logger.Info("program finished gracefully")
	}
}

func realMain(ctx context.Context, conf Config) error {
	slog.Info("starting application")
	srv := server.New(conf.Addr, handler.Routes())
	return srv.ListenAndServeContext(ctx)
}
