package github

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
)

type UseCase struct {
	current CurrentUserRepository
	tokens  TokenRepository
	client  RepositoryClient
	store   RepositoryStore
	now     func() time.Time
}

func NewUseCase(
	current CurrentUserRepository,
	tokens TokenRepository,
	client RepositoryClient,
	store RepositoryStore,
) *UseCase {
	return &UseCase{
		current: current,
		tokens:  tokens,
		client:  client,
		store:   store,
		now:     time.Now,
	}
}

type SyncResult struct {
	Repositories []repositoryanalysis.Repository
	SyncedCount  int
}

func (u *UseCase) SyncRepositories(ctx context.Context, sessionToken string) (SyncResult, error) {
	if sessionToken == "" {
		return SyncResult{}, ErrUnauthenticated
	}

	now := u.now()
	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, now)
	if err != nil {
		return SyncResult{}, err
	}
	if !ok {
		return SyncResult{}, ErrUnauthenticated
	}

	accessToken, ok, err := u.tokens.GitHubAccessToken(ctx, appUser.ID)
	if err != nil {
		return SyncResult{}, err
	}
	if !ok || accessToken == "" {
		return SyncResult{}, ErrMissingGitHubToken
	}

	repositories, err := u.client.ListRepositories(ctx, accessToken, appUser.ID, now)
	if err != nil {
		return SyncResult{}, err
	}

	if err := u.store.UpsertRepositories(ctx, repositories); err != nil {
		return SyncResult{}, err
	}

	return SyncResult{
		Repositories: repositories,
		SyncedCount:  len(repositories),
	}, nil
}

func (u *UseCase) ListRepositories(ctx context.Context, sessionToken string) ([]repositoryanalysis.Repository, error) {
	if sessionToken == "" {
		return nil, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUnauthenticated
	}

	return u.store.ListRepositories(ctx, appUser.ID)
}
