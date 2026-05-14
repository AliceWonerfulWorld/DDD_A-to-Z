import { createContext } from "react";

export const audioSettingsKey = "lang-war.audio-settings";

export interface AudioSettingsContextValue {
  isBgmEnabled: boolean;
  isSeEnabled: boolean;
  toggleBgm: () => void;
  toggleSe: () => void;
}

export const AudioSettingsContext = createContext<AudioSettingsContextValue | null>(null);
