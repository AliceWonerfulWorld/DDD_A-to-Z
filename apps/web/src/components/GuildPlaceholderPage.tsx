import { motion } from "framer-motion";
import { steppedEase } from "../lib/animationUtils";
import { GuildBgm } from "./GuildBgm";
import styles from "./GuildPlaceholderPage.module.css";

interface GuildPlaceholderPageProps {
  title: string;
  caption: string;
  backPath?: string;
  onNavigate: (path: string) => void;
}

export function GuildPlaceholderPage({
  title,
  caption,
  backPath = "/guild",
  onNavigate,
}: GuildPlaceholderPageProps) {
  return (
    <main className={styles.mainContainer}>
      <GuildBgm />

      <motion.section
        className={styles.panel}
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.34, ease: steppedEase(6) }}
      >
        <p className={styles.routeLabel}>ROUTE READY</p>
        <h1 className={styles.title}>{title}</h1>
        <p className={styles.caption}>{caption}</p>
        <motion.button
          className={styles.actionButton}
          type="button"
          whileHover={{ y: -2, scale: 1.02 }}
          whileTap={{ y: 2, scale: 0.98 }}
          onClick={() => onNavigate(backPath)}
        >
          &lt; GUILD BASE
        </motion.button>
      </motion.section>
    </main>
  );
}
