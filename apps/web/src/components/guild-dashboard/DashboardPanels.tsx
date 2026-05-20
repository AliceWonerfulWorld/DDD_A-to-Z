import { AnimatePresence, motion, type Variants } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import { RANKINGS } from "./data";
import type { ActivityLog } from "./types";

const tabContentVariants: Variants = {
  hidden: { opacity: 0, y: 12 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.28, ease: steppedEase(5) },
  },
  exit: {
    opacity: 0,
    y: -10,
    transition: { duration: 0.18, ease: steppedEase(4) },
  },
};

export function ActivityLogPanel({ logs }: { logs: ActivityLog[] }) {
  return (
    <motion.section
      key="activity"
      id="guild-dashboard-activity-panel"
      role="tabpanel"
      aria-labelledby="guild-dashboard-activity-tab"
      variants={tabContentVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      style={{
        height: "100%",
        minHeight: 0,
        border: "2px solid rgba(0, 245, 255, 0.44)",
        background: "rgba(1, 8, 22, 0.74)",
        boxShadow: "inset 0 0 22px rgba(0, 245, 255, 0.12)",
        padding: "clamp(10px, 1.6vw, 18px)",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          gap: "12px",
          color: "#74f7a1",
          fontSize: "clamp(0.54rem, 0.95vw, 0.72rem)",
          lineHeight: 1.5,
          marginBottom: "10px",
        }}
      >
        <span>LIVE ACTIVITY STREAM</span>
        <span>STATUS: ONLINE</span>
      </div>

      <div
        style={{
          display: "grid",
          alignContent: "start",
          gap: "8px",
          height: "calc(100% - 28px)",
        }}
      >
        <AnimatePresence initial={false}>
          {logs.length === 0 ? (
            <motion.div
              key="empty"
              initial={{ opacity: 0 }}
              animate={{ opacity: 0.62 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.22, ease: steppedEase(4) }}
              style={{
                minHeight: "30px",
                color: "rgba(214, 255, 228, 0.72)",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.66rem, 1.22vw, 0.94rem)",
                lineHeight: 1.35,
              }}
            >
              &gt; WAITING FOR GUILD SIGNAL_
            </motion.div>
          ) : null}
          {logs.map((log) => (
            <motion.div
              key={log.id}
              layout
              initial={{ opacity: 0, x: -18 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 18 }}
              transition={{ duration: 0.28, ease: steppedEase(5) }}
              style={{
                display: "grid",
                gridTemplateColumns: "auto 1fr auto",
                alignItems: "baseline",
                gap: "10px",
                minHeight: "30px",
                borderBottom: "1px solid rgba(116, 247, 161, 0.12)",
                color: "#d6ffe4",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.66rem, 1.22vw, 0.94rem)",
                lineHeight: 1.35,
              }}
            >
              <span style={{ color: log.tone }}>&gt;</span>
              <span style={{ overflowWrap: "anywhere" }}>
                <span style={{ color: "#fff8d7" }}>[{log.player}]</span> {log.action}
              </span>
              <span style={{ color: "#ffd966", whiteSpace: "nowrap" }}>
                +{log.cp.toLocaleString()} CP
              </span>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </motion.section>
  );
}

export function RankingsPanel() {
  return (
    <motion.section
      key="rankings"
      id="guild-dashboard-rankings-panel"
      role="tabpanel"
      aria-labelledby="guild-dashboard-rankings-tab"
      variants={tabContentVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      style={{
        height: "100%",
        minHeight: 0,
        display: "grid",
        alignContent: "start",
        gap: "8px",
        overflow: "hidden",
      }}
    >
      {RANKINGS.map((member, index) => (
        <motion.div
          key={member.name}
          initial={{ opacity: 0, y: -12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: index * 0.07, duration: 0.28, ease: steppedEase(5) }}
          style={{
            display: "grid",
            gridTemplateColumns: "minmax(34px, auto) 1fr auto",
            alignItems: "center",
            gap: "clamp(8px, 1.4vw, 16px)",
            minHeight: "clamp(38px, 6.5vh, 54px)",
            border: "2px solid rgba(0, 245, 255, 0.24)",
            background:
              index === 0
                ? "linear-gradient(90deg, rgba(255, 217, 102, 0.22), rgba(1, 8, 22, 0.68))"
                : "rgba(1, 8, 22, 0.58)",
            boxShadow: "inset 0 0 16px rgba(0, 245, 255, 0.08)",
            padding: "8px clamp(8px, 1.5vw, 16px)",
          }}
        >
          <span
            style={{
              color: member.color,
              fontSize: "clamp(0.72rem, 1.45vw, 1rem)",
              lineHeight: 1,
            }}
          >
            #{index + 1}
          </span>
          <span style={{ minWidth: 0 }}>
            <span
              style={{
                display: "block",
                color: "#fff8d7",
                fontSize: "clamp(0.68rem, 1.3vw, 0.95rem)",
                lineHeight: 1.4,
                overflowWrap: "anywhere",
              }}
            >
              {member.name}
            </span>
            <span
              style={{
                display: "block",
                color: "rgba(244, 236, 208, 0.62)",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.58rem, 1.05vw, 0.76rem)",
                lineHeight: 1.35,
                overflowWrap: "anywhere",
              }}
            >
              {member.title}
            </span>
          </span>
          <span
            style={{
              color: "#74f7a1",
              fontSize: "clamp(0.62rem, 1.15vw, 0.86rem)",
              lineHeight: 1.4,
              whiteSpace: "nowrap",
            }}
          >
            {member.cp.toLocaleString()} CP
          </span>
        </motion.div>
      ))}
    </motion.section>
  );
}
