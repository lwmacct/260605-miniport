import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ProjectServicesRoute() {
	return (
		<PortsvcWorkspace
			description="按端口组查看项目、服务和 DIND 容器归属。"
			title="项目服务"
			view="projects"
		/>
	);
}
