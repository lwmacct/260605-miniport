import {
  AppShell as WorkbenchShell,
  ThemeToggle as WorkbenchThemeToggle,
  UserMenu as WorkbenchUserMenu,
  type WorkbenchNavEntry,
} from "@lwmacct/260627-antd-workbench";
import { Space } from "antd";
import { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { appPaths, topNavFromPathname, type TopNavKey } from "../router/navigation";
import { APP_NAME, DISPLAY_VERSION } from "@/shared/config/appConfig";
import { useAuthStateQuery, useLogoutMutation } from "@/modules/auth";

function navItems(admin: boolean): WorkbenchNavEntry[] {
  const items: WorkbenchNavEntry[] = [
    { key: "overview", label: "端口总览" },
    { key: "services", label: "端口分配" },
    { key: "projects", label: "项目服务" },
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
  projects: appPaths.projects,
  overview: appPaths.overview,
  services: appPaths.services,
};

export function AppShell() {
  const navigate = useNavigate();
  const location = useLocation();
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
          <WorkbenchThemeToggle />
          <WorkbenchUserMenu
            user={{ name: user?.username, username: user?.username }}
            onLogout={() => void logoutMutation.mutateAsync()}
            onOpenAccount={() => navigate(appPaths.admin)}
          />
        </Space>
      }
      brand={{
        mark: "M",
        name: APP_NAME,
        subtitle: "端口服务资产管理",
        version: DISPLAY_VERSION,
      }}
      flushContent={isFlushContent}
      nav={navItems(Boolean(user?.admin))}
      selectedNavKey={visibleActiveKey}
      onSelectNav={handleNavigate}
    >
      <Outlet />
    </WorkbenchShell>
  );
}
