import { AnimatePresence, motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import { ActivityLogPanel, RankingsPanel } from "./DashboardPanels";
import type { ActivityLog, GuildTab } from "./types";

interface DashboardMonitorProps {
  activeTab: GuildTab;
  logs: ActivityLog[];
  onSwitchTab: (tab: GuildTab) => void;
  tabs: { id: GuildTab; label: string }[];
}

export function DashboardMonitor({ activeTab, logs, onSwitchTab, tabs }: DashboardMonitorProps) {
  return (
    <motion.section
      initial={{ opacity: 0, scaleY: 0.94 }}
      animate={{ opacity: 1, scaleY: 1 }}
      transition={{ duration: 0.36, ease: steppedEase(6) }}
      aria-label="Guild dashboard monitor"
      style={{
        position: "absolute",
        left: "29.2%",
        top: "16.2%",
        width: "41.6%",
        height: "44.2%",
        display: "grid",
        gridTemplateRows: "auto 1fr",
        gap: "clamp(8px, 1.2vw, 14px)",
        padding: "clamp(14px, 2.1vw, 26px)",
        background: "linear-gradient(180deg, rgba(3, 10, 30, 0.32), rgba(2, 8, 24, 0.58))",
        boxShadow: "inset 0 0 34px rgba(0, 245, 255, 0.1)",
        overflow: "hidden",
      }}
    >
      <header
        style={{
          display: "grid",
          gap: "clamp(8px, 1.2vw, 12px)",
          minWidth: 0,
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "1fr auto auto",
            alignItems: "center",
            gap: "clamp(8px, 1.5vw, 18px)",
            color: "#fff8d7",
            fontSize: "clamp(0.58rem, 1.15vw, 0.86rem)",
            lineHeight: 1.5,
          }}
        >
          <strong style={{ color: "#9be7ff", overflowWrap: "anywhere" }}>TypeScript GUILD</strong>
          <span style={{ color: "#ffd966", whiteSpace: "nowrap" }}>Rank: #3</span>
          <span style={{ color: "#74f7a1", whiteSpace: "nowrap" }}>Total CP: 125,400</span>
        </div>

        <nav
          aria-label="Guild dashboard tabs"
          style={{
            display: "grid",
            gridTemplateColumns: "1fr 1fr",
            gap: "8px",
          }}
        >
          {tabs.map((tab) => {
            const isActive = activeTab === tab.id;

            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => onSwitchTab(tab.id)}
                style={{
                  minHeight: "34px",
                  border: `2px solid ${isActive ? "#00f5ff" : "rgba(0, 245, 255, 0.28)"}`,
                  background: isActive ? "rgba(0, 245, 255, 0.14)" : "rgba(1, 8, 22, 0.46)",
                  color: isActive ? "#fff8d7" : "rgba(244, 236, 208, 0.64)",
                  boxShadow: isActive ? "inset 0 0 16px rgba(0, 245, 255, 0.18)" : "none",
                  cursor: "pointer",
                  fontFamily: "inherit",
                  fontSize: "clamp(0.5rem, 0.95vw, 0.7rem)",
                  lineHeight: 1.4,
                  padding: "6px",
                }}
              >
                [ {tab.label} ]
              </button>
            );
          })}
        </nav>
      </header>

      <div style={{ minHeight: 0, overflow: "hidden" }}>
        <AnimatePresence mode="wait">
          {activeTab === "activity" ? <ActivityLogPanel logs={logs} /> : <RankingsPanel />}
        </AnimatePresence>
      </div>
    </motion.section>
  );
}
