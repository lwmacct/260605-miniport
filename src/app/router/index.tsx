import { Navigate, Outlet, createHashRouter } from "react-router-dom";
import { AppShell } from "../shell/AppShell";
import { appPaths } from "./navigation";
import { OverviewPage } from "../../pages/overview/page";
import { ServicesPage } from "../../pages/services/page";
import { HostsPage } from "../../pages/hosts/page";
import { DependenciesPage } from "../../pages/dependencies/page";

export const router = createHashRouter([
  {
    path: "/",
    element: <Outlet />,
    children: [
      {
        element: <AppShell />,
        children: [
          { index: true, element: <Navigate to={appPaths.overview} replace /> },
          { path: "overview", element: <OverviewPage /> },
          { path: "services", element: <ServicesPage /> },
          { path: "hosts", element: <HostsPage /> },
          { path: "dependencies", element: <DependenciesPage /> },
        ],
      },
      {
        path: "*",
        element: <Navigate to={appPaths.overview} replace />,
      },
    ],
  },
]);
