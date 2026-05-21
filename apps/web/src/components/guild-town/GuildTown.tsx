import { motion, useMotionValue, type PanInfo } from "framer-motion";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { fetchMyGuild } from "../../features/guild/api";
import { findGuildBySlug } from "../../features/guild/guildMaster";
import { getSelectedGuildSlug } from "../../features/guild/membership";
import {
  buyBuilding,
  deployBuilding,
  fetchGuildTownStatus,
  guildSpBalance,
  saveGuildTownPlacements,
  upgradeBuilding,
  type GuildTownStatus,
} from "../../features/guild-town/api";
import {
  useEffect,
  useMemo,
  useRef,
  useState,
  type MouseEvent as ReactMouseEvent,
  type PointerEvent as ReactPointerEvent,
  type WheelEvent,
} from "react";
import { GuildBgm } from "../shared/GuildBgm";
import { BackButton } from "./BackButton";
import { BuildInventory } from "./BuildInventory";
import { BuildingInfoPanel } from "./BuildingInfoPanel";
import { TownMap, type DeploymentPreview } from "./TownMap";
import { TownStatusHeader } from "./TownStatusHeader";
import { steppedEase } from "../../lib/animationUtils";
import { BUILDING_MASTERS, MAX_SCALE, MIN_SCALE, STORE_ANIMATION_MS } from "./townData";
import { clampValue, isPointInsideRect } from "./townMath";
import { isTownRectUnlocked } from "./townUnlock";
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

interface DeploymentDraft {
  x: number;
  y: number;
}

