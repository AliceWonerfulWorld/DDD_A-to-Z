package main

import (
	"log/slog"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guildtownapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guildtown"
	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	profileapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/profile"
	analysisapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/repositoryanalysis"
	spapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/sp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	infragithub "github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/postgres"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
	"gorm.io/gorm"
)

type controllerSet struct {
	auth       *httpapi.AuthController
	repository *httpapi.RepositoryController
	guild      *httpapi.GuildController
	guildTown  *httpapi.GuildTownController
	mypage     *httpapi.MypageController
	profile    *httpapi.ProfileController
	analysis   *httpapi.AnalysisController
	home       *httpapi.HomeController
	sp         *httpapi.SPController
}

func (c controllerSet) registrars() []httpapi.RouteRegistrar {
	return []httpapi.RouteRegistrar{
		c.auth,
		c.repository,
		c.guild,
		c.guildTown,
		c.mypage,
		c.profile,
		c.analysis,
		c.home,
		c.sp,
	}
}

func buildControllers(logger *slog.Logger, db *gorm.DB) (controllerSet, error) {
	settings, err := loadControllerSettings()
	if err != nil {
		return controllerSet{}, err
	}

	tokenCipher, err := security.NewTokenCipher(settings.githubTokenSecret)
	if err != nil {
		return controllerSet{}, err
	}

	oauthClient := infragithub.NewOAuthClient(settings.oauthConfig, nil)
	repositoryClient := infragithub.NewRepositoryClient(nil)

	authStore := postgres.NewAuthStore(db, tokenCipher)
	repositoryStore := postgres.NewRepositoryStore(db)
	contributionPointStore := postgres.NewContributionPointStore(db)
	mypageStore := postgres.NewMyPageStore(db)
	profileStore := postgres.NewProfileStore(db)
	guildStore, err := postgres.NewGuildStore(db)
	if err != nil {
		return controllerSet{}, err
	}
	guildTownStore := postgres.NewGuildTownStore(db)

	authUseCase := authapp.NewUseCase(
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
	cpLedgerIDGenerator := security.NewIDGenerator("cp")
	cpUseCase := contributionpointapp.NewUseCase(contributionPointStore, cpLedgerIDGenerator)
	cpManager := newCPManager(cpUseCase)
	guildUseCase := guildapp.NewUseCaseWithCPTransaction(
		guildStore,
		authStore,
		security.NewIDGenerator("guild_membership"),
		security.NewIDGenerator("guild_cp_contribution"),
		cpUseCase,
		postgres.NewGuildCPContributionTransactioner(db, cpLedgerIDGenerator),
	)
	guildTownUseCase := guildtownapp.NewUseCase(
		guildTownStore,
		authStore,
		guildStore,
		security.NewIDGenerator("guild_town_placement"),
	)
	mypageUseCase := mypageapp.NewUseCase(
		authStore,
		newMypageCPReader(contributionPointStore, mypageStore),
		mypageStore,
	)
	spUseCase := spapp.NewUseCase(authStore, contributionPointStore)
	homeCPProvider := newHomeCPDataProvider(contributionPointStore, mypageStore)
	profileUseCase := profileapp.NewUseCase(
		authStore,
		profileStore,
	)

	analysisUseCase := analysisapp.NewUseCase(
		authStore,
		authStore,
		repositoryClient,
		repositoryStore,
		repositoryClient,
		repositoryClient,
		repositoryClient,
		cpManager,
		cpManager,
		cpManager,
	)

	return controllerSet{
		auth: httpapi.NewAuthController(
			authUseCase,
			logger,
			security.NewSignedValueCodec(settings.cookieSecret),
			settings.cookieSecure,
			settings.frontendURL,
		),
		repository: httpapi.NewRepositoryController(repositoryUseCase, logger),
		guild:      httpapi.NewGuildController(guildUseCase, logger),
		guildTown:  httpapi.NewGuildTownController(guildTownUseCase, logger),
		mypage:     httpapi.NewMypageController(mypageUseCase, logger),
		profile:    httpapi.NewProfileController(profileUseCase, logger),
		analysis:   httpapi.NewAnalysisController(newAnalysisGuard(analysisUseCase, authStore), logger),
		home:       httpapi.NewHomeController(authStore, homeCPProvider, logger),
		sp:         httpapi.NewSPController(spUseCase, logger),
	}, nil
}

type controllerSettings struct {
	oauthConfig       config.GitHubOAuth
	cookieSecret      string
	cookieSecure      bool
	githubTokenSecret string
	frontendURL       string
}

func loadControllerSettings() (controllerSettings, error) {
	oauthConfig, err := config.GitHubOAuthFromEnv()
	if err != nil {
		return controllerSettings{}, err
	}

	cookieSecret, err := config.AuthCookieSecretFromEnv()
	if err != nil {
		return controllerSettings{}, err
	}

	cookieSecure, err := config.AuthCookieSecureFromEnv()
	if err != nil {
		return controllerSettings{}, err
	}

	tokenSecret, err := config.GitHubTokenEncryptionSecretFromEnv()
	if err != nil {
		return controllerSettings{}, err
	}

	return controllerSettings{
		oauthConfig:       oauthConfig,
		cookieSecret:      cookieSecret,
		cookieSecure:      cookieSecure,
		githubTokenSecret: tokenSecret,
		frontendURL:       config.EnvOrDefault("FRONTEND_URL", "http://localhost:5173"),
	}, nil
}
