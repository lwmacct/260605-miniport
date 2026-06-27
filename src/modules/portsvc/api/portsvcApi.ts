import { apiGet, apiSend } from "@/shared/api/client";
import type {
  PortsvcQuery,
  PortsvcSnapshot,
  Meta,
  PortAllocation,
  PortAllocationForm,
  ServiceForm,
  ServiceItem,
} from "../model/portsvcTypes";

function buildQueryString(params: Record<string, string | number | undefined>) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === "") {
      continue;
    }
    search.set(key, String(value));
  }
  const text = search.toString();
  return text ? `?${text}` : "";
}

export async function loadPortsvc(query: PortsvcQuery): Promise<PortsvcSnapshot> {
  const servicesPath = "/api/services" + buildQueryString({
    projectName: query.projectName,
    q: query.serviceQuery,
    sort: query.serviceSort,
    status: query.status,
    userId: query.userId,
  });
  const portsPath = "/api/port-allocations";

  const [meta, services, ports] = await Promise.all([
    apiGet<Meta>("/api/meta"),
    apiGet<ServiceItem[]>(servicesPath),
    apiGet<PortAllocation[]>(portsPath),
  ]);

  return {
    meta,
    ports: ports ?? [],
    services: (services ?? []).map((item) => ({
      ...item,
      dependencies: item.dependencies ?? [],
      repositories: item.repositories ?? [],
      tags: item.tags ?? "",
    })),
  };
}

export function saveService(service: ServiceForm, editingService?: ServiceItem | null) {
  const payload = {
    dependencies: service.dependencies ?? [],
    dindContainer: service.dindContainer ?? "",
    dindIp: service.dindIp ?? "",
    name: service.name ?? "",
    notes: service.notes ?? "",
    owner: service.owner ?? "",
    portAllocationId: service.portAllocationId ?? 0,
    projectName: service.projectName ?? "",
    repositories: service.repositories ?? [],
    status: service.status ?? "planned",
    tags: service.tags ?? "",
    userId: service.userId ?? 0,
  };
  return apiSend<ServiceItem>(
    editingService ? `/api/services/${editingService.id}` : "/api/services",
    {
      method: editingService ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removeService(service: ServiceItem) {
  return apiSend<{ deleted: boolean }>(`/api/services/${service.id}`, { method: "DELETE" });
}

export function savePortAllocation(port: PortAllocationForm, editingPort?: PortAllocation | null) {
  const payload = {
    notes: port.notes ?? "",
    portStart: port.portStart ?? 0,
    status: port.status ?? "available",
    userId: port.userId ?? 0,
  };
  return apiSend<PortAllocation>(
    editingPort ? `/api/port-allocations/${editingPort.id}` : "/api/port-allocations",
    {
      method: editingPort ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removePortAllocation(port: PortAllocation) {
  return apiSend<{ deleted: boolean }>(`/api/port-allocations/${port.id}`, { method: "DELETE" });
}

export function exportServicesURL(query: PortsvcQuery) {
  return "/api/services/export.csv" + buildQueryString({
    projectName: query.projectName,
    q: query.serviceQuery,
    sort: query.serviceSort,
    status: query.status,
    userId: query.userId,
  });
}
