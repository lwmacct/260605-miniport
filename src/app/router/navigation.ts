export const appPaths = {
  admin: "/admin/users",
  console: "/console/overview",
  consoleDependencies: "/console/dependencies",
  consoleOverview: "/console/overview",
  consoleProjects: "/console/projects",
  consoleServices: "/console/services",
  login: "/login",
  register: "/register",
  settings: "/settings/appearance",
} as const;

export type TopNavKey = "admin" | "console" | "settings";

export function topNavFromPathname(pathname: string): TopNavKey {
  if (pathname.startsWith("/admin")) {
    return "admin";
  }
  if (pathname.startsWith("/settings")) {
    return "settings";
  }
  return "console";
}
