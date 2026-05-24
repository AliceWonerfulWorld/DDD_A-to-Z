const homePetStorageKey = "lang-war.pet.home-pet-id";

export function saveHomePetId(petId: string): void {
  try {
    window.localStorage.setItem(homePetStorageKey, petId);
  } catch {
    // Home can still fall back to the current guild pet.
  }
}

export function readHomePetId(): string | null {
  try {
    return window.localStorage.getItem(homePetStorageKey);
  } catch {
    return null;
  }
}
