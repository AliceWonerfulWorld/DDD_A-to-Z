import { motion } from "framer-motion";
import { steppedEase } from "../lib/animationUtils";

interface GuildTownProps {
  onNavigate: (path: string) => void;
  townLevel?: number;
  currentCp?: number;
  nextLevelCp?: number;
  baseSrc?: string;
  mainStructureSrc?: string;
  bonfireSrc?: string;
}

export function GuildTown({
  onNavigate,
  townLevel = 1,
  currentCp = 2500,
  nextLevelCp = 10000,
  baseSrc = "/town/glassfield.png",
  mainStructureSrc = "/town/tent.png",
  bonfireSrc = "/town/bonfire.png",
}: GuildTownProps) {
  const progress = Math.min(100, Math.max(0, (currentCp / nextLevelCp) * 100));

  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        background: "#112b1a",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#fff8d7",
      }}
    >
      <img
        className="pixelated"
        src={baseSrc}
        alt=""
        aria-hidden="true"
        style={{
          position: "absolute",
          inset: 0,
          width: "100%",
          height: "100%",
          objectFit: "cover",
          objectPosition: "center bottom",
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

      <div
        style={{
          position: "absolute",
          left: "50%",
          top: "50%",
          width: "clamp(210px, 29vw, 430px)",
          transform: "translate(-64%, -28%)",
          zIndex: 2,
          pointerEvents: "none",
        }}
      >
        <motion.img
          className="pixelated drop-shadow-[18px_22px_0_rgba(0,0,0,0.28)]"
          src={mainStructureSrc}
          alt=""
          aria-hidden="true"
          initial={{ opacity: 0, y: 18, scale: 0.96 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.36, ease: steppedEase(7) }}
          style={{
            display: "block",
            width: "100%",
            height: "auto",
            mixBlendMode: "screen",
          }}
        />
      </div>

      <div
        style={{
          position: "absolute",
          left: "50%",
          top: "50%",
          width: "clamp(92px, 12vw, 164px)",
          transform: "translate(42%, 62%)",
          zIndex: 3,
          pointerEvents: "none",
        }}
      >
        <div
          aria-hidden="true"
          style={{
            position: "absolute",
            left: "50%",
            top: "46%",
            width: "68%",
            aspectRatio: "1 / 1",
            transform: "translate(-50%, -50%)",
          }}
        >
          <motion.div
            className="bg-orange-500/30 blur-2xl"
            animate={{
              opacity: [0.58, 0.88, 0.66, 0.8, 0.58],
              scale: [0.92, 1.08, 0.98, 1.04, 0.92],
            }}
            transition={{ duration: 1.45, ease: "easeInOut", repeat: Infinity }}
            style={{
              width: "100%",
              height: "100%",
              clipPath: "circle(50% at 50% 50%)",
            }}
          />
        </div>
        <motion.img
          className="pixelated drop-shadow-[10px_14px_0_rgba(0,0,0,0.26)]"
          src={bonfireSrc}
          alt=""
          aria-hidden="true"
          initial={{ opacity: 0, scale: 0.86 }}
          animate={{
            opacity: 1,
            scale: [1, 1.035, 0.98, 1.02, 1],
            y: [0, -2, 1, -1, 0],
          }}
          transition={{
            opacity: { duration: 0.26, ease: steppedEase(5) },
            scale: { duration: 1.25, ease: "easeInOut", repeat: Infinity },
            y: { duration: 1.25, ease: "easeInOut", repeat: Infinity },
          }}
          style={{
            display: "block",
            position: "relative",
            width: "100%",
            height: "auto",
            mixBlendMode: "screen",
          }}
        />
      </div>

      <motion.header
        initial={{ opacity: 0, y: -14 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.32, ease: steppedEase(6) }}
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(68px, 7vw, 86px))",
          left: "50%",
          zIndex: 5,
          width: "min(calc(100vw - 32px), 560px)",
          transform: "translateX(-50%)",
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
            TOWN LEVEL {townLevel}
          </strong>
          <div style={{ minWidth: 0 }}>
            <div
              aria-label={`${currentCp.toLocaleString()} / ${nextLevelCp.toLocaleString()} CP`}
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
              {currentCp.toLocaleString()} / {nextLevelCp.toLocaleString()} CP
            </p>
          </div>
        </div>
      </motion.header>

      <motion.button
        type="button"
        whileHover={{ y: -2, scale: 1.02 }}
        whileTap={{ y: 2, scale: 0.98 }}
        onClick={() => onNavigate("/guild")}
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(14px, 2.2vw, 28px))",
          left: "clamp(14px, 2.2vw, 28px)",
          zIndex: 6,
          minHeight: "42px",
          border: "2px solid rgba(255, 217, 102, 0.78)",
          borderBottomColor: "rgba(96, 62, 22, 0.95)",
          borderRightColor: "rgba(96, 62, 22, 0.95)",
          background: "rgba(3, 10, 24, 0.74)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.56rem, 1.3vw, 0.78rem)",
          lineHeight: 1.5,
          padding: "10px 12px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        &lt; BACK
      </motion.button>

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.08), rgba(0,0,0,0.08) 1px, transparent 1px, transparent 4px)",
          pointerEvents: "none",
          zIndex: 4,
        }}
      />
    </main>
  );
}
