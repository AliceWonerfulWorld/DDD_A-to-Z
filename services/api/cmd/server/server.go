package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/jyogi-web/ddd-a-to-z/gen/go/langwar/home/v1/homev1connect"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/database"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
)

func run(logger *slog.Logger) error {
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, config.DatabaseURLFromEnv())
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database handle: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	controllers, connectHandlers, err := buildControllers(logger, db)
	if err != nil {
		return fmt.Errorf("build controllers: %w", err)
	}

	server := newHTTPServer(logger, controllers, connectHandlers)
	logger.Info("api server listening", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func newHTTPServer(logger *slog.Logger, controllers controllerSet, connectHandlers connectHandlerSet) *http.Server {
	port := config.EnvOrDefault("PORT", "8080")
	mux := httpapi.NewRouter(logger, controllers.registrars()...)

	homePath, homeHandler := homev1connect.NewHomeServiceHandler(connectHandlers.home)
	mux.Handle(homePath, homeHandler)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}
}
