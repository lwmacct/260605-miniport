import { Tag } from "antd";
import { statusOptions } from "./portsvcConstants";
import type { PortAllocation, ServiceItem } from "./portsvcTypes";

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

export function portRange(port?: PortAllocation | null) {
  return port ? `${port.portStart}-${port.portEnd}` : "-";
}

export function buildStats(services: ServiceItem[], ports: PortAllocation[]) {
  const bound = new Set(services.map((item) => item.portAllocationId).filter(Boolean));
  return {
    services: services.length,
    ports: ports.length,
    boundPorts: bound.size,
    freePorts: ports.filter((port) => !bound.has(port.id)).length,
    repositories: services.reduce((sum, item) => sum + item.repositories.length, 0),
    dependencies: services.reduce((sum, item) => sum + item.dependencies.length, 0),
  };
}
