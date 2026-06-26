import type { RouteObject } from "react-router-dom";
import { DependenciesRoute } from "./DependenciesRoute";
import { HostsRoute } from "./HostsRoute";
import { OverviewRoute } from "./OverviewRoute";
import { ServicesRoute } from "./ServicesRoute";

export const inventoryRoutes: RouteObject[] = [
	{ path: "overview", element: <OverviewRoute /> },
	{ path: "services", element: <ServicesRoute /> },
	{ path: "hosts", element: <HostsRoute /> },
	{ path: "dependencies", element: <DependenciesRoute /> },
];
