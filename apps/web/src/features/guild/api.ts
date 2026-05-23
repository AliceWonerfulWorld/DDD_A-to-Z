import { ApiError, apiFetch } from "../../lib/api/client";
import type { GrantedPetAPIResponse } from "../pet/guildGrant";

export interface Guild {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  member_count: number;
  total_contributed_cp: number;
  currentExp?: number;
  current_exp?: number;
  guildLevel?: number;
  guild_experience?: number;
  guild_level?: number;
  isMaxLevel?: boolean;
  is_max_level?: boolean;
  current_guild_level_experience?: number | null;
  next_guild_level_experience?: number | null;
}

export interface GuildMembership {
  id: string;
  user_id?: string;
  joined_at: string;
}

export interface GuildMemberContribution {
  user_id: string;
  name: string;
  total_earned_cp: number;
  total_contributed_cp: number;
  joined_at: string;
}

export interface GuildActivityLog {
  id: string;
  user_id: string;
  player: string;
  type: "commit" | "pull_request";
  repo: string;
  message: string;
  language: string;
  cp: number;
  occurred_at: string;
}

export interface GuildMembershipResponse {
  guild: Guild | null;
  membership?: GuildMembership;
  members?: GuildMemberContribution[];
  granted_pet?: GrantedPetAPIResponse | null;
  grantedPet?: GrantedPetAPIResponse | null;
  pet_already_owned?: boolean;
  petAlreadyOwned?: boolean;
}

export async function fetchGuilds(): Promise<Guild[]> {
  const data = await apiFetch<{ guilds: Guild[] }>("/guilds");
  return data.guilds;
}

export async function fetchMyGuild(): Promise<GuildMembershipResponse | null> {
  try {
    return await apiFetch<GuildMembershipResponse>("/me/guild");
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }
    throw error;
  }
}

export async function joinGuild(guildID: string): Promise<GuildMembershipResponse> {
  return apiFetch<GuildMembershipResponse>(`/guilds/${guildID}/join`, { method: "POST" });
}

export async function leaveGuild(): Promise<void> {
  await apiFetch<void>("/me/guild", { method: "DELETE" });
}

export async function fetchGuildActivityLogs(
  guildID: string,
  limit = 20,
): Promise<GuildActivityLog[]> {
  const encodedGuildID = encodeURIComponent(guildID);
  const params = new URLSearchParams({ limit: String(limit) });
  const data = await apiFetch<{ logs: GuildActivityLog[] }>(
    `/guilds/${encodedGuildID}/activity-logs?${params.toString()}`,
  );
  return data.logs;
}
