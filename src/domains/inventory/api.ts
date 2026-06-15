import { apiGet, apiSend } from "../../shared/api/client";
import type {
  BatchPortGroupUpdate,
  GroupForm,
  Host,
  HostForm,
  InventoryQuery,
  InventorySnapshot,
  Meta,
  PortGroup,
} from "./types";

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
  const hostsPath = "/api/hosts" + buildQueryString({
    environment: query.environment,
    q: query.hostQuery,
    sort: query.hostSort,
  });
  const groupsPath = "/api/port-groups" + buildQueryString({
    hostId: query.hostId,
    q: query.portGroupQuery,
    sort: query.portGroupSort,
    status: query.status,
  });

  const [meta, hosts, groups] = await Promise.all([
    apiGet<Meta>("/api/meta"),
    apiGet<Host[]>(hostsPath),
    apiGet<PortGroup[]>(groupsPath),
  ]);

  return {
    meta,
    hosts: hosts ?? [],
    groups: (groups ?? []).map((group) => ({
      ...group,
      components: group.components ?? [],
      repositories: group.repositories ?? [],
      slots: group.slots ?? [],
      tags: group.tags ?? "",
    })),
  };
}

export function saveHost(host: HostForm, editingHost?: Host | null) {
  const payload = {
    environment: host.environment ?? "",
    ip: host.ip ?? "",
    name: host.name ?? "",
    network: host.network ?? "",
    notes: host.notes ?? "",
  };
  return apiSend<Host>(editingHost ? `/api/hosts/${editingHost.id}` : "/api/hosts", {
    method: editingHost ? "PUT" : "POST",
    body: JSON.stringify(payload),
  });
}

export function removeHost(host: Host) {
  return apiSend<{ deleted: boolean }>(`/api/hosts/${host.id}`, { method: "DELETE" });
}

export function savePortGroup(group: GroupForm, editingGroup?: PortGroup | null) {
  const payload = {
    components: group.components ?? [],
    containerName: group.containerName ?? "",
    dindHost: group.dindHost ?? "",
    hostId: group.hostId ?? 0,
    notes: group.notes ?? "",
    owner: group.owner ?? "",
    portEnd: group.portEnd ?? 0,
    portStart: group.portStart ?? 0,
    repositories: group.repositories ?? [],
    serviceName: group.serviceName ?? "",
    slots: group.slots ?? [],
    status: group.status ?? "planned",
    tags: group.tags ?? "",
  };
  return apiSend<PortGroup>(
    editingGroup ? `/api/port-groups/${editingGroup.id}` : "/api/port-groups",
    {
      method: editingGroup ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removePortGroup(group: PortGroup) {
  return apiSend<{ deleted: boolean }>(`/api/port-groups/${group.id}`, { method: "DELETE" });
}

export function exportPortGroupsURL(query: InventoryQuery) {
  return "/api/exports/port-groups.csv" + buildQueryString({
    hostId: query.hostId,
    q: query.portGroupQuery,
    sort: query.portGroupSort,
    status: query.status,
  });
}

export function batchUpdatePortGroups(ids: number[], changes: BatchPortGroupUpdate) {
  return apiSend<PortGroup[]>("/api/port-groups/batch-update", {
    method: "POST",
    body: JSON.stringify({
      ids,
      ...changes,
    }),
  });
}

export function batchDeletePortGroups(ids: number[]) {
  return apiSend<{ deleted: boolean }>("/api/port-groups/batch-delete", {
    method: "POST",
    body: JSON.stringify({ ids }),
  });
}
