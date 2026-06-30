import { useQuery } from "@tanstack/react-query";
import { fetchAuthState } from "../api/authApi";

export const authKeys = {
  state: ["auth", "state"] as const,
};

export function useAuthStateQuery() {
  return useQuery({
    queryKey: authKeys.state,
    queryFn: fetchAuthState,
    staleTime: 30_000,
  });
}
