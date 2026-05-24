import { useEffect, useMemo, useRef, useState } from "react";
import { calcTerritory, loadHexMap, percentToAxial, type TerritoryCell } from "../../wasm/hexmap";
import type { WarGuild } from "./WarMapData";

const GRID_SIZE = 5;
const HEX_DRAW_SIZE = 4.95;
// 塗り関係
const FILL_BASE = 0.18;
const FILL_RANGE = 0.32;
const FILL_CONTESTED = 0.45;

// pointy-top の6方向隣接 (axial)
const HEX_NEIGHBORS: [number, number][] = [
  [1, 0],
  [1, -1],
  [0, -1],
  [-1, 0],
  [-1, 1],
  [0, 1],
];

function hexVertices(cx: number, cy: number): [number, number][] {
  return Array.from({ length: 6 }, (_, i) => {
    const angle = (Math.PI / 180) * (60 * i + 30);
    return [cx + HEX_DRAW_SIZE * Math.cos(angle), cy + HEX_DRAW_SIZE * Math.sin(angle)] as [
      number,
      number,
    ];
  });
}

function hexPoints(cx: number, cy: number): string {
  return hexVertices(cx, cy)
    .map(([x, y]) => `${x.toFixed(2)},${y.toFixed(2)}`)
    .join(" ");
}

function axialToPercent(q: number, r: number): [number, number] {
  return [50 + GRID_SIZE * (Math.sqrt(3) * q + (Math.sqrt(3) / 2) * r), 50 + GRID_SIZE * 1.5 * r];
}

interface TerritoryOverlayProps {
  guilds: WarGuild[];
}

export function TerritoryOverlay({ guilds }: TerritoryOverlayProps) {
  const [isReady, setIsReady] = useState(false);
  const [cells, setCells] = useState<TerritoryCell[]>([]);
  const loadInitiated = useRef(false);

  useEffect(() => {
    if (loadInitiated.current) return;
    loadInitiated.current = true;
    loadHexMap()
      .then(() => setIsReady(true))
      .catch(() => {
        // hexmap.wasm not present (pre-build); overlay simply stays hidden
      });
  }, []);

  const guildInputs = useMemo(
    () =>
      guilds.map((g) => {
        const [q, r] = percentToAxial(g.x, g.y);
        return { q, r, cp: g.totalCp };
      }),
    [guilds],
  );

  useEffect(() => {
    if (!isReady || guildInputs.length === 0) return;
    setCells(calcTerritory(guildInputs));
  }, [isReady, guildInputs]);

  const cellMap = useMemo(() => {
    const m = new Map<string, TerritoryCell>();
    for (const c of cells) m.set(`${c.q},${c.r}`, c);
    return m;
  }, [cells]);

  // 異なるギルドとの境界エッジを収集
  const borderEdges = useMemo(() => {
    const edges: { x1: number; y1: number; x2: number; y2: number; color: string }[] = [];
    for (const cell of cells) {
      if (cell.guildIdx < 0) continue;
      const guild = guilds[cell.guildIdx];
      if (!guild) continue;
      const [cx, cy] = axialToPercent(cell.q, cell.r);
      const verts = hexVertices(cx, cy);
      HEX_NEIGHBORS.forEach(([dq, dr], i) => {
        const neighbor = cellMap.get(`${cell.q + dq},${cell.r + dr}`);
        const neighborGuild = neighbor ? neighbor.guildIdx : -1;
        if (neighborGuild !== cell.guildIdx) {
          const [x1, y1] = verts[i];
          const [x2, y2] = verts[(i + 1) % 6];
          edges.push({ x1, y1, x2, y2, color: guild.color });
        }
      });
    }
    return edges;
  }, [cells, cellMap, guilds]);

  if (cells.length === 0) return null;

  return (
    <svg
      aria-hidden="true"
      style={{
        position: "absolute",
        inset: 0,
        width: "100%",
        height: "100%",
        pointerEvents: "none",
        zIndex: 1,
      }}
      viewBox="0 0 100 100"
      preserveAspectRatio="none"
    >
      {/* 塗り：同ギルド内はストロークなし */}
      {cells.map((cell) => {
        if (cell.guildIdx < 0) return null;
        const guild = guilds[cell.guildIdx];
        if (!guild) return null;
        const [cx, cy] = axialToPercent(cell.q, cell.r);
        const fillOpacity = cell.contested
          ? FILL_CONTESTED
          : FILL_BASE + (cell.intensity / 100) * FILL_RANGE;
        return (
          <polygon
            key={`${cell.q},${cell.r}`}
            points={hexPoints(cx, cy)}
            fill={guild.color}
            fillOpacity={fillOpacity}
            stroke="none"
          />
        );
      })}
      {/* 境界線：異なるギルド間の辺のみ */}
      {borderEdges.map((e, i) => (
        <line
          key={i}
          x1={e.x1.toFixed(2)}
          y1={e.y1.toFixed(2)}
          x2={e.x2.toFixed(2)}
          y2={e.y2.toFixed(2)}
          stroke={e.color}
          strokeOpacity={0.85}
          strokeWidth={0.25}
          strokeLinecap="round"
        />
      ))}
    </svg>
  );
}
