import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ServiceGroupsRoute() {
  return (
    <PortsvcWorkspace
      description="管理由多个运行环境组成的逻辑服务组或集群。"
      title="服务组"
      view="serviceGroups"
    />
  );
}
