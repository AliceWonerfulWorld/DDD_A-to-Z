import { motion } from "framer-motion";
import { useMemo, useState } from "react";
import type {
  PointerEvent as ReactPointerEvent,
  RefObject,
  WheelEvent as ReactWheelEvent,
} from "react";
import { steppedEase } from "../../lib/animationUtils";
import { BUILDING_MASTERS } from "./townData";
import type { BuildingMaster, BuildingTargetSpLanguage, UserInventoryState } from "./types";

type BuildInventoryTab = "shop" | "inventory";

interface BuildInventoryProps {
  currentGuildLevel: number;
  currentGuildSpLanguage: BuildingTargetSpLanguage;
  inventory: UserInventoryState[];
  inventoryRef: RefObject<HTMLDivElement | null>;
  onBuyBuilding: (building: BuildingMaster) => void;
  onDeployBuilding: (building: BuildingMaster) => void;
  onToggleVisible: () => void;
  stopNestedDrag: (event: ReactPointerEvent<HTMLElement>) => void;
  userCp: number;
  userSp: number;
  visible: boolean;
}

const languageStyles: Record<BuildingTargetSpLanguage, { color: string; label: string }> = {
  Common: { color: "#ffd966", label: "COM" },
  Go: { color: "#00add8", label: "GO" },
  Java: { color: "#f97316", label: "JAVA" },
  Python: { color: "#f7df1e", label: "PY" },
  Rust: { color: "#ff7a1a", label: "RS" },
  TypeScript: { color: "#5cc8ff", label: "TS" },
};

