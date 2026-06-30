import {
  WorkbenchCredentialPage,
  type WorkbenchChallengeResponse,
  type WorkbenchCredentialConfig,
  type WorkbenchCredentialSubmitValues,
} from "@lwmacct/260627-antd-workbench";
import { useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { appPaths } from "@/app/router/navigation";
import { useLoginMutation, useRegisterMutation } from "../model/authMutations";
import { useAuthStateQuery } from "../model/authQueries";
import { createImageChallenge, type AuthChallengeResponse } from "../api/authApi";

interface AuthScreenProps {
  mode: "login" | "register";
}

export function AuthScreen({ mode }: AuthScreenProps) {
  const navigate = useNavigate();
  const authState = useAuthStateQuery();
  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const isRegister = mode === "register";
  const activeMutation = isRegister ? registerMutation : loginMutation;
  const config = authState.data?.config;
  const resetLoginMutation = loginMutation.reset;
  const resetRegisterMutation = registerMutation.reset;

  const credentialConfig: WorkbenchCredentialConfig = {
    challenge: config?.challenge ?? { provider: "image" },
    local: {
      loginEnabled: config?.local.loginEnabled ?? true,
      registrationEnabled: config?.local.registrationEnabled ?? true,
    },
    oauth: false,
  };

  useEffect(() => {
    resetLoginMutation();
    resetRegisterMutation();
  }, [mode, resetLoginMutation, resetRegisterMutation]);

  async function submit(values: WorkbenchCredentialSubmitValues) {
    const challenge = toAuthChallenge(values.challenge);
    if (!challenge || activeMutation.isPending) {
      return;
    }

    await activeMutation.mutateAsync({
      challenge,
      password: values.password,
      username: values.username,
    });
    navigate(appPaths.overview, { replace: true });
  }

  return (
    <WorkbenchCredentialPage
      config={credentialConfig}
      createImageChallenge={createImageChallenge}
      error={activeMutation.error?.message}
      labels={{
        loginDescription: "使用账号进入资产控制台",
        loginTitle: "登录 Miniport",
        registerDescription: "创建账号后进入资产控制台",
        registerTitle: "注册 Miniport",
      }}
      loading={activeMutation.isPending}
      mode={mode}
      renderModeSwitch={({ children, targetMode }) => (
        <Link to={targetMode === "login" ? appPaths.login : appPaths.register}>
          {children}
        </Link>
      )}
      onModeChange={(targetMode) => navigate(targetMode === "login" ? appPaths.login : appPaths.register)}
      onSubmit={submit}
    />
  );
}

function toAuthChallenge(challenge?: WorkbenchChallengeResponse): AuthChallengeResponse | undefined {
  if (!challenge) {
    return undefined;
  }

  switch (challenge.provider) {
    case "image":
    case "hcaptcha":
    case "turnstile":
      return challenge;
    case "recaptcha":
      return undefined;
  }
}
