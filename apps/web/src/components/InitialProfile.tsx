import { useState, useEffect, useRef, useCallback } from "react";
import { motion } from "framer-motion";
import { ParticleBackground } from "./ParticleBackground";

// アセット管理用のMap定義
const SPRITE_ASSETS = {
  // 後からの画像差し替えを容易にする設計
  // 暫定対応: 全言語共通でRustの侍をサンプルとして表示
  RUST_SAMURAI: "/assets/sprites/rust_samurai.png",
} as const;

interface InitialProfileProps {
  onComplete: (username: string) => void;
}

// カクカクした動きを実現するためのカスタムイージング関数
const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

// エンジニア要素を強調するソースコードの断片
const CODE_SNIPPETS = [
  'fn main() {\n  println!("Hello, world!");\n}',
  'public static void main(String[] args) {\n  System.out.println("Hello World");\n}',
  'def hello():\n    print("Hello, world!")',
  'func main() {\n  fmt.Println("Hello, world!")\n}',
  'const sum = (a: number, b: number) => a + b;',
  'SELECT * FROM users WHERE id = 1;',
  'int main() {\n  printf("Hello, World!");\n  return 0;\n}'
];

function CodeRain() {
  return (
    <div style={{ position: "absolute", inset: 0, overflow: "hidden", pointerEvents: "none", zIndex: 0 }}>
      {CODE_SNIPPETS.map((code, i) => (
        <motion.div
          key={i}
          initial={{ y: "110vh", opacity: 0 }}
          animate={{
            y: ["110vh", "-30vh"],
            opacity: [0, 0.15, 0.15, 0],
          }}
          transition={{
            duration: 25 + (i % 4) * 8,
            repeat: Infinity,
            delay: i * 3,
            ease: "linear",
          }}
          style={{
            position: "absolute",
            left: `${8 + (i * 14) % 75}%`,
            fontFamily: "var(--font-dot)",
            fontSize: "0.85rem",
            color: "var(--color-neon-cyan)",
            whiteSpace: "pre",
            filter: "blur(0.5px)",
            textAlign: "left",
          }}
        >
          {code}
        </motion.div>
      ))}
    </div>
  );
}

