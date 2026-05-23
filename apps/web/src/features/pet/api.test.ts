import { beforeEach, describe, expect, it, vi } from "vitest";
import { fetchBattleOpponents, fetchMyPets, startPetBattle, trainPet } from "./api";

function mockFetch(status: number, body: unknown = null) {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(body),
    }),
  );
}

beforeEach(() => {
  vi.unstubAllGlobals();
});

describe("pet api", () => {
  it("fetchMyPets はマイペット画面用データを取得する", async () => {
    mockFetch(200, {
      cpBalance: 120,
      currentGuildPet: null,
      pets: [],
    });

    const result = await fetchMyPets();

    expect(result.cpBalance).toBe(120);
    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/me",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("trainPet は育成対象とステータスを送る", async () => {
    mockFetch(200, { cpBefore: 120, cpAfter: 110, increasedStat: "power", increasedBy: 1 });

    await trainPet("pet/go", "power");

    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/pet%2Fgo/training",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ stat: "power" }),
      }),
    );
  });

  it("fetchBattleOpponents は対戦候補一覧を取得する", async () => {
    mockFetch(200, { opponents: [] });

    const opponents = await fetchBattleOpponents();

    expect(opponents).toEqual([]);
    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/battle-opponents",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("startPetBattle は対戦相手を送る", async () => {
    mockFetch(200, { result: "win", turns: [] });

    await startPetBattle("user_2");

    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/battles",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ opponentUserId: "user_2" }),
      }),
    );
  });
});
