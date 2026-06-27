import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { ServiceItem } from "../model/portsvcTypes";
import { portRange, statusTag } from "../model/portsvcUtils";

type ServicesSectionProps = {
  canManage: boolean;
  onDeleteService: (service: ServiceItem) => void;
  onEditService: (service: ServiceItem) => void;
  onSelectService: (service: ServiceItem) => void;
  services: ServiceItem[];
};

export function ServicesSection({
  canManage,
  onDeleteService,
  onEditService,
  onSelectService,
  services,
}: ServicesSectionProps) {
  const columns: ColumnsType<ServiceItem> = [
    {
      title: "服务",
      dataIndex: "name",
      render: (_, item) => <Typography.Link onClick={() => onSelectService(item)}>{item.name}</Typography.Link>,
    },
    { title: "项目", dataIndex: "projectName", render: (value) => value || "-" },
    { title: "端口组", width: 140, render: (_, item) => portRange(item.portAllocation) },
    { title: "DIND IP", dataIndex: "dindIp", width: 150, render: (value) => value || "-" },
    { title: "容器", dataIndex: "dindContainer", render: (value) => value || "-" },
    {
      title: "仓库",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.repositories.slice(0, 3).map((repo) => (
            <Tag key={`${item.id}-${repo.id ?? repo.url}`}>{repo.name || repo.url}</Tag>
          ))}
          {item.repositories.length > 3 ? <Tag>+{item.repositories.length - 3}</Tag> : null}
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
              <Button icon={<EditOutlined />} onClick={() => onEditService(item)} />
            </Tooltip>
            <Popconfirm title="删除服务" onConfirm={() => onDeleteService(item)}>
              <Button danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        ) : null,
    },
  ];

  return (
    <Card>
      <Table<ServiceItem>
        rowKey="id"
        columns={columns}
        dataSource={services}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 1100 }}
      />
    </Card>
  );
}
