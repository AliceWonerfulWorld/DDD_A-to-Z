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
      const dbAge = now - new Date(dbItems[0].publishedAt).getTime();
      if (dbAge < CACHE_TTL_MS) {
        cache.set(slug, { items: dbItems, fetchedAt: now });
        return dbItems;
      }
    }
  } catch {
    // DB unavailable, fall through to RSS
  }

  const items = await fetchNews(slug);
  if (items.length === 0 && cached) {
    return cached.items;
  }

  cache.set(slug, { items, fetchedAt: now });

  try {
    await saveToDB(slug, items);
  } catch {
    // Cache save failed, non-fatal
  }

  return items;
}

async function loadFromDB(slug: string): Promise<TechNewsItem[]> {
  const rows = await sql`
    SELECT title, url, source, summary, published_at, slug
    FROM tech_news_cache
    WHERE slug = ${slug}
    ORDER BY published_at DESC
    LIMIT 20
  `;
  return rows.map((r) => ({
    title: r.title as string,
    url: r.url as string,
    source: r.source as string,
    summary: r.summary as string,
    publishedAt: r.published_at as string,
    slug: r.slug as string,
  }));
}

async function saveToDB(slug: string, items: TechNewsItem[]): Promise<void> {
  if (items.length === 0) return;

  await sql`DELETE FROM tech_news_cache WHERE slug = ${slug}`;

  const now = new Date().toISOString();
  for (const item of items) {
    await sql`
      INSERT INTO tech_news_cache (slug, title, url, source, summary, published_at, fetched_at)
      VALUES (${item.slug}, ${item.title}, ${item.url}, ${item.source}, ${item.summary}, ${item.publishedAt}, ${now})
    `;
  }
}
