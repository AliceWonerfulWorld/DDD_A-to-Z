import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";
import { SPRITE_ASSETS } from "../../constants/assets";
import type { PetSummary } from "../../features/pet/api";
import { steppedEase } from "../../lib/animationUtils";
import { GopherSprite } from "../shared/GopherSprite";

const gopherTalkLines = [
  "今日もコード日和！",
  "ギルド広場を巡回中。",
  "クエスト、行く？",
  "休憩もだいじ。",
  "CP、ためてこ！",
] as const;

const GOPHER_HITBOX_WIDTH = 132;
const SPRITE_FRAME_WIDTH = 192;
const SPRITE_FRAME_HEIGHT = 208;
const SPRITE_COLUMNS = 8;
const SPRITE_ROWS = 9;
const AUDIO_PANEL_SAFE_SPACE = 124;
const DEFAULT_STAGE_WIDTH = 960;

const homePetSpriteAssets: Record<string, string> = {
  python: SPRITE_ASSETS.PYTHON,
  rust: SPRITE_ASSETS.RUST,
  java: SPRITE_ASSETS.JAVA,
};

function createWalkPath(stageWidth: number) {
  const safeStageWidth = Math.max(stageWidth, 320);
  const minX = Math.min(AUDIO_PANEL_SAFE_SPACE, Math.max(16, safeStageWidth - GOPHER_HITBOX_WIDTH));
  const maxX = Math.max(minX, safeStageWidth - GOPHER_HITBOX_WIDTH - 16);
  const span = maxX - minX;

  return [0, 0.32, 0.2, 0.62, 0.78, 0.56, 0.28, 0].map((ratio) => Math.round(minX + span * ratio));
}

export function WalkingGopher({ onTalk, pet }: { onTalk: () => void; pet?: PetSummary | null }) {
  const gopherButtonRef = useRef<HTMLButtonElement | null>(null);
  const lastXRef = useRef<number | null>(null);
  const talkTimeoutRef = useRef<number | null>(null);
  const [stageWidth, setStageWidth] = useState(DEFAULT_STAGE_WIDTH);
  const [direction, setDirection] = useState<"right" | "left">("right");
  const [talkLine, setTalkLine] = useState<(typeof gopherTalkLines)[number] | null>(null);
  const [reactionCount, setReactionCount] = useState(0);
  const walkPathX = useMemo(() => createWalkPath(stageWidth), [stageWidth]);
  const petAttribute = pet?.attribute.toLowerCase() ?? "go";
  const walkRow = direction === "right" ? 1 : 2;
  const speechBubbleSide =
    direction === "right" ? { left: "104px", right: "auto" } : { left: "auto", right: "96px" };

  useEffect(() => {
    const stageElement = gopherButtonRef.current?.parentElement;
    if (!stageElement) {
      return;
    }

    const updateStageWidth = () => {
      setStageWidth(stageElement.getBoundingClientRect().width);
    };

    updateStageWidth();

    const resizeObserver = new ResizeObserver(updateStageWidth);
    resizeObserver.observe(stageElement);

    return () => {
      resizeObserver.disconnect();
    };
  }, []);

  useEffect(() => {
    return () => {
      if (talkTimeoutRef.current !== null) {
        window.clearTimeout(talkTimeoutRef.current);
      }
    };
  }, []);

  const reactToClick = () => {
    onTalk();

    const nextLine = gopherTalkLines[reactionCount % gopherTalkLines.length];
    setReactionCount((current) => current + 1);
    setTalkLine(nextLine);

    if (talkTimeoutRef.current !== null) {
      window.clearTimeout(talkTimeoutRef.current);
    }

    talkTimeoutRef.current = window.setTimeout(() => {
      setTalkLine(null);
      talkTimeoutRef.current = null;
    }, 2600);
  };

  return (
    <motion.button
      ref={gopherButtonRef}
      type="button"
      aria-label="ペットに話しかける"
      initial={false}
      animate={{
        x: walkPathX,
        y: ["0px", "-10px", "6px", "-14px", "2px", "16px", "8px", "0px"],
        scale: [0.92, 0.88, 0.94, 0.86, 0.9, 0.96, 0.94, 0.92],
      }}
      transition={{
        duration: 26,
        repeat: Infinity,
        ease: "easeInOut",
        times: [0, 0.16, 0.28, 0.44, 0.58, 0.72, 0.88, 1],
      }}
      onUpdate={(latest) => {
        const currentX =
          typeof latest.x === "number" ? latest.x : Number.parseFloat(String(latest.x));
        const lastX = lastXRef.current;

        if (!Number.isFinite(currentX)) {
          return;
        }

        if (lastX !== null) {
          const deltaX = currentX - lastX;
          if (Math.abs(deltaX) > 0.02) {
            setDirection(deltaX > 0 ? "right" : "left");
          }
        }

        lastXRef.current = currentX;
      }}
      style={{
        position: "absolute",
        left: 0,
        bottom: "clamp(6px, 2vh, 18px)",
        width: "92px",
        height: "100px",
        border: 0,
        background: "transparent",
        cursor: "pointer",
        font: "inherit",
        padding: 0,
        zIndex: 4,
      }}
      onClick={reactToClick}
    >
      <AnimatePresence>
        {talkLine && (
          <motion.div
            key={talkLine}
            initial={{ opacity: 0, y: 8, scale: 0.94 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 8, scale: 0.94 }}
            style={{
              position: "absolute",
              ...speechBubbleSide,
              bottom: "46px",
              width: "max-content",
              maxWidth: "min(220px, 34vw)",
              border: "2px solid rgba(255, 215, 0, 0.82)",
              borderBottomColor: "rgba(111, 79, 28, 0.95)",
              borderRightColor: "rgba(111, 79, 28, 0.95)",
              background: "rgba(3, 10, 24, 0.9)",
              boxShadow: "0 0 0 2px rgba(0,0,0,0.72), 4px 4px 0 rgba(0,0,0,0.42)",
              color: "#fff8d7",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.78rem",
              lineHeight: 1.5,
              letterSpacing: "0.04em",
              padding: "8px 10px",
              textAlign: "center",
              whiteSpace: "normal",
              pointerEvents: "none",
              zIndex: 6,
            }}
          >
            {talkLine}
          </motion.div>
        )}
      </AnimatePresence>
      <motion.div
        animate={{ y: [0, -2, 0] }}
        transition={{ duration: 0.45, repeat: Infinity, ease: steppedEase(3) }}
        style={{
          position: "absolute",
          left: 0,
          bottom: 8,
          width: "132px",
          height: "143px",
          transform: "scale(0.62)",
          transformOrigin: "left bottom",
        }}
      >
        <motion.div
          key={reactionCount}
          animate={
            reactionCount === 0
              ? {}
              : {
                  y: [0, -20, 0],
                  rotate: [0, -5, 5, 0],
                }
          }
          transition={{ duration: 0.42, ease: steppedEase(5) }}
        >
          <HomePetSprite attribute={petAttribute} row={walkRow} />
        </motion.div>
      </motion.div>
      <motion.div
        animate={{ scaleX: [1, 0.86, 1], opacity: [0.34, 0.24, 0.34] }}
        transition={{ duration: 0.45, repeat: Infinity, ease: steppedEase(3) }}
        style={{
          position: "absolute",
          left: "18px",
          bottom: "2px",
          width: "62px",
          height: "10px",
          background: "rgba(0,0,0,0.42)",
          filter: "blur(1px)",
        }}
      />
    </motion.button>
  );
}

