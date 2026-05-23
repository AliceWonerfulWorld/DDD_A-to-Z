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

const rawPort = parseInt(process.env.PORT ?? "8082", 10);
const port = !Number.isNaN(rawPort) && Number.isInteger(rawPort) && rawPort >= 1 && rawPort <= 65535
  ? rawPort
  : 8082;

serve({ fetch: app.fetch, port }, () => {
  console.log(`News API listening on port ${port}`);
});
