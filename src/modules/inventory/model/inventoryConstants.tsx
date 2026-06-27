import { AppstoreOutlined, DeploymentUnitOutlined, LinkOutlined, TableOutlined } from "@ant-design/icons";

export const statusOptions = [
  { value: "planned", label: "规划中" },
  { value: "running", label: "运行中" },
  { value: "reserved", label: "保留" },
  { value: "stopped", label: "停用" },
];

export const slotStatusOptions = [
  { value: "empty", label: "空闲" },
  { value: "used", label: "已用" },
  { value: "reserved", label: "保留" },
];

export const componentTypeOptions = [
  { value: "opensource", label: "开源项目" },
  { value: "external", label: "外部组件" },
  { value: "internal", label: "自有组件" },
];

export const repositoryKindOptions = [
  { value: "source", label: "源码" },
  { value: "backend", label: "后端" },
  { value: "frontend", label: "前端" },
  { value: "deploy", label: "部署" },
  { value: "infra", label: "基础设施" },
];

export const navItems = [
  { key: "overview", icon: <AppstoreOutlined />, label: "端口总览" },
  { key: "services", icon: <TableOutlined />, label: "端口分配" },
  { key: "hosts", icon: <DeploymentUnitOutlined />, label: "项目服务" },
  { key: "dependencies", icon: <LinkOutlined />, label: "依赖与仓库" },
];
