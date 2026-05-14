import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import {
  AudioSettingsContext,
  audioSettingsKey,
  type AudioSettingsContextValue,
} from "./audioSettingsContext";

interface StoredAudioSettings {
  bgmEnabled?: unknown;
  seEnabled?: unknown;
}

const defaultAudioSettings = {
  bgmEnabled: true,
  seEnabled: true,
};

const readStoredAudioSettings = () => {
  if (typeof window === "undefined") {
    return defaultAudioSettings;
  }

  try {
    const rawSettings = window.localStorage.getItem(audioSettingsKey);
    if (!rawSettings) {
      return defaultAudioSettings;
    }

    const parsed = JSON.parse(rawSettings) as StoredAudioSettings;

    return {
      bgmEnabled:
        typeof parsed.bgmEnabled === "boolean"
          ? parsed.bgmEnabled
          : defaultAudioSettings.bgmEnabled,
      seEnabled:
        typeof parsed.seEnabled === "boolean" ? parsed.seEnabled : defaultAudioSettings.seEnabled,
    };
  } catch {
    return defaultAudioSettings;
  }
};

interface AudioSettingsProviderProps {
  children: ReactNode;
}

export function AudioSettingsProvider({ children }: AudioSettingsProviderProps) {
  const [settings, setSettings] = useState(readStoredAudioSettings);

  useEffect(() => {
    try {
      window.localStorage.setItem(
        audioSettingsKey,
        JSON.stringify({
          bgmEnabled: settings.bgmEnabled,
          seEnabled: settings.seEnabled,
        }),
      );
    } catch (error) {
      console.error("failed to persist audio settings", error);
    }
  }, [settings.bgmEnabled, settings.seEnabled]);

  const toggleBgm = useCallback(() => {
    setSettings((current) => ({
      ...current,
      bgmEnabled: !current.bgmEnabled,
    }));
  }, []);

  const toggleSe = useCallback(() => {
    setSettings((current) => ({
      ...current,
      seEnabled: !current.seEnabled,
    }));
  }, []);

  const value = useMemo<AudioSettingsContextValue>(
    () => ({
      isBgmEnabled: settings.bgmEnabled,
      isSeEnabled: settings.seEnabled,
      toggleBgm,
      toggleSe,
    }),
    [settings.bgmEnabled, settings.seEnabled, toggleBgm, toggleSe],
  );

  return <AudioSettingsContext.Provider value={value}>{children}</AudioSettingsContext.Provider>;
}
