import { motion } from "framer-motion";

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

const STATUS_MESSAGES = [
  "Connecting to GitHub API...",
  "Fetching recent commits...",
  "Scanning repository languages...",
  "Analyzing contribution patterns...",
  "Calculating Contribution Points...",
  "Syncing with guild data...",
  "Finalizing results...",
];

interface AnalyzingPanelProps {
  progress: number;
  currentMessageIdx: number;
}

export function AnalyzingPanel({ progress, currentMessageIdx }: AnalyzingPanelProps) {
  return (
    <motion.div
      initial={{ scaleY: 0, opacity: 0 }}
      animate={{ scaleY: 1, opacity: 1 }}
      transition={{ duration: 0.4, ease: steppedEase(6) }}
      style={{
        position: "relative",
        zIndex: 2,
        border: "4px solid var(--color-gold)",
        background: "var(--color-navy-light)",
        boxShadow: "0 0 30px rgba(0, 245, 255, 0.12), 8px 8px 0 rgba(0,0,0,0.8)",
        maxWidth: "520px",
        width: "100%",
        padding: "2rem",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
      }}
    >
      <div style={{ textAlign: "center" }}>
        <motion.span
          animate={{ opacity: [1, 0.5, 1] }}
          transition={{ duration: 1.5, repeat: Infinity, ease: steppedEase(3) }}
          style={{
            display: "inline-block",
            fontFamily: "var(--font-press)",
            fontSize: "0.85rem",
            color: "var(--color-neon-cyan)",
            letterSpacing: "0.1em",
          }}
        >
          ▶ GITHUB ACTIVITY SCAN
        </motion.span>
      </div>

      <div
        style={{
          position: "relative",
          height: "120px",
          border: "2px solid rgba(0, 245, 255, 0.3)",
          background: "rgba(0,0,0,0.4)",
          overflow: "hidden",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <motion.div
          animate={{ top: ["0%", "100%"] }}
          transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
          style={{
            position: "absolute",
            left: 0,
            right: 0,
            height: "3px",
            background: "linear-gradient(90deg, transparent, var(--color-neon-cyan), transparent)",
            boxShadow: "0 0 8px var(--color-neon-cyan)",
            zIndex: 1,
          }}
        />
        <span
          style={{
            fontFamily: "var(--font-dot)",
            fontSize: "0.8rem",
            color: "rgba(232, 232, 208, 0.5)",
            letterSpacing: "0.15em",
          }}
        >
          {STATUS_MESSAGES[Math.min(currentMessageIdx, STATUS_MESSAGES.length - 1)]}
        </span>
      </div>

      <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
        <div
          style={{
            width: "100%",
            height: "18px",
            border: "2px solid var(--color-gold)",
            background: "rgba(0,0,0,0.5)",
            position: "relative",
            overflow: "hidden",
          }}
        >
          <motion.div
            animate={{ width: `${progress}%` }}
            transition={{ duration: 0.3, ease: steppedEase(8) }}
            style={{
              height: "100%",
              background: "linear-gradient(90deg, var(--color-gold-dark), var(--color-gold))",
              boxShadow: "0 0 6px var(--color-gold)",
              position: "absolute",
              left: 0,
              top: 0,
            }}
          />
        </div>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            fontFamily: "var(--font-dot)",
            fontSize: "0.7rem",
            color: "var(--color-pixel-white)",
          }}
        >
          <span>SCANNING...</span>
          <span>{Math.round(progress)}%</span>
        </div>
      </div>

      <div style={{ minHeight: "2rem" }}>
        {STATUS_MESSAGES.slice(0, currentMessageIdx + 1).map((msg, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: i === currentMessageIdx ? 1 : 0.4, x: 0 }}
            style={{
              fontFamily: "var(--font-dot)",
              fontSize: "0.75rem",
              color: i === currentMessageIdx ? "var(--color-neon-cyan)" : "rgba(232, 232, 208, 0.3)",
              padding: "2px 0",
            }}
          >
            {i === currentMessageIdx ? "> " : "  "}{msg}
          </motion.div>
        ))}
      </div>
    </motion.div>
  );
}
