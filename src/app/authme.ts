export const authmePath = "/authme";
export const authRefreshEvent = "miniport:auth-refresh";

export function authmeEndpoint(path: string) {
  return `${authmePath}${path}`;
}

export function dispatchAuthRefresh() {
  window.dispatchEvent(new Event(authRefreshEvent));
}
