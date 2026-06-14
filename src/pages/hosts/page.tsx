import { InventoryWorkspace } from "../../domains/inventory/components/InventoryWorkspace";

export function HostsPage() {
  return (
    <InventoryWorkspace
      description="管理 IP 主机、网段和每台主机下的端口组数量。"
      title="主机管理"
      view="hosts"
    />
  );
}
