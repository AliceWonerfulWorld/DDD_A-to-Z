import type { ActivityLog, GuildTab, RankingMember } from "./types";

const PLAYERS = [
  "AkiByte",
  "NullMage",
  "PixelNinja",
  "LoopKnight",
  "TypeSmith",
  "AsyncRogue",
  "BugSlayer",
  "CacheWizard",
];

const LOG_ACTIONS = [
  { action: "merged a pull request", cp: 420, tone: "#74f7a1" },
  { action: "reviewed 3 files", cp: 180, tone: "#9be7ff" },
  { action: "fixed a flaky test", cp: 260, tone: "#ffd966" },
  { action: "shipped guild docs", cp: 150, tone: "#d9b8ff" },
  { action: "closed a critical issue", cp: 540, tone: "#ff9b9b" },
  { action: "refactored a use case", cp: 310, tone: "#74f7a1" },
];

export const INITIAL_LOGS: ActivityLog[] = [
  {
    id: 1005,
    player: "TypeSmith",
    action: "merged a pull request",
    cp: 420,
    tone: "#74f7a1",
  },
  {
    id: 1004,
    player: "PixelNinja",
    action: "fixed a flaky test",
    cp: 260,
    tone: "#ffd966",
  },
  {
    id: 1003,
    player: "LoopKnight",
    action: "reviewed 3 files",
    cp: 180,
    tone: "#9be7ff",
  },
  {
    id: 1002,
    player: "NullMage",
    action: "closed a critical issue",
    cp: 540,
    tone: "#ff9b9b",
  },
  {
    id: 1001,
    player: "AkiByte",
    action: "refactored a use case",
    cp: 310,
    tone: "#74f7a1",
  },
];

export const RANKINGS: RankingMember[] = [
  { name: "TypeSmith", title: "Generic Hero", cp: 35420, color: "#ffd966" },
  { name: "NullMage", title: "Void Debugger", cp: 31980, color: "#d9b8ff" },
  { name: "PixelNinja", title: "UI Shinobi", cp: 28640, color: "#9be7ff" },
  { name: "LoopKnight", title: "Iteration Paladin", cp: 25110, color: "#74f7a1" },
  { name: "AsyncRogue", title: "Promise Runner", cp: 22470, color: "#ff9b9b" },
  { name: "CacheWizard", title: "Memo Sage", cp: 19860, color: "#f4ecd0" },
  { name: "BugSlayer", title: "Regression Breaker", cp: 17420, color: "#f6a6ff" },
];

export const GUILD_TABS: { id: GuildTab; label: string }[] = [
  { id: "activity", label: "ACTIVITY LOG" },
  { id: "rankings", label: "RANKINGS" },
];

export function createLog(id: number): ActivityLog {
  const player = PLAYERS[Math.floor(Math.random() * PLAYERS.length)];
  const log = LOG_ACTIONS[Math.floor(Math.random() * LOG_ACTIONS.length)];

  return {
    id,
    player,
    action: log.action,
    cp: log.cp,
    tone: log.tone,
  };
}
