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

func (r *testRepository) FindPlacementByID(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID) (guildtowndomain.Placement, bool, error) {
	for _, placement := range r.placements {
		if placement.ID == placementID {
			return placement, true, nil
		}
	}

	return guildtowndomain.Placement{}, false, nil
}

func (r *testRepository) ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error {
	r.replaced = placements
	r.placements = placements
	return nil
}

func (r *testRepository) BuyBuilding(ctx context.Context, userID user.ID, guildID guilddomain.ID, building guildtowndomain.BuildingMaster, exp int64, now time.Time) (guilddomain.Guild, error) {
	r.bought = append(r.bought, building.Type)
	r.addedExp += exp
	for index, item := range r.inventory {
		if item.BuildingType == building.Type {
			r.inventory[index].Quantity++
			return guilddomain.NewGuild(guilddomain.Guild{
				ID:              guildID,
				Slug:            "go",
				Name:            "Go",
				Description:     "Go guild",
				Icon:            "GO",
				Color:           "#00acd7",
				SortOrder:       1,
				GuildExperience: r.addedExp,
				CreatedAt:       now,
				UpdatedAt:       now,
			})
		}
	}
	r.inventory = append(r.inventory, guildtowndomain.InventoryItem{
		GuildID:      guildID,
		BuildingType: building.Type,
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
		GuildExperience: r.addedExp,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

func (r *testRepository) CreatePlacement(ctx context.Context, guildID guilddomain.ID, placement guildtowndomain.Placement) error {
	r.created = append(r.created, placement)
	r.placements = append(r.placements, placement)
	return nil
}

func (r *testRepository) UpgradePlacement(ctx context.Context, userID user.ID, guildID guilddomain.ID, placementID guildtowndomain.PlacementID, nextLevel int, cost guildtowndomain.BuildingLevelCost, exp int64, now time.Time) (guilddomain.Guild, error) {
	r.upgradedID = placementID
	r.upgradedNextLevel = nextLevel
	r.addedExp += exp
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
		GuildExperience: r.addedExp,
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
	if repository.replaced[0].Level != 1 {
		t.Fatalf("placement level = %d, 期待値 1", repository.replaced[0].Level)
	}
	if repository.listPlacementsCalls != 1 {
		t.Fatalf("ListPlacements() calls = %d, 期待値 1", repository.listPlacementsCalls)
	}
}

func TestUseCaseSavePlacementsUsesDefaultInventory(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{}
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
	if len(state.Inventory) != 2 {
		t.Fatalf("state inventory length = %d, 期待値 2", len(state.Inventory))
	}
}

func TestUseCaseSavePlacementsAcceptsStoreBuilding(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "plasma-condenser",
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

	_, err := usecase.SavePlacements(context.Background(), "session-token", []SavePlacementCommand{{
		BuildingType: "plasma-condenser",
		X:            12,
		Y:            34,
		Width:        210,
	}})
	if err != nil {
		t.Fatalf("SavePlacements() がエラーを返しました: %v", err)
	}
	if repository.replaced[0].BuildingType != "plasma-condenser" {
		t.Fatalf("replaced building type = %q, 期待値 plasma-condenser", repository.replaced[0].BuildingType)
	}
}

func TestUseCaseSavePlacementsKeepsCanonicalPlacementLevel(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "tent",
			Quantity:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
		placements: []guildtowndomain.Placement{{
			ID:           "placement_1",
			GuildID:      "guild_go",
			BuildingType: "tent",
			Level:        3,
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

	_, err := usecase.SavePlacements(context.Background(), "session-token", []SavePlacementCommand{{
		ID:           "placement_1",
		BuildingType: "tent",
		Level:        5,
		X:            20,
		Y:            40,
		Width:        210,
	}})
	if err != nil {
		t.Fatalf("SavePlacements() がエラーを返しました: %v", err)
	}
	if repository.replaced[0].Level != 3 {
		t.Fatalf("placement level = %d, 期待値 3", repository.replaced[0].Level)
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
		{BuildingType: "bonfire", X: 3, Y: 3, Width: 92},
		{BuildingType: "bonfire", X: 4, Y: 4, Width: 92},
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
	if quantity := inventoryQuantity(state.Inventory, "tent"); quantity != 2 {
		t.Fatalf("tent inventory quantity = %d, 期待値 2", quantity)
	}

	state, err = usecase.BuyBuilding(context.Background(), "session-token", BuyBuildingCommand{BuildingType: "tent"})
	if err != nil {
		t.Fatalf("BuyBuilding() 2回目がエラーを返しました: %v", err)
	}
	if repository.addedExp != guilddomain.BuyBuildingExperience*2 {
		t.Fatalf("added exp after second purchase = %d, 期待値 %d", repository.addedExp, guilddomain.BuyBuildingExperience*2)
	}
	if len(repository.bought) != 2 || repository.bought[0] != "tent" || repository.bought[1] != "tent" {
		t.Fatalf("bought = %v, 期待値 [tent tent]", repository.bought)
	}
	if quantity := inventoryQuantity(state.Inventory, "tent"); quantity != 3 {
		t.Fatalf("tent inventory quantity after second purchase = %d, 期待値 3", quantity)
	}
}

func TestUseCaseDeployBuildingDoesNotAddExp(t *testing.T) {
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

func TestUseCaseDeployBuildingUsesDefaultInventory(t *testing.T) {
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

	_, err := usecase.DeployBuilding(context.Background(), "session-token", DeployBuildingCommand{
		BuildingType: "bonfire",
		X:            12,
		Y:            34,
		Width:        92,
	})
	if err != nil {
		t.Fatalf("DeployBuilding() がエラーを返しました: %v", err)
	}
	if len(repository.created) != 1 {
		t.Fatalf("created placements length = %d, 期待値 1", len(repository.created))
	}
}

func TestUseCaseDeployBuildingRejectsInsufficientInventory(t *testing.T) {
	now := time.Date(2026, 5, 22, 9, 0, 0, 0, time.UTC)
	repository := &testRepository{
		inventory: []guildtowndomain.InventoryItem{{
			GuildID:      "guild_go",
			BuildingType: "tent",
			Quantity:     1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}},
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
		}, {
			ID:           "placement_2",
			GuildID:      "guild_go",
			BuildingType: "tent",
			Level:        1,
			X:            24,
			Y:            48,
			Width:        210,
			ZIndex:       1,
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

	_, err := usecase.DeployBuilding(context.Background(), "session-token", DeployBuildingCommand{
		BuildingType: "tent",
		X:            12,
		Y:            34,
		Width:        210,
	})
	if !errors.Is(err, ErrInsufficientInventory) {
		t.Fatalf("DeployBuilding() error = %v, 期待値 ErrInsufficientInventory", err)
	}
	if len(repository.created) != 0 {
		t.Fatalf("created placements length = %d, 期待値 0", len(repository.created))
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
		NextLevel:   2,
	})
	if err != nil {
		t.Fatalf("UpgradeBuilding() がエラーを返しました: %v", err)
	}
	if repository.addedExp != guilddomain.GuildTownUpgradeLevel2Exp {
		t.Fatalf("upgrade exp = %d, 期待値 %d", repository.addedExp, guilddomain.GuildTownUpgradeLevel2Exp)
	}
	if repository.upgradedNextLevel != 2 {
		t.Fatalf("upgraded next level = %d, 期待値 2", repository.upgradedNextLevel)
	}
}

func TestUseCaseUpgradeBuildingRejectsSkippedLevel(t *testing.T) {
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
	if !errors.Is(err, ErrInvalidPlacementLevel) {
		t.Fatalf("UpgradeBuilding() error = %v, 期待値 ErrInvalidPlacementLevel", err)
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

func inventoryQuantity(inventory []guildtowndomain.InventoryItem, buildingType guildtowndomain.BuildingType) int {
	for _, item := range inventory {
		if item.BuildingType == buildingType {
			return item.Quantity
		}
	}

	return 0
}
