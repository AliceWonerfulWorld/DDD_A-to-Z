package pet

import (
	"context"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// CurrentUserRepository resolves a session token to a User.
type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

// CPBalanceReader provides the user's CP balance.
type CPBalanceReader interface {
	GetBalance(ctx context.Context, userID user.ID) (int64, error)
}

// PetReader provides read-only player pet data.
type PetReader interface {
	ListPetsByUser(ctx context.Context, userID user.ID) ([]PetWithGuild, error)
}

// CurrentGuildReader provides the user's active guild membership.
type CurrentGuildReader interface {
	FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error)
}

// PetWithGuild combines a pet with guild display data needed by the API.
type PetWithGuild struct {
	Pet   petdomain.Pet
	Guild guilddomain.Guild
}
