import { AnimatePresence, motion, type Variants } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { steppedEase } from "../lib/animationUtils";

interface GuildDashboardProps {
  onNavigate: (path: string) => void;
}

type GuildTab = "activity" | "rankings";

interface ActivityLog {
  id: number;
  player: string;
  action: string;
  cp: number;
  tone: string;
}

const PLAYERS = [
  "AkiByte",
  "NullMage",
  "PixelNinja",
  "LoopKnight",
  "TypeSmith",
  "AsyncRogue",
  "BugSlayer",
  "CacheWizard",
];

const LOG_ACTIONS = [
  { action: "merged a pull request", cp: 420, tone: "#74f7a1" },
  { action: "reviewed 3 files", cp: 180, tone: "#9be7ff" },
  { action: "fixed a flaky test", cp: 260, tone: "#ffd966" },
  { action: "shipped guild docs", cp: 150, tone: "#d9b8ff" },
  { action: "closed a critical issue", cp: 540, tone: "#ff9b9b" },
  { action: "refactored a use case", cp: 310, tone: "#74f7a1" },
];

const INITIAL_LOGS: ActivityLog[] = [
  {
    id: 1005,
    player: "TypeSmith",
    action: "merged a pull request",
    cp: 420,
    tone: "#74f7a1",
  },
  {
    id: 1004,
    player: "PixelNinja",
    action: "fixed a flaky test",
    cp: 260,
    tone: "#ffd966",
  },
  {
    id: 1003,
    player: "LoopKnight",
    action: "reviewed 3 files",
    cp: 180,
    tone: "#9be7ff",
  },
  {
    id: 1002,
    player: "NullMage",
    action: "closed a critical issue",
    cp: 540,
    tone: "#ff9b9b",
  },
  {
    id: 1001,
    player: "AkiByte",
    action: "refactored a use case",
    cp: 310,
    tone: "#74f7a1",
  },
];

const RANKINGS = [
  { name: "TypeSmith", title: "Generic Hero", cp: 35420, color: "#ffd966" },
  { name: "NullMage", title: "Void Debugger", cp: 31980, color: "#d9b8ff" },
  { name: "PixelNinja", title: "UI Shinobi", cp: 28640, color: "#9be7ff" },
  { name: "LoopKnight", title: "Iteration Paladin", cp: 25110, color: "#74f7a1" },
  { name: "AsyncRogue", title: "Promise Runner", cp: 22470, color: "#ff9b9b" },
  { name: "CacheWizard", title: "Memo Sage", cp: 19860, color: "#f4ecd0" },
  { name: "BugSlayer", title: "Regression Breaker", cp: 17420, color: "#f6a6ff" },
];

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

function createLog(id: number): ActivityLog {
  const player = PLAYERS[Math.floor(Math.random() * PLAYERS.length)];
  const log = LOG_ACTIONS[Math.floor(Math.random() * LOG_ACTIONS.length)];

  return {
    id,
    player,
    action: log.action,
    cp: log.cp,
    tone: log.tone,
  };
}

function ActivityLogPanel({ logs }: { logs: ActivityLog[] }) {
  return (
    <motion.section
      key="activity"
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

function RankingsPanel() {
  return (
    <motion.section
      key="rankings"
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

export function GuildDashboard({ onNavigate }: GuildDashboardProps) {
  const [activeTab, setActiveTab] = useState<GuildTab>("activity");
  const [logs, setLogs] = useState<ActivityLog[]>(INITIAL_LOGS);

  useEffect(() => {
    let nextId = INITIAL_LOGS[0].id + 1;
    const intervalId = window.setInterval(() => {
      setLogs((current) => [createLog(nextId++), ...current].slice(0, 8));
    }, 2800);

    return () => window.clearInterval(intervalId);
  }, []);

  const tabs = useMemo(
    () => [
      { id: "activity" as const, label: "ACTIVITY LOG" },
      { id: "rankings" as const, label: "RANKINGS" },
    ],
    [],
  );

  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        background: "#07172b",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#f4ecd0",
      }}
    >
      <div
        style={{
          position: "absolute",
          left: "50%",
          top: "50%",
          width: "max(100vw, calc(100svh * 1672 / 941))",
          height: "max(100svh, calc(100vw * 941 / 1672))",
          transform: "translate(-50%, -50%)",
        }}
      >
        <img
          src="/dashboard.png"
          alt=""
          aria-hidden="true"
          style={{
            position: "absolute",
            inset: 0,
            width: "100%",
            height: "100%",
            objectFit: "cover",
            imageRendering: "pixelated",
          }}
        />

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
              <strong style={{ color: "#9be7ff", overflowWrap: "anywhere" }}>
                TypeScript GUILD
              </strong>
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
                    onClick={() => setActiveTab(tab.id)}
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
      </div>

      <button
        type="button"
        onClick={() => onNavigate("/home")}
        style={{
          position: "fixed",
          top: "clamp(14px, 2.2vw, 28px)",
          left: "clamp(14px, 2.2vw, 28px)",
          zIndex: 3,
          border: "2px solid rgba(255, 217, 102, 0.78)",
          background: "rgba(3, 10, 24, 0.72)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.36)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.56rem, 1.3vw, 0.78rem)",
          lineHeight: 1.5,
          padding: "10px 12px",
        }}
      >
        &lt; 拠点に戻る
      </button>

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.09), rgba(0,0,0,0.09) 1px, transparent 1px, transparent 4px)",
          pointerEvents: "none",
          zIndex: 2,
        }}
      />
    </main>
  );
}
