import { Card, Col, Row, Table, Typography } from "antd";
import type { AppStats, ServiceItem } from "../model/portsvcTypes";

type DependenciesSectionProps = {
  services: ServiceItem[];
  stats: AppStats;
};

export function DependenciesSection({ services, stats }: DependenciesSectionProps) {
  const dependencyRows = services.flatMap((service) =>
    service.dependencies.map((item) => ({
      ...item,
      service,
      key: `dependency-${service.id}-${item.id ?? item.name}`,
    })),
  );

  const repositoryRows = services.flatMap((service) =>
    service.repositories.map((item) => ({
      ...item,
      service,
      key: `repo-${service.id}-${item.id ?? item.url}`,
    })),
  );

  return (
    <Row gutter={[14, 14]}>
      <Col xs={24} xl={12}>
        <Card title={`依赖 ${stats.dependencies}`}>
          <Table
            rowKey="key"
            size="small"
            dataSource={dependencyRows}
            pagination={{ pageSize: 8 }}
            columns={[
              { title: "名称", dataIndex: "name" },
              { title: "类型", dataIndex: "type", width: 110 },
              { title: "版本", dataIndex: "version", width: 100 },
              { title: "服务", render: (_, row) => row.service.name },
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
              { title: "服务", render: (_, row) => row.service.name },
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
