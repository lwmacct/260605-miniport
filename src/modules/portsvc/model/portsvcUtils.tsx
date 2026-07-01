import { Tag } from "antd";
import { statusOptions } from "./portsvcConstants";
import type { DependencyAssetItem, PortGroupItem, ServiceGroupItem } from "./portsvcTypes";

export function statusTag(status: string) {
  const colorMap: Record<string, string> = {
    active: "green",
    available: "default",
    planned: "blue",
    reserved: "gold",
    running: "green",
    stopped: "default",
  };
  const labelMap: Record<string, string> = {
    active: "运行中",
  };
  const label = labelMap[status] ?? (statusOptions.find((item) => item.value === status)?.label ?? status);
  return <Tag color={colorMap[status] ?? "default"}>{label}</Tag>;
}

export function runtimeTag(mode: string) {
  return <Tag color={mode === "host" ? "purple" : "cyan"}>{mode === "host" ? "宿主机" : "DIND"}</Tag>;
}

export function splitTags(tags: string) {
  return tags
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

export function portRange(group?: Pick<PortGroupItem, "portStart" | "portEnd"> | null) {
  return group ? `${group.portStart}-${group.portEnd}` : "-";
}

export function buildStats(groups: PortGroupItem[], hosts: unknown[], assets: DependencyAssetItem[], serviceGroups: ServiceGroupItem[] = []) {
  return {
    groups: groups.length,
    hosts: hosts.length,
    runningGroups: groups.filter((item) => item.status === "running").length,
    freeGroups: groups.filter((item) => item.status === "available").length,
    slots: groups.reduce((sum, item) => sum + item.slots.length, 0),
    dependencyAssets: assets.length,
    assetLinks: groups.reduce((sum, item) => sum + item.assetLinks.length, 0),
    serviceGroups: serviceGroups.length,
  };
}
