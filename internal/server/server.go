package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Wrapper around http.Server to avoid cluttering up the main function.
type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

type ServerOption func(*Server)

func New(address string, options ...ServerOption) *Server {
	s := &Server{
		srv: &http.Server{
			Addr:         address,
			Handler:      routes(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		logger: slog.Default(),
	}

	for _, opt := range options {
		opt(s)
	}

	s.logger = s.logger.With(slog.Group("component",
		slog.Group("server", "address", address),
	))

	return s
}

func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

func (s *Server) ListenAndServeContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.logger.Info("starting server")
	errCh := make(chan error)

	select {
	case <-ctx.Done():
		if err := s.srv.Shutdown(context.Background()); err != nil {
			return err
		}
		s.logger.Info("server stopped gracefully")
		return ctx.Err()
	case err := <-errCh:
		s.logger.Warn("server stopped unexpectedly", "error", err)
		return err
	}
}
