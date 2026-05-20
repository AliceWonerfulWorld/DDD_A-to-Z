import { beforeEach, describe, expect, it, vi } from "vitest";
import { fetchGuildActivityLogs, fetchMyGuild, joinGuild, leaveGuild } from "./api";

function mockFetch(status: number, body: unknown = null) {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(body),
    }),
  );
}

beforeEach(() => {
  vi.unstubAllGlobals();
});

describe("guild api", () => {
  it("fetchMyGuild は現在の所属ギルドを取得する", async () => {
    mockFetch(200, {
      guild: {
        id: "guild_go",
        slug: "go",
        name: "Go",
        description: "Go guild",
        icon: "GO",
        color: "#00acd7",
        member_count: 3,
        total_contributed_cp: 120,
      },
      membership: {
        id: "guild_membership_1",
        user_id: "user_1",
        joined_at: "2026-05-18T00:00:00Z",
      },
      members: [
        {
          user_id: "user_1",
          name: "Alice",
          total_earned_cp: 80,
          joined_at: "2026-05-18T00:00:00Z",
        },
      ],
    });

    const result = await fetchMyGuild();

    expect(result?.guild?.id).toBe("guild_go");
    expect(result?.members?.[0]?.total_earned_cp).toBe(80);
    expect(fetch).toHaveBeenCalledWith(
      "/api/me/guild",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("joinGuild は指定ギルドの参加 API を呼ぶ", async () => {
    mockFetch(201, { guild: null });

    await joinGuild("guild_go");

    expect(fetch).toHaveBeenCalledWith(
      "/api/guilds/guild_go/join",
      expect.objectContaining({ method: "POST" }),
    );
  });

  it("leaveGuild は所属ギルド脱退 API を呼ぶ", async () => {
    mockFetch(204);

    await leaveGuild();

    expect(fetch).toHaveBeenCalledWith(
      "/api/me/guild",
      expect.objectContaining({ method: "DELETE" }),
    );
  });

  it("fetchGuildActivityLogs は指定ギルドのアクティビティログを取得する", async () => {
    mockFetch(200, {
      logs: [
        {
          id: "user_1:commit:repo:sha",
          user_id: "user_1",
          player: "Alice",
          type: "commit",
          repo: "jyogi-web/DDD_A-to-Z",
          message: "Add activity logs",
          language: "Go",
          cp: 1,
          occurred_at: "2026-05-20T00:00:00Z",
        },
      ],
    });

    const logs = await fetchGuildActivityLogs("guild/go", 20);

    expect(logs[0]?.message).toBe("Add activity logs");
    expect(fetch).toHaveBeenCalledWith(
      "/api/guilds/guild%2Fgo/activity-logs?limit=20",
      expect.objectContaining({ credentials: "include" }),
    );
  });
});
