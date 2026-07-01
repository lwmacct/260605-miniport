import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function OverviewRoute() {
  return (
    <PortsvcWorkspace
      description="汇总宿主机、端口组、运行环境、服务组和依赖资产数量。"
      title="资源总览"
      view="overview"
    />
  );
}
