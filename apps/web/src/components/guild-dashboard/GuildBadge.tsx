import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";

export function GuildBadge() {
  return (
    <motion.aside
      initial={{ opacity: 0, x: -16 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.32, ease: steppedEase(6) }}
      aria-label="Guild identity"
      style={{
        position: "fixed",
        top: "clamp(78px, 9.4vw, 112px)",
        left: "clamp(14px, 2.2vw, 28px)",
        zIndex: 3,
        display: "grid",
        gridTemplateColumns: "64px minmax(0, 1fr)",
        alignItems: "start",
        gap: "12px",
        width: "min(calc(100vw - 28px), 330px)",
        border: "3px solid rgba(255, 217, 102, 0.8)",
        borderBottomColor: "rgba(96, 62, 22, 0.95)",
        borderRightColor: "rgba(96, 62, 22, 0.95)",
        background: "rgba(3, 10, 24, 0.76)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.72), 6px 6px 0 rgba(0,0,0,0.4), inset 0 0 18px rgba(255,217,102,0.1)",
        color: "#fff8d7",
        padding: "10px 12px",
        backdropFilter: "blur(2px)",
      }}
    >
      <div
        aria-hidden="true"
        style={{
          width: "64px",
          height: "64px",
          display: "grid",
          placeItems: "center",
          border: "3px solid #00f5ff",
          borderBottomColor: "#035a72",
          borderRightColor: "#035a72",
          background:
            "linear-gradient(135deg, rgba(0,245,255,0.28), rgba(255,217,102,0.16)), #061326",
          boxShadow: "inset 0 0 16px rgba(0,245,255,0.22), 3px 3px 0 rgba(0,0,0,0.42)",
          color: "#ffd966",
          fontSize: "1.55rem",
          lineHeight: 1,
        }}
      >
        TS
      </div>
      <div style={{ minWidth: 0 }}>
        <span
          style={{
            display: "block",
            color: "rgba(244, 236, 208, 0.58)",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.66rem",
            lineHeight: 1.5,
          }}
        >
          CURRENT GUILD
        </span>
        <strong
          style={{
            display: "block",
            color: "#fff8d7",
            fontSize: "clamp(0.76rem, 1.5vw, 0.96rem)",
            lineHeight: 1.55,
            overflowWrap: "anywhere",
          }}
        >
          TypeScript Guild
        </strong>
        <span
          aria-label="Current season"
          style={{
            display: "block",
            width: "fit-content",
            marginTop: "8px",
            border: "2px solid rgba(0, 245, 255, 0.66)",
            background: "rgba(1, 8, 22, 0.7)",
            boxShadow: "0 0 12px rgba(0,245,255,0.18), inset 0 0 12px rgba(0,245,255,0.12)",
            color: "#9be7ff",
            fontSize: "0.58rem",
            lineHeight: 1.5,
            padding: "4px 8px",
            textShadow: "0 0 8px rgba(0,245,255,0.72)",
          }}
        >
          SEASON 1
        </span>
      </div>
    </motion.aside>
  );
}
