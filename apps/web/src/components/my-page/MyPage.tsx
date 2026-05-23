import { useState, useMemo, useEffect, type ReactNode } from "react";
import { motion } from "framer-motion";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { useGuardedNavigation } from "../../hooks/useGuardedNavigation";
import {
  fetchMyPage,
  type MyPageResponse,
  type GitHubStats,
  type MyPageBadge,
} from "../../features/mypage/api";
import { updateSelectedBadge } from "../../features/profile/api";
import { updateProfileAPI } from "../../features/profile/api";
import { findGuildBySlug } from "../../features/guild/guildMaster";
import { AvatarPicker } from "../profile/AvatarPicker";
import { NameInput } from "../profile/NameInput";
import styles from "./MyPage.module.css";

interface MyPageProps {
  onNavigate: (path: string) => void;
}

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

const LANG_DISPLAY: Record<string, { icon: string; color: string }> = {
  TypeScript: { icon: "📘", color: "#3178c6" },
  JavaScript: { icon: "📒", color: "#f7df1e" },
  Rust: { icon: "🦀", color: "#ff6b35" },
  Go: { icon: "🐹", color: "#00acd7" },
  Python: { icon: "🐍", color: "#f0c040" },
  Ruby: { icon: "💎", color: "#701516" },
  Java: { icon: "☕", color: "#b07219" },
  Kotlin: { icon: "🅺", color: "#a97bff" },
  Swift: { icon: "🍎", color: "#f05138" },
  "C++": { icon: "⚙️", color: "#f34b7d" },
  C: { icon: "⚙️", color: "#555555" },
  "C#": { icon: "🎯", color: "#178600" },
  PHP: { icon: "🐘", color: "#4f5d95" },
  Shell: { icon: "🐚", color: "#89e051" },
  Dockerfile: { icon: "🐳", color: "#384d54" },
  HTML: { icon: "🌐", color: "#e34c26" },
  CSS: { icon: "🎨", color: "#563d7c" },
  Scala: { icon: "🔥", color: "#c22d40" },
  Dart: { icon: "🎯", color: "#00b4ab" },
  Lua: { icon: "🌙", color: "#000080" },
  Haskell: { icon: "λ", color: "#5e5086" },
};

function langDisplay(name: string): { icon: string; color: string } {
  return LANG_DISPLAY[name] ?? { icon: "◇", color: "#888" };
}

interface LangEntry {
  name: string;
  pct: number;
  count: number;
}

const MOCK_TS_GUILD = findGuildBySlug("typescript");

const MOCK = {
  season: {
    label: "SEASON 1",
    start: "2024/05/01",
    end: "2024/07/31",
    remaining: 52,
  },
  guild: {
    id: MOCK_TS_GUILD?.id ?? "",
    name: MOCK_TS_GUILD?.name ?? "TypeScript",
    slug: MOCK_TS_GUILD?.slug ?? "typescript",
    icon: MOCK_TS_GUILD?.icon ?? "TS",
    color: MOCK_TS_GUILD?.color ?? "#3178c6",
    description:
      MOCK_TS_GUILD?.description ??
      "型の力で安全で堅牢なコードを書く、\nエレガントな戦士たちの集い。",
    member_count: MOCK_TS_GUILD?.memberCount ?? 0,
    rank: 42,
    total_guilds: 156,
    cp: 24680,
    fullName: `${MOCK_TS_GUILD?.name ?? "TypeScript"} GUILD`,
  },
  title: {
    name: "Consistency Master",
    line: "Consistency is key. Daily efforts build the future.",
  },
};

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

function Panel({
  children,
  borderColor,
  className = "",
}: {
  children: ReactNode;
  borderColor?: string;
  className?: string;
}) {
  return (
    <motion.div
      className={`${styles.panel} ${className}`}
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, ease: steppedEase(6) }}
      style={{
        border: `2px solid ${borderColor ?? "rgba(255,255,255,0.08)"}`,
      }}
    >
      {children}
    </motion.div>
  );
}

