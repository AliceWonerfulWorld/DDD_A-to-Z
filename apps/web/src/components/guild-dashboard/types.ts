export type GuildTab = "activity" | "rankings";

export interface ActivityLog {
  id: string;
  player: string;
  action: string;
  cp: number;
  tone: string;
}

export interface RankingMember {
  name: string;
  title: string;
  cp: number;
  color: string;
}
