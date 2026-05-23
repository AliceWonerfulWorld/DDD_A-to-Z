CREATE TABLE tech_news_cache (
  id SERIAL PRIMARY KEY,
  slug TEXT NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  source TEXT NOT NULL DEFAULT '',
  summary TEXT NOT NULL DEFAULT '',
  published_at TIMESTAMPTZ NOT NULL,
  fetched_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX tech_news_cache_slug_fetched_at_idx ON tech_news_cache(slug, fetched_at DESC);