package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/database"
	infragithub "github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/postgres"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, config.DatabaseURLFromEnv())
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get database handle", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	authController, repositoryController, guildController, err := buildControllers(logger, db)
	if err != nil {
		logger.Error("failed to build controllers", "error", err)
		os.Exit(1)
	}

	addr := config.EnvOrDefault("PORT", "8080")
	logger.Info("api server listening", "addr", ":"+addr)

	server := &http.Server{
		Addr:              ":" + addr,
		Handler:           httpapi.NewRouter(logger, authController, repositoryController, guildController),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}

func buildControllers(logger *slog.Logger, db *gorm.DB) (*httpapi.AuthController, *httpapi.RepositoryController, *httpapi.GuildController, error) {
	oauthConfig, err := config.GitHubOAuthFromEnv()
	if err != nil {
		return nil, nil, nil, err
	}
	cookieSecret, err := config.AuthCookieSecretFromEnv()
	if err != nil {
		return nil, nil, nil, err
	}
	cookieSecure, err := config.AuthCookieSecureFromEnv()
	if err != nil {
		return nil, nil, nil, err
	}
	tokenSecret, err := config.GitHubTokenEncryptionSecretFromEnv()
	if err != nil {
		return nil, nil, nil, err
	}
	tokenCipher, err := security.NewTokenCipher(tokenSecret)
	if err != nil {
		return nil, nil, nil, err
	}
	frontendURL := config.EnvOrDefault("FRONTEND_URL", "http://localhost:5173")

	oauthClient := infragithub.NewOAuthClient(oauthConfig, nil)
	repositoryClient := infragithub.NewRepositoryClient(nil)
	authStore := postgres.NewAuthStore(db, tokenCipher)
	repositoryStore := postgres.NewRepositoryStore(db)
	guildStore, err := postgres.NewGuildStore(db)
	if err != nil {
		return nil, nil, nil, err
	}
	usecase := authapp.NewUseCase(
		oauthClient,
		authStore,
		authStore,
		authStore,
		security.NewSecureTokenGenerator(),
	)
	repositoryUseCase := githubapp.NewUseCase(
		authStore,
		authStore,
		repositoryClient,
		repositoryStore,
	)
	guildUseCase := guildapp.NewUseCase(guildStore)

	authController := httpapi.NewAuthController(
		usecase,
		logger,
		security.NewSignedValueCodec(cookieSecret),
		cookieSecure,
		frontendURL,
	)
	repositoryController := httpapi.NewRepositoryController(repositoryUseCase, logger)
	guildController := httpapi.NewGuildController(guildUseCase, logger)

	return authController, repositoryController, guildController, nil
}
