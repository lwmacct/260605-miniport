import { InventoryWorkspace } from "../../domains/inventory/components/InventoryWorkspace";

export function OverviewPage() {
  return (
    <InventoryWorkspace
      description="按 IP 和 10 端口组维护服务、容器、依赖和仓库。"
      title="端口总览"
      view="overview"
    />
  );
}
