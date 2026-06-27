import { Card, Space, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroup } from "../model/inventoryTypes";

type HostsSectionProps = {
  groups: PortGroup[];
  onSelectGroup: (group: PortGroup) => void;
};

export function HostsSection({ groups, onSelectGroup }: HostsSectionProps) {
  const columns: ColumnsType<PortGroup> = [
    {
      title: "项目",
      render: (_, group) => (
        <Space size={[4, 4]} wrap>
          {group.projects.length === 0 ? <Typography.Text type="secondary">未记录项目</Typography.Text> : null}
          {group.projects.map((project) => (
            <Tag key={`${group.id}-${project.name}`}>{project.name}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: "端口分配",
      render: (_, group) => <Typography.Link onClick={() => onSelectGroup(group)}>{group.name}</Typography.Link>,
    },
    { title: "端口组", width: 140, render: (_, group) => `${group.portStart}-${group.portEnd}` },
    { title: "DIND IP", dataIndex: "dindIp", width: 150 },
    { title: "容器", dataIndex: "dindContainer" },
    { title: "用户", dataIndex: "username", width: 120 },
  ];

  return (
    <Card>
      <Table<PortGroup>
        rowKey="id"
        columns={columns}
        dataSource={groups}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 900 }}
      />
    </Card>
  );
}
