import {
  ApiOutlined,
  AppstoreOutlined,
  ClusterOutlined,
  PartitionOutlined,
} from "@ant-design/icons";
import { WorkbenchSectionLayout } from "@lwmacct/260627-antd-workbench";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

type ConsoleSectionKey = "dependencies" | "overview" | "projects" | "services";

const sectionItems = [
  { key: "overview", label: "端口总览", icon: <AppstoreOutlined /> },
  { key: "services", label: "端口分配", icon: <PartitionOutlined /> },
  { key: "projects", label: "项目服务", icon: <ClusterOutlined /> },
  { key: "dependencies", label: "依赖与仓库", icon: <ApiOutlined /> },
] as const;

const sectionKeys = new Set<ConsoleSectionKey>(sectionItems.map((item) => item.key));

function activeSection(pathname: string): ConsoleSectionKey {
  const key = pathname.split("/")[2];

  if (sectionKeys.has(key as ConsoleSectionKey)) {
    return key as ConsoleSectionKey;
  }

  return "overview";
}

export function ConsoleLayout() {
  const location = useLocation();
  const navigate = useNavigate();

  return (
    <WorkbenchSectionLayout
      selectedKey={activeSection(location.pathname)}
      nav={[
        {
          type: "group",
          key: "operations",
          label: "Console",
          children: [...sectionItems],
        },
      ]}
      onSelect={(key) => navigate(`/console/${key}`)}
    >
      <Outlet />
    </WorkbenchSectionLayout>
  );
}
