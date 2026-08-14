import { apiClient, apiData } from "@/shared/api/client";
import type { components, operations } from "@/shared/api/schema.gen";

type Schema = components["schemas"];

export type GitHubStatus = Schema["GithubStatusDTO"];
export type GitHubInstallation = Schema["GithubInstallationDTO"];
export type GitHubRepository = Schema["GithubRepositoryDTO"];
export type GitHubRepositoryFilters = NonNullable<operations["console-list-github-repositories"]["parameters"]["query"]>;

export async function loadGitHubStatus() {
  return apiData(apiClient.GET("/console/github/status"));
}

export async function loadGitHubInstallations() {
  const result = await apiData(apiClient.GET("/console/github/installations"));
  return result.items;
}

export async function loadGitHubRepositories(query = "") {
  const filters: GitHubRepositoryFilters = query.trim() ? { q: query.trim() } : {};
  const result = await apiData(apiClient.GET("/console/github/repositories", { params: { query: filters } }));
  return result.items;
}

export function beginGitHubConnection() {
  return apiData(apiClient.POST("/console/github/connections"));
}

export function syncGitHubInstallation(id: string) {
  return apiData(apiClient.POST("/console/github/installations/sync", { body: { ids: [id] } }));
}
