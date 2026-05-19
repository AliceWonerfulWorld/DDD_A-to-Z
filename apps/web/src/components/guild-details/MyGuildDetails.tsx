import { motion } from "framer-motion";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { PATHS } from "../../constants/paths";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import {
  fetchGuildMembers,
  fetchMyGuild,
  leaveGuild,
  type GuildMemberContribution,
} from "../../features/guild/api";
import { toDisplayGuild, type DisplayGuild } from "../../features/guild/presentation";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { steppedEase } from "../../lib/animationUtils";
import { ApiError } from "../../lib/api/client";
import styles from "./MyGuildDetails.module.css";

interface MyGuildDetailsProps {
  onNavigate: (path: string) => void;
}

interface GuildMember {
  id: string;
  name: string;
  avatar: string;
  totalEarnedCp: number;
  isCurrentUser?: boolean;
}

const panelTransition = { duration: 0.34, ease: steppedEase(6) };

function memberInitials(name: string) {
  const trimmed = name.trim();
  if (!trimmed) {
    return "--";
  }

  return trimmed
    .split(/\s+/)
    .map((part) => part[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
}

function toGuildMember(member: GuildMemberContribution, currentUserID?: string): GuildMember {
  return {
    id: member.user_id,
    name: member.name,
    avatar: member.user_id === currentUserID ? "YOU" : memberInitials(member.name),
    totalEarnedCp: member.total_earned_cp,
    isCurrentUser: member.user_id === currentUserID,
  };
}

export function MyGuildDetails({ onNavigate }: MyGuildDetailsProps) {
  const { isSeEnabled } = useAudioSettings();
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);
  const confirmModalSeRef = useRef<HTMLAudioElement | null>(null);
  const modalCancelSeRef = useRef<HTMLAudioElement | null>(null);
  const modalConfirmSeRef = useRef<HTMLAudioElement | null>(null);
  const [guild, setGuild] = useState<DisplayGuild | null>(null);
  const [members, setMembers] = useState<GuildMember[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLeaving, setIsLeaving] = useState(false);
  const [isLeaveConfirmOpen, setIsLeaveConfirmOpen] = useState(false);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const playSe = useCallback(
    (audio: HTMLAudioElement | null) => {
      if (!audio || !isSeEnabled) {
        return;
      }

      if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
        audio.load();
      }

      audio.currentTime = 0;
      void audio.play().catch(() => {});
    },
    [isSeEnabled],
  );
  const openLeaveConfirm = () => {
    playSe(confirmModalSeRef.current);
    setIsLeaveConfirmOpen(true);
  };
  const cancelLeaveConfirm = () => {
    playSe(modalCancelSeRef.current);
    setIsLeaveConfirmOpen(false);
  };
  const stats = useMemo(() => {
    const displayGuild = guild;

    return [
      { label: "GUILD LEVEL", value: displayGuild ? `LV.${displayGuild.sortOrder + 1}` : "--" },
      {
        label: "LAST SEASON CP",
        value: displayGuild ? displayGuild.previousSeasonCp.toLocaleString() : "--",
      },
      {
        label: "MEMBERS",
        value: displayGuild ? displayGuild.memberCount.toLocaleString() : "--",
      },
    ];
  }, [guild]);

  useEffect(() => {
    let isMounted = true;

    fetchMyGuild()
      .then(async (data) => {
        if (!isMounted) {
          return;
        }

        if (!data?.guild) {
          onNavigate(PATHS.GUILD_SELECT);
          return;
        }

        const displayGuild = toDisplayGuild(data.guild);
        setGuild(displayGuild);
        const memberData =
          data.members && data.members.length > 0
            ? data.members
            : await fetchGuildMembers(displayGuild.id);
        setMembers(memberData.map((member) => toGuildMember(member, data.membership?.user_id)));
      })
      .catch((error) => {
        if (!isMounted) {
          return;
        }

        console.error("failed to fetch my guild details", error);
        setStatusMessage("ギルド情報を取得できませんでした。時間を置いて再度お試しください。");
      })
      .finally(() => {
        if (isMounted) {
          setIsLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [onNavigate]);

  const leaveCurrentGuild = async () => {
    if (isLeaving) {
      return;
    }

    setIsLeaving(true);
    setStatusMessage(null);
    playSe(modalConfirmSeRef.current);
    try {
      await leaveGuild();
      onNavigate(PATHS.GUILD_SELECT);
    } catch (error) {
      if (error instanceof ApiError && error.status === 404) {
        onNavigate(PATHS.GUILD_SELECT);
        return;
      }

      console.error("failed to leave guild", error);
      setStatusMessage("ギルド脱退に失敗しました。少し時間を置いて再度お試しください。");
    } finally {
      setIsLeaving(false);
      setIsLeaveConfirmOpen(false);
    }
  };

  return (
    <main className={styles.main}>
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />
      <audio
        ref={confirmModalSeRef}
        src={AUDIO_ASSETS.se.confirmModal}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={modalCancelSeRef}
        src={AUDIO_ASSETS.se.modalCancel}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={modalConfirmSeRef}
        src={AUDIO_ASSETS.se.modalConfirm}
        preload="none"
        muted={!isSeEnabled}
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
              <div className={styles.emblem} aria-hidden="true" style={{ color: guild?.color }}>
                {isLoading ? "..." : (guild?.icon ?? "--")}
              </div>
              <div>
                <span className={styles.eyebrow}>
                  CURRENT GUILD / {guild?.slug ?? (isLoading ? "SYNCING" : "NONE")}
                </span>
                <h1 className={styles.title} id="my-guild-title">
                  {guild ? `${guild.name} Guild` : "Loading Guild"}
                </h1>
              </div>
            </div>

            <p className={styles.description}>
              {guild?.description ??
                "所属ギルドの作戦情報を同期しています。接続が完了すると現在のギルドが表示されます。"}
            </p>

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
                className={styles.leaveButton}
                type="button"
                disabled={isLeaving || isLoading || !guild}
                whileHover={{ y: -2, scale: 1.02 }}
                whileTap={{ y: 2, scale: 0.98 }}
                onClick={openLeaveConfirm}
              >
                LEAVE GUILD
              </motion.button>
            </div>
            {statusMessage && (
              <p className={styles.statusMessage} role="alert">
                {statusMessage}
              </p>
            )}
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
              <span className={styles.online}>
                {(guild?.memberCount ?? members.length).toLocaleString()} OPERATORS ONLINE
              </span>
            </header>

            <div className={styles.memberList}>
              {members.map((member, index) => (
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
                  </div>
                  <strong className={styles.memberCp}>
                    {member.totalEarnedCp.toLocaleString()} CP
                  </strong>
                </motion.article>
              ))}
              {!isLoading && members.length === 0 && (
                <p className={styles.emptyMembers}>NO OPERATORS ONLINE</p>
              )}
            </div>
          </motion.section>
        </div>
      </div>

      {isLeaveConfirmOpen && (
        <LeaveGuildConfirmDialog
          guildName={guild ? `${guild.name} Guild` : "current guild"}
          isLeaving={isLeaving}
          onCancel={cancelLeaveConfirm}
          onConfirm={() => void leaveCurrentGuild()}
        />
      )}

      <div className={styles.scanline} aria-hidden="true" />
    </main>
  );
}

function LeaveGuildConfirmDialog({
  guildName,
  isLeaving,
  onCancel,
  onConfirm,
}: {
  guildName: string;
  isLeaving: boolean;
  onCancel: () => void;
  onConfirm: () => void;
}) {
  const dialogRef = useRef<HTMLDivElement | null>(null);
  const cancelButtonRef = useRef<HTMLButtonElement | null>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    previousFocusRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
    cancelButtonRef.current?.focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && !isLeaving) {
        event.preventDefault();
        onCancel();
      }
    };

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      previousFocusRef.current?.focus();
    };
  }, [isLeaving, onCancel]);

  return (
    <motion.div
      role="presentation"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className={styles.confirmBackdrop}
      onClick={() => {
        if (!isLeaving) {
          onCancel();
        }
      }}
    >
      <motion.div
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="leave-guild-dialog-title"
        aria-describedby="leave-guild-dialog-description"
        tabIndex={-1}
        initial={{ opacity: 0, y: 18, scale: 0.96 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.18, ease: "easeOut" }}
        className={styles.confirmDialog}
        onClick={(event) => event.stopPropagation()}
      >
        <div className={styles.confirmTitle} id="leave-guild-dialog-title">
          LEAVE GUILD?
        </div>
        <p className={styles.confirmBody} id="leave-guild-dialog-description">
          {guildName} から脱退しますか？
        </p>
        <div className={styles.confirmActions}>
          <motion.button
            ref={cancelButtonRef}
            className={styles.confirmCancelButton}
            type="button"
            disabled={isLeaving}
            whileHover={{ scale: 1.03 }}
            whileTap={{ y: 2, scale: 0.98 }}
            onClick={onCancel}
          >
            CANCEL
          </motion.button>
          <motion.button
            className={styles.confirmLeaveButton}
            type="button"
            disabled={isLeaving}
            whileHover={{ scale: 1.03 }}
            whileTap={{ y: 2, scale: 0.98 }}
            onClick={onConfirm}
          >
            {isLeaving ? "LEAVING..." : "LEAVE"}
          </motion.button>
        </div>
      </motion.div>
    </motion.div>
  );
}
