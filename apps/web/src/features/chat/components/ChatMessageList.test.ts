import { describe, expect, test } from "vitest";
import { getChatMessageAuthorLabel } from "./authorLabel";

describe("getChatMessageAuthorLabel", () => {
  test("user_name がある場合は user_id ではなく user_name を返す", () => {
    expect(
      getChatMessageAuthorLabel({
        id: "msg_1",
        guild_id: "guild_rust",
        user_id: "user_1",
        user_name: "Octo Mage",
        body: "hello",
        created_at: "2026-05-22T00:00:00Z",
      }),
    ).toBe("Octo Mage");
  });

  test("user_name がない場合は user_id に fallback する", () => {
    expect(
      getChatMessageAuthorLabel({
        id: "msg_1",
        guild_id: "guild_rust",
        user_id: "user_1",
        body: "hello",
        created_at: "2026-05-22T00:00:00Z",
      }),
    ).toBe("user_1");
  });
});
