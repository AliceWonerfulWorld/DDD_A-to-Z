import { motion, type MotionValue, type PanInfo } from "framer-motion";
import type { PointerEvent as ReactPointerEvent, RefObject } from "react";
import { steppedEase } from "../../lib/animationUtils";
import type { PlacedItem } from "./types";

interface TownMapProps {
  baseSrc: string;
  dragConstraints: { left: number; right: number; top: number; bottom: number };
  mapRef: RefObject<HTMLDivElement | null>;
  mapX: MotionValue<number>;
  mapY: MotionValue<number>;
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
}

export function TownMap({
  baseSrc,
  dragConstraints,
  mapRef,
  mapX,
  mapY,
  onMoveItem,
  onSelectItem,
  onStoreItem,
  placedItems,
  scale,
  selectedPlacedItemId,
  stopNestedDrag,
  storingPlacedItemIds,
}: TownMapProps) {
  return (
    <motion.div
      ref={mapRef}
      className="absolute left-0 top-0 h-[200vh] w-[200vw] cursor-grab active:cursor-grabbing"
      drag
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

      {placedItems.map((item) => {
        const isSelected = selectedPlacedItemId === item.id;
        const isStoring = storingPlacedItemIds.includes(item.id);

        return (
          <motion.div
            key={item.id}
            animate={{
              opacity: isStoring ? 0 : 1,
              scale: isStoring ? 0.72 : 1,
              y: isStoring ? -22 : 0,
            }}
            transition={{ duration: 0.24, ease: steppedEase(6) }}
            style={{
              position: "absolute",
              left: item.x,
              top: item.y,
              width: item.width,
              height: "fit-content",
              outline: isSelected ? "3px solid rgba(255, 217, 102, 0.82)" : "3px solid transparent",
              outlineOffset: "4px",
              pointerEvents: isStoring ? "none" : "auto",
              transformOrigin: "50% 80%",
              zIndex: isSelected || isStoring ? 10 : 8,
            }}
          >
            <motion.img
              className="pixelated"
              src={item.src}
              alt={item.name}
              drag
              dragSnapToOrigin
              dragElastic={0}
              dragMomentum={false}
              onPointerDown={stopNestedDrag}
              onClick={() => onSelectItem(item.id)}
              onDragEnd={(event, info) => onMoveItem(item, event, info)}
              whileHover={{ scale: 1.02 }}
              whileDrag={{ scale: 1.05, zIndex: 12 }}
              style={{
                display: "block",
                width: "100%",
                height: "auto",
                cursor: "grab",
                touchAction: "none",
                userSelect: "none",
                filter: isSelected
                  ? "drop-shadow(10px 14px 0 rgba(0,0,0,0.3)) drop-shadow(0 0 12px rgba(255,217,102,0.72))"
                  : "drop-shadow(10px 14px 0 rgba(0,0,0,0.3))",
              }}
            />
            {isSelected && (
              <motion.button
                type="button"
                aria-label={`Store ${item.name}`}
                onPointerDown={stopNestedDrag}
                onClick={() => onStoreItem(item)}
                initial={{ opacity: 0, y: 6, scale: 0.92 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.18, ease: steppedEase(4) }}
                whileHover={{ y: -2, backgroundColor: "rgba(255, 217, 102, 0.22)" }}
                whileTap={{ y: 1, scale: 0.96 }}
                style={{
                  position: "absolute",
                  right: "-10px",
                  top: "-44px",
                  minWidth: "74px",
                  minHeight: "34px",
                  border: "2px solid rgba(255, 217, 102, 0.86)",
                  borderBottomColor: "rgba(96, 62, 22, 0.98)",
                  borderRightColor: "rgba(96, 62, 22, 0.98)",
                  background: "rgba(3, 10, 24, 0.9)",
                  boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
                  color: "#fff8d7",
                  cursor: "pointer",
                  fontFamily: "inherit",
                  fontSize: "0.52rem",
                  lineHeight: 1,
                  padding: "8px 9px",
                  textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
                  touchAction: "none",
                }}
              >
                STORE
              </motion.button>
            )}
          </motion.div>
        );
      })}
    </motion.div>
  );
}
