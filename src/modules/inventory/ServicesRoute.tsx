import { InventoryWorkspace } from "./ui/InventoryWorkspace";

export function ServicesRoute() {
	return (
		<InventoryWorkspace
			description="集中查看服务、端口范围、组件和容器分配。"
			title="服务列表"
			view="services"
		/>
	);
}
