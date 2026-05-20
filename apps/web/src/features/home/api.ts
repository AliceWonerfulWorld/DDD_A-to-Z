import { apiFetch } from "../../lib/api/client";

export interface HomeCPData {
  total_cp: number;
  today_cp: number;
  player_level?: number;
  player_level_total_cp?: number;
  next_player_level?: number;
  next_player_level_total_cp?: number;
  next_player_level_remaining?: number;
  lifetime_total_earned_cp?: number;
}

export async function fetchHomeCP(): Promise<HomeCPData> {
  return apiFetch<HomeCPData>("/home");
}
