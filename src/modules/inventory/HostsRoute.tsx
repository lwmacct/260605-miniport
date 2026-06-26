import { InventoryWorkspace } from "./ui/InventoryWorkspace";

export function HostsRoute() {
	return (
		<InventoryWorkspace
			description="管理 IP 主机、网段和每台主机下的端口组数量。"
			title="主机管理"
			view="hosts"
		/>
	);
}
