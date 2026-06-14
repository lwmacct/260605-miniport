import { useQuery } from "@tanstack/react-query";
import { loadInventory } from "./api";
import type { InventoryQuery } from "./types";

export const inventoryKeys = {
  snapshot: (query: InventoryQuery) => ["inventory", "snapshot", query] as const,
};

export function useInventoryQuery(query: InventoryQuery) {
  return useQuery({
    queryKey: inventoryKeys.snapshot(query),
    queryFn: () => loadInventory(query),
  });
}
