import { apiClient, apiData } from "@/shared/api/client";
import type { components, paths } from "@/shared/api/schema.gen";
import type {
  DependencyAssetItem,
  HostForm,
  HostItem,
  PortGroupForm,
  PortGroupItem,
  PortsvcQuery,
  PortsvcSnapshot,
  ServiceGroupForm,
  ServiceGroupItem,
} from "../model/portsvcTypes";

type Schema = components["schemas"];
type APIPath = keyof paths;

export async function loadPortsvc(query: PortsvcQuery): Promise<PortsvcSnapshot> {
  const serviceGroupFilters = { q: query.q, status: query.status };
  const [meta, hosts, portGroups, dependencyAssets, serviceGroups, repositories] = await Promise.all([
    apiData(apiClient.GET("/meta")),
    apiData(apiClient.GET("/console/hosts", { params: { query: {} } })),
    apiData(apiClient.GET("/console/port-groups", { params: { query } })),
    apiData(apiClient.GET("/console/dependency-assets", { params: { query: {} } })),
    apiData(apiClient.GET("/console/service-groups", { params: { query: serviceGroupFilters } })),
    apiData(apiClient.GET("/console/github/repositories", { params: { query: {} } })),
  ]);
  return {
    meta,
    hosts: hosts.items,
    dependencyAssets: dependencyAssets.items,
    repositories: repositories.items,
    portGroups: portGroups.items,
    serviceGroups: serviceGroups.items,
  };
}

export function savePortGroup(group: PortGroupForm, editingGroup?: PortGroupItem | null) {
  const payload = {
    assetLinks: group.assetLinks ?? [],
    repositoryLinks: group.repositoryLinks ?? [],
    hostId: group.hostId ?? "",
    notes: group.notes ?? "",
    portPrefix: group.portPrefix ?? 0,
    environmentName: group.environmentName ?? "",
    environmentOwner: group.environmentOwner ?? "",
    runtimeMode: group.runtimeMode ?? "dind",
    runtimeName: group.runtimeName ?? "",
    serviceIp: group.serviceIp ?? "",
    slots: group.slots ?? [],
    status: group.status ?? "available",
    tags: group.tags ?? "",
  } satisfies Schema["PortGroupCreateDTO"];
  if (editingGroup) {
    const update = { ...payload, id: editingGroup.id } satisfies Schema["PortGroupUpdateDTO"];
    return apiData(apiClient.PUT("/console/port-groups", { body: { items: [update] } }));
  }
  return apiData(apiClient.POST("/console/port-groups", { body: { items: [payload] } }));
}

export function removePortGroup(group: PortGroupItem) {
  return apiData(apiClient.DELETE("/console/port-groups", { body: { ids: [group.id] } }));
}

export function saveServiceGroup(group: ServiceGroupForm, editingGroup?: ServiceGroupItem | null) {
  const payload = {
    description: group.description ?? "",
    kind: group.kind ?? "service",
    name: group.name ?? "",
    notes: group.notes ?? "",
    portGroups: group.portGroups ?? [],
    status: group.status ?? "active",
  } satisfies Schema["ServiceGroupCreateDTO"];
  if (editingGroup) {
    const update = { ...payload, id: editingGroup.id } satisfies Schema["ServiceGroupUpdateDTO"];
    return apiData(apiClient.PUT("/console/service-groups", { body: { items: [update] } }));
  }
  return apiData(apiClient.POST("/console/service-groups", { body: { items: [payload] } }));
}

export function removeServiceGroup(group: ServiceGroupItem) {
  return apiData(apiClient.DELETE("/console/service-groups", { body: { ids: [group.id] } }));
}

export function saveHost(host: HostForm, editingHost?: HostItem | null) {
  const payload = {
    ip: host.ip ?? "",
    name: host.name ?? "",
    notes: host.notes ?? "",
    spec: host.spec ?? "",
    status: host.status ?? "active",
  } satisfies Schema["HostCreateDTO"];
  if (editingHost) {
    const update = { ...payload, id: editingHost.id } satisfies Schema["HostUpdateDTO"];
    return apiData(apiClient.PUT("/console/hosts", { body: { items: [update] } }));
  }
  return apiData(apiClient.POST("/console/hosts", { body: { items: [payload] } }));
}

export function removeHost(host: HostItem) {
  return apiData(apiClient.DELETE("/console/hosts", { body: { ids: [host.id] } }));
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
    provider: asset.provider ?? "manual",
    status: asset.status ?? "active",
    url: asset.url ?? "",
    visibility: asset.visibility ?? "unknown",
  } satisfies Schema["DependencyAssetCreateDTO"];
  if (editingAsset) {
    const update = { ...payload, id: editingAsset.id } satisfies Schema["DependencyAssetUpdateDTO"];
    return apiData(apiClient.PUT("/console/dependency-assets", { body: { items: [update] } }));
  }
  return apiData(apiClient.POST("/console/dependency-assets", { body: { items: [payload] } }));
}

export function removeDependencyAsset(asset: DependencyAssetItem) {
  return apiData(apiClient.DELETE("/console/dependency-assets", { body: { ids: [asset.id] } }));
}

export function exportPortGroupsURL(query: PortsvcQuery) {
  const path: APIPath = "/console/port-groups/export.csv";
  const search = new URLSearchParams();
  if (query.q) search.set("q", query.q);
  if (query.sort) search.set("sort", query.sort);
  if (query.status) search.set("status", query.status);
  const suffix = search.toString();
  return `/api${path}${suffix ? `?${suffix}` : ""}`;
}
