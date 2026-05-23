import type { BattleOpponent, BattleResult, BattleTurnLog, PetSummary } from "./api";

export interface BattleReplayTurn extends BattleTurnLog {
  actorSide: "player" | "enemy";
  targetSide: "player" | "enemy";
  isCritical?: boolean;
  combo?: number;
}

export interface BattleReplay {
  playerPet: PetSummary;
  opponent: BattleOpponent;
  result: BattleResult;
  turns: BattleReplayTurn[];
}

export function toBattleReplay(
  playerPet: PetSummary,
  opponent: BattleOpponent,
  result: BattleResult,
): BattleReplay {
  return {
    playerPet,
    opponent,
    result,
    turns: result.turns.map((turn) => {
      const actorSide = turn.actorPetId === playerPet.id ? "player" : "enemy";
      const targetSide = turn.targetPetId === playerPet.id ? "player" : "enemy";
      return {
        ...turn,
        actorSide,
        targetSide,
        isCritical: turn.damage >= 9 || /会心|critical/i.test(turn.message),
        combo: turn.turn > 1 && actorSide === "player" ? 2 : undefined,
      };
    }),
  };
}

export function buildSampleBattleResult(
  playerPet: PetSummary,
  opponent: BattleOpponent,
): BattleResult {
  const playerHit = Math.max(4, playerPet.power + 2);
  const enemyHit = Math.max(3, opponent.pet.power - 1);
  const finisher = Math.max(playerHit + 2, 9);

  return {
    result: "win",
    turns: [
      {
        turn: 1,
        actorPetId: playerPet.id,
        targetPetId: opponent.pet.id,
        damage: playerHit,
        message: `${playerPet.name} が踏み込んだ！ ${opponent.pet.name} に ${playerHit} ダメージ。`,
      },
      {
        turn: 2,
        actorPetId: opponent.pet.id,
        targetPetId: playerPet.id,
        damage: enemyHit,
        message: `${opponent.pet.name} の反撃。${playerPet.name} は ${enemyHit} ダメージを受けた。`,
      },
      {
        turn: 3,
        actorPetId: playerPet.id,
        targetPetId: opponent.pet.id,
        damage: finisher,
        message: `${playerPet.name} の会心アタック！ ${opponent.pet.name} に ${finisher} ダメージ。`,
      },
    ],
  };
}
