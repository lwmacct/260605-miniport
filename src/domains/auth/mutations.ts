import { useMutation, useQueryClient } from "@tanstack/react-query";
import { defaultAuthConfig, fetchAuthState, login, logout, register } from "./api";
import { authKeys } from "./queries";

export function useLoginMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ username, password }: { username: string; password: string }) =>
      login(username, password),
    onSuccess: async () => {
      queryClient.setQueryData(authKeys.state, await fetchAuthState());
    },
  });
}

export function useRegisterMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ username, password }: { username: string; password: string }) =>
      register(username, password),
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
