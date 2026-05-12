import { motion } from "framer-motion";

interface ContributionItem {
  repo: string;
  type: "commit" | "pull_request";
  message: string;
  language: string;
  cp: number;
  timestamp: string;
}

interface LanguageBreakdown {
  name: string;
  cp: number;
  color: string;
  icon: string;
}

export interface AnalysisResult {
  totalCommits: number;
  totalPRs: number;
  totalCP: number;
  languageBreakdown: LanguageBreakdown[];
  contributions: ContributionItem[];
}

interface CompletePanelProps {
  result: AnalysisResult;
  onContinue: () => void;
}

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

export function CompletePanel({ result, onContinue }: CompletePanelProps) {
  return (
    <motion.div
      initial={{ scaleY: 0, opacity: 0 }}
      animate={{ scaleY: 1, opacity: 1 }}
      transition={{ duration: 0.4, ease: steppedEase(6) }}
      style={{
        position: "relative",
        zIndex: 2,
        border: "4px solid var(--color-gold)",
        background: "var(--color-navy-light)",
        boxShadow: "0 0 30px rgba(0, 245, 255, 0.12), 8px 8px 0 rgba(0,0,0,0.8)",
        maxWidth: "580px",
        width: "100%",
        padding: "2rem",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
        maxHeight: "80vh",
        overflowY: "auto",
      }}
    >
      <div style={{ textAlign: "center" }}>
        <motion.span
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: "spring", stiffness: 200, damping: 10, delay: 0.2 }}
          style={{
            display: "inline-block",
            fontFamily: "var(--font-press)",
            fontSize: "0.8rem",
            color: "var(--color-neon-green)",
            letterSpacing: "0.1em",
          }}
        >
          ★ ANALYSIS COMPLETE
        </motion.span>
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4, duration: 0.5, ease: steppedEase(6) }}
        style={{ textAlign: "center", padding: "1rem 0" }}
      >
        <div
          style={{
            fontFamily: "var(--font-press)",
            fontSize: "0.6rem",
            color: "var(--color-gold)",
            letterSpacing: "0.2em",
            marginBottom: "0.5rem",
          }}
        >
          CONTRIBUTION POINTS
        </div>
        <motion.span
          initial={{ scale: 0.5 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.6, type: "spring", stiffness: 150, damping: 8 }}
          style={{
            fontFamily: "var(--font-press)",
            fontSize: "3.5rem",
            color: "var(--color-neon-cyan)",
            textShadow: "0 0 20px rgba(0, 245, 255, 0.5), 4px 4px 0 rgba(0,0,0,0.5)",
          }}
        >
          {result.totalCP}
        </motion.span>
        <div
          style={{
            marginTop: "0.5rem",
            fontFamily: "var(--font-dot)",
            fontSize: "0.75rem",
            color: "rgba(232, 232, 208, 0.6)",
          }}
        >
          {result.totalCommits} commits / {result.totalPRs} PRs
        </div>
      </motion.div>

      <Divider />

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.7, duration: 0.4 }}
      >
        <SectionLabel text="LANGUAGE BREAKDOWN" />
        <div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "0.8rem" }}>
          {result.languageBreakdown.map((lang, i) => (
            <motion.div
              key={lang.name}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.8 + i * 0.12, duration: 0.4, ease: steppedEase(6) }}
            >
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                  marginBottom: "4px",
                }}
              >
                <span style={{ fontSize: "1rem" }}>{lang.icon}</span>
                <span
                  style={{
                    fontFamily: "var(--font-dot)",
                    fontSize: "0.8rem",
                    color: lang.color,
                    minWidth: "90px",
                  }}
                >
                  {lang.name}
                </span>
                <div
                  style={{
                    flex: 1,
                    height: "14px",
                    border: "2px solid rgba(255,255,255,0.15)",
                    background: "rgba(0,0,0,0.4)",
                    position: "relative",
                    overflow: "hidden",
                  }}
                >
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: result.totalCP === 0 ? "0%" : `${(lang.cp / result.totalCP) * 100}%` }}
                    transition={{ delay: 1 + i * 0.12, duration: 0.6, ease: steppedEase(8) }}
                    style={{
                      height: "100%",
                      background: `linear-gradient(90deg, ${lang.color}80, ${lang.color})`,
                      boxShadow: `0 0 4px ${lang.color}`,
                      position: "absolute",
                      left: 0,
                      top: 0,
                    }}
                  />
                </div>
                <span
                  style={{
                    fontFamily: "var(--font-press)",
                    fontSize: "0.65rem",
                    color: "var(--color-pixel-white)",
                    minWidth: "36px",
                    textAlign: "right",
                  }}
                >
                  +{lang.cp}
                </span>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>

      <Divider />

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.4, duration: 0.4 }}
      >
        <SectionLabel text="RECENT CONTRIBUTIONS" />
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            marginTop: "0.8rem",
            maxHeight: "200px",
            overflowY: "auto",
          }}
        >
          {result.contributions.map((item, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 1.5 + i * 0.05, duration: 0.3 }}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "6px 4px",
                borderBottom: "1px solid rgba(255,255,255,0.05)",
                fontSize: "0.75rem",
              }}
            >
              <span
                style={{
                  fontFamily: "var(--font-press)",
                  fontSize: "0.6rem",
                  color: item.type === "pull_request" ? "var(--color-neon-green)" : "var(--color-gold)",
                  minWidth: "24px",
                  textAlign: "center",
                }}
              >
                +{item.cp}
              </span>
              <span style={{ color: "rgba(232,232,208,0.4)", fontSize: "0.65rem", minWidth: "120px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                {item.repo}
              </span>
              <span style={{ color: "var(--color-pixel-white)", flex: 1, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                {item.message}
              </span>
              <span style={{ color: "rgba(232,232,208,0.3)", fontSize: "0.6rem", minWidth: "50px", textAlign: "right" }}>
                {item.timestamp}
              </span>
            </motion.div>
          ))}
        </div>
      </motion.div>

      <motion.button
        whileHover={{ scale: 1.02, filter: "brightness(1.1)" }}
        whileTap={{ scale: 0.98, y: 4, boxShadow: "0px 0px 0 var(--color-gold-dark)" }}
        onClick={onContinue}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 2, duration: 0.4 }}
        style={{
          marginTop: "0.5rem",
          width: "100%",
          padding: "1rem",
          fontSize: "1rem",
          fontFamily: "var(--font-press)",
          background: "var(--color-gold)",
          color: "#000",
          border: "none",
          boxShadow: "0px 4px 0 var(--color-gold-dark)",
          cursor: "pointer",
          letterSpacing: "0.05em",
        }}
      >
        CONTINUE
      </motion.button>
    </motion.div>
  );
}

function Divider() {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: "0.5rem",
      }}
    >
      <div style={{ flex: 1, height: "1px", background: "rgba(255,215,0,0.2)" }} />
      <span style={{ color: "rgba(255,215,0,0.3)", fontSize: "0.5rem" }}>◆</span>
      <div style={{ flex: 1, height: "1px", background: "rgba(255,215,0,0.2)" }} />
    </div>
  );
}

function SectionLabel({ text }: { text: string }) {
  return (
    <div
      style={{
        fontFamily: "var(--font-dot)",
        fontSize: "0.7rem",
        color: "var(--color-gold)",
        letterSpacing: "0.15em",
        padding: "4px 0",
      }}
    >
      ▶ {text}
    </div>
  );
}
