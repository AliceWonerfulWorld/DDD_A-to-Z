package guildtown

import (
	"context"
	"errors"
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type testRepository struct {
	inventory           []guildtowndomain.InventoryItem
	placements          []guildtowndomain.Placement
	replaced            []guildtowndomain.Placement
	bought              []guildtowndomain.BuildingType
	created             []guildtowndomain.Placement
	upgradedID          guildtowndomain.PlacementID
	upgradedNextLevel   int
	addedExp            int64
	listPlacementsCalls int
}

func (r *testRepository) ListInventory(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.InventoryItem, error) {
	return r.inventory, nil
}

func (r *testRepository) ListPlacements(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.Placement, error) {
	r.listPlacementsCalls++
	return r.placements, nil
}

func (r *testRepository) ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error {
	r.replaced = placements
	r.placements = placements
	return nil
}

func (r *testRepository) BuyBuilding(ctx context.Context, guildID guilddomain.ID, buildingType guildtowndomain.BuildingType, exp int64, now time.Time) (guilddomain.Guild, error) {
	r.bought = append(r.bought, buildingType)
	r.addedExp = exp
	for index, item := range r.inventory {
		if item.BuildingType == buildingType {
			r.inventory[index].Quantity++
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
	}
	r.inventory = append(r.inventory, guildtowndomain.InventoryItem{
		GuildID:      guildID,
		BuildingType: buildingType,
		Quantity:     1,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
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

func (r *testRepository) CreatePlacement(ctx context.Context, guildID guilddomain.ID, placement guildtowndomain.Placement) error {
	r.created = append(r.created, placement)
	r.placements = append(r.placements, placement)
	return nil
}

func (r *testRepository) UpgradePlacement(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID, nextLevel int, exp int64, now time.Time) (guilddomain.Guild, error) {
	r.upgradedID = placementID
	r.upgradedNextLevel = nextLevel
	r.addedExp = exp
	for index, placement := range r.placements {
		if placement.ID == placementID {
			r.placements[index].Level = nextLevel
		}
	}
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

type testCurrentUserRepository struct {
	appUser user.User
	ok      bool
}

func (r testCurrentUserRepository) FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error) {
	return r.appUser, r.ok, nil
}

type testGuildRepository struct {
	membership guilddomain.MembershipWithGuild
	ok         bool
}

func (r testGuildRepository) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	return r.membership, r.ok, nil
}

type testIDGenerator struct{}

func (g testIDGenerator) NewID() (string, error) {
	return "placement_generated", nil
}

func TestUseCaseSavePlacements(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "tent",
			Quantity:     2,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
	}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testGuildRepository{
		membership: testMembershipWithGuild("guild_go", "user_1", now),
		ok:         true,
	}, testIDGenerator{})
	usecase.now = func() time.Time { return now }

	state, err := usecase.SavePlacements(context.Background(), "session-token", []SavePlacementCommand{{
		BuildingType: "tent",
		X:            12,
		Y:            34,
		Width:        210,
	}})
	if err != nil {
		t.Fatalf("SavePlacements() がエラーを返しました: %v", err)
	}
	if len(repository.replaced) != 1 {
		t.Fatalf("replaced length = %d, 期待値 1", len(repository.replaced))
	}
	if repository.replaced[0].ID != "placement_generated" {
		t.Fatalf("generated id = %q, 期待値 placement_generated", repository.replaced[0].ID)
	}
	if len(state.Placements) != 1 {
		t.Fatalf("state placements length = %d, 期待値 1", len(state.Placements))
	}
	if repository.listPlacementsCalls != 0 {
		t.Fatalf("ListPlacements() calls = %d, 期待値 0", repository.listPlacementsCalls)
	}
}

func TestUseCaseSavePlacementsRejectsInsufficientInventory(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	usecase := NewUseCase(&testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "bonfire",
			Quantity:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
	}, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testGuildRepository{
		membership: testMembershipWithGuild("guild_go", "user_1", now),
		ok:         true,
	}, testIDGenerator{})

	_, err := usecase.SavePlacements(context.Background(), "session-token", []SavePlacementCommand{
		{BuildingType: "bonfire", X: 1, Y: 1, Width: 92},
		{BuildingType: "bonfire", X: 2, Y: 2, Width: 92},
	})
	if !errors.Is(err, ErrInsufficientInventory) {
		t.Fatalf("SavePlacements() error = %v, 期待値 ErrInsufficientInventory", err)
	}
}

func TestUseCaseBuyBuildingAddsPersistentExpForEveryPurchase(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "tent",
			Quantity:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
	}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testGuildRepository{
		membership: testMembershipWithGuild("guild_go", "user_1", now),
		ok:         true,
	}, testIDGenerator{})
	usecase.now = func() time.Time { return now }

	state, err := usecase.BuyBuilding(context.Background(), "session-token", BuyBuildingCommand{BuildingType: "tent"})
	if err != nil {
		t.Fatalf("BuyBuilding() がエラーを返しました: %v", err)
	}
	if repository.addedExp != guilddomain.BuyBuildingExperience {
		t.Fatalf("added exp = %d, 期待値 %d", repository.addedExp, guilddomain.BuyBuildingExperience)
	}
	if len(repository.bought) != 1 || repository.bought[0] != "tent" {
		t.Fatalf("bought = %v, 期待値 [tent]", repository.bought)
	}
	if len(state.Inventory) != 1 || state.Inventory[0].Quantity != 2 {
		t.Fatalf("inventory quantity = %+v, 期待値 2", state.Inventory)
	}
}

