import { Navigate, type RouteObject } from "react-router-dom";
import { DependenciesRoute } from "./DependenciesRoute";
import { ProjectServicesRoute } from "./ProjectServicesRoute";
import { OverviewRoute } from "./OverviewRoute";
import { ServicesRoute } from "./ServicesRoute";

export const portsvcConsoleRoutes: RouteObject[] = [
	{ index: true, element: <Navigate to="overview" replace /> },
	{ path: "overview", element: <OverviewRoute /> },
	{ path: "services", element: <ServicesRoute /> },
	{ path: "projects", element: <ProjectServicesRoute /> },
	{ path: "dependencies", element: <DependenciesRoute /> },
];
