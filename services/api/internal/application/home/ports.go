package home

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type ContributionPointReader interface {
	GetBalance(ctx context.Context, userID user.ID) (int64, error)
	GetTodayEarned(ctx context.Context, userID user.ID) (int64, error)
	GetTotalEarned(ctx context.Context, userID user.ID) (int64, error)
}
