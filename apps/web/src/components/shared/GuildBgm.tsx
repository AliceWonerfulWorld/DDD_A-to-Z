import { useEffect, useRef } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";

const GUILD_BGM_VOLUME = 0.36;
const GUILD_BGM_FADE_IN_MS = 620;

interface GuildBgmProps {
  src?: string;
}

export function GuildBgm({ src = AUDIO_ASSETS.bgm.guild }: GuildBgmProps) {
  const { isBgmEnabled } = useAudioSettings();
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const hasAttemptedPlayRef = useRef(false);
  const isUnmountedRef = useRef(false);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    isUnmountedRef.current = false;
    let isUnlocked = false;
    let fadeFrameId: number | undefined;
    audio.volume = 0;

    const fadeInBgm = () => {
      if (isUnmountedRef.current) return;

      const startedAt = performance.now();
      const tick = (now: number) => {
        if (isUnmountedRef.current) return;

        const progress = Math.min(1, (now - startedAt) / GUILD_BGM_FADE_IN_MS);
        audio.volume = GUILD_BGM_VOLUME * progress;

        if (progress < 1) {
          fadeFrameId = window.requestAnimationFrame(tick);
        }
      };

      fadeFrameId = window.requestAnimationFrame(tick);
    };

    const removeUnlockListeners = () => {
      window.removeEventListener("pointerdown", unlockBgm);
      window.removeEventListener("keydown", unlockBgm);
    };

    const playBgm = () => {
      if (isUnlocked) return;

      hasAttemptedPlayRef.current = true;
      void audio
        .play()
        .then(() => {
          if (isUnmountedRef.current) return;

          isUnlocked = true;
          removeUnlockListeners();
          fadeInBgm();
        })
        .catch(() => {
          if (isUnmountedRef.current) return;

          // ブラウザの自動再生制限で止められた場合は、最初のユーザー操作で再試行する。
        });
    };

    const unlockBgm = () => {
      playBgm();
    };

    playBgm();
    window.addEventListener("pointerdown", unlockBgm);
    window.addEventListener("keydown", unlockBgm);

    return () => {
      isUnmountedRef.current = true;
      removeUnlockListeners();
      if (fadeFrameId !== undefined) {
        window.cancelAnimationFrame(fadeFrameId);
      }
      audio.pause();
    };
  }, []);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    audio.muted = !isBgmEnabled;

    if (isBgmEnabled && !hasAttemptedPlayRef.current) {
      hasAttemptedPlayRef.current = true;
      audio.volume = GUILD_BGM_VOLUME;
      void audio.play().catch(() => {});
    }
  }, [isBgmEnabled]);

  return (
    <audio ref={audioRef} src={src} loop preload="auto" muted={!isBgmEnabled} aria-hidden="true" />
  );
}
