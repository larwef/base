package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"
)

// Version injected at compile time.
var version = "No version provided"

type Config struct{}

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx)
	done()
	if err != nil {
		slog.Error("program finished with error", "error", err)
	}
}

func realMain(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting base", "version", version)

	// TODO: Add your code here.

	return nil
}
