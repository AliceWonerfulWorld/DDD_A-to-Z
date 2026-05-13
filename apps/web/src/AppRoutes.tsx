import { Navigate, Route, Routes, useNavigate } from "react-router";
import App from "./App.tsx";
import { ContributionAnalysis } from "./components/ContributionAnalysis.tsx";
import { GuildDashboard } from "./components/GuildDashboard.tsx";
import { Home } from "./components/Home.tsx";
import { InitialProfile } from "./components/InitialProfile.tsx";
import { MyPage } from "./components/MyPage.tsx";

export function AppRoutes() {
  const navigate = useNavigate();

  return (
    <Routes>
      <Route path="/" element={<App />} />
      <Route
        path="/profile"
        element={<InitialProfile onComplete={() => navigate("/analysis")} />}
      />
      <Route
        path="/analysis"
        element={<ContributionAnalysis onComplete={() => navigate("/home")} />}
      />
      <Route path="/home" element={<Home onNavigate={navigate} />} />
      <Route path="/mypage" element={<MyPage onNavigate={navigate} />} />
      <Route path="/guild" element={<GuildDashboard onNavigate={navigate} />} />
      <Route path="/war" element={<Navigate to="/mypage" replace />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
