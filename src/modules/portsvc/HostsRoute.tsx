import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function HostsRoute() {
  return (
    <PortsvcWorkspace
      description="管理承载 DIND 容器或直跑服务的宿主机。"
      title="宿主机"
      view="hosts"
    />
  );
}
