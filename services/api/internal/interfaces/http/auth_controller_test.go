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

func (c fakeGitHubOAuthClient) ExchangeLogin(ctx context.Context, code string) (authapp.GitHubLogin, error) {
	if err := ctx.Err(); err != nil {
		return authapp.GitHubLogin{}, err
	}

	return authapp.GitHubLogin{Profile: c.profile, AccessToken: "github-token"}, nil
}

type fakeTokenGenerator struct {
	tokens []string
}

func (g *fakeTokenGenerator) NewToken() (string, error) {
	token := g.tokens[0]
	g.tokens = g.tokens[1:]
	return token, nil
}

func TestAuthControllerBeginGitHubLogin(t *testing.T) {
	t.Run("GitHubログイン開始で認可URLへリダイレクトする", func(t *testing.T) {
		controller := newTestAuthController(&fakeTokenGenerator{tokens: []string{"state-token"}})
		router := stdhttp.NewServeMux()
		controller.RegisterRoutes(router)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/login", nil)
		router.ServeHTTP(response, request)

		if response.Code != stdhttp.StatusFound {
			t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusFound)
		}

		if got := response.Header().Get("Location"); got != "https://github.test/login?state=state-token" {
			t.Fatalf("Location = %q, 期待値 GitHub 認可URL", got)
		}
		if cookie := findCookie(response.Result().Cookies(), oauthStateCookieName); cookie == nil {
			t.Fatal("OAuth state Cookie が設定されている必要があります")
		} else if cookie.Path != "/" {
			t.Fatalf("OAuth state Cookie Path = %q, 期待値 /", cookie.Path)
		}
	})
}

func TestAuthControllerCompleteGitHubLogin(t *testing.T) {
	t.Run("GitHubログイン完了でセッションCookieを発行する", func(t *testing.T) {
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

		if response.Code != stdhttp.StatusFound {
			t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusFound)
		}
		if got := response.Header().Get("Location"); got != "https://app.test" {
			t.Fatalf("Location = %q, 期待値 https://app.test", got)
		}

		cookie := findCookie(response.Result().Cookies(), sessionCookieName)
		if cookie == nil {
			t.Fatal("lang_war_session Cookie が設定されている必要があります")
		}
		if cookie.Value != "session-token" {
			t.Fatalf("セッションCookie = %q, 期待値 session-token", cookie.Value)
		}
		if !cookie.HttpOnly {
			t.Fatal("セッションCookie は HttpOnly である必要があります")
		}
		if cookie.SameSite != stdhttp.SameSiteLaxMode {
			t.Fatalf("SameSite = %v, 期待値 Lax", cookie.SameSite)
		}

		stateCookie := findCookie(response.Result().Cookies(), oauthStateCookieName)
		if stateCookie == nil {
			t.Fatal("OAuth state Cookie を削除する Set-Cookie が必要です")
		}
		if stateCookie.Path != "/" {
			t.Fatalf("OAuth state Cookie 削除 Path = %q, 期待値 /", stateCookie.Path)
		}
		if stateCookie.MaxAge != -1 {
			t.Fatalf("OAuth state Cookie MaxAge = %d, 期待値 -1", stateCookie.MaxAge)
		}
	})
}

func TestAuthControllerCurrentUser(t *testing.T) {
	t.Run("ログイン済みユーザーを返す", func(t *testing.T) {
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
			t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
		}

		var body struct {
			User struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if body.User.ID == "" {
			t.Fatal("レスポンスのユーザーIDが設定されている必要があります")
		}
		if body.User.Username != "octocat" {
			t.Fatalf("username = %q, 期待値 octocat", body.User.Username)
		}
	})
}

func TestAuthControllerLogout(t *testing.T) {
	t.Run("ログアウトでセッションCookieを削除し現在のユーザーを取得できなくする", func(t *testing.T) {
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

		sessionCookie := findCookie(callbackResponse.Result().Cookies(), sessionCookieName)
		if sessionCookie == nil {
			t.Fatal("lang_war_session Cookie が設定されている必要があります")
		}

		logoutResponse := httptest.NewRecorder()
		logoutRequest := httptest.NewRequest(stdhttp.MethodPost, "/auth/logout", nil)
		logoutRequest.AddCookie(sessionCookie)
		router.ServeHTTP(logoutResponse, logoutRequest)

		if logoutResponse.Code != stdhttp.StatusNoContent {
			t.Fatalf("ステータスコード = %d, 期待値 %d", logoutResponse.Code, stdhttp.StatusNoContent)
		}
		clearedCookie := findCookie(logoutResponse.Result().Cookies(), sessionCookieName)
		if clearedCookie == nil {
			t.Fatal("セッションCookie を削除する Set-Cookie が必要です")
		}
		if clearedCookie.Path != "/" {
			t.Fatalf("セッションCookie 削除 Path = %q, 期待値 /", clearedCookie.Path)
		}
		if clearedCookie.MaxAge != -1 {
			t.Fatalf("セッションCookie MaxAge = %d, 期待値 -1", clearedCookie.MaxAge)
		}

		currentUserResponse := httptest.NewRecorder()
		currentUserRequest := httptest.NewRequest(stdhttp.MethodGet, "/me", nil)
		currentUserRequest.AddCookie(sessionCookie)
		router.ServeHTTP(currentUserResponse, currentUserRequest)

		if currentUserResponse.Code != stdhttp.StatusUnauthorized {
			t.Fatalf("ステータスコード = %d, 期待値 %d", currentUserResponse.Code, stdhttp.StatusUnauthorized)
		}
	})
}

func TestAuthControllerCompleteGitHubLoginRejectsInvalidState(t *testing.T) {
	t.Run("GitHubログイン完了で不正なStateを拒否する", func(t *testing.T) {
		controller := newTestAuthController(&fakeTokenGenerator{})
		router := stdhttp.NewServeMux()
		controller.RegisterRoutes(router)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(stdhttp.MethodGet, "/auth/github/callback?code=github-code&state=missing-state", nil)
		router.ServeHTTP(response, request)

		if response.Code != stdhttp.StatusBadRequest {
			t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusBadRequest)
		}
	})
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
		"https://app.test",
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
