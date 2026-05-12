import { useState, useMemo, type ReactNode } from "react";
import { motion } from "framer-motion";

interface MyPageProps {
  onNavigate: (path: string) => void;
}

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

/* ─── Mock Data ─── */

const MOCK = {
  user: { name: "DevSamurai", title: "Apprentice", rankName: "Code Seeker" },
  season: {
    label: "SEASON 1",
    start: "2024/05/01",
    end: "2024/07/31",
    remaining: 52,
  },
  guild: {
    name: "TypeScript",
    icon: "📘",
    color: "#3178c6",
    fullName: "TypeScript GUILD",
    desc: "型の力で安全で堅牢なコードを書く、\nエレガントな戦士たちの集い。",
    rank: 42,
    total: 156,
    cp: 24680,
  },
  langs: [
    { name: "TypeScript", pct: 42, icon: "📘", color: "#3178c6", atk: 120, def: 95, mag: 130, spd: 85, luk: 70 },
    { name: "Rust", pct: 26, icon: "🦀", color: "#ff6b35", atk: 145, def: 135, mag: 60, spd: 55, luk: 50 },
    { name: "Python", pct: 18, icon: "🐍", color: "#f0c040", atk: 70, def: 60, mag: 140, spd: 80, luk: 90 },
    { name: "Go", pct: 8, icon: "🐹", color: "#00acd7", atk: 90, def: 85, mag: 45, spd: 110, luk: 75 },
    { name: "Other", pct: 6, icon: "🔮", color: "#888", atk: 55, def: 50, mag: 60, spd: 65, luk: 40 },
  ],
  goal: { current: 2150, target: 3300 },
  title: { name: "Consistency Master", line: "Consistency is key. Daily efforts build the future." },
  badges: [
    { id: 1, icon: "🌱", name: "はじめの一歩", desc: "100 Commit達成", color: "#4caf50" },
    { id: 2, icon: "🌿", name: "草生やし名人", desc: "10日連続コミット", color: "#4caf50" },
    { id: 3, icon: "🔄", name: "プルリク職人", desc: "50 PR作成", color: "#00e5ff" },
    { id: 4, icon: "🔀", name: "マージマスター", desc: "100 Merged", color: "#9c27b0" },
    { id: 5, icon: "🌟", name: "コミット王", desc: "500 Commit達成", color: "#f0c040", star: "500" },
    { id: 6, icon: "🔥", name: "ファイヤー", desc: "30日連続コミット", color: "#ff6b35" },
    { id: 7, icon: "🌙", name: "ナイトオウル", desc: "深夜の貢献者", color: "#1a237e" },
    { id: 8, icon: "💎", name: "レジェンド", desc: "1,000 Commit達成", color: "#7b1fa2" },
  ],
  messages: [
    "よく来たな、若きエンジニアよ。\n今日も己のコードを刻め。",
    "知識という名の剣は、使えば使うほど\n研ぎ澄まされていくものだ。",
    "眠れぬ夜の一振りが、やがて\n城を築く礎となる。",
    "焦る必要はない。確かな一歩を\n積み重ねるのだ。",
    "コミットの数だけ、お前は強くなる。\nさあ、今日も戦場へ赴こう。",
  ],
  todayLog: [
    "Fix: resolve type error in UserService",
    "Review: PR #42 — feat: add cache layer",
    "Add: unit tests for auth middleware",
    "Docs: update API endpoint comments",
    "Chore: bump dependencies",
  ],
};

function generateContributions() {
  const weeks = 52;
  const days = 7;
  const data: number[][] = [];
  for (let w = 0; w < weeks; w++) {
    const row: number[] = [];
    for (let d = 0; d < days; d++) {
      const r = Math.random();
      row.push(r < 0.3 ? 0 : r < 0.55 ? 1 : r < 0.75 ? 2 : r < 0.88 ? 3 : r < 0.95 ? 4 : Math.floor(Math.random() * 6) + 5);
    }
    data.push(row);
  }
  return data;
}

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
const DAYS_LABELS = ["Mon", "", "Wed", "", "Fri", "", ""];
const CELL = 14;
const GAP = 3;

