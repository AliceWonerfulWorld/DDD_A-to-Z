import { motion } from "framer-motion";
import type { CSSProperties } from "react";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { PixelSpeakerIcon } from "./PixelSpeakerIcon";

type AudioTogglePanelPosition = "top-left" | "top-right" | "bottom-left" | "bottom-right";

interface AudioTogglePanelProps {
  position?: AudioTogglePanelPosition;
}

const positionStyles: Record<AudioTogglePanelPosition, CSSProperties> = {
  "top-left": {
    top: "clamp(14px, 3vw, 28px)",
    left: "clamp(14px, 3vw, 28px)",
  },
  "top-right": {
    top: "clamp(14px, 3vw, 28px)",
    right: "clamp(14px, 3vw, 28px)",
  },
  "bottom-left": {
    bottom: "clamp(14px, 3vw, 28px)",
    left: "clamp(14px, 3vw, 28px)",
  },
  "bottom-right": {
    right: "clamp(14px, 3vw, 28px)",
    bottom: "clamp(14px, 3vw, 28px)",
  },
};

export function AudioTogglePanel({ position = "top-left" }: AudioTogglePanelProps) {
  const { isBgmEnabled, isSeEnabled, toggleBgm, toggleSe } = useAudioSettings();

  return (
    <motion.div
      initial={{ opacity: 0, y: -12 }}
      animate={{ opacity: 1, y: 0 }}
      style={{
        position: "fixed",
        ...positionStyles[position],
        zIndex: 4,
        display: "grid",
        gap: "8px",
      }}
    >
      <motion.button
        type="button"
        onClick={toggleBgm}
        aria-label={isBgmEnabled ? "BGMをオフにする" : "BGMをオンにする"}
        aria-pressed={isBgmEnabled}
        whileHover={{ scale: 1.04 }}
        whileTap={{ y: 2, scale: 0.98 }}
        style={{
          display: "inline-flex",
          alignItems: "center",
          gap: "8px",
          padding: "9px 12px",
          border: `2px solid ${isBgmEnabled ? "#ffd700" : "#ffffff66"}`,
          boxShadow: `3px 3px 0 rgba(0,0,0,0.72), 0 0 14px ${
            isBgmEnabled ? "rgba(255,215,0,0.3)" : "rgba(255,255,255,0.12)"
          }`,
          background: "rgba(8, 12, 18, 0.72)",
          backdropFilter: "blur(2px)",
          color: isBgmEnabled ? "#fff7dc" : "#ffffff99",
          cursor: "pointer",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.72rem",
          letterSpacing: "0.08em",
          textShadow: "1px 1px 0 #000",
        }}
      >
        <span
          aria-hidden="true"
          style={{
            position: "relative",
            display: "inline-grid",
            placeItems: "center",
            width: "1em",
            height: "1em",
            fontSize: "0.95rem",
            lineHeight: 1,
          }}
        >
          ♪
          {!isBgmEnabled && (
            <span
              aria-hidden="true"
              style={{
                position: "absolute",
                width: "1.25em",
                height: "2px",
                background: "#ffb0aa",
                boxShadow: "1px 1px 0 #000",
                transform: "rotate(-45deg)",
              }}
            />
          )}
        </span>
        <span>{isBgmEnabled ? "BGM ON" : "BGM OFF"}</span>
      </motion.button>

      <motion.button
        type="button"
        onClick={toggleSe}
        aria-label={isSeEnabled ? "SEをオフにする" : "SEをオンにする"}
        aria-pressed={isSeEnabled}
        whileHover={{ scale: 1.04 }}
        whileTap={{ y: 2, scale: 0.98 }}
        style={{
          display: "inline-flex",
          alignItems: "center",
          gap: "8px",
          padding: "9px 12px",
          border: `2px solid ${isSeEnabled ? "#00f5ff" : "#ffffff66"}`,
          boxShadow: `3px 3px 0 rgba(0,0,0,0.72), 0 0 14px ${
            isSeEnabled ? "rgba(0,245,255,0.28)" : "rgba(255,255,255,0.12)"
          }`,
          background: "rgba(8, 12, 18, 0.72)",
          backdropFilter: "blur(2px)",
          color: isSeEnabled ? "#e8ffff" : "#ffffff99",
          cursor: "pointer",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.72rem",
          letterSpacing: "0.08em",
          textShadow: "1px 1px 0 #000",
        }}
      >
        <span
          aria-hidden="true"
          style={{
            position: "relative",
            display: "inline-grid",
            placeItems: "center",
            width: "1.55em",
            height: "1em",
            fontSize: "0.86rem",
            lineHeight: 1,
          }}
        >
          <PixelSpeakerIcon muted={!isSeEnabled} />
        </span>
        <span>{isSeEnabled ? "SE ON" : "SE OFF"}</span>
      </motion.button>
    </motion.div>
  );
}
