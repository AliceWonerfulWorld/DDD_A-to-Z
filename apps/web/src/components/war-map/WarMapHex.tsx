import { motion, useMotionValue, useTransform } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import type { WarGuild } from "./WarMapData";
import { useTransformEffect } from "react-zoom-pan-pinch";

interface WarMapHexProps {
  guild: WarGuild;
  isCurrentGuild: boolean;
  isSelected: boolean;
  onSelect: (guild: WarGuild) => void;
}

export function WarMapHex({ guild, isCurrentGuild, isSelected, onSelect }: WarMapHexProps) {
  const zoomScale = useMotionValue(1);
  const inverseScale = useTransform(zoomScale, (s) => (s < 1 ? 1 / s : 1));

  useTransformEffect(({ state }) => {
    zoomScale.set(state.scale);
  });

  return (
    <motion.button
      type="button"
      aria-label={`Scout ${guild.name} guild`}
      data-war-interactive="true"
      onPointerDown={(event) => event.stopPropagation()}
      onClick={(event) => {
        event.stopPropagation();
        onSelect(guild);
      }}
      initial={{ opacity: 0, scale: 0.72 }}
      animate={{ opacity: 1, scale: isSelected ? 1.04 : 1 }}
      transition={{ duration: 0.26, ease: steppedEase(6) }}
      whileHover={{ scale: isSelected ? 1.06 : 1.04 }}
      whileTap={{ scale: 0.96 }}
      style={{
        position: "absolute",
        left: `${guild.x}%`,
        top: `${guild.y}%`,
        width: "clamp(48px, 4.8vw, 74px)",
        height: "clamp(56px, 5.6vw, 86px)",
        border: 0,
        background: "transparent",
        color: "#fff8d7",
        cursor: "pointer",
        fontFamily: "inherit",
        padding: 0,
        transformOrigin: "50% 50%",
        zIndex: isSelected ? 8 : 6,
        touchAction: "none",
        x: "-50%",
        y: "-50%",
      }}
    >
      <motion.div style={{ width: "100%", height: "100%", scale: inverseScale }}>
        {isCurrentGuild && (
          <span
            style={{
              position: "absolute",
              left: "50%",
              bottom: "calc(100% + 6px)",
              border: `2px solid ${guild.color}`,
              background: "rgba(0, 0, 0, 0.84)",
              boxShadow: `0 0 12px ${guild.color}55`,
              color: guild.accent,
              fontSize: "clamp(0.3rem, 0.6vw, 0.42rem)",
              lineHeight: 1,
              padding: "4px 5px",
              textShadow: "2px 2px 0 rgba(0,0,0,0.76)",
              transform: "translateX(-50%)",
              whiteSpace: "nowrap",
            }}
          >
            MY GUILD
          </span>
        )}
        <svg viewBox="0 0 100 112" aria-hidden="true" style={{ display: "block", width: "100%" }}>
          <defs>
            <linearGradient id={`${guild.id}-hex-fill`} x1="0" x2="1" y1="0" y2="1">
              <stop offset="0%" stopColor="rgba(4, 14, 30, 0.96)" />
              <stop offset="100%" stopColor="rgba(0, 0, 0, 0.78)" />
            </linearGradient>
          </defs>
          <polygon
            points="50 4 92 28 92 78 50 108 8 78 8 28"
            fill={`url(#${guild.id}-hex-fill)`}
            stroke={guild.color}
            strokeWidth={isSelected || isCurrentGuild ? 5 : 3}
            style={{
              filter: `drop-shadow(0 0 8px ${guild.color}) drop-shadow(0 0 18px ${guild.color}88)`,
            }}
          />
          <polygon
            points="50 17 80 35 80 72 50 92 20 72 20 35"
            fill="rgba(255,255,255,0.05)"
            stroke={guild.accent}
            strokeWidth="1.5"
            opacity="0.9"
          />
          <line x1="50" x2="50" y1="18" y2="92" stroke={guild.color} strokeOpacity="0.22" />
          <line x1="21" x2="79" y1="54" y2="54" stroke={guild.color} strokeOpacity="0.22" />
        </svg>
        <span
          style={{
            position: "absolute",
            inset: "0 0 8px",
            display: "grid",
            placeItems: "center",
            color: guild.accent,
            fontSize: "clamp(0.58rem, 1.05vw, 0.84rem)",
            lineHeight: 1,
            overflow: "hidden",
          }}
        >
          <img
            src={`/guild-icons/${guild.mark === "HS" ? "λ" : guild.mark}.png`}
            alt={`${guild.name} icon`}
            style={{
              width: "clamp(28px, 2.8vw, 44px)",
              height: "clamp(28px, 2.8vw, 44px)",
              objectFit: "cover",
              imageRendering: "pixelated",
            }}
          />
        </span>{" "}
        <span
          style={{
            position: "absolute",
            left: "50%",
            top: "calc(100% - 8px)",
            maxWidth: "108px",
            padding: "4px 6px",
            border: `2px solid ${guild.color}`,
            background: "rgba(0, 0, 0, 0.76)",
            boxShadow: `0 0 12px ${guild.color}66`,
            color: "#fff8d7",
            fontSize: "clamp(0.36rem, 0.72vw, 0.48rem)",
            lineHeight: 1.35,
            overflowWrap: "anywhere",
            textShadow: "2px 2px 0 rgba(0,0,0,0.8)",
            transform: "translateX(-50%)",
            whiteSpace: "nowrap",
          }}
        >
          {guild.name}
        </span>
      </motion.div>
    </motion.button>
  );
}
