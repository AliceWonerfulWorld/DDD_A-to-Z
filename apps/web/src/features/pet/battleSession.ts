import type { BattleOpponent, BattleResult, PetSummary } from "./api";

const battleSessionStorageKey = "lang-war.pet.battle-session";

export interface BattleSession {
  playerPet: PetSummary;
  opponent: BattleOpponent;
  result: BattleResult;
}

export function saveBattleSession(session: BattleSession): void {
  try {
    window.sessionStorage.setItem(battleSessionStorageKey, JSON.stringify(session));
  } catch {
    // /battle can still fall back to sample replay data.
  }
}

export function readBattleSession(): BattleSession | null {
  try {
    const raw = window.sessionStorage.getItem(battleSessionStorageKey);
    if (!raw) return null;
    return JSON.parse(raw) as BattleSession;
  } catch {
    return null;
  }
}
