import { Navigate, type RouteObject } from "react-router-dom";
import { SettingsLayout } from "./SettingsLayout";
import { SettingsRoute } from "./SettingsRoute";

export const settingsRoutes: RouteObject = {
  path: "settings",
  element: <SettingsLayout />,
  children: [
    { index: true, element: <Navigate to="appearance" replace /> },
    { path: "appearance", element: <SettingsRoute /> },
  ],
};