export function BuildInventory({
  currentGuildLevel,
  currentGuildSpLanguage,
  inventory,
  inventoryRef,
  onBuyBuilding,
  onDeployBuilding,
  onToggleVisible,
  stopNestedDrag,
  userCp,
  userSp,
  visible,
}: BuildInventoryProps) {
  const [activeTab, setActiveTab] = useState<BuildInventoryTab>("shop");
  const inventoryCountByBuildingId = useMemo(
    () =>
      inventory.reduce<Record<string, number>>((countMap, item) => {
        countMap[item.buildingId] = item.count;
        return countMap;
      }, {}),
    [inventory],
  );
  const ownedInventoryTotal = inventory.reduce((total, item) => total + item.count, 0);

  const stopInventoryWheel = (event: ReactWheelEvent<HTMLElement>) => {
    event.stopPropagation();
  };

  return (
    <motion.aside
      ref={inventoryRef}
      onWheel={stopInventoryWheel}
      initial={{ opacity: 0, x: -18 }}
      animate={{
        opacity: 1,
        x: visible ? 0 : "calc(-100% - 14px)",
      }}
      transition={{ duration: 0.32, ease: steppedEase(6) }}
      aria-label="Build inventory"
      style={{
        position: "fixed",
        left: "clamp(14px, 2vw, 24px)",
        top: "calc(env(safe-area-inset-top, 0px) + 94px)",
        zIndex: 8,
        display: "flex",
        width: "clamp(340px, 28vw, 420px)",
        maxHeight:
          "calc(100vh - env(safe-area-inset-top, 0px) - env(safe-area-inset-bottom, 0px) - 118px)",
        alignItems: "stretch",
        flexDirection: "column",
        gap: "10px",
        overflow: "visible",
        border: visible ? "3px solid rgba(255, 248, 215, 0.8)" : "3px solid transparent",
        borderBottomColor: visible ? "rgba(55, 44, 35, 0.98)" : "transparent",
        borderRightColor: visible ? "rgba(55, 44, 35, 0.98)" : "transparent",
        background: visible ? "rgba(3, 7, 14, 0.9)" : "transparent",
        boxShadow: visible
          ? "0 0 0 2px rgba(0,0,0,0.72), 6px 6px 0 rgba(0,0,0,0.4), inset 0 0 18px rgba(255,248,215,0.08)"
          : "none",
        padding: "10px",
        backdropFilter: "blur(2px)",
      }}
    >
      <motion.button
        type="button"
        aria-label={visible ? "Hide build inventory" : "Show build inventory"}
        aria-expanded={visible}
        onPointerDown={stopNestedDrag}
        onClick={onToggleVisible}
        whileHover={{ x: 2, backgroundColor: "rgba(255, 217, 102, 0.18)" }}
        whileTap={{ x: -1, scale: 0.98 }}
        style={{
          position: "absolute",
          right: "-48px",
          top: "12px",
          width: "42px",
          height: "42px",
          border: "2px solid rgba(255, 217, 102, 0.78)",
          borderBottomColor: "rgba(96, 62, 22, 0.95)",
          borderRightColor: "rgba(96, 62, 22, 0.95)",
          background: "rgba(3, 10, 24, 0.86)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 4px 4px 0 rgba(0,0,0,0.34)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "0.82rem",
          lineHeight: 1,
          padding: "0",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        {visible ? "<<" : ">>"}
      </motion.button>

      {visible && (
        <>
          <div
            role="tablist"
            aria-label="Build inventory tabs"
            onPointerDown={stopNestedDrag}
            style={{
              display: "grid",
              gridTemplateColumns: "1fr 1fr",
              gap: "7px",
            }}
          >
            <InventoryTabButton
              active={activeTab === "shop"}
              label="SHOP"
              onClick={() => setActiveTab("shop")}
            />
            <InventoryTabButton
              active={activeTab === "inventory"}
              badge={ownedInventoryTotal}
              label="INVENTORY"
              onClick={() => setActiveTab("inventory")}
            />
          </div>

          <div
            style={{
              border: "2px solid rgba(116, 247, 161, 0.58)",
              borderBottomColor: "rgba(24, 83, 45, 0.95)",
              borderRightColor: "rgba(24, 83, 45, 0.95)",
              background: "rgba(1, 12, 24, 0.72)",
              boxShadow: "inset 0 0 12px rgba(0,0,0,0.62)",
              color: "#fff8d7",
              padding: "9px 10px",
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            }}
          >
            <p
              style={{
                margin: "0 0 7px",
                color: "#74f7a1",
                fontSize: "0.52rem",
                lineHeight: 1.4,
              }}
            >
              {activeTab === "shop" ? "BUILD SHOP" : "BUILD INVENTORY"}
            </p>
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "1fr 1fr",
                gap: "6px",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.78rem",
                lineHeight: 1.25,
              }}
            >
              <span>GUILD LV.{currentGuildLevel}</span>
              <span style={{ color: "#ffd966", textAlign: "right" }}>
                {activeTab === "shop"
                  ? `${userCp.toLocaleString()} CP`
                  : `${ownedInventoryTotal.toLocaleString()} ITEMS`}
              </span>
            </div>
          </div>

          <div
            onPointerDown={stopNestedDrag}
            style={{
              display: "flex",
              minHeight: 0,
              flexDirection: "column",
              gap: "10px",
              overflowX: "hidden",
              overflowY: "auto",
              paddingRight: "10px",
            }}
          >
            {activeTab === "shop"
              ? BUILDING_MASTERS.map((item) => (
                  <BuildingInventoryCard
                    key={item.id}
                    currentGuildLevel={currentGuildLevel}
                    currentGuildSpLanguage={currentGuildSpLanguage}
                    item={item}
                    onBuy={onBuyBuilding}
                    userCp={userCp}
                    userSp={userSp}
                  />
                ))
              : BUILDING_MASTERS.map((item) => (
                  <BuildingDeployCard
                    key={item.id}
                    count={inventoryCountByBuildingId[item.id] ?? 0}
                    currentGuildSpLanguage={currentGuildSpLanguage}
                    item={item}
                    onDeploy={onDeployBuilding}
                  />
                ))}
          </div>
        </>
      )}
    </motion.aside>
  );
}

interface InventoryTabButtonProps {
  active: boolean;
  badge?: number;
  label: string;
  onClick: () => void;
}

