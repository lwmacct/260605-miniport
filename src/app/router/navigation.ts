export const appPaths = {
  console: "/console/overview",
  consoleDependencies: "/console/dependencies",
  consoleHosts: "/console/hosts",
  consoleOverview: "/console/overview",
  consolePortGroups: "/console/port-groups",
  consoleProjects: "/console/projects",
  consoleServiceGroups: "/console/service-groups",
  settings: "/settings/appearance",
  settingsGitHub: "/settings/github",
} as const;

export type TopNavKey = "console" | "settings";

export function topNavFromPathname(pathname: string): TopNavKey {
  if (pathname.startsWith("/settings")) {
    return "settings";
  }
  return "console";
}
