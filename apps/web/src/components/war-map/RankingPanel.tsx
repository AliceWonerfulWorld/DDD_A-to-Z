import { AnimatePresence, motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import { rankWarGuilds, type WarGuild } from "./WarMapData";

interface RankingPanelProps {
  currentGuildID: string | null;
  guilds: WarGuild[];
  isOpen: boolean;
  isLoaded: boolean;
  onToggle: () => void;
}

const rankColors = ["#ffd966", "#c8f2ff", "#ffad66"];

export function RankingPanel({
  currentGuildID,
  guilds,
  isOpen,
  isLoaded,
  onToggle,
}: RankingPanelProps) {
  const rankedGuilds = rankWarGuilds(guilds);

  return (
    <div
      style={{
        position: "fixed",
        left: "clamp(12px, 2vw, 24px)",
        top: "calc(env(safe-area-inset-top, 0px) + clamp(72px, 7vw, 88px))",
        zIndex: 12,
        pointerEvents: "none",
      }}
    >
      <motion.button
        type="button"
        data-war-interactive="true"
        onPointerDown={(event) => event.stopPropagation()}
        onClick={(event) => {
          event.stopPropagation();
          onToggle();
        }}
        whileHover={{ y: -2, backgroundColor: "rgba(0, 245, 255, 0.14)" }}
        whileTap={{ y: 1, scale: 0.98 }}
        style={{
          minHeight: "42px",
          border: "2px solid rgba(0, 245, 255, 0.86)",
          borderBottomColor: "rgba(0, 49, 64, 0.96)",
          borderRightColor: "rgba(0, 49, 64, 0.96)",
          background: "rgba(2, 10, 24, 0.9)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.76), 0 0 18px rgba(0,245,255,0.28)",
          color: "#d9fbff",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.5rem, 1.5vw, 0.62rem)",
          lineHeight: 1,
          padding: "clamp(8px, 2vw, 12px) clamp(10px, 2.5vw, 14px)",
          pointerEvents: "auto",
          textShadow: "2px 2px 0 rgba(0,0,0,0.8)",
        }}
      >
        [ RANKING ]
      </motion.button>

      <AnimatePresence>
        {isOpen && (
          <motion.aside
            key="ranking-panel"
            data-war-interactive="true"
            initial={{ x: -360, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: -360, opacity: 0 }}
            transition={{ duration: 0.28, ease: steppedEase(6) }}
            style={{
              width: "min(calc(100vw - 24px), 360px)",
              marginTop: "12px",
              border: "3px solid rgba(0, 245, 255, 0.72)",
              borderBottomColor: "rgba(0, 49, 64, 0.96)",
              borderRightColor: "rgba(0, 49, 64, 0.96)",
              background: "linear-gradient(180deg, rgba(1, 8, 20, 0.92), rgba(0, 0, 0, 0.84))",
              boxShadow:
                "0 0 0 2px rgba(0,0,0,0.76), 8px 8px 0 rgba(0,0,0,0.36), inset 0 0 22px rgba(0,245,255,0.1)",
              color: "#fff8d7",
              maxHeight: "min(70vh, 540px)",
              overflow: "auto",
              padding: "14px",
              pointerEvents: "auto",
              backdropFilter: "blur(2px)",
            }}
          >
            <h2
              style={{
                margin: "0 0 12px",
                color: "#00f5ff",
                fontSize: "0.72rem",
                lineHeight: 1.45,
                textShadow: "2px 2px 0 rgba(0,0,0,0.8), 0 0 12px rgba(0,245,255,0.55)",
              }}
            >
              GUILD RANKING
            </h2>
            {rankedGuilds.length === 0 && (
              <p
                style={{
                  margin: 0,
                  color: "rgba(255,248,215,0.72)",
                  fontSize: "0.56rem",
                  lineHeight: 1.7,
                }}
              >
                {isLoaded ? "NO GUILD DATA" : "SYNCING GUILD DATA..."}
              </p>
            )}
            <ol
              style={{
                display: "grid",
                gap: "8px",
                listStyle: "none",
                margin: 0,
                padding: 0,
              }}
            >
              {rankedGuilds.map((guild, index) => {
                const rankColor = rankColors[index] ?? "rgba(255, 248, 215, 0.64)";
                const isCurrentGuild = guild.id === currentGuildID;

                return (
                  <li
                    key={guild.id}
                    style={{
                      display: "grid",
                      gridTemplateColumns: "34px 42px minmax(0, 1fr)",
                      alignItems: "center",
                      gap: "9px",
                      minHeight: "54px",
                      border: `2px solid ${isCurrentGuild ? guild.color : "rgba(255,255,255,0.12)"}`,
                      background: isCurrentGuild
                        ? `linear-gradient(90deg, ${guild.color}38, rgba(0,0,0,0.46))`
                        : "rgba(255,255,255,0.04)",
                      boxShadow: isCurrentGuild ? `inset 0 0 14px ${guild.color}25` : "none",
                      padding: "8px",
                    }}
                  >
                    <span
                      style={{
                        color: rankColor,
                        fontSize: "0.68rem",
                        textShadow: `0 0 10px ${rankColor}`,
                      }}
                    >
                      #{index + 1}
                    </span>
                    <span
                      style={{
                        display: "grid",
                        width: "38px",
                        height: "38px",
                        placeItems: "center",
                        border: `2px solid ${guild.color}`,
                        background: `radial-gradient(circle at 50% 20%, ${guild.color}36, rgba(0,0,0,0.78) 62%)`,
                        color: guild.accent,
                        fontSize: "0.58rem",
                        boxShadow: `0 0 10px ${guild.color}55`,
                        overflow: "hidden",
                      }}
                    >
                      <img
                        src={`/guild-icons/${guild.mark === "HS" ? "λ" : guild.mark}.png`}
                        alt={`${guild.name} icon`}
                        style={{
                          width: "100%",
                          height: "100%",
                          objectFit: "cover",
                          imageRendering: "pixelated",
                        }}
                      />
                    </span>
                    <span style={{ minWidth: 0 }}>
                      <strong
                        style={{
                          display: "block",
                          color: "#fff8d7",
                          fontSize: "0.58rem",
                          lineHeight: 1.45,
                          overflowWrap: "anywhere",
                        }}
                      >
                        {guild.name}
                      </strong>
                      <span
                        style={{
                          display: "block",
                          color: guild.color,
                          fontSize: "0.5rem",
                          lineHeight: 1.45,
                          textShadow: `0 0 10px ${guild.color}66`,
                        }}
                      >
                        {guild.totalCp.toLocaleString()} CP
                      </span>
                      {isCurrentGuild && (
                        <span
                          style={{
                            display: "inline-block",
                            marginTop: "6px",
                            border: `2px solid ${guild.color}`,
                            background: `${guild.color}24`,
                            color: guild.accent,
                            fontSize: "0.42rem",
                            lineHeight: 1,
                            padding: "4px 6px",
                          }}
                        >
                          YOUR GUILD
                        </span>
                      )}
                    </span>
                  </li>
                );
              })}
            </ol>
          </motion.aside>
        )}
      </AnimatePresence>
    </div>
  );
}
