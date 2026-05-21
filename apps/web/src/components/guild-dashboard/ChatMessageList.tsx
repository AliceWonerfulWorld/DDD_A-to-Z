import type { GuildChatMessage } from "../../features/guild/chatApi";

interface ChatMessageListProps {
  messages: GuildChatMessage[];
  currentUserID?: string;
  dense?: boolean;
}

function formatTime(isoString: string): string {
  const date = new Date(isoString);
  return `${String(date.getHours()).padStart(2, "0")}:${String(date.getMinutes()).padStart(2, "0")}`;
}

export function ChatMessageList({ messages, currentUserID, dense = false }: ChatMessageListProps) {
  return (
    <div
      style={{
        minHeight: 0,
        display: "grid",
        alignContent: "start",
        gap: dense ? "10px" : "12px",
        overflow: "auto",
        paddingRight: "4px",
      }}
    >
      {messages.map((message) => {
        const isSelf = currentUserID != null && message.user_id === currentUserID;
        return (
          <div
            key={message.id}
            style={{
              justifySelf: isSelf ? "end" : "start",
              maxWidth: dense ? "88%" : "82%",
              border: `1px solid ${isSelf ? "rgba(255,217,102,0.46)" : "rgba(0,245,255,0.28)"}`,
              background: isSelf
                ? "linear-gradient(180deg, rgba(44,34,12,0.84), rgba(21,14,4,0.74))"
                : "linear-gradient(180deg, rgba(8,25,42,0.82), rgba(1,9,22,0.72))",
              boxShadow: "inset 0 0 14px rgba(0,245,255,0.08)",
              padding: dense ? "10px 12px" : "12px 14px",
            }}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                gap: "10px",
                color: isSelf ? "#ffd966" : "#9be7ff",
                fontFamily: '"DotGothic16", monospace',
                fontSize: dense ? "0.62rem" : "0.68rem",
                lineHeight: 1.4,
              }}
            >
              <span>{message.user_id}</span>
              <span>{formatTime(message.created_at)}</span>
            </div>
            <p
              style={{
                margin: "8px 0 0",
                fontFamily: '"DotGothic16", monospace',
                fontSize: dense ? "0.82rem" : "0.9rem",
                lineHeight: dense ? 1.7 : 1.85,
              }}
            >
              {message.body}
            </p>
          </div>
        );
      })}
    </div>
  );
}
