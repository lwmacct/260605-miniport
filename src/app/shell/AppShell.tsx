import {
  WorkbenchAppearanceButton,
  WorkbenchLanguageToggle,
  WorkbenchShell,
  WorkbenchUserMenu,
  useWorkbenchLocale,
  type WorkbenchNavEntry,
} from "@lwmacct/260627-antd-workbench";
import { Space } from "antd";
import { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { GithubOutlined, KeyOutlined } from "@ant-design/icons";
import { appPaths, topNavFromPathname, type TopNavKey } from "../router/navigation";
import { useAuth } from "../auth";
import { APP_NAME, DISPLAY_VERSION } from "@/shared/config/appConfig";

function navItems(locale: string): WorkbenchNavEntry[] {
  const isZh = locale.startsWith("zh");
  return [
    { key: "console", label: isZh ? "控制台" : "Console" },
    { key: "settings", label: isZh ? "设置" : "Settings" },
  ];
}

const navTargets: Record<TopNavKey, string> = {
  console: appPaths.console,
  settings: appPaths.settings,
};

export function AppShell() {
  const navigate = useNavigate();
  const location = useLocation();
  const { identity, logout } = useAuth();
  const { locale } = useWorkbenchLocale();
  const activeNavKey = topNavFromPathname(location.pathname);
  const isFlushContent = activeNavKey === "console" || activeNavKey === "settings";
  const [optimisticActiveKey, setOptimisticActiveKey] = useState<TopNavKey>();
  const visibleActiveKey = optimisticActiveKey ?? activeNavKey;

  useEffect(() => {
    setOptimisticActiveKey(undefined);
  }, [location.pathname]);

  function handleNavigate(key: string) {
    const navKey = key as TopNavKey;
    const target = navTargets[navKey];
    if (!target) {
      return;
    }
    setOptimisticActiveKey(navKey);
    if (target !== location.pathname) {
      navigate(target);
    }
  }

  return (
    <WorkbenchShell
      account={
        <WorkbenchUserMenu
          user={{
            avatarUrl: identity.avatar_url,
            displayName: identity.name,
            provider: identity.provider === "github" ? "GitHub" : "Access token",
            providerIcon: identity.provider === "github" ? <GithubOutlined /> : <KeyOutlined />,
            username: identity.username,
          }}
          onLogout={logout}
        />
      }
      brand={{
        mark: "M",
        name: APP_NAME,
        version: DISPLAY_VERSION,
      }}
      flushContent={isFlushContent}
      nav={navItems(locale)}
      selectedNavKey={visibleActiveKey}
      utilities={
        <Space>
          <WorkbenchAppearanceButton />
          <WorkbenchLanguageToggle />
        </Space>
      }
      onSelectNav={handleNavigate}
    >
      <Outlet />
    </WorkbenchShell>
  );
}
