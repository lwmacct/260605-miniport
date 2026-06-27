import { Navigate, Outlet, createHashRouter } from "react-router-dom";
import { AppShell } from "../shell/AppShell";
import { appPaths } from "./navigation";
import { ProtectedBoundary } from "./guards";
import { adminRoutes } from "@/modules/admin/routes";
import { authRoutes } from "@/modules/auth/routes";
import { portsvcRoutes } from "@/modules/portsvc/routes";

export const router = createHashRouter([
  {
    path: "/",
    element: <Outlet />,
    children: [
      authRoutes,
      {
        element: <ProtectedBoundary />,
        children: [
          {
            element: <AppShell />,
            children: [
              { index: true, element: <Navigate to={appPaths.overview} replace /> },
              ...portsvcRoutes,
              adminRoutes,
            ],
          },
        ],
      },
      {
        path: "*",
        element: <Navigate to={appPaths.login} replace />,
      },
    ],
  },
]);
