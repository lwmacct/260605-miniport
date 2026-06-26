import { DownOutlined } from "@ant-design/icons";
import { Button, Drawer, Layout, Menu, Space, Typography, type MenuProps } from "antd";
import { useState } from "react";
import styles from "./SectionLayout.module.css";

export interface SectionLayoutProps<Key extends string> {
  activeKey: Key;
  children: React.ReactNode;
  menuItems: MenuProps["items"];
  onChange(key: Key): void;
}

export function SectionLayout<Key extends string>({
  activeKey,
  children,
  menuItems,
  onChange,
}: SectionLayoutProps<Key>) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const currentItem = findMenuItem(menuItems, activeKey);

  function handleChange(key: string) {
    onChange(key as Key);
    setMobileOpen(false);
  }

  return (
    <Layout className={styles.page}>
      <Layout.Sider
        className={`${styles.sider} ${styles.desktopSider}`}
        theme="dark"
        width={208}
      >
        <Menu
          className={styles.menu}
          items={menuItems}
          mode="inline"
          selectedKeys={[activeKey]}
          onClick={({ key }) => handleChange(key)}
        />
      </Layout.Sider>
      <Layout.Content className={styles.content}>
        <div className={styles.mobileNav}>
          <div className={styles.mobileNavCurrent}>
            <Space size={8}>
              {currentItem?.icon}
              <Typography.Text strong>{currentItem?.label ?? activeKey}</Typography.Text>
            </Space>
          </div>
          <Button icon={<DownOutlined />} onClick={() => setMobileOpen(true)}>
            切换分区
          </Button>
        </div>
        {children}
      </Layout.Content>
      <Drawer
        className={styles.mobileDrawer}
        onClose={() => setMobileOpen(false)}
        open={mobileOpen}
        placement="bottom"
        title="切换分区"
      >
        <Menu
          className={styles.drawerMenu}
          items={menuItems}
          mode="inline"
          selectedKeys={[activeKey]}
          onClick={({ key }) => handleChange(key)}
        />
      </Drawer>
    </Layout>
  );
}

interface FlattenMenuItem {
  icon?: React.ReactNode;
  key: string;
  label: React.ReactNode;
}

function findMenuItem(items: MenuProps["items"], activeKey: string): FlattenMenuItem | null {
  for (const item of items ?? []) {
    if (!item) {
      continue;
    }
    if ("children" in item && Array.isArray(item.children)) {
      const child = findMenuItem(item.children, activeKey);
      if (child) {
        return child;
      }
    }
    if ("key" in item && item.key === activeKey && "label" in item) {
      return {
        icon: "icon" in item ? item.icon : undefined,
        key: String(item.key),
        label: item.label,
      };
    }
  }
  return null;
}
