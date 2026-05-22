import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";

interface TownStatusHeaderProps {
  currentCp: number;
  nextLevelCp: number;
  progress: number;
  guildLevel: number;
}

export function TownStatusHeader({
  currentCp,
  nextLevelCp,
  progress,
  guildLevel,
}: TownStatusHeaderProps) {
  return (
    <motion.header
      initial={{ opacity: 0, y: -14 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.32, ease: steppedEase(6) }}
      style={{
        position: "fixed",
        top: "calc(env(safe-area-inset-top, 0px) + clamp(14px, 2.2vw, 28px))",
        right: "clamp(14px, 2.2vw, 28px)",
        zIndex: 5,
        width: "min(calc(100vw - 150px), 720px)",
        border: "3px solid rgba(255, 248, 215, 0.76)",
        borderBottomColor: "rgba(55, 44, 35, 0.96)",
        borderRightColor: "rgba(55, 44, 35, 0.96)",
        background: "rgba(3, 10, 24, 0.72)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.68), 6px 6px 0 rgba(0,0,0,0.34), inset 0 0 18px rgba(116,247,161,0.1)",
        padding: "clamp(10px, 2vw, 14px)",
        backdropFilter: "blur(2px)",
      }}
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "auto 1fr",
          alignItems: "center",
          gap: "clamp(10px, 2vw, 16px)",
        }}
      >
        <strong
          style={{
            color: "#ffd966",
            fontSize: "clamp(0.62rem, 1.45vw, 0.86rem)",
            lineHeight: 1.5,
            whiteSpace: "nowrap",
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          GUILD LEVEL {guildLevel}
        </strong>
        <div style={{ minWidth: 0 }}>
          <div
            aria-label={`${currentCp.toLocaleString()} / ${nextLevelCp.toLocaleString()} EXP`}
            role="progressbar"
            aria-valuemin={0}
            aria-valuemax={nextLevelCp}
            aria-valuenow={currentCp}
            style={{
              height: "14px",
              border: "2px solid rgba(116, 247, 161, 0.72)",
              background: "rgba(1, 8, 22, 0.8)",
              boxShadow: "inset 0 0 10px rgba(0,0,0,0.64)",
              overflow: "hidden",
            }}
          >
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.52, ease: steppedEase(8) }}
              style={{
                height: "100%",
                background:
                  "repeating-linear-gradient(90deg, #74f7a1 0, #74f7a1 10px, #39ff14 10px, #39ff14 20px)",
                boxShadow: "0 0 12px rgba(116,247,161,0.72)",
              }}
            />
          </div>
          <p
            style={{
              margin: "7px 0 0",
              color: "#f4ecd0",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "clamp(0.72rem, 1.45vw, 0.95rem)",
              lineHeight: 1.4,
              textAlign: "right",
            }}
          >
            {currentCp.toLocaleString()} / {nextLevelCp.toLocaleString()} EXP
          </p>
        </div>
      </div>
    </motion.header>
  );
}
