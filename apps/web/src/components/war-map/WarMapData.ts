export interface WarGuild {
  id: string;
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

export const CURRENT_GUILD_ID = "typescript";

export const WAR_GUILDS: WarGuild[] = [
  {
    id: "rust",
    name: "Rust",
    mark: "RS",
    color: "#ff6b35",
    accent: "#ffd0b8",
    x: 25,
    y: 38,
    totalCp: 528400,
    memberCount: 184,
    description: "Memory-safe front line with high defense and fearless refactoring.",
  },
  {
    id: "python",
    name: "Python",
    mark: "PY",
    color: "#ffd43b",
    accent: "#3776ab",
    x: 42,
    y: 56,
    totalCp: 613900,
    memberCount: 246,
    description: "A flexible guild that turns scripts, data, and automation into territory.",
  },
  {
    id: "go",
    name: "Go",
    mark: "GO",
    color: "#00add8",
    accent: "#b8f4ff",
    x: 62,
    y: 40,
    totalCp: 486200,
    memberCount: 172,
    description: "Fast deployment crews holding the network routes with calm concurrency.",
  },
  {
    id: "typescript",
    name: "TypeScript",
    mark: "TS",
    color: "#3178c6",
    accent: "#9ed1ff",
    x: 52,
    y: 68,
    totalCp: 641700,
    memberCount: 258,
    description: "Typed spellcasters building resilient interfaces across the capital grid.",
  },
  {
    id: "java",
    name: "Java",
    mark: "JV",
    color: "#f89820",
    accent: "#ffe0a8",
    x: 72,
    y: 58,
    totalCp: 459600,
    memberCount: 165,
    description: "Veteran fortress engineers with deep reserves and enterprise-scale stamina.",
  },
  {
    id: "haskell",
    name: "Haskell",
    mark: "HS",
    color: "#9b5de5",
    accent: "#e2c8ff",
    x: 36,
    y: 74,
    totalCp: 398300,
    memberCount: 96,
    description: "Abstract tacticians reshaping the battlefield with pure functions.",
  },
  {
    id: "zig",
    name: "Zig",
    mark: "ZG",
    color: "#f7a41d",
    accent: "#fff0b5",
    x: 78,
    y: 28,
    totalCp: 372800,
    memberCount: 88,
    description: "Low-level scouts optimizing supply lines close to the metal.",
  },
];

export const RANKED_WAR_GUILDS = [...WAR_GUILDS].sort((a, b) => b.totalCp - a.totalCp);
