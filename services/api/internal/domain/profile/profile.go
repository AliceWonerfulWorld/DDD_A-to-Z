// Package profile owns the user-profile aggregate.
package profile

import (
	"errors"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

const maxDisplayNameLength = 50

// Profile represents the user-editable profile that is separate from the
// GitHub-linked identity (user.User). It is created once during initial
// on-boarding and can be updated later.
type Profile struct {
	UserID      user.ID
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates a new Profile, validating all invariants.
func New(userID user.ID, displayName string, now time.Time) (Profile, error) {
	if userID == "" {
		return Profile{}, errors.New("user id is required")
	}
	if err := validateDisplayName(displayName); err != nil {
		return Profile{}, err
	}
	return Profile{
		UserID:      userID,
		DisplayName: displayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func validateDisplayName(name string) error {
	if len(name) == 0 {
		return errors.New("display name is required")
	}
	if len([]rune(name)) > maxDisplayNameLength {
		return errors.New("display name must be 50 characters or fewer")
	}
	return nil
}
