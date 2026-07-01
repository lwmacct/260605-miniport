import { EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Col, Empty, Row, Space, Statistic, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { AppStats, PortGroupItem } from "../model/portsvcTypes";
import { portRange, runtimeTag, statusTag } from "../model/portsvcUtils";

type OverviewSectionProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onCreateGroup: () => void;
  onEditGroup: (group: PortGroupItem) => void;
  onSelectGroup: (group: PortGroupItem) => void;
  stats: AppStats;
};

export function OverviewSection({
  canManage,
  groups,
  onCreateGroup,
  onEditGroup,
  onSelectGroup,
  stats,
}: OverviewSectionProps) {
  const columns: ColumnsType<PortGroupItem> = [
    { title: "端口组", width: 140, render: (_, item) => portRange(item) },
    { title: "项目", dataIndex: "projectName", render: (value) => value || "-" },
    { title: "运行", width: 110, render: (_, item) => runtimeTag(item.runtimeMode) },
    { title: "服务 IP", dataIndex: "serviceIp", width: 150, render: (value) => value || "-" },
    { title: "状态", width: 110, render: (_, item) => statusTag(item.status) },
    {
      title: "操作",
      width: 90,
      render: (_, item) =>
        canManage ? (
          <Tooltip title="编辑端口组">
            <Button size="small" icon={<EditOutlined />} onClick={() => onEditGroup(item)} />
          </Tooltip>
        ) : null,
    },
  ];

  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="端口组" value={stats.groups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="宿主机" value={stats.hosts} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="运行中" value={stats.runningGroups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="服务组件" value={stats.slots} /></Card>
        </Col>
      </Row>
      {groups.length === 0 ? (
        <Empty description="还没有端口组">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateGroup} disabled={!canManage}>
            新建端口组
          </Button>
        </Empty>
      ) : (
        <Row gutter={[14, 14]}>
          {groups.map((group) => (
            <Col key={group.id} xs={24} xl={12}>
              <button className="group-tile" onClick={() => onSelectGroup(group)}>
                <span className="group-range">{portRange(group)}</span>
                <span className="group-main">
                  <strong>{group.projectName || group.runtimeName || "未命名项目"}</strong>
                  <small>{group.serviceIp || group.host?.name || group.ownerName}</small>
                </span>
                <span>{statusTag(group.status)}</span>
                <span>
                  {group.slots.slice(0, 2).map((slot) => (
                    <Tag key={slot.id ?? slot.port}>{slot.name || slot.port}</Tag>
                  ))}
                </span>
              </button>
            </Col>
          ))}
        </Row>
      )}
      <Card title="端口组" extra={<Typography.Text type="secondary">{groups.length} 组</Typography.Text>}>
        <Table<PortGroupItem>
          rowKey="id"
          columns={columns}
          dataSource={groups}
          pagination={{ pageSize: 12, showSizeChanger: true }}
          scroll={{ x: 800 }}
        />
      </Card>
    </Space>
  );
}
