import { motion, type MotionValue, type PanInfo } from "framer-motion";
import type { PointerEvent as ReactPointerEvent, RefObject } from "react";
import { TownMapDeploymentPreview } from "./TownMapDeploymentPreview";
import { TownMapLockedOverlay } from "./TownMapLockedOverlay";
import { TownMapPlacedItem } from "./TownMapPlacedItem";
import type { PlacedItem } from "./types";

export interface DeploymentPreview {
  id: string;
  isUnlocked: boolean;
  name: string;
  src: string;
  width: number;
  x: number;
  y: number;
}

interface TownMapProps {
  baseSrc: string;
  currentGuildLevel: number;
  deploymentPreview: DeploymentPreview | null;
  dragConstraints: { left: number; right: number; top: number; bottom: number };
  mapRef: RefObject<HTMLDivElement | null>;
  mapX: MotionValue<number>;
  mapY: MotionValue<number>;
  newlyDeployedItemId: string | null;
  onCancelDeployment: () => void;
  onCommitDeployment: () => void;
  onMoveDeploymentPreview: (event: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => void;
  onMoveItem: (
    item: PlacedItem,
    event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => void;
  onSelectItem: (id: string) => void;
  onStoreItem: (item: PlacedItem) => void;
  placedItems: PlacedItem[];
  scale: number;
  selectedPlacedItemId: string | null;
  stopNestedDrag: (event: ReactPointerEvent<HTMLElement>) => void;
  storingPlacedItemIds: string[];
  unlockClearingLevel: number | null;
}

export function TownMap({
  baseSrc,
  currentGuildLevel,
  deploymentPreview,
  dragConstraints,
  mapRef,
  mapX,
  mapY,
  newlyDeployedItemId,
  onCancelDeployment,
  onCommitDeployment,
  onMoveDeploymentPreview,
  onMoveItem,
  onSelectItem,
  onStoreItem,
  placedItems,
  scale,
  selectedPlacedItemId,
  stopNestedDrag,
  storingPlacedItemIds,
  unlockClearingLevel,
}: TownMapProps) {
  const isDeployMode = deploymentPreview !== null;

  return (
    <motion.div
      ref={mapRef}
      className={`absolute left-0 top-0 h-[200vh] w-[200vw] ${isDeployMode ? "cursor-default" : "cursor-grab active:cursor-grabbing"}`}
      drag={!isDeployMode}
      dragConstraints={dragConstraints}
      dragElastic={0.08}
      dragMomentum={false}
      style={{
        x: mapX,
        y: mapY,
        scale,
        touchAction: "none",
        transformOrigin: "top left",
        userSelect: "none",
      }}
    >
      <img
        className="pixelated"
        src={baseSrc}
        alt=""
        aria-hidden="true"
        draggable={false}
        style={{
          position: "absolute",
          inset: 0,
          width: "100%",
          height: "100%",
          objectFit: "cover",
          objectPosition: "center bottom",
          pointerEvents: "none",
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(180deg, rgba(4, 18, 18, 0.1) 0%, rgba(5, 16, 12, 0.04) 48%, rgba(6, 15, 10, 0.3) 100%)",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      <div
        aria-hidden="true"
        className="bg-[radial-gradient(ellipse_at_center,_transparent_20%,_rgba(0,0,0,0.6)_100%)]"
        style={{
          position: "absolute",
          inset: 0,
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <TownMapLockedOverlay
        currentGuildLevel={currentGuildLevel}
        unlockClearingLevel={unlockClearingLevel}
      />

      {placedItems.map((item) => {
        const isSelected = selectedPlacedItemId === item.id;
        const isStoring = storingPlacedItemIds.includes(item.id);
        const isNewlyDeployed = newlyDeployedItemId === item.id;

        return (
          <TownMapPlacedItem
            key={item.id}
            isDeployMode={isDeployMode}
            isNewlyDeployed={isNewlyDeployed}
            isSelected={isSelected}
            isStoring={isStoring}
            item={item}
            onMoveItem={onMoveItem}
            onSelectItem={onSelectItem}
            onStoreItem={onStoreItem}
            stopNestedDrag={stopNestedDrag}
          />
        );
      })}

      {deploymentPreview && (
        <TownMapDeploymentPreview
          key={deploymentPreview.id}
          deploymentPreview={deploymentPreview}
          onCancelDeployment={onCancelDeployment}
          onCommitDeployment={onCommitDeployment}
          onMoveDeploymentPreview={onMoveDeploymentPreview}
          stopNestedDrag={stopNestedDrag}
        />
      )}
    </motion.div>
  );
}
