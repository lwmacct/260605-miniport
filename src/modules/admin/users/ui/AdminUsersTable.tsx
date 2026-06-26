import { Space, Table, Tag, Typography, type TableColumnsType } from "antd";
import type { AdminUser } from "../api/adminUsersApi";

interface AdminUsersTableProps {
  data: AdminUser[];
  loading?: boolean;
}

export function AdminUsersTable({ data, loading }: AdminUsersTableProps) {
  const columns: TableColumnsType<AdminUser> = [
    { title: "ID", dataIndex: "id", key: "id", width: 80 },
    {
      title: "用户",
      key: "user",
      render: (_, record) => (
        <Space orientation="vertical" size={0}>
          <Typography.Text strong>{record.username}</Typography.Text>
          <Typography.Text type="secondary">{record.displayName || "-"}</Typography.Text>
        </Space>
      ),
    },
    {
      title: "角色",
      dataIndex: "admin",
      key: "admin",
      width: 120,
      render: (admin: boolean) =>
        admin ? <Tag color="blue">admin</Tag> : <Tag>user</Tag>,
    },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (status: string) =>
        status === "active" ? <Tag color="green">启用</Tag> : <Tag color="red">禁用</Tag>,
    },
    {
      title: "禁用时间",
      dataIndex: "disabledAt",
      key: "disabledAt",
      width: 190,
      render: (value?: string) => (value ? new Date(value).toLocaleString() : "-"),
    },
    {
      title: "权限来源",
      key: "source",
      width: 220,
      render: (_, record) =>
        record.admin ? "server.auth.admins 运行时配置" : "普通用户",
    },
  ];

  return (
    <Table<AdminUser>
      columns={columns}
      dataSource={data}
      loading={loading}
      pagination={{
        pageSize: 20,
        showSizeChanger: true,
        showTotal: (total) => `共 ${total} 个用户`,
      }}
      rowKey="id"
      scroll={{ x: 900 }}
      size="small"
    />
  );
}
