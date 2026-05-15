import { useState, useEffect, useCallback } from "react";
import { motion } from "framer-motion";
import { ParticleBackground } from "../shared/ParticleBackground";
import { AnalyzingPanel } from "./AnalyzingPanel";
import { CompletePanel, type AnalysisResult } from "./CompletePanel";

const MOCK_RESULT: AnalysisResult = {
  totalCommits: 23,
  totalPRs: 4,
  totalCP: 43,
  languageBreakdown: [
    { name: "TypeScript", cp: 18, color: "#3178c6", icon: "📘" },
    { name: "Rust", cp: 12, color: "#ff6b35", icon: "🦀" },
    { name: "Python", cp: 8, color: "#3776ab", icon: "🐍" },
    { name: "Go", cp: 5, color: "#00acd7", icon: "🐹" },
  ],
  contributions: [
    {
      repo: "user/awesome-project",
      type: "commit",
      message: "feat: implement OAuth2 login flow",
      language: "TypeScript",
      cp: 1,
      timestamp: "2h ago",
    },
    {
      repo: "user/awesome-project",
      type: "commit",
      message: "refactor: extract auth middleware",
      language: "TypeScript",
      cp: 1,
      timestamp: "3h ago",
    },
    {
      repo: "user/awesome-project",
      type: "commit",
      message: "style: format with prettier",
      language: "TypeScript",
      cp: 1,
      timestamp: "4h ago",
    },
    {
      repo: "user/rust-tool",
      type: "commit",
      message: "fix: handle edge case in parser",
      language: "Rust",
      cp: 1,
      timestamp: "5h ago",
    },
    {
      repo: "user/rust-tool",
      type: "commit",
      message: "docs: add API documentation",
      language: "Rust",
      cp: 1,
      timestamp: "6h ago",
    },
    {
      repo: "user/rust-tool",
      type: "pull_request",
      message: "Add CLI argument parsing",
      language: "Rust",
      cp: 5,
      timestamp: "1d ago",
    },
    {
      repo: "user/data-science",
      type: "commit",
      message: "update: normalize training dataset",
      language: "Python",
      cp: 1,
      timestamp: "1d ago",
    },
    {
      repo: "user/data-science",
      type: "commit",
      message: "feat: add cross-validation split",
      language: "Python",
      cp: 1,
      timestamp: "1d ago",
    },
    {
      repo: "user/go-api",
      type: "commit",
      message: "fix: correct status code on error",
      language: "Go",
      cp: 1,
      timestamp: "2d ago",
    },
    {
      repo: "user/go-api",
      type: "commit",
      message: "feat: add request logging middleware",
      language: "Go",
      cp: 1,
      timestamp: "2d ago",
    },
  ],
};

interface ContributionAnalysisProps {
  onComplete?: () => void;
}

export function ContributionAnalysis({ onComplete }: ContributionAnalysisProps) {
  const [phase, setPhase] = useState<"analyzing" | "complete">("analyzing");
  const [progress, setProgress] = useState(0);
  const [currentMessageIdx, setCurrentMessageIdx] = useState(0);

  useEffect(() => {
    if (phase !== "analyzing") return;

    let alive = true;
    let messageIdx = 0;
    let completeTimer: ReturnType<typeof setTimeout> | null = null;

    const msgTimer = setInterval(() => {
      messageIdx++;
      if (messageIdx < 7) {
        if (alive) setCurrentMessageIdx(messageIdx);
      }
    }, 1200);

    const progTimer = setInterval(() => {
      setProgress((prev) => {
        if (!alive) return prev;
        const next = Math.min(prev + Math.random() * 6 + 2, 100);
        if (next >= 100) {
          clearInterval(msgTimer);
          clearInterval(progTimer);
          completeTimer = setTimeout(() => {
            if (alive) setPhase("complete");
          }, 600);
          return 100;
        }
        return next;
      });
    }, 700);

    return () => {
      alive = false;
      clearInterval(msgTimer);
      clearInterval(progTimer);
      if (completeTimer !== null) clearTimeout(completeTimer);
    };
  }, [phase]);

  const handleContinue = useCallback(() => {
    onComplete?.();
  }, [onComplete]);

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
        <AnalyzingPanel progress={progress} currentMessageIdx={currentMessageIdx} />
      ) : (
        <CompletePanel result={MOCK_RESULT} onContinue={handleContinue} />
      )}
    </div>
  );
}
