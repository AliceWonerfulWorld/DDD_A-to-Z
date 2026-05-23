import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import { getTownUnlockRadiusPercent, getTownUnlockRings } from "./townUnlock";

const LOCKED_LEVEL_COLORS: Record<number, { fog: string; line: string; text: string }> = {
  2: {
    fog: "rgba(86, 170, 255, 0.34)",
    line: "rgba(86, 170, 255, 0.72)",
    text: "#8ed2ff",
  },
  3: {
    fog: "rgba(154, 111, 255, 0.34)",
    line: "rgba(154, 111, 255, 0.72)",
    text: "#c5a7ff",
  },
  4: {
    fog: "rgba(255, 183, 77, 0.33)",
    line: "rgba(255, 183, 77, 0.72)",
    text: "#ffd18a",
  },
  5: {
    fog: "rgba(255, 91, 137, 0.34)",
    line: "rgba(255, 91, 137, 0.72)",
    text: "#ff9fbc",
  },
};
const DEFAULT_LOCKED_LEVEL_COLOR = {
  fog: "rgba(184, 208, 255, 0.34)",
  line: "rgba(184, 208, 255, 0.7)",
  text: "#b8d0ff",
};

interface TownMapLockedOverlayProps {
  currentGuildLevel: number;
  unlockClearingLevel: number | null;
}

export function TownMapLockedOverlay({
  currentGuildLevel,
  unlockClearingLevel,
}: TownMapLockedOverlayProps) {
  const unlockRings = getTownUnlockRings();

  return (
    <>
      <LockedTownFog currentGuildLevel={currentGuildLevel} rings={unlockRings} />
      <LockedAreaLabels currentGuildLevel={currentGuildLevel} rings={unlockRings} />
      {unlockClearingLevel && (
        <UnlockClearingBurst
          level={unlockClearingLevel}
          radiusPercent={getTownUnlockRadiusPercent(unlockClearingLevel)}
        />
      )}
    </>
  );
}

function LockedTownFog({
  currentGuildLevel,
  rings,
}: {
  currentGuildLevel: number;
  rings: ReturnType<typeof getTownUnlockRings>;
}) {
  return (
    <div
      aria-hidden="true"
      style={{
        position: "absolute",
        inset: 0,
        zIndex: 5,
        background: getLockedFogGradient(currentGuildLevel, rings),
        backdropFilter: "blur(1.4px) saturate(0.72)",
        pointerEvents: "none",
      }}
    />
  );
}

function LockedAreaLabels({
  currentGuildLevel,
  rings,
}: {
  currentGuildLevel: number;
  rings: ReturnType<typeof getTownUnlockRings>;
}) {
  return (
    <div
      aria-hidden="true"
      style={{
        position: "absolute",
        inset: 0,
        zIndex: 6,
        pointerEvents: "none",
      }}
    >
      {rings
        .filter((ring) => ring.level > currentGuildLevel)
        .map((ring) => {
          const labelPosition = getLockedAreaLabelPosition(ring.level);
          const lockedColor = getLockedLevelColor(ring.level);

          return (
            <span
              key={ring.level}
              style={{
                position: "absolute",
                left: `${labelPosition.x}%`,
                top: `${labelPosition.y}%`,
                border: `2px solid ${lockedColor.line}`,
                background: "rgba(3, 10, 24, 0.76)",
                boxShadow: `0 0 0 2px rgba(0,0,0,0.54), 0 0 18px ${lockedColor.line}`,
                color: lockedColor.text,
                fontFamily: '"Press Start 2P", "DotGothic16", monospace',
                fontSize: "1rem",
                letterSpacing: 0,
                lineHeight: 1.25,
                padding: "8px 10px",
                textAlign: "center",
                textShadow: `2px 2px 0 rgba(0,0,0,0.88), 0 0 10px ${lockedColor.line}, 0 0 18px rgba(0,0,0,0.8)`,
                transform: "translate(-50%, -50%)",
                whiteSpace: "nowrap",
              }}
            >
              LV.{ring.level}
              <br />
              UNLOCK
            </span>
          );
        })}
    </div>
  );
}

