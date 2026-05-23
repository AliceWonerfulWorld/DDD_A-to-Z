import { AnimatePresence, motion, type Variants } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import { RANKINGS } from "./data";
import type { ActivityLog } from "./types";
import styles from "./DashboardPanels.module.css";

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

export function ActivityLogPanel({ logs, isMobile }: { logs: ActivityLog[]; isMobile?: boolean }) {
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
        display: "flex",
        flexDirection: "column",
        border: isMobile
          ? "1px solid rgba(0, 245, 255, 0.44)"
          : "2px solid rgba(0, 245, 255, 0.44)",
        background: "rgba(1, 8, 22, 0.74)",
        boxShadow: "inset 0 0 22px rgba(0, 245, 255, 0.12)",
        padding: isMobile ? "4px 6px" : "clamp(10px, 1.6vw, 18px)",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          gap: isMobile ? "2px" : "12px",
          color: "#74f7a1",
          fontSize: isMobile ? "0.38rem" : "clamp(0.54rem, 0.95vw, 0.72rem)",
          lineHeight: isMobile ? 1.1 : 1.5,
          marginBottom: isMobile ? "4px" : "10px",
        }}
      >
        <span>LIVE ACTIVITY STREAM</span>
        <span>STATUS: ONLINE</span>
      </div>

      <div
        className={styles.hideScrollbar}
        style={{
          display: "flex",
          flexDirection: "column",
          gap: isMobile ? "4px" : "8px",
          flex: 1,
          minHeight: 0,
          overflowY: "auto",
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
                fontSize: isMobile ? "0.6rem" : "clamp(0.66rem, 1.22vw, 0.94rem)",
                lineHeight: isMobile ? 1.2 : 1.35,
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
                display: "flex",
                alignItems: "flex-start",
                gap: isMobile ? "4px" : "10px",
                padding: isMobile ? "2px 0" : "4px 0",
                borderBottom: "1px solid rgba(116, 247, 161, 0.12)",
                color: "#d6ffe4",
                fontFamily: '"DotGothic16", monospace',
                fontSize: isMobile ? "0.6rem" : "clamp(0.66rem, 1.22vw, 0.94rem)",
                lineHeight: isMobile ? 1.2 : 1.35,
              }}
            >
              <span style={{ color: log.tone, flexShrink: 0 }}>&gt;</span>
              <span
                style={{
                  flex: 1,
                  minWidth: 0,
                  overflowWrap: "anywhere",
                  wordBreak: "break-all",
                }}
              >
                <span style={{ color: "#fff8d7" }}>[{log.player}]</span> {log.action}
              </span>
              <span style={{ color: "#ffd966", whiteSpace: "nowrap", flexShrink: 0 }}>
                +{log.cp.toLocaleString()} CP
              </span>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </motion.section>
  );
}

export function RankingsPanel({ isMobile }: { isMobile?: boolean }) {
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
      className={styles.hideScrollbar}
      style={{
        height: "100%",
        minHeight: 0,
        display: "grid",
        alignContent: "start",
        gap: "8px",
        overflowY: "auto",
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
            alignItems: "start",
            gap: isMobile ? "4px" : "clamp(8px, 1.4vw, 16px)",
            minHeight: isMobile ? "30px" : "clamp(38px, 6.5vh, 54px)",
            border: isMobile
              ? "1px solid rgba(0, 245, 255, 0.24)"
              : "2px solid rgba(0, 245, 255, 0.24)",
            background:
              index === 0
                ? "linear-gradient(90deg, rgba(255, 217, 102, 0.22), rgba(1, 8, 22, 0.68))"
                : "rgba(1, 8, 22, 0.58)",
            boxShadow: "inset 0 0 16px rgba(0, 245, 255, 0.08)",
            padding: isMobile ? "4px 6px" : "8px clamp(8px, 1.5vw, 16px)",
          }}
        >
          <span
            style={{
              color: member.color,
              fontSize: isMobile ? "0.55rem" : "clamp(0.72rem, 1.45vw, 1rem)",
              lineHeight: 1,
            }}
          >
            #{index + 1}
          </span>
          <span style={{ minWidth: 0, overflow: "hidden" }}>
            <span
              style={{
                display: "block",
                color: "#fff8d7",
                fontSize: isMobile ? "0.6rem" : "clamp(0.68rem, 1.3vw, 0.95rem)",
                lineHeight: isMobile ? 1.2 : 1.4,
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
                fontSize: isMobile ? "0.5rem" : "clamp(0.58rem, 1.05vw, 0.76rem)",
                lineHeight: isMobile ? 1.2 : 1.35,
                overflowWrap: "anywhere",
              }}
            >
              {member.title}
            </span>
          </span>
          <span
            style={{
              color: "#74f7a1",
              fontSize: isMobile ? "0.55rem" : "clamp(0.62rem, 1.15vw, 0.86rem)",
              lineHeight: isMobile ? 1.2 : 1.4,
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
