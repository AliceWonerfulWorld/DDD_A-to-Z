import { motion, type Variants } from "framer-motion";
import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { GopherSprite } from "./GopherSprite";

interface HomeProps {
  onNavigate: (path: string) => void;
}

const player = {
  name: "DevSamurai",
  title: "Consistency Master",
  level: 18,
  totalCp: 24680,
  todayCp: 320,
};

const navItems = [
  { label: "WAR MAP", caption: "BATTLE FRONT", path: "/war", accent: "#ff5f56" },
  { label: "GUILD BASE", caption: "COMMUNITY HQ", path: "/guild", accent: "#00f5ff" },
  { label: "MY STATUS", caption: "PLAYER DATA", path: "/mypage", accent: "#ffd700" },
];

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

const gopherTalkLines = [
  "今日もコード日和！",
  "ギルド広場を巡回中。",
  "クエスト、行く？",
  "休憩もだいじ。",
  "CP、ためてこ！",
] as const;

const panelVariants: Variants = {
  hidden: { opacity: 0, y: -14 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.35, ease: steppedEase(6) },
  },
};

const navVariants: Variants = {
  hidden: { opacity: 0, y: 18 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: 0.15 + i * 0.08, duration: 0.35, ease: steppedEase(6) },
  }),
};

function HudPanel({ align = "left", children }: { align?: "left" | "right"; children: ReactNode }) {
  return (
    <motion.section
      variants={panelVariants}
      initial="hidden"
      animate="visible"
      style={{
        width: "min(100%, 360px)",
        border: "3px solid rgba(255, 215, 0, 0.72)",
        borderBottomColor: "rgba(111, 79, 28, 0.95)",
        borderRightColor: "rgba(111, 79, 28, 0.95)",
        background: "rgba(3, 10, 24, 0.72)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.78), 8px 8px 0 rgba(0,0,0,0.45)",
        color: "#f4ecd0",
        padding: "14px 16px",
        textAlign: align,
        backdropFilter: "blur(2px)",
      }}
    >
      {children}
    </motion.section>
  );
}

function LabelValue({ label, value }: { label: string; value: string | number }) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "minmax(78px, max-content) 1fr",
        gap: "12px",
        alignItems: "baseline",
        minHeight: "28px",
      }}
    >
      <span
        style={{
          color: "rgba(244, 236, 208, 0.62)",
          fontSize: "0.62rem",
          lineHeight: 1.5,
        }}
      >
        {label}
      </span>
      <span
        style={{
          color: "#fff8d7",
          fontSize: "clamp(0.74rem, 1.6vw, 0.95rem)",
          lineHeight: 1.5,
          overflowWrap: "anywhere",
        }}
      >
        {value}
      </span>
    </div>
  );
}

function TitleButton({ onClick }: { onClick: () => void }) {
  return (
    <motion.button
      type="button"
      variants={panelVariants}
      initial="hidden"
      animate="visible"
      whileHover={{ y: -2, scale: 1.03 }}
      whileTap={{ y: 2, scale: 0.98 }}
      onClick={onClick}
      style={{
        alignSelf: "flex-start",
        border: "2px solid rgba(244, 236, 208, 0.55)",
        borderBottomColor: "rgba(0,0,0,0.78)",
        borderRightColor: "rgba(0,0,0,0.78)",
        background: "rgba(3, 10, 24, 0.68)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 5px 5px 0 rgba(0,0,0,0.38)",
        color: "#fff8d7",
        cursor: "pointer",
        fontFamily: "inherit",
        fontSize: "clamp(0.56rem, 1.4vw, 0.74rem)",
        lineHeight: 1.5,
        padding: "10px 12px",
      }}
    >
      &lt; TITLE
    </motion.button>
  );
}

