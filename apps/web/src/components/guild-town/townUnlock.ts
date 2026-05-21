import { clampValue } from "./townMath";

const MIN_UNLOCK_RADIUS_PERCENT = 20;
const UNLOCK_RADIUS_STEP_PERCENT = 9;
const MAX_UNLOCK_RADIUS_PERCENT = 54;
const MAX_UNLOCK_LEVEL = 5;

export interface TownMapPoint {
  mapHeight: number;
  mapWidth: number;
  x: number;
  y: number;
}

export function getTownUnlockRadiusPercent(guildLevel: number) {
  return clampValue(
    MIN_UNLOCK_RADIUS_PERCENT + Math.max(0, guildLevel - 1) * UNLOCK_RADIUS_STEP_PERCENT,
    MIN_UNLOCK_RADIUS_PERCENT,
    MAX_UNLOCK_RADIUS_PERCENT,
  );
}

export function getTownUnlockRings(maxLevel = MAX_UNLOCK_LEVEL) {
  return Array.from({ length: maxLevel }, (_, index) => {
    const level = index + 1;

    return {
      level,
      radiusPercent: getTownUnlockRadiusPercent(level),
    };
  });
}

export function isTownPointUnlocked(point: TownMapPoint, guildLevel: number) {
  const mapWidth = Math.max(1, point.mapWidth);
  const mapHeight = Math.max(1, point.mapHeight);
  const radius = getTownUnlockRadiusPercent(guildLevel);
  const dx = ((point.x / mapWidth) * 100 - 50) * 1.08;
  const dy = ((point.y / mapHeight) * 100 - 50) * 0.82;

  return Math.hypot(dx, dy) <= radius;
}
