// アセット管理用のMap定義
// 後からの画像差し替えを容易にする設計
export const SPRITE_ASSETS = {
  // 暫定対応: 全言語共通でRustの侍をサンプルとして表示
  RUST_SAMURAI: "/character/rust_samurai_sheet.png",
  GOPHER: "/character/gopher.webp",
  PYTHON: "/character/python.webp",
  RUST: "/character/rust.webp",
  JAVA: "/character/duke.webp",
} as const;
