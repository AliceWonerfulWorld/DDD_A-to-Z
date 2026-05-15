import type { PanInfo } from "framer-motion";
import type { InventoryItem } from "./types";

export function clampValue(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value));
}

export function isPointInsideRect(point: PanInfo["point"], rect: DOMRect) {
  return (
    point.x >= rect.left && point.x <= rect.right && point.y >= rect.top && point.y <= rect.bottom
  );
}

export function getInventoryMapWidth(item: InventoryItem, viewportWidth: number) {
  return clampValue(viewportWidth * (item.mapWidthVw / 100), item.minMapWidth, item.maxMapWidth);
}
