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
  let playerHP = playerPet.maxHp;
  let enemyHP = opponent.pet.maxHp;
  const turns: BattleTurnLog[] = [];
  const playerStarts = playerPet.speed >= opponent.pet.speed;

  const calculateDamage = (attacker: PetSummary, defender: PetSummary, turn: number) => {
    const base = Math.max(
      2,
      attacker.power + Math.ceil(attacker.speed / 3) - Math.floor(defender.guard / 3),
    );
    const critical = turn % 4 === 0 || attacker.speed - defender.speed >= 4;
    return critical ? base + 4 : base;
  };

  for (let round = 1; round <= 30 && playerHP > 0 && enemyHP > 0; round += 1) {
    const playerTurn = playerStarts ? round % 2 === 1 : round % 2 === 0;
    const attacker = playerTurn ? playerPet : opponent.pet;
    const defender = playerTurn ? opponent.pet : playerPet;
    const damage = calculateDamage(attacker, defender, round);
    const defenderName = playerTurn ? opponent.pet.name : playerPet.name;
    const attackerName = attacker.name;
    const isCritical = damage >= attacker.power + 5;

    if (playerTurn) {
      enemyHP = Math.max(0, enemyHP - damage);
    } else {
      playerHP = Math.max(0, playerHP - damage);
    }

    turns.push({
      turn: round,
      actorPetId: attacker.id,
      targetPetId: defender.id,
      damage,
      message: isCritical
        ? `${attackerName} の会心アタック！ ${defenderName} に ${damage} ダメージ。`
        : `${attackerName} が踏み込んだ！ ${defenderName} に ${damage} ダメージ。`,
    });
  }

  const result = playerHP === enemyHP ? "draw" : playerHP > enemyHP ? "win" : "lose";

  return {
    result,
    turns,
  };
}
