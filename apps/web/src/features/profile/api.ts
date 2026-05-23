import { apiFetch } from "../../lib/api/client";

export async function completeInitialProfileAPI(displayName: string): Promise<void> {
  await apiFetch<void>("/profile/complete", {
    method: "POST",
    body: JSON.stringify({ display_name: displayName }),
  });
}

export type Profile = {
  display_name: string;
  selected_badge_slug: string | null;
};

export async function fetchProfile(): Promise<Profile | null> {
  try {
    const data = await apiFetch<Profile>("/profile");
    return data;
  } catch {
    return null;
  }
}

export async function updateSelectedBadge(badgeSlug: string | null): Promise<void> {
  await apiFetch<void>("/profile/badge", {
    method: "PATCH",
    body: JSON.stringify({ badge_slug: badgeSlug }),
  });
}
