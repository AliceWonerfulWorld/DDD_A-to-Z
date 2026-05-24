// TypeScript wrapper for the MoonBit-compiled hexmap.wasm module.
// Provides territory calculation: given guild positions + CP, returns which
// hex cell each guild controls and how strongly.

export interface TerritoryCell {
  q: number;
  r: number;
  guildIdx: number; // 0-based guild index, or -1 for unclaimed
  intensity: number; // 0-100
  contested: boolean; // true when two guilds are nearly equal influence
}

type HexMapExports = {
  set_guild_count: (n: number) => void;
  set_guild: (idx: number, q: number, r: number, cp: number) => void;
  calc: (radius: number) => void;
  get_cell_count: () => number;
  get_cell_q: (i: number) => number;
  get_cell_r: (i: number) => number;
  get_cell_guild: (i: number) => number;
  get_cell_intensity: (i: number) => number;
  get_cell_contested: (i: number) => number;
};

let _exports: HexMapExports | null = null;

export async function loadHexMap(): Promise<void> {
  const result = await WebAssembly.instantiateStreaming(fetch("/hexmap.wasm"));
  _exports = result.instance.exports as unknown as HexMapExports;
}

// Convert percentage map coordinates to axial hex coords.
// Assumes pointy-top hexes, grid size=5, centered at (50%, 50%).
export function percentToAxial(x: number, y: number, size = 5): [number, number] {
  const cx = x - 50;
  const cy = y - 50;
  const q = Math.round(((Math.sqrt(3) / 3) * cx - (1 / 3) * cy) / size);
  const r = Math.round(((2 / 3) * cy) / size);
  return [q, r];
}

// Calculate territory for a set of guilds.
// guilds: array of { q, r, cp } in axial coordinates
// radius: hex grid radius to evaluate (default 9 → 271 cells)
export function calcTerritory(
  guilds: Array<{ q: number; r: number; cp: number }>,
  radius = 9,
): TerritoryCell[] {
  if (!_exports) return [];

  const ex = _exports;
  ex.set_guild_count(guilds.length);
  for (let i = 0; i < guilds.length; i++) {
    ex.set_guild(i, guilds[i].q, guilds[i].r, Math.max(1, guilds[i].cp));
  }
  ex.calc(radius);

  const count = ex.get_cell_count();
  const cells: TerritoryCell[] = [];
  for (let i = 0; i < count; i++) {
    cells.push({
      q: ex.get_cell_q(i),
      r: ex.get_cell_r(i),
      guildIdx: ex.get_cell_guild(i),
      intensity: ex.get_cell_intensity(i),
      contested: ex.get_cell_contested(i) === 1,
    });
  }
  return cells;
}
