import { useCallback, useRef, useState } from "react";
import { useAudioSettings } from "../features/audio/useAudioSettings";

const loadOnDemand = (audio: HTMLAudioElement) => {
  if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
    audio.load();
  }
};

export function useHomeAudio(onNavigate: (path: string) => void | Promise<void>) {
  const { isSeEnabled } = useAudioSettings();
  const navigationInProgressRef = useRef(false);
  const homeNavSelectSeRef = useRef<HTMLAudioElement | null>(null);
  const confirmModalSeRef = useRef<HTMLAudioElement | null>(null);
  const modalCancelSeRef = useRef<HTMLAudioElement | null>(null);
  const returnTitleSeRef = useRef<HTMLAudioElement | null>(null);
  const gopherTalkSeRef = useRef<HTMLAudioElement | null>(null);
  const pythonPetTalkSeRef = useRef<HTMLAudioElement | null>(null);
  const [audioError, setAudioError] = useState<string | null>(null);

  const playSe = useCallback(
    (audio: HTMLAudioElement | null) => {
      if (!audio || !isSeEnabled) {
        return;
      }

      loadOnDemand(audio);
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

        loadOnDemand(audio);
        audio.currentTime = 0;
        audio.addEventListener("ended", finish, { once: true });
        audio.addEventListener("error", finish, { once: true });
        timeoutId = window.setTimeout(finish, 500);

        void audio.play().catch(finish);
      });
    },
    [isSeEnabled],
  );

  const playModalOpen = useCallback(() => {
    playSe(confirmModalSeRef.current);
  }, [playSe]);

  const playModalCancel = useCallback(() => {
    playSe(modalCancelSeRef.current);
  }, [playSe]);

  const playReturnTitle = useCallback(async () => {
    if (navigationInProgressRef.current) {
      return;
    }

    navigationInProgressRef.current = true;
    try {
      setAudioError(null);
      await playSeUntilEnd(returnTitleSeRef.current);
      await onNavigate("/");
    } catch (error) {
      console.error("failed to return to title from home", error);
      setAudioError("タイトル画面への移動に失敗しました。");
    } finally {
      navigationInProgressRef.current = false;
    }
  }, [onNavigate, playSeUntilEnd]);

  const playGopherTalk = useCallback(() => {
    playSe(gopherTalkSeRef.current);
  }, [playSe]);

  const playPythonPetTalk = useCallback(() => {
    playSe(pythonPetTalkSeRef.current);
  }, [playSe]);

  const playHomeNavSelect = useCallback(
    async (path: string) => {
      if (navigationInProgressRef.current) {
        return;
      }

      navigationInProgressRef.current = true;
      try {
        setAudioError(null);
        await playSeUntilEnd(homeNavSelectSeRef.current);
        await onNavigate(path);
      } catch (error) {
        console.error("failed to navigate from home", error);
        setAudioError("画面移動に失敗しました。");
      } finally {
        navigationInProgressRef.current = false;
      }
    },
    [onNavigate, playSeUntilEnd],
  );

  return {
    audioRefs: {
      homeNavSelectSeRef,
      confirmModalSeRef,
      modalCancelSeRef,
      returnTitleSeRef,
      gopherTalkSeRef,
      pythonPetTalkSeRef,
    },
    audioError,
    isSeEnabled,
    playGopherTalk,
    playPythonPetTalk,
    playHomeNavSelect,
    playModalCancel,
    playModalOpen,
    playReturnTitle,
  };
}
