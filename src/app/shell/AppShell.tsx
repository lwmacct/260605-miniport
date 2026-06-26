import { Layout, Space, type MenuProps } from "antd";
import { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { AppHeader } from "./AppHeader";
import { appPaths, topNavFromPathname, type TopNavKey } from "../router/navigation";
import { APP_NAME, DISPLAY_VERSION } from "@/shared/config/appConfig";
import { useThemeModeContext } from "@/shared/theme/ThemeModeContext";
import { useAuthStateQuery, useLogoutMutation, UserMenu } from "@/modules/auth";
import styles from "./AppShell.module.css";

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
    <Layout className={styles.shell}>
      <AppHeader
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
        activeNavKeys={[visibleActiveKey]}
        brandName={APP_NAME}
        navItems={navItems(Boolean(user?.admin))}
        onNavigate={handleNavigate}
        version={DISPLAY_VERSION}
      />
      <Layout.Content className={isFlushContent ? `${styles.content} ${styles.contentFlush}` : styles.content}>
        <Outlet />
      </Layout.Content>
    </Layout>
  );
}
