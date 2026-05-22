import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";

export function GuildTownLoadingOverlay() {
  return (
    <div
      role="status"
      aria-live="polite"
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 11,
        display: "grid",
        placeItems: "center",
        background: "rgba(3, 7, 14, 0.58)",
        color: "#74f7a1",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        fontSize: "0.72rem",
        letterSpacing: 0,
        textShadow: "0 0 12px rgba(116,247,161,0.78), 2px 2px 0 rgba(0,0,0,0.82)",
      }}
    >
      <span
        style={{
          border: "2px solid rgba(116, 247, 161, 0.78)",
          borderBottomColor: "rgba(24, 83, 45, 0.95)",
          borderRightColor: "rgba(24, 83, 45, 0.95)",
          background: "rgba(1, 12, 24, 0.9)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
          padding: "14px 18px",
        }}
      >
        SYNCING GUILD TOWN...
      </span>
    </div>
  );
}

export function DeployModePanel({
  buildingName,
  onCancel,
}: {
  buildingName: string;
  onCancel: () => void;
}) {
  return (
    <motion.div
      role="status"
      aria-live="polite"
      initial={{ opacity: 0, y: -10, scale: 0.96 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: -8, scale: 0.96 }}
      transition={{ duration: 0.18, ease: steppedEase(5) }}
      style={{
        position: "fixed",
        left: "50%",
        top: "calc(env(safe-area-inset-top, 0px) + 92px)",
        zIndex: 13,
        display: "grid",
        gridTemplateColumns: "minmax(0, 1fr) auto",
        alignItems: "center",
        gap: "12px",
        maxWidth: "min(720px, calc(100vw - 32px))",
        transform: "translateX(-50%)",
        border: "3px solid rgba(116, 247, 161, 0.86)",
        borderBottomColor: "rgba(24, 83, 45, 0.95)",
        borderRightColor: "rgba(24, 83, 45, 0.95)",
        background: "rgba(1, 12, 24, 0.94)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34), 0 0 22px rgba(116,247,161,0.24)",
        color: "#fff8d7",
        fontFamily: '"DotGothic16", monospace',
        fontSize: "0.9rem",
        lineHeight: 1.35,
        padding: "10px 12px",
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      <span style={{ minWidth: 0, overflowWrap: "anywhere" }}>
        {buildingName} をドラッグまたは矢印キーで動かし、BUILD / Enter で配置してください
      </span>
      <button
        type="button"
        onClick={onCancel}
        style={{
          minHeight: "34px",
          border: "2px solid rgba(255, 217, 102, 0.86)",
          borderBottomColor: "rgba(96, 62, 22, 0.98)",
          borderRightColor: "rgba(96, 62, 22, 0.98)",
          background: "rgba(3, 10, 24, 0.9)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
          color: "#ffd966",
          cursor: "pointer",
          fontFamily: '"Press Start 2P", "DotGothic16", monospace',
          fontSize: "0.5rem",
          lineHeight: 1,
          padding: "8px 10px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        CANCEL
      </button>
    </motion.div>
  );
}

export function GuildTownToast({ message }: { message: string }) {
  return (
    <p
      role="alert"
      style={{
        position: "fixed",
        bottom: "calc(env(safe-area-inset-bottom, 0px) + 22px)",
        left: "50%",
        zIndex: 12,
        margin: 0,
        maxWidth: "min(720px, calc(100vw - 32px))",
        transform: "translateX(-50%)",
        border: "2px solid rgba(255, 77, 109, 0.86)",
        borderBottomColor: "rgba(118, 31, 49, 0.95)",
        borderRightColor: "rgba(118, 31, 49, 0.95)",
        background: "rgba(18, 8, 14, 0.94)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
        color: "#ff9aae",
        fontFamily: '"DotGothic16", monospace',
        fontSize: "0.92rem",
        lineHeight: 1.45,
        padding: "10px 14px",
        textAlign: "center",
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      {message}
    </p>
  );
}
