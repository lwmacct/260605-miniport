import type { RouteObject } from "react-router-dom";
import { DependenciesRoute } from "./DependenciesRoute";
import { ProjectServicesRoute } from "./ProjectServicesRoute";
import { OverviewRoute } from "./OverviewRoute";
import { ServicesRoute } from "./ServicesRoute";

export const portsvcRoutes: RouteObject[] = [
	{ path: "overview", element: <OverviewRoute /> },
	{ path: "services", element: <ServicesRoute /> },
	{ path: "projects", element: <ProjectServicesRoute /> },
	{ path: "dependencies", element: <DependenciesRoute /> },
];
