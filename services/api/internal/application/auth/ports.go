package auth

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type GitHubOAuthClient interface {
	AuthCodeURL(state string) string
	ExchangeProfile(ctx context.Context, code string) (user.GitHubProfile, error)
}

type UserRepository interface {
	FindOrCreateByGitHub(ctx context.Context, profile user.GitHubProfile, now time.Time) (user.User, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session Session) error
	Delete(ctx context.Context, sessionToken string) error
}

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type TokenGenerator interface {
	NewToken() (string, error)
}
