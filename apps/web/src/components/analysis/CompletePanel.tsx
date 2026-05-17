import { motion } from "framer-motion";

interface ContributionItem {
  repo: string;
  type?: "commit" | "pull_request";
  message: string;
  language: string;
  cp: number;
  timestamp: string;
}

interface LanguageBreakdown {
  name: string;
  cp: number;
  color?: string;
  icon?: string;
}

export interface AnalysisResult {
  totalCommits: number;
  totalPRs: number;
  totalCP: number;
  totalBalance: number;
  languageBreakdown: LanguageBreakdown[];
  contributions: ContributionItem[];
}

interface CompletePanelProps {
  result: AnalysisResult;
  onContinue: () => void;
}

const LANGUAGE_META: Record<string, { color: string; icon: string }> = {
  TypeScript: { color: "#3178c6", icon: "📘" },
  JavaScript: { color: "#f7df1e", icon: "📒" },
  Rust: { color: "#dea584", icon: "🦀" },
  Go: { color: "#00add8", icon: "🐹" },
  Python: { color: "#3572a5", icon: "🐍" },
  Ruby: { color: "#701516", icon: "💎" },
  Java: { color: "#b07219", icon: "☕" },
  Kotlin: { color: "#a97bff", icon: "🅺" },
  Swift: { color: "#f05138", icon: "🍎" },
  "C++": { color: "#f34b7d", icon: "⚙️" },
  C: { color: "#555555", icon: "⚙️" },
  "C#": { color: "#178600", icon: "🎯" },
  PHP: { color: "#4f5d95", icon: "🐘" },
  Shell: { color: "#89e051", icon: "🐚" },
  Dockerfile: { color: "#384d54", icon: "🐳" },
  HTML: { color: "#e34c26", icon: "🌐" },
  CSS: { color: "#563d7c", icon: "🎨" },
  Scala: { color: "#c22d40", icon: "🔥" },
  Dart: { color: "#00b4ab", icon: "🎯" },
  Lua: { color: "#000080", icon: "🌙" },
  Haskell: { color: "#5e5086", icon: "λ" },
};

function langMeta(name: string): { color: string; icon: string } {
  return LANGUAGE_META[name] ?? { color: "#8b8b8b", icon: "◇" };
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

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.8, duration: 0.5, ease: steppedEase(6) }}
        style={{ textAlign: "center", padding: "0.5rem 0" }}
      >
        <div
          style={{
            fontFamily: "var(--font-press)",
            fontSize: "0.5rem",
            color: "var(--color-pixel-white)",
            letterSpacing: "0.2em",
            marginBottom: "0.3rem",
          }}
        >
          TOTAL CP
        </div>
        <motion.span
          initial={{ scale: 0.5 }}
          animate={{ scale: 1 }}
          transition={{ delay: 1, type: "spring", stiffness: 150, damping: 8 }}
          style={{
            fontFamily: "var(--font-press)",
            fontSize: "2rem",
            color: "var(--color-gold)",
            textShadow: "0 0 20px rgba(255, 215, 0, 0.4), 4px 4px 0 rgba(0,0,0,0.5)",
          }}
        >
          {result.totalBalance}
        </motion.span>
      </motion.div>

      <Divider />

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.1, duration: 0.4 }}
      >
        <SectionLabel text="LANGUAGE BREAKDOWN" />
        <div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "0.8rem" }}>
          {result.languageBreakdown.map((lang, i) => {
            const meta = langMeta(lang.name);
            const color = lang.color ?? meta.color;
            const icon = lang.icon ?? meta.icon;
            return (
              <motion.div
                key={lang.name}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 1.2 + i * 0.12, duration: 0.4, ease: steppedEase(6) }}
              >
                <div
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "8px",
                    marginBottom: "4px",
                  }}
                >
                  <span style={{ fontSize: "1rem" }}>{icon}</span>
                  <span
                    style={{
                      fontFamily: "var(--font-dot)",
                      fontSize: "0.8rem",
                      color: color,
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
                      animate={{
                        width: result.totalCP === 0 ? "0%" : `${(lang.cp / result.totalCP) * 100}%`,
                      }}
                      transition={{ delay: 1.2 + i * 0.12, duration: 0.6, ease: steppedEase(8) }}
                      style={{
                        height: "100%",
                        background: `linear-gradient(90deg, ${color}80, ${color})`,
                        boxShadow: `0 0 4px ${color}`,
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
            );
          })}
        </div>
      </motion.div>

      <Divider />

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.8, duration: 0.4 }}
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
          {result.contributions.map((item, i) => {
            const type = item.type ?? "commit";
            const ts = item.timestamp.length > 10 ? item.timestamp.slice(0, 10) : item.timestamp;
            return (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 1.9 + i * 0.05, duration: 0.3 }}
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
                    color:
                      type === "pull_request" ? "var(--color-neon-green)" : "var(--color-gold)",
                    minWidth: "24px",
                    textAlign: "center",
                  }}
                >
                  +{item.cp}
                </span>
                <span
                  style={{
                    color: "rgba(232,232,208,0.4)",
                    fontSize: "0.65rem",
                    minWidth: "120px",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                >
                  {item.repo}
                </span>
                <span
                  style={{
                    color: "var(--color-pixel-white)",
                    flex: 1,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                >
                  {item.message}
                </span>
                <span
                  style={{
                    color: "rgba(232,232,208,0.3)",
                    fontSize: "0.6rem",
                    minWidth: "50px",
                    textAlign: "right",
                  }}
                >
                  {ts}
                </span>
              </motion.div>
            );
          })}
        </div>
      </motion.div>

      <motion.button
        whileHover={{ scale: 1.02, filter: "brightness(1.1)" }}
        whileTap={{ scale: 0.98, y: 4, boxShadow: "0px 0px 0 var(--color-gold-dark)" }}
        onClick={onContinue}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 2.5, duration: 0.4 }}
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
