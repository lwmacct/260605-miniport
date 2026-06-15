import { Alert, Button, Card, Form, Input, Space, Typography, message } from "antd";
import { useNavigate } from "react-router-dom";
import { appPaths } from "../../app/router/navigation";
import { useLoginMutation, useRegisterMutation } from "../../domains/auth/mutations";
import { useAuthStateQuery } from "../../domains/auth/queries";

type LoginForm = {
  username: string;
  password: string;
};

export function LoginPage() {
  const navigate = useNavigate();
  const authState = useAuthStateQuery();
  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const [form] = Form.useForm<LoginForm>();

  async function handleLogin(values: LoginForm) {
    try {
      await loginMutation.mutateAsync(values);
      navigate(appPaths.overview, { replace: true });
    } catch (error) {
      message.error(error instanceof Error ? error.message : "登录失败");
    }
  }

  async function handleRegister() {
    try {
      const values = await form.validateFields();
      await registerMutation.mutateAsync(values);
      navigate(appPaths.overview, { replace: true });
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message);
      }
    }
  }

  return (
    <main className="app-content" style={{ display: "grid", minHeight: "100vh", placeItems: "center", padding: 24 }}>
      <Card style={{ width: "100%", maxWidth: 420 }}>
        <Space direction="vertical" size={20} style={{ width: "100%" }}>
          <div>
            <Typography.Title level={3} style={{ marginBottom: 8 }}>
              登录 Miniport
            </Typography.Title>
            <Typography.Text type="secondary">
              登录后可查看资产，管理员账号可执行写操作。
            </Typography.Text>
          </div>

          {!authState.data?.config.local.loginEnabled ? (
            <Alert showIcon type="warning" message="当前环境未启用本地账号登录" />
          ) : (
            <Form<LoginForm> form={form} layout="vertical" onFinish={(values) => void handleLogin(values)}>
              <Form.Item name="username" label="用户名" rules={[{ required: true, message: "请输入用户名" }]}>
                <Input autoComplete="username" placeholder="admin" />
              </Form.Item>
              <Form.Item name="password" label="密码" rules={[{ required: true, message: "请输入密码" }]}>
                <Input.Password autoComplete="current-password" placeholder="请输入密码" />
              </Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" loading={loginMutation.isPending}>
                  登录
                </Button>
                {authState.data?.config.local.registrationEnabled ? (
                  <Button onClick={() => void handleRegister()} loading={registerMutation.isPending}>
                    注册并登录
                  </Button>
                ) : null}
              </Space>
            </Form>
          )}
        </Space>
      </Card>
    </main>
  );
}
