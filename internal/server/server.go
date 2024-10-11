package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Addr string `envconfig:"ADDRESS" default:":8080"`
}

// Wrapper around http.Server to avoid cluttering up the main function.
type Server struct {
	logger *slog.Logger
	srv    *http.Server
}

func New(logger *slog.Logger, conf Config, handler http.Handler) *Server {
	s := &Server{
		logger: logger,
		srv: &http.Server{
			Addr:         conf.Addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}

	return s
}

func (s *Server) ListenAndServeContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.logger.Info("starting server")
	errCh := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		if err := s.srv.Shutdown(context.Background()); err != nil {
			return err
		}
		s.logger.Info("server stopped gracefully")
		return ctx.Err()
	case err := <-errCh:
		s.logger.Info("server stopped unexpectedly", "error", err)
		return err
	}
}
