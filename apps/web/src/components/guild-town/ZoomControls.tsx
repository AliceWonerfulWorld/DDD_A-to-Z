import { motion } from "framer-motion";
import { ZOOM_STEP } from "./townData";

interface ZoomControlsProps {
  onZoom: (delta: number) => void;
}

export function ZoomControls({ onZoom }: ZoomControlsProps) {
  return (
    <div
      className="absolute right-6 z-[12] flex flex-col gap-3"
      style={{
        bottom: "calc(env(safe-area-inset-bottom, 0px) + 178px)",
      }}
    >
      <ZoomButton ariaLabel="Zoom in" label="+" onClick={() => onZoom(ZOOM_STEP)} />
      <ZoomButton ariaLabel="Zoom out" label="-" onClick={() => onZoom(-ZOOM_STEP)} />
    </div>
  );
}

function ZoomButton({
  ariaLabel,
  label,
  onClick,
}: {
  ariaLabel: string;
  label: string;
  onClick: () => void;
}) {
  return (
    <motion.button
      type="button"
      aria-label={ariaLabel}
      whileHover={{ y: -2, scale: 1.04, backgroundColor: "rgba(255, 217, 102, 0.2)" }}
      whileTap={{ y: 2, scale: 0.96 }}
      onClick={onClick}
      style={{
        width: "46px",
        height: "46px",
        border: "2px solid rgba(255, 217, 102, 0.78)",
        borderBottomColor: "rgba(96, 62, 22, 0.95)",
        borderRightColor: "rgba(96, 62, 22, 0.95)",
        background: "rgba(3, 10, 24, 0.82)",
        boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
        color: "#fff8d7",
        cursor: "pointer",
        fontFamily: "inherit",
        fontSize: "1.1rem",
        lineHeight: 1,
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      {label}
    </motion.button>
  );
}
