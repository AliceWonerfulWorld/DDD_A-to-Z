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
	PurchaseCP  int64
	PurchaseSP  int64
	TargetSP    string
}

type DefaultInventory struct {
	BuildingType BuildingType
	Quantity     int
}

type BuildingLevelCost struct {
	Level    int
	CP       int64
	SP       int64
	TargetSP string
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
	{
		Type:        "plasma-condenser",
		Name:        "プラズマ・コンデンサ",
		Title:       "プラズマ・コンデンサ",
		Description: "ログイン時の継続的なCP供給を支えるギルドタウンの基礎電源。",
		Src:         "/build-items/plasma-capacitor.jpeg",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   3,
		PurchaseCP:  120,
		TargetSP:    "Common",
	},
	{
		Type:        "hacker-hideout",
		Name:        "ハッカーの隠れ家",
		Title:       "ハッカーの隠れ家",
		Description: "深夜帯の集中作業から得られるCPを強化する秘密拠点。",
		Src:         "/build-items/Hacker_s_Lean-to.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   4,
		PurchaseCP:  160,
		PurchaseSP:  20,
		TargetSP:    "Common",
	},
	{
		Type:        "llm-semantic-compiler",
		Name:        "LLM・セマンティック解析コンパイラ",
		Title:       "LLM・セマンティック解析コンパイラ",
		Description: "コードから得られる全言語のSP獲得効率を底上げするアレイ。",
		Src:         "/build-items/llm-semantic-compiler-removebg.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   5,
		PurchaseCP:  200,
		PurchaseSP:  20,
		TargetSP:    "Common",
	},
	{
		Type:        "ai-pair-programming-pod",
		Name:        "AIペアプログラミング・ポッド",
		Title:       "AIペアプログラミング・ポッド",
		Description: "コミット由来のCPを大きく伸ばすTypeScript特化の開発支援設備。",
		Src:         "/build-items/AI_Pair_Programming_Pod.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   6,
		PurchaseCP:  240,
		PurchaseSP:  60,
		TargetSP:    "TypeScript",
	},
	{
		Type:        "refactoring-lab",
		Name:        "リファクタリング・ラボ",
		Title:       "リファクタリング・ラボ",
		Description: "将来のレビューイベント拡張に備えた品質改善特化の研究施設。",
		Src:         "/build-items/Refactoring_Lab.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   7,
		PurchaseCP:  280,
		PurchaseSP:  80,
		TargetSP:    "TypeScript",
	},
	{
		Type:        "algorithm-arena",
		Name:        "アルゴリズム闘技場",
		Title:       "アルゴリズム闘技場",
		Description: "将来の学習課題検知に備え、精進由来のCPを強化する演習場。",
		Src:         "/build-items/Algorithm_Arena.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   8,
		PurchaseCP:  320,
		PurchaseSP:  100,
		TargetSP:    "Go",
	},
	{
		Type:        "distributed-data-bank",
		Name:        "分散型データバンク",
		Title:       "分散型データバンク",
		Description: "未消費CPにデイリー利息を発生させる分散管理型の保管庫。",
		Src:         "/build-items/Distributed_Data_Bank.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   9,
		PurchaseCP:  400,
		PurchaseSP:  120,
		TargetSP:    "Go",
	},
	{
		Type:        "ci-cd-automation-plant",
		Name:        "CI/CD自動化プラント",
		Title:       "CI/CD自動化プラント",
		Description: "建築に必要なCPコストを圧縮するRust特化の自動化工場。",
		Src:         "/build-items/CI_CD_Automation_Plant.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   10,
		PurchaseCP:  600,
		PurchaseSP:  200,
		TargetSP:    "Rust",
	},
	{
		Type:        "concurrency-tower",
		Name:        "並行処理タワー",
		Title:       "並行処理タワー",
		Description: "同日中の連続コミットに対するコンボCPを高めるGo特化タワー。",
		Src:         "/build-items/Concurrency_Tower.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   11,
		PurchaseCP:  800,
		PurchaseSP:  240,
		TargetSP:    "Go",
	},
	{
		Type:        "cyber-data-core",
		Name:        "サイバー・データコア",
		Title:       "サイバー・データコア",
		Description: "ギルド全員の基本CPを底上げする最上位の共有中枢。",
		Src:         "/build-items/Cyber_Data_Core.png",
		MinMapWidth: 160,
		MapWidthVW:  20,
		MaxMapWidth: 300,
		SortOrder:   12,
		PurchaseCP:  1500,
		PurchaseSP:  600,
		TargetSP:    "Common",
	},
}

var DefaultInventories = []DefaultInventory{
	{BuildingType: "tent", Quantity: 2},
	{BuildingType: "bonfire", Quantity: 3},
}

