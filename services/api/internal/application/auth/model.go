// Package auth owns login use cases and their ports.
package auth

import (
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type Session struct {
	Token     string
	UserID    user.ID
	ExpiresAt time.Time
}

type GitHubLogin struct {
	Profile     user.GitHubProfile
	AccessToken string
}

type LoginStart struct {
	AuthURL        string
	State          string
	StateExpiresAt time.Time
}

type LoginResult struct {
	User    user.User
	Session Session
}
