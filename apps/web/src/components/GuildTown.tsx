import { motion, useMotionValue, type PanInfo } from "framer-motion";
import {
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
  type WheelEvent,
} from "react";
import { steppedEase } from "../lib/animationUtils";

const MIN_SCALE = 0.5;
const MAX_SCALE = 2.5;
const ZOOM_STEP = 0.15;

type InventoryItemType = "tent" | "bonfire";

interface InventoryItem {
  type: InventoryItemType;
  name: string;
  count: number;
  src: string;
  mapWidth: number;
}

interface PlacedItem {
  id: string;
  type: InventoryItemType;
  name: string;
  src: string;
  x: number;
  y: number;
  width: number;
}

const INITIAL_INVENTORY: InventoryItem[] = [
  {
    type: "tent",
    name: "TENT",
    count: 2,
    src: "/town/tent.png",
    mapWidth: 180,
  },
  {
    type: "bonfire",
    name: "BONFIRE",
    count: 3,
    src: "/town/bonfire.png",
    mapWidth: 92,
  },
];

function clampValue(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value));
}

function isPointInsideRect(point: PanInfo["point"], rect: DOMRect) {
  return (
    point.x >= rect.left && point.x <= rect.right && point.y >= rect.top && point.y <= rect.bottom
  );
}

interface GuildTownProps {
  onNavigate: (path: string) => void;
  townLevel?: number;
  currentCp?: number;
  nextLevelCp?: number;
  baseSrc?: string;
  mainStructureSrc?: string;
  bonfireSrc?: string;
}

