export const GUILD_SLUGS = [
  "rust",
  "python",
  "go",
  "typescript",
  "java",
  "haskell",
  "zig",
] as const;

export type GuildSlug = (typeof GUILD_SLUGS)[number];

export const SLUG_FEEDS: Record<GuildSlug, string[]> = {
  rust: [
    "https://blog.rust-lang.org/feed.xml",
    "https://rustacean-station.org/episodes.xml",
  ],
  python: [
    "https://blog.python.org/feeds/posts/default",
    "https://realpython.com/atom.xml",
  ],
  go: ["https://go.dev/blog/feed.atom"],
  typescript: [
    "https://dev.to/feed/tag/typescript",
  ],
  java: [
    "https://inside.java/feed.xml",
    "https://spring.io/blog/feed",
  ],
  haskell: [
    "https://haskellweekly.news/haskell-weekly.atom",
  ],
  zig: [
    "https://ziglang.org/atom.xml",
    "https://dev.to/feed/tag/zig",
  ],
};

export const CACHE_TTL_MS = 30 * 60 * 1000;

export interface TechNewsItem {
  title: string;
  url: string;
  source: string;
  summary: string;
  publishedAt: string;
  slug: string;
}
