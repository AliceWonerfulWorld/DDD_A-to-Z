import { useState, useEffect, useCallback } from "react";
import { motion, type Variants } from "framer-motion";
import { ParticleBackground } from "./components/ParticleBackground";
import { TitleLogo } from "./components/TitleLogo";
import { GitHubLoginButton } from "./components/GitHubLoginButton";
import { InitialProfile } from "./components/InitialProfile";
import { ContributionAnalysis } from "./components/ContributionAnalysis";

// ギルド一覧（チラ見せ用）
const GUILDS = [
  { name: "Rust", color: "#ff6b35", icon: "🦀" },
  { name: "Python", color: "#3776ab", icon: "🐍" },
  { name: "Go", color: "#00acd7", icon: "🐹" },
  { name: "TypeScript", color: "#3178c6", icon: "📘" },
  { name: "Java", color: "#f89820", icon: "☕" },
  { name: "Haskell", color: "#5d4f85", icon: "λ" },
  { name: "Zig", color: "#f7a41d", icon: "⚡" },
] as const;

const containerVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.15, delayChildren: 0.3 },
  },
};

const itemVariants: Variants = {
  hidden: { opacity: 0, y: 30 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: "easeOut" } },
};

/** 時刻（hour）から昼夜を判定する */
function isDaytime(hour: number): boolean {
  return hour >= 6 && hour < 18;
}

