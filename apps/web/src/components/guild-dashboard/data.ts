import type { ActivityLog, GuildTab, RankingMember } from "./types";

export const INITIAL_LOGS: ActivityLog[] = [
  {
    id: "mock-1005",
    player: "UO!",
    action:
      "Commit: feat: ログインユーザーのSP取得API(GET /me/sp)を追加 (#141) * feat: SPシステムを実装し、point_typesをcode+language複合PKに再設計 (#101) - point_typesを(code, language)複合PKに変更し、40以上の言語SPをINSERTのみで管理可能にする - PointType型をstringから{Code, Language}構造体に変更し、SPType(language)コンストラクタを追加 - リポジトリ分析でCP付与と同時に言語別SPを付与するSPEarnerインターフェースを実装 - point_types未登録言語はErrUnsupportedPointTypeで静かにスキップする",
    cp: 840,
    tone: "#ffd966",
  },
  {
    id: "mock-1004",
    player: "TypeSmith",
    action:
      "PR: feat: GuildDashboardにリアルタイムアクティビティストリームを追加 - 10秒ポーリング・AnimatePresenceアニメーション・ActivityLogPanelコンポーネント実装 (#133)",
    cp: 560,
    tone: "#74f7a1",
  },
  {
    id: "mock-1003",
    player: "PixelNinja",
    action:
      "Commit: fix: ブラウザズーム時にログ行が重なるバグを修正 - alignItemsをbaselineからstartに変更、グリッドセルにminWidth:0を追加 (#144)",
    cp: 230,
    tone: "#ff9b9b",
  },
  {
    id: "mock-1002",
    player: "NullMage",
    action:
      "PR: refactor: NewUseCaseでspフィールドを追加し、言語貢献エンティティにSPBalanceを結合、フロントエンド表示コンポーネントを更新 (#138)",
    cp: 620,
    tone: "#74f7a1",
  },
  {
    id: "mock-1001",
    player: "LoopKnight",
    action:
      "Commit: chore: buf.shのパーミッション修正とprotobufスキーマをモジュール化 - generate.shとlint.shに分割してCIワークフローを簡潔化 (#120)",
    cp: 160,
    tone: "#9be7ff",
  },
];

export const RANKINGS: RankingMember[] = [
  { name: "TypeSmith", title: "Generic Hero", cp: 35420, color: "#ffd966" },
  { name: "NullMage", title: "Void Debugger", cp: 31980, color: "#d9b8ff" },
  { name: "PixelNinja", title: "UI Shinobi", cp: 28640, color: "#9be7ff" },
  {
    name: "LoopKnight",
    title: "Iteration Paladin",
    cp: 25110,
    color: "#74f7a1",
  },
  { name: "AsyncRogue", title: "Promise Runner", cp: 22470, color: "#ff9b9b" },
  { name: "CacheWizard", title: "Memo Sage", cp: 19860, color: "#f4ecd0" },
  {
    name: "BugSlayer",
    title: "Regression Breaker",
    cp: 17420,
    color: "#f6a6ff",
  },
];

export const GUILD_TABS: { id: GuildTab; label: string }[] = [
  { id: "activity", label: "ACTIVITY LOG" },
  { id: "rankings", label: "RANKINGS" },
  { id: "season", label: "SEASON" },
];

const PLAYERS = [
  "UO!",
  "TypeSmith",
  "PixelNinja",
  "NullMage",
  "LoopKnight",
  "AsyncRogue",
  "CacheWizard",
  "BugSlayer",
];

const LOG_ACTIONS = [
  {
    action: "Commit: 新機能を実装しました",
    cp: 100,
    tone: "#ffd966",
  },
  {
    action: "PR: 機能を改善しました",
    cp: 150,
    tone: "#74f7a1",
  },
  {
    action: "Commit: バグを修正しました",
    cp: 80,
    tone: "#ff9b9b",
  },
];

export function createLog(id: number): ActivityLog {
  const player = PLAYERS[Math.floor(Math.random() * PLAYERS.length)];
  const log = LOG_ACTIONS[Math.floor(Math.random() * LOG_ACTIONS.length)];

  return {
    id: `mock-${id}`,
    player,
    action: log.action,
    cp: log.cp,
    tone: log.tone,
  };
}
