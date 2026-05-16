import { motion } from "framer-motion";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";

interface BackButtonProps {
  align?: "left" | "right";
  onNavigate: (path: string) => void;
  targetPath?: string;
}

export function BackButton({ align = "left", onNavigate, targetPath = "/guild" }: BackButtonProps) {
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);

  return (
    <>
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />
      <motion.button
        type="button"
        whileHover={{ y: -2, scale: 1.02 }}
        whileTap={{ y: 2, scale: 0.98 }}
        onClick={() => void navigateBackWithSe(targetPath)}
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(14px, 2.2vw, 28px))",
          left: align === "left" ? "clamp(14px, 2.2vw, 28px)" : "auto",
          right: align === "right" ? "clamp(14px, 2.2vw, 28px)" : "auto",
          zIndex: 6,
          minHeight: "42px",
          border: "2px solid rgba(255, 217, 102, 0.78)",
          borderBottomColor: "rgba(96, 62, 22, 0.95)",
          borderRightColor: "rgba(96, 62, 22, 0.95)",
          background: "rgba(3, 10, 24, 0.74)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.56rem, 1.3vw, 0.78rem)",
          lineHeight: 1.5,
          padding: "10px 12px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        &lt; BACK
      </motion.button>
    </>
  );
}
