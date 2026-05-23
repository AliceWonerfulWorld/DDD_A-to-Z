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
  return c.json({ news: items });
});

const port = parseInt(process.env.PORT ?? "8081", 10);

serve({ fetch: app.fetch, port }, () => {
  console.log(`News API listening on port ${port}`);
});
