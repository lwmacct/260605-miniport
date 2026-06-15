export const appPaths = {
  admin: "/admin/users",
  dependencies: "/dependencies",
  hosts: "/hosts",
  login: "/login",
  overview: "/overview",
  register: "/register",
  services: "/services",
} as const;

export type TopNavKey = "admin" | "dependencies" | "hosts" | "overview" | "services";

export function topNavFromPathname(pathname: string): TopNavKey {
  if (pathname.startsWith("/admin")) {
    return "admin";
  }
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
