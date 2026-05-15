import { useMotionValue, type PanInfo } from "framer-motion";
import { AUDIO_ASSETS } from "../features/audio/audioAssets";
import {
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
  type WheelEvent,
} from "react";
import { BackButton } from "./guild-town/BackButton";
import { BuildInventory } from "./guild-town/BuildInventory";
import { BuildingInfoPanel } from "./guild-town/BuildingInfoPanel";
import { TownMap } from "./guild-town/TownMap";
import { TownStatusHeader } from "./guild-town/TownStatusHeader";
import { MAX_SCALE, MIN_SCALE, STORE_ANIMATION_MS, INITIAL_INVENTORY } from "./guild-town/townData";
import { clampValue, getInventoryMapWidth, isPointInsideRect } from "./guild-town/townMath";
import type { InventoryItem, PlacedItem, ViewportSize } from "./guild-town/types";
import { ZoomControls } from "./guild-town/ZoomControls";
import { GuildBgm } from "./GuildBgm";

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
  const [viewport, setViewport] = useState<ViewportSize>({ width: 0, height: 0 });
  const [scale, setScale] = useState(1);
  const [inventory, setInventory] = useState<InventoryItem[]>(INITIAL_INVENTORY);
  const [placedItems, setPlacedItems] = useState<PlacedItem[]>([]);
  const [inventoryVisible, setInventoryVisible] = useState(true);
  const [selectedPlacedItemId, setSelectedPlacedItemId] = useState<string | null>(null);
  const [storingPlacedItemIds, setStoringPlacedItemIds] = useState<string[]>([]);
  const mapRef = useRef<HTMLDivElement>(null);
  const inventoryRef = useRef<HTMLDivElement>(null);
  const seededInitialBuildingsRef = useRef(false);
  const mapX = useMotionValue(0);
  const mapY = useMotionValue(0);
  const progress = Math.min(100, Math.max(0, (currentCp / nextLevelCp) * 100));
  const selectedPlacedItem =
    placedItems.find((placedItem) => placedItem.id === selectedPlacedItemId) ?? null;
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
    if (seededInitialBuildingsRef.current || viewport.width === 0 || viewport.height === 0) {
      return;
    }

    const tent = INITIAL_INVENTORY.find((item) => item.type === "tent");
    const bonfire = INITIAL_INVENTORY.find((item) => item.type === "bonfire");
    if (!tent || !bonfire) return;

    const tentWidth = getInventoryMapWidth(tent, viewport.width);
    const bonfireWidth = getInventoryMapWidth(bonfire, viewport.width);

    setPlacedItems([
      createPlacedItem(tent, {
        id: "initial-tent",
        src: mainStructureSrc,
        width: tentWidth,
        x: viewport.width - tentWidth * 0.64,
        y: viewport.height - tentWidth * 0.28,
      }),
      createPlacedItem(bonfire, {
        id: "initial-bonfire",
        src: bonfireSrc,
        width: bonfireWidth,
        x: viewport.width + bonfireWidth * 0.42,
        y: viewport.height + bonfireWidth * 0.62,
      }),
    ]);
    seededInitialBuildingsRef.current = true;
  }, [bonfireSrc, mainStructureSrc, viewport.height, viewport.width]);

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

    const itemWidth = getInventoryMapWidth(item, viewport.width);
    const dropPoint = getMapDropPoint(info.point, itemWidth);
    if (!dropPoint) return;

    setPlacedItems((currentItems) => [
      ...currentItems,
      createPlacedItem(item, {
        id: `${item.type}-${Date.now()}-${currentItems.length}`,
        src: item.src,
        width: itemWidth,
        x: dropPoint.x,
        y: dropPoint.y,
      }),
    ]);
    setInventoryCount(item.type, -1);
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
    setSelectedPlacedItemId(item.id);
  };

  const handleStorePlacedItem = (item: PlacedItem) => {
    if (storingPlacedItemIds.includes(item.id)) return;

    setStoringPlacedItemIds((currentIds) => [...currentIds, item.id]);
    setSelectedPlacedItemId(null);

    window.setTimeout(() => {
      setPlacedItems((currentItems) =>
        currentItems.filter((placedItem) => placedItem.id !== item.id),
      );
      setInventoryCount(item.type, 1);
      setStoringPlacedItemIds((currentIds) =>
        currentIds.filter((storingItemId) => storingItemId !== item.id),
      );
    }, STORE_ANIMATION_MS);
  };

  const setInventoryCount = (type: InventoryItem["type"], delta: number) => {
    setInventory((currentInventory) =>
      currentInventory.map((inventoryItem) =>
        inventoryItem.type === type
          ? { ...inventoryItem, count: Math.max(0, inventoryItem.count + delta) }
          : inventoryItem,
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
      <GuildBgm src={AUDIO_ASSETS.bgm.guildTown} />

      <TownMap
        baseSrc={baseSrc}
        dragConstraints={dragConstraints}
        mapRef={mapRef}
        mapX={mapX}
        mapY={mapY}
        onMoveItem={handlePlacedItemDragEnd}
        onSelectItem={setSelectedPlacedItemId}
        onStoreItem={handleStorePlacedItem}
        placedItems={placedItems}
        scale={scale}
        selectedPlacedItemId={selectedPlacedItemId}
        stopNestedDrag={stopNestedDrag}
        storingPlacedItemIds={storingPlacedItemIds}
      />

      <TownStatusHeader
        currentCp={currentCp}
        nextLevelCp={nextLevelCp}
        progress={progress}
        townLevel={townLevel}
      />
      <BackButton onNavigate={onNavigate} />
      <BuildInventory
        inventory={inventory}
        inventoryRef={inventoryRef}
        onDragEnd={handleInventoryDragEnd}
        onToggleVisible={() => setInventoryVisible((currentVisible) => !currentVisible)}
        stopNestedDrag={stopNestedDrag}
        visible={inventoryVisible}
      />
      <BuildingInfoPanel item={selectedPlacedItem} onClose={() => setSelectedPlacedItemId(null)} />
      <ZoomControls onZoom={handleZoom} />

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

function createPlacedItem(
  item: InventoryItem,
  placement: { id: string; src: string; width: number; x: number; y: number },
): PlacedItem {
  return {
    id: placement.id,
    type: item.type,
    name: item.name,
    title: item.title,
    description: item.description,
    src: placement.src,
    x: placement.x,
    y: placement.y,
    width: placement.width,
  };
}
