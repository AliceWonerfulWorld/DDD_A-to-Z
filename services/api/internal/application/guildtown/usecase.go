package guildtown

import (
	"context"
	"errors"
	"strings"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
)

var (
	ErrUnauthenticated          = errors.New("unauthenticated")
	ErrActiveMembershipNotFound = errors.New("active guild membership not found")
	ErrUnknownBuildingType      = errors.New("unknown guild town building type")
	ErrGuildNotFound            = errors.New("guild not found")
	ErrInsufficientInventory    = errors.New("insufficient guild town inventory")
	ErrPlacementNotFound        = errors.New("guild town placement not found")
	ErrInvalidPlacementLevel    = errors.New("guild town placement level is invalid")
)

type UseCase struct {
	repository Repository
	current    CurrentUserRepository
	guilds     GuildRepository
	ids        IDGenerator
	now        func() time.Time
}

type TownState struct {
	Guild      guilddomain.Guild
	Buildings  []guildtowndomain.BuildingMaster
	Inventory  []guildtowndomain.InventoryItem
	Placements []guildtowndomain.Placement
}

type SavePlacementCommand struct {
	ID           guildtowndomain.PlacementID
	BuildingType guildtowndomain.BuildingType
	Level        int
	X            float64
	Y            float64
	Width        float64
}

type BuyBuildingCommand struct {
	BuildingType guildtowndomain.BuildingType
}

type DeployBuildingCommand struct {
	ID           guildtowndomain.PlacementID
	BuildingType guildtowndomain.BuildingType
	X            float64
	Y            float64
	Width        float64
}

type UpgradeBuildingCommand struct {
	PlacementID guildtowndomain.PlacementID
	NextLevel   int
}

func NewUseCase(repository Repository, current CurrentUserRepository, guilds GuildRepository, ids IDGenerator) *UseCase {
	if repository == nil {
		panic("guild town repository is required")
	}
	if current == nil {
		panic("current user repository is required")
	}
	if guilds == nil {
		panic("guild repository is required")
	}
	if ids == nil {
		panic("guild town placement id generator is required")
	}

	return &UseCase{repository: repository, current: current, guilds: guilds, ids: ids, now: time.Now}
}

func (u *UseCase) GetTown(ctx context.Context, sessionToken string) (TownState, error) {
	membership, err := u.requireMembership(ctx, sessionToken)
	if err != nil {
		return TownState{}, err
	}

	inventory, err := u.repository.ListInventory(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}
	inventory = inventoryWithDefaults(membership.Membership.GuildID, inventory, u.now())
	placements, err := u.repository.ListPlacements(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}

	return TownState{
		Guild:      membership.Guild,
		Buildings:  guildtowndomain.DefaultBuildingMasters,
		Inventory:  inventory,
		Placements: placements,
	}, nil
}

