import { Navigate, Route, Routes, useLocation, useNavigate } from "react-router";
import App from "./App.tsx";
import { ContributionAnalysis } from "./components/ContributionAnalysis.tsx";
import { GuildBgm } from "./components/GuildBgm.tsx";
import { GuildDashboard } from "./components/GuildDashboard.tsx";
import { GuildTown } from "./components/GuildTown.tsx";
import { Home } from "./components/Home.tsx";
import { HomeBgm } from "./components/HomeBgm.tsx";
import { InitialProfile } from "./components/InitialProfile.tsx";
import { MyPage } from "./components/MyPage.tsx";
import { MyGuildDetails } from "./components/MyGuildDetails.tsx";
import { WarMap } from "./components/war-map/WarMap.tsx";
import { PATHS } from "./constants/paths.ts";
import { fetchMe } from "./features/auth/api.ts";
import { completeInitialProfileAPI } from "./features/profile/api.ts";
import { markInitialProfileCompleted } from "./features/profile/initialProfile.ts";

export function AppRoutes() {
  const navigate = useNavigate();
  const location = useLocation();
  const usesSharedGuildBgm =
    location.pathname === PATHS.GUILD ||
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
        <Route path={PATHS.GUILD} element={<GuildDashboard onNavigate={navigate} />} />
        <Route path={PATHS.GUILD_DETAILS} element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path={PATHS.GUILD_MY_GUILD} element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path={PATHS.GUILD_TOWN} element={<GuildTown onNavigate={navigate} />} />
        <Route path={PATHS.WAR} element={<WarMap onNavigate={navigate} />} />
        <Route path="*" element={<Navigate to={PATHS.ROOT} replace />} />
      </Routes>
    </>
  );
}
