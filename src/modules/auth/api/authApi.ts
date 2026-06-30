import { responseErrorMessage } from "@/shared/api/client";

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

export type ChallengeProvider = "image" | "hcaptcha" | "turnstile";

export interface AuthChallengeConfig {
  provider: ChallengeProvider;
  sitekey?: string;
}

export interface AuthPublicConfig {
  local: {
    loginEnabled: boolean;
    registrationEnabled: boolean;
  };
  challenge: AuthChallengeConfig;
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
  challenge: { provider: "image" },
};

export interface ImageChallenge {
  provider: "image";
  challengeId: string;
  image: string;
  expiresAt: string;
}

export type AuthChallengeResponse =
  | { provider: "image"; challengeId: string; answer: string }
  | { provider: "hcaptcha"; token: string }
  | { provider: "turnstile"; token: string };

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
  const [config, session] = await Promise.all([fetchAuthPublicConfig(), fetchAuthSession()]);
  return { config, session };
}

export async function createImageChallenge(): Promise<ImageChallenge> {
  const response = await fetch("/api/auth/challenges", {
    credentials: "same-origin",
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(await responseErrorMessage(response, "验证码生成失败"));
  }

  return (await response.json()) as ImageChallenge;
}

export async function login(
  username: string,
  password: string,
  challenge: AuthChallengeResponse,
): Promise<void> {
  const response = await fetch("/api/auth/password/login", {
    body: JSON.stringify({ username, password, challenge }),
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(await responseErrorMessage(response, "登录失败"));
  }
}

export async function register(
  username: string,
  password: string,
  challenge: AuthChallengeResponse,
): Promise<void> {
  const response = await fetch("/api/auth/password/register", {
    body: JSON.stringify({ username, password, challenge }),
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