/* ─── Sub-components ─── */

function SectionTitle({ text, color }: { text: string; color?: string }) {
  return (
    <div
      style={{
        fontFamily: '"Press Start 2P", monospace',
        fontSize: "0.9rem",
        color: color ?? "var(--color-gold)",
        letterSpacing: "0.08em",
        padding: "4px 0",
        borderBottom: "1px solid rgba(255,215,0,0.12)",
        marginBottom: "10px",
      }}
    >
      ▸ {text}
    </div>
  );
}

function Panel({ children, borderColor }: { children: ReactNode; borderColor?: string }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, ease: steppedEase(6) }}
      style={{
        border: `2px solid ${borderColor ?? "rgba(255,255,255,0.08)"}`,
        background: "rgba(0,0,0,0.35)",
        padding: "14px",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        overflow: "auto",
      }}
    >
      {children}
    </motion.div>
  );
}

function ProgressBarFill({ pct, color, delay = 0.3 }: { pct: number; color: string; delay?: number }) {
  return (
    <div style={{ height: "100%", position: "relative", overflow: "hidden" }}>
      <motion.div
        initial={{ width: 0 }}
        animate={{ width: `${Math.min(pct, 100)}%` }}
        transition={{ duration: 0.8, delay, ease: steppedEase(8) }}
        style={{
          height: "100%",
          background: `linear-gradient(90deg, ${color}80, ${color})`,
          boxShadow: `0 0 6px ${color}`,
          position: "absolute",
          left: 0,
          top: 0,
        }}
      />
    </div>
  );
}

/* ─── Main Component ─── */

