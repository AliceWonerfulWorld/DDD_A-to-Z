import { useCallback, useRef, useState, type PointerEvent as ReactPointerEvent } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { PATHS } from "../../constants/paths";
import { BackButton } from "../guild-town/BackButton";
import { GuildBgm } from "../shared/GuildBgm";
import { RankingPanel } from "./RankingPanel";
import { ScoutPanel } from "./ScoutPanel";
import { WAR_GUILDS, type WarGuild } from "./WarMapData";
import { WarMapHex } from "./WarMapHex";

interface WarMapProps {
  onNavigate: (path: string) => void;
}

export function WarMap({ onNavigate }: WarMapProps) {
  const { isSeEnabled } = useAudioSettings();
  const [isRankingOpen, setIsRankingOpen] = useState(false);
  const [selectedGuild, setSelectedGuild] = useState<WarGuild | null>(null);
  const guildScoutSeRef = useRef<HTMLAudioElement | null>(null);
  const rankingToggleSeRef = useRef<HTMLAudioElement | null>(null);
  const scoutCloseSeRef = useRef<HTMLAudioElement | null>(null);

  const playSe = useCallback(
    (audio: HTMLAudioElement | null) => {
      if (!audio || !isSeEnabled) {
        return;
      }

      if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
        audio.load();
      }

      audio.currentTime = 0;
      void audio.play().catch(() => {});
    },
    [isSeEnabled],
  );

  const selectGuild = useCallback(
    (guild: WarGuild) => {
      setSelectedGuild((current) => {
        if (current?.id === guild.id) {
          return null;
        }

        playSe(guildScoutSeRef.current);
        return guild;
      });
    },
    [playSe],
  );

  const closeScout = () => {
    setSelectedGuild(null);
  };

  const closeScoutWithSe = useCallback(() => {
    playSe(scoutCloseSeRef.current);
    setSelectedGuild(null);
  }, [playSe]);

  const toggleRankingWithSe = useCallback(() => {
    playSe(rankingToggleSeRef.current);
    setIsRankingOpen((isOpen) => !isOpen);
  }, [playSe]);

  const handleWorldPointerDown = (event: ReactPointerEvent<HTMLElement>) => {
    const target = event.target;
    if (!(target instanceof Element)) return;
    if (target.closest("[data-war-interactive='true']")) return;
    closeScout();
    setIsRankingOpen(false);
  };

  return (
    <main
      className="relative h-screen w-full overflow-hidden"
      onPointerDown={handleWorldPointerDown}
      style={{
        background: "#03101c",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#fff8d7",
      }}
    >
      <GuildBgm src={AUDIO_ASSETS.bgm.warMap} />
      <audio
        ref={guildScoutSeRef}
        src={AUDIO_ASSETS.se.warGuildScout}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={rankingToggleSeRef}
        src={AUDIO_ASSETS.se.modalOpen}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={scoutCloseSeRef}
        src={AUDIO_ASSETS.se.modalCancel}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />

      <div
        className="absolute left-0 top-0 h-screen w-screen"
        style={{
          touchAction: "none",
          transformOrigin: "top left",
          userSelect: "none",
        }}
      >
        <img
          className="pixelated"
          src="/world_map.png"
          alt=""
          aria-hidden="true"
          draggable={false}
          style={{
            position: "absolute",
            inset: 0,
            width: "100%",
            height: "100%",
            objectFit: "cover",
            objectPosition: "center",
            pointerEvents: "none",
          }}
        />

        <div
          aria-hidden="true"
          style={{
            position: "absolute",
            inset: 0,
            background:
              "linear-gradient(180deg, rgba(0,245,255,0.08), rgba(0,0,0,0.03) 46%, rgba(0,0,0,0.32)), radial-gradient(ellipse at center, transparent 34%, rgba(0,0,0,0.58) 100%)",
            pointerEvents: "none",
            zIndex: 2,
          }}
        />

        <div
          aria-hidden="true"
          style={{
            position: "absolute",
            inset: 0,
            backgroundImage:
              "linear-gradient(rgba(0,245,255,0.08) 1px, transparent 1px), linear-gradient(90deg, rgba(0,245,255,0.08) 1px, transparent 1px)",
            backgroundSize: "64px 64px",
            mixBlendMode: "screen",
            opacity: 0.3,
            pointerEvents: "none",
            zIndex: 3,
          }}
        />

        {WAR_GUILDS.map((guild) => (
          <WarMapHex
            key={guild.id}
            guild={guild}
            isSelected={selectedGuild?.id === guild.id}
            onSelect={selectGuild}
          />
        ))}
      </div>

      <RankingPanel isOpen={isRankingOpen} onToggle={toggleRankingWithSe} />
      <ScoutPanel guild={selectedGuild} onClose={closeScoutWithSe} />
      <BackButton onNavigate={onNavigate} targetPath={PATHS.HOME} />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.1), rgba(0,0,0,0.1) 1px, transparent 1px, transparent 4px)",
          pointerEvents: "none",
          zIndex: 4,
        }}
      />
    </main>
  );
}
