import { Navigate, Outlet, createHashRouter } from "react-router-dom";
import { AppShell } from "../shell/AppShell";
import { appPaths } from "./navigation";
import { AdminBoundary, GuestOnlyBoundary, ProtectedBoundary } from "./guards";
import { OverviewPage } from "../../pages/overview/page";
import { ServicesPage } from "../../pages/services/page";
import { HostsPage } from "../../pages/hosts/page";
import { DependenciesPage } from "../../pages/dependencies/page";
import { LoginPage } from "../../pages/login/page";
import { RegisterPage } from "../../pages/register/page";
import { AdminUsersPage } from "../../pages/admin/users.page";

export const router = createHashRouter([
  {
    path: "/",
    element: <Outlet />,
    children: [
      {
        element: <GuestOnlyBoundary />,
        children: [
          { path: "login", element: <LoginPage /> },
          { path: "register", element: <RegisterPage /> },
        ],
      },
      {
        element: <ProtectedBoundary />,
        children: [
          {
            element: <AppShell />,
            children: [
              { index: true, element: <Navigate to={appPaths.overview} replace /> },
              { path: "overview", element: <OverviewPage /> },
              { path: "services", element: <ServicesPage /> },
              { path: "hosts", element: <HostsPage /> },
              { path: "dependencies", element: <DependenciesPage /> },
              {
                path: "admin",
                element: <AdminBoundary />,
                children: [
                  { index: true, element: <Navigate to="users" replace /> },
                  { path: "users", element: <AdminUsersPage /> },
                ],
              },
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
