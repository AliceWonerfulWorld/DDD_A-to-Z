import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { PATHS } from "../../constants/paths";
import { GUILD_MASTERS, type GuildMaster } from "../../features/guild/guildMaster";
import { fetchGuilds, joinGuild as joinGuildAPI } from "../../features/guild/api";
import { hasEverJoinedGuild } from "../../features/guild/membership";
import { toDisplayGuilds } from "../../features/guild/presentation";
import { normalizeGrantedPet, storeGrantedPet } from "../../features/pet/guildGrant";
import { ApiError } from "../../lib/api/client";
import { steppedEase } from "../../lib/animationUtils";

interface GuildSelectionProps {
  onNavigate: (path: string) => void;
}

export function GuildSelection({ onNavigate }: GuildSelectionProps) {
  const [guilds, setGuilds] = useState<GuildMaster[]>(GUILD_MASTERS);
  const [selectedGuild, setSelectedGuild] = useState<GuildMaster>(GUILD_MASTERS[0]);
  const [isJoining, setIsJoining] = useState(false);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const systemMessage = hasEverJoinedGuild()
    ? "次なる戦いの舞台へ。今シーズンの所属を選択してください"
    : "Welcome to the World. Choose your faction.";

  useEffect(() => {
    let isMounted = true;

    fetchGuilds()
      .then((apiGuilds) => {
        const displayGuilds = toDisplayGuilds(apiGuilds);
        if (!isMounted || displayGuilds.length === 0) {
          return;
        }

        setGuilds(displayGuilds);
        setSelectedGuild(displayGuilds[0]);
      })
      .catch((error) => {
        if (!isMounted) {
          return;
        }

        console.error("failed to fetch guilds", error);
        setStatusMessage("ギルド一覧を取得できませんでした。表示中の候補から再試行してください。");
      });

    return () => {
      isMounted = false;
    };
  }, []);

  const joinGuild = async () => {
    if (isJoining) {
      return;
    }

    setIsJoining(true);
    setStatusMessage(null);
    try {
      const result = await joinGuildAPI(selectedGuild.id);
      const grantedPet = result.granted_pet ?? result.grantedPet;
      const petAlreadyOwned = result.pet_already_owned ?? result.petAlreadyOwned ?? false;
      if (grantedPet && !petAlreadyOwned) {
        storeGrantedPet(normalizeGrantedPet(grantedPet));
      }
      onNavigate(PATHS.GUILD);
    } catch (error) {
      if (error instanceof ApiError && error.status === 409) {
        onNavigate(PATHS.GUILD);
        return;
      }

      console.error("failed to join guild", error);
      setStatusMessage("ギルド参加に失敗しました。少し時間を置いて再度お試しください。");
    } finally {
      setIsJoining(false);
    }
  };

  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        background:
          "radial-gradient(circle at 50% 12%, rgba(0,245,255,0.14), transparent 32%), linear-gradient(180deg, rgba(1,5,12,0.96), rgba(0,0,0,0.92))",
        color: "#fff8d7",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
      }}
    >
      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.035) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.035) 1px, transparent 1px)",
          backgroundSize: "36px 36px",
          maskImage: "linear-gradient(180deg, rgba(0,0,0,0.82), transparent)",
          pointerEvents: "none",
        }}
      />
      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.16), rgba(0,0,0,0.16) 1px, transparent 1px, transparent 4px)",
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <section
        aria-labelledby="guild-selection-title"
        style={{
          position: "relative",
          zIndex: 3,
          width: "min(1180px, calc(100% - 28px))",
          minHeight: "100svh",
          margin: "0 auto",
          display: "grid",
          gridTemplateRows: "auto 1fr",
          gap: "18px",
          padding: "clamp(18px, 4vw, 42px) 0",
        }}
      >
        <motion.header
          initial={{ opacity: 0, y: -14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, ease: steppedEase(6) }}
          style={{
            border: "3px solid rgba(255,217,102,0.82)",
            borderBottomColor: "rgba(37, 28, 10, 0.98)",
            borderRightColor: "rgba(37, 28, 10, 0.98)",
            background: "rgba(3, 10, 24, 0.78)",
            boxShadow:
              "0 0 0 2px rgba(0,0,0,0.76), 7px 7px 0 rgba(0,0,0,0.42), inset 0 0 22px rgba(255,217,102,0.08)",
            padding: "16px clamp(14px, 3vw, 24px)",
          }}
        >
          <p
            style={{
              margin: "0 0 10px",
              color: selectedGuild.accent,
              fontSize: "0.55rem",
              lineHeight: 1.5,
              textShadow: `0 0 12px ${selectedGuild.color}`,
            }}
          >
            GUILD ENTRY TERMINAL
          </p>
          <h1
            id="guild-selection-title"
            style={{
              margin: 0,
              color: "#fff8d7",
              fontSize: "clamp(0.9rem, 1.8vw, 1.28rem)",
              lineHeight: 1.7,
              overflowWrap: "anywhere",
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            }}
          >
            {systemMessage}
          </h1>
        </motion.header>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 320px), 1fr))",
            gap: "18px",
            alignItems: "stretch",
          }}
        >
          <GuildList guilds={guilds} selectedGuild={selectedGuild} onSelect={setSelectedGuild} />
          <GuildDetail
            guild={selectedGuild}
            isJoining={isJoining}
            onJoin={() => void joinGuild()}
            statusMessage={statusMessage}
          />
        </div>
      </section>
    </main>
  );
}

