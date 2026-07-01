import { apiGet, apiSend } from "@/shared/api/client";
import type {
  DependencyAssetItem,
  HostForm,
  HostItem,
  Meta,
  PortGroupForm,
  PortGroupItem,
  PortsvcQuery,
  PortsvcSnapshot,
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
  const groupsPath = "/api/port-groups" + buildQueryString({
    q: query.query,
    sort: query.sort,
    status: query.status,
    ownerSubject: query.ownerSubject,
  });

  const [meta, hosts, portGroups, dependencyAssets] = await Promise.all([
    apiGet<Meta>("/api/meta"),
    apiGet<HostItem[]>("/api/hosts"),
    apiGet<PortGroupItem[]>(groupsPath),
    apiGet<DependencyAssetItem[]>("/api/dependency-assets"),
  ]);

  return {
    meta,
    hosts: hosts ?? [],
    dependencyAssets: dependencyAssets ?? [],
    portGroups: (portGroups ?? []).map((item) => ({
      ...item,
      assetLinks: item.assetLinks ?? [],
      slots: item.slots ?? [],
      tags: item.tags ?? "",
    })),
  };
}

export function savePortGroup(group: PortGroupForm, editingGroup?: PortGroupItem | null) {
  const payload = {
    assetLinks: group.assetLinks ?? [],
    hostId: group.hostId ?? "",
    notes: group.notes ?? "",
    ownerSubject: group.ownerSubject ?? "",
    portStart: group.portStart ?? 0,
    projectName: group.projectName ?? "",
    projectOwner: group.projectOwner ?? "",
    runtimeMode: group.runtimeMode ?? "dind",
    runtimeName: group.runtimeName ?? "",
    serviceIp: group.serviceIp ?? "",
    slots: group.slots ?? [],
    status: group.status ?? "available",
    tags: group.tags ?? "",
  };
  return apiSend<PortGroupItem>(
    editingGroup ? `/api/port-groups/${editingGroup.id}` : "/api/port-groups",
    {
      method: editingGroup ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removePortGroup(group: PortGroupItem) {
  return apiSend<{ deleted: boolean }>(`/api/port-groups/${group.id}`, { method: "DELETE" });
}

export function saveHost(host: HostForm, editingHost?: HostItem | null) {
  const payload = {
    ip: host.ip ?? "",
    name: host.name ?? "",
    notes: host.notes ?? "",
    spec: host.spec ?? "",
    status: host.status ?? "active",
  };
  return apiSend<HostItem>(editingHost ? `/api/hosts/${editingHost.id}` : "/api/hosts", {
    method: editingHost ? "PUT" : "POST",
    body: JSON.stringify(payload),
  });
}

export function removeHost(host: HostItem) {
  return apiSend<{ deleted: boolean }>(`/api/hosts/${host.id}`, { method: "DELETE" });
}

export function saveDependencyAsset(asset: Partial<DependencyAssetItem>, editingAsset?: DependencyAssetItem | null) {
  const payload = {
    assetKind: asset.assetKind ?? "component",
    assetType: asset.assetType ?? "middleware",
    controllability: asset.controllability ?? "unknown",
    description: asset.description ?? "",
    externalId: asset.externalId ?? "",
    fullName: asset.fullName ?? "",
    metadata: asset.metadata ?? "",
    name: asset.name ?? "",
    notes: asset.notes ?? "",
    ownerSubject: asset.ownerSubject ?? "",
    provider: asset.provider ?? "manual",
    status: asset.status ?? "active",
    url: asset.url ?? "",
    visibility: asset.visibility ?? "unknown",
  };
  return apiSend<DependencyAssetItem>(
    editingAsset ? `/api/dependency-assets/${editingAsset.id}` : "/api/dependency-assets",
    {
      method: editingAsset ? "PUT" : "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function removeDependencyAsset(asset: DependencyAssetItem) {
  return apiSend<{ deleted: boolean }>(`/api/dependency-assets/${asset.id}`, { method: "DELETE" });
}

export function exportPortGroupsURL(query: PortsvcQuery) {
  return "/api/port-groups/export.csv" + buildQueryString({
    q: query.query,
    sort: query.sort,
    status: query.status,
    ownerSubject: query.ownerSubject,
  });
}
