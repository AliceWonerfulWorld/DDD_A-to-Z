import { motion } from "framer-motion";

interface HomeFirstVisitGuideProps {
  onDismiss: () => void;
}

const guideSteps = [
  { label: "ステータス", text: "左上でプレイヤー情報、右上で獲得CPを確認できます。" },
  { label: "相棒", text: "画面中央のGopherに話しかけると、ちょっとした反応が返ってきます。" },
  { label: "次の行き先", text: "下のメニューから、戦場・ギルド・マイステータスへ移動できます。" },
] as const;

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

export function HomeFirstVisitGuide({ onDismiss }: HomeFirstVisitGuideProps) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.2, ease: steppedEase(3) }}
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 20,
        pointerEvents: "auto",
        background:
          "radial-gradient(circle at 50% 54%, rgba(0,245,255,0.08), rgba(4,8,18,0.62) 48%, rgba(4,8,18,0.82) 100%)",
        overflow: "auto",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "clamp(16px, 4vw, 32px)",
      }}
    >
      <motion.div
        aria-hidden="true"
        initial={{ opacity: 0 }}
        animate={{ opacity: [0, 0.62, 0.36] }}
        transition={{ duration: 1, ease: steppedEase(5) }}
        style={{
          position: "absolute",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 6px, rgba(255,255,255,0.06) 6px, rgba(255,255,255,0.06) 7px)",
          pointerEvents: "none",
        }}
      />

      <motion.div
        initial={{ opacity: 0, y: 16, scale: 0.96 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ delay: 0.18, duration: 0.38, ease: steppedEase(5) }}
        style={{
          width: "min(640px, calc(100vw - 28px))",
          maxHeight: "calc(100svh - 32px)",
          overflowY: "auto",
          position: "relative",
          zIndex: 1,
          border: "3px solid #ffd700",
          background: "rgba(14, 28, 64, 0.94)",
          boxShadow: "0 0 28px rgba(0,245,255,0.25), 7px 7px 0 rgba(0,0,0,0.66)",
          padding: "clamp(16px, 3vw, 24px)",
        }}
      >
        <motion.div
          animate={{ opacity: [1, 0.58, 1] }}
          transition={{ duration: 0.8, repeat: Infinity, ease: steppedEase(3) }}
          style={{
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(1.35rem, 3vw, 1.85rem)",
            color: "#6dff5f",
            letterSpacing: "0.08em",
            marginBottom: "12px",
            textAlign: "center",
          }}
        >
          ホームベース起動
        </motion.div>
        <div
          style={{
            fontFamily: '"DotGothic16", monospace',
            fontSize: "clamp(0.95rem, 2vw, 1.08rem)",
            lineHeight: 1.5,
            color: "#fff8d7",
            marginBottom: "16px",
            textAlign: "center",
          }}
        >
          ここが冒険の拠点です。まずは画面の見方だけ、さくっと確認しましょう。
        </div>

        <div
          style={{
            display: "grid",
            gap: "10px",
            marginBottom: "16px",
          }}
        >
          {guideSteps.map((step, i) => (
            <motion.div
              key={step.label}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.42 + i * 0.16, duration: 0.24, ease: steppedEase(4) }}
              style={{
                display: "grid",
                gridTemplateColumns: "minmax(88px, 0.34fr) 1fr",
                gap: "12px",
                alignItems: "center",
                border: "2px solid rgba(255, 211, 107, 0.62)",
                background: "rgba(4, 8, 18, 0.42)",
                padding: "10px 12px",
              }}
            >
              <div
                style={{
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "1rem",
                  color: "#ffd700",
                  letterSpacing: "0.06em",
                  whiteSpace: "nowrap",
                }}
              >
                ▶ {step.label}
              </div>
              <div
                style={{
                  fontFamily: '"DotGothic16", monospace',
                  fontSize: "0.95rem",
                  lineHeight: 1.45,
                  color: "#fff8d7",
                }}
              >
                {step.text}
              </div>
            </motion.div>
          ))}
        </div>

        <button
          type="button"
          onClick={onDismiss}
          style={{
            width: "100%",
            padding: "12px 16px",
            border: "none",
            background: "#ffd700",
            color: "#0a1024",
            boxShadow: "0 4px 0 #b88900",
            cursor: "pointer",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "1.1rem",
            fontWeight: 700,
            letterSpacing: "0.05em",
          }}
        >
          冒険をはじめる
        </button>
      </motion.div>
    </motion.div>
  );
}
