package sp_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/sp"
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

type stubSPReader struct {
	balances []sp.SPBalance
	err      error
}

func (s *stubSPReader) GetSPBalances(_ context.Context, _ user.ID) ([]sp.SPBalance, error) {
	return s.balances, s.err
}

// --- tests ---

func TestGetSP_EmptyToken(t *testing.T) {
	uc := sp.NewUseCase(&stubCurrentUser{}, &stubSPReader{})

	_, err := uc.GetSP(context.Background(), "")
	if !errors.Is(err, sp.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetSP_SessionNotFound(t *testing.T) {
	uc := sp.NewUseCase(&stubCurrentUser{found: false}, &stubSPReader{})

	_, err := uc.GetSP(context.Background(), "invalid-token")
	if !errors.Is(err, sp.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestGetSP_SessionError(t *testing.T) {
	sessionErr := errors.New("db error")
	uc := sp.NewUseCase(&stubCurrentUser{err: sessionErr}, &stubSPReader{})

	_, err := uc.GetSP(context.Background(), "token")
	if !errors.Is(err, sessionErr) {
		t.Fatalf("expected sessionErr, got %v", err)
	}
}

func TestGetSP_Success(t *testing.T) {
	testUser := user.User{ID: "github_42", GitHubAccount: user.GitHubAccount{GitHubID: 42, Username: "testuser", AvatarURL: "https://example.com/avatar.png"}}
	balances := []sp.SPBalance{
		{Language: "Go", Balance: 150},
		{Language: "TypeScript", Balance: 80},
	}
	uc := sp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubSPReader{balances: balances},
	)

	result, err := uc.GetSP(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 balances, got %d", len(result))
	}
	if result[0].Language != "Go" || result[0].Balance != 150 {
		t.Errorf("unexpected first balance: %+v", result[0])
	}
	if result[1].Language != "TypeScript" || result[1].Balance != 80 {
		t.Errorf("unexpected second balance: %+v", result[1])
	}
}

func TestGetSP_NoSP(t *testing.T) {
	testUser := user.User{ID: "github_42", GitHubAccount: user.GitHubAccount{GitHubID: 42, Username: "testuser", AvatarURL: "https://example.com/avatar.png"}}
	uc := sp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubSPReader{balances: []sp.SPBalance{}},
	)

	result, err := uc.GetSP(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty balances, got %d", len(result))
	}
}

func TestGetSP_SPReaderError(t *testing.T) {
	testUser := user.User{ID: "github_42", GitHubAccount: user.GitHubAccount{GitHubID: 42, Username: "testuser", AvatarURL: "https://example.com/avatar.png"}}
	spErr := errors.New("sp query failed")
	uc := sp.NewUseCase(
		&stubCurrentUser{user: testUser, found: true},
		&stubSPReader{err: spErr},
	)

	_, err := uc.GetSP(context.Background(), "valid-token")
	if !errors.Is(err, spErr) {
		t.Fatalf("expected spErr, got %v", err)
	}
}
