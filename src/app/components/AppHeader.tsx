import { DownOutlined } from "@ant-design/icons";
import { Button, Dropdown, Menu } from "antd";
import type { MenuProps } from "antd";
import type { ReactNode } from "react";
import styles from "./AppHeader.module.css";

interface AppHeaderProps {
  actions?: ReactNode;
  activeKey: string;
  brandName: string;
  navItems: MenuProps["items"];
  onNavigate(key: string): void;
  version?: string;
}

export function AppHeader({ actions, activeKey, brandName, navItems, onNavigate, version }: AppHeaderProps) {
  const activeLabel = findActiveLabel(navItems, activeKey);
  const displayVersion = shortenVersion(version);

  return (
    <header className={styles.header}>
      <div className={styles.brandSlot} aria-label={brandName}>
        <div className={styles.brandMark}>M</div>
        <div className={styles.brandMeta}>
          <strong>{brandName}</strong>
          <span>{displayVersion ?? "端口服务资产管理"}</span>
        </div>
      </div>

      <nav className={styles.navSlot} aria-label="主导航">
        <Menu
          className={styles.fullNav}
          disabledOverflow
          items={navItems}
          mode="horizontal"
          onClick={({ key }) => onNavigate(key)}
          selectedKeys={[activeKey]}
        />
        <Dropdown
          menu={{
            items: navItems,
            onClick: ({ key }) => onNavigate(key),
            selectedKeys: [activeKey],
          }}
          placement="bottom"
          trigger={["click"]}
        >
          <Button className={styles.compactNav} icon={<DownOutlined />}>
            {activeLabel}
          </Button>
        </Dropdown>
      </nav>

      <div className={styles.actionsSlot}>{actions}</div>
    </header>
  );
}

function shortenVersion(version: string | undefined): string | undefined {
  return version?.split(/[+-]/, 1)[0];
}

function findActiveLabel(items: MenuProps["items"], activeKey: string): ReactNode {
  const activeItem = items?.find((item) => item && "key" in item && item.key === activeKey);
  if (activeItem && "label" in activeItem) {
    return activeItem.label;
  }
  return "导航";
}
