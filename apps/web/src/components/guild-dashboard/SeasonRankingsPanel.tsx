import { motion, type Variants } from "framer-motion";
import { useEffect, useState, useCallback } from "react";
import {
  type GuildSeasonMemberRanking,
  type GuildSeasonRanking,
  type Season,
  fetchSeasons,
  fetchGuildSeasonRankings,
  fetchGuildSeasonMemberRankings,
} from "../../features/season/api";
import { findGuildByID } from "../../features/guild/guildMaster";
import { steppedEase } from "../../lib/animationUtils";
import styles from "./DashboardPanels.module.css";

const tabContentVariants: Variants = {
  hidden: { opacity: 0, y: 12 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.28, ease: steppedEase(5) },
  },
  exit: {
    opacity: 0,
    y: -10,
    transition: { duration: 0.18, ease: steppedEase(4) },
  },
};

interface SeasonRankingsPanelProps {
  guildID: string | null;
  isMobile?: boolean;
}

export function SeasonRankingsPanel({ guildID: _guildID, isMobile }: SeasonRankingsPanelProps) {
  const [allSeasons, setAllSeasons] = useState<Season[]>([]);
  const [selectedSeason, setSelectedSeason] = useState<Season | null>(null);
  const [guildRankings, setGuildRankings] = useState<GuildSeasonRanking[]>([]);
  const [expandedGuildID, setExpandedGuildID] = useState<string | null>(null);
  const [membersByGuild, setMembersByGuild] = useState<Record<string, GuildSeasonMemberRanking[]>>(
    {},
  );
  const [loadingMembers, setLoadingMembers] = useState<Record<string, boolean>>({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    const load = async () => {
      try {
        const seasons = await fetchSeasons();
        if (!isMounted) return;
        setAllSeasons(seasons);
        if (seasons.length > 0) {
          setSelectedSeason(seasons[0]);
        }
      } catch (error) {
        if (isMounted) console.error("failed to fetch seasons", error);
      } finally {
        if (isMounted) setIsLoading(false);
      }
    };

    void load();
    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    if (!selectedSeason) return;
    let isMounted = true;

    const loadRankings = async () => {
      try {
        const gr = await fetchGuildSeasonRankings(selectedSeason.number);
        if (!isMounted) return;
        setGuildRankings(gr);
        setExpandedGuildID(null);
        setMembersByGuild({});
      } catch (error) {
        if (isMounted) console.error("failed to fetch guild rankings", error);
      }
    };

    void loadRankings();
    return () => {
      isMounted = false;
    };
  }, [selectedSeason]);

  const toggleGuild = useCallback(
    async (guildID: string) => {
      if (expandedGuildID === guildID) {
        setExpandedGuildID(null);
        return;
      }

      if (!selectedSeason) return;
      setExpandedGuildID(guildID);

      if (!membersByGuild[guildID]) {
        setLoadingMembers((prev) => ({ ...prev, [guildID]: true }));
        try {
          const members = await fetchGuildSeasonMemberRankings(selectedSeason.number, guildID);
          setMembersByGuild((prev) => ({ ...prev, [guildID]: members.slice(0, 5) }));
        } catch (error) {
          console.error("failed to fetch member rankings", error);
        } finally {
          setLoadingMembers((prev) => ({ ...prev, [guildID]: false }));
        }
      }
    },
    [selectedSeason, expandedGuildID, membersByGuild],
  );

  const pill = {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    gap: isMobile ? "2px" : "4px",
    padding: isMobile ? "4px 8px" : "6px clamp(12px, 1.6vw, 20px)",
    border: "1px solid rgba(0, 245, 255, 0.25)",
    borderRadius: "4px",
    cursor: "pointer",
    fontFamily: '"DotGothic16", monospace',
    fontSize: isMobile ? "0.44rem" : "clamp(0.52rem, 0.95vw, 0.68rem)",
    transition: "all 0.18s ease",
    userSelect: "none" as const,
    background: "rgba(1, 8, 22, 0.58)",
    color: "rgba(214, 255, 228, 0.72)",
  };

  const pillActive = {
    ...pill,
    background: "rgba(0, 245, 255, 0.18)",
    borderColor: "rgba(0, 245, 255, 0.6)",
    color: "#00f5ff",
    boxShadow: "0 0 8px rgba(0, 245, 255, 0.25)",
  };

  if (isLoading) {
    return (
      <motion.section
        key="season-loading"
        variants={tabContentVariants}
        initial="hidden"
        animate="visible"
        exit="exit"
        style={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "rgba(214, 255, 228, 0.72)",
          fontFamily: '"DotGothic16", monospace',
          fontSize: isMobile ? "0.6rem" : "clamp(0.66rem, 1.22vw, 0.94rem)",
        }}
      >
        SYNCING SEASON DATA...
      </motion.section>
    );
  }

  if (allSeasons.length === 0) {
    return (
      <motion.section
        key="season-error"
        variants={tabContentVariants}
        initial="hidden"
        animate="visible"
        exit="exit"
        style={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "rgba(214, 255, 228, 0.72)",
          fontFamily: '"DotGothic16", monospace',
          fontSize: isMobile ? "0.6rem" : "clamp(0.66rem, 1.22vw, 0.94rem)",
        }}
      >
        NO SEASONS AVAILABLE
      </motion.section>
    );
  }

  return (
    <motion.section
      key="season"
      variants={tabContentVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      className={styles.hideScrollbar}
      style={{
        height: "100%",
        minHeight: 0,
        display: "flex",
        flexDirection: "column",
        gap: isMobile ? "4px" : "8px",
        overflowY: "auto",
      }}
    >
      <div
        style={{
          display: "flex",
          gap: isMobile ? "4px" : "6px",
          flexWrap: "wrap",
          marginBottom: isMobile ? "4px" : "8px",
        }}
      >
        {allSeasons.map((s) => {
          const active = selectedSeason?.id === s.id;
          return (
            <div key={s.id} style={active ? pillActive : pill} onClick={() => setSelectedSeason(s)}>
              <span>SEASON {s.number}</span>
              {s.is_current && (
                <span style={{ fontSize: "0.7em", color: "#ffd966", marginLeft: "2px" }}>●</span>
              )}
            </div>
          );
        })}
      </div>

      <div
        style={{
          fontSize: isMobile ? "0.4rem" : "clamp(0.48rem, 0.85vw, 0.58rem)",
          color: "rgba(0, 245, 255, 0.7)",
          textShadow: "1px 1px 0 rgba(0,0,0,0.72)",
          letterSpacing: "0.06em",
          marginBottom: isMobile ? "2px" : "4px",
        }}
      >
        {selectedSeason && (
          <span style={{ color: "rgba(244, 236, 208, 0.5)", fontSize: "0.85em" }}>
            {new Date(selectedSeason.starts_at).toLocaleDateString()} –{" "}
            {new Date(selectedSeason.ends_at).toLocaleDateString()}
          </span>
        )}
      </div>

      <div
        style={{
          fontSize: isMobile ? "0.4rem" : "clamp(0.48rem, 0.85vw, 0.58rem)",
          color: "rgba(0, 245, 255, 0.7)",
          textShadow: "1px 1px 0 rgba(0,0,0,0.72)",
          letterSpacing: "0.06em",
          marginBottom: isMobile ? "2px" : "4px",
        }}
      >
        GUILD RANKINGS
      </div>

      {guildRankings.map((ranking, index) => {
        const guild = findGuildByID(ranking.guild_id);
        const accentColor =
          index === 0 ? "#ffd966" : index === 1 ? "#e2e8f0" : index === 2 ? "#b45309" : "#a1a1aa";
        const isExpanded = expandedGuildID === ranking.guild_id;
        const members = membersByGuild[ranking.guild_id];
        const isLoadingMembers = loadingMembers[ranking.guild_id];

        return (
          <div key={ranking.id}>
            <motion.div
              initial={{ opacity: 0, y: -8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05, duration: 0.22, ease: steppedEase(4) }}
              onClick={() => toggleGuild(ranking.guild_id)}
              style={{
                display: "grid",
                gridTemplateColumns: "minmax(28px, auto) 1fr auto",
                alignItems: "center",
                gap: isMobile ? "4px" : "clamp(8px, 1.2vw, 14px)",
                minHeight: isMobile ? "24px" : "clamp(32px, 5.5vh, 44px)",
                border: isMobile
                  ? "1px solid rgba(0, 245, 255, 0.18)"
                  : "2px solid rgba(0, 245, 255, 0.18)",
                background:
                  index === 0
                    ? "linear-gradient(90deg, rgba(255, 217, 102, 0.22), rgba(1, 8, 22, 0.58))"
                    : "rgba(1, 8, 22, 0.48)",
                padding: isMobile ? "3px 6px" : "6px clamp(8px, 1.2vw, 14px)",
                cursor: "pointer",
                transition: "background 0.15s ease",
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = "rgba(0, 245, 255, 0.08)";
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background =
                  index === 0
                    ? "linear-gradient(90deg, rgba(255, 217, 102, 0.22), rgba(1, 8, 22, 0.58))"
                    : "rgba(1, 8, 22, 0.48)";
              }}
            >
              <span
                style={{
                  color: accentColor,
                  fontSize: isMobile ? "0.44rem" : "clamp(0.62rem, 1.2vw, 0.86rem)",
                  lineHeight: 1,
                }}
              >
                #{ranking.rank}
              </span>
              <span
                style={{
                  color: guild?.color ?? "#fff8d7",
                  fontSize: isMobile ? "0.48rem" : "clamp(0.54rem, 1vw, 0.72rem)",
                  overflowWrap: "anywhere",
                }}
              >
                {guild?.icon ?? ""} {guild?.name ?? ranking.guild_id}
              </span>
              <span
                style={{
                  color: "#74f7a1",
                  fontSize: isMobile ? "0.44rem" : "clamp(0.52rem, 0.95vw, 0.66rem)",
                  whiteSpace: "nowrap",
                }}
              >
                {ranking.total_cp.toLocaleString()} CP
              </span>
            </motion.div>

            {isExpanded && (
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: "auto", opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2, ease: steppedEase(4) }}
                style={{ overflow: "hidden" }}
              >
                <div
                  style={{
                    padding: isMobile ? "4px 0 4px 20px" : "6px 0 8px clamp(20px, 3vw, 36px)",
                    display: "flex",
                    flexDirection: "column",
                    gap: isMobile ? "2px" : "4px",
                  }}
                >
                  {isLoadingMembers && (
                    <span
                      style={{
                        color: "rgba(214, 255, 228, 0.5)",
                        fontSize: isMobile ? "0.36rem" : "clamp(0.42rem, 0.72vw, 0.48rem)",
                        fontFamily: '"DotGothic16", monospace',
                      }}
                    >
                      LOADING...
                    </span>
                  )}

                  {members && members.length === 0 && (
                    <span
                      style={{
                        color: "rgba(214, 255, 228, 0.4)",
                        fontSize: isMobile ? "0.36rem" : "clamp(0.42rem, 0.72vw, 0.48rem)",
                        fontFamily: '"DotGothic16", monospace',
                      }}
                    >
                      NO MEMBER DATA
                    </span>
                  )}

                  {members?.map((member, mIndex) => {
                    const rankColor =
                      mIndex === 0
                        ? "#ffd966"
                        : mIndex === 1
                          ? "#e2e8f0"
                          : mIndex === 2
                            ? "#b45309"
                            : "#a1a1aa";

                    return (
                      <motion.div
                        key={member.id}
                        initial={{ opacity: 0, x: -8 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: mIndex * 0.04, duration: 0.18, ease: steppedEase(4) }}
                        style={{
                          display: "grid",
                          gridTemplateColumns: "minmax(24px, auto) 1fr auto",
                          alignItems: "center",
                          gap: isMobile ? "3px" : "clamp(4px, 0.6vw, 8px)",
                          minHeight: isMobile ? "20px" : "clamp(24px, 3.5vh, 32px)",
                          border: "1px solid rgba(244, 236, 208, 0.1)",
                          background: "rgba(1, 8, 22, 0.3)",
                          padding: isMobile ? "2px 6px" : "3px clamp(6px, 0.8vw, 10px)",
                        }}
                      >
                        <span
                          style={{
                            color: rankColor,
                            fontSize: isMobile ? "0.38rem" : "clamp(0.48rem, 0.85vw, 0.6rem)",
                            lineHeight: 1,
                          }}
                        >
                          #{member.rank}
                        </span>
                        <span
                          style={{
                            color: "#fff8d7",
                            fontSize: isMobile ? "0.4rem" : "clamp(0.46rem, 0.85vw, 0.58rem)",
                            overflowWrap: "anywhere",
                          }}
                        >
                          {member.user_name}
                        </span>
                        <span
                          style={{
                            color: "#74f7a1",
                            fontSize: isMobile ? "0.36rem" : "clamp(0.42rem, 0.72vw, 0.5rem)",
                            whiteSpace: "nowrap",
                          }}
                        >
                          {member.contributed_cp.toLocaleString()} CP
                        </span>
                      </motion.div>
                    );
                  })}
                </div>
              </motion.div>
            )}
          </div>
        );
      })}
    </motion.section>
  );
}
