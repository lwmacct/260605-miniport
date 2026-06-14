import { useQuery } from "@tanstack/react-query";
import { loadInventory } from "./api";

export const inventoryKeys = {
  snapshot: ["inventory", "snapshot"] as const,
};

export function useInventoryQuery() {
  return useQuery({
    queryKey: inventoryKeys.snapshot,
    queryFn: loadInventory,
  });
}
