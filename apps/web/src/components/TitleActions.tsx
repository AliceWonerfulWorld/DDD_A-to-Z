import { motion } from "framer-motion";
import { GitHubLoginButton } from "./GitHubLoginButton";

interface TitleActionsProps {
  isLoggedIn: boolean;
  onLogin: () => void;
  onLogoutClick: () => void;
  onStart: () => void;
}

export function TitleActions({ isLoggedIn, onLogin, onLogoutClick, onStart }: TitleActionsProps) {
  if (!isLoggedIn) {
    return <GitHubLoginButton onClick={onLogin} />;
  }

  return (
    <div
      style={{
        display: "inline-flex",
        alignItems: "center",
        flexWrap: "wrap",
        justifyContent: "center",
        gap: "14px",
      }}
    >
      <motion.button
        type="button"
        onClick={onStart}
        whileHover={{ scale: 1.04 }}
        whileTap={{ y: 3, scale: 0.98 }}
        style={{
          padding: "12px 22px",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.9rem",
          fontWeight: "bold",
          background: "#ffd700",
          border: "3px solid #0a0a0a",
          boxShadow: "4px 4px 0 #0a0a0a",
          color: "#0a0a0a",
          cursor: "pointer",
          letterSpacing: "0.08em",
        }}
      >
        START
      </motion.button>
      <motion.button
        type="button"
        onClick={onLogoutClick}
        whileHover={{ scale: 1.04 }}
        whileTap={{ y: 3, scale: 0.98 }}
        style={{
          padding: "12px 18px",
          fontFamily: '"DotGothic16", monospace',
          fontSize: "0.78rem",
          background: "rgba(10,10,10,0.55)",
          border: "2px solid #ff5f56",
          boxShadow: "3px 3px 0 #7a211d",
          color: "#ffb0aa",
          cursor: "pointer",
          letterSpacing: "0.06em",
        }}
      >
        LOGOUT
      </motion.button>
    </div>
  );
}
