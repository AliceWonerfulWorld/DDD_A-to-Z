package http

import (
	"errors"
	"log/slog"
	stdhttp "net/http"
	"time"

	chatapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/chat"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type ChatController struct {
	usecase *chatapp.UseCase
	logger  *slog.Logger
}

func NewChatController(usecase *chatapp.UseCase, logger *slog.Logger) *ChatController {
	return &ChatController{
		usecase: usecase,
		logger:  logger,
	}
}

func (c *ChatController) RegisterRoutes(mux *stdhttp.ServeMux) {
	mux.HandleFunc("POST /guilds/{guildID}/chat-token", c.issueChatToken)
}

func (c *ChatController) issueChatToken(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
		return
	}

	token, err := c.usecase.IssueGuildChatToken(r.Context(), cookie.Value, guilddomain.ID(r.PathValue("guildID")))
	if err != nil {
		c.writeError(w, err)
		return
	}

	if err := writeJSON(w, stdhttp.StatusOK, map[string]any{
		"token":      token.Token,
		"expires_at": token.ExpiresAt.Format(time.RFC3339),
	}); err != nil {
		c.logger.Error("failed to write chat token response", "error", err)
	}
}

func (c *ChatController) writeError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, chatapp.ErrUnauthenticated):
		writeAPIError(w, stdhttp.StatusUnauthorized, "unauthenticated", "unauthenticated", 0, nil)
	case errors.Is(err, chatapp.ErrForbidden):
		writeAPIError(w, stdhttp.StatusForbidden, "forbidden", "forbidden", 0, nil)
	default:
		c.logger.Error("chat token request failed", "error", err)
		writeAPIError(w, stdhttp.StatusInternalServerError, "internal_error", "Internal Server Error", 0, nil)
	}
}
