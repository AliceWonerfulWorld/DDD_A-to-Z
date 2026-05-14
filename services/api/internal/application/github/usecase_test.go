package github

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeCurrentUserRepository struct {
	appUser user.User
	ok      bool
}

func (r fakeCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	if err := ctx.Err(); err != nil {
		return user.User{}, false, err
	}
	return r.appUser, r.ok, nil
}

type fakeTokenRepository struct {
	token string
	ok    bool
	err   error
}

func (r fakeTokenRepository) GitHubAccessToken(ctx context.Context, userID user.ID) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	return r.token, r.ok, r.err
}

type fakeRepositoryClient struct {
	repositories []repositoryanalysis.Repository
	err          error
}

func (c fakeRepositoryClient) ListRepositories(ctx context.Context, accessToken string, userID user.ID, syncedAt time.Time) ([]repositoryanalysis.Repository, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.repositories, nil
}

type fakeRepositoryStore struct {
	upserted []repositoryanalysis.Repository
	listed   []repositoryanalysis.Repository
}

func (s *fakeRepositoryStore) UpsertRepositories(ctx context.Context, repositories []repositoryanalysis.Repository) error {
	s.upserted = repositories
	return nil
}

func (s *fakeRepositoryStore) ListRepositories(ctx context.Context, userID user.ID) ([]repositoryanalysis.Repository, error) {
	return s.listed, nil
}

func TestUseCaseSyncRepositories(t *testing.T) {
	appUser := user.User{ID: "user_1"}
	syncedAt := time.Date(2026, 5, 14, 16, 0, 0, 0, time.UTC)
	repository := repositoryanalysis.Repository{
		GitHubID:        1,
		UserID:          appUser.ID,
		Owner:           "octocat",
		Name:            "hello-world",
		FullName:        "octocat/hello-world",
		DefaultBranch:   "main",
		HTMLURL:         "https://github.com/octocat/hello-world",
		GitHubUpdatedAt: syncedAt,
		SyncedAt:        syncedAt,
	}
	store := &fakeRepositoryStore{}
	usecase := NewUseCase(
		fakeCurrentUserRepository{appUser: appUser, ok: true},
		fakeTokenRepository{token: "github-token", ok: true},
		fakeRepositoryClient{repositories: []repositoryanalysis.Repository{repository}},
		store,
	)
	usecase.now = func() time.Time { return syncedAt }

	result, err := usecase.SyncRepositories(context.Background(), "session-token")
	if err != nil {
		t.Fatalf("SyncRepositories() がエラーを返しました: %v", err)
	}

	if result.SyncedCount != 1 {
		t.Fatalf("SyncedCount = %d, 期待値 1", result.SyncedCount)
	}
	if len(store.upserted) != 1 || store.upserted[0].FullName != "octocat/hello-world" {
		t.Fatalf("upserted repositories = %#v", store.upserted)
	}
}

func TestUseCaseSyncRepositoriesErrors(t *testing.T) {
	t.Run("未認証を返す", func(t *testing.T) {
		usecase := NewUseCase(fakeCurrentUserRepository{}, fakeTokenRepository{}, fakeRepositoryClient{}, &fakeRepositoryStore{})
		if _, err := usecase.SyncRepositories(context.Background(), ""); !errors.Is(err, ErrUnauthenticated) {
			t.Fatalf("error = %v, 期待値 ErrUnauthenticated", err)
		}
	})

	t.Run("GitHub token がない場合はエラーを返す", func(t *testing.T) {
		usecase := NewUseCase(
			fakeCurrentUserRepository{appUser: user.User{ID: "user_1"}, ok: true},
			fakeTokenRepository{},
			fakeRepositoryClient{},
			&fakeRepositoryStore{},
		)
		if _, err := usecase.SyncRepositories(context.Background(), "session-token"); !errors.Is(err, ErrMissingGitHubToken) {
			t.Fatalf("error = %v, 期待値 ErrMissingGitHubToken", err)
		}
	})
}
