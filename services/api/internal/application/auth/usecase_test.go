package auth

import (
	"context"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeGitHubOAuthClient struct {
	authURL string
	profile user.GitHubProfile
}

func (c fakeGitHubOAuthClient) AuthCodeURL(state string) string {
	return c.authURL + "?state=" + state
}

func (c fakeGitHubOAuthClient) ExchangeProfile(ctx context.Context, code string) (user.GitHubProfile, error) {
	if err := ctx.Err(); err != nil {
		return user.GitHubProfile{}, err
	}

	return c.profile, nil
}

type fakeTokenGenerator struct {
	tokens []string
}

func (g *fakeTokenGenerator) NewToken() (string, error) {
	token := g.tokens[0]
	g.tokens = g.tokens[1:]
	return token, nil
}

type fakeUserRepository struct {
	users     map[int64]user.User
	usersByID map[user.ID]user.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users:     map[int64]user.User{},
		usersByID: map[user.ID]user.User{},
	}
}

func (r *fakeUserRepository) FindOrCreateByGitHub(ctx context.Context, profile user.GitHubProfile, now time.Time) (user.User, error) {
	if err := ctx.Err(); err != nil {
		return user.User{}, err
	}

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

	appUser, err := user.NewFromGitHub("user_1", profile, now)
	if err != nil {
		return user.User{}, err
	}
	r.users[profile.GitHubID] = appUser
	r.usersByID[appUser.ID] = appUser
	return appUser, nil
}

func (r *fakeUserRepository) findByID(ctx context.Context, id user.ID) (user.User, bool, error) {
	if err := ctx.Err(); err != nil {
		return user.User{}, false, err
	}

	appUser, ok := r.usersByID[id]
	return appUser, ok, nil
}

type fakeSessionRepository struct {
	sessions []Session
}

func (r *fakeSessionRepository) Save(ctx context.Context, session Session) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.sessions = append(r.sessions, session)
	return nil
}

func (r *fakeSessionRepository) Delete(ctx context.Context, sessionToken string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	for i, session := range r.sessions {
		if session.Token == sessionToken {
			r.sessions = append(r.sessions[:i], r.sessions[i+1:]...)
			return nil
		}
	}

	return nil
}

func (r *fakeSessionRepository) findByToken(ctx context.Context, token string, now time.Time) (Session, bool, error) {
	if err := ctx.Err(); err != nil {
		return Session{}, false, err
	}

	for _, session := range r.sessions {
		if session.Token == token && now.Before(session.ExpiresAt) {
			return session, true, nil
		}
	}

	return Session{}, false, nil
}

type fakeCurrentUserRepository struct {
	users    *fakeUserRepository
	sessions *fakeSessionRepository
}

func (r fakeCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	session, ok, err := r.sessions.findByToken(ctx, sessionToken, now)
	if err != nil {
		return user.User{}, false, err
	}
	if !ok {
		return user.User{}, false, nil
	}

	return r.users.findByID(ctx, session.UserID)
}

func TestUseCaseCompleteGitHubLogin(t *testing.T) {
	t.Run("GitHubログイン完了でユーザーセッションを作成できる", func(t *testing.T) {
		tokens := &fakeTokenGenerator{tokens: []string{"state-token", "session-token"}}
		users := newFakeUserRepository()
		sessions := &fakeSessionRepository{}
		usecase := NewUseCase(
			fakeGitHubOAuthClient{
				authURL: "https://github.com/login/oauth/authorize",
				profile: user.GitHubProfile{
					GitHubID:  123,
					Username:  "octocat",
					AvatarURL: "https://example.com/avatar.png",
				},
			},
			users,
			sessions,
			fakeCurrentUserRepository{users: users, sessions: sessions},
			tokens,
		)
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
		}

		start, err := usecase.BeginGitHubLogin(context.Background())
		if err != nil {
			t.Fatalf("BeginGitHubLogin がエラーを返しました: %v", err)
		}
		if start.State != "state-token" {
			t.Fatalf("state = %q, 期待値 state-token", start.State)
		}

		result, err := usecase.CompleteGitHubLogin(context.Background(), "code")
		if err != nil {
			t.Fatalf("CompleteGitHubLogin がエラーを返しました: %v", err)
		}

		if result.User.ID == "" {
			t.Fatal("ユーザーIDが設定されている必要があります")
		}
		if result.User.GitHubAccount.GitHubID != 123 {
			t.Fatalf("GitHub ID = %d, 期待値 123", result.User.GitHubAccount.GitHubID)
		}
		if result.User.GitHubAccount.Username != "octocat" {
			t.Fatalf("ユーザー名 = %q, 期待値 octocat", result.User.GitHubAccount.Username)
		}
		if result.Session.Token != "session-token" {
			t.Fatalf("セッショントークン = %q, 期待値 session-token", result.Session.Token)
		}
		if result.Session.UserID != result.User.ID {
			t.Fatalf("セッションのユーザーID = %q, 期待値 %q", result.Session.UserID, result.User.ID)
		}
	})
}

func TestUseCaseCurrentUser(t *testing.T) {
	t.Run("セッショントークンから現在のユーザーを取得できる", func(t *testing.T) {
		users := newFakeUserRepository()
		sessions := &fakeSessionRepository{}
		usecase := NewUseCase(
			fakeGitHubOAuthClient{
				authURL: "https://github.com/login/oauth/authorize",
				profile: user.GitHubProfile{
					GitHubID: 1,
					Username: "octocat",
				},
			},
			users,
			sessions,
			fakeCurrentUserRepository{users: users, sessions: sessions},
			&fakeTokenGenerator{tokens: []string{"session-token"}},
		)
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
		}

		login, err := usecase.CompleteGitHubLogin(context.Background(), "code")
		if err != nil {
			t.Fatalf("CompleteGitHubLogin がエラーを返しました: %v", err)
		}

		appUser, err := usecase.CurrentUser(context.Background(), "session-token")
		if err != nil {
			t.Fatalf("CurrentUser がエラーを返しました: %v", err)
		}
		if appUser.ID != login.User.ID {
			t.Fatalf("現在のユーザーID = %q, 期待値 %q", appUser.ID, login.User.ID)
		}
	})
}

func TestUseCaseLogout(t *testing.T) {
	t.Run("セッションを削除してログアウトできる", func(t *testing.T) {
		users := newFakeUserRepository()
		sessions := &fakeSessionRepository{}
		usecase := NewUseCase(
			fakeGitHubOAuthClient{
				authURL: "https://github.com/login/oauth/authorize",
				profile: user.GitHubProfile{
					GitHubID: 1,
					Username: "octocat",
				},
			},
			users,
			sessions,
			fakeCurrentUserRepository{users: users, sessions: sessions},
			&fakeTokenGenerator{tokens: []string{"session-token"}},
		)
		usecase.now = func() time.Time {
			return time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
		}

		if _, err := usecase.CompleteGitHubLogin(context.Background(), "code"); err != nil {
			t.Fatalf("CompleteGitHubLogin がエラーを返しました: %v", err)
		}

		if err := usecase.Logout(context.Background(), "session-token"); err != nil {
			t.Fatalf("Logout がエラーを返しました: %v", err)
		}

		if _, err := usecase.CurrentUser(context.Background(), "session-token"); err != ErrUnauthenticated {
			t.Fatalf("CurrentUser error = %v, 期待値 ErrUnauthenticated", err)
		}
	})
}
