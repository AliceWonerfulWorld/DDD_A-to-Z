export type InventoryItemType = "tent" | "bonfire";

export type BuildingBuffType =
  | "daily"
  | "night"
  | "spBoost"
  | "caffeine"
  | "commit"
  | "refactor"
  | "arena"
  | "interest"
  | "plant"
  | "tower"
  | "core";

export type BuildingTargetSpLanguage = "Go" | "TypeScript" | "Rust" | "Python" | "Java" | "Common";

export interface BuildingLevelStatus {
  level: number;
  upgradeCostCp: number;
  upgradeCostSp: number;
  buffValue: number;
}

export interface BuildingMaster {
  id: string;
  name: string;
  description: string;
  requiredGuildLevel: number;
  buffType: BuildingBuffType;
  targetSpLanguage: BuildingTargetSpLanguage;
  levels: BuildingLevelStatus[];
}

export interface InventoryItem {
  type: InventoryItemType;
  name: string;
  title: string;
  description: string;
  count: number;
  src: string;
  minMapWidth: number;
  mapWidthVw: number;
  maxMapWidth: number;
}

export interface PlacedItem {
  id: string;
  type: InventoryItemType;
  name: string;
  title: string;
  description: string;
  src: string;
  x: number;
  y: number;
  width: number;
}

export interface ViewportSize {
  width: number;
  height: number;
}
