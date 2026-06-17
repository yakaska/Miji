package server

import (
	"Miji/internal/core/logger"
	"Miji/internal/core/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	handler http.Handler
	config  *Config
	logger  *logger.Logger
}

func NewServer(
	config *Config,
	logger *logger.Logger,
	mux *http.ServeMux,
	middlewares ...middleware.Middleware,
) *Server {
	var h http.Handler = mux
	h = middleware.Chain(h, middlewares...)

	return &Server{
		handler: h,
		config:  config,
		logger:  logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Address,
		Handler: s.handler,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		s.logger.Info(
			"Starting server",
			zap.String("address", s.config.Address),
		)

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		s.logger.Warn("Shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			if closeErr := server.Close(); closeErr != nil {
				s.logger.Error("server close error", zap.Error(closeErr))
			}
			return fmt.Errorf("server shutdown: %w", err)
		}
		s.logger.Warn("Server stopped")
		return nil
	}

	return nil
}
