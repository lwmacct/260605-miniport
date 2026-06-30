import { Card, Table, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { ServiceItem } from "../model/portsvcTypes";
import { portRange } from "../model/portsvcUtils";

type ProjectServicesSectionProps = {
  onSelectService: (service: ServiceItem) => void;
  services: ServiceItem[];
};

export function ProjectServicesSection({ onSelectService, services }: ProjectServicesSectionProps) {
  const columns: ColumnsType<ServiceItem> = [
    { title: "项目", dataIndex: "projectName", render: (value) => value || "-" },
    {
      title: "服务",
      render: (_, item) => <Typography.Link onClick={() => onSelectService(item)}>{item.name}</Typography.Link>,
    },
    { title: "端口组", width: 140, render: (_, item) => portRange(item.portAllocation) },
    { title: "DIND IP", dataIndex: "dindIp", width: 150 },
    { title: "容器", dataIndex: "dindContainer" },
    { title: "用户", dataIndex: "ownerName", width: 120 },
  ];

  return (
    <Card>
      <Table<ServiceItem>
        rowKey="id"
        columns={columns}
        dataSource={services}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 900 }}
      />
    </Card>
  );
}
