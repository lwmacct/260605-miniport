import { Navigate, createHashRouter } from "react-router-dom";
import { AppShell } from "../shell/AppShell";
import { appPaths } from "./navigation";
import { consoleRoutes } from "@/modules/console/routes";
import { settingsRoutes } from "@/modules/settings/routes";

export const router = createHashRouter([
  {
    path: "/",
    element: <AppShell />,
    children: [
      { index: true, element: <Navigate to={appPaths.console} replace /> },
      consoleRoutes,
      settingsRoutes,
      {
        path: "*",
        element: <Navigate to={appPaths.console} replace />,
      },
    ],
  },
]);
