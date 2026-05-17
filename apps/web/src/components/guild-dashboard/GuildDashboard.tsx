import { useCallback, useEffect, useRef, useState } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { fetchMyGuild } from "../../features/guild/api";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { PATHS } from "../../constants/paths";
import { DashboardMonitor } from "./DashboardMonitor";
import { createLog, GUILD_TABS, INITIAL_LOGS } from "./data";
import { GuildBadge } from "./GuildBadge";
import { GuildNavigation } from "./GuildNavigation";
import type { GuildTab } from "./types";

interface GuildDashboardProps {
  onNavigate: (path: string) => void;
}

export function GuildDashboard({ onNavigate }: GuildDashboardProps) {
  const { isSeEnabled } = useAudioSettings();
  const [activeTab, setActiveTab] = useState<GuildTab>("activity");
  const [logs, setLogs] = useState(INITIAL_LOGS);
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
    let nextId = INITIAL_LOGS[0].id + 1;
    const intervalId = window.setInterval(() => {
      setLogs((current) => [createLog(nextId++), ...current].slice(0, 8));
    }, 2800);

    return () => window.clearInterval(intervalId);
  }, []);

  useEffect(() => {
    fetchMyGuild()
      .then((data) => {
        if (!data?.guild) {
          onNavigate(PATHS.GUILD_SELECT);
        }
      })
      .catch(console.error);
  }, [onNavigate]);

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
          logs={logs}
          onSwitchTab={switchTab}
          tabs={GUILD_TABS}
        />
      </div>

      <GuildBadge />
      <GuildNavigation onNavigate={onNavigate} />

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
