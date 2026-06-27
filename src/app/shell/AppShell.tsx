import { WorkbenchShell } from "@lwmacct/260627-antd-workbench";
import { Space, type MenuProps } from "antd";
import { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { appPaths, topNavFromPathname, type TopNavKey } from "../router/navigation";
import { APP_NAME, DISPLAY_VERSION } from "@/shared/config/appConfig";
import { useThemeModeContext } from "@/shared/theme/ThemeModeContext";
import { useAuthStateQuery, useLogoutMutation, UserMenu } from "@/modules/auth";

function navItems(admin: boolean): MenuProps["items"] {
  const items: MenuProps["items"] = [
    { key: "overview", label: "端口总览" },
    { key: "services", label: "服务列表" },
    { key: "hosts", label: "主机管理" },
    { key: "dependencies", label: "依赖与仓库" },
  ];
  if (admin) {
    items.push({ key: "admin", label: "管理员" });
  }
  return items;
}

const navTargets: Record<TopNavKey, string> = {
  admin: appPaths.admin,
  dependencies: appPaths.dependencies,
  hosts: appPaths.hosts,
  overview: appPaths.overview,
  services: appPaths.services,
};

export function AppShell() {
  const navigate = useNavigate();
  const location = useLocation();
  const { themeMode, toggleTheme } = useThemeModeContext();
  const authState = useAuthStateQuery();
  const logoutMutation = useLogoutMutation();
  const activeNavKey = topNavFromPathname(location.pathname);
  const isFlushContent = activeNavKey === "admin";
  const [optimisticActiveKey, setOptimisticActiveKey] = useState<TopNavKey>();
  const user = authState.data?.session.user;
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
      actions={
        <Space>
          <UserMenu
            themeMode={themeMode}
            username={user?.username}
            onLogout={() => void logoutMutation.mutateAsync()}
            onOpenAccount={() => navigate(appPaths.admin)}
            onToggleTheme={toggleTheme}
          />
        </Space>
      }
      activeNavKey={visibleActiveKey}
      brand={{
        mark: "M",
        name: APP_NAME,
        subtitle: "端口服务资产管理",
        version: DISPLAY_VERSION,
      }}
      flushContent={isFlushContent}
      navItems={navItems(Boolean(user?.admin))}
      onNavigate={handleNavigate}
    >
      <Outlet />
    </WorkbenchShell>
  );
}
