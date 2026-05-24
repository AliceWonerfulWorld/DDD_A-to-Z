import { useCallback, useEffect, useRef, useState } from "react";
import type { MutableRefObject } from "react";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { fetchMe } from "../../features/auth/api";

const JOURNEY_START_DELAY_MS = 1800;
const INTRO_TEXT = "歓迎しよう、新たな挑戦者よ。\n君のコードネームを教えてくれ。";
const EMPTY_NAME_WARNING_TEXT = "むむっ、名前がないぞ！\nコードネームを入力してくれ。";
const buildNameConfirmText = (name: string) =>
  `そのコードネームは「${name}」でよいのだな？\nよければ、旅立ちの合図をくれ。`;
const buildSendoffText = (name: string) =>
  `よし、「${name}」よ。\nコードの世界へ、いってらっしゃい！`;

interface InitialProfileFlowArgs {
  onComplete: (username: string, avatarUrl: string) => void;
}

function getAudioContext(): typeof window.AudioContext | undefined {
  return (
    window.AudioContext ||
    (window as unknown as { webkitAudioContext?: typeof window.AudioContext }).webkitAudioContext
  );
}

export function useInitialProfileFlow({ onComplete }: InitialProfileFlowArgs) {
  const [username, setUsername] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");
  const [githubAvatarUrl, setGithubAvatarUrl] = useState("");
  const [introText, setIntroText] = useState("");
  const [angryText, setAngryText] = useState("");
  const [confirmText, setConfirmText] = useState("");
  const [sendoffText, setSendoffText] = useState("");
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [isGopherAngry, setIsGopherAngry] = useState(false);
  const [isConfirmingName, setIsConfirmingName] = useState(false);
  const [isSendingOff, setIsSendingOff] = useState(false);

  const audioCtxRef = useRef<AudioContext | null>(null);
  const transitionTimeoutRef = useRef<number | null>(null);
  const angryTimeoutRef = useRef<number | null>(null);
  const introSpeechIntervalRef = useRef<number | null>(null);
  const angrySpeechIntervalRef = useRef<number | null>(null);
  const confirmSpeechIntervalRef = useRef<number | null>(null);
  const sendoffSpeechIntervalRef = useRef<number | null>(null);

  const clearTimer = (timerRef: MutableRefObject<number | null>) => {
    if (timerRef.current !== null) {
      window.clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  };

  const clearIntervalRef = (intervalRef: MutableRefObject<number | null>) => {
    if (intervalRef.current !== null) {
      window.clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  useEffect(() => {
    try {
      const AudioContext = getAudioContext();
      if (AudioContext) {
        audioCtxRef.current = new AudioContext();
      }
    } catch {
      // AudioContext not supported.
    }

    return () => {
      void audioCtxRef.current?.close();
      audioCtxRef.current = null;
      clearTimer(transitionTimeoutRef);
      clearTimer(angryTimeoutRef);
      clearIntervalRef(introSpeechIntervalRef);
      clearIntervalRef(angrySpeechIntervalRef);
      clearIntervalRef(confirmSpeechIntervalRef);
      clearIntervalRef(sendoffSpeechIntervalRef);
    };
  }, []);

  const playTone = useCallback(
    (options: {
      type: OscillatorType;
      frequency: number;
      gain: number;
      duration: number;
      endFrequency?: number;
    }) => {
      void (async () => {
        const ctx = audioCtxRef.current;
        if (!ctx) return;
        if (ctx.state === "suspended") {
          await ctx.resume();
        }
        if (ctx.state !== "running") return;

        const osc = ctx.createOscillator();
        const gain = ctx.createGain();
        osc.type = options.type;
        osc.frequency.setValueAtTime(options.frequency, ctx.currentTime);
        if (options.endFrequency !== undefined) {
          osc.frequency.exponentialRampToValueAtTime(
            options.endFrequency,
            ctx.currentTime + options.duration,
          );
        }
        gain.gain.setValueAtTime(options.gain, ctx.currentTime);
        gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + options.duration);
        osc.connect(gain);
        gain.connect(ctx.destination);
        osc.start();
        osc.stop(ctx.currentTime + options.duration);
      })().catch(() => {
        // Browser autoplay restrictions may still block sound until the next user gesture.
      });
    },
    [],
  );

  const playSpeechBeep = useCallback(() => {
    playTone({
      type: "square",
      frequency: 700 + Math.random() * 100,
      gain: 0.15,
      duration: 0.06,
    });
  }, [playTone]);

  const playNameInputBeep = useCallback(
    (isDeleting: boolean) => {
      playTone({
        type: "square",
        frequency: isDeleting ? 420 : 860,
        gain: 0.12,
        duration: 0.045,
      });
    },
    [playTone],
  );

  const playRejectBeep = useCallback(() => {
    playTone({
      type: "sawtooth",
      frequency: 180,
      endFrequency: 90,
      gain: 0.12,
      duration: 0.16,
    });
  }, [playTone]);

  const playSpeech = useCallback(
    (args: {
      text: string;
      setText: (text: string) => void;
      intervalRef: MutableRefObject<number | null>;
      speedMs: number;
      onDone?: () => void;
    }) => {
      clearIntervalRef(args.intervalRef);
      args.setText("");
      let i = 0;
      args.intervalRef.current = window.setInterval(() => {
        if (i <= args.text.length) {
          args.setText(args.text.slice(0, i));
          if (i < args.text.length && args.text[i] !== " " && args.text[i] !== "\n") {
            playSpeechBeep();
          }
          i++;
        } else {
          clearIntervalRef(args.intervalRef);
          args.onDone?.();
        }
      }, args.speedMs);
    },
    [playSpeechBeep],
  );

  const resetAngryState = useCallback(() => {
    setIsGopherAngry(false);
    setAngryText("");
    clearTimer(angryTimeoutRef);
    clearIntervalRef(angrySpeechIntervalRef);
  }, []);

  const resetConfirmState = useCallback(() => {
    setIsConfirmingName(false);
    setConfirmText("");
    clearIntervalRef(confirmSpeechIntervalRef);
  }, []);

  const resetSendoffState = useCallback(() => {
    setIsSendingOff(false);
    setSendoffText("");
    clearIntervalRef(sendoffSpeechIntervalRef);
  }, []);

  useEffect(() => {
    playSpeech({
      text: INTRO_TEXT,
      setText: setIntroText,
      intervalRef: introSpeechIntervalRef,
      speedMs: 90,
    });

    void fetchMe().then((user) => {
      if (user) {
        setUsername(user.username || "octocat");
        const ghUrl = user.avatar_url || "";
        setGithubAvatarUrl(ghUrl);
        setAvatarUrl(ghUrl);
      }
    });
  }, [playSpeech]);

  const handleUsernameChange = useCallback(
    (nextUsername: string) => {
      if (nextUsername !== username) {
        playNameInputBeep(nextUsername.length < username.length);
      }
      if (nextUsername.trim()) {
        resetAngryState();
      }
      resetConfirmState();
      resetSendoffState();
      setUsername(nextUsername);
    },
    [playNameInputBeep, resetAngryState, resetConfirmState, resetSendoffState, username],
  );

  const startJourney = useCallback(
    (trimmedUsername: string) => {
      setIsTransitioning(true);
      setIsSendingOff(false);

      const titleStartAudio = new Audio(AUDIO_ASSETS.se.titleStart);
      void titleStartAudio.play().catch(() => {
        // Browser audio policies can still reject playback in some environments.
      });

      transitionTimeoutRef.current = window.setTimeout(() => {
        onComplete(trimmedUsername, avatarUrl);
      }, JOURNEY_START_DELAY_MS);
    },
    [onComplete, avatarUrl],
  );

  const handleBeginJourney = useCallback(() => {
    const trimmedUsername = username.trim();
    if (isSendingOff || isTransitioning) return;

    if (!trimmedUsername) {
      setIsGopherAngry(true);
      playSpeech({
        text: EMPTY_NAME_WARNING_TEXT,
        setText: setAngryText,
        intervalRef: angrySpeechIntervalRef,
        speedMs: 82,
      });
      playRejectBeep();
      clearTimer(angryTimeoutRef);
      angryTimeoutRef.current = window.setTimeout(() => {
        resetAngryState();
      }, 3000);
      return;
    }

    resetAngryState();
    setIsConfirmingName(true);
    playSpeech({
      text: buildNameConfirmText(trimmedUsername),
      setText: setConfirmText,
      intervalRef: confirmSpeechIntervalRef,
      speedMs: 82,
    });
  }, [isSendingOff, isTransitioning, playRejectBeep, playSpeech, resetAngryState, username]);

  const handleConfirmNo = useCallback(() => {
    resetConfirmState();
  }, [resetConfirmState]);

  const handleConfirmYes = useCallback(() => {
    const trimmedUsername = username.trim();
    if (!trimmedUsername || isSendingOff || isTransitioning) return;

    resetConfirmState();
    setIsSendingOff(true);
    playSpeech({
      text: buildSendoffText(trimmedUsername),
      setText: setSendoffText,
      intervalRef: sendoffSpeechIntervalRef,
      speedMs: 82,
      onDone: () => startJourney(trimmedUsername),
    });
  }, [isSendingOff, isTransitioning, playSpeech, resetConfirmState, startJourney, username]);

  const dialogueText = isConfirmingName
    ? confirmText
    : isSendingOff
      ? sendoffText
      : isGopherAngry
        ? angryText
        : introText;

  return {
    dialogueText,
    handleBeginJourney,
    handleConfirmNo,
    handleConfirmYes,
    handleUsernameChange,
    isConfirmingName,
    isGopherAngry,
    isSendingOff,
    isTransitioning,
    username,
    avatarUrl,
    githubAvatarUrl,
    setAvatarUrl,
  };
}
