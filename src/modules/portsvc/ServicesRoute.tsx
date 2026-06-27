import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ServicesRoute() {
	return (
		<PortsvcWorkspace
			description="集中查看端口组、DIND IP、项目和容器分配。"
			title="端口分配"
			view="services"
		/>
	);
}
