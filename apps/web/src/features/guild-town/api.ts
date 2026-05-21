import { apiFetch } from "../../lib/api/client";
import type {
  GuildSpLanguage,
  InventoryItem,
  PlacedItem,
  UserInventoryState,
} from "../../components/guild-town/types";

interface GuildTownBuildingResponse {
  type: string;
  name: string;
  title: string;
  description: string;
  src: string;
  min_map_width: number;
  map_width_vw: number;
  max_map_width: number;
}

interface GuildTownInventoryResponse {
  type: string;
  count: number;
}

interface GuildTownPlacementResponse {
  id: string;
  type?: string;
  building_type?: string;
  x: number;
  y: number;
  width: number;
  z_index?: number;
}

interface GuildTownResponse {
  buildings: GuildTownBuildingResponse[];
  inventory: GuildTownInventoryResponse[];
  placements: GuildTownPlacementResponse[];
}

interface HomeCPResponse {
  total_cp: number;
  next_player_level_total_cp?: number;
  player_level?: number;
}

interface SkillPointResponse {
  skill_points: Array<{
    language: string;
    balance: number;
  }>;
}

export interface GuildTownStatus {
  availableItems: InventoryItem[];
  currentCp: number;
  guildLevel: number;
  nextLevelCp: number;
  placedItems: PlacedItem[];
  userInventory: UserInventoryState[];
  userSpMap: Record<string, number>;
}

interface SavePlacementsPayload {
  placements: PlacedItem[];
}

export async function fetchGuildTownStatus(): Promise<GuildTownStatus> {
  const [town, home, skillPoints] = await Promise.all([
    apiFetch<GuildTownResponse>("/me/guild/town"),
    apiFetch<HomeCPResponse>("/home"),
    apiFetch<SkillPointResponse>("/me/sp").catch(() => ({ skill_points: [] })),
  ]);

  const availableItems = town.buildings.map(toInventoryItem);
  const itemByType = new Map(availableItems.map((item) => [item.type, item]));
  const userSpMap = skillPoints.skill_points.reduce<Record<string, number>>((spMap, item) => {
    spMap[item.language] = item.balance;
    return spMap;
  }, {});
  const placementOrderById = new Map(
    town.placements.map((placement) => [placement.id, placement.z_index ?? 0]),
  );

  return {
    availableItems,
    currentCp: home.total_cp,
    guildLevel: home.player_level ?? 1,
    nextLevelCp: home.next_player_level_total_cp ?? Math.max(home.total_cp, 1),
    placedItems: town.placements
      .map((placement) => toPlacedItem(placement, itemByType))
      .sort((a, b) => (placementOrderById.get(a.id) ?? 0) - (placementOrderById.get(b.id) ?? 0)),
    userInventory: toRemainingInventory(town.inventory, town.placements),
    userSpMap,
  };
}

export async function saveGuildTownPlacements(
  payload: SavePlacementsPayload,
): Promise<GuildTownStatus> {
  await apiFetch<GuildTownResponse>("/me/guild/town/placements", {
    body: JSON.stringify({
      placements: payload.placements.map((item) => ({
        id: item.id.startsWith("local-") ? "" : item.id,
        building_type: item.buildingId ?? item.type,
        x: item.x,
        y: item.y,
        width: item.width,
      })),
    }),
    method: "PUT",
  });
  return fetchGuildTownStatus();
}

export async function deployBuilding(payload: SavePlacementsPayload): Promise<GuildTownStatus> {
  return saveGuildTownPlacements(payload);
}

export async function buyBuilding(buildingId: string): Promise<never> {
  throw new Error(`buyBuilding API is not implemented yet: ${buildingId}`);
}

export async function upgradeBuilding(placedBuildingId: string): Promise<never> {
  throw new Error(`upgradeBuilding API is not implemented yet: ${placedBuildingId}`);
}

function toInventoryItem(item: GuildTownBuildingResponse): InventoryItem {
  return {
    type: item.type as InventoryItem["type"],
    name: item.name,
    title: item.title,
    description: item.description,
    count: 0,
    src: item.src,
    minMapWidth: item.min_map_width,
    mapWidthVw: item.map_width_vw,
    maxMapWidth: item.max_map_width,
  };
}

function toPlacedItem(
  placement: GuildTownPlacementResponse,
  itemByType: Map<string, InventoryItem>,
): PlacedItem {
  const type = placement.building_type ?? placement.type ?? "";
  const item = itemByType.get(type);

  return {
    id: placement.id,
    type,
    buildingId: type,
    level: 1,
    name: item?.name ?? type,
    title: item?.title ?? type,
    description: item?.description ?? "",
    src: item?.src ?? "/town/tent.png",
    x: placement.x,
    y: placement.y,
    width: placement.width,
  };
}

function toRemainingInventory(
  inventory: GuildTownInventoryResponse[],
  placements: GuildTownPlacementResponse[],
): UserInventoryState[] {
  const placedCountByType = placements.reduce<Record<string, number>>((countMap, placement) => {
    const type = placement.building_type ?? placement.type ?? "";
    countMap[type] = (countMap[type] ?? 0) + 1;
    return countMap;
  }, {});

  return inventory.map((item) => ({
    buildingId: item.type,
    count: Math.max(0, item.count - (placedCountByType[item.type] ?? 0)),
  }));
}

export function guildSpBalance(
  userSpMap: Record<string, number>,
  language: GuildSpLanguage,
): number {
  return userSpMap[language] ?? 0;
}
