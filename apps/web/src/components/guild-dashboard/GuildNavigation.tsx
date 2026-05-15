import { motion } from "framer-motion";
import { useCallback, useRef } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { steppedEase } from "../../lib/animationUtils";

const GUILD_NAV_SE_TIMEOUT_MS = 500;

interface GuildNavigationProps {
  onNavigate: (path: string) => void | Promise<void>;
}

const destinations = [
  { label: "GUILD DETAILS", path: "/guild/details" },
  { label: "GUILD TOWN", path: "/guild/town" },
];

export function GuildNavigation({ onNavigate }: GuildNavigationProps) {
  const { isSeEnabled } = useAudioSettings();
  const guildNavSelectSeRef = useRef<HTMLAudioElement | null>(null);
  const navigationInProgressRef = useRef(false);

  const playSelectSeUntilEnd = useCallback(() => {
    const audio = guildNavSelectSeRef.current;
    if (!audio || !isSeEnabled) {
      return Promise.resolve();
    }

    return new Promise<void>((resolve) => {
      let timeoutId: number | undefined;

      const finish = () => {
        audio.removeEventListener("ended", finish);
        audio.removeEventListener("error", finish);
        if (timeoutId !== undefined) {
          window.clearTimeout(timeoutId);
        }
        resolve();
      };

      if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
        audio.load();
      }

      audio.currentTime = 0;
      audio.addEventListener("ended", finish, { once: true });
      audio.addEventListener("error", finish, { once: true });
      timeoutId = window.setTimeout(finish, GUILD_NAV_SE_TIMEOUT_MS);

      void audio.play().catch(finish);
    });
  }, [isSeEnabled]);

  const navigateWithSe = useCallback(
    async (path: string) => {
      if (navigationInProgressRef.current) {
        return;
      }

      navigationInProgressRef.current = true;
      try {
        await playSelectSeUntilEnd();
        await onNavigate(path);
      } finally {
        navigationInProgressRef.current = false;
      }
    },
    [onNavigate, playSelectSeUntilEnd],
  );

  return (
    <motion.nav
      initial={{ opacity: 0, y: 18 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.12, duration: 0.34, ease: steppedEase(6) }}
      aria-label="Guild navigation"
      style={{
        position: "fixed",
        left: "50%",
        bottom: "calc(env(safe-area-inset-bottom, 0px) + clamp(18px, 5.2vh, 54px))",
        zIndex: 3,
        display: "grid",
        gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
        gap: "clamp(10px, 2vw, 18px)",
        width: "min(calc(100vw - 28px), 600px)",
        transform: "translateX(-50%)",
      }}
    >
      <audio
        ref={guildNavSelectSeRef}
        src={AUDIO_ASSETS.se.homeNavSelect}
        preload="none"
        aria-hidden="true"
      />
      {destinations.map((destination) => (
        <motion.button
          key={destination.path}
          type="button"
          whileHover={{ y: -3, scale: 1.02 }}
          whileTap={{ y: 2, scale: 0.98 }}
          onClick={() => void navigateWithSe(destination.path)}
          style={{
            minHeight: "clamp(46px, 7.5vh, 66px)",
            border: "3px solid rgba(244, 236, 208, 0.72)",
            borderBottomColor: "rgba(55, 44, 35, 0.98)",
            borderRightColor: "rgba(55, 44, 35, 0.98)",
            background: "linear-gradient(180deg, rgba(18, 19, 24, 0.86), rgba(4, 7, 15, 0.9))",
            boxShadow:
              "0 0 0 2px rgba(0,0,0,0.78), 5px 5px 0 rgba(0,0,0,0.45), inset 0 0 16px rgba(244,236,208,0.08)",
            color: "#fff8d7",
            cursor: "pointer",
            fontFamily: "inherit",
            fontSize: "clamp(0.56rem, 1.35vw, 0.78rem)",
            lineHeight: 1.5,
            padding: "10px 12px",
            textAlign: "center",
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          [ {destination.label} ]
        </motion.button>
      ))}
    </motion.nav>
  );
}
