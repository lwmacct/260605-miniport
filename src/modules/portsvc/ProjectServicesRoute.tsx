import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ProjectServicesRoute() {
  return (
    <PortsvcWorkspace
      description="按 10 端口组查看服务组件、服务 IP 和 DIND/宿主机归属。"
      title="运行环境"
      view="projects"
    />
  );
}
