import { motion } from "framer-motion";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { steppedEase } from "../../lib/animationUtils";
import styles from "./MyGuildDetails.module.css";

interface MyGuildDetailsProps {
  onNavigate: (path: string) => void;
}

interface GuildDetails {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  level: number;
  totalCp: number;
  memberCount: number;
  memberLimit: number;
}

interface GuildMember {
  id: string;
  name: string;
  title: string;
  avatar: string;
  contributionCp: number;
  isCurrentUser?: boolean;
}

const MY_GUILD: GuildDetails = {
  id: "guild_typescript",
  slug: "typescript",
  name: "TypeScript Guild",
  description:
    "型の力でフロントからバックエンドまで支える、適応型のギルド。作戦会議室では、今日の貢献と次の一手がネオンログとして共有されます。",
  icon: "TS",
  color: "#3178c6",
  level: 17,
  totalCp: 125400,
  memberCount: 4,
  memberLimit: 10,
};

const MEMBERS: GuildMember[] = [
  {
    id: "current-user",
    name: "TypeSmith",
    title: "Generic Hero",
    avatar: "YOU",
    contributionCp: 35420,
    isCurrentUser: true,
  },
  {
    id: "null-mage",
    name: "NullMage",
    title: "Void Debugger",
    avatar: "NM",
    contributionCp: 31980,
  },
  {
    id: "pixel-ninja",
    name: "PixelNinja",
    title: "UI Shinobi",
    avatar: "PN",
    contributionCp: 28640,
  },
  {
    id: "loop-knight",
    name: "LoopKnight",
    title: "Iteration Paladin",
    avatar: "LK",
    contributionCp: 25110,
  },
];

const panelTransition = { duration: 0.34, ease: steppedEase(6) };

export function MyGuildDetails({ onNavigate }: MyGuildDetailsProps) {
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);
  const stats = [
    { label: "GUILD LEVEL", value: `LV.${MY_GUILD.level}` },
    { label: "TOTAL CP", value: MY_GUILD.totalCp.toLocaleString() },
    { label: "MEMBERS", value: `${MY_GUILD.memberCount}/${MY_GUILD.memberLimit}` },
  ];

  return (
    <main className={styles.main}>
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />

      <div className={styles.shell}>
        <div className={styles.topBar}>
          <motion.button
            className={styles.backButton}
            type="button"
            whileHover={{ y: -2, scale: 1.02 }}
            whileTap={{ y: 2, scale: 0.98 }}
            onClick={() => void navigateBackWithSe("/guild")}
          >
            &lt; GUILD BASE
          </motion.button>
          <p className={styles.routeLabel}>MY GUILD / OPERATIONS ROOM</p>
        </div>

        <div className={styles.layout}>
          <motion.section
            className={`${styles.panel} ${styles.guildPanel}`}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={panelTransition}
            aria-labelledby="my-guild-title"
          >
            <div className={styles.identity}>
              <div className={styles.emblem} aria-hidden="true" style={{ color: MY_GUILD.color }}>
                {MY_GUILD.icon}
              </div>
              <div>
                <span className={styles.eyebrow}>CURRENT GUILD / {MY_GUILD.slug}</span>
                <h1 className={styles.title} id="my-guild-title">
                  {MY_GUILD.name}
                </h1>
              </div>
            </div>

            <p className={styles.description}>{MY_GUILD.description}</p>

            <div className={styles.stats} aria-label="Guild stats">
              {stats.map((stat) => (
                <div className={styles.statCard} key={stat.label}>
                  <span className={styles.statLabel}>{stat.label}</span>
                  <strong className={styles.statValue}>{stat.value}</strong>
                </div>
              ))}
            </div>

            <div className={styles.actions}>
              <motion.button
                className={styles.copyButton}
                type="button"
                whileHover={{ y: -2, scale: 1.02 }}
                whileTap={{ y: 2, scale: 0.98 }}
              >
                INVITE LINK
              </motion.button>
              <motion.button
                className={styles.leaveButton}
                type="button"
                whileHover={{ y: -2, scale: 1.02 }}
                whileTap={{ y: 2, scale: 0.98 }}
              >
                LEAVE GUILD
              </motion.button>
            </div>
          </motion.section>

          <motion.section
            className={`${styles.panel} ${styles.memberPanel}`}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ ...panelTransition, delay: 0.08 }}
            aria-labelledby="member-list-title"
          >
            <header className={styles.panelHeader}>
              <h2 className={styles.panelTitle} id="member-list-title">
                MEMBER CONTRIBUTION
              </h2>
              <span className={styles.online}>{MEMBERS.length} OPERATORS ONLINE</span>
            </header>

            <div className={styles.memberList}>
              {MEMBERS.map((member, index) => (
                <motion.article
                  className={`${styles.memberRow} ${
                    member.isCurrentUser ? styles.memberRowSelf : ""
                  }`}
                  key={member.id}
                  initial={{ opacity: 0, y: 14 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.12 + index * 0.06, duration: 0.28, ease: steppedEase(5) }}
                  aria-label={member.isCurrentUser ? `${member.name}, current user` : member.name}
                >
                  <div className={styles.avatar} aria-hidden="true">
                    {member.avatar}
                  </div>
                  <div>
                    <strong className={styles.memberName}>{member.name}</strong>
                    <span className={styles.memberTitle}>{member.title}</span>
                  </div>
                  <strong className={styles.memberCp}>
                    {member.contributionCp.toLocaleString()} CP
                  </strong>
                </motion.article>
              ))}
            </div>
          </motion.section>
        </div>
      </div>

      <div className={styles.scanline} aria-hidden="true" />
    </main>
  );
}
