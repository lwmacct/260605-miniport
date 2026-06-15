import { Alert, Card, Table, Tag, Typography } from "antd";
import type { TableColumnsType } from "antd";
import type { AdminUser } from "../../domains/admin-users/api";
import { useAdminUsersQuery } from "../../domains/admin-users/queries";

const userColumns: TableColumnsType<AdminUser> = [
  { title: "用户名", dataIndex: "username", key: "username" },
  { title: "显示名", dataIndex: "displayName", key: "displayName" },
  {
    title: "状态",
    dataIndex: "status",
    key: "status",
    render: (status: string) => (status === "active" ? <Tag color="green">active</Tag> : <Tag>{status}</Tag>),
  },
  {
    title: "管理员",
    dataIndex: "admin",
    key: "admin",
    render: (admin: boolean) => (admin ? <Tag color="blue">admin</Tag> : <Tag>viewer</Tag>),
  },
];

export function AdminUsersPage() {
  const users = useAdminUsersQuery();

  return (
    <Card>
      <Typography.Title level={4}>管理员用户</Typography.Title>
      <Typography.Paragraph type="secondary">
        管理员权限仅来自 `server.auth.admins` 运行时配置。
      </Typography.Paragraph>
      {users.isError ? (
        <Alert showIcon type="error" message={users.error.message} />
      ) : (
        <Table<AdminUser>
          rowKey="id"
          columns={userColumns}
          dataSource={users.data ?? []}
          loading={users.isPending}
          pagination={false}
        />
      )}
    </Card>
  );
}
