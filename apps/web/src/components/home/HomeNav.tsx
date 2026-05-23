import { motion, type Variants } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";
import styles from "./Home.module.css";

interface HomeNavItem {
  label: string;
  caption: string;
  path: string;
  accent: string;
}

const navVariants: Variants = {
  hidden: { opacity: 0, y: 18 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: 0.15 + i * 0.08, duration: 0.35, ease: steppedEase(6) },
  }),
};

export function HomeNav({
  items,
  onNavigate,
}: {
  items: HomeNavItem[];
  onNavigate: (path: string) => void;
}) {
  return (
    <nav
      aria-label="Main navigation"
      className={styles.nav}
    >
      {items.map((item, index) => (
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
            border: `3px solid ${item.accent}`,
            borderBottomColor: "rgba(0,0,0,0.75)",
            borderRightColor: "rgba(0,0,0,0.75)",
            background: "rgba(4, 11, 27, 0.82)",
            boxShadow: `0 0 0 2px rgba(0,0,0,0.78), 6px 6px 0 rgba(0,0,0,0.5), inset 0 0 18px ${item.accent}22`,
            color: "#fff8d7",
            cursor: "pointer",
            fontFamily: "inherit",
            textAlign: "left",
          }}
          className={styles.navButton}
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
  );
}
