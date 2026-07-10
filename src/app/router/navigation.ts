export const appPaths = {
  admin: "/admin/users",
  console: "/console/overview",
  consoleDependencies: "/console/dependencies",
  consoleHosts: "/console/hosts",
  consoleOverview: "/console/overview",
  consolePortGroups: "/console/port-groups",
  consoleProjects: "/console/projects",
  consoleServiceGroups: "/console/service-groups",
  login: "/login",
  register: "/register",
  settings: "/settings/appearance",
  settingsGitHub: "/settings/github",
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
