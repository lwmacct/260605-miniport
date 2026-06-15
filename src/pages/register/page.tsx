import { Navigate } from "react-router-dom";
import { appPaths } from "../../app/router/navigation";
import { AuthScreen } from "../../domains/auth/components/AuthScreen";
import { useAuthStateQuery } from "../../domains/auth/queries";

export function RegisterPage() {
  const authState = useAuthStateQuery();

  if (authState.data && !authState.data.config.local.registrationEnabled) {
    return <Navigate to={appPaths.login} replace />;
  }

  return <AuthScreen mode="register" />;
}
