import { apiFetch } from "../../lib/api/client";

export async function completeInitialProfileAPI(
  displayName: string,
  avatarUrl: string,
): Promise<void> {
  await apiFetch<void>("/profile/complete", {
    method: "POST",
    body: JSON.stringify({ display_name: displayName, avatar_url: avatarUrl }),
  });
}

export async function updateProfileAPI(displayName: string, avatarUrl: string): Promise<void> {
  await apiFetch<void>("/profile", {
    method: "PUT",
    body: JSON.stringify({ display_name: displayName, avatar_url: avatarUrl }),
  });
}

export type Profile = {
  display_name: string;
  avatar_url?: string;
};

export async function fetchProfile(): Promise<Profile | null> {
  try {
    const data = await apiFetch<Profile>("/profile");
    return data;
  } catch {
    // 401 や 404 の場合は null を返す
    return null;
  }
}
