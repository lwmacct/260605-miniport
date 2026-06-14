import { createContext, useContext } from "react";
import type { ThemeMode } from "./theme";

interface ThemeModeContextValue {
  themeMode: ThemeMode;
  toggleTheme(checked: boolean): void;
}

export const ThemeModeContext = createContext<ThemeModeContextValue | null>(null);

export function useThemeModeContext() {
  const value = useContext(ThemeModeContext);
  if (value === null) {
    throw new Error("ThemeModeContext is not available");
  }
  return value;
}
