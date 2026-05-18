import { useState, useEffect, useCallback } from "react";
import { motion } from "framer-motion";
import { ParticleBackground } from "../shared/ParticleBackground";
import { AnalyzingPanel } from "./AnalyzingPanel";
import { CompletePanel, type AnalysisResult } from "./CompletePanel";
import { analyzeContribution } from "../../features/analysis/api";

interface ContributionAnalysisProps {
  onComplete?: () => void;
}

export function ContributionAnalysis({ onComplete }: ContributionAnalysisProps) {
  const [phase, setPhase] = useState<"analyzing" | "complete" | "error">("analyzing");
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState<AnalysisResult | null>(null);

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

    analyzeContribution()
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
    onComplete?.();
  }, [onComplete]);

  const handleRetry = useCallback(() => {
    setProgress(0);
    setPhase("analyzing");
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

      {phase === "analyzing" ? (
        <AnalyzingPanel progress={progress} />
      ) : phase === "error" ? (
        <div
          style={{
            position: "relative",
            zIndex: 2,
            border: "4px solid var(--color-gold)",
            background: "var(--color-navy-light)",
            boxShadow: "0 0 30px rgba(0, 245, 255, 0.12), 8px 8px 0 rgba(0,0,0,0.8)",
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
            Failed to fetch contribution data from GitHub. Make sure you have synced your
            repositories.
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
        <CompletePanel result={result} onContinue={handleContinue} />
      ) : null}
    </div>
  );
}