function ReturnTitleDialog({
  onCancel,
  onConfirm,
}: {
  onCancel: () => void;
  onConfirm: () => void;
}) {
  return (
    <motion.div
      role="presentation"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 10,
        display: "grid",
        placeItems: "center",
        padding: "24px",
        background: "rgba(0, 0, 0, 0.62)",
        backdropFilter: "blur(2px)",
      }}
      onClick={onCancel}
    >
      <motion.div
        role="dialog"
        aria-modal="true"
        aria-labelledby="return-title-dialog-title"
        initial={{ opacity: 0, y: 18, scale: 0.96 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.18, ease: "easeOut" }}
        style={{
          width: "min(100%, 420px)",
          border: "4px solid #ffd700",
          borderBottomColor: "#6f4f1c",
          borderRightColor: "#6f4f1c",
          background: "rgba(3, 10, 24, 0.96)",
          boxShadow: "0 0 0 3px rgba(0,0,0,0.9), 10px 10px 0 rgba(0,0,0,0.58)",
          color: "#fff8d7",
          padding: "22px",
          textAlign: "center",
        }}
        onClick={(event) => event.stopPropagation()}
      >
        <div
          id="return-title-dialog-title"
          style={{
            color: "#ffd700",
            fontSize: "clamp(0.82rem, 3vw, 1.02rem)",
            lineHeight: 1.7,
            marginBottom: "14px",
          }}
        >
          RETURN TO TITLE?
        </div>
        <p
          style={{
            margin: "0 0 22px",
            color: "rgba(255, 248, 215, 0.78)",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "1rem",
            lineHeight: 1.7,
            letterSpacing: "0.04em",
          }}
        >
          タイトル画面に戻りますか？
        </p>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
            gap: "12px",
          }}
        >
          <motion.button
            type="button"
            whileHover={{ scale: 1.03 }}
            whileTap={{ y: 2, scale: 0.98 }}
            onClick={onCancel}
            style={{
              border: "2px solid rgba(255, 255, 255, 0.34)",
              borderBottomColor: "rgba(0,0,0,0.78)",
              borderRightColor: "rgba(0,0,0,0.78)",
              background: "rgba(255, 255, 255, 0.08)",
              boxShadow: "3px 3px 0 rgba(0,0,0,0.5)",
              color: "#fff8d7",
              cursor: "pointer",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.9rem",
              fontWeight: "bold",
              lineHeight: 1.4,
              padding: "11px 12px",
            }}
          >
            CANCEL
          </motion.button>
          <motion.button
            type="button"
            whileHover={{ scale: 1.03 }}
            whileTap={{ y: 2, scale: 0.98 }}
            onClick={onConfirm}
            style={{
              border: "2px solid #fff3a6",
              borderBottomColor: "#6f4f1c",
              borderRightColor: "#6f4f1c",
              background: "#f0c040",
              boxShadow: "3px 3px 0 rgba(0,0,0,0.5)",
              color: "#1b1304",
              cursor: "pointer",
              fontFamily: '"DotGothic16", monospace',
              fontSize: "0.9rem",
              fontWeight: "bold",
              lineHeight: 1.4,
              padding: "11px 12px",
            }}
          >
            TITLE
          </motion.button>
        </div>
      </motion.div>
    </motion.div>
  );
}

