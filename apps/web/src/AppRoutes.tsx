import { Navigate, Route, Routes, useLocation, useNavigate } from "react-router";
import App from "./App.tsx";
import { ContributionAnalysis } from "./components/analysis/ContributionAnalysis.tsx";
import { GuildDashboard } from "./components/guild-dashboard/GuildDashboard.tsx";
import { MyGuildDetails } from "./components/guild-details/MyGuildDetails.tsx";
import { GuildSelection } from "./components/guild-selection/GuildSelection.tsx";
import { GuildTown } from "./components/guild-town/GuildTown.tsx";
import { Home } from "./components/home/Home.tsx";
import { MyPage } from "./components/my-page/MyPage.tsx";
import { InitialProfile } from "./components/profile/InitialProfile.tsx";
import { GuildBgm } from "./components/shared/GuildBgm.tsx";
import { HomeBgm } from "./components/shared/HomeBgm.tsx";
import { WarMap } from "./components/war-map/WarMap.tsx";
import { PATHS } from "./constants/paths.ts";
import { fetchMe } from "./features/auth/api.ts";
import { hasSelectedGuildMembership } from "./features/guild/membership.ts";
import { completeInitialProfileAPI } from "./features/profile/api.ts";
import { markInitialProfileCompleted } from "./features/profile/initialProfile.ts";

export function AppRoutes() {
  const navigate = useNavigate();
  const location = useLocation();
  const usesSharedGuildBgm =
    location.pathname === PATHS.GUILD ||
    location.pathname === PATHS.GUILD_SELECT ||
    location.pathname === PATHS.GUILD_DETAILS ||
    location.pathname === PATHS.GUILD_MY_GUILD;
  const usesSharedHomeBgm = location.pathname === PATHS.HOME || location.pathname === PATHS.MY_PAGE;
  const completeInitialProfile = async (username: string) => {
    if (username.trim() === "") return;

    try {
      const user = await fetchMe();
      if (user) {
        await completeInitialProfileAPI(username);
        markInitialProfileCompleted(user.id);
      }
    } catch (error) {
      console.error("failed to complete initial profile", error);
    } finally {
      navigate(PATHS.ANALYSIS);
    }
  };

  return (
    <>
      {usesSharedHomeBgm && <HomeBgm />}
      {usesSharedGuildBgm && <GuildBgm />}
      <Routes>
        <Route path={PATHS.ROOT} element={<App />} />
        <Route
          path={PATHS.PROFILE}
          element={
            <InitialProfile onComplete={(username) => void completeInitialProfile(username)} />
          }
        />
        <Route
          path={PATHS.ANALYSIS}
          element={<ContributionAnalysis onComplete={() => navigate(PATHS.HOME)} />}
        />
        <Route path={PATHS.HOME} element={<Home onNavigate={navigate} />} />
        <Route path={PATHS.MY_PAGE} element={<MyPage onNavigate={navigate} />} />
        <Route
          path={PATHS.GUILD}
          element={
            hasSelectedGuildMembership() ? (
              <GuildDashboard onNavigate={navigate} />
            ) : (
              <Navigate to={PATHS.GUILD_SELECT} replace />
            )
          }
        />
        <Route
          path={PATHS.GUILD_SELECT}
          element={<GuildSelection onNavigate={(path) => void navigate(path)} />}
        />
        <Route path={PATHS.GUILD_DETAILS} element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path={PATHS.GUILD_MY_GUILD} element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path={PATHS.GUILD_TOWN} element={<GuildTown onNavigate={navigate} />} />
        <Route path={PATHS.WAR} element={<WarMap onNavigate={navigate} />} />
        <Route path="*" element={<Navigate to={PATHS.ROOT} replace />} />
      </Routes>
    </>
  );
}
