import type { InventoryItem } from "./types";

export const MIN_SCALE = 0.5;
export const MAX_SCALE = 2.5;
export const ZOOM_STEP = 0.15;
export const STORE_ANIMATION_MS = 260;

export const INITIAL_INVENTORY: InventoryItem[] = [
  {
    type: "tent",
    name: "TENT",
    title: "旅人のテント",
    description: "ギルドの仲間が遠征前に集う簡易拠点。休息と作戦会議に使われる。",
    count: 2,
    src: "/town/tent.png",
    minMapWidth: 210,
    mapWidthVw: 29,
    maxMapWidth: 430,
  },
  {
    type: "bonfire",
    name: "BONFIRE",
    title: "団らんの焚き火",
    description: "夜のギルドタウンを照らす小さな火。仲間の士気をじんわり温める。",
    count: 3,
    src: "/town/bonfire.png",
    minMapWidth: 92,
    mapWidthVw: 12,
    maxMapWidth: 164,
  },
];
