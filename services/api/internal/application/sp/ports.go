package sp

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// CurrentUserRepository resolves a session token to a User.
// The existing AuthStore already satisfies this interface.
type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

// SPReader provides read-only access to SP balances.
type SPReader interface {
	GetSPBalances(ctx context.Context, userID user.ID) ([]SPBalance, error)
}
