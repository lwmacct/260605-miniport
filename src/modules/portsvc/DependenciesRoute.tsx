import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function DependenciesRoute() {
	return (
		<PortsvcWorkspace
			description="按端口组汇总服务组件、依赖和代码仓库。"
			title="依赖与仓库"
			view="dependencies"
		/>
	);
}
