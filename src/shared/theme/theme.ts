export type ThemeMode = "light" | "dark";

export const themeStorageKey = "miniport-theme";

function isThemeMode(value: string | undefined | null): value is ThemeMode {
  return value === "light" || value === "dark";
}

function preferredTheme(): ThemeMode {
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

export function readInitialTheme(): ThemeMode {
  const bootTheme = document.documentElement.dataset.theme;
  if (isThemeMode(bootTheme)) {
    return bootTheme;
  }

  const stored = window.localStorage.getItem(themeStorageKey);
  if (isThemeMode(stored)) {
    return stored;
  }

  return preferredTheme();
}

export function applyTheme(mode: ThemeMode) {
	document.documentElement.dataset.theme = mode;
	document.documentElement.style.colorScheme = mode;
	window.localStorage.setItem(themeStorageKey, mode);
}
