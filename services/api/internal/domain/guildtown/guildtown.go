// Package guildtown owns guild town building inventory and placement concepts.
package guildtown

import (
	"errors"
	"strings"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

type BuildingType string

type BuildingMaster struct {
	Type        BuildingType
	Name        string
	Title       string
	Description string
	Src         string
	MinMapWidth int
	MapWidthVW  int
	MaxMapWidth int
	SortOrder   int
}

type InventoryItem struct {
	GuildID      guilddomain.ID
	BuildingType BuildingType
	Quantity     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PlacementID string

type Placement struct {
	ID           PlacementID
	GuildID      guilddomain.ID
	BuildingType BuildingType
	Level        int
	X            float64
	Y            float64
	Width        float64
	ZIndex       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var DefaultBuildingMasters = []BuildingMaster{
	{
		Type:        "tent",
		Name:        "TENT",
		Title:       "旅人のテント",
		Description: "ギルドの仲間が遠征前に集う簡易拠点。休息と作戦会議に使われる。",
		Src:         "/town/tent.png",
		MinMapWidth: 210,
		MapWidthVW:  29,
		MaxMapWidth: 430,
		SortOrder:   1,
	},
	{
		Type:        "bonfire",
		Name:        "BONFIRE",
		Title:       "団らんの焚き火",
		Description: "夜のギルドタウンを照らす小さな火。仲間の士気をじんわり温める。",
		Src:         "/town/bonfire.png",
		MinMapWidth: 92,
		MapWidthVW:  12,
		MaxMapWidth: 164,
		SortOrder:   2,
	},
}

func NewInventoryItem(item InventoryItem) (InventoryItem, error) {
	if item.GuildID == "" {
		return InventoryItem{}, errors.New("guild town inventory guild id is required")
	}
	if strings.TrimSpace(string(item.BuildingType)) == "" {
		return InventoryItem{}, errors.New("guild town inventory building type is required")
	}
	if _, ok := FindBuildingMaster(item.BuildingType); !ok {
		return InventoryItem{}, errors.New("guild town inventory building type is unknown")
	}
	if item.Quantity < 0 {
		return InventoryItem{}, errors.New("guild town inventory quantity cannot be negative")
	}
	if item.CreatedAt.IsZero() {
		return InventoryItem{}, errors.New("guild town inventory created at is required")
	}
	if item.UpdatedAt.IsZero() {
		return InventoryItem{}, errors.New("guild town inventory updated at is required")
	}

	return item, nil
}

func NewPlacement(placement Placement) (Placement, error) {
	if placement.ID == "" {
		return Placement{}, errors.New("guild town placement id is required")
	}
	if placement.GuildID == "" {
		return Placement{}, errors.New("guild town placement guild id is required")
	}
	if strings.TrimSpace(string(placement.BuildingType)) == "" {
		return Placement{}, errors.New("guild town placement building type is required")
	}
	if _, ok := FindBuildingMaster(placement.BuildingType); !ok {
		return Placement{}, errors.New("guild town placement building type is unknown")
	}
	if placement.Level == 0 {
		placement.Level = 1
	}
	if placement.Level < 1 || placement.Level > guilddomain.MaxGuildLevel {
		return Placement{}, errors.New("guild town placement level must be between 1 and 5")
	}
	if placement.X < 0 {
		return Placement{}, errors.New("guild town placement x cannot be negative")
	}
	if placement.Y < 0 {
		return Placement{}, errors.New("guild town placement y cannot be negative")
	}
	if placement.Width <= 0 {
		return Placement{}, errors.New("guild town placement width must be positive")
	}
	if placement.ZIndex < 0 {
		return Placement{}, errors.New("guild town placement z index cannot be negative")
	}
	if placement.CreatedAt.IsZero() {
		return Placement{}, errors.New("guild town placement created at is required")
	}
	if placement.UpdatedAt.IsZero() {
		return Placement{}, errors.New("guild town placement updated at is required")
	}

	return placement, nil
}

func FindBuildingMaster(buildingType BuildingType) (BuildingMaster, bool) {
	for _, master := range DefaultBuildingMasters {
		if master.Type == buildingType {
			return master, true
		}
	}

	return BuildingMaster{}, false
}
