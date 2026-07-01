import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem } from "../model/portsvcTypes";
import { portRange, runtimeTag, statusTag } from "../model/portsvcUtils";

type ServicesSectionProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onDeleteGroup: (group: PortGroupItem) => void;
  onEditGroup: (group: PortGroupItem) => void;
  onSelectGroup: (group: PortGroupItem) => void;
};

export function ServicesSection({
  canManage,
  groups,
  onDeleteGroup,
  onEditGroup,
  onSelectGroup,
}: ServicesSectionProps) {
  const columns: ColumnsType<PortGroupItem> = [
    {
      title: "项目",
      dataIndex: "projectName",
      render: (_, item) => (
        <Typography.Link onClick={() => onSelectGroup(item)}>{item.projectName || portRange(item)}</Typography.Link>
      ),
    },
    { title: "端口组", width: 140, render: (_, item) => portRange(item) },
    { title: "运行", width: 110, render: (_, item) => runtimeTag(item.runtimeMode) },
    { title: "服务 IP", dataIndex: "serviceIp", width: 150, render: (value) => value || "-" },
    { title: "宿主机", width: 160, render: (_, item) => item.host?.name || item.host?.ip || "-" },
    { title: "运行标识", dataIndex: "runtimeName", render: (value) => value || "-" },
    {
      title: "槽位",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.slots.slice(0, 4).map((slot) => (
            <Tag key={`${item.id}-${slot.id ?? slot.port}`}>{slot.name || slot.port}</Tag>
          ))}
          {item.slots.length > 4 ? <Tag>+{item.slots.length - 4}</Tag> : null}
        </Space>
      ),
    },
    { title: "状态", width: 110, render: (_, item) => statusTag(item.status) },
    {
      title: "操作",
      width: 120,
      render: (_, item) =>
        canManage ? (
          <Space>
            <Tooltip title="编辑">
              <Button icon={<EditOutlined />} onClick={() => onEditGroup(item)} />
            </Tooltip>
            <Popconfirm title="删除端口组" onConfirm={() => onDeleteGroup(item)}>
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
        scroll={{ x: 1200 }}
      />
    </Card>
  );
}
