import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function DependenciesRoute() {
  return (
    <PortsvcWorkspace
      description="管理自有仓库、开源项目、闭源服务和外部依赖资产。"
      title="依赖资产"
      view="dependencies"
    />
  );
}
