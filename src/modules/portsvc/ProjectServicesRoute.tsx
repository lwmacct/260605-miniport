import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ProjectServicesRoute() {
	return (
		<PortsvcWorkspace
			description="按端口组查看项目、服务组件和 DIND/宿主机归属。"
			title="项目服务"
			view="projects"
		/>
	);
}
