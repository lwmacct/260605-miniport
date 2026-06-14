import { apiGet, apiSend } from "../../shared/api/client";
import type { GroupForm, Host, HostForm, InventorySnapshot, Meta, PortGroup } from "./types";

export async function loadInventory(): Promise<InventorySnapshot> {
  const [meta, hosts, groups] = await Promise.all([
    apiGet<Meta>("/api/meta"),
    apiGet<Host[]>("/api/hosts"),
    apiGet<PortGroup[]>("/api/port-groups"),
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
  return apiSend<{ body: Host }>(editingHost ? `/api/hosts/${editingHost.id}` : "/api/hosts", {
    method: editingHost ? "PUT" : "POST",
    body: JSON.stringify(payload),
  });
}

export function removeHost(host: Host) {
  return apiSend<{ body: { deleted: boolean } }>(`/api/hosts/${host.id}`, { method: "DELETE" });
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
  return apiSend<{ body: PortGroup }>(
    editingGroup ? `/api/port-groups/${editingGroup.id}` : "/api/port-groups",
    {
      method: editingGroup ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removePortGroup(group: PortGroup) {
  return apiSend<{ body: { deleted: boolean } }>(`/api/port-groups/${group.id}`, { method: "DELETE" });
}
