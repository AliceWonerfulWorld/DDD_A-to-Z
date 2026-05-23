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
  selected_badge_slug: string | null;
  avatar_url?: string;
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
