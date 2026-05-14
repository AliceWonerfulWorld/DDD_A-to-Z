import { useContext } from "react";
import { AudioSettingsContext } from "./audioSettingsContext";

export function useAudioSettings() {
  const context = useContext(AudioSettingsContext);

  if (!context) {
    throw new Error("useAudioSettings must be used within AudioSettingsProvider");
  }

  return context;
}
