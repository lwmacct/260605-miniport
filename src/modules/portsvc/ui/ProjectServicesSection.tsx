import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem } from "../model/portsvcTypes";
import { runtimeTag, statusTag } from "../model/portsvcUtils";

type ProjectServicesSectionProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onDeleteGroup: (group: PortGroupItem) => void;
  onEditGroup: (group: PortGroupItem) => void;
  onSelectGroup: (group: PortGroupItem) => void;
};

export function ProjectServicesSection({
  canManage,
  groups,
  onDeleteGroup,
  onEditGroup,
  onSelectGroup,
}: ProjectServicesSectionProps) {
  const columns: ColumnsType<PortGroupItem> = [
    {
      title: "运行环境",
      render: (_, item) => (
        <Typography.Link onClick={() => onSelectGroup(item)}>{item.projectName || "-"}</Typography.Link>
      ),
    },
    { title: "端口组", dataIndex: "portStart", width: 110 },
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
        scroll={{ x: 1080 }}
      />
    </Card>
  );
}