var defaultBuildingLevelCosts = map[BuildingType][]BuildingLevelCost{
	"plasma-condenser": {
		{Level: 1, CP: 120, TargetSP: "Common"},
		{Level: 2, CP: 250, TargetSP: "Common"},
		{Level: 3, CP: 600, TargetSP: "Common"},
		{Level: 4, CP: 1500, TargetSP: "Common"},
		{Level: 5, CP: 3500, TargetSP: "Common"},
	},
	"hacker-hideout": {
		{Level: 1, CP: 160, SP: 20, TargetSP: "Common"},
		{Level: 2, CP: 400, SP: 60, TargetSP: "Common"},
		{Level: 3, CP: 1000, SP: 160, TargetSP: "Common"},
		{Level: 4, CP: 2600, SP: 400, TargetSP: "Common"},
		{Level: 5, CP: 6000, SP: 1000, TargetSP: "Common"},
	},
	"llm-semantic-compiler": {
		{Level: 1, CP: 200, SP: 20, TargetSP: "Common"},
		{Level: 2, CP: 480, SP: 80, TargetSP: "Common"},
		{Level: 3, CP: 1200, SP: 200, TargetSP: "Common"},
		{Level: 4, CP: 3200, SP: 480, TargetSP: "Common"},
		{Level: 5, CP: 8000, SP: 1400, TargetSP: "Common"},
	},
	"ai-pair-programming-pod": {
		{Level: 1, CP: 240, SP: 60, TargetSP: "TypeScript"},
		{Level: 2, CP: 600, SP: 160, TargetSP: "TypeScript"},
		{Level: 3, CP: 1600, SP: 400, TargetSP: "TypeScript"},
		{Level: 4, CP: 4000, SP: 1100, TargetSP: "TypeScript"},
		{Level: 5, CP: 8000, SP: 2800, TargetSP: "TypeScript"},
	},
	"refactoring-lab": {
		{Level: 1, CP: 280, SP: 80, TargetSP: "TypeScript"},
		{Level: 2, CP: 720, SP: 200, TargetSP: "TypeScript"},
		{Level: 3, CP: 1800, SP: 520, TargetSP: "TypeScript"},
		{Level: 4, CP: 4800, SP: 1400, TargetSP: "TypeScript"},
		{Level: 5, CP: 10000, SP: 3400, TargetSP: "TypeScript"},
	},
	"algorithm-arena": {
		{Level: 1, CP: 320, SP: 100, TargetSP: "Go"},
		{Level: 2, CP: 800, SP: 260, TargetSP: "Go"},
		{Level: 3, CP: 2200, SP: 720, TargetSP: "Go"},
		{Level: 4, CP: 6000, SP: 1800, TargetSP: "Go"},
		{Level: 5, CP: 12000, SP: 4000, TargetSP: "Go"},
	},
	"distributed-data-bank": {
		{Level: 1, CP: 400, SP: 120, TargetSP: "Go"},
		{Level: 2, CP: 1000, SP: 320, TargetSP: "Go"},
		{Level: 3, CP: 2800, SP: 880, TargetSP: "Go"},
		{Level: 4, CP: 7200, SP: 2200, TargetSP: "Go"},
		{Level: 5, CP: 14000, SP: 4800, TargetSP: "Go"},
	},
	"ci-cd-automation-plant": {
		{Level: 1, CP: 600, SP: 200, TargetSP: "Rust"},
		{Level: 2, CP: 1400, SP: 480, TargetSP: "Rust"},
		{Level: 3, CP: 3600, SP: 1200, TargetSP: "Rust"},
		{Level: 4, CP: 8800, SP: 3000, TargetSP: "Rust"},
		{Level: 5, CP: 18000, SP: 7200, TargetSP: "Rust"},
	},
	"concurrency-tower": {
		{Level: 1, CP: 800, SP: 240, TargetSP: "Go"},
		{Level: 2, CP: 1800, SP: 600, TargetSP: "Go"},
		{Level: 3, CP: 4400, SP: 1600, TargetSP: "Go"},
		{Level: 4, CP: 10400, SP: 3600, TargetSP: "Go"},
		{Level: 5, CP: 20000, SP: 8000, TargetSP: "Go"},
	},
	"cyber-data-core": {
		{Level: 1, CP: 1500, SP: 600, TargetSP: "Common"},
		{Level: 2, CP: 3600, SP: 1500, TargetSP: "Common"},
		{Level: 3, CP: 7500, SP: 3000, TargetSP: "Common"},
		{Level: 4, CP: 12000, SP: 5400, TargetSP: "Common"},
		{Level: 5, CP: 18000, SP: 7500, TargetSP: "Common"},
	},
}

func FindBuildingLevelCost(buildingType BuildingType, level int) (BuildingLevelCost, bool) {
	for _, cost := range defaultBuildingLevelCosts[buildingType] {
		if cost.Level == level {
			return cost, true
		}
	}

	return BuildingLevelCost{}, false
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

func DefaultInventoryQuantity(buildingType BuildingType) int {
	for _, inventory := range DefaultInventories {
		if inventory.BuildingType == buildingType {
			return inventory.Quantity
		}
	}

	return 0
}

func FindBuildingMaster(buildingType BuildingType) (BuildingMaster, bool) {
	for _, master := range DefaultBuildingMasters {
		if master.Type == buildingType {
			return master, true
		}
	}

	return BuildingMaster{}, false
}
