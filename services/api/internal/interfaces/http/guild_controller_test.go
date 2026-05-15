package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type guildTestRepository struct {
	guilds []guilddomain.Guild
}

func (r guildTestRepository) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	return r.guilds, nil
}

func TestGuildControllerListGuilds(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	controller := NewGuildController(guildapp.NewUseCase(guildTestRepository{
		guilds: []guilddomain.Guild{{
			ID:          "guild_typescript",
			Slug:        "typescript",
			Name:        "TypeScript",
			Description: "型の力で支えるギルド。",
			Icon:        "📘",
			Color:       "#3178c6",
			SortOrder:   5,
			MemberCount: 12,
			CreatedAt:   now,
			UpdatedAt:   now,
		}},
	}), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/guilds", nil)
	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		Guilds []struct {
			ID          string `json:"id"`
			Slug        string `json:"slug"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
			Color       string `json:"color"`
			MemberCount int64  `json:"member_count"`
		} `json:"guilds"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if len(body.Guilds) != 1 {
		t.Fatalf("guilds length = %d, 期待値 1", len(body.Guilds))
	}
	if body.Guilds[0].ID != "guild_typescript" {
		t.Fatalf("id = %q, 期待値 guild_typescript", body.Guilds[0].ID)
	}
	if body.Guilds[0].MemberCount != 12 {
		t.Fatalf("member_count = %d, 期待値 12", body.Guilds[0].MemberCount)
	}
}
