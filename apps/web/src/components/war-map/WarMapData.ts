import type { Guild } from "../../features/guild/api";

export interface WarGuild {
  id: string;
  guildID: string;
  name: string;
  mark: string;
  color: string;
  accent: string;
  x: number;
  y: number;
  totalCp: number;
  memberCount: number;
  description: string;
}

interface WarGuildMapPosition {
  slug: string;
  mark: string;
  accent: string;
  x: number;
  y: number;
}

const WAR_GUILD_POSITIONS: WarGuildMapPosition[] = [
  {
    mark: "RS",
    accent: "#ffd0b8",
    x: 25,
    y: 38,
    slug: "rust",
  },
  {
    mark: "PY",
    accent: "#3776ab",
    x: 42,
    y: 56,
    slug: "python",
  },
  {
    mark: "GO",
    accent: "#b8f4ff",
    x: 62,
    y: 40,
    slug: "go",
  },
  {
    mark: "TS",
    accent: "#9ed1ff",
    x: 52,
    y: 68,
    slug: "typescript",
  },
  {
    mark: "JV",
    accent: "#ffe0a8",
    x: 72,
    y: 58,
    slug: "java",
  },
  {
    mark: "HS",
    accent: "#e2c8ff",
    x: 36,
    y: 74,
    slug: "haskell",
  },
  {
    mark: "ZG",
    accent: "#fff0b5",
    x: 78,
    y: 28,
    slug: "zig",
  },
];

const WAR_GUILD_POSITION_BY_SLUG = new Map(
  WAR_GUILD_POSITIONS.map((position) => [position.slug, position]),
);

export function createWarGuilds(guilds: Guild[]): WarGuild[] {
  return guilds
    .map((guild) => {
      const position = WAR_GUILD_POSITION_BY_SLUG.get(guild.slug);
      if (!position) {
        return null;
      }

      return {
        id: guild.slug,
        guildID: guild.id,
        name: guild.name,
        mark: position.mark,
        color: guild.color,
        accent: position.accent,
        x: position.x,
        y: position.y,
        totalCp: guild.total_contributed_cp,
        memberCount: guild.member_count,
        description: guild.description,
      };
    })
    .filter((guild): guild is WarGuild => guild !== null);
}

export function rankWarGuilds(guilds: WarGuild[]): WarGuild[] {
  return [...guilds].sort((a, b) => {
    if (b.totalCp !== a.totalCp) {
      return b.totalCp - a.totalCp;
    }

    return a.name.localeCompare(b.name);
  });
}

export function findWarGuildByID(
  guilds: WarGuild[],
  guildID: string | null | undefined,
): WarGuild | null {
  if (!guildID) {
    return null;
  }

  return guilds.find((guild) => guild.id === guildID) ?? null;
}
