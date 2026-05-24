package guildtown

import (
	"context"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type Repository interface {
	ListInventory(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.InventoryItem, error)
	ListPlacements(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.Placement, error)
	FindPlacementByID(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID) (guildtowndomain.Placement, bool, error)
	BuyBuilding(ctx context.Context, userID user.ID, guildID guilddomain.ID, building guildtowndomain.BuildingMaster, exp int64, now time.Time) (guilddomain.Guild, error)
	CreatePlacement(ctx context.Context, guildID guilddomain.ID, placement guildtowndomain.Placement) error
	ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error
	UpgradePlacement(ctx context.Context, userID user.ID, guildID guilddomain.ID, placementID guildtowndomain.PlacementID, nextLevel int, cost guildtowndomain.BuildingLevelCost, exp int64, now time.Time) (guilddomain.Guild, error)
}

type CurrentUserRepository interface {
	FindUserBySessionToken(ctx context.Context, sessionToken string, now time.Time) (user.User, bool, error)
}

type GuildRepository interface {
	FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error)
}

type IDGenerator interface {
	NewID() (string, error)
}
