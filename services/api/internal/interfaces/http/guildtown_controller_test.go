package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guildtown"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type guildTownTestRepository struct {
	inventory []guildtowndomain.InventoryItem
}

func (r guildTownTestRepository) ListInventory(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.InventoryItem, error) {
	return r.inventory, nil
}

func (r guildTownTestRepository) ListPlacements(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.Placement, error) {
	return nil, nil
}

func (r guildTownTestRepository) ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error {
	return nil
}

type guildTownTestCurrentUserRepository struct{}

func (r guildTownTestCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return user.User{ID: user.ID("user_1")}, true, nil
}

type guildTownTestGuildRepository struct{}

func (r guildTownTestGuildRepository) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	return guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    userID,
			GuildID:   "guild_go",
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guilddomain.Guild{ID: guilddomain.ID("guild_go"), GuildExperience: 1200},
	}, true, nil
}

type guildTownTestIDGenerator struct{}

func (g guildTownTestIDGenerator) NewID() (string, error) {
	return "placement_1", nil
}

func TestGuildTownControllerSavePlacementsRejectsUnknownFields(t *testing.T) {
	controller := newGuildTownTestController()
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPut, "/me/guild/town/placements", strings.NewReader(`{
		"placements": [],
		"unexpected": true
	}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusBadRequest)
	}
}

func TestGuildTownControllerSavePlacementsRejectsLargeBody(t *testing.T) {
	controller := newGuildTownTestController()
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		stdhttp.MethodPut,
		"/me/guild/town/placements",
		strings.NewReader(strings.Repeat(" ", guildTownPlacementsRequestMaxBytes+1)),
	)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusRequestEntityTooLarge {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusRequestEntityTooLarge)
	}
}

func TestGuildTownControllerGetTownReturnsGuildLevel(t *testing.T) {
	controller := newGuildTownTestController()
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/me/guild/town", nil)
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
	}

	var body struct {
		GuildExperience int64 `json:"guild_experience"`
		GuildLevel      int   `json:"guild_level"`
		NextLevelExp    int64 `json:"next_guild_level_experience"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.GuildExperience != 1200 {
		t.Fatalf("guild_experience = %d, 期待値 1200", body.GuildExperience)
	}
	if body.GuildLevel != 2 {
		t.Fatalf("guild_level = %d, 期待値 2", body.GuildLevel)
	}
	if body.NextLevelExp != 3000 {
		t.Fatalf("next_guild_level_experience = %d, 期待値 3000", body.NextLevelExp)
	}
}

func newGuildTownTestController() *GuildTownController {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	usecase := guildapp.NewUseCase(
		guildTownTestRepository{
			inventory: []guildtowndomain.InventoryItem{{
				GuildID:      guilddomain.ID("guild_go"),
				BuildingType: guildtowndomain.BuildingType("tent"),
				Quantity:     1,
				CreatedAt:    now,
				UpdatedAt:    now,
			}},
		},
		guildTownTestCurrentUserRepository{},
		guildTownTestGuildRepository{},
		guildTownTestIDGenerator{},
	)

	return NewGuildTownController(usecase, slog.New(slog.NewTextHandler(io.Discard, nil)))
}
