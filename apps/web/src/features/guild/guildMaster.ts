export interface GuildMaster {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  accent: string;
  sortOrder: number;
  memberCount: number;
  previousSeasonCp: number;
  oath: string;
  guildExperience?: number;
  guildLevel?: number;
  currentGuildLevelExperience?: number | null;
  nextGuildLevelExperience?: number | null;
  isMaxLevel?: boolean;
}

export const GUILD_MASTERS: GuildMaster[] = [
  {
    id: "guild_rust",
    slug: "rust",
    name: "Rust",
    description:
      "Memory safety is our banner. We build resilient systems, fearless refactors, and tools that refuse to break under pressure.",
    icon: "RS",
    color: "#ff6b35",
    accent: "#ffd0a8",
    sortOrder: 0,
    memberCount: 128,
    previousSeasonCp: 98420,
    oath: "Zero-cost courage",
  },
  {
    id: "guild_python",
    slug: "python",
    name: "Python",
    description:
      "We turn ideas into living prototypes, data pipelines, and automation spells before the enemy has finished compiling.",
    icon: "PY",
    color: "#3776ab",
    accent: "#ffd43b",
    sortOrder: 1,
    memberCount: 214,
    previousSeasonCp: 121760,
    oath: "Readable magic",
  },
  {
    id: "guild_go",
    slug: "go",
    name: "Go",
    description:
      "Small tools, steady services, clear concurrency. We value speed, calm operations, and code that new allies can read at a glance.",
    icon: "GO",
    color: "#00acd7",
    accent: "#b8f7ff",
    sortOrder: 2,
    memberCount: 156,
    previousSeasonCp: 110340,
    oath: "Ship the simple thing",
  },
  {
    id: "guild_typescript",
    slug: "typescript",
    name: "TypeScript",
    description:
      "We guard the frontier between product and platform with typed contracts, sharp UI craft, and refactors that keep moving.",
    icon: "TS",
    color: "#3178c6",
    accent: "#9fd7ff",
    sortOrder: 3,
    memberCount: 189,
    previousSeasonCp: 118900,
    oath: "Types before chaos",
  },
  {
    id: "guild_java",
    slug: "java",
    name: "Java",
    description:
      "Enterprise citadels, durable APIs, and battle-tested runtimes. We hold the long road and keep the old engines roaring.",
    icon: "JV",
    color: "#f89820",
    accent: "#ffe0a6",
    sortOrder: 4,
    memberCount: 97,
    previousSeasonCp: 74220,
    oath: "Stand the test of time",
  },
  {
    id: "guild_haskell",
    slug: "haskell",
    name: "Haskell",
    description:
      "We seek elegant laws in the noise: pure functions, precise abstractions, and strange machines that become obvious after midnight.",
    icon: "λ",
    color: "#8f6bd8",
    accent: "#e1d2ff",
    sortOrder: 5,
    memberCount: 64,
    previousSeasonCp: 58880,
    oath: "Make illegal states vanish",
  },
  {
    id: "guild_zig",
    slug: "zig",
    name: "Zig",
    description:
      "We command the metal directly: explicit control, fearless tooling, and tiny binaries forged for the next generation of systems.",
    icon: "ZG",
    color: "#f7a41d",
    accent: "#fff0b5",
    sortOrder: 6,
    memberCount: 72,
    previousSeasonCp: 61740,
    oath: "No hidden control flow",
  },
];

export function findGuildBySlug(slug: string | null): GuildMaster | null {
  if (!slug) return null;
  return GUILD_MASTERS.find((guild) => guild.slug === slug) ?? null;
}