function WalkingGopher({ onTalk }: { onTalk: () => void }) {
  const lastXRef = useRef<number | null>(null);
  const talkTimeoutRef = useRef<number | null>(null);
  const [direction, setDirection] = useState<"right" | "left">("right");
  const [talkLine, setTalkLine] = useState<(typeof gopherTalkLines)[number] | null>(null);
  const [reactionCount, setReactionCount] = useState(0);
  const walkRow = direction === "right" ? 1 : 2;
  const speechBubbleSide =
    direction === "right"
      ? { left: "104px", right: "auto" }
      : { left: "auto", right: "96px" };

  useEffect(() => {
    return () => {
      if (talkTimeoutRef.current !== null) {
        window.clearTimeout(talkTimeoutRef.current);
      }
    };
  }, []);

  const reactToClick = () => {
    onTalk();

    const nextLine = gopherTalkLines[reactionCount % gopherTalkLines.length];
    setReactionCount((current) => current + 1);
    setTalkLine(nextLine);

    if (talkTimeoutRef.current !== null) {
      window.clearTimeout(talkTimeoutRef.current);
    }

    talkTimeoutRef.current = window.setTimeout(() => {
      setTalkLine(null);
      talkTimeoutRef.current = null;
    }, 2600);
  };

  return (
    <motion.button
      type="button"
      aria-label="Gopher君に話しかける"
      initial={false}
      animate={{
        x: ["16vw", "30vw", "24vw", "50vw", "62vw", "46vw", "28vw", "16vw"],
        y: ["0px", "-10px", "6px", "-14px", "2px", "16px", "8px", "0px"],
        scale: [0.92, 0.88, 0.94, 0.86, 0.9, 0.96, 0.94, 0.92],
      }}
      transition={{
        duration: 26,
        repeat: Infinity,
        ease: "easeInOut",
        times: [0, 0.16, 0.28, 0.44, 0.58, 0.72, 0.88, 1],
      }}
      onUpdate={(latest) => {
        const currentX =
          typeof latest.x === "number" ? latest.x : Number.parseFloat(String(latest.x));
        const lastX = lastXRef.current;

        if (!Number.isFinite(currentX)) {
          return;
        }

        if (lastX !== null) {
          const deltaX = currentX - lastX;
          if (Math.abs(deltaX) > 0.02) {
            setDirection(deltaX > 0 ? "right" : "left");
          }
        }

        lastXRef.current = currentX;
      }}
      style={{
        position: "absolute",
        left: 0,
        bottom: "clamp(6px, 2vh, 18px)",
        width: "92px",
        height: "100px",
        border: 0,
        background: "transparent",
        cursor: "pointer",
        font: "inherit",
        padding: 0,
        zIndex: 4,
      }}
      onClick={reactToClick}
    >
      {talkLine && (
        <motion.div
          key={talkLine}
          initial={{ opacity: 0, y: 8, scale: 0.94 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: 8, scale: 0.94 }}
          style={{
            position: "absolute",
            ...speechBubbleSide,
            bottom: "46px",
            width: "max-content",
            maxWidth: "min(220px, 34vw)",
            border: "2px solid rgba(255, 215, 0, 0.82)",
            borderBottomColor: "rgba(111, 79, 28, 0.95)",
            borderRightColor: "rgba(111, 79, 28, 0.95)",
            background: "rgba(3, 10, 24, 0.9)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.72), 4px 4px 0 rgba(0,0,0,0.42)",
            color: "#fff8d7",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.78rem",
            lineHeight: 1.5,
            letterSpacing: "0.04em",
            padding: "8px 10px",
            textAlign: "center",
            whiteSpace: "normal",
            pointerEvents: "none",
            zIndex: 6,
          }}
        >
          {talkLine}
        </motion.div>
      )}
      <motion.div
        animate={{ y: [0, -2, 0] }}
        transition={{ duration: 0.45, repeat: Infinity, ease: steppedEase(3) }}
        style={{
          position: "absolute",
          left: 0,
          bottom: 8,
          width: "132px",
          height: "143px",
          transform: "scale(0.62)",
          transformOrigin: "left bottom",
        }}
      >
        <motion.div
          key={reactionCount}
          animate={
            reactionCount === 0
              ? {}
              : {
                  y: [0, -20, 0],
                  rotate: [0, -5, 5, 0],
                }
          }
          transition={{ duration: 0.42, ease: steppedEase(5) }}
        >
          <GopherSprite frameCount={8} row={walkRow} />
        </motion.div>
      </motion.div>
      <motion.div
        animate={{ scaleX: [1, 0.86, 1], opacity: [0.34, 0.24, 0.34] }}
        transition={{ duration: 0.45, repeat: Infinity, ease: steppedEase(3) }}
        style={{
          position: "absolute",
          left: "18px",
          bottom: "2px",
          width: "62px",
          height: "10px",
          background: "rgba(0,0,0,0.42)",
          filter: "blur(1px)",
        }}
      />
    </motion.button>
  );
}