export function GuildTown({
  onNavigate,
  townLevel = 1,
  currentCp: initialCurrentCp = 2500,
  nextLevelCp: initialNextLevelCp = 10000,
  baseSrc = "/town/glassfield.png",
}: GuildTownProps) {
  const [viewport, setViewport] = useState<ViewportSize>({ width: 0, height: 0 });
  const [scale, setScale] = useState(1);
  const [placedItems, setPlacedItems] = useState<PlacedItem[]>([]);
  const [availableItems, setAvailableItems] = useState<InventoryItem[]>([]);
  const [inventoryVisible, setInventoryVisible] = useState(true);
  const [selectedPlacedItemId, setSelectedPlacedItemId] = useState<string | null>(null);
  const [deployingBuildingId, setDeployingBuildingId] = useState<string | null>(null);
  const [deploymentDraft, setDeploymentDraft] = useState<DeploymentDraft | null>(null);
  const [newlyDeployedItemId, setNewlyDeployedItemId] = useState<string | null>(null);
  const [unlockClearingLevel, setUnlockClearingLevel] = useState<number | null>(null);
  const [storingPlacedItemIds, setStoringPlacedItemIds] = useState<string[]>([]);
  const [buildFeedbackMessage, setBuildFeedbackMessage] = useState<string | null>(null);
  const [loadErrorMessage, setLoadErrorMessage] = useState<string | null>(null);
  const [isTownLoading, setIsTownLoading] = useState(true);
  const [userCp, setUserCp] = useState(initialCurrentCp);
  const [townNextLevelCp, setTownNextLevelCp] = useState(initialNextLevelCp);
  const [currentGuildLevel, setCurrentGuildLevel] = useState(townLevel);
  const [currentGuildLanguage, setCurrentGuildLanguage] = useState<GuildSpLanguage>(() =>
    getCurrentGuildLanguage(),
  );
  const [userSpMap, setUserSpMap] = useState<Record<string, number>>({});
  const [userInventory, setUserInventory] = useState<UserInventoryState[]>(
    BUILDING_MASTERS.map((building) => ({ buildingId: building.id, count: 0 })),
  );
  const mapRef = useRef<HTMLDivElement>(null);
  const inventoryRef = useRef<HTMLDivElement>(null);
  const previousGuildLevelRef = useRef(currentGuildLevel);
  const mapX = useMotionValue(0);
  const mapY = useMotionValue(0);
  const progress = Math.min(100, Math.max(0, (userCp / townNextLevelCp) * 100));
  const selectedPlacedItem =
    placedItems.find((placedItem) => placedItem.id === selectedPlacedItemId) ?? null;
  const userGuildSp = useMemo(
    () => guildSpBalance(userSpMap, currentGuildLanguage),
    [currentGuildLanguage, userSpMap],
  );
  const inventoryBuildingCatalog = useMemo(
    () => availableItems.map(toInventoryBuildingMaster),
    [availableItems],
  );
  const deployingBuilding =
    deployingBuildingId === null
      ? null
      : (BUILDING_MASTERS.find((building) => building.id === deployingBuildingId) ??
        inventoryBuildingCatalog.find((building) => building.id === deployingBuildingId) ??
        null);
  const deploymentPreview = getDeploymentPreview();
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

    fetchGuildTownStatus()
      .then((status) => {
        if (!isMounted) return;
        applyGuildTownStatus(status);
        setLoadErrorMessage(null);
      })
      .catch((error) => {
        if (!isMounted) return;
        console.error("failed to fetch guild town status", error);
        setLoadErrorMessage("ギルドタウンの読み込みに失敗しました。");
      })
      .finally(() => {
        if (isMounted) {
          setIsTownLoading(false);
        }
      });

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
    const previousGuildLevel = previousGuildLevelRef.current;
    previousGuildLevelRef.current = currentGuildLevel;

    if (currentGuildLevel <= previousGuildLevel) {
      return;
    }

    setUnlockClearingLevel(currentGuildLevel);
    const clearUnlockAnimation = window.setTimeout(() => {
      setUnlockClearingLevel(null);
    }, 1900);

    return () => window.clearTimeout(clearUnlockAnimation);
  }, [currentGuildLevel]);

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

  useEffect(() => {
    if (deployingBuildingId === null) return;

    const handleDeployKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.preventDefault();
        cancelDeployMode();
        return;
      }

      if (event.key === "Enter") {
        event.preventDefault();
        void commitDeployment();
        return;
      }

      const arrowDeltaByKey: Record<string, { x: number; y: number }> = {
        ArrowDown: { x: 0, y: 1 },
        ArrowLeft: { x: -1, y: 0 },
        ArrowRight: { x: 1, y: 0 },
        ArrowUp: { x: 0, y: -1 },
      };
      const arrowDelta = arrowDeltaByKey[event.key];
      if (!arrowDelta) {
        return;
      }

      event.preventDefault();
      moveDeploymentDraftBy(arrowDelta.x, arrowDelta.y, event.shiftKey ? 48 : 16);
    };

    window.addEventListener("keydown", handleDeployKeyDown);
    return () => window.removeEventListener("keydown", handleDeployKeyDown);
  });

  useEffect(() => {
    if (newlyDeployedItemId === null) return;

    const resetAnimationTarget = window.setTimeout(() => {
      setNewlyDeployedItemId(null);
    }, 900);

    return () => window.clearTimeout(resetAnimationTarget);
  }, [newlyDeployedItemId]);

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

  function cancelDeployMode() {
    setDeployingBuildingId(null);
    setDeploymentDraft(null);
    setBuildFeedbackMessage(null);
  }

  const handleTownContextMenu = (event: ReactMouseEvent<HTMLElement>) => {
    if (deployingBuildingId === null) return;

    event.preventDefault();
    cancelDeployMode();
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

  function getMapSize() {
    const mapElement = mapRef.current;

    return {
      height: mapElement?.offsetHeight ?? viewport.height * 2,
      width: mapElement?.offsetWidth ?? viewport.width * 2,
    };
  }

  function getInitialDeploymentDraft(itemWidth: number): DeploymentDraft {
    const mapSize = getMapSize();
    const x = clampValue(
      (-mapX.get() + viewport.width / 2) / scale - itemWidth / 2,
      0,
      Math.max(0, mapSize.width - itemWidth),
    );
    const y = clampValue(
      (-mapY.get() + viewport.height / 2) / scale - itemWidth / 2,
      0,
      Math.max(0, mapSize.height - itemWidth),
    );

    return { x, y };
  }

  function isDeploymentDraftUnlocked(draft: DeploymentDraft, itemWidth: number) {
    const mapSize = getMapSize();

    return isTownRectUnlocked(
      getBuildingUnlockRect({
        itemWidth,
        mapHeight: mapSize.height,
        mapWidth: mapSize.width,
        x: draft.x,
        y: draft.y,
      }),
      currentGuildLevel,
    );
  }

  function moveDeploymentDraftBy(deltaX: number, deltaY: number, step: number) {
    if (!deploymentDraft || viewport.width === 0) return;

    const width = getBuildingMapWidth(viewport.width);
    const mapSize = getMapSize();
    const nextDraft = {
      x: clampValue(deploymentDraft.x + deltaX * step, 0, Math.max(0, mapSize.width - width)),
      y: clampValue(deploymentDraft.y + deltaY * step, 0, Math.max(0, mapSize.height - width)),
    };

    setDeploymentDraft(nextDraft);
    if (!isDeploymentDraftUnlocked(nextDraft, width)) {
      setBuildFeedbackMessage(getLockedDeploymentMessage(currentGuildLevel));
    } else {
      setBuildFeedbackMessage(null);
    }
  }

  function getDeploymentPreview(): DeploymentPreview | null {
    if (!deployingBuilding || !deploymentDraft || viewport.width === 0 || viewport.height === 0) {
      return null;
    }

    const width = getBuildingMapWidth(viewport.width);
    const mapSize = getMapSize();

    return {
      id: deployingBuilding.id,
      isUnlocked: isTownRectUnlocked(
        getBuildingUnlockRect({
          itemWidth: width,
          mapHeight: mapSize.height,
          mapWidth: mapSize.width,
          x: deploymentDraft.x,
          y: deploymentDraft.y,
        }),
        currentGuildLevel,
      ),
      name: deployingBuilding.name,
      src: deployingBuilding.previewSrc ?? "/build-items/plasma-capacitor.jpeg",
      width,
      x: deploymentDraft.x,
      y: deploymentDraft.y,
    };
  }

  const handlePlacedItemDragEnd = (
    item: PlacedItem,
    _event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => {
    const dropPoint = getMapDropPoint(info.point, item.width);
    if (!dropPoint) return;

    if (!isDeploymentDraftUnlocked(dropPoint, item.width)) {
      setBuildFeedbackMessage(getLockedDeploymentMessage(currentGuildLevel));
      return;
    }

    const nextItems = placedItems.map((placedItem) =>
      placedItem.id === item.id ? { ...placedItem, x: dropPoint.x, y: dropPoint.y } : placedItem,
    );
    setPlacedItems(nextItems);
    setSelectedPlacedItemId(item.id);
    setBuildFeedbackMessage(null);
    persistPlacements(nextItems);
  };

  const handleStorePlacedItem = (item: PlacedItem) => {
    if (storingPlacedItemIds.includes(item.id)) return;

    setStoringPlacedItemIds((currentIds) => [...currentIds, item.id]);
    setSelectedPlacedItemId(null);

    window.setTimeout(() => {
      const nextItems = placedItems.filter((placedItem) => placedItem.id !== item.id);
      setPlacedItems(nextItems);
      setStoringPlacedItemIds((currentIds) =>
        currentIds.filter((storingItemId) => storingItemId !== item.id),
      );
      persistPlacements(nextItems);
    }, STORE_ANIMATION_MS);
  };

  const handleBuyBuilding = async (building: BuildingMaster) => {
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

    try {
      await buyBuilding(building.id);
      await reloadGuildTownStatus();
      setBuildFeedbackMessage("");
    } catch (error) {
      console.error("failed to buy guild town building", error);
      setBuildFeedbackMessage("購入APIはまだバックエンドに実装されていません。");
    }
  };

  const handleBeginDeployBuilding = (building: BuildingMaster) => {
    const inventoryItem = userInventory.find((item) => item.buildingId === building.id);
    if (!inventoryItem || inventoryItem.count <= 0) {
      setBuildFeedbackMessage("配置できる建物がインベントリにありません。");
      return;
    }

    setSelectedPlacedItemId(null);
    setDeployingBuildingId(building.id);
    setDeploymentDraft(getInitialDeploymentDraft(getBuildingMapWidth(viewport.width)));
    setBuildFeedbackMessage(null);
  };

  const handleMoveDeploymentPreview = (
    _event: MouseEvent | TouchEvent | PointerEvent,
    info: PanInfo,
  ) => {
    if (!deploymentDraft || viewport.width === 0) return;

    const width = getBuildingMapWidth(viewport.width);
    const dropPoint = getMapDropPoint(info.point, width);
    if (!dropPoint) return;

    setDeploymentDraft(dropPoint);
    if (!isDeploymentDraftUnlocked(dropPoint, width)) {
      setBuildFeedbackMessage(getLockedDeploymentMessage(currentGuildLevel));
    } else {
      setBuildFeedbackMessage(null);
    }
  };

  async function commitDeployment() {
    if (!deployingBuildingId || !deploymentDraft) return;

    const building =
      BUILDING_MASTERS.find((buildingMaster) => buildingMaster.id === deployingBuildingId) ??
      inventoryBuildingCatalog.find(
        (inventoryBuilding) => inventoryBuilding.id === deployingBuildingId,
      );
    const inventoryItem = userInventory.find((item) => item.buildingId === deployingBuildingId);

    if (!building || !inventoryItem || inventoryItem.count <= 0 || viewport.width === 0) {
      setBuildFeedbackMessage("配置できる建物がインベントリにありません。");
      cancelDeployMode();
      return;
    }

    const width = getBuildingMapWidth(viewport.width);
    if (!isDeploymentDraftUnlocked(deploymentDraft, width)) {
      setBuildFeedbackMessage(getLockedDeploymentMessage(currentGuildLevel));
      return;
    }

    const placementX = deploymentDraft.x;
    const placementY = deploymentDraft.y;
    const placedItemId = `local-${building.id}-${Date.now()}`;

    const nextItems = [
      ...placedItems,
      createPlacedBuildingItem(building, {
        id: placedItemId,
        width,
        x: placementX,
        y: placementY,
      }),
    ];
    setUserInventory((currentInventory) =>
      currentInventory.map((item) =>
        item.buildingId === building.id ? { ...item, count: Math.max(0, item.count - 1) } : item,
      ),
    );
    setPlacedItems(nextItems);
    setSelectedPlacedItemId(placedItemId);
    try {
      const status = await deployBuilding({ placements: nextItems });
      applyGuildTownStatus(status);
      const deployedItem = findDeployedItem(
        status.placedItems,
        building.id,
        placementX,
        placementY,
      );
      setNewlyDeployedItemId(deployedItem?.id ?? placedItemId);
      setSelectedPlacedItemId(deployedItem?.id ?? placedItemId);
      setDeployingBuildingId(null);
      setDeploymentDraft(null);
      setBuildFeedbackMessage("");
    } catch (error) {
      console.error("failed to deploy guild town building", error);
      setBuildFeedbackMessage("配置の保存に失敗しました。インベントリ数を確認してください。");
      setDeployingBuildingId(null);
      setDeploymentDraft(null);
      await reloadGuildTownStatus();
    }
  }

  const handleUpgradeBuilding = async (placedItemId: string) => {
    const placedItem = placedItems.find((item) => item.id === placedItemId);
    const building = placedItem?.buildingId
      ? BUILDING_MASTERS.find((buildingMaster) => buildingMaster.id === placedItem.buildingId)
      : undefined;
    if (!placedItem || !building) return;

    const currentLevel = Math.min(Math.max(placedItem.level, 1), building.levels.length);
    const nextLevel = building.levels[currentLevel];
    if (!nextLevel) return;

    if (userCp < nextLevel.upgradeCostCp || userGuildSp < nextLevel.upgradeCostSp) {
      setBuildFeedbackMessage(
        `${building.name}の強化には ${nextLevel.upgradeCostCp.toLocaleString()} CP / ${nextLevel.upgradeCostSp.toLocaleString()} ${building.targetSpLanguage}-SP が必要です。`,
      );
      return;
    }

    try {
      await upgradeBuilding(placedItemId);
      await reloadGuildTownStatus();
      setBuildFeedbackMessage("");
    } catch (error) {
      console.error("failed to upgrade guild town building", error);
      setBuildFeedbackMessage("強化APIはまだバックエンドに実装されていません。");
    }
  };

  const applyGuildTownStatus = (status: GuildTownStatus) => {
    setAvailableItems(status.availableItems);
    setCurrentGuildLevel(status.guildLevel);
    setPlacedItems(status.placedItems);
    setTownNextLevelCp(status.nextLevelCp);
    setUserCp(status.currentCp);
    setUserInventory(status.userInventory);
    setUserSpMap(status.userSpMap);
  };

  const reloadGuildTownStatus = async () => {
    const status = await fetchGuildTownStatus();
    applyGuildTownStatus(status);
  };

  const persistPlacements = (nextItems: PlacedItem[]) => {
    saveGuildTownPlacements({ placements: nextItems })
      .then(applyGuildTownStatus)
      .catch((error) => {
        console.error("failed to save guild town placements", error);
        setBuildFeedbackMessage("配置の保存に失敗しました。");
      });
  };

  return (
    <main
      className="relative h-screen w-full overflow-hidden"
      onContextMenu={handleTownContextMenu}
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
        currentGuildLevel={currentGuildLevel}
        deploymentPreview={deploymentPreview}
        newlyDeployedItemId={newlyDeployedItemId}
        unlockClearingLevel={unlockClearingLevel}
        onMoveItem={handlePlacedItemDragEnd}
        onCancelDeployment={cancelDeployMode}
        onCommitDeployment={() => void commitDeployment()}
        onMoveDeploymentPreview={handleMoveDeploymentPreview}
        onSelectItem={setSelectedPlacedItemId}
        onStoreItem={handleStorePlacedItem}
        placedItems={placedItems}
        scale={scale}
        selectedPlacedItemId={selectedPlacedItemId}
        stopNestedDrag={stopNestedDrag}
        storingPlacedItemIds={storingPlacedItemIds}
      />

      <TownStatusHeader
        currentCp={userCp}
        nextLevelCp={townNextLevelCp}
        progress={progress}
        townLevel={currentGuildLevel}
      />
      <BackButton onNavigate={onNavigate} />
      <BuildInventory
        currentGuildLevel={currentGuildLevel}
        currentGuildLanguage={currentGuildLanguage}
        inventory={userInventory}
        inventoryBuildings={inventoryBuildingCatalog}
        inventoryRef={inventoryRef}
        onBuyBuilding={handleBuyBuilding}
        onDeployBuilding={handleBeginDeployBuilding}
        onToggleVisible={() => setInventoryVisible((currentVisible) => !currentVisible)}
        stopNestedDrag={stopNestedDrag}
        userCp={userCp}
        userGuildSp={userGuildSp}
        visible={inventoryVisible}
      />
      {isTownLoading && <GuildTownLoadingOverlay />}
      {deployingBuilding && (
        <DeployModePanel buildingName={deployingBuilding.name} onCancel={cancelDeployMode} />
      )}
      {loadErrorMessage && <GuildTownToast message={loadErrorMessage} />}
      {buildFeedbackMessage && <GuildTownToast message={buildFeedbackMessage} />}
      <BuildingInfoPanel
        item={selectedPlacedItem}
        onClose={() => setSelectedPlacedItemId(null)}
        onUpgradeBuilding={handleUpgradeBuilding}
        userCp={userCp}
        userGuildSp={userGuildSp}
      />
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

function getBuildingUnlockRect({
  itemWidth,
  mapHeight,
  mapWidth,
  x,
  y,
}: {
  itemWidth: number;
  mapHeight: number;
  mapWidth: number;
  x: number;
  y: number;
}) {
  return {
    height: itemWidth * 0.72,
    mapHeight,
    mapWidth,
    width: itemWidth * 0.82,
    x: x + itemWidth * 0.09,
    y: y + itemWidth * 0.18,
  };
}

function getLockedDeploymentMessage(currentGuildLevel: number) {
  return `このエリアはまだロックされています。ギルドLV.${currentGuildLevel + 1}以上に上げて解放すると配置できます。`;
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
    level: 1,
    name: building.name,
    title: building.name,
    description: building.description,
    src: building.previewSrc ?? "/build-items/plasma-capacitor.jpeg",
    x: placement.x,
    y: placement.y,
    width: placement.width,
  };
}

function findDeployedItem(
  placedItems: PlacedItem[],
  buildingId: string,
  x: number,
  y: number,
): PlacedItem | null {
  return (
    placedItems.find(
      (item) =>
        item.buildingId === buildingId && Math.abs(item.x - x) < 1 && Math.abs(item.y - y) < 1,
    ) ?? null
  );
}

function toInventoryBuildingMaster(item: InventoryItem): BuildingMaster {
  return {
    name: item.name,
    description: item.description,
    id: item.type,
    previewSrc: item.src,
    requiredGuildLevel: 1,
    buffType: "core",
    targetSpLanguage: "Common",
    levels: [{ level: 1, upgradeCostCp: 0, upgradeCostSp: 0, buffValue: 0 }],
  };
}

function GuildTownLoadingOverlay() {
  return (
    <div
      role="status"
      aria-live="polite"
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 11,
        display: "grid",
        placeItems: "center",
        background: "rgba(3, 7, 14, 0.58)",
        color: "#74f7a1",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        fontSize: "0.72rem",
        letterSpacing: 0,
        textShadow: "0 0 12px rgba(116,247,161,0.78), 2px 2px 0 rgba(0,0,0,0.82)",
      }}
    >
      <span
        style={{
          border: "2px solid rgba(116, 247, 161, 0.78)",
          borderBottomColor: "rgba(24, 83, 45, 0.95)",
          borderRightColor: "rgba(24, 83, 45, 0.95)",
          background: "rgba(1, 12, 24, 0.9)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34)",
          padding: "14px 18px",
        }}
      >
        SYNCING GUILD TOWN...
      </span>
    </div>
  );
}