export function GuildTown({
  onNavigate,
  townLevel = 1,
  currentCp = 2500,
  nextLevelCp = 10000,
  baseSrc = "/town/glassfield.png",
  mainStructureSrc = "/town/tent.png",
  bonfireSrc = "/town/bonfire.png",
}: GuildTownProps) {
  const progress = Math.min(100, Math.max(0, (currentCp / nextLevelCp) * 100));
  const [viewport, setViewport] = useState({ width: 0, height: 0 });
  const [scale, setScale] = useState(1);
  const [inventory, setInventory] = useState<InventoryItem[]>(INITIAL_INVENTORY);
  const [placedItems, setPlacedItems] = useState<PlacedItem[]>([]);
  const mapRef = useRef<HTMLDivElement>(null);
  const inventoryRef = useRef<HTMLDivElement>(null);
  const mapX = useMotionValue(0);
  const mapY = useMotionValue(0);
  const dragConstraints = {
    left: Math.min(0, viewport.width - viewport.width * 2 * scale),
    right: 0,
    top: Math.min(0, viewport.height - viewport.height * 2 * scale),
    bottom: 0,
  };

  useEffect(() => {
    const updateViewport = () => {
      setViewport({ width: window.innerWidth, height: window.innerHeight });
    };

    updateViewport();
    window.addEventListener("resize", updateViewport);

    return () => window.removeEventListener("resize", updateViewport);
  }, []);

  useEffect(() => {
    if (viewport.width === 0 || viewport.height === 0) return;

    mapX.set(-viewport.width * 0.5);
    mapY.set(-viewport.height * 0.5);
  }, [mapX, mapY, viewport.height, viewport.width]);

  useEffect(() => {
    mapX.set(clampValue(mapX.get(), dragConstraints.left, dragConstraints.right));
    mapY.set(clampValue(mapY.get(), dragConstraints.top, dragConstraints.bottom));
  }, [
    dragConstraints.bottom,
    dragConstraints.left,
    dragConstraints.right,
    dragConstraints.top,
    mapX,
    mapY,
  ]);

  const handleZoom = (delta: number) => {
    setScale((currentScale) => clampValue(currentScale + delta, MIN_SCALE, MAX_SCALE));
  };

  const handleWheel = (event: WheelEvent<HTMLElement>) => {
    event.preventDefault();
    handleZoom(-event.deltaY * 0.0015);
  };

  const stopNestedDrag = (event: ReactPointerEvent<HTMLElement>) => {
    event.stopPropagation();
  };

  const getMapDropPoint = (point: PanInfo["point"], itemWidth: number) => {
    const mapElement = mapRef.current;
    if (!mapElement) return null;

    const mapRect = mapElement.getBoundingClientRect();
    if (!isPointInsideRect(point, mapRect)) return null;

    const inventoryRect = inventoryRef.current?.getBoundingClientRect();
    if (inventoryRect && isPointInsideRect(point, inventoryRect)) return null;

    const mapWidth = mapRect.width / scale;
    const mapHeight = mapRect.height / scale;
    const x = (point.x - mapRect.left) / scale - itemWidth / 2;
    const y = (point.y - mapRect.top) / scale - itemWidth / 2;

    return {
      x: clampValue(x, 0, Math.max(0, mapWidth - itemWidth)),
      y: clampValue(y, 0, Math.max(0, mapHeight - itemWidth)),
    };
  };

  const handleInventoryDragEnd = (
    item: InventoryItem,
    _event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => {
    if (item.count <= 0) return;

    const dropPoint = getMapDropPoint(info.point, item.mapWidth);
    if (!dropPoint) return;

    setPlacedItems((currentItems) => [
      ...currentItems,
      {
        id: `${item.type}-${Date.now()}-${currentItems.length}`,
        type: item.type,
        name: item.name,
        src: item.src,
        x: dropPoint.x,
        y: dropPoint.y,
        width: item.mapWidth,
      },
    ]);
    setInventory((currentInventory) =>
      currentInventory.map((inventoryItem) =>
        inventoryItem.type === item.type
          ? { ...inventoryItem, count: Math.max(0, inventoryItem.count - 1) }
          : inventoryItem,
      ),
    );
  };

  const handlePlacedItemDragEnd = (
    item: PlacedItem,
    _event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => {
    const dropPoint = getMapDropPoint(info.point, item.width);
    if (!dropPoint) return;

    setPlacedItems((currentItems) =>
      currentItems.map((placedItem) =>
        placedItem.id === item.id ? { ...placedItem, x: dropPoint.x, y: dropPoint.y } : placedItem,
      ),
    );
  };

  return (
    <main
      className="relative h-screen w-full overflow-hidden"
      onWheel={handleWheel}
      style={{
        background: "#112b1a",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#fff8d7",
      }}
    >
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

        <div
          style={{
            position: "absolute",
            left: "50%",
            top: "50%",
            width: "clamp(210px, 29vw, 430px)",
            transform: "translate(-64%, -28%)",
            zIndex: 3,
            pointerEvents: "none",
          }}
        >
          <motion.img
            className="pixelated drop-shadow-[18px_22px_0_rgba(0,0,0,0.28)]"
            src={mainStructureSrc}
            alt=""
            aria-hidden="true"
            initial={{ opacity: 0, y: 18, scale: 0.96 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            transition={{ duration: 0.36, ease: steppedEase(7) }}
            style={{
              display: "block",
              width: "100%",
              height: "auto",
              mixBlendMode: "screen",
            }}
          />
        </div>

        <div
          style={{
            position: "absolute",
            left: "50%",
            top: "50%",
            width: "clamp(92px, 12vw, 164px)",
            transform: "translate(42%, 62%)",
            zIndex: 4,
            pointerEvents: "none",
          }}
        >
          <div
            aria-hidden="true"
            style={{
              position: "absolute",
              left: "50%",
              top: "46%",
              width: "68%",
              aspectRatio: "1 / 1",
              transform: "translate(-50%, -50%)",
            }}
          >
            <motion.div
              className="bg-orange-500/30 blur-2xl"
              animate={{
                opacity: [0.58, 0.88, 0.66, 0.8, 0.58],
                scale: [0.92, 1.08, 0.98, 1.04, 0.92],
              }}
              transition={{ duration: 1.45, ease: "easeInOut", repeat: Infinity }}
              style={{
                width: "100%",
                height: "100%",
                clipPath: "circle(50% at 50% 50%)",
              }}
            />
          </div>
          <motion.img
            className="pixelated drop-shadow-[10px_14px_0_rgba(0,0,0,0.26)]"
            src={bonfireSrc}
            alt=""
            aria-hidden="true"
            initial={{ opacity: 0, scale: 0.86 }}
            animate={{
              opacity: 1,
              scale: [1, 1.035, 0.98, 1.02, 1],
              y: [0, -2, 1, -1, 0],
            }}
            transition={{
              opacity: { duration: 0.26, ease: steppedEase(5) },
              scale: { duration: 1.25, ease: "easeInOut", repeat: Infinity },
              y: { duration: 1.25, ease: "easeInOut", repeat: Infinity },
            }}
            style={{
              display: "block",
              position: "relative",
              width: "100%",
              height: "auto",
              mixBlendMode: "screen",
            }}
          />
        </div>

        {placedItems.map((item) => (
          <motion.img
            key={item.id}
            className="pixelated"
            src={item.src}
            alt={item.name}
            drag
            dragSnapToOrigin
            dragElastic={0}
            dragMomentum={false}
            onPointerDown={stopNestedDrag}
            onDragEnd={(event, info) => handlePlacedItemDragEnd(item, event, info)}
            whileDrag={{ scale: 1.05, zIndex: 12 }}
            style={{
              position: "absolute",
              left: item.x,
              top: item.y,
              width: item.width,
              height: "auto",
              cursor: "grab",
              filter: "drop-shadow(10px 14px 0 rgba(0,0,0,0.3))",
              mixBlendMode: "screen",
              touchAction: "none",
              zIndex: 8,
            }}
          />
        ))}
      </motion.div>

      <motion.header
        initial={{ opacity: 0, y: -14 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.32, ease: steppedEase(6) }}
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(14px, 2.2vw, 28px))",
          right: "clamp(14px, 2.2vw, 28px)",
          zIndex: 5,
          width: "min(calc(100vw - 150px), 720px)",
          border: "3px solid rgba(255, 248, 215, 0.76)",
          borderBottomColor: "rgba(55, 44, 35, 0.96)",
          borderRightColor: "rgba(55, 44, 35, 0.96)",
          background: "rgba(3, 10, 24, 0.72)",
          boxShadow:
            "0 0 0 2px rgba(0,0,0,0.68), 6px 6px 0 rgba(0,0,0,0.34), inset 0 0 18px rgba(116,247,161,0.1)",
          padding: "clamp(10px, 2vw, 14px)",
          backdropFilter: "blur(2px)",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "auto 1fr",
            alignItems: "center",
            gap: "clamp(10px, 2vw, 16px)",
          }}
        >
          <strong
            style={{
              color: "#ffd966",
              fontSize: "clamp(0.62rem, 1.45vw, 0.86rem)",
              lineHeight: 1.5,
              whiteSpace: "nowrap",
              textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
            }}
          >
            TOWN LEVEL {townLevel}
          </strong>
          <div style={{ minWidth: 0 }}>
            <div
              aria-label={`${currentCp.toLocaleString()} / ${nextLevelCp.toLocaleString()} CP`}
              role="progressbar"
              aria-valuemin={0}
              aria-valuemax={nextLevelCp}
              aria-valuenow={currentCp}
              style={{
                height: "14px",
                border: "2px solid rgba(116, 247, 161, 0.72)",
                background: "rgba(1, 8, 22, 0.8)",
                boxShadow: "inset 0 0 10px rgba(0,0,0,0.64)",
                overflow: "hidden",
              }}
            >
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${progress}%` }}
                transition={{ duration: 0.52, ease: steppedEase(8) }}
                style={{
                  height: "100%",
                  background:
                    "repeating-linear-gradient(90deg, #74f7a1 0, #74f7a1 10px, #39ff14 10px, #39ff14 20px)",
                  boxShadow: "0 0 12px rgba(116,247,161,0.72)",
                }}
              />
            </div>
            <p
              style={{
                margin: "7px 0 0",
                color: "#f4ecd0",
                fontFamily: '"DotGothic16", monospace',
                fontSize: "clamp(0.72rem, 1.45vw, 0.95rem)",
                lineHeight: 1.4,
                textAlign: "right",
              }}
            >
              {currentCp.toLocaleString()} / {nextLevelCp.toLocaleString()} CP
            </p>
          </div>
        </div>
      </motion.header>

      <motion.button
        type="button"
        whileHover={{ y: -2, scale: 1.02 }}
        whileTap={{ y: 2, scale: 0.98 }}
        onClick={() => onNavigate("/guild")}
        style={{
          position: "fixed",
          top: "calc(env(safe-area-inset-top, 0px) + clamp(14px, 2.2vw, 28px))",
          left: "clamp(14px, 2.2vw, 28px)",
          zIndex: 6,
          minHeight: "42px",
          border: "2px solid rgba(255, 217, 102, 0.78)",
          borderBottomColor: "rgba(96, 62, 22, 0.95)",
          borderRightColor: "rgba(96, 62, 22, 0.95)",
          background: "rgba(3, 10, 24, 0.74)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
          color: "#fff8d7",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.56rem, 1.3vw, 0.78rem)",
          lineHeight: 1.5,
          padding: "10px 12px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        &lt; BACK
      </motion.button>

      <motion.div
        ref={inventoryRef}
        initial={{ opacity: 0, y: 18 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.32, ease: steppedEase(6) }}
        style={{
          position: "fixed",
          left: "50%",
          bottom: "calc(env(safe-area-inset-bottom, 0px) + clamp(14px, 2vw, 24px))",
          transform: "translateX(-50%)",
          zIndex: 8,
          display: "flex",
          width: "min(calc(100vw - 160px), 620px)",
          minHeight: "104px",
          alignItems: "stretch",
          gap: "10px",
          overflowX: "auto",
          border: "3px solid rgba(255, 248, 215, 0.8)",
          borderBottomColor: "rgba(55, 44, 35, 0.98)",
          borderRightColor: "rgba(55, 44, 35, 0.98)",
          background: "rgba(3, 7, 14, 0.88)",
          boxShadow:
            "0 0 0 2px rgba(0,0,0,0.72), 6px 6px 0 rgba(0,0,0,0.4), inset 0 0 18px rgba(255,248,215,0.08)",
          padding: "10px",
          backdropFilter: "blur(2px)",
        }}
      >
        {inventory.map((item) => {
          const isAvailable = item.count > 0;

          return (
            <motion.div
              key={item.type}
              role="button"
              aria-label={`${item.name} inventory item. ${item.count} remaining.`}
              aria-disabled={!isAvailable}
              tabIndex={isAvailable ? 0 : -1}
              drag={isAvailable}
              dragSnapToOrigin
              dragElastic={0}
              dragMomentum={false}
              onPointerDown={stopNestedDrag}
              onDragEnd={(event, info) => handleInventoryDragEnd(item, event, info)}
              whileHover={
                isAvailable ? { y: -2, backgroundColor: "rgba(255, 217, 102, 0.12)" } : undefined
              }
              whileTap={isAvailable ? { y: 2, scale: 0.98 } : undefined}
              whileDrag={{ scale: 1.06, zIndex: 20 }}
              style={{
                position: "relative",
                display: "grid",
                gridTemplateRows: "52px auto",
                width: "116px",
                flex: "0 0 116px",
                alignItems: "center",
                justifyItems: "center",
                gap: "6px",
                border: "2px solid rgba(116, 247, 161, 0.62)",
                borderBottomColor: "rgba(24, 83, 45, 0.95)",
                borderRightColor: "rgba(24, 83, 45, 0.95)",
                background: isAvailable ? "rgba(1, 12, 24, 0.78)" : "rgba(18, 18, 18, 0.68)",
                boxShadow: "inset 0 0 12px rgba(0,0,0,0.62)",
                color: isAvailable ? "#fff8d7" : "rgba(255, 248, 215, 0.42)",
                cursor: isAvailable ? "grab" : "not-allowed",
                fontFamily: "inherit",
                fontSize: "0.52rem",
                lineHeight: 1.35,
                padding: "8px 7px",
                textAlign: "center",
                textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
                touchAction: "none",
              }}
            >
              <img
                className="pixelated"
                src={item.src}
                alt=""
                aria-hidden="true"
                draggable={false}
                style={{
                  display: "block",
                  maxWidth: "64px",
                  maxHeight: "52px",
                  opacity: isAvailable ? 1 : 0.38,
                  pointerEvents: "none",
                  filter: "drop-shadow(4px 5px 0 rgba(0,0,0,0.34))",
                  mixBlendMode: "screen",
                }}
              />
              <span>{item.name}</span>
              <span
                aria-hidden="true"
                style={{
                  position: "absolute",
                  right: "6px",
                  top: "6px",
                  minWidth: "24px",
                  border: "2px solid rgba(255, 217, 102, 0.78)",
                  background: isAvailable ? "rgba(62, 26, 8, 0.92)" : "rgba(8, 8, 8, 0.92)",
                  color: isAvailable ? "#ffd966" : "rgba(255, 248, 215, 0.46)",
                  fontSize: "0.52rem",
                  lineHeight: 1,
                  padding: "4px 5px",
                  boxShadow: "0 0 0 1px rgba(0,0,0,0.72)",
                }}
              >
                x{item.count}
              </span>
            </motion.div>
          );
        })}
      </motion.div>

      <div
        className="absolute right-6 z-[6] flex flex-col gap-3"
        style={{
          bottom: "calc(env(safe-area-inset-bottom, 0px) + 142px)",
        }}
      >
        <motion.button
          type="button"
          aria-label="Zoom in"
          whileHover={{ y: -2, scale: 1.04, backgroundColor: "rgba(255, 217, 102, 0.2)" }}
          whileTap={{ y: 2, scale: 0.96 }}
          onClick={() => handleZoom(ZOOM_STEP)}
          style={{
            width: "46px",
            height: "46px",
            border: "2px solid rgba(255, 217, 102, 0.78)",
            borderBottomColor: "rgba(96, 62, 22, 0.95)",
            borderRightColor: "rgba(96, 62, 22, 0.95)",
            background: "rgba(3, 10, 24, 0.82)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
            color: "#fff8d7",
            cursor: "pointer",
            fontFamily: "inherit",
            fontSize: "1.1rem",
            lineHeight: 1,
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          +
        </motion.button>
        <motion.button
          type="button"
          aria-label="Zoom out"
          whileHover={{ y: -2, scale: 1.04, backgroundColor: "rgba(255, 217, 102, 0.2)" }}
          whileTap={{ y: 2, scale: 0.96 }}
          onClick={() => handleZoom(-ZOOM_STEP)}
          style={{
            width: "46px",
            height: "46px",
            border: "2px solid rgba(255, 217, 102, 0.78)",
            borderBottomColor: "rgba(96, 62, 22, 0.95)",
            borderRightColor: "rgba(96, 62, 22, 0.95)",
            background: "rgba(3, 10, 24, 0.82)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.62), 5px 5px 0 rgba(0,0,0,0.34)",
            color: "#fff8d7",
            cursor: "pointer",
            fontFamily: "inherit",
            fontSize: "1.1rem",
            lineHeight: 1,
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          -
        </motion.button>
      </div>

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.08), rgba(0,0,0,0.08) 1px, transparent 1px, transparent 4px)",
          pointerEvents: "none",
          zIndex: 4,
        }}
      />
    </main>
  );
}
