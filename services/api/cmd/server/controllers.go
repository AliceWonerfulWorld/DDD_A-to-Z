package main

import (
	"log/slog"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	chatapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/chat"
	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guildtownapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guildtown"
	homeapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/home"
	mypageapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	profileapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/profile"
	analysisapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/repositoryanalysis"
	spapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/sp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	infragithub "github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/postgres"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
	connectapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/connect"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
	"gorm.io/gorm"
)

type connectHandlerSet struct {
	home *connectapi.HomeHandler
}

type controllerSet struct {
	auth       *httpapi.AuthController
	repository *httpapi.RepositoryController
	guild      *httpapi.GuildController
	guildTown  *httpapi.GuildTownController
	mypage     *httpapi.MypageController
	pet        *httpapi.PetController
	profile    *httpapi.ProfileController
	analysis   *httpapi.AnalysisController
	home       *httpapi.HomeController
	sp         *httpapi.SPController
	chat       *httpapi.ChatController
}

func (c controllerSet) registrars() []httpapi.RouteRegistrar {
	return []httpapi.RouteRegistrar{
		c.auth,
		c.repository,
		c.guild,
		c.guildTown,
		c.mypage,
		c.pet,
		c.profile,
		c.analysis,
		c.home,
		c.sp,
		c.chat,
	}
}

func buildControllers(logger *slog.Logger, db *gorm.DB) (controllerSet, connectHandlerSet, error) {
	settings, err := loadControllerSettings()
	if err != nil {
		return controllerSet{}, connectHandlerSet{}, err
	}

	tokenCipher, err := security.NewTokenCipher(settings.githubTokenSecret)
	if err != nil {
		return controllerSet{}, connectHandlerSet{}, err
	}

	oauthClient := infragithub.NewOAuthClient(settings.oauthConfig, nil)
	repositoryClient := infragithub.NewRepositoryClient(nil)

	authStore := postgres.NewAuthStore(db, tokenCipher)
	repositoryStore := postgres.NewRepositoryStore(db)
	contributionPointStore := postgres.NewContributionPointStore(db)
	mypageStore := postgres.NewMyPageStore(db)
	petStore := postgres.NewPetStore(db)
	profileStore := postgres.NewProfileStore(db)
	guildStore, err := postgres.NewGuildStore(db)
	if err != nil {
		return controllerSet{}, connectHandlerSet{}, err
	}
	guildTownStore := postgres.NewGuildTownStore(db)
	chatStore := postgres.NewChatStore(db, guildStore)

	authUseCase := authapp.NewUseCase(
		oauthClient,
		authStore,
		authStore,
		authStore,
		security.NewObfuscatedTokenGenerator(
			security.NewSecureTokenGenerator(),
			security.NewAwkTextMixer(),
			settings.tokenMixerSalt,
		),
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
	guildUseCase := guildapp.NewUseCaseWithPetAndCPTransaction(
		guildStore,
		authStore,
		security.NewIDGenerator("guild_membership"),
		security.NewIDGenerator("pet"),
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
		repositoryClient,
		authStore,
		mypageapp.NewGuildMembershipReader(guildStore),
		mypageapp.NewProfileReader(profileStore),
	)
	petUseCase := petapp.NewUseCase(
		authStore,
		newMypageCPReader(contributionPointStore, mypageStore),
		petStore,
		guildStore,
	)
	spUseCase := spapp.NewUseCase(authStore, contributionPointStore)
	homeCPProvider := newHomeCPDataProvider(contributionPointStore, mypageStore)
	homeUseCase := homeapp.NewUseCase(authStore, homeCPProvider)
	profileUseCase := profileapp.NewUseCase(
		authStore,
		profileStore,
	)

	analysisUseCase := analysisapp.NewUseCaseWithContributionStore(
		authStore,
		authStore,
		repositoryClient,
		repositoryStore,
		repositoryStore,
		repositoryClient,
		repositoryClient,
		repositoryClient,
		cpManager,
		cpManager,
		cpManager,
	)
	chatUseCase := chatapp.NewUseCase(
		authStore,
		chatStore,
		security.NewSecureTokenGenerator(),
		security.NewSHA256Hasher(),
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
			pet:        httpapi.NewPetController(petUseCase, logger),
			profile:    httpapi.NewProfileController(profileUseCase, logger),
			analysis:   httpapi.NewAnalysisController(newAnalysisGuard(analysisUseCase, authStore), logger),
			home:       httpapi.NewHomeController(homeUseCase, logger),
			sp:         httpapi.NewSPController(spUseCase, logger),
			chat:       httpapi.NewChatController(chatUseCase, logger),
		},
		connectHandlerSet{
			home: connectapi.NewHomeHandler(homeUseCase),
		},
		nil
}

type controllerSettings struct {
	oauthConfig       config.GitHubOAuth
	cookieSecret      string
	cookieSecure      bool
	githubTokenSecret string
	frontendURL       string
	tokenMixerSalt    string
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
		tokenMixerSalt:    config.EnvOrDefault("TOKEN_MIXER_SALT", cookieSecret),
	}, nil
}
