import { Navigate } from "react-router-dom";
import { appPaths } from "@/app/router/navigation";
import { useAuthStateQuery } from "./model/authQueries";
import { AuthScreen } from "./ui/AuthScreen";

export function RegisterRoute() {
	const authState = useAuthStateQuery();

	if (authState.data && !authState.data.config.local.registrationEnabled) {
		return <Navigate to={appPaths.login} replace />;
	}

	return <AuthScreen mode="register" />;
}
