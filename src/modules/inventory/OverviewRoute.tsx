import { InventoryWorkspace } from "./ui/InventoryWorkspace";

export function OverviewRoute() {
	return (
		<InventoryWorkspace
			description="按 IP 和 10 端口组维护服务、容器、依赖和仓库。"
			title="端口总览"
			view="overview"
		/>
	);
}
