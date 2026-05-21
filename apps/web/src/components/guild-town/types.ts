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

export const GUILD_LANGUAGES = [
  "Go",
  "TypeScript",
  "Rust",
  "Python",
  "Java",
  "Haskell",
  "Zig",
  "Common",
] as const;

export type GuildSpLanguage = (typeof GUILD_LANGUAGES)[number];

export type BuildingTargetSpLanguage = GuildSpLanguage;

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
  previewSrc?: string;
  requiredGuildLevel: number;
  buffType: BuildingBuffType;
  targetSpLanguage: BuildingTargetSpLanguage;
  levels: BuildingLevelStatus[];
}

export interface UserInventoryState {
  buildingId: string;
  count: number;
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
  type: string;
  buildingId?: string;
  level: number;
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
