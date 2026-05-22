package chat

import (
	"context"
	"errors"
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// --- test doubles ---

type testCurrentUser struct {
	user user.User
	ok   bool
	err  error
}

func (r testCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return r.user, r.ok, r.err
}

type testRepo struct {
	membership    *guilddomain.MembershipWithGuild
	membershipErr error
	inserted      *ChatToken
	insertErr     error
}

func (r testRepo) FindMembershipByUserAndGuild(_ context.Context, userID user.ID, guildID guilddomain.ID) (guilddomain.MembershipWithGuild, bool, error) {
	if r.membershipErr != nil {
		return guilddomain.MembershipWithGuild{}, false, r.membershipErr
	}
	if r.membership == nil || r.membership.Membership.UserID != userID || r.membership.Guild.ID != guildID {
		return guilddomain.MembershipWithGuild{}, false, nil
	}
	return *r.membership, true, nil
}

func (r *testRepo) InsertChatToken(_ context.Context, token ChatToken) error {
	if r.insertErr != nil {
		return r.insertErr
	}
	r.inserted = &token
	return nil
}

type testTokenGen struct{ token string }

func (g testTokenGen) NewToken() (string, error) { return g.token, nil }

type testHasher struct{}

func (h testHasher) Hash(token string) string { return "hash:" + token }

// --- helpers ---

func fixedUser(id user.ID) user.User {
	return user.User{ID: id}
}

func membership(userID user.ID, guildID guilddomain.ID) guilddomain.MembershipWithGuild {
	return guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{UserID: userID},
		Guild:      guilddomain.Guild{ID: guildID},
	}
}

func newUseCase(current testCurrentUser, repo *testRepo) *UseCase {
	uc := NewUseCase(current, repo, testTokenGen{token: "raw-token"}, testHasher{})
	uc.now = func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }
	return uc
}

// --- tests ---

func TestIssueGuildChatToken_Success(t *testing.T) {
	uid := user.ID("user_001")
	gid := guilddomain.ID("guild_rust")
	repo := &testRepo{membership: ptr(membership(uid, gid))}
	uc := newUseCase(testCurrentUser{user: fixedUser(uid), ok: true}, repo)

	got, err := uc.IssueGuildChatToken(context.Background(), "session", gid)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Token != "raw-token" {
		t.Errorf("Token = %q, want %q", got.Token, "raw-token")
	}
	if got.TokenHash != "hash:raw-token" {
		t.Errorf("TokenHash = %q, want %q", got.TokenHash, "hash:raw-token")
	}
	if got.UserID != string(uid) {
		t.Errorf("UserID = %q, want %q", got.UserID, uid)
	}
	if got.GuildID != string(gid) {
		t.Errorf("GuildID = %q, want %q", got.GuildID, gid)
	}
	if got.ExpiresAt != uc.now().Add(chatTokenTTL) {
		t.Errorf("ExpiresAt = %v, want %v", got.ExpiresAt, uc.now().Add(chatTokenTTL))
	}
	if repo.inserted == nil {
		t.Fatal("InsertChatToken was not called")
	}
}

func TestIssueGuildChatToken_EmptySessionToken(t *testing.T) {
	repo := &testRepo{}
	uc := newUseCase(testCurrentUser{ok: true}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "", "guild_rust")

	if !errors.Is(err, ErrUnauthenticated) {
		t.Errorf("err = %v, want ErrUnauthenticated", err)
	}
}

func TestIssueGuildChatToken_EmptyGuildID(t *testing.T) {
	uid := user.ID("user_001")
	repo := &testRepo{}
	uc := newUseCase(testCurrentUser{user: fixedUser(uid), ok: true}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "session", "")

	if !errors.Is(err, ErrForbidden) {
		t.Errorf("err = %v, want ErrForbidden", err)
	}
}

func TestIssueGuildChatToken_SessionNotFound(t *testing.T) {
	repo := &testRepo{}
	uc := newUseCase(testCurrentUser{ok: false}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "session", "guild_rust")

	if !errors.Is(err, ErrUnauthenticated) {
		t.Errorf("err = %v, want ErrUnauthenticated", err)
	}
}

func TestIssueGuildChatToken_NotMember(t *testing.T) {
	uid := user.ID("user_001")
	// メンバーシップが存在しない
	repo := &testRepo{membership: nil}
	uc := newUseCase(testCurrentUser{user: fixedUser(uid), ok: true}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "session", "guild_rust")

	if !errors.Is(err, ErrForbidden) {
		t.Errorf("err = %v, want ErrForbidden", err)
	}
}

func TestIssueGuildChatToken_DifferentGuild(t *testing.T) {
	uid := user.ID("user_001")
	// guild_python に所属しているユーザーが guild_rust のトークンを要求
	repo := &testRepo{membership: ptr(membership(uid, "guild_python"))}
	uc := newUseCase(testCurrentUser{user: fixedUser(uid), ok: true}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "session", "guild_rust")

	if !errors.Is(err, ErrForbidden) {
		t.Errorf("err = %v, want ErrForbidden", err)
	}
}

func TestIssueGuildChatToken_InsertError(t *testing.T) {
	uid := user.ID("user_001")
	gid := guilddomain.ID("guild_rust")
	dbErr := errors.New("db error")
	repo := &testRepo{membership: ptr(membership(uid, gid)), insertErr: dbErr}
	uc := newUseCase(testCurrentUser{user: fixedUser(uid), ok: true}, repo)

	_, err := uc.IssueGuildChatToken(context.Background(), "session", gid)

	if !errors.Is(err, dbErr) {
		t.Errorf("err = %v, want %v", err, dbErr)
	}
}

func ptr[T any](v T) *T { return &v }
