package http

import (
	"context"
	"encoding/json"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	chatapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/chat"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// --- test doubles ---

type chatTestCurrentUser struct {
	appUser user.User
	ok      bool
}

func (r chatTestCurrentUser) FindUserBySessionToken(_ context.Context, _ string, _ time.Time) (user.User, bool, error) {
	return r.appUser, r.ok, nil
}

type chatTestRepo struct {
	membership *guilddomain.MembershipWithGuild
	insertErr  error
}

func (r chatTestRepo) FindMembershipByUserAndGuild(_ context.Context, userID user.ID, guildID guilddomain.ID) (guilddomain.MembershipWithGuild, bool, error) {
	if r.membership == nil || r.membership.Membership.UserID != userID || r.membership.Guild.ID != guildID {
		return guilddomain.MembershipWithGuild{}, false, nil
	}
	return *r.membership, true, nil
}

func (r *chatTestRepo) InsertChatToken(_ context.Context, _ chatapp.ChatToken) error {
	return r.insertErr
}

type chatTestTokenGen struct{}

func (g chatTestTokenGen) NewToken() (string, error) { return "test-raw-token", nil }

type chatTestHasher struct{}

func (h chatTestHasher) Hash(token string) string { return "hash:" + token }

func newChatController(current chatTestCurrentUser, repo *chatTestRepo) *ChatController {
	uc := chatapp.NewUseCase(current, repo, chatTestTokenGen{}, chatTestHasher{})
	return NewChatController(uc, slog.Default())
}

func chatRequest(method, path, cookie string) *stdhttp.Request {
	r := httptest.NewRequest(method, path, nil)
	if cookie != "" {
		r.Header.Set("Cookie", sessionCookieName+"="+cookie)
	}
	return r
}

// --- tests ---

func TestChatController_IssueChatToken_Success(t *testing.T) {
	uid := user.ID("user_001")
	gid := guilddomain.ID("guild_rust")
	membership := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{UserID: uid},
		Guild:      guilddomain.Guild{ID: gid},
	}
	repo := &chatTestRepo{membership: &membership}
	c := newChatController(chatTestCurrentUser{appUser: user.User{ID: uid}, ok: true}, repo)

	mux := stdhttp.NewServeMux()
	c.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := chatRequest("POST", "/guilds/guild_rust/chat-token", "session")
	r.SetPathValue("guildID", "guild_rust")
	mux.ServeHTTP(w, r)

	if w.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, stdhttp.StatusOK, w.Body.String())
	}

	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["token"] != "test-raw-token" {
		t.Errorf("token = %v, want %q", body["token"], "test-raw-token")
	}
	if body["expires_at"] == "" || body["expires_at"] == nil {
		t.Errorf("expires_at is empty")
	}
}

func TestChatController_IssueChatToken_NoCookie(t *testing.T) {
	repo := &chatTestRepo{}
	c := newChatController(chatTestCurrentUser{}, repo)

	mux := stdhttp.NewServeMux()
	c.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := chatRequest("POST", "/guilds/guild_rust/chat-token", "")
	r.SetPathValue("guildID", "guild_rust")
	mux.ServeHTTP(w, r)

	if w.Code != stdhttp.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, stdhttp.StatusUnauthorized)
	}
}

func TestChatController_IssueChatToken_SessionNotFound(t *testing.T) {
	repo := &chatTestRepo{}
	c := newChatController(chatTestCurrentUser{ok: false}, repo)

	mux := stdhttp.NewServeMux()
	c.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := chatRequest("POST", "/guilds/guild_rust/chat-token", "invalid-session")
	r.SetPathValue("guildID", "guild_rust")
	mux.ServeHTTP(w, r)

	if w.Code != stdhttp.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, stdhttp.StatusUnauthorized)
	}
}

func TestChatController_IssueChatToken_NotMember(t *testing.T) {
	uid := user.ID("user_001")
	// メンバーシップなし
	repo := &chatTestRepo{membership: nil}
	c := newChatController(chatTestCurrentUser{appUser: user.User{ID: uid}, ok: true}, repo)

	mux := stdhttp.NewServeMux()
	c.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := chatRequest("POST", "/guilds/guild_rust/chat-token", "session")
	r.SetPathValue("guildID", "guild_rust")
	mux.ServeHTTP(w, r)

	if w.Code != stdhttp.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, stdhttp.StatusForbidden)
	}
}

func TestChatController_IssueChatToken_DifferentGuild(t *testing.T) {
	uid := user.ID("user_001")
	// guild_python 所属ユーザーが guild_rust のトークンを要求
	membership := guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{UserID: uid},
		Guild:      guilddomain.Guild{ID: "guild_python"},
	}
	repo := &chatTestRepo{membership: &membership}
	c := newChatController(chatTestCurrentUser{appUser: user.User{ID: uid}, ok: true}, repo)

	mux := stdhttp.NewServeMux()
	c.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := chatRequest("POST", "/guilds/guild_rust/chat-token", "session")
	r.SetPathValue("guildID", "guild_rust")
	mux.ServeHTTP(w, r)

	if w.Code != stdhttp.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, stdhttp.StatusForbidden)
	}
}
