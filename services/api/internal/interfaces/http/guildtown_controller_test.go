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
	inventory  []guildtowndomain.InventoryItem
	placements []guildtowndomain.Placement
}

func (r guildTownTestRepository) ListInventory(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.InventoryItem, error) {
	return r.inventory, nil
}

func (r guildTownTestRepository) ListPlacements(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.Placement, error) {
	return r.placements, nil
}

func (r guildTownTestRepository) FindPlacementByID(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID) (guildtowndomain.Placement, bool, error) {
	for _, placement := range r.placements {
		if placement.ID == placementID {
			return placement, true, nil
		}
	}

	return guildtowndomain.Placement{}, false, nil
}

func (r guildTownTestRepository) ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error {
	return nil
}

func (r guildTownTestRepository) BuyBuilding(ctx context.Context, userID user.ID, guildID guilddomain.ID, building guildtowndomain.BuildingMaster, exp int64, now time.Time) (guilddomain.Guild, error) {
	return guilddomain.NewGuild(guilddomain.Guild{
		ID:              guildID,
		Slug:            "go",
		Name:            "Go",
		Description:     "Go guild",
		Icon:            "GO",
		Color:           "#00acd7",
		SortOrder:       1,
		GuildExperience: exp,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

func (r guildTownTestRepository) CreatePlacement(ctx context.Context, guildID guilddomain.ID, placement guildtowndomain.Placement) error {
	return nil
}

func (r guildTownTestRepository) UpgradePlacement(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID, nextLevel int, exp int64, now time.Time) (guilddomain.Guild, error) {
	return guilddomain.NewGuild(guilddomain.Guild{
		ID:              guildID,
		Slug:            "go",
		Name:            "Go",
		Description:     "Go guild",
		Icon:            "GO",
		Color:           "#00acd7",
		SortOrder:       1,
		GuildExperience: exp,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

type guildTownTestCurrentUserRepository struct{}

func (r guildTownTestCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return user.User{ID: user.ID("user_1")}, true, nil
}

type guildTownTestGuildRepository struct{}

func (r guildTownTestGuildRepository) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	guild, err := guilddomain.NewGuild(guilddomain.Guild{
		ID:              guilddomain.ID("guild_go"),
		Slug:            "go",
		Name:            "Go",
		Description:     "Go guild",
		Icon:            "GO",
		Color:           "#00acd7",
		SortOrder:       1,
		GuildExperience: 5200,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return guilddomain.MembershipWithGuild{}, false, err
	}

	return guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    userID,
			GuildID:   "guild_go",
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guild,
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

func TestGuildTownControllerSavePlacementsAcceptsAPIPrefix(t *testing.T) {
	controller := newGuildTownTestController()
	router := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)), controller)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPut, "/api/me/guild/town/placements", strings.NewReader(`{
		"placements": []
	}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusOK)
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
	if body.GuildExperience != 5200 {
		t.Fatalf("guild_experience = %d, 期待値 5200", body.GuildExperience)
	}
	if body.GuildLevel != 2 {
		t.Fatalf("guild_level = %d, 期待値 2", body.GuildLevel)
	}
	if body.NextLevelExp != 20000 {
		t.Fatalf("next_guild_level_experience = %d, 期待値 20000", body.NextLevelExp)
	}
}

func TestGuildTownControllerBuyBuildingAcceptsCamelCaseBuildingID(t *testing.T) {
	controller := newGuildTownTestController()
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/me/guild/town/buildings", strings.NewReader(`{"buildingId":"tent"}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusCreated {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusCreated)
	}

	var body struct {
		CurrentExp int64 `json:"currentExp"`
		GuildLevel int   `json:"guildLevel"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	if body.CurrentExp != 300 {
		t.Fatalf("currentExp = %d, 期待値 300", body.CurrentExp)
	}
	if body.GuildLevel != 1 {
		t.Fatalf("guildLevel = %d, 期待値 1", body.GuildLevel)
	}
}

func TestGuildTownControllerBuyBuildingRejectsMultipleJSONValues(t *testing.T) {
	controller := newGuildTownTestController()
	router := stdhttp.NewServeMux()
	controller.RegisterRoutes(router)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodPost, "/me/guild/town/buildings", strings.NewReader(`{"buildingId":"tent"}{"extra":1}`))
	request.AddCookie(&stdhttp.Cookie{Name: sessionCookieName, Value: "session-token"})

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("ステータスコード = %d, 期待値 %d", response.Code, stdhttp.StatusBadRequest)
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
