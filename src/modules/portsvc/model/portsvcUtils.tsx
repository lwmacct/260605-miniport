import { Tag } from "antd";
import { statusOptions } from "./portsvcConstants";
import type { PortGroupItem } from "./portsvcTypes";

export function statusTag(status: string) {
  const colorMap: Record<string, string> = {
    active: "green",
    available: "default",
    planned: "blue",
    reserved: "gold",
    running: "green",
    stopped: "default",
  };
  const label = status === "available" ? "可用" : (statusOptions.find((item) => item.value === status)?.label ?? status);
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

export function buildStats(groups: PortGroupItem[], hosts: unknown[]) {
  return {
    groups: groups.length,
    hosts: hosts.length,
    runningGroups: groups.filter((item) => item.status === "running").length,
    freeGroups: groups.filter((item) => item.status === "available").length,
    slots: groups.reduce((sum, item) => sum + item.slots.length, 0),
    repositories: groups.reduce((sum, item) => sum + item.repositories.length, 0),
    dependencies: groups.reduce((sum, item) => sum + item.dependencies.length, 0),
  };
}
