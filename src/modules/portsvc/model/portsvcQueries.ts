import { useQuery } from "@tanstack/react-query";
import { loadPortsvc } from "../api/portsvcApi";
import type { PortsvcQuery } from "./portsvcTypes";

export const portsvcKeys = {
  snapshot: (query: PortsvcQuery) => ["portsvc", "snapshot", query] as const,
};

export function usePortsvcQuery(query: PortsvcQuery) {
  return useQuery({
    queryKey: portsvcKeys.snapshot(query),
    queryFn: () => loadPortsvc(query),
  });
}
