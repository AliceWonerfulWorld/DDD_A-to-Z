package guildtown

import (
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

func TestNewInventoryItemRejectsUnknownBuildingType(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)

	_, err := NewInventoryItem(InventoryItem{
		GuildID:      guilddomain.ID("guild_go"),
		BuildingType: BuildingType("unknown"),
		Quantity:     1,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err == nil {
		t.Fatal("NewInventoryItem() error = nil, 期待値 error")
	}
}

func TestNewPlacementRejectsUnknownBuildingType(t *testing.T) {
	now := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)

	_, err := NewPlacement(Placement{
		ID:           PlacementID("placement_1"),
		GuildID:      guilddomain.ID("guild_go"),
		BuildingType: BuildingType("unknown"),
		X:            1,
		Y:            1,
		Width:        100,
		ZIndex:       0,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err == nil {
		t.Fatal("NewPlacement() error = nil, 期待値 error")
	}
}

func TestFindBuildingLevelCost(t *testing.T) {
	cost, ok := FindBuildingLevelCost("plasma-condenser", 2)
	if !ok {
		t.Fatal("FindBuildingLevelCost() ok = false, 期待値 true")
	}
	if cost.CP != 250 || cost.SP != 0 || cost.TargetSP != "Common" {
		t.Fatalf("cost = %+v, 期待値 CP=250 SP=0 TargetSP=Common", cost)
	}

	cost, ok = FindBuildingLevelCost("concurrency-tower", 5)
	if !ok {
		t.Fatal("FindBuildingLevelCost() ok = false, 期待値 true")
	}
	if cost.CP != 20000 || cost.SP != 8000 || cost.TargetSP != "Go" {
		t.Fatalf("cost = %+v, 期待値 CP=20000 SP=8000 TargetSP=Go", cost)
	}

	if _, ok := FindBuildingLevelCost("tent", 2); ok {
		t.Fatal("FindBuildingLevelCost(tent, 2) ok = true, 期待値 false")
	}
}
