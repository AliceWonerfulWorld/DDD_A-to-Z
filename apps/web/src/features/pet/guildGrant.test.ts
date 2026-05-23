import { beforeEach, describe, expect, it, vi } from "vitest";
import { consumeGrantedPet, normalizeGrantedPet, storeGrantedPet } from "./guildGrant";

beforeEach(() => {
  vi.restoreAllMocks();
  const storage = new Map<string, string>();
  vi.stubGlobal("window", {
    sessionStorage: {
      getItem: (key: string) => storage.get(key) ?? null,
      setItem: (key: string, value: string) => storage.set(key, value),
      removeItem: (key: string) => storage.delete(key),
    },
  });
});

describe("guildGrant", () => {
  it("snake_case のギルド参加レスポンスを画面用に変換する", () => {
    expect(
      normalizeGrantedPet({
        id: "pet_1",
        guild_id: "guild_go",
        attribute: "go",
        created_at: "2026-05-23T00:00:00Z",
      }),
    ).toEqual({
      id: "pet_1",
      guildId: "guild_go",
      attribute: "go",
      createdAt: "2026-05-23T00:00:00Z",
    });
  });

  it("camelCase のギルド参加レスポンスはそのまま返す", () => {
    const pet = {
      id: "pet_1",
      guildId: "guild_go",
      attribute: "go",
      createdAt: "2026-05-23T00:00:00Z",
    };

    expect(normalizeGrantedPet(pet)).toEqual(pet);
  });

  it("獲得ペットは一度だけ取り出せる", () => {
    storeGrantedPet({
      id: "pet_1",
      guildId: "guild_go",
      attribute: "go",
      createdAt: "2026-05-23T00:00:00Z",
    });

    expect(consumeGrantedPet()?.guildId).toBe("guild_go");
    expect(consumeGrantedPet()).toBeNull();
  });
});
