import { PortsvcWorkspace } from "./ui/PortsvcWorkspace";

export function PortGroupsRoute() {
  return (
    <PortsvcWorkspace
      description="按 10 个端口为一组统计 1000-5999 的端口组占用。"
      title="端口组"
      view="portGroups"
    />
  );
}
