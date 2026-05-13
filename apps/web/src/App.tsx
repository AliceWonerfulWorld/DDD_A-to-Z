import { useState, useEffect } from "react";
import { motion, type Variants } from "framer-motion";
import { useNavigate } from "react-router";
import { ParticleBackground } from "./components/ParticleBackground";
import { TitleLogo } from "./components/TitleLogo";
import { GitHubLoginButton } from "./components/GitHubLoginButton";
import { beginLogin, fetchMe, logout } from "./features/auth/api";
import type { CurrentUser } from "./features/auth/types";
import { hasCompletedInitialProfile } from "./features/profile/initialProfile";

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
  const navigate = useNavigate();
  const [isDay, setIsDay] = useState(() => isDaytime(new Date().getHours()));
  const [currentUser, setCurrentUser] = useState<CurrentUser | null>(null);
  const [isInitialProfileCompleted, setIsInitialProfileCompleted] = useState(false);
  const [isLogoutDialogOpen, setIsLogoutDialogOpen] = useState(false);

  // 毎分チェックして日没・夜明けをリアルタイム反映
  useEffect(() => {
    const tick = () => setIsDay(isDaytime(new Date().getHours()));
    const id = setInterval(tick, 60_000);
    return () => clearInterval(id);
  }, []);

  // ログイン状態を確認
  useEffect(() => {
    fetchMe()
      .then((user) => {
        setCurrentUser(user);
        setIsInitialProfileCompleted(user ? hasCompletedInitialProfile(user.id) : false);
      })
      .catch(() => {});
  }, []);

  const handleLogin = beginLogin;
  const confirmLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error("failed to logout", error);
    } finally {
      setCurrentUser(null);
      setIsInitialProfileCompleted(false);
      setIsLogoutDialogOpen(false);
    }
  };
  const handleStart = () => {
    navigate(isInitialProfileCompleted ? "/home" : "/profile");
  };
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

      {currentUser && (
        <motion.div
          initial={{ opacity: 0, y: -12 }}
          animate={{ opacity: 1, y: 0 }}
          style={{
            position: "fixed",
            top: "clamp(14px, 3vw, 28px)",
            right: "clamp(14px, 3vw, 28px)",
            zIndex: 4,
            display: "inline-flex",
            alignItems: "center",
            gap: "10px",
            maxWidth: "min(calc(100vw - 28px), 360px)",
            padding: "10px 16px",
            border: "2px solid #39ff14",
            boxShadow: "3px 3px 0 #1a7a00",
            background: "rgba(4, 18, 12, 0.82)",
            backdropFilter: "blur(2px)",
          }}
        >
          <img
            src={currentUser.avatar_url}
            alt={currentUser.username}
            style={{ width: 28, height: 28, borderRadius: "50%", flex: "0 0 auto" }}
          />
          <span
            style={{
              minWidth: 0,
              overflow: "hidden",
              textOverflow: "ellipsis",
              whiteSpace: "nowrap",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.85rem",
              color: "#39ff14",
              letterSpacing: "0.05em",
            }}
          >
            {currentUser.username}
          </span>
          <span
            style={{
              flex: "0 0 auto",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.68rem",
              color: "#39ff1480",
            }}
          >
            ▶ LOGGED IN
          </span>
        </motion.div>
      )}

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

        {/* GitHubログインボタン / ログイン済み表示 */}
        <motion.div variants={itemVariants} style={{ marginBottom: "4rem" }}>
          {currentUser ? (
            <div
              style={{
                display: "inline-flex",
                alignItems: "center",
                flexWrap: "wrap",
                justifyContent: "center",
                gap: "14px",
              }}
            >
              <motion.button
                type="button"
                onClick={handleStart}
                whileHover={{ scale: 1.04 }}
                whileTap={{ y: 3, scale: 0.98 }}
                style={{
                  padding: "12px 22px",
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.9rem",
                  fontWeight: "bold",
                  background: "#ffd700",
                  border: "3px solid #0a0a0a",
                  boxShadow: "4px 4px 0 #0a0a0a",
                  color: "#0a0a0a",
                  cursor: "pointer",
                  letterSpacing: "0.08em",
                }}
              >
                START
              </motion.button>
              <motion.button
                type="button"
                onClick={() => {
                  setIsLogoutDialogOpen(true);
                }}
                whileHover={{ scale: 1.04 }}
                whileTap={{ y: 3, scale: 0.98 }}
                style={{
                  padding: "12px 18px",
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.78rem",
                  background: "rgba(10,10,10,0.55)",
                  border: "2px solid #ff5f56",
                  boxShadow: "3px 3px 0 #7a211d",
                  color: "#ffb0aa",
                  cursor: "pointer",
                  letterSpacing: "0.06em",
                }}
              >
                LOGOUT
              </motion.button>
            </div>
          ) : (
            <GitHubLoginButton onClick={handleLogin} />
          )}
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
          <motion.div
            variants={itemVariants}
            style={{ marginTop: "2rem", display: "flex", gap: "8px", justifyContent: "center" }}
          >
            <button
              onClick={() => navigate("/profile")}
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
              onClick={() => {
                navigate("/analysis");
              }}
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
            <button
              onClick={() => {
                navigate("/mypage");
              }}
              style={{
                padding: "8px 16px",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "0.8rem",
                background: "transparent",
                border: "1px solid #39ff1440",
                color: "#39ff1480",
                cursor: "pointer",
              }}
            >
              [DEV] マイページ
            </button>
          </motion.div>
        )}
      </motion.main>

      {isLogoutDialogOpen && (
        <motion.div
          role="presentation"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          style={{
            position: "fixed",
            inset: 0,
            zIndex: 10,
            display: "grid",
            placeItems: "center",
            padding: "24px",
            background: "rgba(0, 0, 0, 0.62)",
            backdropFilter: "blur(2px)",
          }}
          onClick={() => setIsLogoutDialogOpen(false)}
        >
          <motion.div
            role="dialog"
            aria-modal="true"
            aria-labelledby="logout-dialog-title"
            initial={{ opacity: 0, y: 18, scale: 0.96 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            transition={{ duration: 0.18, ease: "easeOut" }}
            style={{
              width: "min(100%, 420px)",
              border: "4px solid #ffd700",
              borderBottomColor: "#6f4f1c",
              borderRightColor: "#6f4f1c",
              background: "rgba(3, 10, 24, 0.96)",
              boxShadow: "0 0 0 3px rgba(0,0,0,0.9), 10px 10px 0 rgba(0,0,0,0.58)",
              color: "#fff8d7",
              padding: "22px",
              textAlign: "center",
            }}
            onClick={(event) => event.stopPropagation()}
          >
            <div
              id="logout-dialog-title"
              style={{
                color: "#ffd700",
                fontFamily: '"Press Start 2P", "DotGothic16", monospace',
                fontSize: "clamp(0.92rem, 3vw, 1.1rem)",
                lineHeight: 1.7,
                marginBottom: "14px",
              }}
            >
              LOGOUT?
            </div>
            <p
              style={{
                margin: "0 0 22px",
                color: "rgba(255, 248, 215, 0.78)",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "1rem",
                lineHeight: 1.7,
                letterSpacing: "0.04em",
              }}
            >
              本当にログアウトしますか？
            </p>
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
                gap: "12px",
              }}
            >
              <motion.button
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ y: 2, scale: 0.98 }}
                onClick={() => setIsLogoutDialogOpen(false)}
                style={{
                  border: "2px solid rgba(255, 255, 255, 0.34)",
                  borderBottomColor: "rgba(0,0,0,0.78)",
                  borderRightColor: "rgba(0,0,0,0.78)",
                  background: "rgba(255, 255, 255, 0.08)",
                  boxShadow: "3px 3px 0 rgba(0,0,0,0.5)",
                  color: "#fff8d7",
                  cursor: "pointer",
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.9rem",
                  fontWeight: "bold",
                  lineHeight: 1.4,
                  padding: "11px 12px",
                }}
              >
                CANCEL
              </motion.button>
              <motion.button
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ y: 2, scale: 0.98 }}
                onClick={() => {
                  void confirmLogout();
                }}
                style={{
                  border: "2px solid #ffb0aa",
                  borderBottomColor: "#7a211d",
                  borderRightColor: "#7a211d",
                  background: "#ff5f56",
                  boxShadow: "3px 3px 0 rgba(0,0,0,0.5)",
                  color: "#180403",
                  cursor: "pointer",
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.9rem",
                  fontWeight: "bold",
                  lineHeight: 1.4,
                  padding: "11px 12px",
                }}
              >
                LOGOUT
              </motion.button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </div>
  );
}

export default App;
