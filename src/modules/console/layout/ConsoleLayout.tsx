import {
  ApiOutlined,
  AppstoreOutlined,
  ClusterOutlined,
  CloudServerOutlined,
  PartitionOutlined,
  ShareAltOutlined,
} from "@ant-design/icons";
import { WorkbenchSectionLayout } from "@lwmacct/260627-antd-workbench";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

type ConsoleSectionKey = "dependencies" | "hosts" | "overview" | "port-groups" | "projects" | "service-groups";

const overviewItems = [
  { key: "overview", label: "资源总览", icon: <AppstoreOutlined /> },
] as const;

const resourceItems = [
  { key: "hosts", label: "宿主机", icon: <CloudServerOutlined /> },
  { key: "port-groups", label: "端口组", icon: <PartitionOutlined /> },
  { key: "projects", label: "运行环境", icon: <ClusterOutlined /> },
  { key: "service-groups", label: "服务组", icon: <ShareAltOutlined /> },
] as const;

const dependencyItems = [
  { key: "dependencies", label: "依赖资产", icon: <ApiOutlined /> },
] as const;

const sectionKeys = new Set<ConsoleSectionKey>(
  [...overviewItems, ...resourceItems, ...dependencyItems].map((item) => item.key),
);

function activeSection(pathname: string): ConsoleSectionKey {
  const key = pathname.split("/")[2];

  if (key === "services") {
    return "projects";
  }

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
          key: "overview-group",
          label: "总览",
          children: [...overviewItems],
        },
        {
          type: "group",
          key: "resource-management",
          label: "资源管理",
          children: [...resourceItems],
        },
        {
          type: "group",
          key: "dependency-management",
          label: "依赖管理",
          children: [...dependencyItems],
        },
      ]}
      onSelect={(key) => navigate(`/console/${key}`)}
    >
      <Outlet />
    </WorkbenchSectionLayout>
  );
}
