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
	ErrInsufficientInventory    = errors.New("insufficient guild town inventory")
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
	X            float64
	Y            float64
	Width        float64
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
		placement, err := guildtowndomain.NewPlacement(guildtowndomain.Placement{
			ID:           id,
			GuildID:      membership.Membership.GuildID,
			BuildingType: command.BuildingType,
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