function App() {
  const [isDay, setIsDay] = useState(() => isDaytime(new Date().getHours()));
  const [isProfileMode, setIsProfileMode] = useState(() =>
    new URLSearchParams(window.location.search).get("profile") === "1",
  );
  const isAnalysisMode =
    new URLSearchParams(window.location.search).get("analysis") === "1";

  // 毎分チェックして日没・夜明けをリアルタイム反映
  useEffect(() => {
    const tick = () => setIsDay(isDaytime(new Date().getHours()));
    const id = setInterval(tick, 60_000);
    return () => clearInterval(id);
  }, []);

  const handleLogin = useCallback(() => {
    // TODO: GitHub OAuth 連携実装後にここでリダイレクト
    window.location.href = "/api/auth/github";
  }, []);

  if (isProfileMode) {
    return (
      <InitialProfile
        onComplete={() => {
          window.location.href = "/?analysis=1";
        }}
      />
    );
  }

  if (isAnalysisMode) {
    return (
      <ContributionAnalysis
        onComplete={() => {
          window.location.href = "/";
        }}
      />
    );
  }

  const bgImage = isDay ? "url('/pixel-town-day.png')" : "url('/pixel-town-night.png')";
  const overlay = isDay
    ? "linear-gradient(180deg, rgba(20,40,80,0.45) 0%, rgba(20,40,80,0.22) 50%, rgba(20,40,80,0.55) 100%)"
    : "linear-gradient(180deg, rgba(10,10,30,0.72) 0%, rgba(10,10,30,0.45) 50%, rgba(10,10,30,0.82) 100%)";

  return (
    <div
      style={{
        minHeight: "100svh",
        backgroundImage: bgImage,
        backgroundSize: "cover",
        backgroundPosition: "center bottom",
        backgroundRepeat: "no-repeat",
        position: "relative",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        padding: "2rem 1rem",
        overflow: "hidden",
        /* 背景画像の切り替えをフェードさせる */
        transition: "background-image 2s ease",
      }}
    >
      {/* グラデーションオーバーレイ（昼夜で調整） */}
      <motion.div
        aria-hidden="true"
        animate={{ opacity: 1 }}
        key={isDay ? "day-overlay" : "night-overlay"}
        initial={{ opacity: 0 }}
        transition={{ duration: 1.5 }}
        style={{
          position: "fixed",
          inset: 0,
          background: overlay,
          pointerEvents: "none",
          zIndex: 0,
        }}
      />

      {/* 背景パーティクル */}
      <ParticleBackground />

      {/* スキャンライン風オーバーレイ */}
      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0,0,0,0.04) 2px, rgba(0,0,0,0.04) 4px)",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      {/* メインコンテンツ */}
      <motion.main
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        style={{
          position: "relative",
          zIndex: 2,
          textAlign: "center",
          maxWidth: "860px",
          width: "100%",
        }}
      >
        {/* シーズンバッジ */}
        <motion.div variants={itemVariants} style={{ marginBottom: "2rem" }}>
          <span
            style={{
              display: "inline-block",
              padding: "4px 12px",
              border: "2px solid #ffd700",
              color: "#ffd700",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.75rem",
              letterSpacing: "0.15em",
              boxShadow: "3px 3px 0 #997e00",
            }}
          >
            ▶ SEASON 1 — NOW OPEN
          </span>
        </motion.div>

        {/* タイトルロゴ */}
        <motion.div variants={itemVariants} style={{ marginBottom: "1.5rem" }}>
          <TitleLogo />
        </motion.div>

        {/* キャッチコピー */}
        <motion.p
          variants={itemVariants}
          style={{
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(1rem, 2.5vw, 1.4rem)",
            color: "#e8e8d0",
            marginBottom: "3rem",
            lineHeight: 1.8,
            letterSpacing: "0.05em",
          }}
        >
          君のコードが、推しの力になる。
        </motion.p>

        {/* GitHubログインボタン */}
        <motion.div variants={itemVariants} style={{ marginBottom: "4rem" }}>
          <GitHubLoginButton onClick={handleLogin} />
        </motion.div>

        {/* 区切り線 */}
        <motion.div
          variants={itemVariants}
          style={{
            display: "flex",
            alignItems: "center",
            gap: "1rem",
            marginBottom: "2rem",
          }}
        >
          <div style={{ flex: 1, height: "2px", background: "#ffffff20" }} />
          <span
            style={{
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.7rem",
              color: "#ffffff60",
              letterSpacing: "0.2em",
            }}
          >
            SELECT YOUR GUILD
          </span>
          <div style={{ flex: 1, height: "2px", background: "#ffffff20" }} />
        </motion.div>

        {/* ギルド一覧（チラ見せ） */}
        <motion.div
          variants={itemVariants}
          style={{
            display: "flex",
            flexWrap: "wrap",
            gap: "12px",
            justifyContent: "center",
          }}
        >
          {GUILDS.map((guild, i) => (
            <motion.div
              key={guild.name}
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.8 + i * 0.08, duration: 0.4 }}
              whileHover={{ y: -4, scale: 1.05 }}
              style={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                gap: "6px",
                padding: "12px 16px",
                border: `2px solid ${guild.color}`,
                boxShadow: `3px 3px 0 ${guild.color}80`,
                background: `${guild.color}15`,
                cursor: "pointer",
                minWidth: "90px",
              }}
            >
              <span style={{ fontSize: "1.5rem" }}>{guild.icon}</span>
              <span
                style={{
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.7rem",
                  color: guild.color,
                  letterSpacing: "0.05em",
                }}
              >
                {guild.name}
              </span>
            </motion.div>
          ))}
        </motion.div>

        {/* フッター注記 */}
        <motion.p
          variants={itemVariants}
          style={{
            marginTop: "3rem",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.65rem",
            color: "#ffffff40",
            letterSpacing: "0.1em",
          }}
        >
          ※ ギルドの選択は後で変更できます。まずはGitHubでログインしてください。
        </motion.p>

        {import.meta.env.DEV && (
          <motion.div variants={itemVariants} style={{ marginTop: "2rem", display: "flex", gap: "8px", justifyContent: "center" }}>
            <button
              onClick={() => setIsProfileMode(true)}
              style={{
                padding: "8px 16px",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.8rem",
                background: "transparent",
                border: "1px solid #ffffff40",
                color: "#ffffff80",
                cursor: "pointer",
              }}
            >
              [DEV] ユーザー登録
            </button>
            <button
              onClick={() => { window.location.href = "/?analysis=1"; }}
              style={{
                padding: "8px 16px",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.8rem",
                background: "transparent",
                border: "1px solid #00f5ff40",
                color: "#00f5ff80",
                cursor: "pointer",
              }}
            >
              [DEV] 解析画面
            </button>
          </motion.div>
        )}
      </motion.main>
    </div>
  );
}

export default App;