export function InitialProfile({ onComplete }: InitialProfileProps) {
  const [username, setUsername] = useState("octocat"); // GitHubからの取得名を想定
  const [displayedText, setDisplayedText] = useState("");
  const fullText = "歓迎しよう、新たな挑戦者よ。\n君のコードネームを教えてくれ。";

  // Web Audio APIによるピコピコ音の準備
  const audioCtxRef = useRef<AudioContext | null>(null);

  useEffect(() => {
    try {
      const AudioContext = window.AudioContext || (window as unknown as { webkitAudioContext?: typeof window.AudioContext }).webkitAudioContext;
      if (AudioContext) {
        audioCtxRef.current = new AudioContext();
      }
    } catch {
      // AudioContext not supported
    }
  }, []);

  const playBeep = useCallback(() => {
    if (!audioCtxRef.current || audioCtxRef.current.state === "suspended") return;
    try {
      const ctx = audioCtxRef.current;
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();
      osc.type = "square";
      // 少しランダムに周波数を揺らしてレトロな「喋り声」感を出す
      osc.frequency.setValueAtTime(700 + Math.random() * 100, ctx.currentTime);
      gain.gain.setValueAtTime(0.03, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.06);
      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start();
      osc.stop(ctx.currentTime + 0.06);
    } catch {
      // ignore
    }
  }, []);

  // タイプライターエフェクト
  useEffect(() => {
    let i = 0;
    const interval = setInterval(() => {
      if (i <= fullText.length) {
        setDisplayedText(fullText.slice(0, i));
        
        // スペースや改行以外の文字が表示されるタイミングで音を鳴らす
        if (i < fullText.length && fullText[i] !== " " && fullText[i] !== "\n") {
          playBeep();
        }
        i++;
      } else {
        clearInterval(interval);
      }
    }, 90); // 1文字ずつ表示 (ピコピコ感を出すため少し遅めの90msに設定)
    return () => clearInterval(interval);
  }, [fullText, playBeep]);

  return (
    <div
      className="flex flex-col items-center justify-center min-h-[100svh] p-4 relative overflow-hidden"
      style={{
        background: "radial-gradient(circle at 50% 120%, var(--color-navy-mid) 0%, var(--color-navy) 60%, #050510 100%)",
        fontFamily: "var(--font-dot)",
        color: "var(--color-pixel-white)",
      }}
    >
      {/* 冒険の始まりを演出するパーティクルと光 */}
      <ParticleBackground />

      {/* エンジニア感を強調するソースコードの断片 */}
      <CodeRain />

      {/* 魔法陣のような光の演出（下部から） */}
      <motion.div
        animate={{ opacity: [0.3, 0.6, 0.3] }}
        transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
        style={{
          position: "absolute",
          bottom: "-20vh",
          left: "50%",
          transform: "translateX(-50%)",
          width: "150vw",
          height: "60vh",
          background: "radial-gradient(ellipse at center, rgba(0, 245, 255, 0.15) 0%, transparent 70%)",
          pointerEvents: "none",
          zIndex: 0,
        }}
      />

      {/* 背景のスキャンライン（Appと共通） */}
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

      {/* ウィンドウ展開アニメーション */}
      <motion.div
        initial={{ scaleY: 0, opacity: 0 }}
        animate={{ scaleY: 1, opacity: 1 }}
        transition={{ duration: 0.4, ease: steppedEase(6) }}
        style={{
          position: "relative",
          zIndex: 2,
          border: "4px solid var(--color-gold)",
          background: "var(--color-navy-light)",
          boxShadow: "0 0 30px rgba(0, 245, 255, 0.15), 8px 8px 0 rgba(0,0,0,0.8)",
          maxWidth: "500px",
          width: "100%",
          padding: "2rem",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: "2rem",
        }}
      >
        {/* ガイドテキスト */}
        <div style={{ minHeight: "3.5rem", fontSize: "1.1rem", whiteSpace: "pre-wrap", lineHeight: "1.6", width: "100%" }}>
          {displayedText}
          <motion.span
            animate={{ opacity: [1, 0] }}
            transition={{ duration: 0.8, repeat: Infinity, ease: steppedEase(2) }}
          >
            _
          </motion.span>
        </div>

        {/* アバター表示エリア */}
        <div style={{ position: "relative", width: "128px", height: "128px" }}>
          {/* 背後の回転魔法陣的な演出（レトロ風） */}
          <motion.div
            animate={{ rotate: 360 }}
            transition={{ duration: 15, repeat: Infinity, ease: "linear" }}
            style={{
              position: "absolute",
              top: "-20%",
              left: "-20%",
              width: "140%",
              height: "140%",
              border: "2px dashed var(--color-gold)",
              borderRadius: "50%",
              opacity: 0.15,
              zIndex: 0,
            }}
          />

          {/* カクカク浮遊するアバター */}
          <motion.div
            animate={{ y: [0, -12, 0] }}
            transition={{
              duration: 1.2,
              repeat: Infinity,
              ease: steppedEase(4),
            }}
            style={{
              width: "100%",
              height: "100%",
              border: "2px solid rgba(255,255,255,0.1)",
              background: "rgba(0,0,0,0.4)",
              boxShadow: "inset 0 0 10px rgba(0,0,0,0.5)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <img
              src={SPRITE_ASSETS.RUST_SAMURAI}
              alt="Rust Samurai"
              style={{
                width: "100%",
                height: "100%",
                objectFit: "contain",
                imageRendering: "pixelated",
              }}
              onError={(e) => {
                // 画像がない場合のフォールバック（暫定対応）
                e.currentTarget.style.display = "none";
                const parent = e.currentTarget.parentElement;
                if (parent && !parent.querySelector('.fallback-text')) {
                  const span = document.createElement('span');
                  span.className = 'fallback-text';
                  span.style.color = '#666';
                  span.style.fontSize = '0.8rem';
                  span.textContent = 'NO SPRITE';
                  parent.appendChild(span);
                }
              }}
            />
          </motion.div>
          {/* 影 */}
          <motion.div
            animate={{ scale: [1, 0.8, 1], opacity: [0.5, 0.2, 0.5] }}
            transition={{
              duration: 1.2,
              repeat: Infinity,
              ease: steppedEase(4),
            }}
            style={{
              position: "absolute",
              bottom: "-20px",
              left: "50%",
              transform: "translateX(-50%)",
              width: "60px",
              height: "10px",
              background: "rgba(0,0,0,0.6)",
              borderRadius: "50%", // 影だけは少し丸みを持たせる（レトロゲームでよくある手法）
            }}
          />
        </div>

        {/* ユーザー名入力 */}
        <div style={{ width: "100%" }}>
          <label
            style={{
              display: "block",
              marginBottom: "0.5rem",
              fontSize: "0.8rem",
              color: "var(--color-gold)",
              letterSpacing: "0.1em",
            }}
          >
            ▶ ENTER YOUR NAME
          </label>
          <div style={{ position: "relative" }}>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={{
                width: "100%",
                padding: "0.8rem",
                fontSize: "1.2rem",
                fontFamily: "var(--font-dot)",
                background: "rgba(0,0,0,0.5)",
                color: "var(--color-pixel-white)",
                border: "2px solid rgba(255,255,255,0.4)",
                outline: "none",
                textAlign: "center",
              }}
              onFocus={(e) => (e.target.style.borderColor = "var(--color-gold)")}
              onBlur={(e) => (e.target.style.borderColor = "rgba(255,255,255,0.4)")}
            />
          </div>
        </div>

        {/* 決定ボタン */}
        <motion.button
          whileHover={{ scale: 1.02, filter: "brightness(1.1)" }}
          whileTap={{ scale: 0.98, y: 4, boxShadow: "0px 0px 0 var(--color-gold-dark)" }}
          onClick={() => onComplete(username)}
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
            cursor: "pointer",
            letterSpacing: "0.05em",
          }}
        >
          BEGIN JOURNEY
        </motion.button>
      </motion.div>
    </div>
  );
}
