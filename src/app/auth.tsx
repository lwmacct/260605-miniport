import { GithubOutlined, KeyOutlined } from "@ant-design/icons";
import {
  WorkbenchAccessDeniedPage,
  WorkbenchOAuthSignInPage,
  WorkbenchTokenSignInPage,
} from "@lwmacct/260627-antd-workbench";
import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import { APP_NAME } from "@/shared/config/appConfig";
import { authmeEndpoint, authRefreshEvent } from "./authme";
import { loadSession, type AuthIdentity, type AuthMethod, type AuthMethodID, type SessionState } from "./session";

type AuthState =
  | SessionState
  | { status: "signing-in"; method: AuthMethodID; methods: AuthMethod[] }
  | { status: "invalid-token"; methods: AuthMethod[] };

type AuthContextValue = {
  identity: AuthIdentity;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error("useAuth must be used inside AuthBoundary");
  }
  return value;
}

export function AuthBoundary({ children, initialSession }: { children: ReactNode; initialSession: SessionState }) {
  const [state, setState] = useState<AuthState>(initialSession);
  const [logoutLoading, setLogoutLoading] = useState(false);

  const retrySession = useCallback(async () => {
    setState(await loadSession());
  }, []);

  useEffect(() => {
    const update = () => void retrySession();
    window.addEventListener(authRefreshEvent, update);
    return () => window.removeEventListener(authRefreshEvent, update);
  }, [retrySession]);

  const logout = useCallback(async () => {
    const methods = "methods" in state ? state.methods : undefined;
    if (!methods) {
      return;
    }
    setLogoutLoading(true);
    try {
      const response = await fetch(authmeEndpoint("/session"), {
        credentials: "same-origin",
        method: "DELETE",
      });
      setState(response.ok ? { status: "signed-out", methods } : { status: "unavailable", methods });
    } catch {
      setState({ status: "unavailable", methods });
    } finally {
      setLogoutLoading(false);
    }
  }, [state]);

  const oidcLogin = useCallback((methods: AuthMethod[]) => {
    setState({ status: "signing-in", method: "github", methods });
    const returnTo = window.location.pathname + window.location.search + window.location.hash;
    window.requestAnimationFrame(() => {
      window.location.assign(`${authmeEndpoint("/login/github")}?return_to=${encodeURIComponent(returnTo)}`);
    });
  }, []);

  const tokenLogin = useCallback(async (token: string, methods: AuthMethod[]) => {
    setState({ status: "signing-in", method: "token", methods });
    try {
      const response = await fetch(authmeEndpoint("/login/token"), {
        body: JSON.stringify({ token }),
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        method: "POST",
      });
      if (!response.ok) {
        setState(response.status === 401 ? { status: "invalid-token", methods } : { status: "unavailable", methods });
        return;
      }
      await retrySession();
    } catch {
      setState({ status: "unavailable", methods });
    }
  }, [retrySession]);

  const value = useMemo(
    () => state.status === "authenticated" && state.access === "granted" ? { identity: state.identity, logout } : null,
    [logout, state],
  );

  if (state.status === "authenticated" && state.access === "denied") {
    return (
      <WorkbenchAccessDeniedPage
        brand={{ mark: "M", name: APP_NAME }}
        identity={{
          avatarUrl: state.identity.avatar_url,
          displayName: state.identity.name,
          provider: state.identity.provider === "github" ? "GitHub" : "Access token",
          providerIcon: state.identity.provider === "github" ? <GithubOutlined /> : <KeyOutlined />,
          username: state.identity.username,
        }}
        logoutLoading={logoutLoading}
        onLogout={() => void logout()}
      />
    );
  }

  if (!value) {
    const methods = "methods" in state ? state.methods : undefined;
    const tokenEnabled = methods?.some((method) => method.id === "token") ?? false;
    const oidcEnabled = methods?.some((method) => method.id === "github") ?? false;
    if (methods && tokenEnabled) {
      return (
        <WorkbenchTokenSignInPage
          brand={{ description: "使用访问凭据进入控制台", mark: "M", name: APP_NAME }}
          error={state.status === "unavailable" ? "认证服务暂时不可用" : state.status === "invalid-token" ? "访问令牌无效" : undefined}
          loading={state.status === "signing-in" && state.method === "token"}
          oauth={oidcEnabled ? {
            pendingProvider: state.status === "signing-in" && state.method === "github" ? "github" : undefined,
            providers: [{ label: "GitHub", provider: "github" }],
            onSelectProvider: () => oidcLogin(methods),
          } : undefined}
          retry={state.status === "unavailable"}
          onRetry={state.status === "unavailable" ? retrySession : undefined}
          onSubmit={({ token }) => tokenLogin(token, methods)}
        />
      );
    }
    return (
      <WorkbenchOAuthSignInPage
        brand={{ description: "使用授权身份进入控制台", mark: "M", name: APP_NAME }}
        error={state.status === "unavailable" ? "认证服务暂时不可用" : undefined}
        hint={state.status === "signed-out" ? "仅允许已配置的操作员访问" : undefined}
        pendingProvider={state.status === "signing-in" && state.method === "github" ? "github" : undefined}
        providers={[{ disabled: !oidcEnabled, label: "GitHub", provider: "github" }]}
        retry={state.status === "unavailable"}
        onRetry={state.status === "unavailable" ? retrySession : undefined}
        onSelectProvider={() => methods && oidcLogin(methods)}
      />
    );
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
