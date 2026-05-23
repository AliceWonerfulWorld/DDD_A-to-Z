import { CACHE_TTL_MS, type GuildSlug, type TechNewsItem } from "./constants";
import { fetchNews } from "./rss-fetcher";
import sql from "./db";

const cache = new Map<string, { items: TechNewsItem[]; fetchedAt: number }>();

export async function getTechNews(slug: GuildSlug): Promise<TechNewsItem[]> {
  const now = Date.now();

  const cached = cache.get(slug);
  if (cached && now - cached.fetchedAt < CACHE_TTL_MS) {
    return cached.items;
  }

  try {
    const dbItems = await loadFromDB(slug);
    if (dbItems.length > 0) {
      const dbAge = now - new Date(dbItems[0].fetchedAt ?? dbItems[0].publishedAt).getTime();
      if (dbAge < CACHE_TTL_MS) {
        cache.set(slug, { items: dbItems, fetchedAt: now });
        return dbItems;
      }
    }
  } catch (err) {
    console.warn(`[tech-news] failed to load from DB for slug=${slug}`, err instanceof Error ? err.message : err);
  }

  const items = await fetchNews(slug);
  if (items.length === 0 && cached) {
    return cached.items;
  }

  cache.set(slug, { items, fetchedAt: now });

  try {
    await saveToDB(slug, items);
  } catch (err) {
    console.warn(`[tech-news] failed to save to DB for slug=${slug} count=${items.length}`, err instanceof Error ? err.message : err);
  }

  return items;
}

interface DbRow {
  title: string;
  url: string;
  source: string;
  summary: string;
  published_at: string;
  fetched_at: string;
  slug: string;
}

async function loadFromDB(slug: string): Promise<TechNewsItem[]> {
  const rows = await sql`
    SELECT title, url, source, summary, published_at, fetched_at, slug
    FROM tech_news_cache
    WHERE slug = ${slug}
    ORDER BY published_at DESC
    LIMIT 20
  `;
  return (rows as unknown as DbRow[]).map((r) => ({
    title: r.title,
    url: r.url,
    source: r.source,
    summary: r.summary,
    publishedAt: r.published_at,
    fetchedAt: r.fetched_at,
    slug: r.slug as GuildSlug,
  }));
}

async function saveToDB(slug: string, items: TechNewsItem[]): Promise<void> {
  if (items.length === 0) return;

  await sql.begin(async (tx) => {
    await tx`DELETE FROM tech_news_cache WHERE slug = ${slug}`;
    const now = new Date().toISOString();
    for (const item of items) {
      await tx`
        INSERT INTO tech_news_cache (slug, title, url, source, summary, published_at, fetched_at)
        VALUES (${slug}, ${item.title}, ${item.url}, ${item.source}, ${item.summary}, ${item.publishedAt}, ${now})
      `;
    }
  });
}