import type { GuildChatMessage } from "../api";

export function getChatMessageAuthorLabel(message: GuildChatMessage): string {
  return message.user_name ?? message.user_id;
}
