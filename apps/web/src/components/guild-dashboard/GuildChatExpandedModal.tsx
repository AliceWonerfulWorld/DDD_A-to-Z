import { AnimatePresence, motion } from "framer-motion";
import type { Channel } from "phoenix";
import type { GuildChatMessage } from "../../features/guild/chatApi";
import { steppedEase } from "../../lib/animationUtils";
import { ChatComposer } from "./ChatComposer";
import { ChatMessageList } from "./ChatMessageList";

interface GuildChatExpandedModalProps {
  isOpen: boolean;
  messages: GuildChatMessage[];
  channel: Channel | null;
  onMinimize: () => void;
  onClose: () => void;
}

export function GuildChatExpandedModal({
  isOpen,
  messages,
  channel,
  onMinimize,
  onClose,
}: GuildChatExpandedModalProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            key="guild-chat-expanded-backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2, ease: steppedEase(4) }}
            onClick={onClose}
            style={{
              position: "fixed",
              inset: 0,
              zIndex: 4,
              background: "rgba(1, 6, 16, 0.44)",
              backdropFilter: "blur(4px)",
            }}
          />
          <motion.section
            key="guild-chat-expanded-modal"
            initial={{ opacity: 0, scale: 0.96, y: 18 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.98, y: 12 }}
            transition={{ duration: 0.28, ease: steppedEase(5) }}
            aria-labelledby="guild-chat-expanded-title"
            style={{
              position: "fixed",
              inset: "clamp(20px, 4vw, 40px)",
              zIndex: 5,
              display: "grid",
              gridTemplateRows: "auto auto minmax(0, 1fr) auto",
              gap: "14px",
              border: "2px solid rgba(0, 245, 255, 0.76)",
              borderLeftColor: "rgba(0, 245, 255, 0.92)",
              background:
                "linear-gradient(180deg, rgba(3, 10, 24, 0.82), rgba(1, 8, 22, 0.92)), repeating-linear-gradient(180deg, rgba(130,234,255,0.05), rgba(130,234,255,0.05) 1px, transparent 1px, transparent 4px)",
              boxShadow:
                "-8px 0 24px rgba(0,243,255,0.24), 0 0 0 2px rgba(0,0,0,0.64), inset 0 0 28px rgba(0,245,255,0.1)",
              color: "#f4ecd0",
              padding: "clamp(16px, 2vw, 24px)",
              overflow: "hidden",
            }}
          >
            <header
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "start",
                gap: "12px",
              }}
            >
              <div>
                <span
                  style={{
                    display: "block",
                    marginBottom: "8px",
                    color: "rgba(155, 231, 255, 0.72)",
                    fontFamily: '"DotGothic16", monospace',
                    fontSize: "0.68rem",
                    lineHeight: 1.4,
                  }}
                >
                  ARCHIVE RELAY
                </span>
                <h2
                  id="guild-chat-expanded-title"
                  style={{
                    margin: 0,
                    color: "#d9fbff",
                    fontSize: "clamp(1rem, 1.6vw, 1.28rem)",
                    lineHeight: 1.45,
                    textShadow: "0 0 12px rgba(0,245,255,0.42)",
                  }}
                >
                  GUILD CHAT LOG
                </h2>
              </div>

              <div style={{ display: "flex", gap: "8px" }}>
                <motion.button
                  type="button"
                  whileHover={{ y: -1, scale: 1.03 }}
                  whileTap={{ y: 1, scale: 0.98 }}
                  onClick={onMinimize}
                  style={controlButtonStyle}
                >
                  [ MINIMIZE ]
                </motion.button>
                <motion.button
                  type="button"
                  whileHover={{ y: -1, scale: 1.03 }}
                  whileTap={{ y: 1, scale: 0.98 }}
                  onClick={onClose}
                  aria-label="Close guild chat log"
                  style={{ ...controlButtonStyle, minWidth: "44px", padding: "0 14px" }}
                >
                  ×
                </motion.button>
              </div>
            </header>

            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                gap: "12px",
                flexWrap: "wrap",
                color: "#74f7a1",
                fontSize: "0.64rem",
                lineHeight: 1.5,
                borderTop: "1px solid rgba(0, 245, 255, 0.18)",
                borderBottom: "1px solid rgba(0, 245, 255, 0.18)",
                padding: "10px 0",
              }}
            >
              <span>FULL ARCHIVE MODE</span>
              <span>{messages.length} LOG ENTRIES AVAILABLE</span>
            </div>

            <ChatMessageList messages={messages} />
            <ChatComposer channel={channel} placeholder="write a message to the whole guild..." />
          </motion.section>
        </>
      )}
    </AnimatePresence>
  );
}

const controlButtonStyle: React.CSSProperties = {
  minHeight: "42px",
  border: "2px solid rgba(0, 245, 255, 0.68)",
  borderBottomColor: "rgba(2, 54, 72, 0.96)",
  borderRightColor: "rgba(2, 54, 72, 0.96)",
  background: "rgba(3, 12, 24, 0.84)",
  color: "#d9fbff",
  cursor: "pointer",
  fontFamily: "inherit",
  fontSize: "0.56rem",
  lineHeight: 1,
  padding: "0 12px",
};
