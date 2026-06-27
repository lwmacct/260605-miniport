import { PlusOutlined } from "@ant-design/icons";
import { Button, Card, Col, Empty, Row, Space, Statistic, Tag } from "antd";
import type { AppStats, PortGroup } from "../model/inventoryTypes";
import { statusTag } from "../model/inventoryUtils";

type OverviewSectionProps = {
  canManage: boolean;
  groups: PortGroup[];
  onCreateGroup: () => void;
  onSelectGroup: (group: PortGroup) => void;
  stats: AppStats;
};

export function OverviewSection({
  canManage,
  groups,
  onCreateGroup,
  onSelectGroup,
  stats,
}: OverviewSectionProps) {
  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="端口分配" value={stats.allocations} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="用户" value={stats.users} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="已用端口" value={stats.usedSlots} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="空闲槽位" value={stats.emptySlots} /></Card>
        </Col>
      </Row>
      {groups.length === 0 ? (
        <Empty description="还没有端口分配">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateGroup} disabled={!canManage}>
            新建端口分配
          </Button>
        </Empty>
      ) : (
        <Row gutter={[14, 14]}>
          {groups.map((group) => (
            <Col key={group.id} xs={24} xl={12}>
              <button className="group-tile" onClick={() => onSelectGroup(group)}>
                <span className="group-range">{group.portStart}-{group.portEnd}</span>
                <span className="group-main">
                  <strong>{group.name}</strong>
                  <small>{group.dindIp || group.dindContainer || group.username}</small>
                </span>
                <span>{statusTag(group.status)}</span>
                <span>
                  {group.projects.slice(0, 2).map((project) => (
                    <Tag key={project.name}>{project.name}</Tag>
                  ))}
                </span>
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
            </Col>
          ))}
        </Row>
      )}
    </Space>
  );
}
