import { Card, Col, Row, Table, Typography } from "antd";
import type { AppStats, PortGroup } from "../model/inventoryTypes";

type DependenciesSectionProps = {
  groups: PortGroup[];
  stats: AppStats;
};

export function DependenciesSection({ groups, stats }: DependenciesSectionProps) {
  const componentRows = groups.flatMap((group) =>
    group.components.map((item) => ({
      ...item,
      group,
      key: `component-${group.id}-${item.id ?? item.name}`,
    })),
  );

  const repositoryRows = groups.flatMap((group) =>
    group.repositories.map((item) => ({
      ...item,
      group,
      key: `repo-${group.id}-${item.id ?? item.url}`,
    })),
  );

  return (
    <Row gutter={[14, 14]}>
      <Col xs={24} xl={12}>
        <Card title={`组件 ${stats.components}`}>
          <Table
            rowKey="key"
            size="small"
            dataSource={componentRows}
            pagination={{ pageSize: 8 }}
            columns={[
              { title: "名称", dataIndex: "name" },
              { title: "类型", dataIndex: "type", width: 110 },
              { title: "版本", dataIndex: "version", width: 100 },
              { title: "服务", render: (_, row) => row.group.serviceName },
            ]}
          />
        </Card>
      </Col>
      <Col xs={24} xl={12}>
        <Card title={`仓库 ${stats.repositories}`}>
          <Table
            rowKey="key"
            size="small"
            dataSource={repositoryRows}
            pagination={{ pageSize: 8 }}
            columns={[
              { title: "名称", dataIndex: "name" },
              { title: "类型", dataIndex: "kind", width: 100 },
              { title: "服务", render: (_, row) => row.group.serviceName },
              {
                title: "地址",
                dataIndex: "url",
                render: (url) => <Typography.Text copyable ellipsis>{url}</Typography.Text>,
              },
            ]}
          />
        </Card>
      </Col>
    </Row>
  );
}
