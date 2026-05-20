package mypage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// --- stubs ---

type stubCurrentUser struct {
	user  user.User
	found bool
	err   error
}

func (s *stubCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return s.user, s.found, s.err
}

type stubCPReader struct {
	balance     int64
	totalEarned int64
	totalSpent  int64
	err         error
}

func (s *stubCPReader) GetBalance(_ context.Context, _ user.ID) (int64, error) {
	return s.balance, s.err
}

func (s *stubCPReader) GetTotalEarned(_ context.Context, _ user.ID) (int64, error) {
	return s.totalEarned, s.err
}

func (s *stubCPReader) GetTotalSpent(_ context.Context, _ user.ID) (int64, error) {
	return s.totalSpent, s.err
}

type stubRepoReader struct {
	summary mypage.RepositorySummary
	err     error
}

func (s *stubRepoReader) GetRepositorySummary(_ context.Context, _ user.ID, _ int) (mypage.RepositorySummary, error) {
	return s.summary, s.err
}

type stubGitHubStatsReader struct {
	stats *mypage.GitHubStats
	err   error
}

func (s *stubGitHubStatsReader) FetchStats(_ context.Context, _, _ string) (*mypage.GitHubStats, error) {
	return s.stats, s.err
}

type stubGitHubTokenRepo struct {
	token string
	found bool
	err   error
}

func (s *stubGitHubTokenRepo) GitHubAccessToken(_ context.Context, _ user.ID) (string, bool, error) {
	return s.token, s.found, s.err
}

type stubGuildMembershipReader struct {
	guild *mypage.GuildInfo
	total int
	err   error
}

func (s *stubGuildMembershipReader) GetGuildMembership(_ context.Context, _ user.ID) (*mypage.GuildInfo, error) {
	return s.guild, s.err
}

func (s *stubGuildMembershipReader) GetTotalGuilds(_ context.Context) (int, error) {
	return s.total, s.err
}

// --- tests ---

func TestGetMyPage_EmptyToken(t *testing.T) {
	uc := mypage.NewUseCase(
		&stubCurrentUser{},
		&stubCPReader{},
		&stubRepoReader{},
		&stubGitHubStatsReader{},
		&stubGitHubTokenRepo{},
		&stubGuildMembershipReader{},
	)

	_, err := uc.GetMyPage(context.Background(), "")
	if !errors.Is(err, mypage.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetMyPage_SessionNotFound(t *testing.T) {
	uc := mypage.NewUseCase(
		&stubCurrentUser{found: false},
		&stubCPReader{},
		&stubRepoReader{},
		&stubGitHubStatsReader{},
		&stubGitHubTokenRepo{},
		&stubGuildMembershipReader{},
	)

	_, err := uc.GetMyPage(context.Background(), "invalid-token")
	if !errors.Is(err, mypage.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetMyPage_Success(t *testing.T) {
	testUser := user.User{
		ID: "github_42",
		GitHubAccount: user.GitHubAccount{
			GitHubID:  42,
			Username:  "testuser",
			AvatarURL: "https://example.com/avatar.png",
		},
	}

	uc := mypage.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPReader{balance: 100, totalEarned: 200, totalSpent: 100},
		&stubRepoReader{summary: mypage.RepositorySummary{
			TotalCount:      3,
			LanguageSummary: map[string]int{"Go": 2, "TypeScript": 1},
			Recent: []mypage.RecentRepository{
				{GitHubID: 1, FullName: "testuser/repo1", Language: "Go", HTMLURL: "https://github.com/testuser/repo1"},
			},
		}},
		&stubGitHubStatsReader{},
		&stubGitHubTokenRepo{},
		&stubGuildMembershipReader{},
	)

	result, err := uc.GetMyPage(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.User.ID != "github_42" {
		t.Errorf("expected user id github_42, got %s", result.User.ID)
	}
	if result.CP.Balance != 100 {
		t.Errorf("expected balance 100, got %d", result.CP.Balance)
	}
	if result.CP.TotalEarned != 200 {
		t.Errorf("expected total earned 200, got %d", result.CP.TotalEarned)
	}
	if result.CP.TotalSpent != 100 {
		t.Errorf("expected total spent 100, got %d", result.CP.TotalSpent)
	}
	if result.Repositories.TotalCount != 3 {
		t.Errorf("expected 3 repositories, got %d", result.Repositories.TotalCount)
	}
	if result.Guild != nil {
		t.Errorf("expected nil guild, got %v", result.Guild)
	}
}

func TestGetMyPage_CPError(t *testing.T) {
	testUser := user.User{ID: "github_42", GitHubAccount: user.GitHubAccount{GitHubID: 42, Username: "test", AvatarURL: "https://example.com"}}
	cpErr := errors.New("db connection failed")

	uc := mypage.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPReader{err: cpErr},
		&stubRepoReader{},
		&stubGitHubStatsReader{},
		&stubGitHubTokenRepo{},
		&stubGuildMembershipReader{},
	)

	_, err := uc.GetMyPage(context.Background(), "valid-token")
	if !errors.Is(err, cpErr) {
		t.Fatalf("expected cpErr, got %v", err)
	}
}

func TestGetMyPage_RepoError(t *testing.T) {
	testUser := user.User{ID: "github_42", GitHubAccount: user.GitHubAccount{GitHubID: 42, Username: "test", AvatarURL: "https://example.com"}}
	repoErr := errors.New("repo query failed")

	uc := mypage.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubCPReader{balance: 0, totalEarned: 0, totalSpent: 0},
		&stubRepoReader{err: repoErr},
		&stubGitHubStatsReader{},
		&stubGitHubTokenRepo{},
		&stubGuildMembershipReader{},
	)

	_, err := uc.GetMyPage(context.Background(), "valid-token")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got %v", err)
	}
}
