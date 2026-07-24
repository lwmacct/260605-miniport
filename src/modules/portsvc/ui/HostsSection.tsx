import { EditOutlined } from "@ant-design/icons";
import { WorkbenchPanel } from "@lwmacct/260627-antd-workbench";
import { Button, Space, Table, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { HostItem, PortGroupItem } from "../model/portsvcTypes";
import { statusTag } from "../model/portsvcUtils";

type HostsSectionProps = {
  groups: PortGroupItem[];
  hosts: HostItem[];
  onEditHost: (host: HostItem) => void;
};

export function HostsSection({ groups, hosts, onEditHost }: HostsSectionProps) {
  const groupCountByHost = new Map<string, number>();
  for (const group of groups) {
    if (!group.hostId) {
      continue;
    }
    groupCountByHost.set(group.hostId, (groupCountByHost.get(group.hostId) ?? 0) + 1);
  }

  const columns: ColumnsType<HostItem> = [
    {
      title: "名称",
      dataIndex: "name",
      render: (value) => value || "-",
    },
    {
      title: "IP",
      dataIndex: "ip",
      width: 160,
      render: (value) => value || "-",
    },
    {
      title: "规格",
      dataIndex: "spec",
      width: 120,
      render: (value) => value || "-",
    },
    {
      title: "端口组",
      width: 110,
      render: (_, item) => groupCountByHost.get(item.id) ?? 0,
    },
    {
      title: "状态",
      dataIndex: "status",
      width: 110,
      render: (value) => statusTag(value),
    },
    {
      title: "备注",
      dataIndex: "notes",
      render: (value) => (value ? <Typography.Text ellipsis>{value}</Typography.Text> : "-"),
    },
    {
      title: "操作",
      width: 80,
      render: (_, item) => (
        <Space>
          <Tooltip title="编辑宿主机">
            <Button size="small" icon={<EditOutlined />} onClick={() => onEditHost(item)} />
          </Tooltip>
        </Space>
      ),
    },
  ];

  return (
    <WorkbenchPanel>
      <Table<HostItem>
        rowKey="id"
        columns={columns}
        dataSource={hosts}
        pagination={{ pageSize: 12 }}
        scroll={{ x: 900 }}
      />
    </WorkbenchPanel>
  );
}
