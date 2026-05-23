import { Hono } from "hono";
import { cors } from "hono/cors";
import { serve } from "@hono/node-server";
import { GUILD_SLUGS, type GuildSlug } from "./constants";
import { getTechNews } from "./tech-news-service";

const app = new Hono();

app.use("*", cors());

app.get("/healthz", (c) => c.json({ status: "ok" }));

app.get("/guilds/:slug/tech-news", async (c) => {
  const slug = c.req.param("slug");

  if (!GUILD_SLUGS.includes(slug as GuildSlug)) {
    return c.json(
      { error: "invalid_slug", message: "Invalid guild slug" },
      400,
    );
  }

  const items = await getTechNews(slug as GuildSlug);
  return c.json({
    news: items.map((item) => ({
      title: item.title,
      url: item.url,
      source: item.source,
      summary: item.summary,
      published_at: item.publishedAt,
      slug: item.slug,
    })),
  });
});

const portEnv = process.env.PORT ?? "";
const port = /^\d+$/.test(portEnv)
  ? (() => {
      const n = Number(portEnv);
      return n >= 1 && n <= 65535 ? n : 8082;
    })()
  : 8082;

serve({ fetch: app.fetch, port }, () => {
  console.log(`News API listening on port ${port}`);
});
