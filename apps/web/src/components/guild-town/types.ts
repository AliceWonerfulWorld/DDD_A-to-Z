export type InventoryItemType = "tent" | "bonfire";

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
