import { motion } from "framer-motion";
import { useCallback, useEffect, useRef, useState } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import {
  fetchGuildActivityLogs,
  fetchMyGuild,
  type GuildActivityLog,
} from "../../features/guild/api";
import { connectChat, type GuildChatMessage } from "../../features/chat/api";
import type { ChatConnection } from "../../features/chat/api";
import { toDisplayGuild, type DisplayGuild } from "../../features/guild/presentation";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { steppedEase } from "../../lib/animationUtils";
import { PATHS } from "../../constants/paths";
import { GuildChatExpandedModal } from "../../features/chat/components/GuildChatExpandedModal";
import { GuildChatOverlay } from "../../features/chat/components/GuildChatOverlay";
import { DashboardMonitor } from "./DashboardMonitor";
import { GUILD_TABS, INITIAL_LOGS } from "./data";
import { GuildBadge } from "./GuildBadge";
import { GuildNavigation } from "./GuildNavigation";
import type { ActivityLog, GuildTab } from "./types";

const CHAT_SERVICE_URL = import.meta.env.VITE_CHAT_SERVICE_URL ?? "ws://localhost:4000";

interface GuildDashboardProps {
  onNavigate: (path: string) => void;
}

type ChatView = "closed" | "compact" | "expanded";

const ACTIVITY_LOG_LIMIT = 20;
const ACTIVITY_LOG_POLLING_MS = 10_000;

function toActivityLog(log: GuildActivityLog): ActivityLog {
  const prefix = log.type === "pull_request" ? "PR" : "Commit";

  return {
    id: log.id,
    player: log.player,
    action: `${prefix}: ${log.message}`,
    cp: log.cp,
    tone: log.type === "pull_request" ? "#74f7a1" : "#ffd966",
  };
}

