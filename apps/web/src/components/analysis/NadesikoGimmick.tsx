import { useEffect, useRef, useState } from "react";
import { GopherSprite } from "../shared/GopherSprite";

export function NadesikoGimmick() {
  const containerRef = useRef<HTMLDivElement>(null);
  const [gopherY, setGopherY] = useState(90);
  const [gameState, setGameState] = useState("待機");

  useEffect(() => {
    // @ts-ignore
    window.syncNakoState = (y: number, state: string) => {
      setGopherY(y);
      setGameState(state);
    };

    let isUnmounted = false;

    const code = `
「#nako-canvas」へ描画開始

GY＝90
GVY＝0
BX＝400
BY＝90
B速度＝6
スコア＝0
状態＝「待機」

「window.nakoClickFlag = 0」をJS実行

●ループとは
  描画クリア
  
  クリックは「window.nakoClickFlag」をJS実行
  もしクリックが1ならば
    「window.nakoClickFlag = 0」をJS実行
もし状態が「待機」ならば
  状態＝「実行」
  BX＝400
  スコア＝0
  B速度＝6
    違えば
      もし状態が「実行」ならば
        もしGY≧90ならば
          GVY＝-11

        ここまで
      違えば
        状態＝「待機」
        GY＝90
        BX＝400
        スコア＝0
      ここまで
    ここまで
  ここまで
  
  もし状態が「待機」ならば
    「white」に塗色設定
    「14px sans-serif」に描画フォント設定
    「クリックでスタート！」を[0, 50]に文字描画
  違えば
    もし状態が「実行」ならば
      GY＝GY＋GVY
      GVY＝GVY＋0.9
      もしGY＞90ならば
        GY＝90
        GVY＝0
      ここまで
      
      BX＝BX－B速度
      もしBX＜-20ならば
        BX＝400
        スコア＝スコア＋1
        乱数結果は2の乱数
        もし乱数結果が0ならば
          BY＝90
        違えば
          BY＝60
        ここまで
        B速度＝B速度＋0.4
      ここまで
      
      ゴーファー右＝20＋40
      ゴーファー下＝GY＋40
      バグ右＝BX＋20
      バグ下＝BY＋20
      
      X重なり＝0
      もしBX＜ゴーファー右ならば
        もし20＜バグ右ならば
          X重なり＝1
        ここまで
      ここまで
      
      Y重なり＝0
      もしBY＜ゴーファー下ならば
        もしGY＜バグ下ならば
          Y重なり＝1
        ここまで
      ここまで
      
      もしX重なりが1ならば
        もしY重なりが1ならば
          状態＝「終了」
        ここまで
      ここまで
    違えば
      「#ff5555」に塗色設定
      「14px sans-serif」に描画フォント設定
      「解析が遅れます！（衝突）」を[0, 50]に文字描画
    ここまで
  ここまで

  「#ff5555」に塗色設定
  [BX, BY, 20, 20]に四角描画
  
  「white」に塗色設定
  「12px sans-serif」に描画フォント設定
  スコア文字列は「SCORE: 」＆スコア
  スコア文字列を[0, 10]に文字描画

  JSコマンドは「window.syncNakoState(」＆GY＆「, '」＆状態＆「')」
  JSコマンドをJS実行
ここまで

「ループ」を0.033秒毎
`;
    const startNako = async () => {
      // @ts-ignore
      const nako = window.navigator.nako3;
      if (nako) {
        try {
          await nako.run(code);

          const handleGlobalClick = () => {
            if (isUnmounted) return;
            // @ts-ignore
            window.nakoClickFlag = 1;
            setTimeout(() => {
              // @ts-ignore
              window.nakoClickFlag = 0;
            }, 50);
          };
          window.addEventListener("pointerdown", handleGlobalClick);

          // 終了時にタイマーなどを停止するクリーンアップ
          // @ts-ignore
          window.__nakoCleanup = () => {
            window.removeEventListener("pointerdown", handleGlobalClick);
            try {
              nako.run("全タイマー停止");
            } catch (err) {}
          };
        } catch (e) {
          console.error("Nadesiko script error:", e);
        }
      } else {
        if (!isUnmounted) {
          setTimeout(startNako, 100);
        }
      }
    };

    void startNako();

    return () => {
      isUnmounted = true;
      // @ts-ignore
      if (window.__nakoCleanup) window.__nakoCleanup();
    };
  }, []);

  return (
    <div
      ref={containerRef}
      style={{
        position: "fixed",
        bottom: 0,
        left: 0,
        zIndex: 50,
        pointerEvents: "none",
      }}
    >
      <canvas
        id="nako-canvas"
        width={400}
        height={150}
        style={{
          cursor: "pointer",
          pointerEvents: "auto",
        }}
      />

      {/* 画面左下に配置するゴーファーくん */}
      <div
        style={{
          position: "absolute",
          left: "15px",
          top: gopherY - 15 + "px",
          pointerEvents: "none",
          transform: "scale(0.35)",
          transformOrigin: "top left",
        }}
      >
        <GopherSprite row={gameState === "終了" ? 2 : 0} />
      </div>
    </div>
  );
}
