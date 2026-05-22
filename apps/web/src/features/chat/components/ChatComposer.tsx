import { motion } from "framer-motion";
import type { Channel } from "phoenix";
import { useRef, useState } from "react";

interface ChatComposerProps {
  channel: Channel | null;
  placeholder?: string;
}

export function ChatComposer({
  channel,
  placeholder = "broadcast your next move...",
}: ChatComposerProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [sendError, setSendError] = useState(false);

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const input = inputRef.current;
    if (!input || !channel) {
      return;
    }
    const body = input.value.trim();
    if (!body) {
      return;
    }
    setSendError(false);
    const savedBody = body;
    channel.push("new_message", { body }).receive("error", () => {
      setSendError(true);
      if (input) {
        input.value = savedBody;
      }
    });
    input.value = "";
  };

  return (
    <form
      onSubmit={handleSubmit}
      style={{
        display: "grid",
        gridTemplateColumns: "minmax(0, 1fr) auto",
        gap: "10px",
        borderTop: "1px solid rgba(0, 245, 255, 0.18)",
        paddingTop: "12px",
      }}
    >
      {sendError && (
        <p
          role="alert"
          style={{
            gridColumn: "1 / -1",
            margin: 0,
            color: "#ff6b6b",
            fontSize: "0.68rem",
            fontFamily: '"DotGothic16", monospace',
          }}
        >
          SEND FAILED — RETRY
        </p>
      )}
      <input
        ref={inputRef}
        type="text"
        name="guild-chat"
        placeholder={placeholder}
        autoComplete="off"
        disabled={channel == null}
        style={{
          width: "100%",
          minHeight: "42px",
          border: "1px solid rgba(0, 245, 255, 0.34)",
          background: "rgba(0, 8, 20, 0.72)",
          color: "#f4ecd0",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.76rem",
          padding: "0 12px",
        }}
      />
      <motion.button
        type="submit"
        disabled={channel == null}
        whileHover={{ y: -1, scale: 1.02 }}
        whileTap={{ y: 1, scale: 0.98 }}
        style={{
          minHeight: "42px",
          border: "2px solid rgba(0, 245, 255, 0.68)",
          borderBottomColor: "rgba(2, 54, 72, 0.96)",
          borderRightColor: "rgba(2, 54, 72, 0.96)",
          background: "rgba(3, 12, 24, 0.84)",
          color: "#d9fbff",
          cursor: channel != null ? "pointer" : "not-allowed",
          fontFamily: "inherit",
          fontSize: "0.58rem",
          lineHeight: 1,
          padding: "0 12px",
        }}
      >
        SEND
      </motion.button>
    </form>
  );
}
