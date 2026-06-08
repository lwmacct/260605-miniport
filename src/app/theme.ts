import { theme } from "antd";
import type { ThemeConfig } from "antd";

export type ThemeMode = "light" | "dark";

export const themeStorageKey = "miniport-theme";

export function readStoredTheme(): ThemeMode {
  const stored = window.localStorage.getItem(themeStorageKey);
  if (stored === "light" || stored === "dark") {
    return stored;
  }
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
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
        colorTextBase: "#d4d4d4",
        colorTextSecondary: "#9cdcfe",
      },
      components: {
        Layout: {
          headerBg: "#252526",
          siderBg: "#252526",
        },
        Menu: {
          darkItemBg: "#252526",
          darkSubMenuItemBg: "#252526",
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
      colorPrimary: "#007acc",
      colorBgBase: "#ffffff",
      colorBgLayout: "#f3f3f3",
      colorBgContainer: "#ffffff",
      colorBorder: "#d4d4d4",
      colorTextBase: "#1f2328",
      colorTextSecondary: "#616161",
    },
    components: {
      Layout: {
        headerBg: "#ffffff",
        siderBg: "#f3f3f3",
      },
      Menu: {
        itemBg: "#f3f3f3",
        itemColor: "#424242",
        itemHoverBg: "#e8e8e8",
        itemSelectedBg: "#e5f1fb",
        itemSelectedColor: "#007acc",
      },
    },
  };
}
