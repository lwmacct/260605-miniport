import { Button, Card, Empty, Space, Tag, Typography } from "antd";
import { EditOutlined, PlusOutlined } from "@ant-design/icons";
import type { AppStats, Host, PortGroup } from "../types";
import { statusTag } from "../utils";
import { Metric } from "../components/Metric";

type OverviewSectionProps = {
  stats: AppStats;
  hosts: Host[];
  groupsByHost: Array<{ host: Host; groups: PortGroup[] }>;
  onCreateHost: () => void;
  onEditHost: (host: Host) => void;
  onSelectGroup: (group: PortGroup) => void;
};

export function OverviewSection({
  stats,
  hosts,
  groupsByHost,
  onCreateHost,
  onEditHost,
  onSelectGroup,
}: OverviewSectionProps) {
  return (
    <Space orientation="vertical" size={16} className="content-stack">
      <div className="metric-grid">
        <Card><Metric label="主机" value={stats.hosts} /></Card>
        <Card><Metric label="端口组" value={stats.groups} /></Card>
        <Card><Metric label="已分配端口" value={stats.usedSlots} /></Card>
        <Card><Metric label="空闲槽位" value={stats.emptySlots} /></Card>
      </div>
      {hosts.length === 0 ? (
        <Empty description="还没有主机">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateHost}>
            新建主机
          </Button>
        </Empty>
      ) : (
        <div className="host-grid">
          {groupsByHost.map(({ host, groups }) => (
            <Card
              key={host.id}
              title={
                <Space>
                  <span>{host.ip}</span>
                  {host.environment ? <Tag>{host.environment}</Tag> : null}
                </Space>
              }
              extra={<Button size="small" icon={<EditOutlined />} onClick={() => onEditHost(host)} />}
            >
              <Typography.Text type="secondary">{host.name || host.network || "未命名主机"}</Typography.Text>
              <div className="group-stack">
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
              </div>
            </Card>
          ))}
        </div>
      )}
    </Space>
  );
}
