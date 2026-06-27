import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function OverviewRoute() {
	return (
		<PortsvcWorkspace
			description="按用户全局管理 10000-59999 内的 10 端口组分配。"
			title="端口总览"
			view="overview"
		/>
	);
}
