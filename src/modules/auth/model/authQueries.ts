import { useQuery } from "@tanstack/react-query";
import { createImageChallenge, fetchAuthState } from "../api/authApi";

export const authKeys = {
  state: ["auth", "state"] as const,
  imageChallenge: (resetKey: number) => ["auth", "imageChallenge", resetKey] as const,
};

export function useAuthStateQuery() {
  return useQuery({
    queryKey: authKeys.state,
    queryFn: fetchAuthState,
    staleTime: 30_000,
  });
}

export function useImageChallengeQuery(resetKey: number, enabled: boolean) {
  return useQuery({
    enabled,
    queryKey: authKeys.imageChallenge(resetKey),
    queryFn: createImageChallenge,
    staleTime: 30_000,
    retry: false,
  });
}
