import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip } from "antd";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import type { ColumnsType } from "antd/es/table";
import type { PortGroup } from "../types";
import { statusTag } from "../utils";

type ServicesSectionProps = {
  groups: PortGroup[];
  onSelectGroup: (group: PortGroup) => void;
  onEditGroup: (group: PortGroup) => void;
  onDeleteGroup: (group: PortGroup) => void;
};

export function ServicesSection({ groups, onSelectGroup, onEditGroup, onDeleteGroup }: ServicesSectionProps) {
  const columns: ColumnsType<PortGroup> = [
    {
      title: "服务",
      dataIndex: "serviceName",
      render: (_, group) => (
        <Button type="link" className="table-link" onClick={() => onSelectGroup(group)}>
          {group.serviceName}
        </Button>
      ),
    },
    {
      title: "IP",
      width: 150,
      render: (_, group) => group.host?.ip ?? "-",
    },
    {
      title: "端口组",
      width: 140,
      render: (_, group) => `${group.portStart}-${group.portEnd}`,
    },
    {
      title: "容器",
      dataIndex: "containerName",
      render: (value) => value || "-",
    },
    {
      title: "状态",
      width: 110,
      render: (_, group) => statusTag(group.status),
    },
    {
      title: "组件",
      render: (_, group) => (
        <Space size={[4, 4]} wrap>
          {group.components.slice(0, 4).map((item) => (
            <Tag key={`${group.id}-${item.name}`}>{item.name}</Tag>
          ))}
          {group.components.length > 4 ? <Tag>+{group.components.length - 4}</Tag> : null}
        </Space>
      ),
    },
    {
      title: "操作",
      width: 120,
      render: (_, group) => (
        <Space>
          <Tooltip title="编辑">
            <Button icon={<EditOutlined />} onClick={() => onEditGroup(group)} />
          </Tooltip>
          <Popconfirm title="删除端口组" onConfirm={() => onDeleteGroup(group)}>
            <Button danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card>
      <Table<PortGroup>
        rowKey="id"
        columns={columns}
        dataSource={groups}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 980 }}
      />
    </Card>
  );
}
