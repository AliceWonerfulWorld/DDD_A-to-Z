import { motion } from "framer-motion";
import { steppedEase } from "../../lib/animationUtils";

interface ProfileActionsProps {
  isConfirmingName: boolean;
  isSendingOff: boolean;
  isTransitioning: boolean;
  onBeginJourney: () => void;
  onConfirmNo: () => void;
  onConfirmYes: () => void;
}

export function ProfileActions({
  isConfirmingName,
  isSendingOff,
  isTransitioning,
  onBeginJourney,
  onConfirmNo,
  onConfirmYes,
}: ProfileActionsProps) {
  if (isConfirmingName) {
    return <ConfirmActions onConfirmNo={onConfirmNo} onConfirmYes={onConfirmYes} />;
  }

  const isDisabled = isSendingOff || isTransitioning;

  return (
    <motion.button
      whileHover={isDisabled ? undefined : { scale: 1.02, filter: "brightness(1.1)" }}
      whileTap={
        isDisabled
          ? undefined
          : { scale: 0.98, y: 4, boxShadow: "0px 0px 0 var(--color-gold-dark)" }
      }
      onClick={onBeginJourney}
      disabled={isDisabled}
      animate={
        isTransitioning
          ? {
              backgroundColor: ["var(--color-gold)", "#ffffff", "var(--color-gold)", "#ffffff"],
              boxShadow: [
                "0px 4px 0 var(--color-gold-dark)",
                "0px 0px 0 var(--color-gold-dark)",
                "0px 4px 0 var(--color-gold-dark)",
                "0px 0px 0 var(--color-gold-dark)",
              ],
            }
          : undefined
      }
      transition={{ duration: 0.52, ease: steppedEase(4) }}
      style={{
        marginTop: "1rem",
        width: "100%",
        padding: "1rem",
        fontSize: "1.1rem",
        fontFamily: "var(--font-press)",
        background: "var(--color-gold)",
        color: "#000",
        border: "none",
        boxShadow: "0px 4px 0 var(--color-gold-dark)",
        cursor: isDisabled ? "not-allowed" : "pointer",
        letterSpacing: "0.05em",
        opacity: isDisabled ? 0.75 : 1,
      }}
    >
      {isTransitioning ? "START!" : isSendingOff ? "GOOD LUCK!" : "BEGIN JOURNEY"}
    </motion.button>
  );
}

function ConfirmActions({
  onConfirmNo,
  onConfirmYes,
}: {
  onConfirmNo: () => void;
  onConfirmYes: () => void;
}) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "1fr 1fr",
        gap: "0.75rem",
        width: "100%",
      }}
    >
      <motion.button
        whileHover={{ scale: 1.02, filter: "brightness(1.1)" }}
        whileTap={{ scale: 0.98, y: 4, boxShadow: "0px 0px 0 var(--color-gold-dark)" }}
        onClick={onConfirmYes}
        style={{
          ...actionButtonStyle,
          background: "var(--color-gold)",
          color: "#000",
          border: "none",
          boxShadow: "0px 4px 0 var(--color-gold-dark)",
        }}
      >
        YES
      </motion.button>
      <motion.button
        whileHover={{ scale: 1.02, filter: "brightness(1.1)" }}
        whileTap={{ scale: 0.98, y: 4, boxShadow: "0px 0px 0 rgba(0,0,0,0.8)" }}
        onClick={onConfirmNo}
        style={{
          ...actionButtonStyle,
          background: "var(--color-navy)",
          color: "var(--color-pixel-white)",
          border: "2px solid rgba(255,255,255,0.45)",
          boxShadow: "0px 4px 0 rgba(0,0,0,0.8)",
        }}
      >
        NO
      </motion.button>
    </div>
  );
}

const actionButtonStyle = {
  marginTop: "1rem",
  width: "100%",
  padding: "1rem",
  fontSize: "0.9rem",
  fontFamily: "var(--font-press)",
  cursor: "pointer",
  letterSpacing: "0.05em",
} as const;