export function MyPage({ onNavigate }: MyPageProps) {
  const contribs = useMemo(() => generateContributions(), []);
  const [tooltip, setTooltip] = useState<{ x: number; y: number; date: string; count: number } | null>(null);

  const streak = useMemo(() => {
    let max = 0, cur = 0;
    for (const row of contribs) {
      for (const c of row) {
        if (c > 0) { cur++; max = Math.max(max, cur); } else cur = 0;
      }
    }
    return max;
  }, [contribs]);

  const monthLabels = useMemo(() => {
    const labels: { week: number; label: string }[] = [];
    for (let w = 0; w < 52; w++) {
      const d = new Date(); d.setDate(d.getDate() - (52 - w) * 7);
      const m = d.getMonth();
      if (!labels.length || labels[labels.length - 1].label !== MONTHS[m]) {
        labels.push({ week: w, label: MONTHS[m] });
      }
    }
    return labels;
  }, []);

  const guild = MOCK.guild;
  const gColor = guild.color;

  return (
    <div
      className="flex flex-col min-h-svh relative overflow-hidden"
      style={{
        background: "radial-gradient(ellipse at 50% 0%, #0d1b2a 0%, #0a0a1a 50%, #050510 100%)",
        fontFamily: '"Press Start 2P", monospace',
        color: "#e8e8d0",
      }}
    >
      {/* City silhouette */}
      <div
        aria-hidden="true"
        style={{
          position: "fixed", inset: 0,
          background: "url('/pixel-town-night.png') center bottom / cover no-repeat",
          opacity: 0.1, pointerEvents: "none", zIndex: 0,
        }}
      />
      {/* Scanline */}
      <div
        aria-hidden="true"
        style={{
          position: "fixed", inset: 0,
          backgroundImage: "repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0,0,0,0.04) 2px, rgba(0,0,0,0.04) 4px)",
          pointerEvents: "none", zIndex: 1,
        }}
      />

      {/* ─── Header ─── */}
      <motion.header
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.4, ease: steppedEase(6) }}
        style={{
          position: "relative", zIndex: 3,
          display: "flex", alignItems: "center", justifyContent: "space-between",
          padding: "10px 24px",
          borderBottom: "2px solid rgba(240,192,64,0.3)",
          background: "rgba(0,0,0,0.6)",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <span style={{ fontSize: "2rem" }}>👑</span>
          <span style={{ fontSize: "1rem", color: "#f0c040", letterSpacing: "0.1em" }}>MY PAGE</span>
          <span style={{ fontSize: "1rem", color: "rgba(240,192,64,0.3)" }}>{">"}</span>
          <span style={{ fontSize: "0.9rem", color: "rgba(232,232,208,0.5)" }}>ENGINEER STATUS</span>
        </div>
        <button
          onClick={() => onNavigate("/")}
          style={{
            fontFamily: '"Press Start 2P", monospace',
            fontSize: "0.8rem", color: "rgba(232,232,208,0.4)",
            background: "none", border: "1px solid rgba(255,255,255,0.12)",
            padding: "6px 12px", cursor: "pointer",
          }}
        >
          LOGOUT [→]
        </button>
      </motion.header>

      {/* ─── Main Content ─── */}
      <div
        style={{
          position: "relative", zIndex: 2,
          flex: 1,
          display: "flex", flexDirection: "column", gap: "14px",
          padding: "16px 20px",
          maxWidth: "1280px", width: "100%", margin: "0 auto",
          overflow: "hidden",
        }}
      >
        {/* ═══ Top Row: 3 columns ═══ */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: "14px", flex: 1, minHeight: 0 }}>
          {/* Left: Adventurer Profile */}
          <Panel borderColor="rgba(240,192,64,0.3)">
            <SectionTitle text="PROFILE" color="#f0c040" />
            <div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontFamily: '"Press Start 2P", monospace', fontSize: "1rem", color: "#e8e8d0" }}>{MOCK.user.name}</div>
                <div style={{ fontSize: "0.8rem", color: "rgba(232,232,208,0.4)", marginTop: "4px", fontFamily: '"Press Start 2P", monospace' }}>
                  {MOCK.user.title}
                </div>
                <div style={{ marginTop: "10px", border: "1px solid rgba(240,192,64,0.15)", background: "rgba(240,192,64,0.04)", padding: "8px" }}>
                  <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.3)", fontFamily: '"Press Start 2P", monospace' }}>TITLE</div>
                  <div style={{ display: "flex", alignItems: "center", gap: "4px", marginTop: "4px" }}>
                    <span style={{ fontSize: "1.4rem" }}>👑</span>
                    <span style={{ fontSize: "0.8rem", color: "#f0c040", fontFamily: '"Press Start 2P", monospace' }}>{MOCK.user.rankName}</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Season info */}
            <div style={{ marginTop: "12px", padding: "10px", border: "1px solid rgba(156,39,176,0.3)", background: "rgba(156,39,176,0.06)" }}>
              <div style={{ display: "flex", alignItems: "center", gap: "6px", marginBottom: "6px" }}>
                <span style={{ fontSize: "0.7rem", background: "rgba(156,39,176,0.5)", color: "#e8e8d0", padding: "2px 6px", fontFamily: '"Press Start 2P", monospace' }}>
                  {MOCK.season.label}
                </span>
              </div>
              <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.4)", fontFamily: '"Press Start 2P", monospace' }}>
                {MOCK.season.start} 〜 {MOCK.season.end}
              </div>
              <div style={{ marginTop: "6px", display: "flex", alignItems: "baseline", gap: "6px" }}>
                <span style={{ fontSize: "2rem", color: "#f0c040", fontFamily: '"Press Start 2P", monospace' }}>{MOCK.season.remaining}</span>
                <span style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.4)", fontFamily: '"Press Start 2P", monospace', marginLeft: "8px" }}>DAYS LEFT</span>
              </div>
            </div>
          </Panel>

          {/* Center: Guild */}
          <Panel borderColor={`${gColor}40`}>
            <SectionTitle text="GUILD" color={gColor} />
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: "8px" }}>
              {/* Emblem */}
              <motion.div
                animate={{ scaleY: [1, 1.03, 1] }}
                transition={{ duration: 2, repeat: Infinity, ease: steppedEase(4) }}
                style={{
                  width: "100px", height: "110px",
                  border: `2px solid ${gColor}`,
                  background: `${gColor}08`,
                  display: "flex", flexDirection: "column",
                  alignItems: "center", justifyContent: "center",
                  gap: "2px",
                  position: "relative",
                }}
              >
                <span style={{ fontSize: "3.6rem" }}>📖</span>
                <div style={{ fontSize: "0.6rem", color: gColor, fontFamily: '"Press Start 2P", monospace', letterSpacing: "0.2em" }}>VOL.3</div>
                {/* Laurel decoration */}
                <span style={{ position: "absolute", top: "-4px", fontSize: "1.2rem", color: "#f0c040" }}>🏅</span>
              </motion.div>

              <div style={{ fontFamily: '"Press Start 2P", monospace', fontSize: "0.9rem", color: "#00e5ff", letterSpacing: "0.1em", textAlign: "center", lineHeight: 1.4 }}>
                {MOCK.guild.fullName}
              </div>

              <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                <span style={{ fontSize: "0.7rem", background: "rgba(156,39,176,0.5)", color: "#e8e8d0", padding: "2px 6px", fontFamily: '"Press Start 2P", monospace' }}>
                  🏷 MEMBER
                </span>
                <span style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.3)", fontFamily: '"Press Start 2P", monospace' }}>since Season 1</span>
              </div>

              <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.4)", textAlign: "center", lineHeight: 1.6, fontFamily: '"Press Start 2P", monospace' }}>
                {MOCK.guild.desc}
              </div>
            </div>

            {/* Guild stats */}
            <div style={{ display: "flex", flexDirection: "column", gap: "14px", marginTop: "12px", padding: "12px 8px", borderTop: "1px solid rgba(255,255,255,0.06)", background: "rgba(0,0,0,0.2)" }}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-end" }}>
                <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.4)", fontFamily: '"Press Start 2P", monospace' }}>RANK</div>
                <div style={{ textAlign: "right" }}>
                  <span style={{ fontSize: "1.4rem", color: "#f0c040", fontFamily: '"Press Start 2P", monospace' }}>#{MOCK.guild.rank}</span>
                  <span style={{ fontSize: "0.6rem", color: "rgba(232,232,208,0.3)", fontFamily: '"Press Start 2P", monospace', marginLeft: "8px" }}>/ {MOCK.guild.total} GUILDS</span>
                </div>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-end" }}>
                <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.4)", fontFamily: '"Press Start 2P", monospace' }}>GUILD CP</div>
                <div style={{ textAlign: "right" }}>
                  <span style={{ fontSize: "1.2rem", color: "#00e5ff", fontFamily: '"Press Start 2P", monospace' }}>{MOCK.guild.cp.toLocaleString()}</span>
                  <div style={{ fontSize: "0.5rem", color: "rgba(232,232,208,0.3)", fontFamily: '"Press Start 2P", monospace', marginTop: "4px" }}>Contribution</div>
                </div>
              </div>
            </div>

            <button
              onClick={() => onNavigate("/?guild=1")}
              style={{
                marginTop: "10px", width: "100%", padding: "10px",
                fontFamily: '"Press Start 2P", monospace', fontSize: "0.8rem",
                color: "#00e5ff", border: `1px solid ${gColor}40`,
                background: `${gColor}08`, cursor: "pointer",
              }}
            >
              VIEW DETAILS ▶
            </button>
          </Panel>

          {/* Right: GitHub Grass */}
          <Panel borderColor="rgba(76,175,80,0.3)">
            <SectionTitle text="GITHUB GRASS" color="#4caf50" />
            <div>
              {/* Month labels */}
              <div style={{ display: "flex", marginLeft: "42px", marginBottom: "4px", height: "16px", position: "relative" }}>
                {monthLabels.map((m, i) => (
                  <div key={i} style={{ position: "absolute", left: `${m.week * (CELL + GAP)}px`, fontSize: "0.7rem", color: "rgba(232,232,208,0.25)", fontFamily: '"Press Start 2P", monospace' }}>
                    {m.label}
                  </div>
                ))}
              </div>
              <div style={{ display: "flex", gap: `${GAP}px` }}>
                <div style={{ display: "flex", flexDirection: "column", gap: `${GAP}px`, paddingTop: "2px" }}>
                  {DAYS_LABELS.map((d, i) => (
                    <div key={i} style={{ height: `${CELL}px`, display: "flex", alignItems: "center", fontSize: "0.7rem", color: "rgba(232,232,208,0.2)", fontFamily: '"Press Start 2P", monospace' }}>
                      {d}
                    </div>
                  ))}
                </div>
                <div style={{ display: "flex", gap: `${GAP}px` }}>
                  {contribs.map((week, wi) => (
                    <div key={wi} style={{ display: "flex", flexDirection: "column", gap: `${GAP}px` }}>
                      {week.map((count, di) => {
                        const d = new Date(); d.setDate(d.getDate() - (52 - wi) * 7 + di);
                        const ds = `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`;
                        return (
                          <div
                            key={`${wi}-${di}`}
                            onMouseEnter={(e) => setTooltip({ x: e.clientX, y: e.clientY, date: ds, count })}
                            onMouseLeave={() => setTooltip(null)}
                            style={{
                              width: `${CELL}px`, height: `${CELL}px`,
                              background: count === 0 ? "rgba(255,255,255,0.03)" : `rgba(76,175,80,${0.08 + count * 0.1})`,
                              border: "1px solid rgba(255,255,255,0.04)",
                              cursor: "pointer",
                            }}
                          />
                        );
                      })}
                    </div>
                  ))}
                </div>
              </div>
              {/* Legend + Streak */}
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginTop: "8px" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
                  <span style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.25)", fontFamily: '"Press Start 2P", monospace' }}>Less</span>
                  {[0, 1, 2, 3, 4].map((v) => (
                    <div key={v} style={{ width: "8px", height: "8px", background: v === 0 ? "rgba(255,255,255,0.03)" : `rgba(76,175,80,${0.08 + v * 0.12})`, border: "1px solid rgba(255,255,255,0.04)" }} />
                  ))}
                  <span style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.25)", fontFamily: '"Press Start 2P", monospace' }}>More</span>
                </div>
                <div style={{ fontSize: "0.7rem", color: "#f0c040", fontFamily: '"Press Start 2P", monospace' }}>
                  MAX STREAK: <span style={{ color: "#4caf50" }}>{streak}</span> DAYS
                </div>
              </div>
            </div>
            {/* Tooltip */}
            {tooltip && (
              <div style={{ position: "fixed", left: tooltip.x + 10, top: tooltip.y - 28, background: "rgba(0,0,0,0.92)", border: "1px solid #f0c040", padding: "4px 8px", fontSize: "0.8rem", color: "#e8e8d0", pointerEvents: "none", zIndex: 100, fontFamily: '"Press Start 2P", monospace' }}>
                {tooltip.date}: {tooltip.count} commits
              </div>
            )}
          </Panel>
        </div>

        {/* ═══ Bottom Row: 3 columns ═══ */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: "14px", flex: 1, minHeight: 0 }}>
          {/* Left: Language Status */}
          <Panel borderColor="rgba(191,0,255,0.3)">
            <SectionTitle text="LANGUAGES" color="#bf00ff" />
            <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
              {MOCK.langs.map((lang, i) => (
                <motion.div
                  key={lang.name}
                  initial={{ opacity: 0, x: -8 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.3 + i * 0.08, duration: 0.35, ease: steppedEase(6) }}
                >
                  <div style={{ display: "flex", alignItems: "center", gap: "6px", marginBottom: "3px" }}>
                    <span>{lang.icon}</span>
                    <span style={{ fontSize: "0.8rem", color: lang.color, minWidth: "72px", fontFamily: '"Press Start 2P", monospace' }}>{lang.name}</span>
                    <div style={{ flex: 1, height: "8px", border: "1px solid rgba(255,255,255,0.06)", background: "rgba(0,0,0,0.4)", position: "relative", overflow: "hidden" }}>
                      <ProgressBarFill pct={lang.pct} color={lang.color} delay={0.4 + i * 0.08} />
                    </div>
                    <span style={{ fontSize: "0.7rem", color: lang.color, minWidth: "28px", textAlign: "right", fontFamily: '"Press Start 2P", monospace' }}>{lang.pct}%</span>
                  </div>
                  <div style={{ display: "flex", gap: "6px", marginLeft: "22px", flexWrap: "wrap" }}>
                    {(["atk", "def", "mag", "spd", "luk"] as const).map((s) => (
                      <span key={s} style={{ fontSize: "0.6rem", color: "rgba(232,232,208,0.2)", fontFamily: '"Press Start 2P", monospace' }}>
                        {s.toUpperCase()} {lang[s]}
                      </span>
                    ))}
                  </div>
                </motion.div>
              ))}
            </div>
          </Panel>

          {/* Center: Goal + Title */}
          <Panel borderColor="rgba(240,192,64,0.3)">
            <SectionTitle text="GOALS" color="#f0c040" />
            <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "8px" }}>
              <span style={{ fontSize: "2rem" }}>🏁</span>
              <span style={{ fontSize: "0.8rem", color: "rgba(232,232,208,0.6)", fontFamily: '"Press Start 2P", monospace' }}>
                TARGET: {MOCK.goal.target.toLocaleString()} COMMITS
              </span>
            </div>
            <div style={{ height: "16px", border: "2px solid rgba(255,255,255,0.1)", background: "rgba(0,0,0,0.4)", position: "relative", overflow: "hidden", marginBottom: "6px" }}>
              <ProgressBarFill pct={(MOCK.goal.current / MOCK.goal.target) * 100} color="#4caf50" delay={0.4} />
            </div>
            <div style={{ fontSize: "0.8rem", color: "#4caf50", textAlign: "right", fontFamily: '"Press Start 2P", monospace' }}>
              {MOCK.goal.current.toLocaleString()} / {MOCK.goal.target.toLocaleString()}
            </div>

            <div style={{ marginTop: "16px", paddingTop: "12px", borderTop: "1px solid rgba(255,255,255,0.06)" }}>
              <SectionTitle text="TITLE" color="#f0c040" />
              <div style={{ textAlign: "center", padding: "8px" }}>
                <span style={{ fontSize: "2.4rem" }}>👑</span>
                <div style={{ fontSize: "0.9rem", color: "#f0c040", marginTop: "4px", fontFamily: '"Press Start 2P", monospace' }}>{MOCK.title.name}</div>
                <div style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.35)", marginTop: "6px", lineHeight: 1.6, fontFamily: '"Press Start 2P", monospace' }}>
                  "{MOCK.title.line}"
                </div>
              </div>
            </div>
          </Panel>

          {/* Right: (empty or mini stats) */}
          <Panel borderColor="rgba(0,229,255,0.2)">
            <SectionTitle text="QUICK STATS" color="#00e5ff" />
            <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
              <MiniStat label="Today's Commits" value="8" />
              <MiniStat label="This Week" value="34" />
              <MiniStat label="Total PRs" value="127" />
              <MiniStat label="Open Issues" value="3" />
              <MiniStat label="Streak" value={`${streak} days`} />
            </div>
          </Panel>
        </div>

      </div>
    </div>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
      <span style={{ fontSize: "0.7rem", color: "rgba(232,232,208,0.3)", fontFamily: '"Press Start 2P", monospace' }}>{label}</span>
      <span style={{ fontSize: "0.9rem", color: "#00e5ff", fontFamily: '"Press Start 2P", monospace' }}>{value}</span>
    </div>
  );
}
