import { ApiError, apiFetch } from "../../lib/api/client";

export interface Guild {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  member_count: number;
}

export interface GuildMembership {
  id: string;
  joined_at: string;
}

export interface GuildMembershipResponse {
  guild: Guild | null;
  membership?: GuildMembership;
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
