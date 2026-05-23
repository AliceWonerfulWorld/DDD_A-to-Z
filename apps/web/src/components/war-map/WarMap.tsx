import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
} from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { fetchGuilds, fetchMyGuild } from "../../features/guild/api";
import { PATHS } from "../../constants/paths";
import { BackButton } from "../guild-town/BackButton";
import { GuildBgm } from "../shared/GuildBgm";
import { RankingPanel } from "./RankingPanel";
import { ScoutPanel } from "./ScoutPanel";
import { createWarGuilds, findWarGuildByID, type WarGuild } from "./WarMapData";
import { WarMapHex } from "./WarMapHex";
import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

interface WarMapProps {
  onNavigate: (path: string) => void;
}

export function WarMap({ onNavigate }: WarMapProps) {
  const { isSeEnabled } = useAudioSettings();
  const [currentGuildID, setCurrentGuildID] = useState<string | null>(null);
  const [isCurrentGuildLoaded, setIsCurrentGuildLoaded] = useState(false);
  const [isGuildsLoaded, setIsGuildsLoaded] = useState(false);
  const [warGuilds, setWarGuilds] = useState<WarGuild[]>([]);
  const [isRankingOpen, setIsRankingOpen] = useState(false);
  const [selectedGuild, setSelectedGuild] = useState<WarGuild | null>(null);
  const guildScoutSeRef = useRef<HTMLAudioElement | null>(null);
  const rankingToggleSeRef = useRef<HTMLAudioElement | null>(null);
  const scoutCloseSeRef = useRef<HTMLAudioElement | null>(null);

  useEffect(() => {
    let isMounted = true;

    fetchGuilds()
      .then((guilds) => {
        if (!isMounted) {
          return;
        }

        setWarGuilds(createWarGuilds(guilds));
      })
      .catch((error) => {
        if (!isMounted) {
          return;
        }

        console.error("failed to fetch guilds for war map", error);
        setWarGuilds([]);
      })
      .finally(() => {
        if (isMounted) {
          setIsGuildsLoaded(true);
        }
      });

    fetchMyGuild()
      .then((data) => {
        if (!isMounted) {
          return;
        }

        const nextGuildID = data?.guild?.slug ?? null;
        setCurrentGuildID(nextGuildID);
      })
      .catch((error) => {
        if (!isMounted) {
          return;
        }

        console.error("failed to fetch my guild for war map", error);
        setCurrentGuildID(null);
      })
      .finally(() => {
        if (isMounted) {
          setIsCurrentGuildLoaded(true);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    setSelectedGuild((current) => {
      if (current) {
        return findWarGuildByID(warGuilds, current.id);
      }

      return findWarGuildByID(warGuilds, currentGuildID);
    });
  }, [currentGuildID, warGuilds]);

  const currentGuild = findWarGuildByID(warGuilds, currentGuildID);

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

      if (typeof window !== "undefined" && window.innerWidth <= 768) {
        setIsRankingOpen(false);
      }
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
    setIsRankingOpen((isOpen) => {
      if (!isOpen && typeof window !== "undefined" && window.innerWidth <= 768) {
        setSelectedGuild(null);
      }
      return !isOpen;
    });
  }, [playSe]);

  const handleWorldPointerDown = (event: ReactPointerEvent<HTMLElement>) => {
    const target = event.target;
    if (!(target instanceof Element)) return;
    if (target.closest("[data-war-interactive='true']")) return;
    closeScout();
    setIsRankingOpen(false);
  };

  const transformComponentRef = useRef<any>(null);

  const [isMobile, setIsMobile] = useState(
    typeof window !== "undefined" ? window.innerWidth <= 768 : false,
  );

  useEffect(() => {
    let lastWidth = window.innerWidth;
    let timeoutId: number;

    const handleResize = () => {
      const currentWidth = window.innerWidth;
      setIsMobile(currentWidth <= 768);

      if (currentWidth !== lastWidth) {
        lastWidth = currentWidth;
        clearTimeout(timeoutId);
        timeoutId = window.setTimeout(() => {
          if (transformComponentRef.current) {
            transformComponentRef.current.resetTransform(0);
          }
        }, 200);
      }
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
      clearTimeout(timeoutId);
    };
  }, []);

  const MapContent = () => (
    <div
      style={{
        width: "100%",
        height: "100%",
        position: "relative",
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

      {warGuilds.map((guild) => (
        <WarMapHex
          key={guild.id}
          guild={guild}
          isCurrentGuild={guild.id === currentGuildID}
          isSelected={selectedGuild?.id === guild.id}
          onSelect={selectGuild}
        />
      ))}
    </div>
  );

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

      <div style={{ width: "100vw", height: "100svh", overflow: "hidden" }}>
        <TransformWrapper
          ref={transformComponentRef}
          initialScale={1}
          minScale={1}
          maxScale={3}
          centerOnInit={true}
          limitToBounds={true}
        >
          <TransformComponent
            wrapperStyle={{ width: "100%", height: "100%" }}
            contentStyle={{ width: "max(100vw, 177.78svh)", height: "max(100svh, 56.25vw)" }}
          >
            <MapContent />
          </TransformComponent>
        </TransformWrapper>
      </div>

      <RankingPanel
        currentGuildID={currentGuildID}
        guilds={warGuilds}
        isOpen={isRankingOpen}
        isLoaded={isGuildsLoaded}
        onToggle={toggleRankingWithSe}
      />
      <ScoutPanel
        guild={selectedGuild}
        isCurrentGuild={selectedGuild?.id === currentGuildID}
        onClose={closeScoutWithSe}
      />
      <BackButton onNavigate={onNavigate} targetPath={PATHS.HOME} />

      <div
        data-war-interactive="true"
        style={{
          position: "fixed",
          left: "50%",
          top: "calc(env(safe-area-inset-top, 0px) + 14px)",
          zIndex: 12,
          border: `3px solid ${currentGuild?.color ?? "rgba(255,217,102,0.72)"}`,
          borderBottomColor: "rgba(28, 20, 8, 0.98)",
          borderRightColor: "rgba(28, 20, 8, 0.98)",
          background: "rgba(3, 10, 24, 0.88)",
          boxShadow:
            "0 0 0 2px rgba(0,0,0,0.76), 6px 6px 0 rgba(0,0,0,0.36), inset 0 0 18px rgba(255,255,255,0.06)",
          color: "#fff8d7",
          minWidth: "min(360px, calc(100vw - 210px))",
          maxWidth: "calc(100vw - 210px)",
          padding: "10px 14px",
          textAlign: "center",
          textShadow: "2px 2px 0 rgba(0,0,0,0.8)",
          transform: "translateX(-50%)",
        }}
      >
        <p
          style={{
            margin: "0 0 6px",
            color: currentGuild?.accent ?? "#ffd966",
            fontSize: "0.44rem",
            lineHeight: 1.4,
          }}
        >
          CURRENT BANNER
        </p>
        <p
          style={{
            margin: 0,
            color: "#fff8d7",
            fontSize: isMobile ? "0.44rem" : "clamp(0.48rem, 1vw, 0.66rem)",
            lineHeight: 1.5,
            whiteSpace: "nowrap",
            overflow: "hidden",
            textOverflow: "ellipsis",
          }}
        >
          {!isCurrentGuildLoaded && "Synchronizing guild data..."}
          {isCurrentGuildLoaded &&
            currentGuild &&
            `${currentGuild.name} Guild deployed on this front.`}
          {isCurrentGuildLoaded &&
            !currentGuild &&
            "No guild selected yet. Join a guild to claim a front."}
        </p>
      </div>

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
