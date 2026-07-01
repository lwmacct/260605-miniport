import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function PortGroupsRoute() {
  return (
    <PortsvcWorkspace
      description="按 10 个端口为一组统计 10000-59999 的端口组占用。"
      title="端口组"
      view="portGroups"
    />
  );
}
