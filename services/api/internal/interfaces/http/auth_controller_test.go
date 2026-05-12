package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/memory"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/security"
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

func TestAuthControllerBeginGitHubLoginRedirects(t *testing.T) {
	controller := newTestAuthController(&fakeTokenGenerator{tokens: []string{"state-token"}})
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/login", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusFound {
		t.Fatalf("status code = %d, want %d", response.Code, stdhttp.StatusFound)
	}

	if got := response.Header().Get("Location"); got != "https://github.test/login?state=state-token" {
		t.Fatalf("Location = %q, want GitHub auth URL", got)
	}
	if cookie := findCookie(response.Result().Cookies(), oauthStateCookieName); cookie == nil {
		t.Fatal("oauth state cookie should be set")
	}
}

func TestAuthControllerCompleteGitHubLoginIssuesSessionCookie(t *testing.T) {
	controller := newTestAuthController(&fakeTokenGenerator{tokens: []string{"state-token", "session-token"}})
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	loginResponse := httptest.NewRecorder()
	loginRequest := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/login", nil)
	router.ServeHTTP(loginResponse, loginRequest)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/callback?code=github-code&state=state-token", nil)
	request.AddCookie(findCookie(loginResponse.Result().Cookies(), oauthStateCookieName))
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, stdhttp.StatusOK)
	}

	cookie := findCookie(response.Result().Cookies(), sessionCookieName)
	if cookie == nil {
		t.Fatal("lang_war_session cookie should be set")
	}
	if cookie.Value != "session-token" {
		t.Fatalf("session cookie = %q, want session-token", cookie.Value)
	}
	if !cookie.HttpOnly {
		t.Fatal("session cookie should be HttpOnly")
	}
	if cookie.SameSite != stdhttp.SameSiteLaxMode {
		t.Fatalf("SameSite = %v, want Lax", cookie.SameSite)
	}

	var body struct {
		User struct {
			ID        string `json:"id"`
			GitHubID  int64  `json:"github_id"`
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.User.ID == "" {
		t.Fatal("response user id should be set")
	}
	if body.User.GitHubID != 123 {
		t.Fatalf("github_id = %d, want 123", body.User.GitHubID)
	}
	if body.User.Username != "octocat" {
		t.Fatalf("username = %q, want octocat", body.User.Username)
	}
	if body.User.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("avatar_url = %q, want test avatar URL", body.User.AvatarURL)
	}
}

func TestAuthControllerCurrentUserReturnsLoggedInUser(t *testing.T) {
	controller := newTestAuthController(&fakeTokenGenerator{tokens: []string{"state-token", "session-token"}})
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	loginResponse := httptest.NewRecorder()
	loginRequest := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/login", nil)
	router.ServeHTTP(loginResponse, loginRequest)

	callbackResponse := httptest.NewRecorder()
	callbackRequest := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/callback?code=github-code&state=state-token", nil)
	callbackRequest.AddCookie(findCookie(loginResponse.Result().Cookies(), oauthStateCookieName))
	router.ServeHTTP(callbackResponse, callbackRequest)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/me", nil)
	request.AddCookie(findCookie(callbackResponse.Result().Cookies(), sessionCookieName))
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.User.ID == "" {
		t.Fatal("response user id should be set")
	}
	if body.User.Username != "octocat" {
		t.Fatalf("username = %q, want octocat", body.User.Username)
	}
}

func TestAuthControllerCompleteGitHubLoginRejectsInvalidState(t *testing.T) {
	controller := newTestAuthController(&fakeTokenGenerator{})
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/callback?code=github-code&state=missing-state", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("status code = %d, want %d", response.Code, stdhttp.StatusBadRequest)
	}
}

func newTestAuthController(tokens *fakeTokenGenerator) *AuthController {
	users := memory.NewUserRepository()
	sessions := memory.NewSessionRepository()
	usecase := authapp.NewUseCase(
		fakeGitHubOAuthClient{
			authURL: "https://github.test/login",
			profile: user.GitHubProfile{
				GitHubID:  123,
				Username:  "octocat",
				AvatarURL: "https://example.com/avatar.png",
			},
		},
		users,
		sessions,
		memory.NewCurrentUserRepository(users, sessions),
		tokens,
	)

	return NewAuthController(
		usecase,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		security.NewSignedValueCodec("test-secret"),
		false,
	)
}

func findCookie(cookies []*stdhttp.Cookie, name string) *stdhttp.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}
