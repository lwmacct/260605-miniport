import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tooltip } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { Host, PortGroup } from "../model/inventoryTypes";

type HostsSectionProps = {
  canManage: boolean;
  groups: PortGroup[];
  hosts: Host[];
  onDeleteHost: (host: Host) => void;
  onEditHost: (host: Host) => void;
};

export function HostsSection({ canManage, groups, hosts, onDeleteHost, onEditHost }: HostsSectionProps) {
  const columns: ColumnsType<Host> = [
    { title: "IP", dataIndex: "ip", width: 160 },
    { title: "名称", dataIndex: "name" },
    { title: "网段", dataIndex: "network" },
    { title: "环境", dataIndex: "environment", width: 120 },
    {
      title: "端口组",
      width: 100,
      render: (_, host) => groups.filter((group) => group.hostId === host.id).length,
    },
    {
      title: "操作",
      width: 120,
      render: (_, host) =>
        canManage ? (
          <Space>
            <Tooltip title="编辑">
              <Button icon={<EditOutlined />} onClick={() => onEditHost(host)} />
            </Tooltip>
            <Popconfirm title="删除主机" onConfirm={() => onDeleteHost(host)}>
              <Button danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        ) : null,
    },
  ];

  return (
    <Card>
      <Table<Host>
        rowKey="id"
        columns={columns}
        dataSource={hosts}
        pagination={false}
        scroll={{ x: 760 }}
      />
    </Card>
  );
}
