import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import { PATHS } from "../../constants/paths";
import { fetchTechNews, type TechNewsItem } from "../../features/tech-news/api";
import { steppedEase } from "../../lib/animationUtils";

const TECH_NEWS_POLLING_MS = 5 * 60 * 1000;
const MAX_VISIBLE_ITEMS = 5;

interface GuildTechNewsPanelProps {
  guildSlug: string | null;
  onNavigate: (path: string) => void;
}

export function GuildTechNewsPanel({ guildSlug, onNavigate }: GuildTechNewsPanelProps) {
  const [items, setItems] = useState<TechNewsItem[]>([]);

  useEffect(() => {
    if (!guildSlug) return;

    let cancelled = false;
    const intervalID = window.setInterval(() => {
      fetchTechNews(guildSlug).then((data) => {
        if (!cancelled) setItems(data);
      });
    }, TECH_NEWS_POLLING_MS);

    fetchTechNews(guildSlug).then((data) => {
      if (!cancelled) setItems(data);
    });

    return () => {
      cancelled = true;
      window.clearInterval(intervalID);
    };
  }, [guildSlug]);

  if (items.length === 0) return null;

  const visibleItems = items.slice(0, MAX_VISIBLE_ITEMS);

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2, duration: 0.34, ease: steppedEase(5) }}
      aria-label="Guild tech news"
      style={{
        position: "fixed",
        bottom: "calc(env(safe-area-inset-bottom, 0px) + clamp(106px, 10vh, 128px))",
        left: "clamp(14px, 2.2vw, 28px)",
        zIndex: 3,
        width: "min(320px, calc(50vw - 48px))",
        maxHeight: "min(360px, 45vh)",
        border: "2px solid rgba(0, 245, 255, 0.35)",
        borderBottomColor: "rgba(2, 54, 72, 0.88)",
        borderRightColor: "rgba(2, 54, 72, 0.88)",
        background: "rgba(3, 10, 24, 0.82)",
        boxShadow:
          "0 0 0 2px rgba(0,0,0,0.64), 0 0 14px rgba(0,245,255,0.1), inset 0 0 10px rgba(0,245,255,0.05)",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          padding: "clamp(6px, 0.8vh, 10px) clamp(10px, 1.2vw, 14px)",
          borderBottom: "1px solid rgba(0, 245, 255, 0.2)",
          fontSize: "clamp(0.44rem, 0.85vw, 0.56rem)",
          color: "rgba(0, 245, 255, 0.7)",
          textShadow: "1px 1px 0 rgba(0,0,0,0.72)",
          letterSpacing: "0.08em",
        }}
      >
        TECH NEWS
      </div>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: 0,
          overflow: "hidden",
          flex: 1,
        }}
      >
        {visibleItems.map((item) => (
          <a
            key={item.url}
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            style={{
              display: "block",
              padding: "clamp(5px, 0.7vh, 9px) clamp(10px, 1.2vw, 14px)",
              textDecoration: "none",
              color: "#c8f4ff",
              fontSize: "clamp(0.4rem, 0.78vw, 0.52rem)",
              lineHeight: 1.45,
              textShadow: "1px 1px 0 rgba(0,0,0,0.72)",
              borderBottom: "1px solid rgba(0, 245, 255, 0.08)",
              transition: "background 0.12s ease",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = "rgba(0, 245, 255, 0.06)";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = "transparent";
            }}
          >
            {item.title}
          </a>
        ))}
      </div>

      <button
        type="button"
        onClick={() => {
          (onNavigate as (path: string, opts?: { state?: Record<string, string> }) => void)(
            PATHS.GUILD_TECH_NEWS,
            { state: { guildSlug: guildSlug ?? "" } },
          );
        }}
        style={{
          border: "none",
          borderTop: "1px solid rgba(0, 245, 255, 0.2)",
          background: "rgba(0, 245, 255, 0.06)",
          color: "rgba(0, 245, 255, 0.8)",
          cursor: "pointer",
          fontFamily: "inherit",
          fontSize: "clamp(0.4rem, 0.78vw, 0.5rem)",
          lineHeight: 1.5,
          padding: "clamp(5px, 0.7vh, 8px) clamp(10px, 1.2vw, 14px)",
          textShadow: "1px 1px 0 rgba(0,0,0,0.72)",
          letterSpacing: "0.06em",
          transition: "background 0.12s ease",
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.background = "rgba(0, 245, 255, 0.14)";
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.background = "rgba(0, 245, 255, 0.06)";
        }}
      >
        [ MORE NEWS → ]
      </button>
    </motion.div>
  );
}
