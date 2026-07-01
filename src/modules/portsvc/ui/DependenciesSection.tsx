import { Card, Col, Row, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { AppStats, PortGroupItem } from "../model/portsvcTypes";

type DependenciesSectionProps = {
  groups: PortGroupItem[];
  stats: AppStats;
};

export function DependenciesSection({ groups, stats }: DependenciesSectionProps) {
  const columns: ColumnsType<PortGroupItem> = [
    { title: "项目", dataIndex: "projectName", render: (value) => value || "-" },
    {
      title: "服务组件",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.slots.map((slot) => (
            <Tag key={slot.id ?? slot.port}>{slot.name}:{slot.port}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: "仓库",
      render: (_, item) => (
        <Space direction="vertical" size={2}>
          {item.repositories.length === 0 ? <Typography.Text type="secondary">-</Typography.Text> : null}
          {item.repositories.map((repo) => (
            <Typography.Text key={repo.id ?? repo.url} copyable ellipsis>
              {repo.name || repo.kind}: {repo.url}
            </Typography.Text>
          ))}
        </Space>
      ),
    },
    {
      title: "依赖",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.dependencies.map((dep) => (
            <Tag key={dep.id ?? `${dep.name}-${dep.version}`}>
              {dep.name}
              {dep.version ? ` ${dep.version}` : ""}
            </Tag>
          ))}
        </Space>
      ),
    },
  ];

  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={8}>
          <Card><Statistic title="服务组件" value={stats.slots} /></Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card><Statistic title="仓库" value={stats.repositories} /></Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card><Statistic title="依赖" value={stats.dependencies} /></Card>
        </Col>
      </Row>
      <Card>
        <Table<PortGroupItem>
          rowKey="id"
          columns={columns}
          dataSource={groups}
          pagination={{ pageSize: 12 }}
          scroll={{ x: 1000 }}
        />
      </Card>
    </Space>
  );
}
