import { useQuery } from "@tanstack/react-query";
import { loadInventory } from "../api/inventoryApi";
import type { InventoryQuery } from "./inventoryTypes";

export const inventoryKeys = {
  snapshot: (query: InventoryQuery) => ["inventory", "snapshot", query] as const,
};

export function useInventoryQuery(query: InventoryQuery) {
  return useQuery({
    queryKey: inventoryKeys.snapshot(query),
    queryFn: () => loadInventory(query),
  });
}
