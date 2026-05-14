import { useState, useEffect } from "react";
import { motion, type Variants } from "framer-motion";
import { useNavigate } from "react-router";
import { ParticleBackground } from "./components/ParticleBackground";
import { TitleLogo } from "./components/TitleLogo";
import { LogoutConfirmDialog } from "./components/LogoutConfirmDialog";
import { GuildPreviewList } from "./components/GuildPreviewList";
import { TitleActions } from "./components/TitleActions";
import { TitleUserBadge } from "./components/TitleUserBadge";
import { AudioTogglePanel } from "./components/AudioTogglePanel";
import { AUDIO_ASSETS } from "./features/audio/audioAssets";
import { beginLogin, fetchMe, logout } from "./features/auth/api";
import type { CurrentUser } from "./features/auth/types";
import { hasCompletedInitialProfile } from "./features/profile/initialProfile";
import { useTitleAudio } from "./hooks/useTitleAudio";

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
  const {
    audioRefs,
    isBgmEnabled,
    isSeEnabled,
    playModalCancel,
    playModalConfirm,
    playModalOpen,
    playTitleStart,
  } = useTitleAudio();
  const [isDay, setIsDay] = useState(() => isDaytime(new Date().getHours()));
  const [currentUser, setCurrentUser] = useState<CurrentUser | null>(null);
  const [isInitialProfileCompleted, setIsInitialProfileCompleted] = useState(false);
  const [isLogoutDialogOpen, setIsLogoutDialogOpen] = useState(false);
  const [isStarting, setIsStarting] = useState(false);

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
  const handleStart = async () => {
    if (isStarting) {
      return;
    }

    setIsStarting(true);
    await playTitleStart();
    navigate(isInitialProfileCompleted ? "/home" : "/profile");
  };
  const closeLogoutDialog = () => setIsLogoutDialogOpen(false);
  const cancelLogout = () => {
    playModalCancel();
    closeLogoutDialog();
  };
  const submitLogout = () => {
    playModalConfirm();
    void confirmLogout();
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
      <audio
        ref={audioRefs.titleBgmRef}
        src={AUDIO_ASSETS.bgm.title}
        loop
        preload="auto"
        muted={!isBgmEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.confirmModalSeRef}
        src={AUDIO_ASSETS.se.confirmModal}
        preload="auto"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.modalCancelSeRef}
        src={AUDIO_ASSETS.se.modalCancel}
        preload="auto"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.modalConfirmSeRef}
        src={AUDIO_ASSETS.se.modalConfirm}
        preload="auto"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.titleStartSeRef}
        src={AUDIO_ASSETS.se.titleStart}
        preload="auto"
        muted={!isSeEnabled}
        aria-hidden="true"
      />

      <AudioTogglePanel />

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

      {currentUser && <TitleUserBadge user={currentUser} />}

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
          <TitleActions
            isLoggedIn={currentUser !== null}
            onLogin={handleLogin}
            onLogoutClick={() => setIsLogoutDialogOpen(true)}
            onStart={handleStart}
            isStarting={isStarting}
          />
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
          <div
            style={{
              flex: 1,
              height: "1px",
              background: "linear-gradient(90deg, transparent, rgba(255,247,220,0.7))",
              boxShadow: "0 0 8px rgba(255,247,220,0.22)",
            }}
          />
          <span
            style={{
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.72rem",
              color: "#fff7dc",
              letterSpacing: "0.2em",
              padding: "4px 10px",
              border: "1px solid rgba(255,247,220,0.28)",
              background: "rgba(8, 12, 18, 0.52)",
              boxShadow: "0 0 0 1px rgba(0,0,0,0.55), 0 0 14px rgba(255,247,220,0.16)",
              textShadow: "1px 1px 0 #000, 0 0 8px rgba(255,247,220,0.38)",
            }}
          >
            SELECT YOUR GUILD
          </span>
          <div
            style={{
              flex: 1,
              height: "1px",
              background: "linear-gradient(90deg, rgba(255,247,220,0.7), transparent)",
              boxShadow: "0 0 8px rgba(255,247,220,0.22)",
            }}
          />
        </motion.div>

        <motion.div variants={itemVariants}>
          <GuildPreviewList />
        </motion.div>

        {import.meta.env.DEV && (
          <motion.div
            variants={itemVariants}
            style={{ marginTop: "3rem", display: "flex", gap: "8px", justifyContent: "center" }}
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

      <LogoutConfirmDialog
        isOpen={isLogoutDialogOpen}
        onOpen={playModalOpen}
        onCancel={cancelLogout}
        onConfirm={submitLogout}
      />
    </div>
  );
}

export default App;
