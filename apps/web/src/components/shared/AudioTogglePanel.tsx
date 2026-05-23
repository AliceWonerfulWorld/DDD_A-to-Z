import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import type { CSSProperties } from "react";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { PixelSpeakerIcon } from "./PixelSpeakerIcon";

type AudioTogglePanelPosition = "top-left" | "top-right" | "bottom-left" | "bottom-right";

interface AudioTogglePanelProps {
  position?: AudioTogglePanelPosition;
  inlineOnMobile?: boolean;
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

const PC_BREAKPOINT = 768;

export function AudioTogglePanel({
  position = "top-left",
  inlineOnMobile = false,
}: AudioTogglePanelProps) {
  const { isBgmEnabled, isSeEnabled, toggleBgm, toggleSe } = useAudioSettings();
  const [isPcSize, setIsPcSize] = useState(
    typeof window !== "undefined" ? window.innerWidth >= PC_BREAKPOINT : true,
  );
  const [isExpanded, setIsExpanded] = useState(isPcSize);

  useEffect(() => {
    const handleResize = () => {
      const isNowPc = window.innerWidth >= PC_BREAKPOINT;
      setIsPcSize(isNowPc);
      setIsExpanded(isNowPc);
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const toggleExpanded = () => {
    if (!isPcSize) {
      setIsExpanded(!isExpanded);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: -12 }}
      animate={{ opacity: 1, y: 0 }}
      style={
        !isPcSize && inlineOnMobile
          ? {
              display: "flex",
              flexDirection: "column", // Expand downwards
              gap: "8px",
              alignItems: "flex-start", // Align items correctly
              position: "relative",
              zIndex: 4,
            }
          : {
              position: "fixed",
              ...positionStyles[position],
              ...(!isPcSize && position.includes("bottom") ? { bottom: "150px" } : {}),
              zIndex: 4,
              display: "flex",
              flexDirection: "column",
              gap: "8px",
              alignItems: position.includes("right") ? "flex-end" : "flex-start",
            }
      }
    >
      {/* モバイルのみ：メインオーディオボタン */}
      {!isPcSize && (
        <motion.button
          type="button"
          onClick={toggleExpanded}
          aria-label="オーディオ設定"
          aria-expanded={isExpanded}
          whileHover={{ scale: 1.08 }}
          whileTap={{ scale: 0.96 }}
          style={{
            display: "inline-flex",
            alignItems: "center",
            justifyContent: "center",
            width: "44px",
            height: "44px",
            minHeight: "44px",
            minWidth: "44px",
            border: "2px solid #00f5ff",
            boxShadow: "3px 3px 0 rgba(0,0,0,0.72), 0 0 14px rgba(0,245,255,0.3)",
            background: "rgba(8, 12, 18, 0.82)",
            backdropFilter: "blur(2px)",
            color: "#00f5ff",
            cursor: "pointer",
            fontFamily: '"DotGothic16", monospace',
            textShadow: "1px 1px 0 #000",
            fontSize: "1.2rem",
            lineHeight: 1,
            transition: "all 0.2s ease",
          }}
        >
          <span aria-hidden="true">🔊</span>
        </motion.button>
      )}

      {/* PC時は常に表示、モバイルは展開時のみ表示 */}
      <AnimatePresence>
        {isExpanded && (
          <motion.div
            initial={!isPcSize ? { opacity: 0, scale: 0.95 } : undefined}
            animate={!isPcSize ? { opacity: 1, scale: 1 } : undefined}
            exit={!isPcSize ? { opacity: 0, scale: 0.95 } : undefined}
            transition={{ duration: 0.2 }}
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "6px",
            }}
          >
            {/* BGM ボタン */}
            <motion.button
              type="button"
              onClick={() => {
                toggleBgm();
              }}
              aria-label={isBgmEnabled ? "BGMをオフにする" : "BGMをオンにする"}
              aria-pressed={isBgmEnabled}
              whileHover={{ scale: 1.04 }}
              whileTap={{ y: 1, scale: 0.98 }}
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "6px",
                padding: "8px 10px",
                minHeight: "40px",
                border: `2px solid ${isBgmEnabled ? "#ffd700" : "#ffffff66"}`,
                boxShadow: `3px 3px 0 rgba(0,0,0,0.72), 0 0 14px ${
                  isBgmEnabled ? "rgba(255,215,0,0.3)" : "rgba(255,255,255,0.12)"
                }`,
                background: "rgba(8, 12, 18, 0.82)",
                backdropFilter: "blur(2px)",
                color: isBgmEnabled ? "#fff7dc" : "#ffffff99",
                cursor: "pointer",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.65rem, 2vw, 0.72rem)",
                letterSpacing: "0.08em",
                textShadow: "1px 1px 0 #000",
                whiteSpace: "nowrap",
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
              <span>{isBgmEnabled ? "BGM" : "OFF"}</span>
            </motion.button>

            {/* SE ボタン */}
            <motion.button
              type="button"
              onClick={() => {
                toggleSe();
              }}
              aria-label={isSeEnabled ? "SEをオフにする" : "SEをオンにする"}
              aria-pressed={isSeEnabled}
              whileHover={{ scale: 1.04 }}
              whileTap={{ y: 1, scale: 0.98 }}
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "6px",
                padding: "8px 10px",
                minHeight: "40px",
                border: `2px solid ${isSeEnabled ? "#00f5ff" : "#ffffff66"}`,
                boxShadow: `3px 3px 0 rgba(0,0,0,0.72), 0 0 14px ${
                  isSeEnabled ? "rgba(0,245,255,0.28)" : "rgba(255,255,255,0.12)"
                }`,
                background: "rgba(8, 12, 18, 0.82)",
                backdropFilter: "blur(2px)",
                color: isSeEnabled ? "#e8ffff" : "#ffffff99",
                cursor: "pointer",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.65rem, 2vw, 0.72rem)",
                letterSpacing: "0.08em",
                textShadow: "1px 1px 0 #000",
                whiteSpace: "nowrap",
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
              <span>{isSeEnabled ? "SE" : "OFF"}</span>
            </motion.button>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
}
