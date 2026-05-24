import { useState, useEffect, useCallback, useRef } from "react";
import { motion } from "framer-motion";
import { ParticleBackground } from "../shared/ParticleBackground";
import { AnalyzingPanel } from "./AnalyzingPanel";
import { CompletePanel, type AnalysisResult } from "./CompletePanel";
import { analyzeContribution } from "../../features/analysis/api";
import { NadesikoGimmick } from "../analysis/NadesikoGimmick";

interface ContributionAnalysisProps {
  onComplete?: () => void;
}

const ANALYSIS_BOOT_DELAY_MS = 2800;
const BOOT_MESSAGES = [
  "INSERTING QUEST LOG",
  "CHECKING GITHUB RELIC",
  "WARMING SCAN CRYSTAL",
  "READY",
];
const ARRIVAL_TILES = Array.from({ length: 48 }, (_, i) => ({
  col: i % 8,
  row: Math.floor(i / 8),
}));
const DEPARTURE_DELAY_MS = 1900;
const DEPARTURE_SPARKS = [
  { x: 12, y: 18, size: 6, delay: 0.08, color: "var(--color-neon-cyan)" },
  { x: 20, y: 74, size: 4, delay: 0.34, color: "var(--color-gold)" },
  { x: 28, y: 38, size: 5, delay: 0.18, color: "var(--color-pixel-white)" },
  { x: 36, y: 58, size: 7, delay: 0.3, color: "var(--color-neon-cyan)" },
  { x: 44, y: 22, size: 4, delay: 0.12, color: "var(--color-gold)" },
  { x: 52, y: 70, size: 5, delay: 0.42, color: "var(--color-pixel-white)" },
  { x: 60, y: 34, size: 6, delay: 0.24, color: "var(--color-neon-cyan)" },
  { x: 68, y: 64, size: 4, delay: 0.28, color: "var(--color-gold)" },
  { x: 76, y: 28, size: 5, delay: 0.16, color: "var(--color-pixel-white)" },
  { x: 84, y: 80, size: 6, delay: 0.38, color: "var(--color-neon-cyan)" },
] as const;

const steppedEase = (steps: number) => (t: number) =>
  Math.floor(t * steps) / steps;

function AnalysisArrivalWipe() {
  return (
    <motion.div
      initial={{ opacity: 1 }}
      animate={{ opacity: 0 }}
      transition={{ delay: 1.05, duration: 0.12, ease: steppedEase(2) }}
      className="fixed inset-0 z-[9999] pointer-events-none overflow-hidden"
    >
      <motion.div
        initial={{ opacity: 1 }}
        animate={{ opacity: [1, 1, 0] }}
        transition={{ duration: 0.72, ease: steppedEase(5) }}
        className="absolute inset-0 bg-white"
      />

      <div
        aria-hidden="true"
        style={{
          position: "absolute",
          inset: 0,
          display: "grid",
          gridTemplateColumns: "repeat(8, 1fr)",
          gridTemplateRows: "repeat(6, 1fr)",
        }}
      >
        {ARRIVAL_TILES.map((tile, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 1 }}
            animate={{ opacity: 0 }}
            transition={{
              delay: 0.18 + (tile.row + Math.abs(tile.col - 3.5)) * 0.045,
              duration: 0.08,
              ease: steppedEase(1),
            }}
            style={{
              background:
                i % 3 === 0
                  ? "var(--color-neon-cyan)"
                  : "var(--color-pixel-white)",
            }}
          />
        ))}
      </div>

      <motion.div
        initial={{ opacity: 0, scaleX: 0 }}
        animate={{ opacity: [0, 1, 0], scaleX: [0, 1, 0] }}
        transition={{ delay: 0.72, duration: 0.42, ease: steppedEase(5) }}
        style={{
          position: "absolute",
          top: "50%",
          left: 0,
          right: 0,
          height: "18px",
          background: "var(--color-pixel-white)",
          boxShadow:
            "0 18px 0 var(--color-neon-cyan), 0 -18px 0 var(--color-gold)",
          transform: "translateY(-50%)",
        }}
      />
    </motion.div>
  );
}

