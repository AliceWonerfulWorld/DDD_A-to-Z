import { motion } from "framer-motion";
import { SPRITE_ASSETS } from "../../constants/assets";
import { steppedEase } from "../../lib/animationUtils";

interface GopherSpriteProps {
  className?: string;
  frameCount?: number;
  row?: number;
  style?: React.CSSProperties;
}

const FRAME_WIDTH = 192;
const FRAME_HEIGHT = 208;
const COLUMNS = 8;
const ROWS = 9;
const IDLE_FRAMES = 6;
const DISPLAY_WIDTH = 132;
const DISPLAY_HEIGHT = 143;

export function GopherSprite({
  className = "",
  frameCount = IDLE_FRAMES,
  row = 0,
  style,
}: GopherSpriteProps) {
  const scale = DISPLAY_WIDTH / FRAME_WIDTH;
  const displaySheetWidth = Math.round(FRAME_WIDTH * COLUMNS * scale);
  const displaySheetHeight = Math.round(FRAME_HEIGHT * ROWS * scale);
  const displayFrameWidth = Math.round(DISPLAY_WIDTH);
  const displayFrameHeight = Math.round(DISPLAY_HEIGHT);
  const frameStep = Math.round(FRAME_WIDTH * scale);
  const totalMoveX = frameStep * frameCount;
  const rowOffsetY = Math.round(FRAME_HEIGHT * row * scale);

  return (
    <motion.div
      className={className}
      animate={{ backgroundPositionX: ["0px", `-${totalMoveX}px`] }}
      transition={{
        duration: 0.9,
        repeat: Infinity,
        ease: steppedEase(frameCount),
      }}
      style={{
        width: `${displayFrameWidth}px`,
        height: `${displayFrameHeight}px`,
        backgroundImage: `url(${SPRITE_ASSETS.GOPHER})`,
        backgroundPositionY: `-${rowOffsetY}px`,
        backgroundRepeat: "no-repeat",
        backgroundSize: `${displaySheetWidth}px ${displaySheetHeight}px`,
        imageRendering: "pixelated",
        ...style,
      }}
    />
  );
}
