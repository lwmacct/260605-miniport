import { apiGet, apiSend } from "@/shared/api/client";

export type GitHubStatus = {
  enabled: boolean;
  appSlug: string;
  installUrl: string;
};

export type GitHubInstallation = {
  id: string;
  githubInstallationId: number;
  accountId: number;
  accountLogin: string;
  accountType: string;
  avatarUrl: string;
  repositorySelection: "all" | "selected" | string;
  status: string;
  suspendedAt?: string;
  lastSyncedAt?: string;
  lastSyncError?: string;
};

export type GitHubRepository = {
  id: string;
  installationId: string;
  githubRepositoryId: number;
  ownerLogin: string;
  name: string;
  fullName: string;
  htmlUrl: string;
  description: string;
  defaultBranch: string;
  visibility: string;
  private: boolean;
  fork: boolean;
  archived: boolean;
  disabled: boolean;
  state: string;
  pushedAt?: string;
  remoteUpdatedAt?: string;
  lastSeenAt: string;
};

export function loadGitHubStatus() {
  return apiGet<GitHubStatus>("/api/github/status");
}

export function loadGitHubInstallations() {
  return apiGet<GitHubInstallation[]>("/api/github/installations");
}

export function loadGitHubRepositories(query = "") {
  const search = new URLSearchParams();
  if (query.trim()) {
    search.set("q", query.trim());
  }
  const suffix = search.size ? `?${search.toString()}` : "";
  return apiGet<GitHubRepository[]>(`/api/github/repositories${suffix}`);
}

export function beginGitHubConnection() {
  return apiSend<{ url: string }>("/api/github/connections", { method: "POST" });
}

export function syncGitHubInstallation(id: string) {
  return apiSend<GitHubInstallation>(`/api/github/installations/${id}/sync`, { method: "POST" });
}
