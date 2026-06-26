import { Navigate, type RouteObject } from "react-router-dom";
import { AdminBoundary } from "@/app/router/guards";
import { AdminLayout } from "./layout/AdminLayout";
import { AuditRoute } from "./audit/AuditRoute";
import { IntegrationsRoute } from "./integrations/IntegrationsRoute";
import { SystemRoute } from "./system/SystemRoute";
import { AdminUsersPage } from "./users/ui/AdminUsersPage";

export const adminRoutes: RouteObject = {
	path: "admin",
	element: <AdminBoundary />,
	children: [
		{
			element: <AdminLayout />,
			children: [
				{ index: true, element: <Navigate to="users" replace /> },
				{ path: "users", element: <AdminUsersPage /> },
				{ path: "audit", element: <AuditRoute /> },
				{ path: "system", element: <SystemRoute /> },
				{ path: "integrations", element: <IntegrationsRoute /> },
			],
		},
	],
};
