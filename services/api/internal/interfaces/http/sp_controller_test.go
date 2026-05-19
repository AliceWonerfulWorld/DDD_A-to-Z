package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/sp"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	httpapi "github.com/jyogi-web/ddd-a-to-z/services/api/internal/interfaces/http"
)

// --- stubs ---

type stubSPCurrentUser struct {
	u     user.User
	found bool
	err   error
}

func (s *stubSPCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return s.u, s.found, s.err
}

type stubSPReader struct {
	balances []sp.SPBalance
	err      error
}

func (s *stubSPReader) GetSPBalances(_ context.Context, _ user.ID) ([]sp.SPBalance, error) {
	return s.balances, s.err
}

func newSPMux(currentUser *stubSPCurrentUser, reader *stubSPReader) *http.ServeMux {
	uc := sp.NewUseCase(currentUser, reader)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	controller := httpapi.NewSPController(uc, logger)
	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)
	return mux
}

// --- tests ---

func TestSPController_NoCookie(t *testing.T) {
	mux := newSPMux(&stubSPCurrentUser{}, &stubSPReader{})

	req := httptest.NewRequest("GET", "/me/sp", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	errSection, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatal("expected error section in response")
	}
	if errSection["code"] != "unauthenticated" {
		t.Errorf("expected code unauthenticated, got %v", errSection["code"])
	}
}

func TestSPController_InvalidSession(t *testing.T) {
	mux := newSPMux(&stubSPCurrentUser{found: false}, &stubSPReader{})

	req := httptest.NewRequest("GET", "/me/sp", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "invalid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestSPController_Success(t *testing.T) {
	testUser := user.User{ID: "github_99", GitHubAccount: user.GitHubAccount{GitHubID: 99, Username: "octocat", AvatarURL: "https://example.com/avatar.png"}}
	mux := newSPMux(
		&stubSPCurrentUser{u: testUser, found: true},
		&stubSPReader{balances: []sp.SPBalance{
			{Language: "Go", Balance: 150},
			{Language: "TypeScript", Balance: 80},
		}},
	)

	req := httptest.NewRequest("GET", "/me/sp", nil)
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

	items, ok := body["skill_points"].([]any)
	if !ok {
		t.Fatal("expected skill_points array in response")
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 skill_points, got %d", len(items))
	}

	first := items[0].(map[string]any)
	if first["language"] != "Go" {
		t.Errorf("expected language Go, got %v", first["language"])
	}
	if first["balance"].(float64) != 150 {
		t.Errorf("expected balance 150, got %v", first["balance"])
	}
}

func TestSPController_EmptySP(t *testing.T) {
	testUser := user.User{ID: "github_99", GitHubAccount: user.GitHubAccount{GitHubID: 99, Username: "octocat", AvatarURL: "https://example.com/avatar.png"}}
	mux := newSPMux(
		&stubSPCurrentUser{u: testUser, found: true},
		&stubSPReader{balances: []sp.SPBalance{}},
	)

	req := httptest.NewRequest("GET", "/me/sp", nil)
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

	items, ok := body["skill_points"].([]any)
	if !ok {
		t.Fatal("expected skill_points array in response")
	}
	if len(items) != 0 {
		t.Errorf("expected empty skill_points, got %d items", len(items))
	}
}

func TestSPController_InternalError(t *testing.T) {
	testUser := user.User{ID: "github_99", GitHubAccount: user.GitHubAccount{GitHubID: 99, Username: "octocat", AvatarURL: "https://example.com/avatar.png"}}
	mux := newSPMux(
		&stubSPCurrentUser{u: testUser, found: true},
		&stubSPReader{err: errors.New("db error")},
	)

	req := httptest.NewRequest("GET", "/me/sp", nil)
	req.AddCookie(&http.Cookie{Name: "lang_war_session", Value: "valid-token"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
