import { useCallback, useRef } from "react";
import { AUDIO_ASSETS } from "../features/audio/audioAssets";
import { useAudioSettings } from "../features/audio/useAudioSettings";

const BACK_NAVIGATION_SE_TIMEOUT_MS = 500;

type Navigate = (path: string) => void | Promise<void>;

export const BACK_NAVIGATION_SE_SRC = AUDIO_ASSETS.se.modalCancel;

export function useBackNavigationSe(onNavigate: Navigate) {
  const { isSeEnabled } = useAudioSettings();
  const backNavigationSeRef = useRef<HTMLAudioElement | null>(null);
  const navigationInProgressRef = useRef(false);

  const playBackNavigationSe = useCallback(() => {
    const audio = backNavigationSeRef.current;
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
      timeoutId = window.setTimeout(finish, BACK_NAVIGATION_SE_TIMEOUT_MS);

      void audio.play().catch(finish);
    });
  }, [isSeEnabled]);

  const navigateBackWithSe = useCallback(
    async (path: string) => {
      if (navigationInProgressRef.current) {
        return;
      }

      navigationInProgressRef.current = true;
      try {
        await playBackNavigationSe();
        await onNavigate(path);
      } finally {
        navigationInProgressRef.current = false;
      }
    },
    [onNavigate, playBackNavigationSe],
  );

  return {
    backNavigationSeRef,
    navigateBackWithSe,
  };
}
