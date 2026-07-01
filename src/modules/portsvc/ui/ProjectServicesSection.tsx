import { Card, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem } from "../model/portsvcTypes";
import { portRange, runtimeTag } from "../model/portsvcUtils";

type ProjectServicesSectionProps = {
  groups: PortGroupItem[];
  onSelectGroup: (group: PortGroupItem) => void;
};

export function ProjectServicesSection({ groups, onSelectGroup }: ProjectServicesSectionProps) {
  const columns: ColumnsType<PortGroupItem> = [
    {
      title: "项目",
      render: (_, item) => (
        <Typography.Link onClick={() => onSelectGroup(item)}>{item.projectName || "-"}</Typography.Link>
      ),
    },
    { title: "端口组", width: 140, render: (_, item) => portRange(item) },
    { title: "运行", width: 120, render: (_, item) => runtimeTag(item.runtimeMode) },
    { title: "服务 IP", dataIndex: "serviceIp", width: 150 },
    { title: "宿主机", width: 160, render: (_, item) => item.host?.name || item.host?.ip || "-" },
    {
      title: "服务组件",
      render: (_, item) =>
        item.slots.length ? item.slots.map((slot) => <Tag key={slot.id ?? slot.port}>{slot.name}:{slot.port}</Tag>) : "-",
    },
    { title: "用户", dataIndex: "ownerName", width: 120 },
  ];

  return (
    <Card>
      <Table<PortGroupItem>
        rowKey="id"
        columns={columns}
        dataSource={groups}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 1000 }}
      />
    </Card>
  );
}
