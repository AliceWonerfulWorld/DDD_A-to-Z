import type { GrantedPet } from "./api";

const grantedPetStorageKey = "lang-war.pet.last-granted";

interface GrantedPetAPIResponseSnake {
  id: string;
  guild_id: string;
  attribute: string;
  created_at: string;
}

interface GrantedPetAPIResponseCamel {
  id: string;
  guildId: string;
  attribute: string;
  createdAt: string;
}

export type GrantedPetAPIResponse = GrantedPetAPIResponseSnake | GrantedPetAPIResponseCamel;

export function normalizeGrantedPet(pet: GrantedPetAPIResponse): GrantedPet {
  if ("guild_id" in pet) {
    return {
      id: pet.id,
      guildId: pet.guild_id,
      attribute: pet.attribute,
      createdAt: pet.created_at,
    };
  }

  return pet;
}

export function storeGrantedPet(pet: GrantedPet): void {
  try {
    window.sessionStorage.setItem(grantedPetStorageKey, JSON.stringify(pet));
  } catch {
    // The next screen can still render without the one-shot celebration.
  }
}

export function consumeGrantedPet(): GrantedPet | null {
  try {
    const raw = window.sessionStorage.getItem(grantedPetStorageKey);
    if (!raw) return null;
    window.sessionStorage.removeItem(grantedPetStorageKey);
    return JSON.parse(raw) as GrantedPet;
  } catch {
    return null;
  }
}
