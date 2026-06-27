import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function DependenciesRoute() {
	return (
		<PortsvcWorkspace
			description="按服务汇总组件和代码仓库，方便做依赖梳理。"
			title="依赖与仓库"
			view="dependencies"
		/>
	);
}
