import type { GuildChatMessage } from "./api";

export const GUILD_CHAT_MESSAGES: GuildChatMessage[] = [
  {
    id: "guild-msg-1",
    guild_id: "guild_mock",
    user_id: "user_null_mage",
    user_name: "NullMage",
    body: "西門ルートの監視、こっちで継続中。",
    created_at: "2026-05-22T22:14:00Z",
  },
  {
    id: "guild-msg-2",
    guild_id: "guild_mock",
    user_id: "user_pixel_ninja",
    user_name: "PixelNinja",
    body: "次の演出差し替え、3 分で反映いける。",
    created_at: "2026-05-22T22:15:00Z",
  },
  {
    id: "guild-msg-3",
    guild_id: "guild_mock",
    user_id: "user_type_smith",
    user_name: "TypeSmith",
    body: "了解、UI の更新後に activity log と同期する。",
    created_at: "2026-05-22T22:16:00Z",
  },
  {
    id: "guild-msg-4",
    guild_id: "guild_mock",
    user_id: "user_loop_knight",
    user_name: "LoopKnight",
    body: "南側ルートに review 待ちが残ってる。終わり次第こっち戻る。",
    created_at: "2026-05-22T22:18:00Z",
  },
  {
    id: "guild-msg-5",
    guild_id: "guild_mock",
    user_id: "user_aki_byte",
    user_name: "AkiByte",
    body: "最新の CI ログ、警告だけだった。デプロイ進めてよさそう。",
    created_at: "2026-05-22T22:21:00Z",
  },
  {
    id: "guild-msg-6",
    guild_id: "guild_mock",
    user_id: "user_type_smith",
    user_name: "TypeSmith",
    body: "ありがとう。次は guild details 側の導線もまとめて見る。",
    created_at: "2026-05-22T22:23:00Z",
  },
];
