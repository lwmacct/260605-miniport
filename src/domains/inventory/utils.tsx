import { Tag } from "antd";
import { statusOptions } from "./constants";
import type { PortSlot } from "./types";

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

export function buildStats(groups: PortGroupLike[], hosts: HostLike[]) {
  const slots = groups.flatMap((group) => group.slots ?? []);
  return {
    hosts: hosts.length,
    groups: groups.length,
    usedSlots: slots.filter((slot) => slot.status !== "empty").length,
    emptySlots: slots.filter((slot) => slot.status === "empty").length,
    components: groups.reduce((sum, group) => sum + (group.components ?? []).length, 0),
    repositories: groups.reduce((sum, group) => sum + (group.repositories ?? []).length, 0),
  };
}

type HostLike = { id: number };
type PortGroupLike = {
  components?: unknown[];
  repositories?: unknown[];
  slots?: Array<{ status: string }>;
};
