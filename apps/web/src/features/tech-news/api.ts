export interface TechNewsItem {
  title: string;
  url: string;
  source: string;
  summary: string;
  published_at: string;
}

export async function fetchTechNews(slug: string): Promise<TechNewsItem[]> {
  const res = await fetch(`/news-api/guilds/${slug}/tech-news`, {
    credentials: "include",
  });

  if (!res.ok) {
    if (res.status === 400) {
      return [];
    }
    return [];
  }

  const data = await res.json();
  return data.news ?? [];
}