export function GuildDashboard({ onNavigate }: GuildDashboardProps) {
  const { isSeEnabled } = useAudioSettings();
  const [activeTab, setActiveTab] = useState<GuildTab>("activity");
  const [chatView, setChatView] = useState<ChatView>("closed");
  const [logs, setLogs] = useState<ActivityLog[]>(import.meta.env.DEV ? INITIAL_LOGS : []);
  const [currentGuild, setCurrentGuild] = useState<DisplayGuild | null>(null);
  const [isCurrentGuildLoaded, setIsCurrentGuildLoaded] = useState(false);
  const [chatMessages, setChatMessages] = useState<GuildChatMessage[]>([]);
  const chatConnectionRef = useRef<ChatConnection | null>(null);
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);
  const tabSwitchSeRef = useRef<HTMLAudioElement | null>(null);

  const playTabSwitchSe = useCallback(() => {
    const audio = tabSwitchSeRef.current;
    if (!audio || !isSeEnabled) {
      return;
    }

    if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
      audio.load();
    }

    audio.currentTime = 0;
    void audio.play().catch(() => {});
  }, [isSeEnabled]);

  const switchTab = useCallback(
    (tab: GuildTab) => {
      if (activeTab === tab) {
        return;
      }

      playTabSwitchSe();
      setActiveTab(tab);
    },
    [activeTab, playTabSwitchSe],
  );

  useEffect(() => {
    let isMounted = true;

    fetchMyGuild()
      .then((data) => {
        if (!isMounted) {
          return;
        }

        if (!data?.guild) {
          onNavigate(PATHS.GUILD_SELECT);
          return;
        }

        setCurrentGuild(toDisplayGuild(data.guild));
      })
      .catch((error) => {
        if (!isMounted) {
          return;
        }

        console.error("failed to fetch my guild for dashboard", error);
      })
      .finally(() => {
        if (isMounted) {
          setIsCurrentGuildLoaded(true);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [onNavigate]);

  useEffect(() => {
    if (!currentGuild) {
      return;
    }

    let isMounted = true;

    const startChat = async () => {
      try {
        const connection = await connectChat(
          currentGuild.id,
          CHAT_SERVICE_URL,
          (msg) => {
            if (isMounted) {
              setChatMessages((prev) => [...prev, msg]);
            }
          },
          (msgs) => {
            if (isMounted) {
              setChatMessages(msgs);
            }
          },
          (reason) => {
            if (isMounted) {
              console.error("chat join error", reason);
            }
          },
        );
        if (isMounted) {
          chatConnectionRef.current = connection;
        } else {
          connection.disconnect();
        }
      } catch (error) {
        if (isMounted) {
          console.error("failed to connect to chat", error);
        }
      }
    };

    void startChat();

    return () => {
      isMounted = false;
      chatConnectionRef.current?.disconnect();
      chatConnectionRef.current = null;
      setChatMessages([]);
    };
  }, [currentGuild]);

  useEffect(() => {
    // DEV環境ではAPIを叩かずINITIAL_LOGSをそのまま表示する
    if (import.meta.env.DEV) {
      return;
    }

    if (!currentGuild || activeTab !== "activity") {
      return;
    }

    let isMounted = true;
    let intervalID: number | undefined;

    const loadLogs = () => {
      if (document.visibilityState === "hidden") {
        return;
      }

      fetchGuildActivityLogs(currentGuild.id, ACTIVITY_LOG_LIMIT)
        .then((activityLogs) => {
          if (isMounted) {
            setLogs(activityLogs.map(toActivityLog));
          }
        })
        .catch((error) => {
          if (isMounted) {
            console.error("failed to fetch guild activity logs", error);
          }
        });
    };

    loadLogs();
    intervalID = window.setInterval(loadLogs, ACTIVITY_LOG_POLLING_MS);

    return () => {
      isMounted = false;
      if (intervalID !== undefined) {
        window.clearInterval(intervalID);
      }
    };
  }, [activeTab, currentGuild]);

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
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />
      <audio
        ref={tabSwitchSeRef}
        src={AUDIO_ASSETS.se.buttonClick}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />

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

        <DashboardMonitor
          activeTab={activeTab}
          guild={currentGuild}
          isGuildLoading={!isCurrentGuildLoaded}
          logs={logs}
          onSwitchTab={switchTab}
          tabs={GUILD_TABS}
        />
      </div>

      <GuildBadge guild={currentGuild} isLoading={!isCurrentGuildLoaded} />
      <GuildNavigation onNavigate={onNavigate} />
      <motion.button
        type="button"
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.14, duration: 0.28, ease: steppedEase(5) }}
        whileHover={{ y: -2, scale: 1.02 }}
        whileTap={{ y: 1, scale: 0.98 }}
        onClick={() => setChatView((current) => (current === "closed" ? "compact" : "closed"))}
        aria-expanded={chatView !== "closed"}
        aria-controls={
          chatView === "expanded" ? "guild-chat-expanded-title" : "guild-chat-overlay-title"
        }
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(88px, 8vw, 112px))",
          right: "clamp(16px, 2.4vw, 32px)",
          zIndex: 4,
          minHeight: "44px",
          border: "2px solid rgba(0, 245, 255, 0.82)",
          borderBottomColor: "rgba(2, 54, 72, 0.96)",
          borderRightColor: "rgba(2, 54, 72, 0.96)",
          background: "rgba(3, 10, 24, 0.82)",
          boxShadow:
            "0 0 0 2px rgba(0,0,0,0.62), 0 0 16px rgba(0,245,255,0.24), inset 0 0 14px rgba(0,245,255,0.1)",
          color: "#d9fbff",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.54rem, 1vw, 0.72rem)",
          lineHeight: 1.45,
          padding: "10px 12px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        {chatView === "closed" ? "[ COMM LINK ]" : "[ CLOSE CHAT ]"}
      </motion.button>
      <GuildChatOverlay
        isOpen={chatView === "compact"}
        messages={chatMessages.slice(-4)}
        channel={chatConnectionRef.current?.channel ?? null}
        onExpand={() => setChatView("expanded")}
        onClose={() => setChatView("closed")}
      />
      <GuildChatExpandedModal
        isOpen={chatView === "expanded"}
        messages={chatMessages}
        channel={chatConnectionRef.current?.channel ?? null}
        onMinimize={() => setChatView("compact")}
        onClose={() => setChatView("closed")}
      />

      <button
        type="button"
        onClick={() => void navigateBackWithSe("/home")}
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
        &lt; BACK
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
