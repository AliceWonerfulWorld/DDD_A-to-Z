import { motion } from "framer-motion";
import { SPRITE_ASSETS } from "../../constants/assets";

interface RustSamuraiProps {
  className?: string;
  style?: React.CSSProperties;
}

const steppedEase = (steps: number) => (t: number) => Math.floor(t * steps) / steps;

export function RustSamurai({ className = "", style }: RustSamuraiProps) {
  // ユーザーが求めるスケール（2.0 = 128x256）
  const scale = 2.0;

  // 元画像の正しい解像度（866x288）と、1コマの正しいサイズ（144x288）
  // 866pxのうち、左右に1pxずつの余白があると推測されるため、それを考慮した厳密な計算を行います

  // 表示上のサイズ（全て整数に強制）
  const displayFrameWidth = Math.round(64 * scale); // 128
  const displayFrameHeight = Math.round(128 * scale); // 256

  // CSSとして指定する厳密なピクセル値（サブピクセルを防ぐため全て整数化）
  // 元画像の左余白(1px)のスケール後の値も四捨五入して整数にする
  const renderRatio = displayFrameHeight / 288;
  const startX = Math.round(1 * renderRatio);

  // 移動距離も完璧な整数（1コマの幅 × 6）
  const totalMoveX = displayFrameWidth * 6; // 768

  // 背景の全体幅（本来は866pxですが、余白を削って768pxに押し込めるのではなく、
  // 1コマ分の比率から逆算して完璧な整数幅を割り出します）
  // 144px : 866px = 128px : bgWidth -> bgWidth = 866 * (128/144) = 769.77... -> 770px
  const bgWidth = Math.round(866 * renderRatio);

  return (
    <motion.div
      className={className}
      // サブピクセル（小数）を完全に排除した整数ピクセルのみでアニメーション
      animate={{ backgroundPositionX: [`-${startX}px`, `-${startX + totalMoveX}px`] }}
      transition={{
        duration: 0.8,
        repeat: Infinity,
        ease: steppedEase(6),
      }}
      style={{
        width: `${displayFrameWidth}px`,
        height: `${displayFrameHeight}px`,
        backgroundImage: `url(${SPRITE_ASSETS.RUST_SAMURAI})`,
        // 背景サイズも整数化
        backgroundSize: `${bgWidth}px ${displayFrameHeight}px`,
        backgroundRepeat: "no-repeat",
        imageRendering: "pixelated",
        ...style,
      }}
    />
  );
}
