import type { GroupForm, Host, HostForm, Meta, PortGroup } from "./types";

export async function requestJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    let detail = `${res.status} ${res.statusText}`;
    try {
      const payload = await res.json();
      detail = payload.detail || payload.title || detail;
    } catch {
      // Keep the HTTP status text when the backend returns no JSON body.
    }
    throw new Error(detail);
  }
  return (await res.json()) as T;
}

export async function loadInventory() {
  const [meta, hosts, groups] = await Promise.all([
    requestJSON<Meta>("/api/meta"),
    requestJSON<Host[]>("/api/hosts"),
    requestJSON<PortGroup[]>("/api/port-groups"),
  ]);
  const normalizedHosts = hosts ?? [];
  const normalizedGroups = groups ?? [];

  return {
    meta,
    hosts: normalizedHosts,
    groups: normalizedGroups.map((group) => ({
      ...group,
      tags: group.tags ?? "",
      slots: group.slots ?? [],
      components: group.components ?? [],
      repositories: group.repositories ?? [],
    })),
  };
}

export function saveHost(host: HostForm, editingHost?: Host | null) {
  const payload = {
    ip: host.ip ?? "",
    name: host.name ?? "",
    network: host.network ?? "",
    environment: host.environment ?? "",
    notes: host.notes ?? "",
  };
  return requestJSON<Host>(editingHost ? `/api/hosts/${editingHost.id}` : "/api/hosts", {
    method: editingHost ? "PUT" : "POST",
    body: JSON.stringify(payload),
  });
}

export function removeHost(host: Host) {
  return requestJSON<{ deleted: boolean }>(`/api/hosts/${host.id}`, { method: "DELETE" });
}

export function savePortGroup(group: GroupForm, editingGroup?: PortGroup | null) {
  const payload = {
    hostId: group.hostId ?? 0,
    portStart: group.portStart ?? 0,
    portEnd: group.portEnd ?? 0,
    serviceName: group.serviceName ?? "",
    containerName: group.containerName ?? "",
    dindHost: group.dindHost ?? "",
    status: group.status ?? "planned",
    owner: group.owner ?? "",
    tags: group.tags ?? "",
    notes: group.notes ?? "",
    slots: group.slots ?? [],
    components: group.components ?? [],
    repositories: group.repositories ?? [],
  };
  return requestJSON<PortGroup>(editingGroup ? `/api/port-groups/${editingGroup.id}` : "/api/port-groups", {
    method: editingGroup ? "PUT" : "POST",
    body: JSON.stringify(payload),
  });
}

export function removePortGroup(group: PortGroup) {
  return requestJSON<{ deleted: boolean }>(`/api/port-groups/${group.id}`, { method: "DELETE" });
}
