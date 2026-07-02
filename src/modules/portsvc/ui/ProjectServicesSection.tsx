import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem, ServiceGroupItem } from "../model/portsvcTypes";
import { runtimeTag, statusTag } from "../model/portsvcUtils";

type ProjectServicesSectionProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onDeleteGroup: (group: PortGroupItem) => void;
  onEditGroup: (group: PortGroupItem) => void;
  onSelectGroup: (group: PortGroupItem) => void;
  serviceGroups: ServiceGroupItem[];
};

export function ProjectServicesSection({
  canManage,
  groups,
  onDeleteGroup,
  onEditGroup,
  onSelectGroup,
  serviceGroups,
}: ProjectServicesSectionProps) {
  const serviceGroupsByPortGroup = new Map<string, ServiceGroupItem[]>();
  for (const serviceGroup of serviceGroups) {
    for (const link of serviceGroup.portGroups) {
      const current = serviceGroupsByPortGroup.get(link.portGroupId) ?? [];
      current.push(serviceGroup);
      serviceGroupsByPortGroup.set(link.portGroupId, current);
    }
  }

  const columns: ColumnsType<PortGroupItem> = [
    {
      title: "运行环境",
      render: (_, item) => (
        <Typography.Link onClick={() => onSelectGroup(item)}>{item.environmentName || "-"}</Typography.Link>
      ),
    },
    { title: "端口组", dataIndex: "portPrefix", width: 110 },
    {
      title: "服务组",
      width: 220,
      render: (_, item) => {
        const relatedGroups = serviceGroupsByPortGroup.get(item.id) ?? [];
        return relatedGroups.length ? (
          <Space size={[4, 4]} wrap>
            {relatedGroups.map((group) => (
              <Tag key={group.id ?? group.name}>{group.name}</Tag>
            ))}
          </Space>
        ) : "-";
      },
    },
    { title: "运行", width: 120, render: (_, item) => runtimeTag(item.runtimeMode) },
    { title: "服务 IP", dataIndex: "serviceIp", width: 150 },
    { title: "宿主机", width: 160, render: (_, item) => item.host?.name || item.host?.ip || "-" },
    {
      title: "服务组件",
      render: (_, item) =>
        item.slots.length ? item.slots.map((slot) => <Tag key={slot.id ?? slot.port}>{slot.name}:{slot.port}</Tag>) : "-",
    },
    { title: "状态", width: 110, render: (_, item) => statusTag(item.status) },
    { title: "用户", dataIndex: "ownerName", width: 120 },
    {
      title: "操作",
      width: 120,
      render: (_, item) =>
        canManage ? (
          <Space>
            <Tooltip title="编辑">
              <Button icon={<EditOutlined />} onClick={() => onEditGroup(item)} />
            </Tooltip>
            <Popconfirm title="删除运行环境" onConfirm={() => onDeleteGroup(item)}>
              <Button danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        ) : null,
    },
  ];

  return (
    <Card>
      <Table<PortGroupItem>
        rowKey="id"
        columns={columns}
        dataSource={groups}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 1280 }}
      />
    </Card>
  );
}
