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
  buildingId?: string;
  type?: string;
  building_type?: string;
  level?: number;
  x: number;
  y: number;
  width: number;
  z_index?: number;
}

interface GuildTownResponse {
  guildLevel?: number;
  currentExp?: number;
  current_exp?: number;
  guild_level?: number;
  guild_experience?: number;
  current_guild_level_experience?: number;
  next_guild_level_experience?: number;
  buildings: GuildTownBuildingResponse[];
  inventory: GuildTownInventoryResponse[];
  placements: GuildTownPlacementResponse[];
}

interface HomeCPResponse {
  total_cp: number;
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
  currentGuildLevelExperience: number;
  guildExperience: number;
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
    currentGuildLevelExperience: town.current_guild_level_experience ?? 0,
    guildExperience: town.currentExp ?? town.current_exp ?? town.guild_experience ?? 0,
    guildLevel: town.guildLevel ?? town.guild_level ?? 1,
    nextLevelCp:
      town.next_guild_level_experience ??
      Math.max(town.currentExp ?? town.current_exp ?? town.guild_experience ?? 0, 1),
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
        level: item.level,
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

export async function buyBuilding(buildingId: string): Promise<GuildTownStatus> {
  await apiFetch<GuildTownResponse>("/me/guild/town/buildings", {
    body: JSON.stringify({ buildingId }),
    method: "POST",
  });
  return fetchGuildTownStatus();
}

export async function upgradeBuilding(
  placedBuildingId: string,
  nextLevel: number,
): Promise<GuildTownStatus> {
  await apiFetch<GuildTownResponse>(
    `/me/guild/town/placements/${encodeURIComponent(placedBuildingId)}/upgrade`,
    {
      body: JSON.stringify({ nextLevel }),
      method: "PATCH",
    },
  );
  return fetchGuildTownStatus();
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
  const type = placement.building_type ?? placement.buildingId ?? placement.type ?? "";
  const item = itemByType.get(type);

  return {
    id: placement.id,
    type,
    buildingId: type,
    level: placement.level ?? 1,
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
    const type = placement.building_type ?? placement.buildingId ?? placement.type ?? "";
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
