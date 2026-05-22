import { clampValue } from "./townMath";
import type { TownMapRect } from "./townUnlock";
import { getRequiredTownUnlockLevel, isTownRectUnlocked } from "./townUnlock";

export interface DeploymentDraft {
  x: number;
  y: number;
}

export interface TownMapSize {
  height: number;
  width: number;
}

export function getBuildingMapWidth(viewportWidth: number) {
  return clampValue(viewportWidth * 0.14, 112, 220);
}

export function getBuildingUnlockRect({
  itemWidth,
  mapHeight,
  mapWidth,
  x,
  y,
}: {
  itemWidth: number;
  mapHeight: number;
  mapWidth: number;
  x: number;
  y: number;
}): TownMapRect {
  return {
    height: itemWidth * 0.72,
    mapHeight,
    mapWidth,
    width: itemWidth * 0.82,
    x: x + itemWidth * 0.09,
    y: y + itemWidth * 0.18,
  };
}

export function isPlacementUnlocked({
  draft,
  guildLevel,
  itemWidth,
  mapSize,
}: {
  draft: DeploymentDraft;
  guildLevel: number;
  itemWidth: number;
  mapSize: TownMapSize;
}) {
  return isTownRectUnlocked(
    getBuildingUnlockRect({
      itemWidth,
      mapHeight: mapSize.height,
      mapWidth: mapSize.width,
      x: draft.x,
      y: draft.y,
    }),
    guildLevel,
  );
}

export function getRequiredPlacementLevel({
  draft,
  itemWidth,
  mapSize,
}: {
  draft: DeploymentDraft;
  itemWidth: number;
  mapSize: TownMapSize;
}) {
  return getRequiredTownUnlockLevel(
    getBuildingUnlockRect({
      itemWidth,
      mapHeight: mapSize.height,
      mapWidth: mapSize.width,
      x: draft.x,
      y: draft.y,
    }),
  );
}
