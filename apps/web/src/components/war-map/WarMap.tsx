import { animate, motion, useMotionValue } from "framer-motion";
import {
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
  type WheelEvent,
} from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { steppedEase } from "../../lib/animationUtils";
import { BackButton } from "../guild-town/BackButton";
import { GuildBgm } from "../GuildBgm";
import { RankingPanel } from "./RankingPanel";
import { ScoutPanel } from "./ScoutPanel";
import { WAR_GUILDS, type WarGuild } from "./WarMapData";
import { WarMapHex } from "./WarMapHex";

interface WarMapProps {
  onNavigate: (path: string) => void;
}

interface ViewportSize {
  width: number;
  height: number;
}

const MIN_SCALE = 0.82;
const MAX_SCALE = 1.55;
const FOCUS_SCALE = 1.22;

export function WarMap({ onNavigate }: WarMapProps) {
  const [viewport, setViewport] = useState<ViewportSize>({ width: 0, height: 0 });
  const [scale, setScale] = useState(1);
  const [isRankingOpen, setIsRankingOpen] = useState(false);
  const [selectedGuild, setSelectedGuild] = useState<WarGuild | null>(null);
  const mapRef = useRef<HTMLDivElement>(null);
  const mapX = useMotionValue(0);
  const mapY = useMotionValue(0);
  const dragConstraints = getDragConstraints(viewport, scale);

  useEffect(() => {
    const updateViewport = () => {
      setViewport({ width: window.innerWidth, height: window.innerHeight });
    };

    updateViewport();
    window.addEventListener("resize", updateViewport);

    return () => window.removeEventListener("resize", updateViewport);
  }, []);

  useEffect(() => {
    if (viewport.width === 0 || viewport.height === 0) return;

    mapX.set(0);
    mapY.set(0);
  }, [mapX, mapY, viewport.height, viewport.width]);

  useEffect(() => {
    mapX.set(clampValue(mapX.get(), dragConstraints.left, dragConstraints.right));
    mapY.set(clampValue(mapY.get(), dragConstraints.top, dragConstraints.bottom));
  }, [
    dragConstraints.bottom,
    dragConstraints.left,
    dragConstraints.right,
    dragConstraints.top,
    mapX,
    mapY,
  ]);

  const focusGuild = (guild: WarGuild) => {
    const nextScale = Math.max(scale, FOCUS_SCALE);
    const target = getFocusedPosition(guild, viewport, nextScale);

    setSelectedGuild(guild);
    setScale(nextScale);
    animate(mapX, target.x, { duration: 0.55, ease: steppedEase(10) });
    animate(mapY, target.y, { duration: 0.55, ease: steppedEase(10) });
  };

  const closeScout = () => {
    setSelectedGuild(null);
  };

  const handleZoom = (delta: number) => {
    setScale((currentScale) => clampValue(currentScale + delta, MIN_SCALE, MAX_SCALE));
  };

  const handleWheel = (event: WheelEvent<HTMLElement>) => {
    event.preventDefault();
    handleZoom(-event.deltaY * 0.0015);
  };

  const handleWorldPointerDown = (event: ReactPointerEvent<HTMLElement>) => {
    const target = event.target;
    if (!(target instanceof Element)) return;
    if (target.closest("[data-war-interactive='true']")) return;

    closeScout();
  };

  return (
    <main
      className="relative h-screen w-full overflow-hidden"
      onPointerDown={handleWorldPointerDown}
      onWheel={handleWheel}
      style={{
        background: "#03101c",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#fff8d7",
      }}
    >
      <GuildBgm src={AUDIO_ASSETS.bgm.guildTown} />

      <motion.div
        ref={mapRef}
        className="absolute left-0 top-0 h-screen w-screen cursor-grab active:cursor-grabbing"
        drag
        dragConstraints={dragConstraints}
        dragElastic={0.08}
        dragMomentum={false}
        style={{
          x: mapX,
          y: mapY,
          scale,
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
            onSelect={focusGuild}
          />
        ))}
      </motion.div>

      <RankingPanel isOpen={isRankingOpen} onToggle={() => setIsRankingOpen((isOpen) => !isOpen)} />
      <ScoutPanel guild={selectedGuild} onClose={closeScout} />
      <BackButton onNavigate={onNavigate} />

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

function getFocusedPosition(guild: WarGuild, viewport: ViewportSize, scale: number) {
  const mapWidth = viewport.width;
  const mapHeight = viewport.height;
  const targetX = (mapWidth * guild.x) / 100;
  const targetY = (mapHeight * guild.y) / 100;
  const constraints = getDragConstraints(viewport, scale);

  return {
    x: clampValue(viewport.width * 0.5 - targetX * scale, constraints.left, constraints.right),
    y: clampValue(viewport.height * 0.5 - targetY * scale, constraints.top, constraints.bottom),
  };
}

function getDragConstraints(viewport: ViewportSize, scale: number) {
  return {
    left: Math.min(0, viewport.width - viewport.width * scale),
    right: 0,
    top: Math.min(0, viewport.height - viewport.height * scale),
    bottom: 0,
  };
}

function clampValue(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max);
}
