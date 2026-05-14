package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type UserRepository struct {
	mu        sync.Mutex
	users     map[int64]user.User
	usersByID map[user.ID]user.User
	nextID    int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:     map[int64]user.User{},
		usersByID: map[user.ID]user.User{},
		nextID:    1,
	}
}

func (r *UserRepository) FindOrCreateByGitHub(ctx context.Context, login authapp.GitHubLogin, now time.Time) (user.User, error) {
	if err := ctx.Err(); err != nil {
		return user.User{}, err
	}
	profile := login.Profile

	r.mu.Lock()
	defer r.mu.Unlock()

	appUser, ok := r.users[profile.GitHubID]
	if ok {
		updatedUser, err := appUser.LinkGitHub(profile, now)
		if err != nil {
			return user.User{}, err
		}
		r.users[profile.GitHubID] = updatedUser
		r.usersByID[updatedUser.ID] = updatedUser
		return updatedUser, nil
	}

	appUser, err := user.NewFromGitHub(user.ID(fmt.Sprintf("user_%d", r.nextID)), profile, now)
	if err != nil {
		return user.User{}, err
	}
	r.nextID++
	r.users[profile.GitHubID] = appUser
	r.usersByID[appUser.ID] = appUser

	return appUser, nil
}

func (r *UserRepository) findByID(ctx context.Context, id user.ID) (user.User, bool, error) {
	if err := ctx.Err(); err != nil {
		return user.User{}, false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	appUser, ok := r.usersByID[id]
	return appUser, ok, nil
}

type SessionRepository struct {
	mu       sync.Mutex
	sessions map[string]authapp.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{sessions: map[string]authapp.Session{}}
}

func (r *SessionRepository) Save(ctx context.Context, session authapp.Session) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionToken string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, sessionToken)
	return nil
}

func (r *SessionRepository) findByToken(ctx context.Context, token string) (authapp.Session, bool, error) {
	if err := ctx.Err(); err != nil {
		return authapp.Session{}, false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[token]
	if !ok {
		return authapp.Session{}, false, nil
	}

	return session, true, nil
}

type CurrentUserRepository struct {
	users    *UserRepository
	sessions *SessionRepository
}

func NewCurrentUserRepository(users *UserRepository, sessions *SessionRepository) *CurrentUserRepository {
	return &CurrentUserRepository{
		users:    users,
		sessions: sessions,
	}
}

func (r *CurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	session, ok, err := r.sessions.findByToken(ctx, sessionToken)
	if err != nil {
		return user.User{}, false, err
	}
	if !ok || !now.Before(session.ExpiresAt) {
		return user.User{}, false, nil
	}

	return r.users.findByID(ctx, session.UserID)
}
