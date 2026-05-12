import { motion, type Variants } from "framer-motion";
import type { ReactNode } from "react";

interface HomeProps {
  onNavigate: (path: string) => void;
}

const player = {
  name: "DevSamurai",
  title: "Consistency Master",
  level: 18,
  totalCp: 24680,
  todayCp: 320,
};

const navItems = [
  { label: "WAR MAP", caption: "BATTLE FRONT", path: "/?war=1", accent: "#ff5f56" },
  { label: "GUILD BASE", caption: "COMMUNITY HQ", path: "/?guild=1", accent: "#00f5ff" },
  { label: "MY STATUS", caption: "PLAYER DATA", path: "/?mypage=1", accent: "#ffd700" },
];

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

const panelVariants: Variants = {
  hidden: { opacity: 0, y: -14 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.35, ease: steppedEase(6) },
  },
};

const navVariants: Variants = {
  hidden: { opacity: 0, y: 18 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: 0.15 + i * 0.08, duration: 0.35, ease: steppedEase(6) },
  }),
};

function HudPanel({ align = "left", children }: { align?: "left" | "right"; children: ReactNode }) {
  return (
    <motion.section
      variants={panelVariants}
      initial="hidden"
      animate="visible"
      style={{
        width: "min(100%, 360px)",
        border: "3px solid rgba(255, 215, 0, 0.72)",
        borderBottomColor: "rgba(111, 79, 28, 0.95)",
        borderRightColor: "rgba(111, 79, 28, 0.95)",
        background: "rgba(3, 10, 24, 0.72)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.78), 8px 8px 0 rgba(0,0,0,0.45)",
        color: "#f4ecd0",
        padding: "14px 16px",
        textAlign: align,
        backdropFilter: "blur(2px)",
      }}
    >
      {children}
    </motion.section>
  );
}

function LabelValue({ label, value }: { label: string; value: string | number }) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "minmax(78px, max-content) 1fr",
        gap: "12px",
        alignItems: "baseline",
        minHeight: "28px",
      }}
    >
      <span
        style={{
          color: "rgba(244, 236, 208, 0.62)",
          fontSize: "0.62rem",
          lineHeight: 1.5,
        }}
      >
        {label}
      </span>
      <span
        style={{
          color: "#fff8d7",
          fontSize: "clamp(0.74rem, 1.6vw, 0.95rem)",
          lineHeight: 1.5,
          overflowWrap: "anywhere",
        }}
      >
        {value}
      </span>
    </div>
  );
}

export function Home({ onNavigate }: HomeProps) {
  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        backgroundImage:
          "linear-gradient(180deg, rgba(4, 8, 18, 0.25), rgba(4, 8, 18, 0.56)), url('/home_lunch.png')",
        backgroundSize: "cover",
        backgroundPosition: "center center",
        backgroundRepeat: "no-repeat",
        imageRendering: "pixelated",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#f4ecd0",
      }}
    >
      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.1), rgba(0,0,0,0.1) 1px, transparent 1px, transparent 4px)",
          mixBlendMode: "multiply",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          boxShadow: "inset 0 0 90px rgba(0,0,0,0.7)",
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <div
        style={{
          position: "relative",
          zIndex: 3,
          minHeight: "100svh",
          display: "grid",
          gridTemplateRows: "auto 1fr auto",
          padding: "clamp(16px, 3vw, 32px)",
          gap: "20px",
        }}
      >
        <header
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
            gap: "18px",
            flexWrap: "wrap",
          }}
        >
          <HudPanel>
            <div
              style={{
                color: "#ffd700",
                fontSize: "0.64rem",
                lineHeight: 1.6,
                marginBottom: "10px",
              }}
            >
              PLAYER INFO
            </div>
            <div style={{ display: "grid", gap: "6px" }}>
              <LabelValue label="NAME" value={player.name} />
              <LabelValue label="TITLE" value={player.title} />
              <LabelValue label="LEVEL" value={`LV.${player.level}`} />
            </div>
          </HudPanel>

          <HudPanel align="right">
            <div
              style={{
                color: "#00f5ff",
                fontSize: "0.64rem",
                lineHeight: 1.6,
                marginBottom: "10px",
              }}
            >
              CONTRIBUTION POWER
            </div>
            <div style={{ display: "grid", gap: "6px" }}>
              <LabelValue label="TOTAL CP" value={player.totalCp.toLocaleString()} />
              <LabelValue label="TODAY CP" value={`+${player.todayCp.toLocaleString()}`} />
            </div>
          </HudPanel>
        </header>

        <section
          aria-label="Character placement area"
          style={{
            minHeight: "clamp(220px, 42vh, 520px)",
          }}
        />

        <nav
          aria-label="Main navigation"
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 220px), 1fr))",
            gap: "14px",
            width: "min(100%, 900px)",
            justifySelf: "center",
          }}
        >
          {navItems.map((item, index) => (
            <motion.button
              key={item.label}
              custom={index}
              variants={navVariants}
              initial="hidden"
              animate="visible"
              whileHover={{ y: -4, scale: 1.02 }}
              whileTap={{ y: 2, scale: 0.98 }}
              onClick={() => onNavigate(item.path)}
              style={{
                minHeight: "86px",
                border: `3px solid ${item.accent}`,
                borderBottomColor: "rgba(0,0,0,0.75)",
                borderRightColor: "rgba(0,0,0,0.75)",
                background: "rgba(4, 11, 27, 0.82)",
                boxShadow: `0 0 0 2px rgba(0,0,0,0.78), 6px 6px 0 rgba(0,0,0,0.5), inset 0 0 18px ${item.accent}22`,
                color: "#fff8d7",
                cursor: "pointer",
                fontFamily: "inherit",
                padding: "16px 18px",
                textAlign: "left",
              }}
            >
              <span
                style={{
                  display: "block",
                  color: item.accent,
                  fontSize: "clamp(0.88rem, 2vw, 1.08rem)",
                  lineHeight: 1.5,
                  marginBottom: "8px",
                  overflowWrap: "anywhere",
                }}
              >
                {item.label}
              </span>
              <span
                style={{
                  display: "block",
                  color: "rgba(244, 236, 208, 0.66)",
                  fontSize: "0.6rem",
                  lineHeight: 1.5,
                }}
              >
                {item.caption}
              </span>
            </motion.button>
          ))}
        </nav>
      </div>
    </main>
  );
}
