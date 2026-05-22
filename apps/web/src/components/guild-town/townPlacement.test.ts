import { describe, expect, it } from "vitest";
import {
  getBuildingMapWidth,
  getBuildingUnlockRect,
  getRequiredPlacementLevel,
  isPlacementUnlocked,
} from "./townPlacement";

describe("townPlacement", () => {
  it("building width is clamped for small and large viewports", () => {
    expect(getBuildingMapWidth(320)).toBe(112);
    expect(getBuildingMapWidth(1200)).toBeCloseTo(168);
    expect(getBuildingMapWidth(2400)).toBe(220);
  });

  it("uses a smaller unlock rect than the visual building bounds", () => {
    expect(
      getBuildingUnlockRect({
        itemWidth: 200,
        mapHeight: 1000,
        mapWidth: 1600,
        x: 100,
        y: 120,
      }),
    ).toEqual({
      height: 144,
      mapHeight: 1000,
      mapWidth: 1600,
      width: 164,
      x: 118,
      y: 156,
    });
  });

  it("reports whether a placement is inside the current unlocked town area", () => {
    const mapSize = { height: 1000, width: 1000 };

    expect(
      isPlacementUnlocked({
        draft: { x: 450, y: 450 },
        guildLevel: 1,
        itemWidth: 100,
        mapSize,
      }),
    ).toBe(true);
    expect(
      isPlacementUnlocked({
        draft: { x: 860, y: 860 },
        guildLevel: 1,
        itemWidth: 100,
        mapSize,
      }),
    ).toBe(false);
    expect(
      getRequiredPlacementLevel({
        draft: { x: 860, y: 860 },
        itemWidth: 100,
        mapSize,
      }),
    ).toBe(3);
  });
});
