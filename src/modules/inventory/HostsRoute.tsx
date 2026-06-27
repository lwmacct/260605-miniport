import { InventoryWorkspace } from "./ui/InventoryWorkspace";

export function HostsRoute() {
	return (
		<InventoryWorkspace
			description="按端口组查看项目、服务和 DIND 容器归属。"
			title="项目服务"
			view="hosts"
		/>
	);
}
