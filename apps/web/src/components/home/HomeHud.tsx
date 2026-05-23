import { motion, type Variants } from "framer-motion";
import type { ReactNode } from "react";
import { steppedEase } from "../../lib/animationUtils";
import { AudioTogglePanel } from "../shared/AudioTogglePanel";
import styles from "./Home.module.css";

interface PlayerSummary {
  name: string;
  title: string;
  level: number;
  levelTotalCp: number;
  totalCp: number;
  todayCp: number;
  nextLevel: number;
  nextLevelTotalCp: number;
  nextLevelRemainingCp: number;
  lifetimeTotalEarnedCp: number;
}

export interface GuildSummary {
  name: string;
  icon: string;
  rank: string;
  accent: string;
}

const panelVariants: Variants = {
  hidden: { opacity: 0, y: -14 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.35, ease: steppedEase(6) },
  },
};

function HudPanel({
  align = "left",
  children,
  className = "",
}: {
  align?: "left" | "right";
  children: ReactNode;
  className?: string;
}) {
  return (
    <motion.section
      className={`${styles.hudPanel} ${className}`}
      variants={panelVariants}
      initial="hidden"
      animate="visible"
      style={{ textAlign: align }}
    >
      {children}
    </motion.section>
  );
}

function LabelValue({
  label,
  value,
}: {
  label: string;
  value: string | number;
}) {
  return (
    <div className={styles.labelRow}>
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

function ExpBar({ player }: { player: PlayerSummary }) {
  const levelRange = Math.max(1, player.nextLevelTotalCp - player.levelTotalCp);
  const earnedInLevel = Math.max(
    0,
    player.lifetimeTotalEarnedCp - player.levelTotalCp,
  );
  const pct = Math.min(100, Math.max(0, (earnedInLevel / levelRange) * 100));

  return (
    <div className={styles.expBar}>
      <div
        aria-label={`EXP ${earnedInLevel.toLocaleString()} / ${levelRange.toLocaleString()} CP`}
        role="progressbar"
        aria-valuemin={0}
        aria-valuemax={levelRange}
        aria-valuenow={earnedInLevel}
        style={{
          height: "12px",
          border: "2px solid rgba(255, 215, 0, 0.5)",
          borderBottomColor: "rgba(111, 79, 28, 0.95)",
          borderRightColor: "rgba(111, 79, 28, 0.95)",
          background: "rgba(0,0,0,0.45)",
          boxShadow: "inset 0 0 0 1px rgba(255,255,255,0.08)",
          overflow: "hidden",
        }}
      >
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${pct}%` }}
          transition={{ duration: 0.45, ease: steppedEase(8) }}
          style={{
            height: "100%",
            background: "linear-gradient(90deg, #00f5ff, #ffd700)",
            boxShadow: "0 0 10px rgba(0, 245, 255, 0.75)",
          }}
        />
      </div>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          gap: "8px",
          color: "rgba(255, 248, 215, 0.72)",
          fontSize: "0.54rem",
          lineHeight: 1.5,
        }}
      >
        <span>{earnedInLevel.toLocaleString()} CP</span>
        <span>{levelRange.toLocaleString()} CP</span>
      </div>
    </div>
  );
}

function TitleButton({ onClick, className = "" }: { onClick: () => void; className?: string }) {
  return (
    <motion.button
      className={className}
      type="button"
      variants={panelVariants}
      initial="hidden"
      animate="visible"
      whileHover={{ y: -2, scale: 1.03 }}
      whileTap={{ y: 2, scale: 0.98 }}
      onClick={onClick}
      style={{
        alignSelf: "flex-start",
        border: "2px solid rgba(244, 236, 208, 0.55)",
        borderBottomColor: "rgba(0,0,0,0.78)",
        borderRightColor: "rgba(0,0,0,0.78)",
        background: "rgba(3, 10, 24, 0.68)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 5px 5px 0 rgba(0,0,0,0.38)",
        color: "#fff8d7",
        cursor: "pointer",
        fontFamily: "inherit",
        fontSize: "clamp(0.56rem, 1.4vw, 0.74rem)",
        lineHeight: 1.5,
        padding: "10px 12px",
      }}
    >
      &lt; TITLE
    </motion.button>
  );
}

function GuildEmblem({ accent, icon }: { accent: string; icon: string }) {
  const isValidIcon = /^[A-Z0-9λ]+$/.test(icon);
  return (
    <div
      aria-hidden="true"
      style={{
        width: "48px",
        height: "48px",
        display: "grid",
        placeItems: "center",
        flex: "0 0 auto",
        border: `3px solid ${accent}`,
        borderBottomColor: "#035a72",
        borderRightColor: "#035a72",
        background:
          "linear-gradient(135deg, rgba(49,120,198,0.34), rgba(255,217,102,0.16)), #061326",
        boxShadow: `0 0 0 2px rgba(0,0,0,0.78), inset 0 0 16px ${accent}44, 3px 3px 0 rgba(0,0,0,0.42)`,
        color: "#ffd966",
        fontSize: "1rem",
        lineHeight: 1,
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        overflow: "hidden",
      }}
    >
      {isValidIcon ? (
        <img
          src={`/guild-icons/${icon}.png`}
          alt="Guild icon"
          style={{
            width: "100%",
            height: "100%",
            objectFit: "cover",
            imageRendering: "pixelated",
          }}
        />
      ) : (
        <span>{icon}</span>
      )}
    </div>
  );
}

function GuildMembership({
  guild,
  className = "",
}: {
  guild: GuildSummary | null | undefined;
  className?: string;
}) {
  const accent = guild?.accent ?? "#f4ecd0";
  const icon = guild?.icon ?? "--";
  const name = guild === undefined ? "確認中" : (guild?.name ?? "未所属");
  const rank =
    guild === undefined
      ? "Guild status loading..."
      : (guild?.rank ?? "Guild Base から所属先を選択");

  return (
    <HudPanel className={className}>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "12px",
        }}
      >
        <GuildEmblem accent={accent} icon={icon} />
        <div style={{ minWidth: 0 }}>
          <div
            style={{
              color: "rgba(244, 236, 208, 0.62)",
              fontSize: "0.58rem",
              lineHeight: 1.5,
              marginBottom: "5px",
            }}
          >
            GUILD
          </div>
          <div
            style={{
              color: accent,
              fontSize: "clamp(0.74rem, 1.7vw, 0.92rem)",
              lineHeight: 1.45,
              overflowWrap: "anywhere",
            }}
          >
            {name}
          </div>
          <div
            style={{
              color: "rgba(255, 248, 215, 0.72)",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.72rem",
              lineHeight: 1.5,
              marginTop: "2px",
            }}
          >
            {rank}
          </div>
        </div>
      </div>
    </HudPanel>
  );
}

export function HomeHud({
  guild,
  onReturnTitle,
  player,
}: {
  guild: GuildSummary | null | undefined;
  onReturnTitle: () => void;
  player: PlayerSummary;
}) {
  return (
    <header className={styles.hudHeader}>
      <div className={styles.hudLeft}>
        <HudPanel className={styles.playerInfoPanel}>
          <div
            className={styles.sectionHeader}
            style={{ color: "#ffd700" }}
          >
            PLAYER INFO
          </div>
          <div className={styles.infoGrid}>
            <LabelValue label="NAME" value={player.name} />
            <LabelValue label="TITLE" value={player.title} />
            <LabelValue label="LEVEL" value={`LV.${player.level}`} />
            <LabelValue
              label="NEXT"
              value={`${player.nextLevelRemainingCp.toLocaleString()} CP → LV.${player.nextLevel}`}
            />
            <ExpBar player={player} />
          </div>
        </HudPanel>

        <GuildMembership guild={guild} className={styles.guildPanel} />

        <div
          className={styles.titleBtn}
          style={{
            display: "flex",
            gap: "12px",
            alignItems: "flex-start",
          }}
        >
          <AudioTogglePanel position="bottom-left" inlineOnMobile={true} />
          <TitleButton onClick={onReturnTitle} />
        </div>
      </div>

      <HudPanel align="right" className={styles.contribPanel}>
        <div
          className={styles.sectionHeader}
          style={{ color: "#00f5ff" }}
        >
          CONTRIBUTION POINT
        </div>
        <div className={styles.infoGrid}>
          <LabelValue
            label="TOTAL CP"
            value={player.totalCp.toLocaleString()}
          />
          <LabelValue
            label="TODAY CP"
            value={`+${player.todayCp.toLocaleString()}`}
          />
          <LabelValue
            label="EXP"
            value={`${player.lifetimeTotalEarnedCp.toLocaleString()} / ${player.nextLevelTotalCp.toLocaleString()}`}
          />
        </div>
      </HudPanel>
    </header>
  );
}
