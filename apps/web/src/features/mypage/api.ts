import { apiFetch } from "../../lib/api/client";

export interface MyPageUser {
  id: string;
  github_id: number;
  username: string;
  avatar_url: string;
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
}

export interface MyPageResponse {
  user: MyPageUser;
  contribution_points: MyPageContributionPoints;
  repositories: MyPageRepositories;
  guild: MyPageGuild | null;
}

export async function fetchMyPage(): Promise<MyPageResponse> {
  return apiFetch<MyPageResponse>("/mypage");
}
