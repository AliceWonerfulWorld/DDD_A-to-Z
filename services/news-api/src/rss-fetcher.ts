import Parser from "rss-parser";
import {
  SLUG_FEEDS,
  type GuildSlug,
  type TechNewsItem,
} from "./constants";

const parser = new Parser({
  timeout: 15000,
  headers: {
    "User-Agent": "LangWar/1.0",
  },
});

export async function fetchNews(slug: GuildSlug): Promise<TechNewsItem[]> {
  const feedUrls = SLUG_FEEDS[slug];
  if (!feedUrls) {
    return [];
  }

  const seen = new Set<string>();
  const items: TechNewsItem[] = [];

  for (const url of feedUrls) {
    if (items.length >= 20) break;

    try {
      const feed = await parser.parseURL(url);
      console.log(`[rss-fetcher] ${slug}: fetched ${feed.items?.length ?? 0} items from ${url}`);

      for (const entry of feed.items ?? []) {
        if (items.length >= 20) break;
        if (!entry.link) continue;
        if (!entry.title) continue;
        let normalizedUrl: string;
        try {
          const parsed = new URL(entry.link);
          if (parsed.protocol !== "http:" && parsed.protocol !== "https:") continue;
          normalizedUrl = parsed.href;
        } catch {
          continue;
        }
        if (seen.has(normalizedUrl)) continue;
        seen.add(normalizedUrl);

        const summary = entry.contentSnippet
          ? entry.contentSnippet.length > 300
            ? entry.contentSnippet.slice(0, 300) + "..."
            : entry.contentSnippet
          : "";

        const publishedAt =
          entry.pubDate ?? entry.isoDate ?? new Date().toISOString();

        items.push({
          title: entry.title,
          url: normalizedUrl,
          source: feed.title ?? url,
          summary,
          publishedAt,
          slug,
        });
      }
    } catch (err) {
      console.error(`[rss-fetcher] ${slug}: failed to fetch ${url}`, err instanceof Error ? err.message : err);
    }
  }

  if (items.length === 0) {
    console.warn(`[rss-fetcher] ${slug}: no items collected from any feed`);
  }

  return items;
}
