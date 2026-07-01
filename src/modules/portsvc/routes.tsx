import { Navigate, type RouteObject } from "react-router-dom";
import { DependenciesRoute } from "./DependenciesRoute";
import { HostsRoute } from "./HostsRoute";
import { ProjectServicesRoute } from "./ProjectServicesRoute";
import { OverviewRoute } from "./OverviewRoute";
import { PortGroupsRoute } from "./PortGroupsRoute";
import { ServiceGroupsRoute } from "./ServiceGroupsRoute";

export const portsvcConsoleRoutes: RouteObject[] = [
  { index: true, element: <Navigate to="overview" replace /> },
  { path: "overview", element: <OverviewRoute /> },
  { path: "hosts", element: <HostsRoute /> },
  { path: "port-groups", element: <PortGroupsRoute /> },
  { path: "projects", element: <ProjectServicesRoute /> },
  { path: "service-groups", element: <ServiceGroupsRoute /> },
  { path: "dependencies", element: <DependenciesRoute /> },
  { path: "services", element: <Navigate to="../projects" replace /> },
];
