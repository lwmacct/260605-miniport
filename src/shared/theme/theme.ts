import { theme } from "antd";
import type { ThemeConfig } from "antd";

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

const sharedTokens = {
  borderRadius: 6,
  fontFamily:
    '-apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif',
};

export function appTheme(mode: ThemeMode): ThemeConfig {
  if (mode === "dark") {
    return {
      algorithm: theme.darkAlgorithm,
      token: {
        ...sharedTokens,
        colorPrimary: "#007acc",
        colorBgBase: "#1e1e1e",
        colorBgLayout: "#1e1e1e",
        colorBgContainer: "#252526",
        colorBorder: "#3c3c3c",
        colorTextBase: "#cccccc",
        colorTextSecondary: "#858585",
      },
      components: {
        Card: {
          borderRadiusLG: 4,
          headerHeight: 40,
        },
        Layout: {
          headerBg: "#252526",
          siderBg: "#252526",
        },
        Menu: {
          darkItemBg: "#252526",
          darkSubMenuItemBg: "#252526",
          darkItemColor: "#cccccc",
          darkItemHoverBg: "#2a2d2e",
          darkItemHoverColor: "#ffffff",
          darkItemSelectedBg: "#094771",
          darkItemSelectedColor: "#ffffff",
        },
      },
    };
  }

  return {
    algorithm: theme.defaultAlgorithm,
    token: {
      ...sharedTokens,
      colorPrimary: "#0969da",
      colorBgBase: "#ffffff",
      colorBgLayout: "#f3f3f3",
      colorBgContainer: "#ffffff",
      colorBorder: "#d0d7de",
      colorTextBase: "#24292f",
      colorTextSecondary: "#6e7781",
    },
    components: {
      Card: {
        borderRadiusLG: 4,
        headerHeight: 40,
      },
      Layout: {
        headerBg: "#ffffff",
        siderBg: "#f3f3f3",
      },
      Menu: {
        itemBg: "#f3f3f3",
        itemColor: "#6e7781",
        itemHoverBg: "#e8eef7",
        itemSelectedBg: "#e8eef7",
        itemSelectedColor: "#0969da",
      },
    },
  };
}
