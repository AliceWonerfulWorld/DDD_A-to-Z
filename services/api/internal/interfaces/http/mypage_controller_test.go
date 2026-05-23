package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"
	"os"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
)

// --- stubs ---

type stubMypageCurrentUser struct {
	user  user.User
	found bool
	err   error
}

func (s *stubMypageCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return s.user, s.found, s.err
}

type stubMypageCPReader struct {
	balance     int64
	totalEarned int64
	totalSpent  int64
}

func (s *stubMypageCPReader) GetBalance(_ context.Context, _ user.ID) (int64, error) {
	return s.balance, nil
}

func (s *stubMypageCPReader) GetTotalEarned(_ context.Context, _ user.ID) (int64, error) {
	return s.totalEarned, nil
}

func (s *stubMypageCPReader) GetTotalSpent(_ context.Context, _ user.ID) (int64, error) {
	return s.totalSpent, nil
}

type stubMypageRepoReader struct {
	summary mypage.RepositorySummary
}

func (s *stubMypageRepoReader) GetRepositorySummary(_ context.Context, _ user.ID, _ int) (mypage.RepositorySummary, error) {
	return s.summary, nil
}

type stubMypageGitHubStatsReader struct {
	stats *mypage.GitHubStats
	err   error
}

func (s *stubMypageGitHubStatsReader) FetchStats(_ context.Context, _, _ string) (*mypage.GitHubStats, error) {
	return s.stats, s.err
}

type stubMypageGitHubTokenRepo struct {
	token string
	found bool
	err   error
}

func (s *stubMypageGitHubTokenRepo) GitHubAccessToken(_ context.Context, _ user.ID) (string, bool, error) {
	return s.token, s.found, s.err
}

type stubMypageGuildMembershipReader struct {
	guild *mypage.GuildInfo
	total int
}

func (s *stubMypageGuildMembershipReader) GetGuildMembership(_ context.Context, _ user.ID) (*mypage.GuildInfo, error) {
	return s.guild, nil
}

func (s *stubMypageGuildMembershipReader) GetTotalGuilds(_ context.Context) (int, error) {
	return s.total, nil
}

type stubMypageBadgeReader struct {
	badges []mypage.BadgeSummary
}

func (s *stubMypageBadgeReader) ListUserBadges(_ context.Context, _ user.ID) ([]mypage.BadgeSummary, error) {
	return s.badges, nil
}

// --- tests ---

func TestMypageController_NoCookie(t *testing.T) {
	uc := mypage.NewUseCase(
		&stubMypageCurrentUser{},
		&stubMypageCPReader{},
		&stubMypageRepoReader{},
		&stubMypageGitHubStatsReader{},
		&stubMypageGitHubTokenRepo{},
		&stubMypageGuildMembershipReader{},
		&stubMypageBadgeReader{},
		nil,
	)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	controller := httpapi.NewMypageController(uc, logger)

	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/mypage", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMypageController_Success(t *testing.T) {
	testUser := user.User{
		ID: "github_99",
		GitHubAccount: user.GitHubAccount{
			GitHubID:  99,
			Username:  "octocat",
			AvatarURL: "https://example.com/avatar.png",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	uc := mypage.NewUseCase(
		&stubMypageCurrentUser{user: testUser, found: true},
		&stubMypageCPReader{balance: 50, totalEarned: 100, totalSpent: 50},
		&stubMypageRepoReader{summary: mypage.RepositorySummary{
			TotalCount:      2,
			LanguageSummary: map[string]int{"Go": 2},
			Recent: []mypage.RecentRepository{
				{GitHubID: 1, FullName: "octocat/repo", Language: "Go", HTMLURL: "https://github.com/octocat/repo"},
			},
		}},
		&stubMypageGitHubStatsReader{},
		&stubMypageGitHubTokenRepo{},
		&stubMypageGuildMembershipReader{},
		&stubMypageBadgeReader{},
		nil,
	)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	controller := httpapi.NewMypageController(uc, logger)

	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/mypage", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	userSection, ok := body["user"].(map[string]any)
	if !ok {
		t.Fatal("expected user section in response")
	}
	if userSection["username"] != "octocat" {
		t.Errorf("expected username octocat, got %v", userSection["username"])
	}

	cpSection, ok := body["contribution_points"].(map[string]any)
	if !ok {
		t.Fatal("expected contribution_points section in response")
	}
	if cpSection["balance"].(float64) != 50 {
		t.Errorf("expected balance 50, got %v", cpSection["balance"])
	}

	if body["guild"] != nil {
		t.Errorf("expected null guild, got %v", body["guild"])
	}
}
