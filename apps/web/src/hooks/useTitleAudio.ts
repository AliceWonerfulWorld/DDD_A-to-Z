import { useCallback, useEffect, useRef } from "react";
import { useAudioSettings } from "../features/audio/useAudioSettings";

// SE playback safety timeout to avoid blocking navigation if the ended event never fires.
const SE_PLAYBACK_TIMEOUT_MS = 700;

export function useTitleAudio() {
  const { isBgmEnabled, isSeEnabled } = useAudioSettings();
  const titleBgmRef = useRef<HTMLAudioElement | null>(null);
  const confirmModalSeRef = useRef<HTMLAudioElement | null>(null);
  const modalCancelSeRef = useRef<HTMLAudioElement | null>(null);
  const modalConfirmSeRef = useRef<HTMLAudioElement | null>(null);
  const titleStartSeRef = useRef<HTMLAudioElement | null>(null);
  const hasAttemptedBgmPlayRef = useRef(false);

  useEffect(() => {
    const audio = titleBgmRef.current;
    if (!audio) {
      return;
    }

    let isUnlocked = false;
    audio.volume = 0.42;

    const removeUnlockListeners = () => {
      window.removeEventListener("pointerdown", unlockBgm);
      window.removeEventListener("keydown", unlockBgm);
    };

    const playBgm = () => {
      if (isUnlocked) {
        return;
      }

      hasAttemptedBgmPlayRef.current = true;

      void audio
        .play()
        .then(() => {
          isUnlocked = true;
          removeUnlockListeners();
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
      audio.pause();
    };
  }, []);

  useEffect(() => {
    if (titleBgmRef.current) {
      titleBgmRef.current.muted = !isBgmEnabled;

      if (isBgmEnabled && !hasAttemptedBgmPlayRef.current) {
        hasAttemptedBgmPlayRef.current = true;
        void titleBgmRef.current.play().catch(() => {});
      }
    }
  }, [isBgmEnabled]);

  useEffect(() => {
    const seAudios = [
      confirmModalSeRef.current,
      modalCancelSeRef.current,
      modalConfirmSeRef.current,
      titleStartSeRef.current,
    ];

    for (const audio of seAudios) {
      if (audio) {
        audio.muted = !isSeEnabled;
      }
    }
  }, [isSeEnabled]);

  const playSe = useCallback(
    (audio: HTMLAudioElement | null) => {
      if (!audio || !isSeEnabled) {
        return;
      }

      audio.currentTime = 0;
      void audio.play().catch(() => {
        // Browser autoplay restrictions can still block sound in unusual navigation paths.
      });
    },
    [isSeEnabled],
  );

  const playSeUntilEnd = useCallback(
    (audio: HTMLAudioElement | null) => {
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

        audio.currentTime = 0;
        audio.addEventListener("ended", finish, { once: true });
        audio.addEventListener("error", finish, { once: true });
        timeoutId = window.setTimeout(finish, SE_PLAYBACK_TIMEOUT_MS);

        void audio.play().catch(finish);
      });
    },
    [isSeEnabled],
  );

  const playModalCancel = useCallback(() => {
    playSe(modalCancelSeRef.current);
  }, [playSe]);

  const playModalConfirm = useCallback(() => {
    playSe(modalConfirmSeRef.current);
  }, [playSe]);

  const playModalOpen = useCallback(() => {
    playSe(confirmModalSeRef.current);
  }, [playSe]);

  const playTitleStart = useCallback(() => {
    return playSeUntilEnd(titleStartSeRef.current);
  }, [playSeUntilEnd]);

  return {
    audioRefs: {
      titleBgmRef,
      confirmModalSeRef,
      modalCancelSeRef,
      modalConfirmSeRef,
      titleStartSeRef,
    },
    isBgmEnabled,
    isSeEnabled,
    playModalCancel,
    playModalConfirm,
    playModalOpen,
    playTitleStart,
  };
}
