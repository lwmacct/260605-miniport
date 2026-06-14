import { MoonOutlined, SunOutlined } from "@ant-design/icons";
import { Layout, Switch, Tooltip, type MenuProps } from "antd";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { AppHeader } from "./AppHeader";
import { appPaths, topNavFromPathname, type TopNavKey } from "../router/navigation";
import { APP_NAME, DISPLAY_VERSION } from "../../shared/config/appConfig";
import { useThemeModeContext } from "../../shared/theme/ThemeModeContext";
import styles from "./AppShell.module.css";

function navItems(): MenuProps["items"] {
  return [
    { key: "overview", label: "端口总览" },
    { key: "services", label: "服务列表" },
    { key: "hosts", label: "主机管理" },
    { key: "dependencies", label: "依赖与仓库" },
  ];
}

const navTargets: Record<TopNavKey, string> = {
  dependencies: appPaths.dependencies,
  hosts: appPaths.hosts,
  overview: appPaths.overview,
  services: appPaths.services,
};

export function AppShell() {
  const navigate = useNavigate();
  const location = useLocation();
  const { themeMode, toggleTheme } = useThemeModeContext();
  const activeNavKey = topNavFromPathname(location.pathname);

  return (
    <Layout className={styles.shell}>
      <AppHeader
        actions={
          <Tooltip title={themeMode === "dark" ? "切换到明亮模式" : "切换到暗色模式"}>
            <Switch
              checked={themeMode === "dark"}
              checkedChildren={<MoonOutlined />}
              unCheckedChildren={<SunOutlined />}
              onChange={toggleTheme}
            />
          </Tooltip>
        }
        activeNavKeys={[activeNavKey]}
        brandName={APP_NAME}
        navItems={navItems()}
        onNavigate={(key) => navigate(navTargets[key as TopNavKey])}
        version={DISPLAY_VERSION}
      />
      <Layout.Content className={styles.content}>
        <Outlet />
      </Layout.Content>
    </Layout>
  );
}
