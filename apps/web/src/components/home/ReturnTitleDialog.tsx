import { motion } from "framer-motion";
import { useEffect, useRef } from "react";

const getTabbableElements = (element: HTMLElement) =>
  Array.from(
    element.querySelectorAll<HTMLElement>(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])',
    ),
  ).filter((target) => !target.hasAttribute("disabled") && !target.getAttribute("aria-hidden"));

export function ReturnTitleDialog({
  onCancel,
  onConfirm,
}: {
  onCancel: () => void;
  onConfirm: () => void;
}) {
  const dialogRef = useRef<HTMLDivElement | null>(null);
  const cancelButtonRef = useRef<HTMLButtonElement | null>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    previousFocusRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
    cancelButtonRef.current?.focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.preventDefault();
        onCancel();
        return;
      }

      if (event.key !== "Tab" || !dialogRef.current) {
        return;
      }

      const tabbableElements = getTabbableElements(dialogRef.current);
      const firstElement = tabbableElements[0];
      const lastElement = tabbableElements.at(-1);

      if (!firstElement || !lastElement) {
        event.preventDefault();
        dialogRef.current.focus();
        return;
      }

      if (event.shiftKey && document.activeElement === firstElement) {
        event.preventDefault();
        lastElement.focus();
        return;
      }

      if (!event.shiftKey && document.activeElement === lastElement) {
        event.preventDefault();
        firstElement.focus();
      }
    };

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      previousFocusRef.current?.focus();
    };
  }, [onCancel]);

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
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="return-title-dialog-title"
        tabIndex={-1}
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
            ref={cancelButtonRef}
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
