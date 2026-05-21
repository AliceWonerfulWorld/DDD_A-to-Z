import { Channel, Socket } from "phoenix";
import { apiFetch } from "../../lib/api/client";

export interface GuildChatMessage {
  id: string;
  guild_id: string;
  user_id: string;
  body: string;
  created_at: string;
}

interface ChatTokenResponse {
  token: string;
  expires_at: string;
}

export async function fetchChatToken(guildID: string): Promise<ChatTokenResponse> {
  return apiFetch<ChatTokenResponse>(`/guilds/${guildID}/chat-token`, { method: "POST" });
}

export interface ChatConnection {
  socket: Socket;
  channel: Channel;
  disconnect: () => void;
}

export async function connectChat(
  guildID: string,
  chatServiceUrl: string,
  onMessage: (msg: GuildChatMessage) => void,
  onHistory: (msgs: GuildChatMessage[]) => void,
  onError: (reason: string) => void,
): Promise<ChatConnection> {
  const { token } = await fetchChatToken(guildID);

  const socket = new Socket(`${chatServiceUrl}/socket`, { params: { token } });
  socket.connect();

  const channel = socket.channel(`guild:${guildID}`, {});

  channel
    .join()
    .receive("ok", (resp: unknown) => {
      const r = resp as { messages: GuildChatMessage[] };
      onHistory(r.messages ?? []);
    })
    .receive("error", (resp: unknown) => {
      const r = resp as { reason?: string };
      onError(r.reason ?? "join failed");
    });

  channel.on("new_message", (msg: unknown) => {
    onMessage(msg as GuildChatMessage);
  });

  const disconnect = () => {
    channel.leave();
    socket.disconnect();
  };

  return { socket, channel, disconnect };
}
