import { useMotionValue, type PanInfo } from "framer-motion";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { fetchMyGuild } from "../../features/guild/api";
import { findGuildBySlug } from "../../features/guild/guildMaster";
import { getSelectedGuildSlug } from "../../features/guild/membership";
import {
  useEffect,
  useRef,
  useState,
  type PointerEvent as ReactPointerEvent,
  type WheelEvent,
} from "react";
import { GuildBgm } from "../shared/GuildBgm";
import { BackButton } from "./BackButton";
import { BuildInventory } from "./BuildInventory";
import { BuildingInfoPanel } from "./BuildingInfoPanel";
import { TownMap } from "./TownMap";
import { TownStatusHeader } from "./TownStatusHeader";
import {
  BUILDING_MASTERS,
  MAX_SCALE,
  MIN_SCALE,
  STORE_ANIMATION_MS,
  INITIAL_INVENTORY,
} from "./townData";
import { clampValue, getInventoryMapWidth, isPointInsideRect } from "./townMath";
import { GUILD_LANGUAGES } from "./types";
import type {
  BuildingMaster,
  GuildSpLanguage,
  InventoryItem,
  PlacedItem,
  UserInventoryState,
  ViewportSize,
} from "./types";
import { ZoomControls } from "./ZoomControls";

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
  const [placedItems, setPlacedItems] = useState<PlacedItem[]>([]);
  const [inventoryVisible, setInventoryVisible] = useState(true);
  const [selectedPlacedItemId, setSelectedPlacedItemId] = useState<string | null>(null);
  const [storingPlacedItemIds, setStoringPlacedItemIds] = useState<string[]>([]);
  const [buildFeedbackMessage, setBuildFeedbackMessage] = useState<string | null>(null);
  const [userCp, setUserCp] = useState(1200);
  const [currentGuildLanguage, setCurrentGuildLanguage] = useState<GuildSpLanguage>(() =>
    getCurrentGuildLanguage(),
  );
  const [userGuildSp, setUserGuildSp] = useState(500);
  const [userInventory, setUserInventory] = useState<UserInventoryState[]>(
    BUILDING_MASTERS.map((building) => ({ buildingId: building.id, count: 0 })),
  );
  const mapRef = useRef<HTMLDivElement>(null);
  const inventoryRef = useRef<HTMLDivElement>(null);
  const seededInitialBuildingsRef = useRef(false);
  const mapX = useMotionValue(0);
  const mapY = useMotionValue(0);
  const progress = Math.min(100, Math.max(0, (currentCp / nextLevelCp) * 100));
  const currentGuildLevel = 3;
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
    let isMounted = true;

    fetchMyGuild()
      .then((data) => {
        if (!isMounted || !data?.guild || !isGuildSpLanguage(data.guild.name)) {
          return;
        }

        setCurrentGuildLanguage(data.guild.name);
      })
      .catch((error) => {
        if (isMounted) {
          console.error("failed to fetch current guild for guild town", error);
        }
      });

    return () => {
      isMounted = false;
    };
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
      setStoringPlacedItemIds((currentIds) =>
        currentIds.filter((storingItemId) => storingItemId !== item.id),
      );
      if (item.buildingId) {
        setUserInventory((currentInventory) =>
          currentInventory.map((inventoryItem) =>
            inventoryItem.buildingId === item.buildingId
              ? { ...inventoryItem, count: inventoryItem.count + 1 }
              : inventoryItem,
          ),
        );
      }
    }, STORE_ANIMATION_MS);
  };

  const handleBuyBuilding = (building: BuildingMaster) => {
    const firstLevel = building.levels[0];
    const canBuy =
      currentGuildLevel >= building.requiredGuildLevel &&
      userCp >= firstLevel.upgradeCostCp &&
      userGuildSp >= firstLevel.upgradeCostSp;

    if (!canBuy) {
      const failureMessage = getBuyFailureMessage({
        building,
        currentGuildLanguage,
        currentGuildLevel,
        firstLevel,
        userCp,
        userGuildSp,
      });

      console.debug("failed to buy guild town building", {
        buildingId: building.id,
        currentGuildLanguage,
        currentGuildLevel,
        requiredGuildLevel: building.requiredGuildLevel,
        requiredCp: firstLevel.upgradeCostCp,
        requiredGuildSp: firstLevel.upgradeCostSp,
        userCp,
        userGuildSp,
      });
      setBuildFeedbackMessage(failureMessage);
      return;
    }

    setUserCp((currentValue) => currentValue - firstLevel.upgradeCostCp);
    setUserGuildSp((currentValue) => currentValue - firstLevel.upgradeCostSp);
    setUserInventory((currentInventory) =>
      currentInventory.map((inventoryItem) =>
        inventoryItem.buildingId === building.id
          ? { ...inventoryItem, count: inventoryItem.count + 1 }
          : inventoryItem,
      ),
    );
  };

  const handleDeployBuilding = (building: BuildingMaster) => {
    const inventoryItem = userInventory.find((item) => item.buildingId === building.id);
    if (
      !inventoryItem ||
      inventoryItem.count <= 0 ||
      viewport.width === 0 ||
      viewport.height === 0
    ) {
      return;
    }

    const width = getBuildingMapWidth(viewport.width);
    const mapWidth = viewport.width * 2;
    const mapHeight = viewport.height * 2;
    const x = clampValue(
      (-mapX.get() + viewport.width / 2) / scale - width / 2,
      0,
      mapWidth - width,
    );
    const y = clampValue(
      (-mapY.get() + viewport.height / 2) / scale - width / 2,
      0,
      mapHeight - width,
    );
    const placedItemId = `${building.id}-${Date.now()}`;

    setUserInventory((currentInventory) =>
      currentInventory.map((item) =>
        item.buildingId === building.id ? { ...item, count: Math.max(0, item.count - 1) } : item,
      ),
    );
    setPlacedItems((currentItems) => [
      ...currentItems,
      createPlacedBuildingItem(building, {
        id: placedItemId,
        width,
        x,
        y,
      }),
    ]);
    setSelectedPlacedItemId(placedItemId);
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
        currentGuildLevel={currentGuildLevel}
        currentGuildLanguage={currentGuildLanguage}
        inventory={userInventory}
        inventoryRef={inventoryRef}
        onBuyBuilding={handleBuyBuilding}
        onDeployBuilding={handleDeployBuilding}
        onToggleVisible={() => setInventoryVisible((currentVisible) => !currentVisible)}
        stopNestedDrag={stopNestedDrag}
        userCp={userCp}
        userGuildSp={userGuildSp}
        visible={inventoryVisible}
      />
      {buildFeedbackMessage && (
        <p
          role="alert"
          style={{
            position: "fixed",
            bottom: "calc(env(safe-area-inset-bottom, 0px) + 22px)",
            left: "50%",
            zIndex: 12,
            margin: 0,
            maxWidth: "min(720px, calc(100vw - 32px))",
            transform: "translateX(-50%)",
            border: "2px solid rgba(255, 77, 109, 0.86)",
            borderBottomColor: "rgba(118, 31, 49, 0.95)",
            borderRightColor: "rgba(118, 31, 49, 0.95)",
            background: "rgba(18, 8, 14, 0.94)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
            color: "#ff9aae",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.92rem",
            lineHeight: 1.45,
            padding: "10px 14px",
            textAlign: "center",
            textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
          }}
        >
          {buildFeedbackMessage}
        </p>
      )}
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

function getBuildingMapWidth(viewportWidth: number) {
  return clampValue(viewportWidth * 0.14, 112, 220);
}

function getBuyFailureMessage({
  building,
  currentGuildLanguage,
  currentGuildLevel,
  firstLevel,
  userCp,
  userGuildSp,
}: {
  building: BuildingMaster;
  currentGuildLanguage: GuildSpLanguage;
  currentGuildLevel: number;
  firstLevel: BuildingMaster["levels"][number];
  userCp: number;
  userGuildSp: number;
}) {
  if (currentGuildLevel < building.requiredGuildLevel) {
    return `ギルドLV.${building.requiredGuildLevel}で解放されます。`;
  }

  if (userCp < firstLevel.upgradeCostCp) {
    return `${building.name}の購入には ${firstLevel.upgradeCostCp.toLocaleString()} CP が必要です。`;
  }

  if (userGuildSp < firstLevel.upgradeCostSp) {
    return `${building.name}の購入には ${firstLevel.upgradeCostSp.toLocaleString()} ${currentGuildLanguage}-SP が必要です。`;
  }

  return `${building.name}を購入できませんでした。`;
}

function getCurrentGuildLanguage(): GuildSpLanguage {
  const selectedGuild = findGuildBySlug(getSelectedGuildSlug());
  if (!selectedGuild || !isGuildSpLanguage(selectedGuild.name)) {
    return "Common";
  }

  return selectedGuild.name;
}

function isGuildSpLanguage(language: string): language is GuildSpLanguage {
  return GUILD_LANGUAGES.includes(language as GuildSpLanguage);
}

function createPlacedBuildingItem(
  building: BuildingMaster,
  placement: { id: string; width: number; x: number; y: number },
): PlacedItem {
  return {
    id: placement.id,
    type: building.id,
    buildingId: building.id,
    name: building.name,
    title: building.name,
    description: building.description,
    src: building.previewSrc ?? "/build-items/plasma-capacitor.jpeg",
    x: placement.x,
    y: placement.y,
    width: placement.width,
  };
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
