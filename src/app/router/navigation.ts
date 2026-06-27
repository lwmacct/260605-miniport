export const appPaths = {
  admin: "/admin/users",
  dependencies: "/dependencies",
  projects: "/projects",
  login: "/login",
  overview: "/overview",
  register: "/register",
  services: "/services",
} as const;

export type TopNavKey = "admin" | "dependencies" | "projects" | "overview" | "services";

export function topNavFromPathname(pathname: string): TopNavKey {
  if (pathname.startsWith("/admin")) {
    return "admin";
  }
  if (pathname.startsWith("/services")) {
    return "services";
  }
  if (pathname.startsWith("/projects")) {
    return "projects";
  }
  if (pathname.startsWith("/dependencies")) {
    return "dependencies";
  }
  return "overview";
}
