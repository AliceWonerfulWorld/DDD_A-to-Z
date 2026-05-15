import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import { steppedEase } from "../../lib/animationUtils";
import { GopherSprite } from "../shared/GopherSprite";

const gopherTalkLines = [
  "今日もコード日和！",
  "ギルド広場を巡回中。",
  "クエスト、行く？",
  "休憩もだいじ。",
  "CP、ためてこ！",
] as const;

export function WalkingGopher({ onTalk }: { onTalk: () => void }) {
  const lastXRef = useRef<number | null>(null);
  const talkTimeoutRef = useRef<number | null>(null);
  const [direction, setDirection] = useState<"right" | "left">("right");
  const [talkLine, setTalkLine] = useState<(typeof gopherTalkLines)[number] | null>(null);
  const [reactionCount, setReactionCount] = useState(0);
  const walkRow = direction === "right" ? 1 : 2;
  const speechBubbleSide =
    direction === "right" ? { left: "104px", right: "auto" } : { left: "auto", right: "96px" };

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
      type="button"
      aria-label="Gopher君に話しかける"
      initial={false}
      animate={{
        x: ["16vw", "30vw", "24vw", "50vw", "62vw", "46vw", "28vw", "16vw"],
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
          <GopherSprite frameCount={8} row={walkRow} />
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
