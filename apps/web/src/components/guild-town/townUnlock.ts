const UNLOCK_RADIUS_BY_LEVEL = [32, 46, 60, 74, 88] as const;
const MAX_UNLOCK_LEVEL = 5;

export interface TownMapPoint {
  mapHeight: number;
  mapWidth: number;
  x: number;
  y: number;
}

export interface TownMapRect {
  height: number;
  mapHeight: number;
  mapWidth: number;
  width: number;
  x: number;
  y: number;
}

export function getTownUnlockRadiusPercent(guildLevel: number) {
  const levelIndex = Math.min(
    Math.max(0, Math.floor(guildLevel) - 1),
    UNLOCK_RADIUS_BY_LEVEL.length - 1,
  );

  return UNLOCK_RADIUS_BY_LEVEL[levelIndex];
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

export function isTownRectUnlocked(rect: TownMapRect, guildLevel: number) {
  const insetX = Math.min(rect.width * 0.18, 24);
  const insetY = Math.min(rect.height * 0.18, 24);
  const checkPoints: TownMapPoint[] = [
    {
      mapHeight: rect.mapHeight,
      mapWidth: rect.mapWidth,
      x: rect.x + rect.width / 2,
      y: rect.y + rect.height / 2,
    },
    { mapHeight: rect.mapHeight, mapWidth: rect.mapWidth, x: rect.x + insetX, y: rect.y + insetY },
    {
      mapHeight: rect.mapHeight,
      mapWidth: rect.mapWidth,
      x: rect.x + rect.width - insetX,
      y: rect.y + insetY,
    },
    {
      mapHeight: rect.mapHeight,
      mapWidth: rect.mapWidth,
      x: rect.x + insetX,
      y: rect.y + rect.height - insetY,
    },
    {
      mapHeight: rect.mapHeight,
      mapWidth: rect.mapWidth,
      x: rect.x + rect.width - insetX,
      y: rect.y + rect.height - insetY,
    },
  ];

  return checkPoints.every((point) => isTownPointUnlocked(point, guildLevel));
}

export function getRequiredTownUnlockLevel(rect: TownMapRect) {
  const requiredRing = getTownUnlockRings().find((ring) => isTownRectUnlocked(rect, ring.level));

  return requiredRing?.level ?? MAX_UNLOCK_LEVEL;
}
