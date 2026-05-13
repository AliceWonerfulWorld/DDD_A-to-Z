import { apiFetch, ApiError } from "../../lib/api/client";
import type { CurrentUser } from "./types";

type MeResponse = { user: CurrentUser };

// 未ログインなら null を返す
export async function fetchMe(): Promise<CurrentUser | null> {
  try {
    const data = await apiFetch<MeResponse>("/me");
    return data.user;
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) return null;
    throw e;
  }
}

export function beginLogin(): void {
  window.location.href = "/api/auth/github/login";
}

export async function logout(): Promise<void> {
  await apiFetch<void>("/auth/logout", { method: "POST" });
}
