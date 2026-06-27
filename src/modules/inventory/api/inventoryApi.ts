import { apiGet, apiSend } from "@/shared/api/client";
import type {
  BatchPortGroupUpdate,
  GroupForm,
  InventoryQuery,
  InventorySnapshot,
  Meta,
  PortGroup,
} from "../model/inventoryTypes";

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

export async function loadInventory(query: InventoryQuery): Promise<InventorySnapshot> {
  const groupsPath = "/api/inventory/allocations" + buildQueryString({
    dindIp: query.dindIp,
    projectName: query.projectName,
    q: query.portGroupQuery,
    sort: query.portGroupSort,
    status: query.status,
    userId: query.userId,
  });

  const [meta, groups] = await Promise.all([
    apiGet<Meta>("/api/meta"),
    apiGet<PortGroup[]>(groupsPath),
  ]);

  return {
    meta,
    groups: (groups ?? []).map((group) => ({
      ...group,
      components: group.components ?? [],
      projects: group.projects ?? [],
      repositories: group.repositories ?? [],
      slots: group.slots ?? [],
      tags: group.tags ?? "",
    })),
  };
}

export function savePortGroup(group: GroupForm, editingGroup?: PortGroup | null) {
  const payload = {
    components: group.components ?? [],
    dindContainer: group.dindContainer ?? "",
    dindIp: group.dindIp ?? "",
    name: group.name ?? "",
    notes: group.notes ?? "",
    owner: group.owner ?? "",
    portStart: group.portStart ?? 0,
    projects: group.projects ?? [],
    repositories: group.repositories ?? [],
    slots: group.slots ?? [],
    status: group.status ?? "planned",
    tags: group.tags ?? "",
    userId: group.userId ?? 0,
  };
  return apiSend<PortGroup>(
    editingGroup ? `/api/inventory/allocations/${editingGroup.id}` : "/api/inventory/allocations",
    {
      method: editingGroup ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removePortGroup(group: PortGroup) {
  return apiSend<{ deleted: boolean }>(`/api/inventory/allocations/${group.id}`, { method: "DELETE" });
}

export function exportPortGroupsURL(query: InventoryQuery) {
  return "/api/inventory/exports/allocations.csv" + buildQueryString({
    dindIp: query.dindIp,
    projectName: query.projectName,
    q: query.portGroupQuery,
    sort: query.portGroupSort,
    status: query.status,
    userId: query.userId,
  });
}

export function batchUpdatePortGroups(ids: number[], changes: BatchPortGroupUpdate) {
  return apiSend<PortGroup[]>("/api/inventory/allocations/batch-update", {
    method: "POST",
    body: JSON.stringify({
      ids,
      ...changes,
    }),
  });
}

export function batchDeletePortGroups(ids: number[]) {
  return apiSend<{ deleted: boolean }>("/api/inventory/allocations/batch-delete", {
    method: "POST",
    body: JSON.stringify({ ids }),
  });
}
