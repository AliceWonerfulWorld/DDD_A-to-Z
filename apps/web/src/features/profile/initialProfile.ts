const completedProfileUserIDsKey = "lang-war.initial-profile.completed-user-ids";

function readCompletedUserIDs(): string[] {
  try {
    const value = window.localStorage.getItem(completedProfileUserIDsKey);
    if (!value) return [];

    const parsed = JSON.parse(value);
    if (!Array.isArray(parsed)) return [];

    return parsed.filter((id): id is string => typeof id === "string" && id.length > 0);
  } catch {
    return [];
  }
}

export function hasCompletedInitialProfile(userID: string): boolean {
  return readCompletedUserIDs().includes(userID);
}

export function markInitialProfileCompleted(userID: string): void {
  const completedUserIDs = new Set(readCompletedUserIDs());
  completedUserIDs.add(userID);
  try {
    window.localStorage.setItem(completedProfileUserIDsKey, JSON.stringify([...completedUserIDs]));
  } catch {
    // Ignore storage failures so profile completion can continue.
  }
}
