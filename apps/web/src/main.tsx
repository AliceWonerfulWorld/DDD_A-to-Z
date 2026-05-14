import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import "./index.css";
import { AppRoutes } from "./AppRoutes.tsx";
import { AudioSettingsProvider } from "./features/audio/AudioSettingsProvider.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <AudioSettingsProvider>
        <AppRoutes />
      </AudioSettingsProvider>
    </BrowserRouter>
  </StrictMode>,
);
