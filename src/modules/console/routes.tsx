import type { RouteObject } from "react-router-dom";
import { portsvcConsoleRoutes } from "@/modules/portsvc/routes";
import { ConsoleLayout } from "./layout/ConsoleLayout";

export const consoleRoutes: RouteObject = {
  path: "console",
  element: <ConsoleLayout />,
  children: portsvcConsoleRoutes,
};
