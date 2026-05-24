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

  it("non-2xx response は rejection として伝播する", async () => {
    mockFetch(400, { message: "bad request" });

    await expect(fetchMyPets()).rejects.toThrow();

    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/me",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("trainPet は育成対象とステータスを送る", async () => {
    mockFetch(200, {
      pet: {
        id: "pet/go",
        guildId: "guild_go",
        guildName: "Go",
        name: "Gopher",
        species: "gopher",
        attribute: "Go",
        level: 1,
        exp: 0,
        maxHp: 35,
        power: 7,
        guard: 5,
        speed: 7,
        acquiredAt: "2026-05-23T00:00:00Z",
      },
      spentCp: 10,
      cpBalance: 110,
    });

    const result = await trainPet("pet/go", "power");

    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/pet%2Fgo/train",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ stat: "power" }),
      }),
    );
    expect(result).toEqual(
      expect.objectContaining({
        cpBefore: 120,
        cpAfter: 110,
        increasedStat: "power",
        increasedBy: 1,
      }),
    );
  });

  it("fetchBattleOpponents は対戦候補一覧を取得する", async () => {
    mockFetch(200, {
      opponents: [
        {
          petId: "pet_rust",
          guildId: "guild_rust",
          guildName: "Rust",
          displayName: "Rust Challenger",
          name: "Ferris",
          species: "crab",
          attribute: "Rust",
          level: 1,
          maxHp: 45,
          power: 7,
          guard: 7,
          speed: 3,
        },
      ],
    });

    const opponents = await fetchBattleOpponents();

    expect(opponents[0]).toEqual(
      expect.objectContaining({
        userId: "pet_rust",
        petId: "pet_rust",
        playerName: "Rust Challenger",
        pet: expect.objectContaining({ id: "pet_rust", guildName: "Rust" }),
      }),
    );
    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/battle/opponents",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("startPetBattle は対戦相手を送る", async () => {
    mockFetch(200, { result: "loss", turns: [] });

    const result = await startPetBattle("pet_go", "pet_rust");

    expect(fetch).toHaveBeenCalledWith(
      "/api/pets/pet_go/battle",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ opponentPetId: "pet_rust" }),
      }),
    );
    expect(result.result).toBe("lose");
  });
});
