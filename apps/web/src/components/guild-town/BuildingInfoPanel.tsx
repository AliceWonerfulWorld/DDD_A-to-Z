import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import { steppedEase } from "../../lib/animationUtils";
import { BUILDING_MASTERS } from "./townData";
import type { BuildingBuffType, GuildSpLanguage, PlacedItem } from "./types";

interface BuildingInfoPanelProps {
  item: PlacedItem | null;
  onClose: () => void;
  onUpgradeBuilding: (placedItemId: string) => void;
  userCp: number;
  userGuildSp: number;
}

export function BuildingInfoPanel({
  item,
  onClose,
  onUpgradeBuilding,
  userCp,
  userGuildSp,
}: BuildingInfoPanelProps) {
  const [viewportWidth, setViewportWidth] = useState(() => getViewportWidth());

  useEffect(() => {
    const updateViewportWidth = () => setViewportWidth(getViewportWidth());

    updateViewportWidth();
    window.addEventListener("resize", updateViewportWidth);

    return () => window.removeEventListener("resize", updateViewportWidth);
  }, []);

  if (!item) return null;

  const isCompactLayout = viewportWidth < 760;
  const building = item.buildingId
    ? BUILDING_MASTERS.find((buildingMaster) => buildingMaster.id === item.buildingId)
    : undefined;
  const currentLevelIndex = Math.min(Math.max(item.level, 1), 5) - 1;
  const currentLevel = building?.levels[currentLevelIndex];
  const nextLevel = building?.levels[currentLevelIndex + 1];
  const isMaxLevel = Boolean(building && !nextLevel);
  const isCpShort = Boolean(nextLevel && userCp < nextLevel.upgradeCostCp);
  const isSpShort = Boolean(nextLevel && userGuildSp < nextLevel.upgradeCostSp);
  const canUpgrade = Boolean(building && nextLevel && !isCpShort && !isSpShort);
  const targetSpLanguage = building?.targetSpLanguage ?? "Common";

  return (
    <motion.section
      key={item.id}
      initial={{ opacity: 0, y: 18, scale: 0.98 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: 12, scale: 0.98 }}
      transition={{ duration: 0.28, ease: steppedEase(6) }}
      aria-live="polite"
      style={{
        position: "fixed",
        left: "clamp(16px, 5.4vw, 84px)",
        right: isCompactLayout ? "16px" : "clamp(78px, 7vw, 112px)",
        bottom: "calc(env(safe-area-inset-bottom, 0px) + 8px)",
        zIndex: 8,
        display: "grid",
        gridTemplateColumns: isCompactLayout
          ? "72px minmax(0, 1fr)"
          : "72px minmax(260px, 1fr) minmax(260px, 0.68fr)",
        minHeight: isCompactLayout ? "auto" : "88px",
        alignItems: "center",
        gap: isCompactLayout ? "10px 12px" : "12px",
        border: "3px solid rgba(255, 248, 215, 0.82)",
        borderBottomColor: "rgba(55, 44, 35, 0.98)",
        borderRightColor: "rgba(55, 44, 35, 0.98)",
        background: "linear-gradient(180deg, rgba(4, 10, 22, 0.94), rgba(3, 7, 14, 0.9))",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.76), 7px 7px 0 rgba(0,0,0,0.36), inset 0 0 22px rgba(116,247,161,0.09)",
        color: "#fff8d7",
        padding: "10px 12px",
        backdropFilter: "blur(2px)",
      }}
    >
      <motion.button
        type="button"
        aria-label="Close building info"
        onClick={onClose}
        whileHover={{ y: -1, backgroundColor: "rgba(255, 217, 102, 0.16)" }}
        whileTap={{ y: 1, scale: 0.96 }}
        style={{
          position: "absolute",
          right: "8px",
          top: "8px",
          width: "28px",
          height: "28px",
          border: "2px solid rgba(255, 217, 102, 0.72)",
          borderBottomColor: "rgba(96, 62, 22, 0.95)",
          borderRightColor: "rgba(96, 62, 22, 0.95)",
          background: "rgba(3, 10, 24, 0.78)",
          boxShadow: "0 0 0 1px rgba(0,0,0,0.62)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "0.64rem",
          lineHeight: 1,
          padding: 0,
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        x
      </motion.button>

      <div
        aria-hidden="true"
        style={{
          display: "grid",
          width: "72px",
          height: "64px",
          placeItems: "center",
          border: "2px solid rgba(116, 247, 161, 0.58)",
          background: "rgba(1, 12, 24, 0.72)",
          boxShadow: "inset 0 0 14px rgba(0,0,0,0.68)",
        }}
      >
        <img
          className="pixelated"
          src={item.src}
          alt=""
          draggable={false}
          style={{
            display: "block",
            maxWidth: "56px",
            maxHeight: "52px",
            filter: "drop-shadow(4px 5px 0 rgba(0,0,0,0.34))",
          }}
        />
      </div>
      <div style={{ minWidth: 0 }}>
        <p
          style={{
            margin: "0 0 4px",
            color: "#74f7a1",
            fontSize: "0.52rem",
            lineHeight: 1.4,
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          BUILDING DATA
        </p>
        <h2
          style={{
            margin: "0 0 5px",
            color: "#ffd966",
            fontSize: "clamp(0.72rem, 1.6vw, 0.95rem)",
            lineHeight: 1.5,
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          {item.title}
        </h2>
        <p
          style={{
            margin: 0,
            color: "#f4ecd0",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(0.76rem, 1.18vw, 0.92rem)",
            lineHeight: 1.38,
          }}
        >
          {item.description}
        </p>
      </div>
      {building && currentLevel ? (
        <div
          style={{
            display: "grid",
            gap: "7px",
            gridColumn: isCompactLayout ? "1 / -1" : undefined,
            minWidth: 0,
            paddingRight: "20px",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "8px",
              minWidth: 0,
              flexWrap: "wrap",
            }}
          >
            <motion.span
              key={`${item.id}-${item.level}`}
              initial={{ scale: 1.34, y: -3 }}
              animate={{ scale: 1, y: 0 }}
              transition={{ type: "spring", stiffness: 520, damping: 14 }}
              style={{
                display: "inline-grid",
                minWidth: "86px",
                placeItems: "center",
                border: "2px solid rgba(116, 247, 161, 0.9)",
                background: "rgba(1, 22, 32, 0.84)",
                boxShadow: "0 0 16px rgba(116,247,161,0.38), inset 0 0 12px rgba(116,247,161,0.14)",
                color: "#74f7a1",
                fontSize: "0.62rem",
                lineHeight: 1,
                padding: "7px 8px",
                textShadow: "0 0 8px rgba(116,247,161,0.86), 2px 2px 0 rgba(0,0,0,0.72)",
              }}
            >
              [ LV.{currentLevel.level} ]
            </motion.span>
            <span
              style={{
                color: isMaxLevel ? "#ffd966" : "#5cc8ff",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.76rem",
                lineHeight: 1.3,
                overflowWrap: "anywhere",
                textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
              }}
            >
              {isMaxLevel ? "MAX OUTPUT" : `Next Lv.${nextLevel?.level}`}
            </span>
          </div>

          <div
            style={{
              display: "grid",
              gap: "4px",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.72rem",
              lineHeight: 1.3,
            }}
          >
            <span style={{ color: "#f4ecd0" }}>
              現在の効果: {formatBuffEffect(building.buffType, currentLevel.buffValue)}
            </span>
            {nextLevel ? (
              <span style={{ color: "#5cc8ff" }}>
                -&gt; Next Lv.{nextLevel.level}:{" "}
                {formatBuffEffect(building.buffType, nextLevel.buffValue)}
              </span>
            ) : (
              <span style={{ color: "#ffd966" }}>強化上限に到達しています。</span>
            )}
          </div>

          {nextLevel ? (
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "minmax(0, 1fr) minmax(0, 1fr)",
                gap: "6px",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.62rem",
                lineHeight: 1.25,
              }}
            >
              <UpgradeCostPill
                isShort={isCpShort}
                label={`${nextLevel.upgradeCostCp.toLocaleString()} CP`}
                tone="#ffd966"
              />
              <UpgradeCostPill
                isShort={isSpShort}
                label={`${nextLevel.upgradeCostSp.toLocaleString()} ${targetSpLanguage}-SP`}
                tone={getLanguageTone(targetSpLanguage)}
              />
            </div>
          ) : null}

          <motion.button
            type="button"
            disabled={!canUpgrade}
            onClick={() => onUpgradeBuilding(item.id)}
            whileHover={
              canUpgrade ? { y: -1, backgroundColor: "rgba(4, 83, 46, 0.96)" } : undefined
            }
            whileTap={canUpgrade ? { y: 1, scale: 0.98 } : undefined}
            style={{
              width: "100%",
              border: `2px solid ${canUpgrade ? "rgba(116, 247, 161, 0.82)" : "rgba(92, 92, 92, 0.78)"}`,
              borderBottomColor: canUpgrade ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
              borderRightColor: canUpgrade ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
              background: canUpgrade ? "rgba(4, 67, 37, 0.9)" : "rgba(37, 37, 37, 0.92)",
              boxShadow: canUpgrade
                ? "0 0 14px rgba(116,247,161,0.24), inset 0 0 10px rgba(116,247,161,0.1)"
                : "inset 0 0 10px rgba(0,0,0,0.52)",
              color: canUpgrade ? "#74f7a1" : "rgba(255, 248, 215, 0.38)",
              cursor: canUpgrade ? "pointer" : "not-allowed",
              fontFamily: "inherit",
              fontSize: "0.54rem",
              lineHeight: 1,
              padding: "8px 8px",
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            }}
          >
            {isMaxLevel ? "MAX LEVEL" : "UPGRADE"}
          </motion.button>
        </div>
      ) : null}
    </motion.section>
  );
}

interface UpgradeCostPillProps {
  isShort: boolean;
  label: string;
  tone: string;
}

function UpgradeCostPill({ isShort, label, tone }: UpgradeCostPillProps) {
  return (
    <span
      className={isShort ? "text-red-500" : undefined}
      style={{
        display: "grid",
        minWidth: 0,
        placeItems: "center",
        overflow: "hidden",
        border: `2px solid ${isShort ? "#ef4444" : tone}`,
        background: "rgba(1, 8, 22, 0.76)",
        boxShadow: `inset 0 0 10px rgba(0,0,0,0.58), 0 0 10px ${isShort ? "#ef444455" : `${tone}44`}`,
        color: isShort ? "#ef4444" : tone,
        padding: "6px 5px",
        textAlign: "center",
        overflowWrap: "anywhere",
      }}
      title={label}
    >
      {label}
    </span>
  );
}

function formatBuffEffect(buffType: BuildingBuffType, buffValue: number) {
  const labelByType: Record<BuildingBuffType, string> = {
    arena: "精進CP",
    caffeine: "カフェインCP",
    commit: "GitHubコミットCP",
    core: "基本CP",
    daily: "ログインCP",
    interest: "未消費CP利息",
    night: "深夜作業CP",
    plant: "建築CPコスト削減",
    refactor: "レビュー改善CP",
    spBoost: "全言語SP",
    tower: "連続コミットCP",
  };

  if (buffType === "daily") {
    return `${labelByType[buffType]} +${buffValue.toLocaleString()}`;
  }

  if (buffType === "interest") {
    return `${labelByType[buffType]} +${formatPercent(buffValue, 2)}`;
  }

  if (buffType === "plant") {
    return `${labelByType[buffType]} -${formatPercent(buffValue)}`;
  }

  return `${labelByType[buffType]} +${formatPercent(buffValue)}`;
}

function formatPercent(value: number, maximumFractionDigits = 0) {
  return `${(value * 100).toLocaleString(undefined, {
    maximumFractionDigits,
  })}%`;
}

function getLanguageTone(language: GuildSpLanguage) {
  const toneByLanguage: Record<GuildSpLanguage, string> = {
    Common: "#ffd966",
    Go: "#00add8",
    Haskell: "#8f6bd8",
    Java: "#f97316",
    Python: "#f7df1e",
    Rust: "#ff7a1a",
    TypeScript: "#5cc8ff",
    Zig: "#f7a41d",
  };

  return toneByLanguage[language];
}

function getViewportWidth() {
  if (typeof window === "undefined") {
    return 1024;
  }

  return window.innerWidth;
}
