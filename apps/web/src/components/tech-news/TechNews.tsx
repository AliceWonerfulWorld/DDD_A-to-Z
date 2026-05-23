import { motion } from "framer-motion";
import { useCallback, useEffect, useRef, useState } from "react";
import { PATHS } from "../../constants/paths";
import { fetchTechNews, type TechNewsItem } from "../../features/tech-news/api";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { useAudioSettings } from "../../features/audio/useAudioSettings";
import { BACK_NAVIGATION_SE_SRC, useBackNavigationSe } from "../../hooks/useBackNavigationSe";
import { steppedEase } from "../../lib/animationUtils";
import { findGuildBySlug, GUILD_MASTERS, type GuildMaster } from "../../features/guild/guildMaster";
import styles from "./TechNews.module.css";

interface TechNewsProps {
  onNavigate: (path: string) => void;
}

const NEWS_POLLING_MS = 5 * 60 * 1000;

function formatPublishedAt(iso: string): string {
  const d = new Date(iso);
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHr = Math.floor(diffMin / 60);
  if (diffHr < 24) return `${diffHr}h ago`;
  const diffDay = Math.floor(diffHr / 24);
  if (diffDay < 30) return `${diffDay}d ago`;
  return d.toLocaleDateString();
}

export function TechNews({ onNavigate }: TechNewsProps) {
  const { isSeEnabled } = useAudioSettings();
  const { backNavigationSeRef, navigateBackWithSe } = useBackNavigationSe(onNavigate);
  const selectSeRef = useRef<HTMLAudioElement | null>(null);

  const [selectedSlug, setSelectedSlug] = useState<string>(GUILD_MASTERS[0].slug);
  const [items, setItems] = useState<TechNewsItem[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const latestRequestId = useRef(0);

  const loadNews = useCallback((slug: string) => {
    const requestId = ++latestRequestId.current;
    setIsLoading(true);
    setError(null);
    fetchTechNews(slug)
      .then((data) => {
        if (requestId !== latestRequestId.current) return;
        setItems(data);
      })
      .catch((err) => {
        if (requestId !== latestRequestId.current) return;
        console.error("failed to fetch tech news", err);
        setError("FAILED TO LOAD NEWS");
        setItems([]);
      })
      .finally(() => {
        if (requestId !== latestRequestId.current) return;
        setIsLoading(false);
      });
  }, []);

  useEffect(() => {
    loadNews(selectedSlug);
  }, [selectedSlug, loadNews]);

  useEffect(() => {
    let intervalID: number | undefined;
    if (items.length > 0) {
      intervalID = window.setInterval(() => {
        loadNews(selectedSlug);
      }, NEWS_POLLING_MS);
    }
    return () => {
      if (intervalID !== undefined) {
        window.clearInterval(intervalID);
      }
    };
  }, [selectedSlug, items.length, loadNews]);

  const playSelectSe = useCallback(() => {
    const audio = selectSeRef.current;
    if (!audio || !isSeEnabled) return;
    if (audio.preload === "none" && audio.readyState === HTMLMediaElement.HAVE_NOTHING) {
      audio.load();
    }
    audio.currentTime = 0;
    void audio.play().catch(() => {});
  }, [isSeEnabled]);

  const handleSlugChange = useCallback(
    (slug: string) => {
      if (slug === selectedSlug) return;
      playSelectSe();
      setSelectedSlug(slug);
    },
    [selectedSlug, playSelectSe],
  );

  const currentGuild: GuildMaster | null = findGuildBySlug(selectedSlug);
  const accentColor = currentGuild?.color ?? "#f4ecd0";

  return (
    <main className={styles.page}>
      <audio
        ref={backNavigationSeRef}
        src={BACK_NAVIGATION_SE_SRC}
        preload="none"
        aria-hidden="true"
      />
      <audio
        ref={selectSeRef}
        src={AUDIO_ASSETS.se.buttonClick}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />

      <button
        type="button"
        className={styles.backButton}
        onClick={() => void navigateBackWithSe(PATHS.GUILD)}
      >
        &lt; BACK
      </button>

      <motion.h1
        className={styles.title}
        initial={{ opacity: 0, y: -12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.06, duration: 0.28, ease: steppedEase(5) }}
      >
        TECH NEWS
      </motion.h1>

      <nav className={styles.slugNav}>
        {GUILD_MASTERS.map((guild) => (
          <button
            key={guild.slug}
            type="button"
            className={`${styles.slugButton} ${guild.slug === selectedSlug ? styles.slugButtonActive : ""}`}
            style={
              guild.slug === selectedSlug
                ? {
                    borderColor: accentColor,
                    color: accentColor,
                    boxShadow: `0 0 12px ${accentColor}44, inset 0 0 10px ${accentColor}22`,
                  }
                : undefined
            }
            aria-pressed={guild.slug === selectedSlug}
            onClick={() => handleSlugChange(guild.slug)}
          >
            {guild.icon} {guild.name.toUpperCase()}
          </button>
        ))}
      </nav>

      <motion.section
        className={styles.newsList}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.18, duration: 0.32, ease: steppedEase(5) }}
      >
        {isLoading && items.length === 0 && (
          <div className={styles.statusMessage}>SYNCING NEWS FEED...</div>
        )}
        {error && <div className={styles.errorMessage}>{error}</div>}
        {!isLoading && !error && items.length === 0 && (
          <div className={styles.statusMessage}>NO NEWS AVAILABLE</div>
        )}
        {items.map((item) => (
          <a
            key={item.url}
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className={styles.newsCard}
            style={{ borderLeftColor: accentColor }}
          >
            <div className={styles.newsHeader}>
              <span className={styles.newsTitle}>{item.title}</span>
              <span className={styles.newsSource} style={{ color: accentColor }}>
                {item.source}
              </span>
            </div>
            {item.summary && <p className={styles.newsSummary}>{item.summary}</p>}
            <span className={styles.newsTime}>{formatPublishedAt(item.published_at)}</span>
          </a>
        ))}
      </motion.section>

      <div aria-hidden="true" className={styles.scanlines} />
    </main>
  );
}
