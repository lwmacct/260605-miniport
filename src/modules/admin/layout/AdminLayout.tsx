import {
  ApiOutlined,
  AuditOutlined,
  SettingOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { SectionLayout as WorkbenchSectionLayout } from "@lwmacct/260627-antd-workbench";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

type AdminSectionKey = "audit" | "integrations" | "system" | "users";

const sectionItems = [
  { key: "users", label: "用户管理", icon: <TeamOutlined /> },
  { key: "audit", label: "安全审计", icon: <AuditOutlined /> },
  { key: "system", label: "系统设置", icon: <SettingOutlined /> },
  { key: "integrations", label: "集成配置", icon: <ApiOutlined /> },
] as const;

const sectionKeys = new Set<AdminSectionKey>(sectionItems.map((item) => item.key));

function activeSection(pathname: string): AdminSectionKey {
  const key = pathname.split("/")[2];

  if (sectionKeys.has(key as AdminSectionKey)) {
    return key as AdminSectionKey;
  }

  return "users";
}

export function AdminLayout() {
  const location = useLocation();
  const navigate = useNavigate();

  return (
    <WorkbenchSectionLayout
      selectedKey={activeSection(location.pathname)}
      nav={[
        {
          type: "group",
          key: "system-management",
          label: "系统管理",
          children: [...sectionItems],
        },
      ]}
      onSelect={(key) => navigate(`/admin/${key}`)}
    >
      <Outlet />
    </WorkbenchSectionLayout>
  );
}
