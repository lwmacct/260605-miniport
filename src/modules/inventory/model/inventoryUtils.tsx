import { Tag } from "antd";
import { statusOptions } from "./inventoryConstants";
import type { PortSlot } from "./inventoryTypes";

export function statusTag(status: string) {
  const colorMap: Record<string, string> = {
    planned: "blue",
    reserved: "gold",
    running: "green",
    stopped: "default",
  };
  const label = statusOptions.find((item) => item.value === status)?.label ?? status;
  return <Tag color={colorMap[status] ?? "default"}>{label}</Tag>;
}

export function splitTags(tags: string) {
  return tags
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

export function makeSlots(portStart?: number | null, portEnd?: number | null, existing?: PortSlot[]) {
  if (!portStart || !portEnd || portEnd < portStart) {
    return [];
  }
  const byPort = new Map((existing ?? []).map((slot) => [slot.port, slot]));
  return Array.from({ length: portEnd - portStart + 1 }, (_, index) => {
    const port = portStart + index;
    return {
      port,
      name: byPort.get(port)?.name ?? "",
      protocol: byPort.get(port)?.protocol ?? "tcp",
      purpose: byPort.get(port)?.purpose ?? "",
      status: byPort.get(port)?.status ?? "empty",
      notes: byPort.get(port)?.notes ?? "",
    };
  });
}

export function buildStats(groups: PortGroupLike[]) {
  const slots = groups.flatMap((group) => group.slots ?? []);
  return {
    allocations: groups.length,
    users: new Set(groups.map((group) => group.username).filter(Boolean)).size,
    usedSlots: slots.filter((slot) => slot.status !== "empty").length,
    emptySlots: slots.filter((slot) => slot.status === "empty").length,
    projects: groups.reduce((sum, group) => sum + (group.projects ?? []).length, 0),
    components: groups.reduce((sum, group) => sum + (group.components ?? []).length, 0),
    repositories: groups.reduce((sum, group) => sum + (group.repositories ?? []).length, 0),
  };
}

type PortGroupLike = {
  components?: unknown[];
  projects?: unknown[];
  repositories?: unknown[];
  slots?: Array<{ status: string }>;
  username?: string;
};
