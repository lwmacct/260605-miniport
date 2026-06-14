import { InventoryWorkspace } from "../../domains/inventory/components/InventoryWorkspace";

export function ServicesPage() {
  return (
    <InventoryWorkspace
      description="集中查看服务、端口范围、组件和容器分配。"
      title="服务列表"
      view="services"
    />
  );
}
