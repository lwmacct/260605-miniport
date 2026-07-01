import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function ServicesRoute() {
	return (
		<PortsvcWorkspace
			description="集中查看端口组、运行位置、服务 IP 和端口槽位。"
			title="端口分配"
			view="services"
		/>
	);
}
