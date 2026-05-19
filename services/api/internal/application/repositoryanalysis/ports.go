package repositoryanalysis

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
}

type GitHubCommitClient interface {
	ListCommits(ctx context.Context, accessToken, owner, repo, author string, since time.Time) ([]repositoryanalysis.CommitItem, error)
}

type GitHubPRClient interface {
	ListPullRequests(ctx context.Context, accessToken, owner, repo, author string, since time.Time) ([]repositoryanalysis.PullRequestItem, error)
}

type GitHubLanguageClient interface {
	ListLanguages(ctx context.Context, accessToken, owner, repo string) (map[string]int64, error)
}

type CPEarner interface {
	Earn(ctx context.Context, userID user.ID, amount int64, reason, sourceType, sourceID string) error
}

type SPEarner interface {
	EarnSP(ctx context.Context, userID user.ID, language string, amount int64, reason, sourceType, sourceID string) error
}

type CPBalanceProvider interface {
	GetBalance(ctx context.Context, userID user.ID) (int64, error)
	GetLastAnalyzedAt(ctx context.Context, userID user.ID) (*time.Time, error)
	UpdateLastAnalyzedAt(ctx context.Context, userID user.ID, at time.Time) error
}
