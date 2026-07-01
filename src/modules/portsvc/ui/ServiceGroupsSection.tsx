import { EditOutlined } from "@ant-design/icons";
import { Button, Card, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { ServiceGroupItem } from "../model/portsvcTypes";
import { statusTag } from "../model/portsvcUtils";

type ServiceGroupsSectionProps = {
  canManage: boolean;
  onEditServiceGroup: (group: ServiceGroupItem) => void;
  onSelectServiceGroup: (group: ServiceGroupItem) => void;
  serviceGroups: ServiceGroupItem[];
};

export function ServiceGroupsSection({
  canManage,
  onEditServiceGroup,
  onSelectServiceGroup,
  serviceGroups,
}: ServiceGroupsSectionProps) {
  const columns: ColumnsType<ServiceGroupItem> = [
    {
      title: "服务组",
      dataIndex: "name",
      render: (value, item) => (
        <Typography.Link onClick={() => onSelectServiceGroup(item)}>{value || "-"}</Typography.Link>
      ),
    },
    { title: "类型", dataIndex: "kind", width: 110, render: (value) => <Tag>{value}</Tag> },
    {
      title: "运行环境",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.portGroups.length === 0 ? <Typography.Text type="secondary">-</Typography.Text> : null}
          {item.portGroups.map((link) => (
            <Tag key={link.id ?? link.portGroupId}>
              {link.portGroup?.portStart ?? link.portGroupId}
              {link.role ? ` · ${link.role}` : ""}
            </Tag>
          ))}
        </Space>
      ),
    },
    { title: "状态", dataIndex: "status", width: 110, render: (value) => statusTag(value) },
    { title: "用户", dataIndex: "ownerName", width: 120, render: (value) => value || "-" },
    {
      title: "说明",
      dataIndex: "description",
      render: (value) => (value ? <Typography.Text ellipsis>{value}</Typography.Text> : "-"),
    },
    {
      title: "操作",
      width: 80,
      render: (_, item) =>
        canManage ? (
          <Tooltip title="编辑服务组">
            <Button size="small" icon={<EditOutlined />} onClick={() => onEditServiceGroup(item)} />
          </Tooltip>
        ) : null,
    },
  ];

  return (
    <Card>
      <Table<ServiceGroupItem>
        rowKey={(item) => item.id ?? item.name}
        columns={columns}
        dataSource={serviceGroups}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 1000 }}
      />
    </Card>
  );
}