function ProgressBarFill({
  pct,
  color,
  delay = 0.3,
}: {
  pct: number;
  color: string;
  delay?: number;
}) {
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
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);
  const guardedNavigate = useGuardedNavigation(onNavigate);
  const [mypageData, setMypageData] = useState<MyPageResponse | null>(null);
  const [apiError, setApiError] = useState(false);
  const [selectedBadgeSlug, setSelectedBadgeSlug] = useState<string | null>(null);
  const [badgeError, setBadgeError] = useState<string | null>(null);
  const [isEditingProfile, setIsEditingProfile] = useState(false);
  const [editName, setEditName] = useState("");
  const [editAvatar, setEditAvatar] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  const loadMyPage = () => {
    fetchMyPage()
      .then((data) => {
        setMypageData(data);
        setSelectedBadgeSlug(data.selected_badge_slug);
        setApiError(false);
      })
      .catch(() => setApiError(true));
  };

  useEffect(() => {
    loadMyPage();
  }, []);

  const handleBadgeSelect = async (slug: string | null) => {
    try {
      await updateSelectedBadge(slug);
      setSelectedBadgeSlug(slug);
      setBadgeError(null);
    } catch {
      setBadgeError("Failed to update selected badge");
    }
  };

  const effectiveSelectedSlug = selectedBadgeSlug;

  const activeBadge = useMemo(() => {
    if (!mypageData || !effectiveSelectedSlug) return null;
    return mypageData.badges.find((b) => b.slug === effectiveSelectedSlug) ?? null;
  }, [mypageData, effectiveSelectedSlug]);

  const handleEditProfileClick = () => {
    if (mypageData) {
      setEditName(mypageData.user.username);
      setEditAvatar(mypageData.user.avatar_url);
    }
    setIsEditingProfile(true);
  };

  const handleSaveProfileClick = async () => {
    if (!editName.trim() || isSaving) return;
    setIsSaving(true);
    try {
      await updateProfileAPI(editName, editAvatar);
      setIsEditingProfile(false);
      loadMyPage();
    } catch (e) {
      console.error("failed to update profile", e);
    } finally {
      setIsSaving(false);
    }
  };

  const langEntries: LangEntry[] = useMemo(() => {
    if (!mypageData) return [];
    const summary = mypageData.repositories.language_summary;
    const total = mypageData.repositories.total_count;
    return Object.entries(summary)
      .map(([name, count]) => ({
        name,
        count,
        pct: total > 0 ? Math.round((count / total) * 100) : 0,
      }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 6);
  }, [mypageData]);

  const guild = useMemo(() => {
    const apiGuild = mypageData?.guild;
    if (!apiGuild) return MOCK.guild;
    const master = findGuildBySlug(apiGuild.slug);
    return master ? { ...apiGuild, icon: master.icon, color: master.color } : apiGuild;
  }, [mypageData]);
  const gColor = guild.color ?? "#3178c6";

  return (
    <div
      className="flex flex-col min-h-svh relative overflow-hidden"
      style={{
        background: "radial-gradient(ellipse at 50% 0%, #0d1b2a 0%, #0a0a1a 50%, #050510 100%)",
        fontFamily: '"Press Start 2P", monospace',
        color: "#e8e8d0",
      }}
    >
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />
      {/* City silhouette */}
      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          background: "url('/pixel-town-night.png') center bottom / cover no-repeat",
          opacity: 0.1,
          pointerEvents: "none",
          zIndex: 0,
        }}
      />
      {/* Scanline */}
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

      {/* ─── Header ─── */}
      <motion.header
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.4, ease: steppedEase(6) }}
        style={{
          position: "relative",
          zIndex: 3,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "10px 24px",
          borderBottom: "2px solid rgba(240,192,64,0.3)",
          background: "rgba(0,0,0,0.6)",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <span style={{ fontSize: "2rem" }}>👑</span>
          <span style={{ fontSize: "1rem", color: "#f0c040", letterSpacing: "0.1em" }}>
            MY PAGE
          </span>
          <span style={{ fontSize: "1rem", color: "rgba(240,192,64,0.3)" }}>{">"}</span>
          <span style={{ fontSize: "0.9rem", color: "rgba(232,232,208,0.5)" }}>
            ENGINEER STATUS
          </span>
        </div>
        <button
          onClick={() => void navigateBackWithSe("/home")}
          style={{
            fontFamily: '"Press Start 2P", monospace',
            fontSize: "0.72rem",
            color: "#1b1304",
            background: "#f0c040",
            border: "2px solid #fff3a6",
            borderBottomColor: "#6f4f1c",
            borderRightColor: "#6f4f1c",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.72), 4px 4px 0 rgba(0,0,0,0.48)",
            padding: "9px 14px",
            cursor: "pointer",
            lineHeight: 1.4,
          }}
        >
          BACK [→]
        </button>
      </motion.header>

      {/* ─── Main Content ─── */}
      <div className={styles.mainContainer}>
        {/* ═══ Top Row: 3 columns ═══ */}
        <div className={styles.topRow}>
          {/* Left: Adventurer Profile */}
          <Panel borderColor="rgba(240,192,64,0.3)">
            <SectionTitle text="PROFILE" color="#f0c040" />
            {isEditingProfile ? (
              <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                <AvatarPicker
                  avatarUrl={editAvatar}
                  githubAvatarUrl={mypageData?.user.github_avatar_url}
                  disabled={isSaving}
                  onChange={setEditAvatar}
                />
                <NameInput username={editName} onChange={setEditName} disabled={isSaving} />
                <div
                  style={{
                    display: "flex",
                    gap: "10px",
                    justifyContent: "flex-end",
                    marginTop: "10px",
                  }}
                >
                  <button
                    onClick={() => setIsEditingProfile(false)}
                    disabled={isSaving}
                    style={{
                      fontFamily: '"Press Start 2P", monospace',
                      fontSize: "0.7rem",
                      padding: "8px 12px",
                      background: "rgba(255,255,255,0.1)",
                      color: "#fff",
                      border: "1px solid rgba(255,255,255,0.3)",
                      cursor: "pointer",
                    }}
                  >
                    CANCEL
                  </button>
                  <button
                    onClick={() => void handleSaveProfileClick()}
                    disabled={isSaving}
                    style={{
                      fontFamily: '"Press Start 2P", monospace',
                      fontSize: "0.7rem",
                      padding: "8px 12px",
                      background: "#f0c040",
                      color: "#1b1304",
                      border: "1px solid #fff3a6",
                      cursor: "pointer",
                    }}
                  >
                    SAVE
                  </button>
                </div>
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "15px" }}>
                  {mypageData?.user.avatar_url && (
                    <div
                      style={{
                        width: "64px",
                        height: "64px",
                        border: "2px solid #f0c040",
                        padding: "2px",
                        background: "rgba(0,0,0,0.5)",
                      }}
                    >
                      <img
                        src={mypageData.user.avatar_url}
                        alt="avatar"
                        style={{
                          width: "100%",
                          height: "100%",
                          objectFit: "cover",
                          imageRendering: "pixelated",
                        }}
                      />
                    </div>
                  )}
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div
                      style={{
                        fontFamily: '"Press Start 2P", monospace',
                        fontSize: "1rem",
                        color: "#e8e8d0",
                      }}
                    >
                      {mypageData?.user.username ?? "-"}
                    </div>
                  </div>
                </div>

                <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                  <span style={{ fontSize: "1.4rem" }}>👑</span>
                  <span
                    style={{
                      fontSize: "0.8rem",
                      color: "#f0c040",
                      fontFamily: '"Press Start 2P", monospace',
                    }}
                  >
                    {activeBadge ? activeBadge.name : "-"}
                  </span>
                </div>
                {activeBadge && (
                  <div
                    style={{
                      fontSize: "0.6rem",
                      color: "rgba(232,232,208,0.5)",
                      fontFamily: '"Press Start 2P", monospace',
                      lineHeight: 1.4,
                      marginTop: "4px",
                    }}
                  >
                    {activeBadge.description}
                  </div>
                )}

                <button
                  onClick={handleEditProfileClick}
                  style={{
                    alignSelf: "flex-start",
                    fontFamily: '"Press Start 2P", monospace',
                    fontSize: "0.6rem",
                    color: "#1b1304",
                    background: "#f0c040",
                    border: "2px solid #fff3a6",
                    borderBottomColor: "#6f4f1c",
                    borderRightColor: "#6f4f1c",
                    boxShadow: "0 0 0 2px rgba(0,0,0,0.72), 4px 4px 0 rgba(0,0,0,0.48)",
                    padding: "6px 10px",
                    cursor: "pointer",
                    marginTop: "4px",
                  }}
                >
                  EDIT PROFILE
                </button>
              </div>
            )}

            {/* Season info */}
            <div
              style={{
                marginTop: "12px",
                padding: "10px",
                border: "1px solid rgba(156,39,176,0.3)",
                background: "rgba(156,39,176,0.06)",
              }}
            >
              <div
                style={{ display: "flex", alignItems: "center", gap: "6px", marginBottom: "6px" }}
              >
                <span
                  style={{
                    fontSize: "0.7rem",
                    background: "rgba(156,39,176,0.5)",
                    color: "#e8e8d0",
                    padding: "2px 6px",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  {MOCK.season.label}
                </span>
              </div>
              <div
                style={{
                  fontSize: "0.7rem",
                  color: "rgba(232,232,208,0.4)",
                  fontFamily: '"Press Start 2P", monospace',
                }}
              >
                {MOCK.season.start} 〜 {MOCK.season.end}
              </div>
              <div
                style={{ marginTop: "6px", display: "flex", alignItems: "baseline", gap: "6px" }}
              >
                <span
                  style={{
                    fontSize: "2rem",
                    color: "#f0c040",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  {MOCK.season.remaining}
                </span>
                <span
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.4)",
                    fontFamily: '"Press Start 2P", monospace',
                    marginLeft: "8px",
                  }}
                >
                  DAYS LEFT
                </span>
              </div>
            </div>
          </Panel>

          {/* Center: Guild */}
          <Panel borderColor={`${gColor}40`}>
            <SectionTitle text="GUILD" color={gColor} />
            <div
              style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: "8px" }}
            >
              {/* Emblem */}
              <motion.div
                animate={{ scaleY: [1, 1.03, 1] }}
                transition={{ duration: 2, repeat: Infinity, ease: steppedEase(4) }}
                style={{
                  width: "100px",
                  height: "100px",
                  border: `2px solid ${gColor}`,
                  background: `${gColor}08`,
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  justifyContent: "center",
                  gap: "2px",
                  position: "relative",
                  overflow: "hidden",
                }}
              >
                {(() => {
                  const icon = guild.icon ?? "--";
                  const isValidIcon = /^[A-Z0-9λ]+$/.test(icon);
                  return isValidIcon ? (
                    <img
                      src={`/guild-icons/${icon}.png`}
                      alt={`${guild.name} guild icon`}
                      style={{
                        width: "100%",
                        height: "100%",
                        objectFit: "cover",
                        imageRendering: "pixelated",
                      }}
                    />
                  ) : (
                    <span style={{ fontSize: "3.6rem" }}>{icon}</span>
                  );
                })()}
                {/* Laurel decoration */}
                <span
                  style={{
                    position: "absolute",
                    top: "-4px",
                    fontSize: "1.2rem",
                    color: "#f0c040",
                  }}
                >
                  🏅
                </span>
              </motion.div>

              <div
                style={{
                  fontFamily: '"Press Start 2P", monospace',
                  fontSize: "0.9rem",
                  color: "#00e5ff",
                  letterSpacing: "0.1em",
                  textAlign: "center",
                  lineHeight: 1.4,
                }}
              >
                {guild.name}
              </div>

              <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                <span
                  style={{
                    fontSize: "0.7rem",
                    background: "rgba(156,39,176,0.5)",
                    color: "#e8e8d0",
                    padding: "2px 6px",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  🏷 MEMBER
                </span>
                <span
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.3)",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  since Season 1
                </span>
              </div>

              <div
                style={{
                  fontSize: "0.7rem",
                  color: "rgba(232,232,208,0.4)",
                  textAlign: "center",
                  lineHeight: 1.6,
                  fontFamily: '"Press Start 2P", monospace',
                }}
              >
                {guild.description}
              </div>
            </div>

            {/* Guild stats */}
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "14px",
                marginTop: "12px",
                padding: "12px 8px",
                borderTop: "1px solid rgba(255,255,255,0.06)",
                background: "rgba(0,0,0,0.2)",
              }}
            >
              <div
                style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-end" }}
              >
                <div
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.4)",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  RANK
                </div>
                <div style={{ textAlign: "right" }}>
                  <span
                    style={{
                      fontSize: "1.4rem",
                      color: "#f0c040",
                      fontFamily: '"Press Start 2P", monospace',
                    }}
                  >
                    #{guild.rank ?? "-"}
                  </span>
                  <span
                    style={{
                      fontSize: "0.6rem",
                      color: "rgba(232,232,208,0.3)",
                      fontFamily: '"Press Start 2P", monospace',
                      marginLeft: "8px",
                    }}
                  >
                    / {guild.total_guilds ?? "-"} GUILDS
                  </span>
                </div>
              </div>
              <div
                style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-end" }}
              >
                <div
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.4)",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  GUILD CP
                </div>
                <div style={{ textAlign: "right" }}>
                  <span
                    style={{
                      fontSize: "1.2rem",
                      color: "#00e5ff",
                      fontFamily: '"Press Start 2P", monospace',
                    }}
                  >
                    {(guild.cp ?? 0).toLocaleString()}
                  </span>
                  <div
                    style={{
                      fontSize: "0.5rem",
                      color: "rgba(232,232,208,0.3)",
                      fontFamily: '"Press Start 2P", monospace',
                      marginTop: "4px",
                    }}
                  >
                    Contribution
                  </div>
                </div>
              </div>
              <div
                style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-end" }}
              >
                <div
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.4)",
                    fontFamily: '"Press Start 2P", monospace',
                  }}
                >
                  MEMBERS
                </div>
                <div style={{ textAlign: "right" }}>
                  <span
                    style={{
                      fontSize: "1.2rem",
                      color: "#00e5ff",
                      fontFamily: '"Press Start 2P", monospace',
                    }}
                  >
                    {guild.member_count?.toLocaleString() ?? "-"}
                  </span>
                </div>
              </div>
            </div>

            <button
              onClick={() => void guardedNavigate("/guild")}
              style={{
                marginTop: "10px",
                width: "100%",
                padding: "10px",
                fontFamily: '"Press Start 2P", monospace',
                fontSize: "0.8rem",
                color: "#00e5ff",
                border: `1px solid ${gColor}40`,
                background: `${gColor}08`,
                cursor: "pointer",
              }}
            >
              VIEW DETAILS ▶
            </button>
          </Panel>

          {/* Right: Engineer Status */}
          <Panel borderColor="rgba(0,229,255,0.2)">
            <SectionTitle text="ENGINEER STATUS" color="#00e5ff" />
            {mypageData === undefined ? (
              <div
                style={{
                  fontSize: "0.7rem",
                  color: "rgba(232,232,208,0.3)",
                  fontFamily: '"Press Start 2P", monospace',
                  textAlign: "center",
                  padding: "20px 0",
                }}
              >
                Loading...
              </div>
            ) : mypageData?.github_stats ? (
              <GitHubStatsPanel stats={mypageData.github_stats} />
            ) : (
              <div
                style={{
                  fontSize: "0.7rem",
                  color: "rgba(232,232,208,0.3)",
                  fontFamily: '"Press Start 2P", monospace',
                  textAlign: "center",
                  padding: "20px 0",
                }}
              >
                {apiError ? "Failed to load" : "No data"}
              </div>
            )}
          </Panel>
        </div>

        {/* ═══ Bottom Row: 2 columns ═══ */}
        <div className={styles.bottomRow}>
          {/* Left: Language Status */}
          <Panel borderColor="rgba(191,0,255,0.3)">
            <SectionTitle text="LANGUAGES" color="#bf00ff" />
            <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
              {mypageData === undefined ? (
                <div
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.3)",
                    fontFamily: '"Press Start 2P", monospace',
                    textAlign: "center",
                    padding: "20px 0",
                  }}
                >
                  Loading...
                </div>
              ) : langEntries.length === 0 ? (
                <div
                  style={{
                    fontSize: "0.7rem",
                    color: "rgba(232,232,208,0.3)",
                    fontFamily: '"Press Start 2P", monospace',
                    textAlign: "center",
                    padding: "20px 0",
                  }}
                >
                  {apiError ? "Failed to load" : "No data"}
                </div>
              ) : (
                langEntries.map((lang, i) => {
                  const meta = langDisplay(lang.name);
                  return (
                    <motion.div
                      key={lang.name}
                      initial={{ opacity: 0, x: -8 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: 0.3 + i * 0.08, duration: 0.35, ease: steppedEase(6) }}
                    >
                      <div
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: "6px",
                        }}
                      >
                        <span style={{ fontSize: "0.9rem" }}>{meta.icon}</span>
                        <span
                          style={{
                            fontSize: "0.75rem",
                            color: meta.color,
                            minWidth: "60px",
                            fontFamily: '"Press Start 2P", monospace',
                          }}
                        >
                          {lang.name}
                        </span>
                        <div
                          style={{
                            flex: 1,
                            height: "6px",
                            border: "1px solid rgba(255,255,255,0.06)",
                            background: "rgba(0,0,0,0.4)",
                            position: "relative",
                            overflow: "hidden",
                          }}
                        >
                          <ProgressBarFill
                            pct={lang.pct}
                            color={meta.color}
                            delay={0.4 + i * 0.08}
                          />
                        </div>
                        <span
                          style={{
                            fontSize: "0.65rem",
                            color: meta.color,
                            minWidth: "32px",
                            textAlign: "right",
                            fontFamily: '"Press Start 2P", monospace',
                          }}
                        >
                          {lang.pct}%
                        </span>
                      </div>
                    </motion.div>
                  );
                })
              )}
            </div>
          </Panel>

          {/* Right: Quick Stats */}
          <Panel borderColor="rgba(0,229,255,0.2)">
            <SectionTitle text="CONTRIBUTION POINTS" color="#00e5ff" />
            {mypageData ? (
              <div style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
                <MiniStat
                  icon="💰"
                  label="Balance"
                  value={mypageData.contribution_points.balance.toLocaleString()}
                />
                <MiniStat
                  icon="📈"
                  label="Total Earned"
                  value={mypageData.contribution_points.total_earned.toLocaleString()}
                />
                <MiniStat
                  icon="📉"
                  label="Total Spent"
                  value={mypageData.contribution_points.total_spent.toLocaleString()}
                />
              </div>
            ) : (
              <div
                style={{
                  fontSize: "0.7rem",
                  color: "rgba(232,232,208,0.3)",
                  fontFamily: '"Press Start 2P", monospace',
                  textAlign: "center",
                  padding: "20px 0",
                }}
              >
                {apiError ? "Failed to load" : "Loading..."}
              </div>
            )}
          </Panel>
        </div>

        {/* ═══ Badges Row ═══ */}
        {mypageData && mypageData.badges && mypageData.badges.length > 0 && (
          <Panel borderColor="rgba(240,192,64,0.3)">
            <SectionTitle text="BADGES" color="#f0c040" />
            <div style={{ display: "flex", flexWrap: "wrap", gap: "10px" }}>
              {mypageData.badges.map((badge) => (
                <BadgeCard
                  key={badge.slug}
                  badge={badge}
                  selected={badge.slug === effectiveSelectedSlug}
                  onSelect={handleBadgeSelect}
                />
              ))}
            </div>
            {badgeError && (
              <div
                style={{
                  marginTop: "8px",
                  fontSize: "0.6rem",
                  color: "#ff6b6b",
                  fontFamily: '"Press Start 2P", monospace',
                }}
              >
                {badgeError}
              </div>
            )}
          </Panel>
        )}
      </div>
    </div>
  );
}