func TestUseCaseDeployBuildingDoesNotAddExp(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testGuildRepository{
		membership: testMembershipWithGuild("guild_go", "user_1", now),
		ok:         true,
	}, testIDGenerator{})
	usecase.now = func() time.Time { return now }

	state, err := usecase.DeployBuilding(context.Background(), "session-token", DeployBuildingCommand{
		BuildingType: "tent",
		X:            12,
		Y:            34,
		Width:        210,
	})
	if err != nil {
		t.Fatalf("DeployBuilding() がエラーを返しました: %v", err)
	}
	if repository.addedExp != 0 {
		t.Fatalf("deploy added exp = %d, 期待値 0", repository.addedExp)
	}
	if len(repository.created) != 1 {
		t.Fatalf("created placements length = %d, 期待値 1", len(repository.created))
	}
	if len(state.Placements) != 1 {
		t.Fatalf("state placements length = %d, 期待値 1", len(state.Placements))
	}
}

func TestUseCaseUpgradeBuildingAddsExpByNextLevel(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		placements: []guildtowndomain.Placement{{
			ID:           "placement_1",
			GuildID:      "guild_go",
			BuildingType: "tent",
			Level:        1,
			X:            12,
			Y:            34,
			Width:        210,
			ZIndex:       0,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
	}
	usecase := NewUseCase(repository, testCurrentUserRepository{
		appUser: user.User{ID: "user_1"},
		ok:      true,
	}, testGuildRepository{
		membership: testMembershipWithGuild("guild_go", "user_1", now),
		ok:         true,
	}, testIDGenerator{})
	usecase.now = func() time.Time { return now }

	_, err := usecase.UpgradeBuilding(context.Background(), "session-token", UpgradeBuildingCommand{
		PlacementID: "placement_1",
		NextLevel:   5,
	})
	if err != nil {
		t.Fatalf("UpgradeBuilding() がエラーを返しました: %v", err)
	}
	if repository.addedExp != guilddomain.GuildTownUpgradeLevel5Exp {
		t.Fatalf("upgrade exp = %d, 期待値 %d", repository.addedExp, guilddomain.GuildTownUpgradeLevel5Exp)
	}
	if repository.upgradedNextLevel != 5 {
		t.Fatalf("upgraded next level = %d, 期待値 5", repository.upgradedNextLevel)
	}
}

func testMembershipWithGuild(guildID guilddomain.ID, userID user.ID, now time.Time) guilddomain.MembershipWithGuild {
	return guilddomain.MembershipWithGuild{
		Membership: guilddomain.Membership{
			ID:        "membership_1",
			UserID:    userID,
			GuildID:   guildID,
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Guild: guilddomain.Guild{ID: guildID},
	}
}
