import { Navigate, type RouteObject } from "react-router-dom";
import { AdminBoundary } from "@/app/router/guards";
import { AdminUsersPage } from "./users/ui/AdminUsersPage";

export const adminRoutes: RouteObject = {
	path: "admin",
	element: <AdminBoundary />,
	children: [
		{ index: true, element: <Navigate to="users" replace /> },
		{ path: "users", element: <AdminUsersPage /> },
	],
};
