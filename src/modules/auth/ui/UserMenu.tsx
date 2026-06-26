import { LogoutOutlined, MoonOutlined, SettingOutlined, SunOutlined, UserOutlined } from "@ant-design/icons";
import { Avatar, Button, Dropdown, type MenuProps } from "antd";

interface UserMenuProps {
	onOpenAccount?(): void;
	onToggleTheme?(): void;
	themeMode?: "light" | "dark";
	username?: string;
	onLogout(): void;
}

export function UserMenu({ username, themeMode, onLogout, onOpenAccount, onToggleTheme }: UserMenuProps) {
	const initial = (username || "?").trim().slice(0, 1).toUpperCase();
	const items: MenuProps["items"] = [
    {
      disabled: true,
      icon: <UserOutlined />,
      key: "user",
      label: username || "未命名用户",
	},
	{ type: "divider" },
	onOpenAccount
		? {
				icon: <SettingOutlined />,
				key: "settings",
				label: "设置",
				onClick: onOpenAccount,
			}
		: null,
	onToggleTheme
		? {
				icon: themeMode === "dark" ? <MoonOutlined /> : <SunOutlined />,
				key: "theme",
				label: themeMode === "dark" ? "切换浅色模式" : "切换深色模式",
				onClick: onToggleTheme,
			}
		: null,
	{
		icon: <LogoutOutlined />,
      key: "logout",
      label: "退出登录",
      onClick: onLogout,
    },
  ];

  return (
    <Dropdown menu={{ items }} placement="bottomRight" trigger={["click"]}>
      <Button aria-label="用户菜单" shape="circle" type="text">
        <Avatar size={28} style={{ background: "var(--app-accent)" }}>
          {initial}
        </Avatar>
      </Button>
    </Dropdown>
  );
}
