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
import { fetchMe } from "./features/auth/api.ts";
import { markInitialProfileCompleted } from "./features/profile/initialProfile.ts";

export function AppRoutes() {
  const navigate = useNavigate();
  const location = useLocation();
  const usesSharedGuildBgm =
    location.pathname === "/guild" ||
    location.pathname === "/guild/details" ||
    location.pathname === "/guild/my-guild";
  const usesSharedHomeBgm = location.pathname === "/home" || location.pathname === "/mypage";
  const completeInitialProfile = async (username: string) => {
    if (username.trim() === "") return;

    try {
      const user = await fetchMe();
      if (user) {
        markInitialProfileCompleted(user.id);
      }
    } catch (error) {
      console.error("failed to complete initial profile", error);
    } finally {
      navigate("/analysis");
    }
  };

  return (
    <>
      {usesSharedHomeBgm && <HomeBgm />}
      {usesSharedGuildBgm && <GuildBgm />}
      <Routes>
        <Route path="/" element={<App />} />
        <Route
          path="/profile"
          element={
            <InitialProfile onComplete={(username) => void completeInitialProfile(username)} />
          }
        />
        <Route
          path="/analysis"
          element={<ContributionAnalysis onComplete={() => navigate("/home")} />}
        />
        <Route path="/home" element={<Home onNavigate={navigate} />} />
        <Route path="/mypage" element={<MyPage onNavigate={navigate} />} />
        <Route path="/guild" element={<GuildDashboard onNavigate={navigate} />} />
        <Route path="/guild/details" element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path="/guild/my-guild" element={<MyGuildDetails onNavigate={navigate} />} />
        <Route path="/guild/town" element={<GuildTown onNavigate={navigate} />} />
        <Route path="/war" element={<Navigate to="/mypage" replace />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  );
}
