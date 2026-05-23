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
    throw new Error(`fetchTechNews: ${res.status} ${res.statusText}`);
  }

  const data = await res.json();
  if (!Array.isArray(data.news)) {
    throw new Error("fetchTechNews: invalid response, news is not an array");
  }
  return data.news;
}