function GuildList({
  guilds,
  selectedGuild,
  onSelect,
}: {
  guilds: GuildMaster[];
  selectedGuild: GuildMaster;
  onSelect: (guild: GuildMaster) => void;
}) {
  return (
    <motion.aside
      initial={{ opacity: 0, x: -24 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.35, ease: steppedEase(6) }}
      aria-label="Guild list"
      style={{
        minHeight: 0,
        border: "3px solid rgba(0,245,255,0.64)",
        borderBottomColor: "rgba(2, 12, 20, 0.98)",
        borderRightColor: "rgba(2, 12, 20, 0.98)",
        background: "rgba(0,0,0,0.78)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.76), 7px 7px 0 rgba(0,0,0,0.38), inset 0 0 26px rgba(0,245,255,0.08)",
        padding: "12px",
        overflow: "hidden",
      }}
    >
      <p
        style={{
          margin: "0 0 12px",
          color: "#00f5ff",
          fontSize: "0.52rem",
          lineHeight: 1.45,
          textShadow: "0 0 12px rgba(0,245,255,0.78)",
        }}
      >
        FORMATION LIST
      </p>
      <div
        style={{
          display: "grid",
          gap: "10px",
          maxHeight: "calc(100svh - 218px)",
          overflowY: "auto",
          paddingRight: "4px",
        }}
      >
        {guilds.map((guild, index) => {
          const isSelected = guild.slug === selectedGuild.slug;

          return (
            <motion.button
              key={guild.slug}
              type="button"
              initial={{ opacity: 0, x: -18 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.08 + index * 0.04, duration: 0.25 }}
              whileHover={{ x: 3 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => onSelect(guild)}
              aria-pressed={isSelected}
              style={{
                display: "grid",
                gridTemplateColumns: "54px minmax(0, 1fr) auto",
                alignItems: "center",
                gap: "10px",
                minHeight: "68px",
                border: `2px solid ${isSelected ? guild.color : "rgba(255,255,255,0.12)"}`,
                background: isSelected
                  ? `linear-gradient(90deg, ${guild.color}34, rgba(4,11,27,0.84))`
                  : "rgba(4, 11, 27, 0.7)",
                boxShadow: isSelected
                  ? `0 0 16px ${guild.color}55, inset 0 0 18px ${guild.color}22`
                  : "inset 0 0 12px rgba(0,0,0,0.42)",
                color: "#fff8d7",
                cursor: "pointer",
                fontFamily: "inherit",
                padding: "9px 10px",
                textAlign: "left",
              }}
            >
              <GuildMark guild={guild} size={44} />
              <span style={{ minWidth: 0 }}>
                <span
                  style={{
                    display: "block",
                    color: isSelected ? guild.accent : "#fff8d7",
                    fontSize: "0.62rem",
                    lineHeight: 1.45,
                    overflowWrap: "anywhere",
                    textShadow: isSelected ? `0 0 10px ${guild.color}` : "none",
                  }}
                >
                  {guild.name}
                </span>
                <span
                  style={{
                    display: "block",
                    marginTop: "6px",
                    color: "rgba(244,236,208,0.62)",
                    fontFamily: '"DotGothic16", monospace',
                    fontSize: "0.78rem",
                    lineHeight: 1.35,
                  }}
                >
                  {guild.oath}
                </span>
              </span>
              <span
                style={{
                  color: guild.accent,
                  fontSize: "0.52rem",
                  lineHeight: 1.35,
                  textAlign: "right",
                  whiteSpace: "nowrap",
                }}
              >
                {guild.memberCount}
              </span>
            </motion.button>
          );
        })}
      </div>
    </motion.aside>
  );
}

function GuildDetail({
  guild,
  isJoining,
  onJoin,
  statusMessage,
}: {
  guild: GuildMaster;
  isJoining: boolean;
  onJoin: () => void;
  statusMessage: string | null;
}) {
  return (
    <motion.article
      key={guild.slug}
      initial={{ opacity: 0, x: 24 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.32, ease: steppedEase(6) }}
      aria-live="polite"
      style={{
        display: "grid",
        gridTemplateRows: "auto auto 1fr auto",
        minHeight: "min(650px, calc(100svh - 170px))",
        border: `3px solid ${guild.color}`,
        borderBottomColor: "rgba(20, 18, 12, 0.98)",
        borderRightColor: "rgba(20, 18, 12, 0.98)",
        background: "linear-gradient(180deg, rgba(4, 10, 22, 0.94), rgba(0,0,0,0.88))",
        boxShadow: `0 0 0 2px rgba(0,0,0,0.76), 8px 8px 0 rgba(0,0,0,0.4), inset 0 0 30px ${guild.color}22`,
        padding: "clamp(16px, 3vw, 26px)",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "minmax(88px, 124px) minmax(0, 1fr)",
          gap: "16px",
          alignItems: "center",
        }}
      >
        <GuildMark guild={guild} size={112} />
        <div style={{ minWidth: 0 }}>
          <p
            style={{
              margin: "0 0 10px",
              color: guild.accent,
              fontSize: "0.52rem",
              lineHeight: 1.45,
              textShadow: `0 0 12px ${guild.color}`,
            }}
          >
            SELECTED FACTION
          </p>
          <h2
            style={{
              margin: 0,
              color: "#ffd966",
              fontSize: "clamp(1rem, 2.4vw, 1.7rem)",
              lineHeight: 1.45,
              overflowWrap: "anywhere",
              textShadow: `2px 2px 0 rgba(0,0,0,0.72), 0 0 18px ${guild.color}`,
            }}
          >
            {guild.name} Guild
          </h2>
        </div>
      </div>

      <section style={{ marginTop: "22px" }}>
        <p
          style={{
            margin: "0 0 10px",
            color: guild.accent,
            fontSize: "0.52rem",
            lineHeight: 1.45,
          }}
        >
          MANIFESTO
        </p>
        <p
          style={{
            margin: 0,
            color: "#f4ecd0",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "1rem",
            lineHeight: 1.7,
          }}
        >
          {guild.description}
        </p>
      </section>

      <dl
        style={{
          alignSelf: "end",
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(150px, 1fr))",
          gap: "12px",
          margin: "24px 0",
        }}
      >
        <Stat label="MEMBERS" value={guild.memberCount.toLocaleString()} color={guild.color} />
        <Stat
          label="LAST SEASON CP"
          value={guild.previousSeasonCp.toLocaleString()}
          color={guild.color}
        />
        <Stat label="GUILD CODE" value={guild.slug.toUpperCase()} color={guild.color} />
      </dl>

      <motion.button
        type="button"
        disabled={isJoining}
        onClick={onJoin}
        whileHover={{
          y: -3,
          boxShadow: `0 0 0 2px rgba(0,0,0,0.78), 0 0 22px ${guild.color}, 0 0 46px ${guild.color}88, inset 0 0 24px ${guild.color}44`,
        }}
        whileTap={{ y: 1, scale: 0.98 }}
        style={{
          width: "100%",
          minHeight: "72px",
          border: `3px solid ${guild.color}`,
          borderBottomColor: "rgba(0,0,0,0.82)",
          borderRightColor: "rgba(0,0,0,0.82)",
          background: `linear-gradient(90deg, ${guild.color}38, rgba(0,0,0,0.78))`,
          boxShadow: `0 0 0 2px rgba(0,0,0,0.78), 7px 7px 0 rgba(0,0,0,0.42), inset 0 0 18px ${guild.color}22`,
          color: "#fff8d7",
          cursor: isJoining ? "wait" : "pointer",
          opacity: isJoining ? 0.72 : 1,
          fontFamily: "inherit",
          fontSize: "clamp(0.78rem, 1.5vw, 1rem)",
          lineHeight: 1.5,
          textShadow: `2px 2px 0 rgba(0,0,0,0.72), 0 0 14px ${guild.color}`,
        }}
      >
        {isJoining ? "JOINING..." : "JOIN GUILD"}
      </motion.button>
      {statusMessage && (
        <p
          role="alert"
          style={{
            margin: "14px 0 0",
            color: "#ffd966",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.9rem",
            lineHeight: 1.6,
          }}
        >
          {statusMessage}
        </p>
      )}
    </motion.article>
  );
}