function InventoryTabButton({ active, badge, label, onClick }: InventoryTabButtonProps) {
  return (
    <button
      type="button"
      role="tab"
      aria-selected={active}
      onClick={onClick}
      style={{
        minHeight: "38px",
        border: `2px solid ${active ? "rgba(116, 247, 161, 0.86)" : "rgba(84, 96, 112, 0.72)"}`,
        borderBottomColor: active ? "rgba(24, 83, 45, 0.95)" : "rgba(24, 31, 42, 0.95)",
        borderRightColor: active ? "rgba(24, 83, 45, 0.95)" : "rgba(24, 31, 42, 0.95)",
        background: active ? "rgba(4, 67, 37, 0.9)" : "rgba(10, 16, 28, 0.82)",
        boxShadow: active
          ? "0 0 14px rgba(116,247,161,0.28), inset 0 0 10px rgba(116,247,161,0.12)"
          : "inset 0 0 10px rgba(0,0,0,0.54)",
        color: active ? "#74f7a1" : "rgba(244, 236, 208, 0.58)",
        cursor: "pointer",
        fontFamily: "inherit",
        fontSize: "0.52rem",
        lineHeight: 1,
        padding: "8px 6px",
        textShadow: active
          ? "0 0 8px rgba(116,247,161,0.72), 2px 2px 0 rgba(0,0,0,0.72)"
          : "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      {label}
      {badge !== undefined && badge > 0 ? (
        <span style={{ marginLeft: "6px", color: "#ffd966" }}>x{badge}</span>
      ) : null}
    </button>
  );
}

interface BuildingInventoryCardProps {
  currentGuildLevel: number;
  currentGuildSpLanguage: BuildingTargetSpLanguage;
  item: BuildingMaster;
  onBuy: (building: BuildingMaster) => void;
  userCp: number;
  userSp: number;
}

function BuildingInventoryCard({
  currentGuildLevel,
  currentGuildSpLanguage,
  item,
  onBuy,
  userCp,
  userSp,
}: BuildingInventoryCardProps) {
  const firstLevel = item.levels[0];
  const languageStyle = languageStyles[currentGuildSpLanguage];
  const isLocked = currentGuildLevel < item.requiredGuildLevel;
  const isCpShort = userCp < firstLevel.upgradeCostCp;
  const isSpShort = userSp < firstLevel.upgradeCostSp;
  const canBuild = !isLocked && !isCpShort && !isSpShort;

  return (
    <motion.article
      aria-label={`${item.name}. Requires guild level ${item.requiredGuildLevel}.`}
      aria-disabled={!canBuild}
      whileHover={!isLocked ? { y: -2, backgroundColor: "rgba(255, 217, 102, 0.12)" } : undefined}
      whileTap={canBuild ? { y: 2, scale: 0.98 } : undefined}
      style={{
        position: "relative",
        display: "grid",
        gridTemplateRows: "auto 1fr auto",
        minHeight: "196px",
        gap: "9px",
        overflow: "hidden",
        border: `2px solid ${isLocked ? "rgba(255, 77, 109, 0.48)" : "rgba(116, 247, 161, 0.62)"}`,
        borderBottomColor: isLocked ? "rgba(118, 31, 49, 0.95)" : "rgba(24, 83, 45, 0.95)",
        borderRightColor: isLocked ? "rgba(118, 31, 49, 0.95)" : "rgba(24, 83, 45, 0.95)",
        background: canBuild ? "rgba(1, 12, 24, 0.78)" : "rgba(18, 18, 18, 0.72)",
        boxShadow: "inset 0 0 12px rgba(0,0,0,0.62)",
        color: isLocked ? "rgba(255, 248, 215, 0.48)" : "#fff8d7",
        fontFamily: "inherit",
        fontSize: "0.52rem",
        lineHeight: 1.45,
        padding: "10px",
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        touchAction: "none",
        pointerEvents: isLocked ? "none" : "auto",
      }}
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "64px minmax(0, 1fr)",
          alignItems: "center",
          gap: "9px",
          minWidth: 0,
        }}
      >
        <div
          aria-hidden="true"
          style={{
            display: "grid",
            width: "64px",
            height: "64px",
            placeItems: "center",
            border: `2px solid ${languageStyle.color}`,
            background: item.previewSrc
              ? "rgba(1, 8, 22, 0.78)"
              : "radial-gradient(circle at 50% 40%, rgba(255,255,255,0.16), rgba(0,0,0,0.18) 42%, rgba(0,0,0,0.62))",
            boxShadow: `0 0 14px ${languageStyle.color}55, inset 0 0 12px rgba(0,0,0,0.7)`,
            color: languageStyle.color,
            fontSize: "0.58rem",
            overflow: "hidden",
          }}
        >
          {item.previewSrc ? (
            <img
              className="pixelated"
              src={item.previewSrc}
              alt=""
              draggable={false}
              style={{
                display: "block",
                width: "100%",
                height: "100%",
                objectFit: "contain",
                pointerEvents: "none",
              }}
            />
          ) : (
            languageStyle.label
          )}
        </div>
        <div style={{ minWidth: 0 }}>
          <p
            style={{
              margin: "0 0 5px",
              color: languageStyle.color,
              fontSize: "0.48rem",
              lineHeight: 1.35,
            }}
          >
            LV.{item.requiredGuildLevel} / {item.buffType.toUpperCase()}
          </p>
          <h2
            style={{
              margin: 0,
              color: isLocked ? "rgba(255, 248, 215, 0.5)" : "#ffd966",
              fontSize: "0.56rem",
              lineHeight: 1.45,
              overflowWrap: "anywhere",
            }}
          >
            {item.name}
          </h2>
        </div>
      </div>

      <p
        style={{
          margin: 0,
          color: isLocked ? "rgba(244, 236, 208, 0.48)" : "#f4ecd0",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.8rem",
          lineHeight: 1.35,
        }}
      >
        {item.description}
      </p>

      <div
        style={{
          display: "grid",
          gap: "8px",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "minmax(0, 1fr) minmax(0, 1fr)",
            gap: "6px",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.68rem",
            lineHeight: 1.2,
          }}
        >
          <CostPill
            isShort={isCpShort}
            label={`${firstLevel.upgradeCostCp.toLocaleString()} CP`}
            tone="#ffd966"
          />
          <CostPill
            isShort={isSpShort}
            label={`${firstLevel.upgradeCostSp.toLocaleString()} ${currentGuildSpLanguage}-SP`}
            tone={languageStyle.color}
          />
        </div>

        <button
          type="button"
          disabled={!canBuild}
          onClick={() => onBuy(item)}
          style={{
            width: "100%",
            border: `2px solid ${canBuild ? "rgba(116, 247, 161, 0.74)" : "rgba(92, 92, 92, 0.78)"}`,
            borderBottomColor: canBuild ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
            borderRightColor: canBuild ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
            background: canBuild ? "rgba(4, 67, 37, 0.9)" : "rgba(37, 37, 37, 0.92)",
            boxShadow: canBuild
              ? "0 0 12px rgba(116,247,161,0.2), inset 0 0 10px rgba(116,247,161,0.1)"
              : "inset 0 0 10px rgba(0,0,0,0.52)",
            color: canBuild ? "#74f7a1" : "rgba(255, 248, 215, 0.38)",
            cursor: canBuild ? "pointer" : "not-allowed",
            fontFamily: "inherit",
            fontSize: "0.56rem",
            lineHeight: 1,
            padding: "9px 8px",
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          BUY
        </button>
      </div>

      {isLocked && (
        <div
          aria-hidden="true"
          className="absolute inset-0 z-10 flex flex-col items-center justify-center bg-black/80"
          style={{
            gap: "8px",
            padding: "14px",
            pointerEvents: "none",
            textAlign: "center",
          }}
        >
          <span
            style={{
              color: "#ff4d6d",
              fontSize: "0.72rem",
              lineHeight: 1.6,
              textShadow: "0 0 10px rgba(255,77,109,0.82), 2px 2px 0 rgba(0,0,0,0.8)",
            }}
          >
            REQUIRES
            <br />
            GUILD LV.{item.requiredGuildLevel}
          </span>
        </div>
      )}
    </motion.article>
  );
}

interface BuildingDeployCardProps {
  count: number;
  currentGuildSpLanguage: BuildingTargetSpLanguage;
  item: BuildingMaster;
  onDeploy: (building: BuildingMaster) => void;
}

function BuildingDeployCard({
  count,
  currentGuildSpLanguage,
  item,
  onDeploy,
}: BuildingDeployCardProps) {
  const languageStyle = languageStyles[currentGuildSpLanguage];
  const canDeploy = count > 0;

  return (
    <motion.article
      aria-label={`${item.name}. Inventory count ${count}.`}
      aria-disabled={!canDeploy}
      whileHover={canDeploy ? { y: -2, backgroundColor: "rgba(116, 247, 161, 0.1)" } : undefined}
      whileTap={canDeploy ? { y: 2, scale: 0.98 } : undefined}
      style={{
        display: "grid",
        gridTemplateRows: "auto 1fr auto",
        minHeight: "178px",
        gap: "9px",
        opacity: canDeploy ? 1 : 0.46,
        overflow: "hidden",
        border: `2px solid ${canDeploy ? "rgba(116, 247, 161, 0.62)" : "rgba(84, 96, 112, 0.62)"}`,
        borderBottomColor: canDeploy ? "rgba(24, 83, 45, 0.95)" : "rgba(24, 31, 42, 0.95)",
        borderRightColor: canDeploy ? "rgba(24, 83, 45, 0.95)" : "rgba(24, 31, 42, 0.95)",
        background: canDeploy ? "rgba(1, 12, 24, 0.78)" : "rgba(18, 18, 18, 0.72)",
        boxShadow: "inset 0 0 12px rgba(0,0,0,0.62)",
        color: "#fff8d7",
        fontFamily: "inherit",
        fontSize: "0.52rem",
        lineHeight: 1.45,
        padding: "10px",
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "64px minmax(0, 1fr) auto",
          alignItems: "center",
          gap: "9px",
          minWidth: 0,
        }}
      >
        <div
          aria-hidden="true"
          style={{
            display: "grid",
            width: "64px",
            height: "64px",
            placeItems: "center",
            border: `2px solid ${languageStyle.color}`,
            background: "rgba(1, 8, 22, 0.78)",
            boxShadow: `0 0 14px ${languageStyle.color}55, inset 0 0 12px rgba(0,0,0,0.7)`,
            color: languageStyle.color,
            fontSize: "0.58rem",
            overflow: "hidden",
          }}
        >
          {item.previewSrc ? (
            <img
              className="pixelated"
              src={item.previewSrc}
              alt=""
              draggable={false}
              style={{
                display: "block",
                width: "100%",
                height: "100%",
                objectFit: "contain",
                pointerEvents: "none",
              }}
            />
          ) : (
            languageStyle.label
          )}
        </div>
        <div style={{ minWidth: 0 }}>
          <p
            style={{
              margin: "0 0 5px",
              color: languageStyle.color,
              fontSize: "0.48rem",
              lineHeight: 1.35,
            }}
          >
            {item.buffType.toUpperCase()}
          </p>
          <h2
            style={{
              margin: 0,
              color: "#ffd966",
              fontSize: "0.56rem",
              lineHeight: 1.45,
              overflowWrap: "anywhere",
            }}
          >
            {item.name}
          </h2>
        </div>
        <span
          aria-label={`Owned count ${count}`}
          style={{
            minWidth: "42px",
            border: "2px solid rgba(255, 217, 102, 0.78)",
            background: "rgba(1, 8, 22, 0.76)",
            boxShadow: "0 0 10px rgba(255,217,102,0.22), inset 0 0 10px rgba(0,0,0,0.58)",
            color: "#ffd966",
            padding: "6px 5px",
            textAlign: "center",
          }}
        >
          x{count}
        </span>
      </div>

      <p
        style={{
          margin: 0,
          color: "#f4ecd0",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.8rem",
          lineHeight: 1.35,
        }}
      >
        {item.description}
      </p>

      <button
        type="button"
        disabled={!canDeploy}
        onClick={() => onDeploy(item)}
        style={{
          width: "100%",
          border: `2px solid ${canDeploy ? "rgba(116, 247, 161, 0.74)" : "rgba(92, 92, 92, 0.78)"}`,
          borderBottomColor: canDeploy ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
          borderRightColor: canDeploy ? "rgba(24, 83, 45, 0.95)" : "rgba(32, 32, 32, 0.95)",
          background: canDeploy ? "rgba(4, 67, 37, 0.9)" : "rgba(37, 37, 37, 0.92)",
          boxShadow: canDeploy
            ? "0 0 12px rgba(116,247,161,0.2), inset 0 0 10px rgba(116,247,161,0.1)"
            : "inset 0 0 10px rgba(0,0,0,0.52)",
          color: canDeploy ? "#74f7a1" : "rgba(255, 248, 215, 0.38)",
          cursor: canDeploy ? "pointer" : "not-allowed",
          fontFamily: "inherit",
          fontSize: "0.56rem",
          lineHeight: 1,
          padding: "9px 8px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        SET
      </button>
    </motion.article>
  );
}

interface CostPillProps {
  isShort: boolean;
  label: string;
  tone: string;
}

function CostPill({ isShort, label, tone }: CostPillProps) {
  return (
    <span
      className={isShort ? "text-red-500" : undefined}
      style={{
        border: `2px solid ${isShort ? "#ef4444" : tone}`,
        background: "rgba(1, 8, 22, 0.76)",
        boxShadow: `inset 0 0 10px rgba(0,0,0,0.58), 0 0 10px ${isShort ? "#ef444455" : `${tone}44`}`,
        color: isShort ? "#ef4444" : tone,
        display: "grid",
        lineHeight: 1.25,
        minWidth: 0,
        overflow: "hidden",
        padding: "6px 5px",
        placeItems: "center",
        textAlign: "center",
        overflowWrap: "anywhere",
      }}
      title={label}
    >
      {label}
    </span>
  );
}
