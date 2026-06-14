export const appPaths = {
  dependencies: "/dependencies",
  hosts: "/hosts",
  overview: "/overview",
  services: "/services",
} as const;

export type TopNavKey = "dependencies" | "hosts" | "overview" | "services";

export function topNavFromPathname(pathname: string): TopNavKey {
  if (pathname.startsWith("/services")) {
    return "services";
  }
  if (pathname.startsWith("/hosts")) {
    return "hosts";
  }
  if (pathname.startsWith("/dependencies")) {
    return "dependencies";
  }
  return "overview";
}
