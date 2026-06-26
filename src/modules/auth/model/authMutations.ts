import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  defaultAuthConfig,
  fetchAuthState,
  login,
  logout,
  register,
  type AuthChallengeResponse,
} from "../api/authApi";
import { authKeys } from "./authQueries";

export function useLoginMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      challenge,
      password,
      username,
    }: {
      challenge: AuthChallengeResponse;
      password: string;
      username: string;
    }) => login(username, password, challenge),
    onSuccess: async () => {
      queryClient.setQueryData(authKeys.state, await fetchAuthState());
    },
  });
}

export function useRegisterMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      challenge,
      password,
      username,
    }: {
      challenge: AuthChallengeResponse;
      password: string;
      username: string;
    }) => register(username, password, challenge),
    onSuccess: async () => {
      queryClient.setQueryData(authKeys.state, await fetchAuthState());
    },
  });
}

export function useLogoutMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.setQueryData(authKeys.state, {
        config: defaultAuthConfig,
        session: { authenticated: false },
      });
    },
  });
}
