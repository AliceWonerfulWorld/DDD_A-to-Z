import { apiFetch } from "../../lib/api/client";

export interface MyPageUser {
  id: string;
  github_id: number;
  username: string;
  avatar_url: string;
  github_avatar_url: string;
  created_at: string;
}

export interface MyPageContributionPoints {
  balance: number;
  total_earned: number;
  total_spent: number;
}

export interface MyPageRecentRepo {
  github_id: number;
  full_name: string;
  language: string;
  html_url: string;
  pushed_at: string | null;
}

export interface MyPageRepositories {
  total_count: number;
  language_summary: Record<string, number>;
  recent: MyPageRecentRepo[];
}

export interface MyPageGuild {
  id: string;
  name: string;
  slug: string;
  icon: string;
  color: string;
  description: string;
  member_count: number;
  rank: number;
  total_guilds: number;
  cp: number;
}

export interface GitHubStats {
  total_stars: number;
  total_prs: number;
  total_issues: number;
  contributed_to: number;
  public_repos: number;
  github_created_at: string;
  yearly_commits: number;
  yearly_contributions: number;
}

export interface MyPageBadge {
  slug: string;
  name: string;
  description: string;
  icon: string;
  earned_at: string;
}

export interface MyPageResponse {
  user: MyPageUser;
  contribution_points: MyPageContributionPoints;
  repositories: MyPageRepositories;
  guild: MyPageGuild | null;
  github_stats: GitHubStats | null;
  badges: MyPageBadge[];
  selected_badge_slug: string | null;
}

export async function fetchMyPage(): Promise<MyPageResponse> {
  return apiFetch<MyPageResponse>("/mypage");
}
