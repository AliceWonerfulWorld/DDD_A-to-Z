import { apiFetch } from "../../lib/api/client";

export type PetTrainingStat = "hp" | "power" | "guard" | "speed";

export interface PetSummary {
  id: string;
  guildId: string;
  guildName: string;
  name: string;
  species: string;
  attribute: string;
  level: number;
  exp: number;
  maxHp: number;
  power: number;
  guard: number;
  speed: number;
  acquiredAt: string;
}

export interface MyPetsResponse {
  cpBalance: number;
  currentGuildPet: PetSummary | null;
  pets: PetSummary[];
}

export interface GrantedPet {
  id: string;
  guildId: string;
  attribute: string;
  createdAt: string;
}

export interface TrainingResult {
  pet: PetSummary;
  cpBefore: number;
  cpAfter: number;
  increasedStat: PetTrainingStat;
  increasedBy: number;
}

export interface BattleOpponent {
  userId: string;
  petId: string;
  playerName: string;
  pet: PetSummary;
}

export interface BattleTurnLog {
  turn: number;
  actorPetId: string;
  targetPetId: string;
  damage: number;
  message: string;
}

export interface BattleResult {
  result: "win" | "lose" | "draw";
  turns: BattleTurnLog[];
}

interface PetTrainingApiResponse {
  pet: PetSummary;
  spentCp: number;
  cpBalance: number;
}

interface BattleOpponentApiResponse {
  petId: string;
  guildId: string;
  guildName: string;
  DisplayName: string;
  name: string;
  species: string;
  attribute: string;
  level: number;
  maxHp: number;
  power: number;
  guard: number;
  speed: number;
}

interface BattleResultApiResponse extends Omit<BattleResult, "result"> {
  result: "win" | "loss" | "draw";
}

export const PET_TRAINING_COSTS: Record<
  PetTrainingStat,
  { label: string; amount: number; cost: number }
> = {
  hp: { label: "HP", amount: 5, cost: 20 },
  power: { label: "Power", amount: 1, cost: 10 },
  guard: { label: "Guard", amount: 1, cost: 10 },
  speed: { label: "Speed", amount: 1, cost: 10 },
};

export async function fetchMyPets(): Promise<MyPetsResponse> {
  return apiFetch<MyPetsResponse>("/pets/me");
}

export async function trainPet(petId: string, stat: PetTrainingStat): Promise<TrainingResult> {
  const result = await apiFetch<PetTrainingApiResponse>(
    `/pets/${encodeURIComponent(petId)}/train`,
    {
      method: "POST",
      body: JSON.stringify({ stat }),
    },
  );
  return {
    pet: result.pet,
    cpBefore: result.cpBalance + result.spentCp,
    cpAfter: result.cpBalance,
    increasedStat: stat,
    increasedBy: PET_TRAINING_COSTS[stat].amount,
  };
}

export async function fetchBattleOpponents(): Promise<BattleOpponent[]> {
  const data = await apiFetch<{ opponents: BattleOpponentApiResponse[] }>("/pets/battle/opponents");
  return data.opponents.map((opponent) => ({
    userId: opponent.petId,
    petId: opponent.petId,
    playerName: opponent.DisplayName, // ← API から取得
    pet: {
      id: opponent.petId,
      guildId: opponent.guildId,
      guildName: opponent.guildName,
      name: opponent.name,
      species: opponent.species,
      attribute: opponent.attribute,
      level: opponent.level,
      exp: 0,
      maxHp: opponent.maxHp,
      power: opponent.power,
      guard: opponent.guard,
      speed: opponent.speed,
      acquiredAt: "",
    },
  }));
}

export async function startPetBattle(petId: string, opponentPetId: string): Promise<BattleResult> {
  const result = await apiFetch<BattleResultApiResponse>(
    `/pets/${encodeURIComponent(petId)}/battle`,
    {
      method: "POST",
      body: JSON.stringify({ opponentPetId }),
    },
  );
  return {
    ...result,
    result: result.result === "loss" ? "lose" : result.result,
  };
}