function HomePetSprite({ attribute, row }: { attribute: string; row: number }) {
  if (attribute === "go") {
    return <GopherSprite frameCount={8} row={row} />;
  }
  const spriteAsset = homePetSpriteAssets[attribute];
  if (spriteAsset) {
    return <LanguageHomePetSprite asset={spriteAsset} row={row} />;
  }

  const label = attribute.slice(0, 2).toUpperCase();
  return (
    <div
      style={{
        display: "grid",
        width: "132px",
        height: "143px",
        placeItems: "center",
        border: "4px solid rgba(255, 248, 215, 0.62)",
        background:
          "linear-gradient(180deg, rgba(116, 247, 161, 0.22), rgba(0, 245, 255, 0.14)), rgba(20, 38, 54, 0.74)",
        color: "#fff8d7",
        fontSize: "1.8rem",
        boxShadow: "0 0 0 2px rgba(0, 0, 0, 0.7), inset 0 0 22px rgba(255, 255, 255, 0.08)",
      }}
    >
      {label}
    </div>
  );
}

function LanguageHomePetSprite({ asset, row }: { asset: string; row: number }) {
  const scale = 132 / SPRITE_FRAME_WIDTH;
  const displaySheetWidth = Math.round(SPRITE_FRAME_WIDTH * SPRITE_COLUMNS * scale);
  const displaySheetHeight = Math.round(SPRITE_FRAME_HEIGHT * SPRITE_ROWS * scale);
  const frameStep = Math.round(SPRITE_FRAME_WIDTH * scale);
  const rowOffsetY = Math.round(SPRITE_FRAME_HEIGHT * row * scale);

  return (
    <motion.div
      animate={{ backgroundPositionX: ["0px", `-${frameStep * 8}px`] }}
      transition={{ duration: 0.9, repeat: Infinity, ease: steppedEase(8) }}
      style={{
        width: "132px",
        height: "143px",
        backgroundImage: `url(${asset})`,
        backgroundPositionY: `-${rowOffsetY}px`,
        backgroundRepeat: "no-repeat",
        backgroundSize: `${displaySheetWidth}px ${displaySheetHeight}px`,
        imageRendering: "pixelated",
      }}
    />
  );
}
