import { useEffect, useState } from "react";
import { applyTheme, readInitialTheme, type ThemeMode } from "./theme";

export function useThemeMode() {
  const [themeMode, setThemeMode] = useState<ThemeMode>(() => readInitialTheme());

  useEffect(() => {
    applyTheme(themeMode);
  }, [themeMode]);

	function toggleTheme() {
		setThemeMode((current) => (current === "dark" ? "light" : "dark"));
	}

  return { themeMode, toggleTheme };
}
