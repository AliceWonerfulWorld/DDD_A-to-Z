package repositoryanalysis

import (
	"context"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type fakeCurrentUserRepo struct {
	u user.User
}

func (r *fakeCurrentUserRepo) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return r.u, true, nil
}

type fakeTokenRepo struct {
	token string
}

func (r *fakeTokenRepo) GitHubAccessToken(_ context.Context, _ user.ID) (string, bool, error) {
	return r.token, true, nil
}

type fakeRepoClient struct {
	repos     []repositoryanalysis.Repository
	languages map[string]int64
}

func (r *fakeRepoClient) ListRepositories(_ context.Context, _ string, _ user.ID, _ time.Time) ([]repositoryanalysis.Repository, error) {
	return r.repos, nil
}

func (r *fakeRepoClient) ListCommits(_ context.Context, _, _, _, _ string, _ time.Time) ([]repositoryanalysis.CommitItem, error) {
	commits := make([]repositoryanalysis.CommitItem, 10)
	for i := range commits {
		commits[i] = repositoryanalysis.CommitItem{Message: "fix", Committed: time.Now()}
	}
	return commits, nil
}

func (r *fakeRepoClient) ListPullRequests(_ context.Context, _, _, _, _ string, _ time.Time) ([]repositoryanalysis.PullRequestItem, error) {
	return nil, nil
}

func (r *fakeRepoClient) ListLanguages(_ context.Context, _, _, _ string) (map[string]int64, error) {
	return r.languages, nil
}

type fakeRepoStore struct{}

func (s *fakeRepoStore) UpsertRepositories(_ context.Context, _ []repositoryanalysis.Repository) error {
	return nil
}

type fakeCPEarner struct {
	earned int64
}

func (e *fakeCPEarner) Earn(_ context.Context, _ user.ID, amount int64, _, _, _ string) error {
	e.earned += amount
	return nil
}

type fakeSPEarner struct {
	earned map[string]int64
}

func (e *fakeSPEarner) EarnSP(_ context.Context, _ user.ID, language string, amount int64, _, _, _ string) error {
	if e.earned == nil {
		e.earned = map[string]int64{}
	}
	e.earned[language] += amount
	return nil
}

type fakeCPBalanceProvider struct{}

func (p *fakeCPBalanceProvider) GetBalance(_ context.Context, _ user.ID) (int64, error) {
	return 0, nil
}

func (p *fakeCPBalanceProvider) GetLastAnalyzedAt(_ context.Context, _ user.ID) (*time.Time, error) {
	return nil, nil
}

func (p *fakeCPBalanceProvider) UpdateLastAnalyzedAt(_ context.Context, _ user.ID, _ time.Time) error {
	return nil
}

func newTestUseCase(repos []repositoryanalysis.Repository, langs map[string]int64, cp *fakeCPEarner, sp *fakeSPEarner) *UseCase {
	repoClient := &fakeRepoClient{repos: repos, languages: langs}
	appUser := user.User{
		ID:            "user_1",
		GitHubAccount: user.GitHubAccount{Username: "testuser"},
	}
	uc := NewUseCase(
		&fakeCurrentUserRepo{u: appUser},
		&fakeTokenRepo{token: "token"},
		repoClient,
		&fakeRepoStore{},
		repoClient,
		repoClient,
		repoClient,
		cp,
		sp,
		&fakeCPBalanceProvider{},
	)
	uc.now = func() time.Time { return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) }
	return uc
}

func testRepo(language string) repositoryanalysis.Repository {
	return repositoryanalysis.Repository{
		GitHubID:        1,
		FullName:        "owner/" + language,
		Owner:           "owner",
		Name:            language,
		Language:        language,
		DefaultBranch:   "main",
		HTMLURL:         "https://github.com/owner/" + language,
		GitHubUpdatedAt: time.Now(),
		SyncedAt:        time.Now(),
	}
}

func TestAnalyzeSPEarned(t *testing.T) {
	t.Run("GoリポジトリのコミットでGo SPが付与される", func(t *testing.T) {
		cp := &fakeCPEarner{}
		sp := &fakeSPEarner{}
		repos := []repositoryanalysis.Repository{testRepo("Go")}
		uc := newTestUseCase(repos, map[string]int64{"Go": 1000}, cp, sp)

		_, err := uc.Analyze(context.Background(), "session")
		if err != nil {
			t.Fatalf("Analyze() がエラーを返しました: %v", err)
		}
		if sp.earned["Go"] == 0 {
			t.Error("Go SP が付与されていません")
		}
	})

	t.Run("TypeScriptリポジトリのコミットでTypeScript SPが付与される", func(t *testing.T) {
		cp := &fakeCPEarner{}
		sp := &fakeSPEarner{}
		repos := []repositoryanalysis.Repository{testRepo("TypeScript")}
		uc := newTestUseCase(repos, map[string]int64{"TypeScript": 1000}, cp, sp)

		_, err := uc.Analyze(context.Background(), "session")
		if err != nil {
			t.Fatalf("Analyze() がエラーを返しました: %v", err)
		}
		if sp.earned["TypeScript"] == 0 {
			t.Error("TypeScript SP が付与されていません")
		}
	})

	t.Run("複数言語のリポジトリで言語別にSPが付与される", func(t *testing.T) {
		cp := &fakeCPEarner{}
		sp := &fakeSPEarner{}
		repos := []repositoryanalysis.Repository{testRepo("Go")}
		uc := newTestUseCase(repos, map[string]int64{"Go": 600, "TypeScript": 400}, cp, sp)

		_, err := uc.Analyze(context.Background(), "session")
		if err != nil {
			t.Fatalf("Analyze() がエラーを返しました: %v", err)
		}
		if sp.earned["Go"] == 0 {
			t.Error("Go SP が付与されていません")
		}
		if sp.earned["TypeScript"] == 0 {
			t.Error("TypeScript SP が付与されていません")
		}
		totalSP := sp.earned["Go"] + sp.earned["TypeScript"]
		if totalSP != cp.earned {
			t.Errorf("SP合計 %d が CP %d と一致しません", totalSP, cp.earned)
		}
	})

	t.Run("SP付与量の合計はCP付与量と等しい", func(t *testing.T) {
		cp := &fakeCPEarner{}
		sp := &fakeSPEarner{}
		repos := []repositoryanalysis.Repository{testRepo("Go")}
		uc := newTestUseCase(repos, map[string]int64{"Go": 1}, cp, sp)

		_, err := uc.Analyze(context.Background(), "session")
		if err != nil {
			t.Fatalf("Analyze() がエラーを返しました: %v", err)
		}
		var totalSP int64
		for _, v := range sp.earned {
			totalSP += v
		}
		if totalSP != cp.earned {
			t.Errorf("SP合計 %d が CP %d と一致しません", totalSP, cp.earned)
		}
	})
}