export function Home({ onNavigate }: HomeProps) {
  const [isReturnTitleDialogOpen, setIsReturnTitleDialogOpen] = useState(false);
  const homeBgmRef = useRef<HTMLAudioElement | null>(null);
  const confirmModalSeRef = useRef<HTMLAudioElement | null>(null);
  const modalCancelSeRef = useRef<HTMLAudioElement | null>(null);
  const returnTitleSeRef = useRef<HTMLAudioElement | null>(null);
  const gopherTalkSeRef = useRef<HTMLAudioElement | null>(null);

  useEffect(() => {
    const audio = homeBgmRef.current;
    if (!audio) {
      return;
    }

    let isUnlocked = false;
    audio.volume = 0.34;

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

  const playSe = useCallback((audio: HTMLAudioElement | null) => {
    if (!audio) {
      return;
    }

    audio.currentTime = 0;
    void audio.play().catch(() => {
      // Browser autoplay restrictions can still block sound in unusual navigation paths.
    });
  }, []);

  const playSeUntilEnd = useCallback((audio: HTMLAudioElement | null) => {
    if (!audio) {
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
      timeoutId = window.setTimeout(finish, 500);

      void audio.play().catch(finish);
    });
  }, []);

  const cancelReturnTitle = () => {
    playSe(modalCancelSeRef.current);
    setIsReturnTitleDialogOpen(false);
  };

  const openReturnTitleDialog = () => {
    playSe(confirmModalSeRef.current);
    setIsReturnTitleDialogOpen(true);
  };

  const confirmReturnTitle = async () => {
    await playSeUntilEnd(returnTitleSeRef.current);
    onNavigate("/");
  };

  const playGopherTalk = () => {
    playSe(gopherTalkSeRef.current);
  };

  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        backgroundImage:
          "linear-gradient(180deg, rgba(4, 8, 18, 0.25), rgba(4, 8, 18, 0.56)), url('/home_lunch.png')",
        backgroundSize: "cover",
        backgroundPosition: "center center",
        backgroundRepeat: "no-repeat",
        imageRendering: "pixelated",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#f4ecd0",
      }}
    >
      <audio
        ref={homeBgmRef}
        src="/bgm/home_bgm.ogg"
        loop
        preload="auto"
        aria-hidden="true"
      />
      <audio
        ref={confirmModalSeRef}
        src="/SE/confirm-modal.wav"
        preload="auto"
        aria-hidden="true"
      />
      <audio
        ref={modalCancelSeRef}
        src="/SE/modal-cancel.wav"
        preload="auto"
        aria-hidden="true"
      />
      <audio
        ref={returnTitleSeRef}
        src="/SE/return-title.wav"
        preload="auto"
        aria-hidden="true"
      />
      <audio
        ref={gopherTalkSeRef}
        src="/SE/gopher-talk.wav"
        preload="auto"
        aria-hidden="true"
      />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.1), rgba(0,0,0,0.1) 1px, transparent 1px, transparent 4px)",
          mixBlendMode: "multiply",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          boxShadow: "inset 0 0 90px rgba(0,0,0,0.7)",
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <div
        style={{
          position: "relative",
          zIndex: 3,
          minHeight: "100svh",
          display: "grid",
          gridTemplateRows: "auto 1fr auto",
          padding: "clamp(16px, 3vw, 32px)",
          gap: "20px",
        }}
      >
        <header
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
            gap: "18px",
            flexWrap: "wrap",
          }}
        >
          <div
            style={{
              display: "grid",
              gap: "12px",
              width: "min(100%, 360px)",
            }}
          >
            <HudPanel>
              <div
                style={{
                  color: "#ffd700",
                  fontSize: "0.64rem",
                  lineHeight: 1.6,
                  marginBottom: "10px",
                }}
              >
                PLAYER INFO
              </div>
              <div style={{ display: "grid", gap: "6px" }}>
                <LabelValue label="NAME" value={player.name} />
                <LabelValue label="TITLE" value={player.title} />
                <LabelValue label="LEVEL" value={`LV.${player.level}`} />
              </div>
            </HudPanel>

            <TitleButton onClick={openReturnTitleDialog} />
          </div>

          <HudPanel align="right">
            <div
              style={{
                color: "#00f5ff",
                fontSize: "0.64rem",
                lineHeight: 1.6,
                marginBottom: "10px",
              }}
            >
              CONTRIBUTION POINT
            </div>
            <div style={{ display: "grid", gap: "6px" }}>
              <LabelValue label="TOTAL CP" value={player.totalCp.toLocaleString()} />
              <LabelValue label="TODAY CP" value={`+${player.todayCp.toLocaleString()}`} />
            </div>
          </HudPanel>
        </header>

        <section
          aria-label="Character placement area"
          style={{
            position: "relative",
            minHeight: "clamp(220px, 42vh, 520px)",
            overflow: "visible",
          }}
        >
          <WalkingGopher onTalk={playGopherTalk} />
        </section>

        <nav
          aria-label="Main navigation"
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(min(100%, 220px), 1fr))",
            gap: "14px",
            width: "min(100%, 900px)",
            justifySelf: "center",
          }}
        >
          {navItems.map((item, index) => (
            <motion.button
              key={item.label}
              custom={index}
              variants={navVariants}
              initial="hidden"
              animate="visible"
              whileHover={{ y: -4, scale: 1.02 }}
              whileTap={{ y: 2, scale: 0.98 }}
              onClick={() => onNavigate(item.path)}
              style={{
                minHeight: "86px",
                border: `3px solid ${item.accent}`,
                borderBottomColor: "rgba(0,0,0,0.75)",
                borderRightColor: "rgba(0,0,0,0.75)",
                background: "rgba(4, 11, 27, 0.82)",
                boxShadow: `0 0 0 2px rgba(0,0,0,0.78), 6px 6px 0 rgba(0,0,0,0.5), inset 0 0 18px ${item.accent}22`,
                color: "#fff8d7",
                cursor: "pointer",
                fontFamily: "inherit",
                padding: "16px 18px",
                textAlign: "left",
              }}
            >
              <span
                style={{
                  display: "block",
                  color: item.accent,
                  fontSize: "clamp(0.88rem, 2vw, 1.08rem)",
                  lineHeight: 1.5,
                  marginBottom: "8px",
                  overflowWrap: "anywhere",
                }}
              >
                {item.label}
              </span>
              <span
                style={{
                  display: "block",
                  color: "rgba(244, 236, 208, 0.66)",
                  fontSize: "0.6rem",
                  lineHeight: 1.5,
                }}
              >
                {item.caption}
              </span>
            </motion.button>
          ))}
        </nav>
      </div>

      {isReturnTitleDialogOpen && (
        <ReturnTitleDialog
          onCancel={cancelReturnTitle}
          onConfirm={confirmReturnTitle}
        />
      )}
    </main>
  );
}