function getLockedFogGradient(
  currentGuildLevel: number,
  rings: ReturnType<typeof getTownUnlockRings>,
) {
  const lockedRings = rings.filter((ring) => ring.level > currentGuildLevel);
  const currentRadius =
    rings.find((ring) => ring.level === currentGuildLevel)?.radiusPercent ??
    getTownUnlockRadiusPercent(currentGuildLevel);

  if (lockedRings.length === 0) {
    return "radial-gradient(ellipse at center, transparent 0%, transparent 100%)";
  }

  const stops = [
    "transparent 0%",
    `transparent ${currentRadius}%`,
    `${getLockedLevelColor(lockedRings[0].level).fog} ${currentRadius + 1}%`,
  ];

  for (const ring of lockedRings) {
    const color = getLockedLevelColor(ring.level).fog;
    const nextRing = lockedRings.find((candidate) => candidate.level === ring.level + 1);
    const endRadius = ring.radiusPercent;

    stops.push(`${color} ${Math.max(currentRadius + 1, endRadius)}%`);

    if (nextRing) {
      stops.push(`${getLockedLevelColor(nextRing.level).fog} ${endRadius + 1}%`);
    } else {
      stops.push("rgba(11, 20, 32, 0.72) 100%");
    }
  }

  return `radial-gradient(ellipse at center, ${stops.join(", ")})`;
}

function getLockedAreaLabelPosition(level: number) {
  const positionByLevel: Record<number, { x: number; y: number }> = {
    2: { x: 74, y: 36 }, // Light blue ring (radius 32-46)
    3: { x: 17, y: 69 }, // Purple ring (radius 46-60)
    4: { x: 16, y: 16 }, // Orange ring (radius 60-74)
    5: { x: 90, y: 90 }, // Pink ring (radius 74-88)
  };

  return positionByLevel[level] ?? { x: 76, y: 50 };
}

function getLockedLevelColor(level: number) {
  return LOCKED_LEVEL_COLORS[level] ?? DEFAULT_LOCKED_LEVEL_COLOR;
}

function UnlockClearingBurst({ level, radiusPercent }: { level: number; radiusPercent: number }) {
  return (
    <motion.div
      aria-hidden="true"
      initial={{ opacity: 0.95, scale: 0.9 }}
      animate={{ opacity: 0, scale: 1.16 }}
      transition={{ duration: 1.75, ease: steppedEase(12) }}
      style={{
        position: "absolute",
        left: "50%",
        top: "50%",
        zIndex: 9,
        width: `${(radiusPercent * 2) / 1.08}%`,
        height: `${(radiusPercent * 2) / 0.82}%`,
        transform: "translate(-50%, -50%)",
        borderRadius: "50%",
        background:
          "radial-gradient(ellipse at center, rgba(255,255,255,0.42) 0%, rgba(116,247,161,0.28) 38%, rgba(184,208,255,0.5) 66%, rgba(184,208,255,0) 100%)",
        boxShadow: "0 0 34px rgba(116,247,161,0.52), 0 0 72px rgba(184,208,255,0.4)",
        pointerEvents: "none",
      }}
    >
      <motion.span
        initial={{ opacity: 0, y: 10, scale: 0.92 }}
        animate={{ opacity: [0, 1, 1, 0], y: [10, 0, 0, -8], scale: [0.92, 1, 1, 1] }}
        transition={{ duration: 1.55, ease: steppedEase(10) }}
        style={{
          position: "absolute",
          left: "50%",
          top: "18%",
          transform: "translateX(-50%)",
          border: "3px solid rgba(116,247,161,0.86)",
          background: "rgba(1, 12, 24, 0.92)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 0 0 22px rgba(116,247,161,0.48)",
          color: "#74f7a1",
          fontFamily: '"Press Start 2P", "DotGothic16", monospace',
          fontSize: "0.58rem",
          lineHeight: 1,
          padding: "8px 10px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.78)",
          whiteSpace: "nowrap",
        }}
      >
        LV.{level} AREA OPEN
      </motion.span>
    </motion.div>
  );
}
