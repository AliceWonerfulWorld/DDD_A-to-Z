package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

var (
	ErrMissingCode     = errors.New("auth code is required")
	ErrInvalidState    = errors.New("oauth state is invalid or expired")
	ErrUnauthenticated = errors.New("unauthenticated")
)

type UseCase struct {
	github     GitHubOAuthClient
	users      UserRepository
	sessions   SessionRepository
	current    CurrentUserRepository
	tokens     TokenGenerator
	now        func() time.Time
	stateTTL   time.Duration
	sessionTTL time.Duration
}

func NewUseCase(
	github GitHubOAuthClient,
	users UserRepository,
	sessions SessionRepository,
	current CurrentUserRepository,
	tokens TokenGenerator,
) *UseCase {
	return &UseCase{
		github:     github,
		users:      users,
		sessions:   sessions,
		current:    current,
		tokens:     tokens,
		now:        time.Now,
		stateTTL:   10 * time.Minute,
		sessionTTL: 24 * time.Hour,
	}
}

func (u *UseCase) BeginGitHubLogin(ctx context.Context) (LoginStart, error) {
	if err := ctx.Err(); err != nil {
		return LoginStart{}, err
	}

	state, err := u.tokens.NewToken()
	if err != nil {
		return LoginStart{}, err
	}

	expiresAt := u.now().Add(u.stateTTL)

	return LoginStart{
		AuthURL:        u.github.AuthCodeURL(state),
		State:          state,
		StateExpiresAt: expiresAt,
	}, nil
}

func (u *UseCase) CompleteGitHubLogin(ctx context.Context, code string) (LoginResult, error) {
	if code == "" {
		return LoginResult{}, ErrMissingCode
	}

	profile, err := u.github.ExchangeProfile(ctx, code)
	if err != nil {
		return LoginResult{}, err
	}

	now := u.now()
	appUser, err := u.users.FindOrCreateByGitHub(ctx, profile, now)
	if err != nil {
		return LoginResult{}, err
	}

	token, err := u.tokens.NewToken()
	if err != nil {
		return LoginResult{}, err
	}

	session := Session{
		Token:     token,
		UserID:    appUser.ID,
		ExpiresAt: now.Add(u.sessionTTL),
	}
	if err := u.sessions.Save(ctx, session); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		User:    appUser,
		Session: session,
	}, nil
}

func (u *UseCase) CurrentUser(ctx context.Context, sessionToken string) (user.User, error) {
	if sessionToken == "" {
		return user.User{}, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return user.User{}, err
	}
	if !ok {
		return user.User{}, ErrUnauthenticated
	}

	return appUser, nil
}

func (u *UseCase) Logout(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	return u.sessions.Delete(ctx, sessionToken)
}
