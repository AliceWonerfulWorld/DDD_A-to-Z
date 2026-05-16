package profile

import (
	"context"
	"time"

	domainprofile "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/profile"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// CurrentUserRepository resolves a session token to the authenticated user.
// The existing AuthStore already satisfies this interface.
type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

// ProfileRepository persists and retrieves user profiles.
type ProfileRepository interface {
	// Save creates or updates the profile for the given user.
	Save(ctx context.Context, p domainprofile.Profile) error
	// FindByUserID returns the profile for the given user, or (_, false, nil) if not found.
	FindByUserID(ctx context.Context, userID user.ID) (domainprofile.Profile, bool, error)
}
