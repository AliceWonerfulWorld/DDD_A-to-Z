package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/database"
	infragithub "github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/postgres"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, config.DatabaseURLFromEnv())
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	authController := buildAuthController(logger, db)

	addr := config.EnvOrDefault("PORT", "8080")
	logger.Info("api server listening", "addr", ":"+addr)

	if err := http.ListenAndServe(":"+addr, httpapi.NewRouter(logger, authController)); err != nil {
		logger.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}

func buildAuthController(logger *slog.Logger, db *sql.DB) *httpapi.AuthController {
	oauthClient := infragithub.NewOAuthClient(config.GitHubOAuthFromEnv(), nil)
	authStore := postgres.NewAuthStore(db)
	usecase := authapp.NewUseCase(
		oauthClient,
		authStore,
		authStore,
		authStore,
		security.NewSecureTokenGenerator(),
	)

	return httpapi.NewAuthController(
		usecase,
		logger,
		security.NewSignedValueCodec(config.AuthCookieSecretFromEnv()),
		config.AuthCookieSecureFromEnv(),
	)
}