function MiniStat({ label, value, icon }: { label: string; value: string; icon?: string }) {
  return (
    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
      <span
        style={{
          fontSize: "0.7rem",
          color: "rgba(232,232,208,0.8)",
          fontFamily: '"Press Start 2P", monospace',
        }}
      >
        {icon ? `${icon} ${label}` : label}
      </span>
      <span
        style={{ fontSize: "0.9rem", color: "#00e5ff", fontFamily: '"Press Start 2P", monospace' }}
      >
        {value}
      </span>
    </div>
  );
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  return `${y}/${m}`;
}

const currentYear = new Date().getFullYear();

function GitHubStatsPanel({ stats }: { stats: GitHubStats }) {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
      <MiniStat icon="⭐" label="Total Stars" value={stats.total_stars.toLocaleString()} />
      <MiniStat
        icon="📝"
        label={`${currentYear} Commits`}
        value={stats.yearly_commits.toLocaleString()}
      />
      <MiniStat icon="🔀" label="Total PRs" value={stats.total_prs.toLocaleString()} />
      <MiniStat icon="🐛" label="Total Issues" value={stats.total_issues.toLocaleString()} />
      <MiniStat icon="📦" label="Public Repos" value={stats.public_repos.toLocaleString()} />
      <MiniStat icon="📅" label="GitHub Started" value={formatDate(stats.github_created_at)} />
      <MiniStat
        icon="🎯"
        label={`${currentYear} Contributions`}
        value={stats.yearly_contributions.toLocaleString()}
      />
    </div>
  );
}

