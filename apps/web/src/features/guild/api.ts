import { ApiError, apiFetch } from "../../lib/api/client";

export interface Guild {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  member_count: number;
  total_contributed_cp: number;
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
  joined_at: string;
}

export interface GuildMembershipResponse {
  guild: Guild | null;
  membership?: GuildMembership;
  members?: GuildMemberContribution[];
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
