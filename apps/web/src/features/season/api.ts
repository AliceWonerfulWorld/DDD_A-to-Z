import { apiFetch } from "../../lib/api/client";

export interface Season {
  id: string;
  number: number;
  starts_at: string;
  ends_at: string;
  is_current: boolean;
}

export interface GuildSeasonRanking {
  id: string;
  season_id: string;
  guild_id: string;
  total_cp: number;
  rank: number;
  member_count: number;
}

export interface GuildSeasonMemberRanking {
  id: string;
  season_id: string;
  guild_id: string;
  user_id: string;
  user_name: string;
  contributed_cp: number;
  rank: number;
}

export async function fetchSeasons(): Promise<Season[]> {
  const data = await apiFetch<{ seasons: Season[] }>("/seasons");
  return data.seasons;
}

export async function fetchCurrentSeason(): Promise<Season> {
  const data = await apiFetch<Season>("/seasons/current");
  return data;
}

export async function fetchSeasonByNumber(number: number): Promise<Season> {
  const data = await apiFetch<Season>(`/seasons/${number}`);
  return data;
}

export async function fetchGuildSeasonRankings(
  seasonNumber: number,
): Promise<GuildSeasonRanking[]> {
  const data = await apiFetch<{ rankings: GuildSeasonRanking[] }>(
    `/seasons/${seasonNumber}/guild-rankings`,
  );
  return data.rankings;
}

export async function fetchGuildSeasonMemberRankings(
  seasonNumber: number,
  guildID: string,
): Promise<GuildSeasonMemberRanking[]> {
  const encodedGuildID = encodeURIComponent(guildID);
  const data = await apiFetch<{ rankings: GuildSeasonMemberRanking[] }>(
    `/seasons/${seasonNumber}/guilds/${encodedGuildID}/member-rankings`,
  );
  return data.rankings;
}
