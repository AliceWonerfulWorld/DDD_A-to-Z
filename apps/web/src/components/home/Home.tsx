import { useEffect, useState } from "react";
import { HomeHud } from "./HomeHud";
import type { GuildSummary } from "./HomeHud";
import { HomeFirstVisitGuide } from "./HomeFirstVisitGuide";
import { HomeNav } from "./HomeNav";
import { ReturnTitleDialog } from "./ReturnTitleDialog";
import { WalkingGopher } from "./WalkingGopher";
import { PATHS } from "../../constants/paths";
import styles from "./Home.module.css";
import { useHomeAudio } from "../../hooks/useHomeAudio";
import { AUDIO_ASSETS } from "../../features/audio/audioAssets";
import { fetchMyGuild } from "../../features/guild/api";
import { toDisplayGuild } from "../../features/guild/presentation";
import { fetchProfile, type Profile } from "../../features/profile/api";
import { fetchHomeCP, type HomeCPData } from "../../features/home/api";
import { fetchMyPets, type PetSummary } from "../../features/pet/api";
import { readHomePetId } from "../../features/pet/homePet";

interface HomeProps {
  onNavigate: (path: string) => void | Promise<void>;
}

const defaultPlayer = {
  name: "DevSamurai",
  title: "Consistency Master",
  level: 1,
  levelTotalCp: 0,
  totalCp: 0,
  todayCp: 0,
  nextLevel: 2,
  nextLevelTotalCp: 100,
  nextLevelRemainingCp: 100,
  lifetimeTotalEarnedCp: 0,
};

function playerLevelFromTotalEarned(totalEarned: number) {
  if (totalEarned <= 0) {
    return 1;
  }

  return Math.floor(Math.sqrt(totalEarned / 100)) + 1;
}

function totalEarnedForPlayerLevel(level: number) {
  if (level <= 1) {
    return 0;
  }

  return (level - 1) ** 2 * 100;
}

const navItems = [
  { label: "REPO ANALYSIS", caption: "GITHUB SCAN", path: PATHS.ANALYSIS, accent: "#7cff6b" },
  { label: "WAR MAP", caption: "BATTLE FRONT", path: PATHS.WAR, accent: "#ff5f56" },
  { label: "GUILD BASE", caption: "COMMUNITY HQ", path: PATHS.GUILD, accent: "#00f5ff" },
  { label: "MY PET", caption: "GROW & BATTLE", path: PATHS.PETS, accent: "#74f7a1" },
  { label: "MY STATUS", caption: "PLAYER DATA", path: PATHS.MY_PAGE, accent: "#ffd700" },
];
const homeFirstVisitGuideKey = "lang-war.home.first-visit-guide-seen";

