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
  return apiFetch<TrainingResult>(`/pets/${encodeURIComponent(petId)}/training`, {
    method: "POST",
    body: JSON.stringify({ stat }),
  });
}

export async function fetchBattleOpponents(): Promise<BattleOpponent[]> {
  const data = await apiFetch<{ opponents: BattleOpponent[] }>("/pets/battle-opponents");
  return data.opponents;
}

export async function startPetBattle(opponentUserId: string): Promise<BattleResult> {
  return apiFetch<BattleResult>("/pets/battles", {
    method: "POST",
    body: JSON.stringify({ opponentUserId }),
  });
}
