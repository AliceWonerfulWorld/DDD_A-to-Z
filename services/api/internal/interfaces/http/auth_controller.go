package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	authapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/auth"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

const (
	oauthStateCookieName = "lang_war_oauth_state"
	oauthStateCookiePath = "/"
	sessionCookieName    = "lang_war_session"
)

type StateCodec interface {
	Sign(value string, expiresAt time.Time) (string, error)
	Verify(signedValue string, now time.Time) (string, error)
}

type AuthController struct {
	usecase      *authapp.UseCase
	logger       *slog.Logger
	stateCodec   StateCodec
	cookieSecure bool
	frontendURL  string
	now          func() time.Time
}

func NewAuthController(
	usecase *authapp.UseCase,
	logger *slog.Logger,
	stateCodec StateCodec,
	cookieSecure bool,
	frontendURL string,
) *AuthController {
	return &AuthController{
		usecase:      usecase,
		logger:       logger,
		stateCodec:   stateCodec,
		cookieSecure: cookieSecure,
		frontendURL:  frontendURL,
		now:          time.Now,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/github/login", c.beginGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", c.completeGitHubLogin)
	mux.HandleFunc("POST /auth/logout", c.logout)
	mux.HandleFunc("GET /me", c.currentUser)
}

func (c *AuthController) beginGitHubLogin(w http.ResponseWriter, r *http.Request) {
	start, err := c.usecase.BeginGitHubLogin(r.Context())
	if err != nil {
		c.logger.Error("failed to begin github login", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	signedState, err := c.stateCodec.Sign(start.State, start.StateExpiresAt)
	if err != nil {
		c.logger.Error("failed to sign oauth state", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    signedState,
		Path:     oauthStateCookiePath,
		Expires:  start.StateExpiresAt,
		HttpOnly: true,
		Secure:   c.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, start.AuthURL, http.StatusFound)
}

func (c *AuthController) completeGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if !c.validOAuthState(r) {
		c.clearOAuthStateCookie(w)
		c.writeLoginError(w, authapp.ErrInvalidState)
		return
	}

	result, err := c.usecase.CompleteGitHubLogin(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		c.clearOAuthStateCookie(w)
		c.writeLoginError(w, err)
		return
	}
	c.clearOAuthStateCookie(w)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    result.Session.Token,
		Path:     "/",
		Expires:  result.Session.ExpiresAt,
		HttpOnly: true,
		Secure:   c.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, c.frontendURL, http.StatusFound)
}

func (c *AuthController) currentUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		c.writeCurrentUserError(w, authapp.ErrUnauthenticated)
		return
	}

	appUser, err := c.usecase.CurrentUser(r.Context(), cookie.Value)
	if err != nil {
		c.writeCurrentUserError(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, userResponse(appUser)); err != nil {
		c.logger.Error("failed to write current user response", "error", err)
	}
}

func (c *AuthController) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		if err := c.usecase.Logout(r.Context(), cookie.Value); err != nil {
			c.logger.Error("failed to logout", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	c.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) writeLoginError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrMissingCode), errors.Is(err, authapp.ErrInvalidState):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		c.logger.Error("failed to complete github login", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *AuthController) writeCurrentUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authapp.ErrUnauthenticated):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		c.logger.Error("failed to get current user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *AuthController) validOAuthState(r *http.Request) bool {
	state := r.URL.Query().Get("state")
	if state == "" {
		return false
	}

	cookie, err := r.Cookie(oauthStateCookieName)
	if err != nil {
		return false
	}

	signedState, err := c.stateCodec.Verify(cookie.Value, c.now())
	return err == nil && signedState == state
}

func (c *AuthController) clearOAuthStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     oauthStateCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AuthController) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func userResponse(appUser user.User) map[string]any {
	return map[string]any{
		"user": map[string]any{
			"id":         appUser.ID,
			"github_id":  appUser.GitHubAccount.GitHubID,
			"username":   appUser.GitHubAccount.Username,
			"avatar_url": appUser.GitHubAccount.AvatarURL,
		},
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	return encoder.Encode(body)
}