function AnalysisBootPanel() {
  return (
    <motion.div
      initial={{ scaleY: 0, opacity: 0 }}
      animate={{ scaleY: 1, opacity: 1 }}
      transition={{ duration: 0.35, ease: steppedEase(6) }}
      style={{
        position: "relative",
        zIndex: 2,
        border: "4px solid var(--color-gold)",
        background: "var(--color-navy-light)",
        boxShadow:
          "0 0 30px rgba(0, 245, 255, 0.12), 8px 8px 0 rgba(0,0,0,0.8)",
        maxWidth: "560px",
        width: "100%",
        padding: "2rem",
        display: "flex",
        flexDirection: "column",
        gap: "1.5rem",
      }}
    >
      <div style={{ textAlign: "center" }}>
        <motion.span
          animate={{ opacity: [1, 0, 1] }}
          transition={{
            duration: 0.32,
            repeat: Infinity,
            ease: steppedEase(2),
          }}
          style={{
            display: "inline-block",
            fontFamily: "var(--font-press)",
            fontSize: "0.85rem",
            color: "var(--color-gold)",
            letterSpacing: "0.08em",
          }}
        >
          ▶ ANALYSIS TERMINAL
        </motion.span>
      </div>

      <div
        style={{
          position: "relative",
          height: "150px",
          border: "2px solid rgba(255, 211, 107, 0.65)",
          background: "rgba(0,0,0,0.48)",
          overflow: "hidden",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <motion.div
          initial={{ y: "-100%" }}
          animate={{ y: ["-100%", "120%"] }}
          transition={{
            duration: 1.15,
            repeat: Infinity,
            ease: steppedEase(8),
          }}
          style={{
            position: "absolute",
            inset: "0 0 auto 0",
            height: "36px",
            background:
              "repeating-linear-gradient(0deg, rgba(0,245,255,0.24) 0 4px, rgba(255,255,255,0.12) 4px 8px)",
          }}
        />

        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: [0, 1, 1, 0.7], scale: [0.8, 1, 1.08, 1] }}
          transition={{ duration: 1.2, ease: steppedEase(5) }}
          style={{
            width: "92px",
            height: "92px",
            border: "4px solid var(--color-neon-cyan)",
            boxShadow: "8px 8px 0 rgba(0,0,0,0.8)",
            background:
              "linear-gradient(45deg, rgba(0,245,255,0.16) 25%, transparent 25%, transparent 50%, rgba(0,245,255,0.16) 50%, rgba(0,245,255,0.16) 75%, transparent 75%)",
            backgroundSize: "16px 16px",
          }}
        />
      </div>

      <div style={{ minHeight: "5.8rem" }}>
        {BOOT_MESSAGES.map((message, i) => (
          <motion.div
            key={message}
            initial={{ opacity: 0, x: -8 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{
              delay: 0.85 + i * 0.36,
              duration: 0.18,
              ease: steppedEase(4),
            }}
            style={{
              fontFamily: "var(--font-dot)",
              fontSize: "0.8rem",
              color:
                i === BOOT_MESSAGES.length - 1
                  ? "var(--color-neon-green)"
                  : "var(--color-pixel-white)",
              padding: "3px 0",
            }}
          >
            &gt; {message}
            <motion.span
              animate={{ opacity: [0, 1, 0] }}
              transition={{
                duration: 0.35,
                repeat: Infinity,
                ease: steppedEase(2),
              }}
            >
              _
            </motion.span>
          </motion.div>
        ))}
      </div>
    </motion.div>
  );
}

function AnalysisDepartureWipe() {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.08, ease: steppedEase(1) }}
      className="fixed inset-0 z-[9999] pointer-events-none overflow-hidden"
    >
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: [0, 0.75, 0.95] }}
        transition={{ duration: 0.75, ease: steppedEase(5) }}
        className="absolute inset-0"
        style={{
          background:
            "radial-gradient(circle at 50% 50%, rgba(0,245,255,0.18), rgba(13,25,54,0.9) 42%, #050510 100%)",
        }}
      />

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: [0, 0.5, 0.2] }}
        transition={{ duration: 0.85, ease: steppedEase(5) }}
        className="absolute inset-0"
        style={{
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 7px, rgba(255,255,255,0.06) 7px, rgba(255,255,255,0.06) 8px)",
        }}
      />

      <motion.div
        initial={{ scaleX: 0, opacity: 0 }}
        animate={{ scaleX: [0, 1, 0.32], opacity: [0, 1, 0] }}
        transition={{ duration: 0.82, ease: steppedEase(7) }}
        style={{
          position: "absolute",
          top: "50%",
          left: "8vw",
          right: "8vw",
          height: "8px",
          background: "rgba(255,255,255,0.95)",
          boxShadow:
            "0 0 18px rgba(255,255,255,0.65), 0 14px 0 rgba(0,245,255,0.5), 0 -14px 0 rgba(255,211,107,0.45)",
          transform: "translateY(-50%)",
          transformOrigin: "center",
        }}
      />

      {DEPARTURE_SPARKS.map((spark, i) => (
        <motion.div
          key={i}
          aria-hidden="true"
          initial={{ opacity: 0, x: "-18vw", scale: 0.8 }}
          animate={{ opacity: [0, 1, 0], x: "18vw", scale: [0.8, 1, 0.8] }}
          transition={{
            delay: spark.delay,
            duration: 0.72,
            ease: steppedEase(6),
          }}
          style={{
            position: "absolute",
            left: `${spark.x}%`,
            top: `${spark.y}%`,
            width: `${spark.size}px`,
            height: `${spark.size}px`,
            background: spark.color,
            boxShadow: `0 0 12px ${spark.color}`,
          }}
        />
      ))}

      <motion.div
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: [0, 1, 1, 0], y: [12, 0, 0, -4] }}
        transition={{ delay: 0.38, duration: 1.08, ease: steppedEase(6) }}
        style={{
          position: "absolute",
          left: 0,
          right: 0,
          top: "calc(50% + 28px)",
          padding: "0 1rem",
          fontFamily: "var(--font-press)",
          fontSize: "clamp(0.72rem, 2vw, 1.1rem)",
          color: "var(--color-pixel-white)",
          letterSpacing: "0.08em",
          textAlign: "center",
          textShadow: "0 0 12px rgba(0,245,255,0.8), 3px 3px 0 rgba(0,0,0,0.8)",
          whiteSpace: "nowrap",
        }}
      >
        OPENING HOME GATE
      </motion.div>
    </motion.div>
  );
}