function BadgeCard({
  badge,
  selected,
  onSelect,
}: {
  badge: MyPageBadge;
  selected: boolean;
  onSelect: (slug: string | null) => void;
}) {
  return (
    <motion.button
      initial={{ opacity: 0, scale: 0.85 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.35, ease: steppedEase(6) }}
      onClick={() => onSelect(selected ? null : badge.slug)}
      style={{
        display: "flex",
        alignItems: "center",
        gap: "8px",
        width: "140px",
        border: selected ? "2px solid #f0c040" : "1px solid rgba(240,192,64,0.25)",
        background: selected ? "rgba(240,192,64,0.15)" : "rgba(240,192,64,0.06)",
        padding: "8px 10px",
        cursor: "pointer",
        fontFamily: '"Press Start 2P", monospace',
      }}
    >
      <span style={{ fontSize: "1.3rem" }}>{badge.icon}</span>
      <div
        style={{
          fontSize: "0.7rem",
          color: selected ? "#f0c040" : "rgba(232,232,208,0.7)",
          fontFamily: '"Press Start 2P", monospace',
          lineHeight: 1.4,
          overflow: "hidden",
          textOverflow: "ellipsis",
          whiteSpace: "nowrap",
        }}
      >
        {badge.name}
      </div>
    </motion.button>
  );
}
