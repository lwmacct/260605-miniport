import createClient from "openapi-fetch";
import { dispatchAuthRefresh } from "@/app/authme";
import type { components, paths } from "./schema.gen";

export type APIProblemDetail = components["schemas"]["ErrorDetail"];
type GeneratedAPIProblem = components["schemas"]["ErrorModel"];

export type APIProblem = Omit<GeneratedAPIProblem, "detail" | "errors" | "status" | "title"> & {
  detail: string;
  errors?: APIProblemDetail[] | null;
  status: number;
  title: string;
};

const problemMessages: Record<string, string> = {
  "urn:problem:github-disabled": "GitHub 集成未启用",
  "urn:problem:github-resource-not-found": "GitHub 资源不存在",
  "urn:problem:internal-server-error": "服务暂时不可用",
  "urn:problem:invalid-github-signature": "GitHub Webhook 签名无效",
  "urn:problem:invalid-github-state": "GitHub 授权状态无效",
  "urn:problem:invalid-portsvc-request": "请求参数无效",
  "urn:problem:portsvc-resource-not-found": "资源不存在",
  "urn:problem:portsvc-resource-conflict": "资源已存在或正在使用",
};

export class APIError extends Error {
  readonly type: string;
  readonly status: number;
  readonly problem: APIProblem;

  constructor(problem: APIProblem) {
    super(problemMessages[problem.type] ?? problem.detail);
    this.name = "APIError";
    this.type = problem.type;
    this.status = problem.status;
    this.problem = problem;
  }
}

export const apiClient = createClient<paths>({
  baseUrl: "/api",
  credentials: "same-origin",
  headers: { Accept: "application/json" },
});

interface APIResponse<T> {
  data?: T;
  error?: unknown;
  response: Response;
}

export async function apiData<T>(request: Promise<APIResponse<T>>): Promise<T> {
  const result = await request;
  if (result.data !== undefined) {
    return result.data;
  }
  if (result.response.status === 401 || result.response.status === 403) {
    dispatchAuthRefresh();
  }
  throw new APIError(normalizeProblem(result.error, result.response.status));
}

function normalizeProblem(value: unknown, status: number): APIProblem {
  const fallback: APIProblem = {
    type: "urn:problem:invalid-error-response",
    title: "Request Failed",
    status,
    detail: `请求失败，状态码：${status}`,
  };
  if (!value || typeof value !== "object") {
    return fallback;
  }
  const problem = value as Partial<GeneratedAPIProblem>;
  if (typeof problem.title !== "string" || typeof problem.detail !== "string") {
    return fallback;
  }
  return {
    ...problem,
    type: typeof problem.type === "string" ? problem.type : "about:blank",
    status,
    title: problem.title,
    detail: problem.detail,
  };
}
