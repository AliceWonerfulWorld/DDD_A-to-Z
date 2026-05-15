import { motion } from "framer-motion";
import type { CurrentUser } from "../../features/auth/types";

interface TitleUserBadgeProps {
  user: CurrentUser;
}

export function TitleUserBadge({ user }: TitleUserBadgeProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: -12 }}
      animate={{ opacity: 1, y: 0 }}
      style={{
        position: "fixed",
        top: "clamp(14px, 3vw, 28px)",
        right: "clamp(14px, 3vw, 28px)",
        zIndex: 4,
        display: "inline-flex",
        alignItems: "center",
        gap: "10px",
        maxWidth: "min(calc(100vw - 28px), 360px)",
        padding: "10px 16px",
        border: "2px solid #39ff14",
        boxShadow: "3px 3px 0 #1a7a00",
        background: "rgba(4, 18, 12, 0.82)",
        backdropFilter: "blur(2px)",
      }}
    >
      <img
        src={user.avatar_url}
        alt={user.username}
        style={{ width: 28, height: 28, borderRadius: "50%", flex: "0 0 auto" }}
      />
      <span
        style={{
          minWidth: 0,
          overflow: "hidden",
          textOverflow: "ellipsis",
          whiteSpace: "nowrap",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.85rem",
          color: "#39ff14",
          letterSpacing: "0.05em",
        }}
      >
        {user.username}
      </span>
      <span
        style={{
          flex: "0 0 auto",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.68rem",
          color: "#39ff1480",
        }}
      >
        ▶ LOGGED IN
      </span>
    </motion.div>
  );
}
