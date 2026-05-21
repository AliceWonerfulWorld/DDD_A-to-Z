import { motion, type MotionValue, type PanInfo } from "framer-motion";
import type { CSSProperties, PointerEvent as ReactPointerEvent, RefObject } from "react";
import { steppedEase } from "../../lib/animationUtils";
import type { PlacedItem } from "./types";
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

export interface DeploymentPreview {
  id: string;
  isUnlocked: boolean;
  name: string;
  src: string;
  width: number;
  x: number;
  y: number;
}

interface TownMapProps {
  baseSrc: string;
  currentGuildLevel: number;
  deploymentPreview: DeploymentPreview | null;
  dragConstraints: { left: number; right: number; top: number; bottom: number };
  mapRef: RefObject<HTMLDivElement | null>;
  mapX: MotionValue<number>;
  mapY: MotionValue<number>;
  newlyDeployedItemId: string | null;
  onCancelDeployment: () => void;
  onCommitDeployment: () => void;
  onMoveDeploymentPreview: (event: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => void;
  onMoveItem: (
    item: PlacedItem,
    event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => void;
  onSelectItem: (id: string) => void;
  onStoreItem: (item: PlacedItem) => void;
  placedItems: PlacedItem[];
  scale: number;
  selectedPlacedItemId: string | null;
  stopNestedDrag: (event: ReactPointerEvent<HTMLElement>) => void;
  storingPlacedItemIds: string[];
  unlockClearingLevel: number | null;
}

export function TownMap({
  baseSrc,
  currentGuildLevel,
  deploymentPreview,
  dragConstraints,
  mapRef,
  mapX,
  mapY,
  newlyDeployedItemId,
  onCancelDeployment,
  onCommitDeployment,
  onMoveDeploymentPreview,
  onMoveItem,
  onSelectItem,
  onStoreItem,
  placedItems,
  scale,
  selectedPlacedItemId,
  stopNestedDrag,
  storingPlacedItemIds,
  unlockClearingLevel,
}: TownMapProps) {
  const isDeployMode = deploymentPreview !== null;
  const unlockRings = getTownUnlockRings();

  return (
    <motion.div
      ref={mapRef}
      className={`absolute left-0 top-0 h-[200vh] w-[200vw] ${isDeployMode ? "cursor-default" : "cursor-grab active:cursor-grabbing"}`}
      drag={!isDeployMode}
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
        src={baseSrc}
        alt=""
        aria-hidden="true"
        draggable={false}
        style={{
          position: "absolute",
          inset: 0,
          width: "100%",
          height: "100%",
          objectFit: "cover",
          objectPosition: "center bottom",
          pointerEvents: "none",
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(180deg, rgba(4, 18, 18, 0.1) 0%, rgba(5, 16, 12, 0.04) 48%, rgba(6, 15, 10, 0.3) 100%)",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      <div
        aria-hidden="true"
        className="bg-[radial-gradient(ellipse_at_center,_transparent_20%,_rgba(0,0,0,0.6)_100%)]"
        style={{
          position: "absolute",
          inset: 0,
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <LockedTownFog currentGuildLevel={currentGuildLevel} rings={unlockRings} />
      <LockedAreaLabels currentGuildLevel={currentGuildLevel} rings={unlockRings} />
      {unlockClearingLevel && (
        <UnlockClearingBurst
          level={unlockClearingLevel}
          radiusPercent={getTownUnlockRadiusPercent(unlockClearingLevel)}
        />
      )}

      {placedItems.map((item) => {
        const isSelected = selectedPlacedItemId === item.id;
        const isStoring = storingPlacedItemIds.includes(item.id);
        const isNewlyDeployed = newlyDeployedItemId === item.id;

        return (
          <motion.div
            key={item.id}
            initial={
              isNewlyDeployed
                ? {
                    opacity: 0,
                    scale: 0.82,
                    y: -96,
                  }
                : false
            }
            animate={{
              opacity: isStoring ? 0 : 1,
              scale: isStoring ? 0.72 : 1,
              y: isStoring ? -22 : 0,
            }}
            transition={{ duration: 0.24, ease: steppedEase(6) }}
            style={{
              position: "absolute",
              left: item.x,
              top: item.y,
              width: item.width,
              height: "fit-content",
              outline: isSelected ? "3px solid rgba(255, 217, 102, 0.82)" : "3px solid transparent",
              outlineOffset: "4px",
              pointerEvents: isStoring || isDeployMode ? "none" : "auto",
              transformOrigin: "50% 80%",
              zIndex: isSelected || isStoring ? 10 : 8,
            }}
          >
            <motion.img
              className="pixelated"
              src={item.src}
              alt={item.name}
              drag
              dragSnapToOrigin
              dragElastic={0}
              dragMomentum={false}
              onPointerDown={stopNestedDrag}
              onClick={() => onSelectItem(item.id)}
              onDragEnd={(event, info) => onMoveItem(item, event, info)}
              whileHover={{ scale: 1.02 }}
              whileDrag={{ scale: 1.05, zIndex: 12 }}
              style={{
                display: "block",
                width: "100%",
                height: "auto",
                cursor: "grab",
                touchAction: "none",
                userSelect: "none",
                filter: isSelected
                  ? "drop-shadow(10px 14px 0 rgba(0,0,0,0.3)) drop-shadow(0 0 12px rgba(255,217,102,0.72))"
                  : "drop-shadow(10px 14px 0 rgba(0,0,0,0.3))",
              }}
            />
            {isSelected && (
              <motion.button
                type="button"
                aria-label={`Store ${item.name}`}
                onPointerDown={stopNestedDrag}
                onClick={() => onStoreItem(item)}
                initial={{ opacity: 0, y: 6, scale: 0.92 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.18, ease: steppedEase(4) }}
                whileHover={{ y: -2, backgroundColor: "rgba(255, 217, 102, 0.22)" }}
                whileTap={{ y: 1, scale: 0.96 }}
                style={{
                  position: "absolute",
                  right: "-10px",
                  top: "-44px",
                  minWidth: "74px",
                  minHeight: "34px",
                  border: "2px solid rgba(255, 217, 102, 0.86)",
                  borderBottomColor: "rgba(96, 62, 22, 0.98)",
                  borderRightColor: "rgba(96, 62, 22, 0.98)",
                  background: "rgba(3, 10, 24, 0.9)",
                  boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
                  color: "#fff8d7",
                  cursor: "pointer",
                  fontFamily: "inherit",
                  fontSize: "0.52rem",
                  lineHeight: 1,
                  padding: "8px 9px",
                  textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
                  touchAction: "none",
                }}
              >
                STORE
              </motion.button>
            )}
          </motion.div>
        );
      })}

      {deploymentPreview && (
        <motion.div
          key={deploymentPreview.id}
          aria-label={`${deploymentPreview.name} deployment preview`}
          aria-roledescription="deployment preview"
          drag
          dragElastic={0}
          dragMomentum={false}
          dragSnapToOrigin
          onDragEnd={onMoveDeploymentPreview}
          onPointerDown={stopNestedDrag}
          initial={{ opacity: 0, scale: 0.92, y: -18 }}
          animate={{
            opacity: deploymentPreview.isUnlocked ? 0.68 : 0.42,
            scale: deploymentPreview.isUnlocked ? 1 : 0.96,
            y: 0,
          }}
          transition={{ duration: 0.18, ease: steppedEase(5) }}
          whileDrag={{ opacity: 0.76, scale: 1.03, zIndex: 13 }}
          style={{
            position: "absolute",
            left: deploymentPreview.x,
            top: deploymentPreview.y,
            width: deploymentPreview.width,
            cursor: "grab",
            filter: deploymentPreview.isUnlocked
              ? "drop-shadow(10px 14px 0 rgba(0,0,0,0.28)) drop-shadow(0 0 18px rgba(116,247,161,0.58))"
              : "grayscale(0.65) drop-shadow(10px 14px 0 rgba(0,0,0,0.36)) drop-shadow(0 0 18px rgba(255,77,109,0.5))",
            touchAction: "none",
            transformOrigin: "50% 80%",
            userSelect: "none",
            zIndex: 12,
          }}
        >
          <img
            className="pixelated"
            src={deploymentPreview.src}
            alt=""
            draggable={false}
            style={{
              display: "block",
              width: "100%",
              height: "auto",
              pointerEvents: "none",
            }}
          />
          <div
            aria-hidden="true"
            style={{
              position: "absolute",
              left: "50%",
              bottom: "-14px",
              width: "72%",
              height: "16px",
              transform: "translateX(-50%)",
              border: `2px solid ${deploymentPreview.isUnlocked ? "rgba(116,247,161,0.74)" : "rgba(255,77,109,0.72)"}`,
              background: deploymentPreview.isUnlocked
                ? "rgba(116,247,161,0.12)"
                : "rgba(255,77,109,0.12)",
              boxShadow: deploymentPreview.isUnlocked
                ? "0 0 14px rgba(116,247,161,0.34)"
                : "0 0 14px rgba(255,77,109,0.32)",
            }}
          />
          <div
            style={{
              position: "absolute",
              left: "50%",
              top: "calc(100% + 12px)",
              transform: "translateX(-50%)",
              border: `2px solid ${deploymentPreview.isUnlocked ? "rgba(116,247,161,0.78)" : "rgba(255,77,109,0.78)"}`,
              background: "rgba(3, 10, 24, 0.9)",
              boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
              color: deploymentPreview.isUnlocked ? "#74f7a1" : "#ff9aae",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.72rem",
              lineHeight: 1,
              padding: "6px 8px",
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
              whiteSpace: "nowrap",
            }}
          >
            {deploymentPreview.isUnlocked ? "BUILD READY" : "LOCKED AREA"}
          </div>
        </motion.div>
      )}

      {deploymentPreview && (
        <div
          onPointerDown={stopNestedDrag}
          style={{
            position: "absolute",
            left: deploymentPreview.x + deploymentPreview.width + 12,
            top: deploymentPreview.y,
            zIndex: 13,
            display: "grid",
            gap: "8px",
            width: "104px",
          }}
        >
          <button
            type="button"
            onClick={onCommitDeployment}
            style={deployActionButtonStyle(deploymentPreview.isUnlocked)}
          >
            BUILD
          </button>
          <button type="button" onClick={onCancelDeployment} style={deployActionButtonStyle(false)}>
            CANCEL
          </button>
        </div>
      )}
    </motion.div>
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
    2: { x: 62, y: 30 },
    3: { x: 76, y: 44 },
    4: { x: 30, y: 18 },
    5: { x: 83, y: 82 },
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

function deployActionButtonStyle(isPrimary: boolean): CSSProperties {
  return {
    minHeight: "34px",
    border: `2px solid ${isPrimary ? "rgba(116, 247, 161, 0.86)" : "rgba(255, 217, 102, 0.76)"}`,
    borderBottomColor: isPrimary ? "rgba(24, 83, 45, 0.95)" : "rgba(96, 62, 22, 0.98)",
    borderRightColor: isPrimary ? "rgba(24, 83, 45, 0.95)" : "rgba(96, 62, 22, 0.98)",
    background: isPrimary ? "rgba(4, 67, 37, 0.92)" : "rgba(3, 10, 24, 0.9)",
    boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
    color: isPrimary ? "#74f7a1" : "#ffd966",
    cursor: "pointer",
    fontFamily: '"Press Start 2P", "DotGothic16", monospace',
    fontSize: "0.46rem",
    lineHeight: 1,
    padding: "8px 9px",
    textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
  };
}