export function Home({ onNavigate }: HomeProps) {
  const [isReturnTitleDialogOpen, setIsReturnTitleDialogOpen] = useState(false);
  const [isFirstVisitGuideOpen, setIsFirstVisitGuideOpen] = useState(false);
  const [guild, setGuild] = useState<GuildSummary | null | undefined>(undefined);
  const navigateFromHome = (path: string) => {
    if (path === PATHS.GUILD && guild === undefined) {
      return;
    }
    if (path === PATHS.GUILD && guild === null) {
      onNavigate(PATHS.GUILD_SELECT);
      return;
    }

    onNavigate(path);
  };
  const {
    audioRefs,
    audioError,
    isSeEnabled,
    playGopherTalk,
    playPythonPetTalk,
    playHomeNavSelect,
    playModalCancel,
    playModalOpen,
    playReturnTitle,
  } = useHomeAudio(navigateFromHome);

  const [profile, setProfile] = useState<Profile | null>(null);
  const [cpData, setCpData] = useState<HomeCPData | null>(null);
  const [homePet, setHomePet] = useState<PetSummary | null>(null);
  const lifetimeTotalEarnedCp = cpData?.lifetime_total_earned_cp ?? cpData?.total_cp ?? 0;
  const playerLevel = cpData?.player_level ?? playerLevelFromTotalEarned(lifetimeTotalEarnedCp);
  const playerLevelTotalCp =
    cpData?.player_level_total_cp ?? totalEarnedForPlayerLevel(playerLevel);
  const nextPlayerLevel = cpData?.next_player_level ?? playerLevel + 1;
  const nextPlayerLevelTotalCp =
    cpData?.next_player_level_total_cp ?? totalEarnedForPlayerLevel(nextPlayerLevel);
  const nextPlayerLevelRemainingCp =
    cpData?.next_player_level_remaining ??
    Math.max(0, nextPlayerLevelTotalCp - lifetimeTotalEarnedCp);

  useEffect(() => {
    try {
      setIsFirstVisitGuideOpen(window.localStorage.getItem(homeFirstVisitGuideKey) !== "true");
    } catch {
      setIsFirstVisitGuideOpen(true);
    }
  }, []);

  useEffect(() => {
    fetchProfile().then(setProfile).catch(console.error);
  }, []);

  useEffect(() => {
    fetchHomeCP().then(setCpData).catch(console.error);
  }, []);

  useEffect(() => {
    fetchMyPets()
      .then((data) => {
        const homePetId = readHomePetId();
        const selectedPet =
          data.pets.find((pet) => pet.id === homePetId) ?? data.currentGuildPet ?? data.pets[0];
        setHomePet(selectedPet ?? null);
      })
      .catch((error) => {
        console.error(error);
        setHomePet(null);
      });
  }, []);

  useEffect(() => {
    fetchMyGuild()
      .then((data) => {
        if (!data?.guild) {
          setGuild(null);
          return;
        }

        const displayGuild = toDisplayGuild(data.guild);
        setGuild({
          name: `${displayGuild.name} Guild`,
          icon: displayGuild.icon,
          rank: "Member",
          accent: displayGuild.accent,
        });
      })
      .catch((error) => {
        console.error(error);
        setGuild(null);
      });
  }, []);

  const cancelReturnTitle = () => {
    playModalCancel();
    setIsReturnTitleDialogOpen(false);
  };

  const openReturnTitleDialog = () => {
    playModalOpen();
    setIsReturnTitleDialogOpen(true);
  };

  const dismissFirstVisitGuide = () => {
    try {
      window.localStorage.setItem(homeFirstVisitGuideKey, "true");
    } catch {
      // The guide can still close even if localStorage is unavailable.
    }
    setIsFirstVisitGuideOpen(false);
  };

  return (
    <main
      style={{
        minHeight: "100svh",
        position: "relative",
        overflow: "hidden",
        backgroundImage:
          "linear-gradient(180deg, rgba(4, 8, 18, 0.25), rgba(4, 8, 18, 0.56)), url('/home_lunch.png')",
        backgroundSize: "cover",
        backgroundPosition: "center center",
        backgroundRepeat: "no-repeat",
        imageRendering: "pixelated",
        fontFamily: '"Press Start 2P", "DotGothic16", monospace',
        color: "#f4ecd0",
      }}
    >
      <audio
        ref={audioRefs.homeNavSelectSeRef}
        src={AUDIO_ASSETS.se.homeNavSelect}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.confirmModalSeRef}
        src={AUDIO_ASSETS.se.confirmModal}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.modalCancelSeRef}
        src={AUDIO_ASSETS.se.modalCancel}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.returnTitleSeRef}
        src={AUDIO_ASSETS.se.returnTitle}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.gopherTalkSeRef}
        src={AUDIO_ASSETS.se.gopherTalk}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />
      <audio
        ref={audioRefs.pythonPetTalkSeRef}
        src={AUDIO_ASSETS.se.pythonPetTalk}
        preload="none"
        muted={!isSeEnabled}
        aria-hidden="true"
      />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          backgroundImage:
            "repeating-linear-gradient(0deg, rgba(0,0,0,0.1), rgba(0,0,0,0.1) 1px, transparent 1px, transparent 4px)",
          mixBlendMode: "multiply",
          pointerEvents: "none",
          zIndex: 1,
        }}
      />

      <div
        aria-hidden="true"
        style={{
          position: "fixed",
          inset: 0,
          boxShadow: "inset 0 0 90px rgba(0,0,0,0.7)",
          pointerEvents: "none",
          zIndex: 2,
        }}
      />

      <div
        style={{
          position: "relative",
          zIndex: 3,
          minHeight: "100svh",
          display: "grid",
          gridTemplateRows: "auto minmax(96px, 1fr) auto",
          padding: "clamp(12px, 2.6vmin, 32px)",
          gap: "clamp(12px, 2svh, 20px)",
        }}
      >
        <HomeHud
          guild={guild}
          player={{
            ...defaultPlayer,
            name: profile?.display_name || defaultPlayer.name,
            level: playerLevel,
            levelTotalCp: playerLevelTotalCp,
            totalCp: cpData?.total_cp ?? defaultPlayer.totalCp,
            todayCp: cpData?.today_cp ?? defaultPlayer.todayCp,
            nextLevel: nextPlayerLevel,
            nextLevelTotalCp: nextPlayerLevelTotalCp,
            nextLevelRemainingCp: nextPlayerLevelRemainingCp,
            lifetimeTotalEarnedCp: lifetimeTotalEarnedCp,
          }}
          onReturnTitle={openReturnTitleDialog}
        />

        <section aria-label="Character placement area" className={styles.characterArea}>
          <WalkingGopher
            pet={homePet}
            onTalk={
              homePet?.attribute.toLowerCase() === "python" ? playPythonPetTalk : playGopherTalk
            }
          />
        </section>

        <HomeNav items={navItems} onNavigate={playHomeNavSelect} />
      </div>

      {isReturnTitleDialogOpen && (
        <ReturnTitleDialog onCancel={cancelReturnTitle} onConfirm={playReturnTitle} />
      )}

      {isFirstVisitGuideOpen && <HomeFirstVisitGuide onDismiss={dismissFirstVisitGuide} />}

      {audioError && (
        <div
          role="alert"
          style={{
            position: "fixed",
            right: "clamp(14px, 3vw, 28px)",
            bottom: "clamp(14px, 3vw, 28px)",
            zIndex: 12,
            maxWidth: "min(320px, calc(100vw - 28px))",
            border: "2px solid #ff5f56",
            borderBottomColor: "rgba(0,0,0,0.82)",
            borderRightColor: "rgba(0,0,0,0.82)",
            background: "rgba(18, 6, 10, 0.94)",
            boxShadow: "0 0 0 2px rgba(0,0,0,0.78), 5px 5px 0 rgba(0,0,0,0.48)",
            color: "#fff8d7",
            fontFamily: '"DotGothic16", monospace',
            fontSize: "0.86rem",
            lineHeight: 1.6,
            padding: "10px 12px",
          }}
        >
          {audioError}
        </div>
      )}
    </main>
  );
}
