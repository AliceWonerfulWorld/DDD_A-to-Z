import type { Guild } from "./api";
import { findGuildBySlug, GUILD_MASTERS, type GuildMaster } from "./guildMaster";

const fallbackGuildMeta = {
  accent: "#fff8d7",
  previousSeasonCp: 0,
  oath: "New banner raised",
};

export type DisplayGuild = GuildMaster;

export function toDisplayGuild(guild: Guild): DisplayGuild {
  const master = findGuildBySlug(guild.slug);

  return {
    id: guild.id,
    slug: guild.slug,
    name: guild.name,
    description: guild.description,
    icon: guild.icon,
    color: guild.color,
    accent: master?.accent ?? fallbackGuildMeta.accent,
    sortOrder: master?.sortOrder ?? GUILD_MASTERS.length,
    memberCount: guild.member_count,
    previousSeasonCp: master?.previousSeasonCp ?? fallbackGuildMeta.previousSeasonCp,
    oath: master?.oath ?? fallbackGuildMeta.oath,
  };
}

export function toDisplayGuilds(guilds: Guild[]): DisplayGuild[] {
  return guilds.map(toDisplayGuild).sort((left, right) => left.sortOrder - right.sortOrder);
}