function DeployModePanel({
  buildingName,
  onCancel,
}: {
  buildingName: string;
  onCancel: () => void;
}) {
  return (
    <motion.div
      role="status"
      aria-live="polite"
      initial={{ opacity: 0, y: -10, scale: 0.96 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: -8, scale: 0.96 }}
      transition={{ duration: 0.18, ease: steppedEase(5) }}
      style={{
        position: "fixed",
        left: "50%",
        top: "calc(env(safe-area-inset-top, 0px) + 92px)",
        zIndex: 13,
        display: "grid",
        gridTemplateColumns: "minmax(0, 1fr) auto",
        alignItems: "center",
        gap: "12px",
        maxWidth: "min(720px, calc(100vw - 32px))",
        transform: "translateX(-50%)",
        border: "3px solid rgba(116, 247, 161, 0.86)",
        borderBottomColor: "rgba(24, 83, 45, 0.95)",
        borderRightColor: "rgba(24, 83, 45, 0.95)",
        background: "rgba(1, 12, 24, 0.94)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.68), 4px 4px 0 rgba(0,0,0,0.34), 0 0 22px rgba(116,247,161,0.24)",
        color: "#fff8d7",
        fontFamily: '"DotGothic16", monospace',
        fontSize: "0.9rem",
        lineHeight: 1.35,
        padding: "10px 12px",
        textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
      }}
    >
      <span style={{ minWidth: 0, overflowWrap: "anywhere" }}>
        {buildingName} をドラッグまたは矢印キーで動かし、BUILD / Enter で配置してください
      </span>
      <button
        type="button"
        onClick={onCancel}
        style={{
          minHeight: "34px",
          border: "2px solid rgba(255, 217, 102, 0.86)",
          borderBottomColor: "rgba(96, 62, 22, 0.98)",
          borderRightColor: "rgba(96, 62, 22, 0.98)",
          background: "rgba(3, 10, 24, 0.9)",
          boxShadow: "0 0 0 2px rgba(0,0,0,0.58), 3px 3px 0 rgba(0,0,0,0.28)",
          color: "#ffd966",
          cursor: "pointer",
          fontFamily: '"Press Start 2P", "DotGothic16", monospace',
          fontSize: "0.5rem",
          lineHeight: 1,
          padding: "8px 10px",
          textShadow: "2px 2px 0 rgba(0,0,0,0.72)",
        }}
      >
        CANCEL
      </button>
    </motion.div>
  );
}

function GuildTownToast({ message }: { message: string }) {
  return (
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
      {message}
    </p>
  );
}
