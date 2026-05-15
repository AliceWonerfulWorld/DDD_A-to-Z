import { motion, type PanInfo } from "framer-motion";
import type { PointerEvent as ReactPointerEvent, RefObject } from "react";
import { steppedEase } from "../../lib/animationUtils";
import type { InventoryItem } from "./types";

interface BuildInventoryProps {
  inventory: InventoryItem[];
  inventoryRef: RefObject<HTMLDivElement | null>;
  onDragEnd: (
    item: InventoryItem,
    event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => void;
  onToggleVisible: () => void;
  stopNestedDrag: (event: ReactPointerEvent<HTMLElement>) => void;
  visible: boolean;
}

export function BuildInventory({
  inventory,
  inventoryRef,
  onDragEnd,
  onToggleVisible,
  stopNestedDrag,
  visible,
}: BuildInventoryProps) {
  return (
    <motion.aside
      ref={inventoryRef}
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
        width: "148px",
        maxHeight:
          "calc(100vh - env(safe-area-inset-top, 0px) - env(safe-area-inset-bottom, 0px) - 118px)",
        alignItems: "stretch",
        flexDirection: "column",
        gap: "10px",
        overflow: "visible",
        border: visible ? "3px solid rgba(255, 248, 215, 0.8)" : "3px solid transparent",
        borderBottomColor: visible ? "rgba(55, 44, 35, 0.98)" : "transparent",
        borderRightColor: visible ? "rgba(55, 44, 35, 0.98)" : "transparent",
        background: visible ? "rgba(3, 7, 14, 0.88)" : "transparent",
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

      {visible &&
        inventory.map((item) => {
          const isAvailable = item.count > 0;

          return (
            <motion.div
              key={item.type}
              role="button"
              aria-label={`${item.name} inventory item. ${item.count} remaining.`}
              aria-disabled={!isAvailable}
              tabIndex={isAvailable ? 0 : -1}
              drag={isAvailable}
              dragSnapToOrigin
              dragElastic={0}
              dragMomentum={false}
              onPointerDown={stopNestedDrag}
              onDragEnd={(event, info) => onDragEnd(item, event, info)}
              whileHover={
                isAvailable ? { y: -2, backgroundColor: "rgba(255, 217, 102, 0.12)" } : undefined
              }
              whileTap={isAvailable ? { y: 2, scale: 0.98 } : undefined}
              whileDrag={{ scale: 1.06, zIndex: 20 }}
              style={{
                position: "relative",
                display: "grid",
                gridTemplateRows: "58px auto",
                width: "100%",
                minHeight: "116px",
                alignItems: "center",
                justifyItems: "center",
                gap: "6px",
                border: "2px solid rgba(116, 247, 161, 0.62)",
                borderBottomColor: "rgba(24, 83, 45, 0.95)",
                borderRightColor: "rgba(24, 83, 45, 0.95)",
                background: isAvailable ? "rgba(1, 12, 24, 0.78)" : "rgba(18, 18, 18, 0.68)",
                boxShadow: "inset 0 0 12px rgba(0,0,0,0.62)",
                color: isAvailable ? "#fff8d7" : "rgba(255, 248, 215, 0.42)",
                cursor: isAvailable ? "grab" : "not-allowed",
                fontFamily: "inherit",
                fontSize: "0.52rem",
                lineHeight: 1.35,
                padding: "8px 7px",
                textAlign: "center",
                textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
                touchAction: "none",
              }}
            >
              <img
                className="pixelated"
                src={item.src}
                alt=""
                aria-hidden="true"
                draggable={false}
                style={{
                  display: "block",
                  maxWidth: "74px",
                  maxHeight: "58px",
                  opacity: isAvailable ? 1 : 0.38,
                  pointerEvents: "none",
                  filter: "drop-shadow(4px 5px 0 rgba(0,0,0,0.34))",
                  mixBlendMode: "screen",
                }}
              />
              <span>{item.name}</span>
              <span
                aria-hidden="true"
                style={{
                  position: "absolute",
                  right: "6px",
                  top: "6px",
                  minWidth: "24px",
                  border: "2px solid rgba(255, 217, 102, 0.78)",
                  background: isAvailable ? "rgba(62, 26, 8, 0.92)" : "rgba(8, 8, 8, 0.92)",
                  color: isAvailable ? "#ffd966" : "rgba(255, 248, 215, 0.46)",
                  fontSize: "0.52rem",
                  lineHeight: 1,
                  padding: "4px 5px",
                  boxShadow: "0 0 0 1px rgba(0,0,0,0.72)",
                }}
              >
                x{item.count}
              </span>
            </motion.div>
          );
        })}
    </motion.aside>
  );
}
