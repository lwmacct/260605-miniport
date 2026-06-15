import { useQuery } from "@tanstack/react-query";
import { fetchAuthState } from "./api";

export const authKeys = {
  state: ["auth", "state"] as const,
};

export function useAuthStateQuery() {
  return useQuery({
    queryKey: authKeys.state,
    queryFn: fetchAuthState,
  });
}
