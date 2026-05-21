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