function GuildMark({ guild, size }: { guild: GuildMaster; size: number }) {
  const isValidIcon = /^[A-Z0-9λ]+$/.test(guild.icon);
  return (
    <span
      aria-hidden="true"
      style={{
        display: "grid",
        width: `${size}px`,
        height: `${size}px`,
        placeItems: "center",
        border: `2px solid ${guild.color}`,
        background: `radial-gradient(circle at 50% 20%, ${guild.color}36, rgba(0,0,0,0.78) 62%)`,
        boxShadow: `inset 0 0 18px rgba(0,0,0,0.68), 0 0 16px ${guild.color}66`,
        color: guild.accent,
        fontSize: `${Math.max(0.76, size / 44)}rem`,
        lineHeight: 1,
        textShadow: `0 0 12px ${guild.color}`,
        overflow: "hidden",
      }}
    >
      {isValidIcon ? (
        <img
          src={`/guild-icons/${guild.icon}.png`}
          alt={`${guild.name} icon`}
          style={{
            width: "100%",
            height: "100%",
            objectFit: "cover",
            imageRendering: "pixelated",
          }}
        />
      ) : (
        <span>{guild.icon}</span>
      )}
    </span>
  );
}

function Stat({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div
      style={{
        border: "2px solid rgba(255,255,255,0.12)",
        background: "rgba(255,255,255,0.045)",
        boxShadow: "inset 0 0 14px rgba(0,0,0,0.4)",
        padding: "12px",
      }}
    >
      <dt
        style={{
          color,
          fontSize: "0.48rem",
          lineHeight: 1.35,
          marginBottom: "8px",
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
