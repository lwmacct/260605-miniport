import { PlusOutlined } from "@ant-design/icons";
import { Button, Col, Progress, Row, Space, Statistic, Table } from "antd";
import Card from "antd/es/card/Card";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem } from "../model/portsvcTypes";

const PORT_PREFIX_MIN = 1000;
const PORT_PREFIX_MAX = 5999;
const SEGMENT_SIZE = 1000;
const SEGMENT_STARTS = [1000, 2000, 3000, 4000, 5000] as const;

type SegmentUsage = {
  key: number;
  segment: string;
  range: string;
  totalGroups: number;
  usedGroups: number;
  freeGroups: number;
  usagePercent: number;
  firstAvailable?: number;
};

type PortGroupsUsageSectionProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onCreateGroup: () => void;
};

function findFirstAvailable(start: number, end: number, usedPrefixes: Set<number>) {
  for (let portPrefix = start; portPrefix <= end; portPrefix++) {
    if (!usedPrefixes.has(portPrefix)) {
      return portPrefix;
    }
  }
  return undefined;
}

function buildSegmentUsage(groups: PortGroupItem[]): SegmentUsage[] {
  const usedPrefixes = new Set(
    groups
      .map((group) => group.portPrefix)
      .filter((portPrefix) => portPrefix >= PORT_PREFIX_MIN && portPrefix <= PORT_PREFIX_MAX),
  );

  return SEGMENT_STARTS.map((start) => {
    const end = Math.min(start + SEGMENT_SIZE - 1, PORT_PREFIX_MAX);
    const totalGroups = end - start + 1;
    const usedGroups = [...usedPrefixes].filter((portPrefix) => portPrefix >= start && portPrefix <= end).length;
    const freeGroups = totalGroups - usedGroups;
    const usagePercent = Number(((usedGroups / totalGroups) * 100).toFixed(1));
    const firstAvailable = findFirstAvailable(start, end, usedPrefixes);

    return {
      key: start,
      segment: `${start}-${end}`,
      range: `${start}-${end}`,
      totalGroups,
      usedGroups,
      freeGroups,
      usagePercent,
      firstAvailable,
    };
  });
}

export function PortGroupsUsageSection({ canManage, groups, onCreateGroup }: PortGroupsUsageSectionProps) {
  const segments = buildSegmentUsage(groups);
  const totalGroups = segments.reduce((sum, item) => sum + item.totalGroups, 0);
  const usedGroups = segments.reduce((sum, item) => sum + item.usedGroups, 0);
  const freeGroups = totalGroups - usedGroups;
  const usagePercent = Number(((usedGroups / totalGroups) * 100).toFixed(1));

  const columns: ColumnsType<SegmentUsage> = [
    { title: "端口组段", dataIndex: "segment", width: 120 },
    { title: "端口组范围", dataIndex: "range", width: 180 },
    { title: "总端口组", dataIndex: "totalGroups", width: 120 },
    { title: "已使用", dataIndex: "usedGroups", width: 110 },
    { title: "剩余可用", dataIndex: "freeGroups", width: 120 },
    {
      title: "占用率",
      width: 220,
      render: (_, item) => (
        <Progress percent={item.usagePercent} size="small" status={item.freeGroups === 0 ? "exception" : "normal"} />
      ),
    },
    {
      title: "首个可用端口组",
      width: 160,
      render: (_, item) => item.firstAvailable ?? "-",
    },
  ];

  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="端口组容量" value={totalGroups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="已使用端口组" value={usedGroups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="剩余端口组" value={freeGroups} /></Card>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <Card><Statistic title="整体占用率" value={usagePercent} suffix="%" /></Card>
        </Col>
      </Row>
      <Card
        title="端口组使用情况"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateGroup} disabled={!canManage}>
            新建端口组
          </Button>
        }
      >
        <Table<SegmentUsage>
          rowKey="key"
          columns={columns}
          dataSource={segments}
          pagination={false}
          scroll={{ x: 1000 }}
        />
      </Card>
    </Space>
  );
}
