import { beforeEach, describe, expect, it, vi } from "vitest";
import { readBattleSession, saveBattleSession } from "./battleSession";
import { sampleCurrentPet, sampleOpponents } from "./sampleData";
import { buildSampleBattleResult } from "./battleReplay";

beforeEach(() => {
  vi.restoreAllMocks();
  const storage = new Map<string, string>();
  vi.stubGlobal("window", {
    sessionStorage: {
      getItem: (key: string) => storage.get(key) ?? null,
      setItem: (key: string, value: string) => storage.set(key, value),
    },
  });
});

describe("battleSession", () => {
  it("battle session を保存して読み出せる", () => {
    const opponent = sampleOpponents[0]!;
    const result = buildSampleBattleResult(sampleCurrentPet, opponent);

    saveBattleSession({ playerPet: sampleCurrentPet, opponent, result });

    expect(readBattleSession()?.opponent.userId).toBe(opponent.userId);
  });
});