export function ContributionAnalysis({
  onComplete,
}: ContributionAnalysisProps) {
  const [phase, setPhase] = useState<
    "booting" | "analyzing" | "complete" | "departing" | "error"
  >("booting");
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState<AnalysisResult | null>(null);
  const promiseRef = useRef<Promise<AnalysisResult> | null>(null);
  const departureTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    return () => {
      if (departureTimerRef.current !== null) {
        clearTimeout(departureTimerRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (phase !== "booting") return;

    const bootTimer = setTimeout(() => {
      setPhase("analyzing");
    }, ANALYSIS_BOOT_DELAY_MS);

    return () => clearTimeout(bootTimer);
  }, [phase]);

  useEffect(() => {
    if (phase !== "analyzing") return;

    let alive = true;
    let completeTimer: ReturnType<typeof setTimeout> | null = null;

    const advanceProgress = () => {
      setProgress((prev) => {
        if (!alive) return prev;
        const next = prev + Math.random() * 4 + 1;
        return Math.min(next, 99);
      });
    };

    const progTimer = setInterval(advanceProgress, 1400);

    if (!promiseRef.current) {
      promiseRef.current = analyzeContribution();
    }
    const promise = promiseRef.current;

    promise
      .then((data) => {
        if (!alive) return;
        setResult(data);
        clearInterval(progTimer);
        setProgress(100);
        completeTimer = setTimeout(() => {
          if (alive) setPhase("complete");
        }, 600);
      })
      .catch(() => {
        if (!alive) return;
        clearInterval(progTimer);
        if (alive) setPhase("error");
      });

    return () => {
      alive = false;
      clearInterval(progTimer);
      if (completeTimer !== null) clearTimeout(completeTimer);
    };
  }, [phase]);

  const handleContinue = useCallback(() => {
    if (phase === "departing") return;

    setPhase("departing");
    departureTimerRef.current = setTimeout(() => {
      onComplete?.();
    }, DEPARTURE_DELAY_MS);
  }, [onComplete, phase]);

  const handleRetry = useCallback(() => {
    setProgress(0);
    setPhase("booting");
    promiseRef.current = null;
    if (departureTimerRef.current !== null) {
      clearTimeout(departureTimerRef.current);
      departureTimerRef.current = null;
    }
  }, []);

  return (
    <div
      className="flex flex-col items-center justify-center min-h-svh p-4 relative overflow-hidden"
      style={{
        background:
          "radial-gradient(circle at 50% 120%, var(--color-navy-mid) 0%, var(--color-navy) 60%, #050510 100%)",
        fontFamily: "var(--font-dot)",
        color: "var(--color-pixel-white)",
      }}
    >
      <ParticleBackground />

      <motion.div
        animate={{ opacity: [0.2, 0.5, 0.2] }}
        transition={{ duration: 5, repeat: Infinity, ease: "easeInOut" }}
        style={{
          position: "absolute",
          bottom: "-15vh",
          left: "50%",
          transform: "translateX(-50%)",
          width: "120vw",
          height: "50vh",
          background:
            "radial-gradient(ellipse at center, rgba(0, 245, 255, 0.1) 0%, transparent 70%)",
          pointerEvents: "none",
          zIndex: 0,
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "absolute",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0,0,0,0.04) 2px, rgba(0,0,0,0.04) 4px)",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      {phase === "booting" ? (
        <>
          <AnalysisBootPanel />
          <AnalysisArrivalWipe />
        </>
      ) : phase === "analyzing" ? (
        <AnalyzingPanel progress={progress} />
      ) : phase === "error" ? (
        <div
          style={{
            position: "relative",
            zIndex: 2,
            border: "4px solid var(--color-gold)",
            background: "var(--color-navy-light)",
            boxShadow:
              "0 0 30px rgba(0, 245, 255, 0.12), 8px 8px 0 rgba(0,0,0,0.8)",
            maxWidth: "520px",
            width: "100%",
            padding: "2rem",
            display: "flex",
            flexDirection: "column",
            gap: "1.5rem",
            textAlign: "center",
          }}
        >
          <span
            style={{
              fontFamily: "var(--font-press)",
              fontSize: "1rem",
              color: "var(--color-gold)",
            }}
          >
            ✗ ANALYSIS FAILED
          </span>
          <span style={{ fontSize: "0.8rem", color: "rgba(232,232,208,0.6)" }}>
            Failed to fetch contribution data from GitHub. Make sure you have
            synced your repositories.
          </span>
          <button
            onClick={handleRetry}
            style={{
              padding: "0.8rem 2rem",
              fontFamily: "var(--font-press)",
              fontSize: "0.85rem",
              background: "var(--color-gold)",
              color: "#000",
              border: "none",
              boxShadow: "0px 4px 0 var(--color-gold-dark)",
              cursor: "pointer",
            }}
          >
            RETRY
          </button>
        </div>
      ) : result ? (
        <>
          <CompletePanel
            isDeparting={phase === "departing"}
            result={result}
            onContinue={handleContinue}
          />
          {phase === "departing" && <AnalysisDepartureWipe />}
        </>
      ) : null}
      <NadesikoGimmick />
    </div>
  );
}
