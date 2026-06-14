import { EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Col, Empty, Row, Space, Statistic, Tag, Typography } from "antd";
import type { AppStats, Host, PortGroup } from "../types";
import { statusTag } from "../utils";

type OverviewSectionProps = {
  groupsByHost: Array<{ host: Host; groups: PortGroup[] }>;
  hosts: Host[];
  onCreateHost: () => void;
  onEditHost: (host: Host) => void;
  onSelectGroup: (group: PortGroup) => void;
  stats: AppStats;
};

export function OverviewSection({
  groupsByHost,
  hosts,
  onCreateHost,
  onEditHost,
  onSelectGroup,
  stats,
}: OverviewSectionProps) {
  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="主机" value={stats.hosts} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="端口组" value={stats.groups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="已分配端口" value={stats.usedSlots} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="空闲槽位" value={stats.emptySlots} /></Card>
        </Col>
      </Row>
      {hosts.length === 0 ? (
        <Empty description="还没有主机">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateHost}>
            新建主机
          </Button>
        </Empty>
      ) : (
        <Row gutter={[14, 14]}>
          {groupsByHost.map(({ host, groups }) => (
            <Col key={host.id} xs={24} xl={12}>
              <Card
                title={
                  <Space>
                    <span>{host.ip}</span>
                    {host.environment ? <Tag>{host.environment}</Tag> : null}
                  </Space>
                }
                extra={<Button size="small" icon={<EditOutlined />} onClick={() => onEditHost(host)} />}
              >
                <Typography.Text type="secondary">{host.name || host.network || "未命名主机"}</Typography.Text>
                <Space direction="vertical" size={8} className="content-stack host-groups">
                  {groups.length === 0 ? (
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="无端口组" />
                  ) : (
                    groups.map((group) => (
                      <button key={group.id} className="group-tile" onClick={() => onSelectGroup(group)}>
                        <span className="group-range">{group.portStart}-{group.portEnd}</span>
                        <span className="group-main">
                          <strong>{group.serviceName}</strong>
                          <small>{group.containerName || group.dindHost || "未填写容器"}</small>
                        </span>
                        <span>{statusTag(group.status)}</span>
                        <span className="slot-strip">
                          {group.slots.map((slot) => (
                            <i
                              key={slot.port}
                              className={`slot-dot slot-${slot.status}`}
                              title={`${slot.port} ${slot.name || ""}`}
                            />
                          ))}
                        </span>
                      </button>
                    ))
                  )}
                </Space>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </Space>
  );
}
