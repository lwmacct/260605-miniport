export class APIError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message);
    this.name = "APIError";
  }
}

async function responseErrorMessage(response: Response, fallback: string): Promise<string> {
  try {
    const payload = (await response.json()) as { detail?: unknown; title?: unknown };
    if (typeof payload.detail === "string" && payload.detail) {
      return payload.detail;
    }
    if (typeof payload.title === "string" && payload.title) {
      return payload.title;
    }
  } catch {
    // Keep fallback when the response is not JSON.
  }
  return fallback;
}

export async function apiGet<T>(path: string): Promise<T> {
  const response = await fetch(path, {
    credentials: "same-origin",
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new APIError(
      await responseErrorMessage(response, `Request failed with status ${response.status}`),
      response.status,
    );
  }

  return response.json() as Promise<T>;
}

export async function apiSend<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    credentials: "same-origin",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    ...init,
  });

  if (!response.ok) {
    throw new APIError(
      await responseErrorMessage(response, `Request failed with status ${response.status}`),
      response.status,
    );
  }

  return response.json() as Promise<T>;
}
