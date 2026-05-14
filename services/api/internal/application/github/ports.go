package github

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type TokenRepository interface {
	GitHubAccessToken(ctx context.Context, userID user.ID) (string, bool, error)
}

type RepositoryClient interface {
	ListRepositories(ctx context.Context, accessToken string, userID user.ID, syncedAt time.Time) ([]repositoryanalysis.Repository, error)
}

type RepositoryStore interface {
	UpsertRepositories(ctx context.Context, repositories []repositoryanalysis.Repository) error
	ListRepositories(ctx context.Context, userID user.ID) ([]repositoryanalysis.Repository, error)
}
