import { motion } from "framer-motion";
import { steppedEase } from "../lib/animationUtils";

interface GuildPlaceholderPageProps {
  title: string;
  caption: string;
  onNavigate: (path: string) => void;
}

export function GuildPlaceholderPage({ title, caption, onNavigate }: GuildPlaceholderPageProps) {
  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        display: "grid",
        placeItems: "center",
        padding: "24px",
        backgroundImage:
          "linear-gradient(180deg, rgba(4, 8, 18, 0.28), rgba(4, 8, 18, 0.68)), url('/pixel-town-night.png')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        imageRendering: "pixelated",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#f4ecd0",
      }}
    >
      <motion.section
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.34, ease: steppedEase(6) }}
        style={{
          width: "min(100%, 560px)",
          border: "3px solid rgba(0, 245, 255, 0.78)",
          borderBottomColor: "rgba(2, 54, 72, 0.96)",
          borderRightColor: "rgba(2, 54, 72, 0.96)",
          background: "rgba(3, 10, 24, 0.84)",
          boxShadow:
            "0 0 0 2px rgba(0,0,0,0.76), 8px 8px 0 rgba(0,0,0,0.45), inset 0 0 22px rgba(0,245,255,0.12)",
          padding: "clamp(20px, 5vw, 34px)",
          textAlign: "center",
        }}
      >
        <p
          style={{
            margin: "0 0 14px",
            color: "#74f7a1",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(0.74rem, 2.3vw, 1rem)",
            lineHeight: 1.6,
          }}
        >
          ROUTE READY
        </p>
        <h1
          style={{
            margin: "0 0 18px",
            color: "#fff8d7",
            fontSize: "clamp(1rem, 4vw, 1.5rem)",
            lineHeight: 1.7,
            overflowWrap: "anywhere",
          }}
        >
          {title}
        </h1>
        <p
          style={{
            margin: "0 0 28px",
            color: "rgba(244, 236, 208, 0.74)",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(0.9rem, 2.6vw, 1.1rem)",
            lineHeight: 1.8,
          }}
        >
          {caption}
        </p>
        <motion.button
          type="button"
          whileHover={{ y: -2, scale: 1.02 }}
          whileTap={{ y: 2, scale: 0.98 }}
          onClick={() => onNavigate("/guild")}
          style={{
            border: "2px solid rgba(255, 217, 102, 0.78)",
            borderBottomColor: "rgba(96, 62, 22, 0.95)",
            borderRightColor: "rgba(96, 62, 22, 0.95)",
            background: "rgba(3, 10, 24, 0.8)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 5px 5px 0 rgba(0,0,0,0.38)",
            color: "#fff8d7",
            cursor: "pointer",
            fontFamily: "inherit",
            fontSize: "clamp(0.58rem, 1.8vw, 0.78rem)",
            lineHeight: 1.5,
            padding: "10px 12px",
          }}
        >
          &lt; GUILD BASE
        </motion.button>
      </motion.section>
    </main>
  );
}
