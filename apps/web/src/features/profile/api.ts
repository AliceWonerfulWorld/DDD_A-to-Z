import { apiFetch } from "../../lib/api/client";

export async function completeInitialProfileAPI(displayName: string): Promise<void> {
  await apiFetch<void>("/profile/complete", {
    method: "POST",
    body: JSON.stringify({ display_name: displayName }),
  });
}

export type Profile = {
  display_name: string;
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
