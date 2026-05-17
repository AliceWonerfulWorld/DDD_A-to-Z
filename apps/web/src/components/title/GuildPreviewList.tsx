import { motion } from "framer-motion";

const GUILDS = [
  { name: "Rust", color: "#ff6b35", icon: "🦀" },
  { name: "Python", color: "#3776ab", icon: "🐍" },
  { name: "Go", color: "#00acd7", icon: "🐹" },
  { name: "TypeScript", color: "#3178c6", icon: "📘" },
  { name: "Java", color: "#f89820", icon: "☕" },
  { name: "Haskell", color: "#5d4f85", icon: "λ" },
  { name: "Zig", color: "#f7a41d", icon: "⚡" },
] as const;

export function GuildPreviewList() {
  return (
    <motion.div
      style={{
        display: "flex",
        flexWrap: "wrap",
        gap: "12px",
        justifyContent: "center",
      }}
    >
      {GUILDS.map((guild, i) => (
        <motion.div
          key={guild.name}
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.8 + i * 0.08, duration: 0.4 }}
          whileHover={{ y: -4, scale: 1.05 }}
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            gap: "10px",
            position: "relative",
            overflow: "hidden",
            padding: "16px 18px",
            border: `1px solid ${guild.color}cc`,
            boxShadow: `0 0 0 1px rgba(0,0,0,0.72), 3px 3px 0 rgba(0,0,0,0.58), 0 0 18px ${guild.color}2f, inset 0 0 24px rgba(0,0,0,0.42)`,
            background: `linear-gradient(180deg, rgba(9, 13, 20, 0.7) 0%, rgba(9, 13, 20, 0.5) 100%), radial-gradient(circle at 50% 0%, ${guild.color}38 0%, transparent 58%)`,
            backdropFilter: "blur(1.5px)",
            cursor: "pointer",
            minHeight: "112px",
            minWidth: "104px",
          }}
        >
          <span
            aria-hidden="true"
            style={{
              position: "absolute",
              top: 0,
              left: "14px",
              right: "14px",
              height: "2px",
              background: guild.color,
              boxShadow: `0 0 10px ${guild.color}`,
            }}
          />
          <span
            style={{
              position: "relative",
              zIndex: 1,
              fontSize: "2rem",
              lineHeight: 1,
              filter: `drop-shadow(0 2px 0 rgba(0,0,0,0.9)) drop-shadow(0 0 10px ${guild.color}55)`,
            }}
          >
            {guild.icon}
          </span>
          <span
            style={{
              position: "relative",
              zIndex: 1,
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.76rem",
              color: "#fff7dc",
              letterSpacing: "0.07em",
              textShadow: `1px 1px 0 #000, 0 0 10px ${guild.color}`,
            }}
          >
            {guild.name}
          </span>
        </motion.div>
      ))}
    </motion.div>
  );
}
