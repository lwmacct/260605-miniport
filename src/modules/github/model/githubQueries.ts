import { useQuery } from "@tanstack/react-query";
import { loadGitHubInstallations, loadGitHubRepositories, loadGitHubStatus } from "../api/githubApi";

export const githubKeys = {
  all: ["github"] as const,
  status: ["github", "status"] as const,
  installations: ["github", "installations"] as const,
  repositories: (query = "") => ["github", "repositories", query] as const,
};

export function useGitHubStatusQuery() {
  return useQuery({ queryKey: githubKeys.status, queryFn: loadGitHubStatus });
}

export function useGitHubInstallationsQuery() {
  return useQuery({ queryKey: githubKeys.installations, queryFn: loadGitHubInstallations });
}

export function useGitHubRepositoriesQuery(query = "") {
  return useQuery({ queryKey: githubKeys.repositories(query), queryFn: () => loadGitHubRepositories(query) });
}
