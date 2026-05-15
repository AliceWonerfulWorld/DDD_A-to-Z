import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import type { PlacedItem } from "./types";

interface BuildingInfoPanelProps {
  item: PlacedItem | null;
  onClose: () => void;
}

export function BuildingInfoPanel({ item, onClose }: BuildingInfoPanelProps) {
  if (!item) return null;

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
        left: "50%",
        bottom: "calc(env(safe-area-inset-bottom, 0px) + clamp(14px, 2vw, 24px))",
        zIndex: 9,
        display: "grid",
        gridTemplateColumns: "76px minmax(0, 1fr)",
        width: "min(calc(100vw - 210px), 720px)",
        minHeight: "96px",
        alignItems: "center",
        gap: "14px",
        transform: "translateX(-50%)",
        border: "3px solid rgba(255, 248, 215, 0.82)",
        borderBottomColor: "rgba(55, 44, 35, 0.98)",
        borderRightColor: "rgba(55, 44, 35, 0.98)",
        background: "linear-gradient(180deg, rgba(4, 10, 22, 0.94), rgba(3, 7, 14, 0.9))",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.76), 7px 7px 0 rgba(0,0,0,0.36), inset 0 0 22px rgba(116,247,161,0.09)",
        color: "#fff8d7",
        padding: "12px 14px",
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
          width: "76px",
          height: "68px",
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
            maxWidth: "58px",
            maxHeight: "54px",
            filter: "drop-shadow(4px 5px 0 rgba(0,0,0,0.34))",
          }}
        />
      </div>
      <div style={{ minWidth: 0 }}>
        <p
          style={{
            margin: "0 0 6px",
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
            margin: "0 0 7px",
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
            fontSize: "clamp(0.78rem, 1.45vw, 0.98rem)",
            lineHeight: 1.45,
          }}
        >
          {item.description}
        </p>
      </div>
    </motion.section>
  );
}
