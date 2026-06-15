import { LockOutlined, UserOutlined } from "@ant-design/icons";
import { Alert, Button, Card, Form, Input, Space, Typography } from "antd";
import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { appPaths } from "../../../app/router/navigation";
import { useLoginMutation, useRegisterMutation } from "../mutations";
import { useAuthStateQuery } from "../queries";
import type { AuthChallengeResponse } from "../api";
import { ChallengeField } from "./ChallengeField";
import styles from "./AuthScreen.module.css";

interface AuthScreenProps {
  mode: "login" | "register";
}

interface AuthFormValues {
  confirmPassword?: string;
  password: string;
  username: string;
}

export function AuthScreen({ mode }: AuthScreenProps) {
  const [form] = Form.useForm<AuthFormValues>();
  const navigate = useNavigate();
  const authState = useAuthStateQuery();
  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const [challenge, setChallenge] = useState<AuthChallengeResponse>();
  const [challengeError, setChallengeError] = useState("");
  const [challengeResetKey, setChallengeResetKey] = useState(0);
  const isRegister = mode === "register";
  const activeMutation = isRegister ? registerMutation : loginMutation;
  const config = authState.data?.config;
  const localLoginEnabled = config?.local.loginEnabled ?? true;
  const registrationEnabled = config?.local.registrationEnabled ?? true;
  const visibleError = challengeError || activeMutation.error?.message || "";
  const resetLoginMutation = loginMutation.reset;
  const resetRegisterMutation = registerMutation.reset;

  const resetChallenge = useCallback(() => {
    setChallenge(undefined);
    setChallengeError("");
    setChallengeResetKey((value) => value + 1);
  }, []);

  useEffect(() => {
    form.resetFields();
    resetLoginMutation();
    resetRegisterMutation();
    resetChallenge();
  }, [form, mode, resetChallenge, resetLoginMutation, resetRegisterMutation]);

  async function submit(values: AuthFormValues) {
    if (!challenge || activeMutation.isPending) {
      return;
    }

    if (isRegister && values.password !== values.confirmPassword) {
      form.setFields([{ name: "confirmPassword", errors: ["两次输入的密码不一致"] }]);
      return;
    }

    try {
      await activeMutation.mutateAsync({
        challenge,
        password: values.password,
        username: values.username,
      });
      navigate(appPaths.overview, { replace: true });
    } catch {
      resetChallenge();
    } finally {
      form.setFieldValue("password", "");
      form.setFieldValue("confirmPassword", "");
    }
  }

  return (
    <main className={styles.page}>
      <Card className={styles.panel}>
        <Space orientation="vertical" size={4} className={styles.header}>
          <Typography.Title level={1}>{isRegister ? "注册 Miniport" : "登录 Miniport"}</Typography.Title>
          <Typography.Text type="secondary">
            {isRegister ? "创建账号后进入资产控制台" : "使用账号进入资产控制台"}
          </Typography.Text>
        </Space>

        {visibleError ? <Alert className={styles.alert} showIcon message={visibleError} type="error" /> : null}

        {!localLoginEnabled ? (
          <Alert showIcon type="warning" message="当前环境未启用本地账号登录" />
        ) : (
          <Form<AuthFormValues>
            form={form}
            layout="vertical"
            onFinish={(values) => void submit(values)}
            requiredMark={false}
          >
            <Form.Item name="username" label="用户名" rules={[{ required: true, message: "请输入用户名" }]}>
              <Input autoComplete="username" prefix={<UserOutlined />} placeholder="admin" />
            </Form.Item>
            <Form.Item name="password" label="密码" rules={[{ required: true, message: "请输入密码" }]}>
              <Input.Password
                autoComplete={isRegister ? "new-password" : "current-password"}
                prefix={<LockOutlined />}
                placeholder="请输入密码"
              />
            </Form.Item>
            {isRegister ? (
              <Form.Item
                name="confirmPassword"
                label="确认密码"
                rules={[{ required: true, message: "请再次输入密码" }]}
              >
                <Input.Password
                  autoComplete="new-password"
                  aria-label="确认密码"
                  prefix={<LockOutlined />}
                  placeholder="再次输入密码"
                />
              </Form.Item>
            ) : null}
            <Form.Item label="验证码" required>
              <ChallengeField
                config={config?.challenge ?? { provider: "image" }}
                disabled={activeMutation.isPending}
                onChange={setChallenge}
                onError={setChallengeError}
                resetKey={challengeResetKey}
              />
            </Form.Item>
            <Button
              block
              disabled={!challenge || (isRegister && !registrationEnabled)}
              htmlType="submit"
              loading={activeMutation.isPending}
              type="primary"
            >
              {isRegister ? "注册并登录" : "登录"}
            </Button>
          </Form>
        )}

        {localLoginEnabled && registrationEnabled ? (
          <Typography.Paragraph className={styles.modeSwitch}>
            {isRegister ? "已有账号？" : "还没有账号？"}
            <Link to={isRegister ? appPaths.login : appPaths.register}>
              {isRegister ? "返回登录" : "创建账号"}
            </Link>
          </Typography.Paragraph>
        ) : localLoginEnabled ? (
          <Typography.Paragraph className={styles.modeSwitch} type="secondary">
            注册已关闭，请联系管理员创建账号
          </Typography.Paragraph>
        ) : null}
      </Card>
    </main>
  );
}
