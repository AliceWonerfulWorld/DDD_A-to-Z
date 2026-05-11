import { motion } from "framer-motion";
import { SPRITE_ASSETS } from "../constants/assets";

interface RustSamuraiProps {
  className?: string;
  style?: React.CSSProperties;
}

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

export function RustSamurai({ className = "", style }: RustSamuraiProps) {
  // キャラクターのサイズ倍率（1.0 = 64x128）を2.0に変更し、さらに大きく表示
  const scale = 2.0;
  const frameWidth = 64 * scale;
  const frameHeight = 128 * scale;
  const totalWidth = frameWidth * 6; // 6コマ分の総幅
  
  return (
    <motion.div
      className={className}
      // 本来の画像は 6コマ で縦長（1コマの比率が1:2）です。
      // 指定した scale 倍率に合わせてアニメーションの移動距離も計算します。
      animate={{ backgroundPositionX: ["0px", `-${totalWidth}px`] }}
      transition={{
        duration: 0.8,
        repeat: Infinity,
        ease: steppedEase(6),
      }}
      style={{
        width: `${frameWidth}px`, // スケール適用後の1コマの幅
        height: `${frameHeight}px`, // スケール適用後の高さ
        backgroundImage: `url(${SPRITE_ASSETS.RUST_SAMURAI})`,
        backgroundSize: "auto 100%", // 高さを基準にアスペクト比を維持（潰れないようにする）
        backgroundRepeat: "no-repeat",
        imageRendering: "pixelated", // ドット絵がぼやけないようにする
        ...style,
      }}
    />
  );
}
