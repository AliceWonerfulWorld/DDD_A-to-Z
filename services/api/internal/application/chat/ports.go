package chat

import (
	"context"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type Repository interface {
	FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error)
	InsertChatToken(ctx context.Context, token ChatToken) error
}

type TokenGenerator interface {
	NewToken() (string, error)
}

type Hasher interface {
	Hash(token string) string
}
