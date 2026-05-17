import { AnimatePresence, motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import type { WarGuild } from "./WarMapData";

interface ScoutPanelProps {
  guild: WarGuild | null;
  isCurrentGuild?: boolean;
  onClose: () => void;
}

export function ScoutPanel({ guild, isCurrentGuild = false, onClose }: ScoutPanelProps) {
  return (
    <AnimatePresence>
      {guild && (
        <motion.aside
          key={guild.id}
          data-war-interactive="true"
          initial={{ x: 420, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          exit={{ x: 420, opacity: 0 }}
          transition={{ duration: 0.3, ease: steppedEase(6) }}
          style={{
            position: "fixed",
            right: "clamp(12px, 2vw, 24px)",
            top: "calc(env(safe-area-inset-top, 0px) + clamp(78px, 9vw, 96px))",
            zIndex: 12,
            width: "min(calc(100vw - 24px), 392px)",
            border: `3px solid ${guild.color}`,
            borderBottomColor: "rgba(20, 18, 12, 0.98)",
            borderRightColor: "rgba(20, 18, 12, 0.98)",
            background: "linear-gradient(180deg, rgba(2, 8, 20, 0.94), rgba(0, 0, 0, 0.86))",
            boxShadow: `0 0 0 2px rgba(0,0,0,0.76), 8px 8px 0 rgba(0,0,0,0.36), inset 0 0 24px ${guild.color}22`,
            color: "#fff8d7",
            padding: "16px",
            pointerEvents: "auto",
            backdropFilter: "blur(2px)",
          }}
        >
          <motion.button
            type="button"
            aria-label="Close scout panel"
            onPointerDown={(event) => event.stopPropagation()}
            onClick={(event) => {
              event.stopPropagation();
              onClose();
            }}
            whileHover={{ y: -1, backgroundColor: "rgba(255, 217, 102, 0.16)" }}
            whileTap={{ y: 1, scale: 0.96 }}
            style={{
              position: "absolute",
              right: "10px",
              top: "10px",
              width: "30px",
              height: "30px",
              border: "2px solid rgba(255, 217, 102, 0.72)",
              borderBottomColor: "rgba(96, 62, 22, 0.95)",
              borderRightColor: "rgba(96, 62, 22, 0.95)",
              background: "rgba(3, 10, 24, 0.78)",
              boxShadow: "0 0 0 1px rgba(0,0,0,0.62)",
              color: "#fff8d7",
              cursor: "pointer",
              fontFamily: "inherit",
              fontSize: "0.66rem",
              lineHeight: 1,
              padding: 0,
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            }}
          >
            x
          </motion.button>

          <p
            style={{
              margin: "0 0 10px",
              color: guild.accent,
              fontSize: "0.52rem",
              lineHeight: 1.45,
              textShadow: `0 0 12px ${guild.color}`,
            }}
          >
            SCOUT REPORT
          </p>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "72px minmax(0, 1fr)",
              gap: "12px",
              alignItems: "center",
              marginBottom: "14px",
              paddingRight: "34px",
            }}
          >
            <div
              aria-hidden="true"
              style={{
                display: "grid",
                width: "72px",
                height: "64px",
                placeItems: "center",
                border: `2px solid ${guild.color}`,
                background: "rgba(0,0,0,0.44)",
                boxShadow: `inset 0 0 14px rgba(0,0,0,0.68), 0 0 12px ${guild.color}55`,
                color: guild.accent,
                fontSize: "0.9rem",
              }}
            >
              {guild.mark}
            </div>
            <div style={{ minWidth: 0 }}>
              <h2
                style={{
                  margin: "0 0 8px",
                  color: "#ffd966",
                  fontSize: "clamp(0.78rem, 1.8vw, 1rem)",
                  lineHeight: 1.45,
                  overflowWrap: "anywhere",
                  textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
                }}
              >
                {guild.name} GUILD
              </h2>
              <span
                style={{
                  display: "inline-block",
                  border: `2px solid ${guild.color}`,
                  background: `${guild.color}24`,
                  color: guild.accent,
                  fontSize: "0.48rem",
                  lineHeight: 1,
                  padding: "6px 7px",
                }}
              >
                {guild.color}
              </span>
              {isCurrentGuild && (
                <span
                  style={{
                    display: "inline-block",
                    marginLeft: "8px",
                    border: `2px solid ${guild.color}`,
                    background: `${guild.color}24`,
                    color: guild.accent,
                    fontSize: "0.48rem",
                    lineHeight: 1,
                    padding: "6px 7px",
                  }}
                >
                  YOUR GUILD
                </span>
              )}
            </div>
          </div>

          <p
            style={{
              margin: "0 0 16px",
              color: "#f4ecd0",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "clamp(0.82rem, 1.45vw, 0.98rem)",
              lineHeight: 1.55,
            }}
          >
            {guild.description}
          </p>

          <dl
            style={{
              display: "grid",
              gridTemplateColumns: "1fr 1fr",
              gap: "10px",
              margin: 0,
            }}
          >
            <Stat label="TOTAL CP" value={guild.totalCp.toLocaleString()} color={guild.color} />
            <Stat label="MEMBERS" value={guild.memberCount.toLocaleString()} color={guild.color} />
          </dl>
        </motion.aside>
      )}
    </AnimatePresence>
  );
}

function Stat({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div
      style={{
        border: "2px solid rgba(255,255,255,0.12)",
        background: "rgba(255,255,255,0.04)",
        padding: "10px",
      }}
    >
      <dt
        style={{
          color,
          fontSize: "0.48rem",
          lineHeight: 1.35,
          marginBottom: "7px",
          textShadow: `0 0 10px ${color}66`,
        }}
      >
        {label}
      </dt>
      <dd
        style={{
          color: "#fff8d7",
          fontSize: "0.66rem",
          lineHeight: 1.35,
          margin: 0,
          overflowWrap: "anywhere",
        }}
      >
        {value}
      </dd>
    </div>
  );
}
