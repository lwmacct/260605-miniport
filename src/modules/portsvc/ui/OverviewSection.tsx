import { EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Col, Empty, Row, Space, Statistic, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { AppStats, PortAllocation, ServiceItem } from "../model/portsvcTypes";
import { portRange, statusTag } from "../model/portsvcUtils";

type OverviewSectionProps = {
  canManage: boolean;
  onCreateService: () => void;
  onEditPort: (port: PortAllocation) => void;
  onSelectService: (service: ServiceItem) => void;
  ports: PortAllocation[];
  services: ServiceItem[];
  stats: AppStats;
};

export function OverviewSection({
  canManage,
  onCreateService,
  onEditPort,
  onSelectService,
  ports,
  services,
  stats,
}: OverviewSectionProps) {
  const portColumns: ColumnsType<PortAllocation> = [
    {
      title: "端口组",
      width: 150,
      render: (_, item) => `${item.portStart}-${item.portEnd}`,
    },
    {
      title: "状态",
      width: 110,
      render: (_, item) => statusTag(item.status),
    },
    {
      title: "备注",
      dataIndex: "notes",
      render: (value) => value || "-",
    },
    {
      title: "操作",
      width: 90,
      render: (_, item) =>
        canManage ? (
          <Tooltip title="编辑端口组">
            <Button size="small" icon={<EditOutlined />} onClick={() => onEditPort(item)} />
          </Tooltip>
        ) : null,
    },
  ];

  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="服务" value={stats.services} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="端口组" value={stats.ports} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="已绑定端口组" value={stats.boundPorts} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="空闲端口组" value={stats.freePorts} /></Card>
        </Col>
      </Row>
      {services.length === 0 ? (
        <Empty description="还没有服务">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateService} disabled={!canManage}>
            新建服务
          </Button>
        </Empty>
      ) : (
        <Row gutter={[14, 14]}>
          {services.map((service) => (
            <Col key={service.id} xs={24} xl={12}>
              <button className="group-tile" onClick={() => onSelectService(service)}>
                <span className="group-range">{portRange(service.portAllocation)}</span>
                <span className="group-main">
                  <strong>{service.name}</strong>
                  <small>{service.projectName || service.dindIp || service.ownerName}</small>
                </span>
                <span>{statusTag(service.status)}</span>
                <span>
                  {service.repositories.slice(0, 2).map((repo) => (
                    <Tag key={repo.id ?? repo.url}>{repo.name || repo.kind}</Tag>
                  ))}
                </span>
              </button>
            </Col>
          ))}
        </Row>
      )}
      <Card
        title="端口组"
        extra={
          <Typography.Text type="secondary">
            {ports.length} 组
          </Typography.Text>
        }
      >
        <Table<PortAllocation>
          rowKey="id"
          columns={portColumns}
          dataSource={ports}
          pagination={{ pageSize: 12, showSizeChanger: true }}
          scroll={{ x: 720 }}
        />
      </Card>
    </Space>
  );
}
