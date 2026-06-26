import {
  ApiOutlined,
  AuditOutlined,
  SettingOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { SectionLayout } from "@/shared/layouts/SectionLayout";

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
    <SectionLayout
      activeKey={activeSection(location.pathname)}
      menuItems={[
        {
          type: "group",
          label: "系统管理",
          children: [...sectionItems],
        },
      ]}
      onChange={(key) => navigate(`/admin/${key}`)}
    >
      <Outlet />
    </SectionLayout>
  );
}
