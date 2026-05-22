import { motion, type PanInfo } from "framer-motion";
import type { CSSProperties, PointerEvent as ReactPointerEvent } from "react";
import { steppedEase } from "../../lib/animationUtils";
import type { DeploymentPreview } from "./TownMap";

interface TownMapDeploymentPreviewProps {
  deploymentPreview: DeploymentPreview;
  onCancelDeployment: () => void;
  onCommitDeployment: () => void;
  onMoveDeploymentPreview: (event: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => void;
  stopNestedDrag: (event: ReactPointerEvent<HTMLElement>) => void;
}

export function TownMapDeploymentPreview({
  deploymentPreview,
  onCancelDeployment,
  onCommitDeployment,
  onMoveDeploymentPreview,
  stopNestedDrag,
}: TownMapDeploymentPreviewProps) {
  return (
    <>
      <motion.div
        aria-label={`${deploymentPreview.name} deployment preview`}
        aria-roledescription="deployment preview"
        drag
        dragElastic={0}
        dragMomentum={false}
        dragSnapToOrigin
        onDragEnd={onMoveDeploymentPreview}
        onPointerDown={stopNestedDrag}
        initial={{ opacity: 0, scale: 0.92, y: -18 }}
        animate={{
          opacity: deploymentPreview.isUnlocked ? 0.68 : 0.42,
          scale: deploymentPreview.isUnlocked ? 1 : 0.96,
          y: 0,
        }}
        transition={{ duration: 0.18, ease: steppedEase(5) }}
        whileDrag={{ opacity: 0.76, scale: 1.03, zIndex: 13 }}
        style={{
          position: "absolute",
          left: deploymentPreview.x,
          top: deploymentPreview.y,
          width: deploymentPreview.width,
          cursor: "grab",
          filter: deploymentPreview.isUnlocked
            ? "drop-shadow(10px 14px 0 rgba(0,0,0,0.28)) drop-shadow(0 0 18px rgba(116,247,161,0.58))"
            : "grayscale(0.65) drop-shadow(10px 14px 0 rgba(0,0,0,0.36)) drop-shadow(0 0 18px rgba(255,77,109,0.5))",
          touchAction: "none",
          transformOrigin: "50% 80%",
          userSelect: "none",
          zIndex: 12,
        }}
      >
        <img
          className="pixelated"
          src={deploymentPreview.src}
          alt=""
          draggable={false}
          style={{
            display: "block",
            width: "100%",
            height: "auto",
            pointerEvents: "none",
          }}
        />
        <div
          aria-hidden="true"
          style={{
            position: "absolute",
            left: "50%",
            bottom: "-14px",
            width: "72%",
            height: "16px",
            transform: "translateX(-50%)",
            border: `2px solid ${deploymentPreview.isUnlocked ? "rgba(116,247,161,0.74)" : "rgba(255,77,109,0.72)"}`,
            background: deploymentPreview.isUnlocked
              ? "rgba(116,247,161,0.12)"
              : "rgba(255,77,109,0.12)",
            boxShadow: deploymentPreview.isUnlocked
              ? "0 0 14px rgba(116,247,161,0.34)"
              : "0 0 14px rgba(255,77,109,0.32)",
          }}
        />
        <div
          style={{
            position: "absolute",
            left: "50%",
            top: "calc(100% + 12px)",
            transform: "translateX(-50%)",
            border: `2px solid ${deploymentPreview.isUnlocked ? "rgba(116,247,161,0.78)" : "rgba(255,77,109,0.78)"}`,
            background: "rgba(3, 10, 24, 0.9)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
            color: deploymentPreview.isUnlocked ? "#74f7a1" : "#ff9aae",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.72rem",
            lineHeight: 1,
            padding: "6px 8px",
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            whiteSpace: "nowrap",
          }}
        >
          {deploymentPreview.isUnlocked ? "BUILD READY" : "LOCKED AREA"}
        </div>
      </motion.div>

      <div
        onPointerDown={stopNestedDrag}
        style={{
          position: "absolute",
          left: deploymentPreview.x + deploymentPreview.width + 12,
          top: deploymentPreview.y,
          zIndex: 13,
          display: "grid",
          gap: "8px",
          width: "104px",
        }}
      >
        <button
          type="button"
          onClick={onCommitDeployment}
          style={deployActionButtonStyle(deploymentPreview.isUnlocked)}
        >
          BUILD
        </button>
        <button type="button" onClick={onCancelDeployment} style={deployActionButtonStyle(false)}>
          CANCEL
        </button>
      </div>
    </>
  );
}

function deployActionButtonStyle(isPrimary: boolean): CSSProperties {
  return {
    minHeight: "34px",
    border: `2px solid ${isPrimary ? "rgba(116, 247, 161, 0.86)" : "rgba(255, 217, 102, 0.76)"}`,
    borderBottomColor: isPrimary ? "rgba(24, 83, 45, 0.95)" : "rgba(96, 62, 22, 0.98)",
    borderRightColor: isPrimary ? "rgba(24, 83, 45, 0.95)" : "rgba(96, 62, 22, 0.98)",
    background: isPrimary ? "rgba(4, 67, 37, 0.92)" : "rgba(3, 10, 24, 0.9)",
    boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
    color: isPrimary ? "#74f7a1" : "#ffd966",
    cursor: "pointer",
    fontFamily: '"Press Start 2P", "DotGothic16", monospace',
    fontSize: "0.46rem",
    lineHeight: 1,
    padding: "8px 9px",
    textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
  };
}
