import { motion } from "framer-motion";

export function TitleLogo() {
  return (
    <motion.div
      animate={{ y: [0, -15, 0] }}
      transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
      style={{ display: "inline-block" }}
    >
      {/* ピクセル風シャドウ付きタイトル */}
      <h1
        style={{
          fontFamily: '"Press Start 2P", monospace',
          fontSize: "clamp(2rem, 6vw, 4.5rem)",
          color: "#ffd700",
          margin: 0,
          letterSpacing: "0.08em",
          lineHeight: 1.2,
          /* ボカシなしのソリッドシャドウ */
          textShadow: [
            "4px 4px 0 #ccac00",
            "8px 8px 0 #997e00",
            "0 0 40px rgba(255,215,0,0.4)",
          ].join(", "),
        }}
      >
        LANG WAR
      </h1>
    </motion.div>
  );
}
