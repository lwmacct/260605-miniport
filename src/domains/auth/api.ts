import { responseErrorMessage } from "../../shared/api/client";

export interface AuthSession {
  authenticated: boolean;
  expiresAt?: string;
  user?: {
    id: number;
    username: string;
    displayName: string;
    status: string;
    admin: boolean;
  };
}

export interface AuthPublicConfig {
  local: {
    loginEnabled: boolean;
    registrationEnabled: boolean;
  };
}

export interface AuthState {
  config: AuthPublicConfig;
  session: AuthSession;
}

export const defaultAuthConfig: AuthPublicConfig = {
  local: {
    loginEnabled: true,
    registrationEnabled: true,
  },
};

export async function fetchAuthSession(): Promise<AuthSession> {
  const response = await fetch("/api/auth/me", {
    credentials: "same-origin",
  });

  if (!response.ok) {
    return { authenticated: false };
  }

  return (await response.json()) as AuthSession;
}

export async function fetchAuthPublicConfig(): Promise<AuthPublicConfig> {
  const response = await fetch("/api/auth/config", {
    credentials: "same-origin",
  });

  if (!response.ok) {
    return defaultAuthConfig;
  }

  return (await response.json()) as AuthPublicConfig;
}

export async function fetchAuthState(): Promise<AuthState> {
  const [session, config] = await Promise.all([
    fetchAuthSession(),
    fetchAuthPublicConfig(),
  ]);
  return { config, session };
}

export async function login(username: string, password: string): Promise<void> {
  const response = await fetch("/api/auth/password/login", {
    body: JSON.stringify({ username, password }),
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(await responseErrorMessage(response, "登录失败"));
  }
}

export async function register(username: string, password: string): Promise<void> {
  const response = await fetch("/api/auth/password/register", {
    body: JSON.stringify({ username, password }),
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(await responseErrorMessage(response, "注册失败"));
  }
}

export async function logout(): Promise<void> {
  await fetch("/api/auth/logout", {
    credentials: "same-origin",
    method: "POST",
  });
}
