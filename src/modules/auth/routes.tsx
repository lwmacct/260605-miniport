import type { RouteObject } from "react-router-dom";
import { GuestOnlyBoundary } from "@/app/router/guards";
import { LoginRoute } from "./LoginRoute";
import { RegisterRoute } from "./RegisterRoute";

export const authRoutes: RouteObject = {
	element: <GuestOnlyBoundary />,
	children: [
		{ path: "login", element: <LoginRoute /> },
		{ path: "register", element: <RegisterRoute /> },
	],
};
