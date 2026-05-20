import { describe, expect, it } from "vitest";
import { createWarGuilds, rankWarGuilds } from "./WarMapData";
import type { Guild } from "../../features/guild/api";

describe("WarMapData", () => {
  it("API のギルド集計値を war map 表示用データへ変換する", () => {
    const guilds: Guild[] = [
      {
        id: "guild_go",
        slug: "go",
        name: "Go",
        description: "Go guild",
        icon: "GO",
        color: "#00acd7",
        member_count: 3,
        total_contributed_cp: 120,
      },
      {
        id: "guild_unknown",
        slug: "unknown",
        name: "Unknown",
        description: "Unknown guild",
        icon: "??",
        color: "#ffffff",
        member_count: 99,
        total_contributed_cp: 999,
      },
    ];

    const warGuilds = createWarGuilds(guilds);

    expect(warGuilds).toHaveLength(1);
    expect(warGuilds[0]).toMatchObject({
      id: "go",
      guildID: "guild_go",
      name: "Go",
      totalCp: 120,
      memberCount: 3,
    });
  });

  it("総 CP の降順でランキングする", () => {
    const rankedGuilds = rankWarGuilds(
      createWarGuilds([
        {
          id: "guild_go",
          slug: "go",
          name: "Go",
          description: "Go guild",
          icon: "GO",
          color: "#00acd7",
          member_count: 3,
          total_contributed_cp: 120,
        },
        {
          id: "guild_rust",
          slug: "rust",
          name: "Rust",
          description: "Rust guild",
          icon: "RS",
          color: "#ff6b35",
          member_count: 2,
          total_contributed_cp: 240,
        },
      ]),
    );

    expect(rankedGuilds.map((guild) => guild.id)).toEqual(["rust", "go"]);
  });
});