func (u *UseCase) SavePlacements(ctx context.Context, sessionToken string, commands []SavePlacementCommand) (TownState, error) {
	membership, err := u.requireMembership(ctx, sessionToken)
	if err != nil {
		return TownState{}, err
	}

	inventory, err := u.repository.ListInventory(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}
	inventory = inventoryWithDefaults(membership.Membership.GuildID, inventory, u.now())
	existingPlacements, err := u.repository.ListPlacements(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}
	levelByPlacementID := make(map[guildtowndomain.PlacementID]int, len(existingPlacements))
	for _, placement := range existingPlacements {
		levelByPlacementID[placement.ID] = placement.Level
	}

	owned := make(map[guildtowndomain.BuildingType]int, len(inventory))
	for _, item := range inventory {
		owned[item.BuildingType] = item.Quantity
	}

	now := u.now()
	used := map[guildtowndomain.BuildingType]int{}
	placements := make([]guildtowndomain.Placement, 0, len(commands))
	for index, command := range commands {
		if _, ok := guildtowndomain.FindBuildingMaster(command.BuildingType); !ok {
			return TownState{}, ErrUnknownBuildingType
		}
		used[command.BuildingType]++
		if used[command.BuildingType] > owned[command.BuildingType] {
			return TownState{}, ErrInsufficientInventory
		}

		id := command.ID
		if id == "" {
			generatedID, err := u.ids.NewID()
			if err != nil {
				return TownState{}, err
			}
			id = guildtowndomain.PlacementID(generatedID)
		}
		level := levelByPlacementID[id]
		if level == 0 {
			level = 1
		}
		placement, err := guildtowndomain.NewPlacement(guildtowndomain.Placement{
			ID:           id,
			GuildID:      membership.Membership.GuildID,
			BuildingType: command.BuildingType,
			Level:        level,
			X:            command.X,
			Y:            command.Y,
			Width:        command.Width,
			ZIndex:       index,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if err != nil {
			return TownState{}, err
		}
		placements = append(placements, placement)
	}

	if err := u.repository.ReplacePlacements(ctx, membership.Membership.GuildID, placements); err != nil {
		return TownState{}, err
	}

	return TownState{
		Guild:      membership.Guild,
		Buildings:  guildtowndomain.DefaultBuildingMasters,
		Inventory:  inventory,
		Placements: placements,
	}, nil
}

func (u *UseCase) BuyBuilding(ctx context.Context, sessionToken string, command BuyBuildingCommand) (TownState, error) {
	membership, err := u.requireMembership(ctx, sessionToken)
	if err != nil {
		return TownState{}, err
	}
	building, ok := guildtowndomain.FindBuildingMaster(command.BuildingType)
	if !ok {
		return TownState{}, ErrUnknownBuildingType
	}

	updatedGuild, err := u.repository.BuyBuilding(
		ctx,
		membership.Membership.UserID,
		membership.Membership.GuildID,
		building,
		guilddomain.BuyBuildingExperience,
		u.now(),
	)
	if err != nil {
		return TownState{}, err
	}

	return u.getTownForGuild(ctx, updatedGuild)
}

func (u *UseCase) DeployBuilding(ctx context.Context, sessionToken string, command DeployBuildingCommand) (TownState, error) {
	membership, err := u.requireMembership(ctx, sessionToken)
	if err != nil {
		return TownState{}, err
	}
	if _, ok := guildtowndomain.FindBuildingMaster(command.BuildingType); !ok {
		return TownState{}, ErrUnknownBuildingType
	}

	inventory, err := u.repository.ListInventory(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}
	inventory = inventoryWithDefaults(membership.Membership.GuildID, inventory, u.now())
	placements, err := u.repository.ListPlacements(ctx, membership.Membership.GuildID)
	if err != nil {
		return TownState{}, err
	}
	owned := 0
	for _, item := range inventory {
		if item.BuildingType == command.BuildingType {
			owned = item.Quantity
			break
		}
	}
	placed := 0
	for _, placement := range placements {
		if placement.BuildingType == command.BuildingType {
			placed++
		}
	}
	if placed >= owned {
		return TownState{}, ErrInsufficientInventory
	}

	id := command.ID
	if id == "" {
		generatedID, err := u.ids.NewID()
		if err != nil {
			return TownState{}, err
		}
		id = guildtowndomain.PlacementID(generatedID)
	}

	now := u.now()
	placement, err := guildtowndomain.NewPlacement(guildtowndomain.Placement{
		ID:           id,
		GuildID:      membership.Membership.GuildID,
		BuildingType: command.BuildingType,
		Level:        1,
		X:            command.X,
		Y:            command.Y,
		Width:        command.Width,
		ZIndex:       0,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return TownState{}, err
	}

	if err := u.repository.CreatePlacement(ctx, membership.Membership.GuildID, placement); err != nil {
		return TownState{}, err
	}

	return u.getTownForGuild(ctx, membership.Guild)
}

func (u *UseCase) UpgradeBuilding(ctx context.Context, sessionToken string, command UpgradeBuildingCommand) (TownState, error) {
	membership, err := u.requireMembership(ctx, sessionToken)
	if err != nil {
		return TownState{}, err
	}
	if command.PlacementID == "" || command.NextLevel < 2 || command.NextLevel > guilddomain.MaxGuildLevel {
		return TownState{}, ErrInvalidPlacementLevel
	}
	placement, ok, err := u.repository.FindPlacementByID(ctx, membership.Membership.GuildID, command.PlacementID)
	if err != nil {
		return TownState{}, err
	}
	if !ok {
		return TownState{}, ErrPlacementNotFound
	}
	if command.NextLevel != placement.Level+1 {
		return TownState{}, ErrInvalidPlacementLevel
	}
	cost, ok := guildtowndomain.FindBuildingLevelCost(placement.BuildingType, command.NextLevel)
	if !ok {
		cost = guildtowndomain.BuildingLevelCost{Level: command.NextLevel}
	}

	exp := guilddomain.CalculateUpgradeExp(command.NextLevel)
	if exp <= 0 {
		return TownState{}, ErrInvalidPlacementLevel
	}

	updatedGuild, err := u.repository.UpgradePlacement(
		ctx,
		membership.Membership.UserID,
		membership.Membership.GuildID,
		command.PlacementID,
		command.NextLevel,
		cost,
		exp,
		u.now(),
	)
	if err != nil {
		return TownState{}, err
	}

	return u.getTownForGuild(ctx, updatedGuild)
}

func (u *UseCase) requireMembership(ctx context.Context, sessionToken string) (guilddomain.MembershipWithGuild, error) {
	if strings.TrimSpace(sessionToken) == "" {
		return guilddomain.MembershipWithGuild{}, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}
	if !ok {
		return guilddomain.MembershipWithGuild{}, ErrUnauthenticated
	}

	membership, ok, err := u.guilds.FindActiveMembershipByUserID(ctx, appUser.ID)
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}
	if !ok {
		return guilddomain.MembershipWithGuild{}, ErrActiveMembershipNotFound
	}

	return membership, nil
}

func (u *UseCase) getTownForGuild(ctx context.Context, guild guilddomain.Guild) (TownState, error) {
	inventory, err := u.repository.ListInventory(ctx, guild.ID)
	if err != nil {
		return TownState{}, err
	}
	inventory = inventoryWithDefaults(guild.ID, inventory, u.now())
	placements, err := u.repository.ListPlacements(ctx, guild.ID)
	if err != nil {
		return TownState{}, err
	}

	return TownState{
		Guild:      guild,
		Buildings:  guildtowndomain.DefaultBuildingMasters,
		Inventory:  inventory,
		Placements: placements,
	}, nil
}

func inventoryWithDefaults(guildID guilddomain.ID, inventory []guildtowndomain.InventoryItem, now time.Time) []guildtowndomain.InventoryItem {
	merged := make([]guildtowndomain.InventoryItem, 0, len(inventory)+len(guildtowndomain.DefaultInventories))
	itemByType := make(map[guildtowndomain.BuildingType]int, len(inventory)+len(guildtowndomain.DefaultInventories))

	for _, item := range inventory {
		itemByType[item.BuildingType] = len(merged)
		merged = append(merged, item)
	}

	for _, defaultInventory := range guildtowndomain.DefaultInventories {
		index, ok := itemByType[defaultInventory.BuildingType]
		if ok {
			if merged[index].Quantity < defaultInventory.Quantity {
				merged[index].Quantity = defaultInventory.Quantity
			}
			continue
		}

		merged = append(merged, guildtowndomain.InventoryItem{
			GuildID:      guildID,
			BuildingType: defaultInventory.BuildingType,
			Quantity:     defaultInventory.Quantity,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	return merged
}
