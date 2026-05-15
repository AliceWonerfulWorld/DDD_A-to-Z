import { useEffect, useRef } from "react";
import { AUDIO_ASSETS } from "../features/audio/audioAssets";
import { useAudioSettings } from "../features/audio/useAudioSettings";

const HOME_BGM_VOLUME = 0.34;
const HOME_BGM_FADE_IN_MS = 520;

export function HomeBgm() {
  const { isBgmEnabled } = useAudioSettings();
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const hasAttemptedPlayRef = useRef(false);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    let isUnlocked = false;
    let fadeFrameId: number | undefined;
    audio.volume = 0;

    const fadeInBgm = () => {
      const startedAt = performance.now();
      const tick = (now: number) => {
        const progress = Math.min(1, (now - startedAt) / HOME_BGM_FADE_IN_MS);
        audio.volume = HOME_BGM_VOLUME * progress;

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
          isUnlocked = true;
          removeUnlockListeners();
          fadeInBgm();
        })
        .catch(() => {
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
      audio.volume = HOME_BGM_VOLUME;
      void audio.play().catch(() => {});
    }
  }, [isBgmEnabled]);

  return (
    <audio
      ref={audioRef}
      src={AUDIO_ASSETS.bgm.home}
      loop
      preload="auto"
      muted={!isBgmEnabled}
      aria-hidden="true"
    />
  );
}
