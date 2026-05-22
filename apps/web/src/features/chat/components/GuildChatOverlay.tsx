import { AnimatePresence, motion } from "framer-motion";
import type { Channel } from "phoenix";
import type { GuildChatMessage } from "../api";
import { steppedEase } from "../../../lib/animationUtils";
import { ChatComposer } from "./ChatComposer";
import { ChatMessageList } from "./ChatMessageList";

interface GuildChatOverlayProps {
  isOpen: boolean;
  messages: GuildChatMessage[];
  channel: Channel | null;
  onExpand: () => void;
  onClose: () => void;
}

export function GuildChatOverlay({
  isOpen,
  messages,
  channel,
  onExpand,
  onClose,
}: GuildChatOverlayProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <motion.aside
          key="guild-chat-overlay"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.24, ease: steppedEase(5) }}
          aria-labelledby="guild-chat-overlay-title"
          style={{
            position: "fixed",
            top: "calc(env(safe-area-inset-top, 0px) + clamp(136px, 12vw, 164px))",
            right: "clamp(16px, 2.4vw, 32px)",
            zIndex: 4,
            width: "min(calc(100vw - 32px), 360px)",
            minHeight: "420px",
            display: "grid",
            gridTemplateRows: "auto minmax(0, 1fr) auto",
            gap: "12px",
            border: "2px solid rgba(0, 245, 255, 0.74)",
            borderLeftColor: "rgba(0, 245, 255, 0.92)",
            background:
              "linear-gradient(180deg, rgba(3, 10, 24, 0.76), rgba(1, 8, 22, 0.88)), repeating-linear-gradient(180deg, rgba(130,234,255,0.05), rgba(130,234,255,0.05) 1px, transparent 1px, transparent 4px)",
            boxShadow:
              "-5px 0 15px rgba(0,243,255,0.3), 0 0 0 2px rgba(0,0,0,0.64), inset 0 0 24px rgba(0,245,255,0.09)",
            color: "#f4ecd0",
            backdropFilter: "blur(10px)",
            padding: "16px",
            overflow: "hidden",
          }}
        >
          <header
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "start",
              gap: "12px",
              borderBottom: "1px solid rgba(0, 245, 255, 0.2)",
              paddingBottom: "12px",
            }}
          >
            <div>
              <span
                style={{
                  display: "block",
                  marginBottom: "8px",
                  color: "rgba(155, 231, 255, 0.72)",
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.64rem",
                  lineHeight: 1.4,
                }}
              >
                COMM RELAY
              </span>
              <h2
                id="guild-chat-overlay-title"
                style={{
                  margin: 0,
                  color: "#d9fbff",
                  fontSize: "clamp(0.82rem, 1.25vw, 0.98rem)",
                  lineHeight: 1.45,
                  textShadow: "0 0 12px rgba(0,245,255,0.4)",
                }}
              >
                GUILD CHAT
              </h2>
            </div>

            <div style={{ display: "flex", gap: "8px" }}>
              <motion.button
                type="button"
                whileHover={{ y: -1, scale: 1.03 }}
                whileTap={{ y: 1, scale: 0.98 }}
                onClick={onExpand}
                style={controlButtonStyle}
              >
                [ EXPAND ]
              </motion.button>
              <motion.button
                type="button"
                whileHover={{ y: -1, scale: 1.03 }}
                whileTap={{ y: 1, scale: 0.98 }}
                onClick={onClose}
                aria-label="Close guild chat"
                style={{ ...controlButtonStyle, minWidth: "40px", padding: "0 14px" }}
              >
                ×
              </motion.button>
            </div>
          </header>

          <ChatMessageList messages={messages} dense />
          <ChatComposer channel={channel} />
        </motion.aside>
      )}
    </AnimatePresence>
  );
}

const controlButtonStyle: React.CSSProperties = {
  minWidth: "112px",
  minHeight: "40px",
  border: "2px solid rgba(0, 245, 255, 0.68)",
  borderBottomColor: "rgba(2, 54, 72, 0.96)",
  borderRightColor: "rgba(2, 54, 72, 0.96)",
  background: "rgba(3, 12, 24, 0.84)",
  color: "#d9fbff",
  cursor: "pointer",
  fontFamily: "inherit",
  fontSize: "0.54rem",
  lineHeight: 1.1,
  padding: "0 10px",
  whiteSpace: "nowrap",
};
