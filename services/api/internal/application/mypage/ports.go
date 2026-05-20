package mypage

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// GitHubTokenRepository provides the user's GitHub access token.
type GitHubTokenRepository interface {
	GitHubAccessToken(ctx context.Context, userID user.ID) (string, bool, error)
}

// CurrentUserRepository resolves a session token to a User.
// The existing AuthStore already satisfies this interface.
type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

// ContributionPointReader provides read-only access to CP data.
type ContributionPointReader interface {
	// GetBalance returns the current CP balance.
	GetBalance(ctx context.Context, userID user.ID) (int64, error)
	// GetTotalEarned returns the lifetime total of earned CP.
	GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
	// GetTotalSpent returns the lifetime total of spent CP (as a positive value).
	GetTotalSpent(ctx context.Context, userID user.ID) (int64, error)
}

// RepositorySummaryReader provides a summarized view of repositories.
type RepositorySummaryReader interface {
	GetRepositorySummary(ctx context.Context, userID user.ID, recentLimit int) (RepositorySummary, error)
}

// GitHubStatsReader fetches aggregate GitHub user statistics.
type GitHubStatsReader interface {
	FetchStats(ctx context.Context, accessToken, username string) (*GitHubStats, error)
}
