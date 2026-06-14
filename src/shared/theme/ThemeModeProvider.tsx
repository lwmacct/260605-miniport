import type { PropsWithChildren } from "react";
import { ThemeModeContext } from "./ThemeModeContext";
import { useThemeMode } from "./useThemeMode";

export function ThemeModeProvider({ children }: PropsWithChildren) {
  const value = useThemeMode();

  return <ThemeModeContext.Provider value={value}>{children}</ThemeModeContext.Provider>;
}
