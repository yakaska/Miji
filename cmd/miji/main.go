package main

import (
	"Miji/internal/core/db"
	"Miji/internal/core/logger"
	"Miji/internal/core/middleware"
	"Miji/internal/core/server"
	"Miji/internal/links"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func run(ctx context.Context, getenv func(string) string, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logCfg, err := logger.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("logger config: %w", err)
	}
	log, err := logger.NewLogger(logCfg)
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer log.Close()

	log.Debug("Initializing database")
	dbCfg := db.NewConfigMust()
	dbPool, err := db.NewConnectionPool(ctx, dbCfg)
	if err != nil {
		log.Fatal("Failed to init connection pool", zap.Error(err))
	}
	defer dbPool.Close()

	log.Debug("Starting Miji - Url Shortener")

	linksHandler := links.NewHTTPHandler(
		links.NewService(links.NewPostgresRepository(dbPool)),
	)

	mux := http.NewServeMux()
	server.AddRoutes(mux, linksHandler)

	srv := server.NewServer(
		server.NewConfigMust(),
		log,
		mux,
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Panic(),
		middleware.Trace(),
	)

	return srv.Start(ctx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
