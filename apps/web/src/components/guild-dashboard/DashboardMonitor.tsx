import { AnimatePresence, motion } from "framer-motion";
import type { DisplayGuild } from "../../features/guild/presentation";
import { steppedEase } from "../../lib/animationUtils";
import { ActivityLogPanel, RankingsPanel } from "./DashboardPanels";
import type { ActivityLog, GuildTab } from "./types";

interface DashboardMonitorProps {
  activeTab: GuildTab;
  guild: DisplayGuild | null;
  isGuildLoading: boolean;
  logs: ActivityLog[];
  onSwitchTab: (tab: GuildTab) => void;
  tabs: { id: GuildTab; label: string }[];
}

export function DashboardMonitor({
  activeTab,
  guild,
  isGuildLoading,
  logs,
  onSwitchTab,
  tabs,
}: DashboardMonitorProps) {
  const guildName = isGuildLoading ? "SYNCING GUILD" : guild ? `${guild.name} GUILD` : "NO GUILD";
  const guildRank = guild ? `Rank: #${guild.sortOrder + 1}` : "Rank: --";
  const guildLevel = guild?.guildLevel ?? 1;
  const guildExperience = guild?.guildExperience ?? 0;
  const currentLevelExperience = guild?.currentGuildLevelExperience ?? 0;
  const nextLevelExperience = guild?.nextGuildLevelExperience;
  const isMaxGuildLevel =
    guild?.isMaxLevel ??
    (guild
      ? guildLevel >= 5 || nextLevelExperience === null || nextLevelExperience === undefined
      : false);
  const levelRange = Math.max(1, (nextLevelExperience ?? 5000) - currentLevelExperience);
  const progressInLevel = Math.min(
    levelRange,
    Math.max(0, guildExperience - currentLevelExperience),
  );
  const progressPercent = isMaxGuildLevel ? 100 : Math.round((progressInLevel / levelRange) * 100);
  const progressValue = isMaxGuildLevel
    ? 100
    : Math.min(100, Math.max(0, (progressInLevel / levelRange) * 100));

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
            gridTemplateColumns: "1fr auto",
            alignItems: "center",
            gap: "clamp(8px, 1.5vw, 18px)",
            color: "#fff8d7",
            fontSize: "clamp(0.58rem, 1.15vw, 0.86rem)",
            lineHeight: 1.5,
          }}
        >
          <strong style={{ color: guild?.accent ?? "#9be7ff", overflowWrap: "anywhere" }}>
            {guildName}
          </strong>
          <span style={{ color: "#ffd966", whiteSpace: "nowrap" }}>{guildRank}</span>
        </div>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "auto 1fr auto",
            alignItems: "center",
            gap: "10px",
            color: "#d9fbff",
            fontSize: "clamp(0.48rem, 0.82vw, 0.62rem)",
            lineHeight: 1.4,
          }}
        >
          <span style={{ color: "#ffd966", whiteSpace: "nowrap" }}>EXP</span>
          <div
            role="progressbar"
            aria-label="Guild level experience"
            aria-valuemin={0}
            aria-valuemax={levelRange}
            aria-valuenow={isMaxGuildLevel ? levelRange : progressInLevel}
            aria-valuetext={
              isMaxGuildLevel
                ? "MAX"
                : `${progressInLevel.toLocaleString()} of ${levelRange.toLocaleString()} (${progressPercent}%)`
            }
            style={{
              height: "10px",
              border: "2px solid rgba(116, 247, 161, 0.6)",
              background: "rgba(1, 8, 22, 0.72)",
              boxShadow: "inset 0 0 8px rgba(0,0,0,0.64)",
              overflow: "hidden",
            }}
          >
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${progressValue}%` }}
              transition={{ duration: 0.42, ease: steppedEase(8) }}
              style={{
                height: "100%",
                background:
                  "repeating-linear-gradient(90deg, #74f7a1 0, #74f7a1 8px, #39ff14 8px, #39ff14 16px)",
                boxShadow: "0 0 10px rgba(116,247,161,0.68)",
              }}
            />
          </div>
          <span style={{ color: "#f4ecd0", whiteSpace: "nowrap" }}>
            {isMaxGuildLevel
              ? "MAX"
              : `${guildExperience.toLocaleString()} / ${(nextLevelExperience ?? 5000).toLocaleString()}`}
          </span>
        </div>

        <nav
          aria-label="Guild dashboard tabs"
          role="tablist"
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
                id={`guild-dashboard-${tab.id}-tab`}
                role="tab"
                aria-controls={`guild-dashboard-${tab.id}-panel`}
                aria-selected={isActive}
                tabIndex={isActive ? 0 : -1}
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
