import type { BattleOpponent, BattleResult, PetSummary } from "./api";

const battleSessionStorageKey = "lang-war.pet.battle-session";

export interface BattleSession {
  playerPet: PetSummary;
  opponent: BattleOpponent;
  result: BattleResult;
}

function hasObjectShape(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function hasPetSummaryShape(value: unknown): value is PetSummary {
  return (
    hasObjectShape(value) &&
    typeof value.id === "string" &&
    typeof value.name === "string" &&
    typeof value.attribute === "string" &&
    typeof value.maxHp === "number"
  );
}

function hasBattleResultShape(value: unknown): value is BattleResult {
  return (
    hasObjectShape(value) &&
    (value.result === "win" || value.result === "lose" || value.result === "draw") &&
    Array.isArray(value.turns)
  );
}

function validateBattleSession(value: unknown): value is BattleSession {
  return (
    hasObjectShape(value) &&
    hasPetSummaryShape(value.playerPet) &&
    hasObjectShape(value.opponent) &&
    typeof value.opponent.userId === "string" &&
    typeof value.opponent.playerName === "string" &&
    hasPetSummaryShape(value.opponent.pet) &&
    hasBattleResultShape(value.result)
  );
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
    const parsed: unknown = JSON.parse(raw);
    return validateBattleSession(parsed) ? parsed : null;
  } catch {
    return null;
  }
}
