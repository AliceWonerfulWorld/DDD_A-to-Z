import { useCallback, useEffect, useRef, useState } from "react";

export function useTitleAudio() {
  const titleBgmRef = useRef<HTMLAudioElement | null>(null);
  const confirmModalSeRef = useRef<HTMLAudioElement | null>(null);
  const modalCancelSeRef = useRef<HTMLAudioElement | null>(null);
  const modalConfirmSeRef = useRef<HTMLAudioElement | null>(null);
  const [isBgmMuted, setIsBgmMuted] = useState(false);

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
      titleBgmRef.current.muted = isBgmMuted;
    }
  }, [isBgmMuted]);

  const playSe = useCallback((audio: HTMLAudioElement | null) => {
    if (!audio) {
      return;
    }

    audio.currentTime = 0;
    void audio.play().catch(() => {
      // Browser autoplay restrictions can still block sound in unusual navigation paths.
    });
  }, []);

  const toggleBgm = useCallback(() => {
    setIsBgmMuted((current) => {
      const shouldMute = !current;

      if (!shouldMute) {
        void titleBgmRef.current?.play().catch(() => {});
      }

      return shouldMute;
    });
  }, []);

  const playModalCancel = useCallback(() => {
    playSe(modalCancelSeRef.current);
  }, [playSe]);

  const playModalConfirm = useCallback(() => {
    playSe(modalConfirmSeRef.current);
  }, [playSe]);

  const playModalOpen = useCallback(() => {
    playSe(confirmModalSeRef.current);
  }, [playSe]);

  return {
    audioRefs: {
      titleBgmRef,
      confirmModalSeRef,
      modalCancelSeRef,
      modalConfirmSeRef,
    },
    isBgmMuted,
    playModalCancel,
    playModalConfirm,
    playModalOpen,
    toggleBgm,
  };
}
