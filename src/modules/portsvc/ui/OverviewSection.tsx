import { PlusOutlined } from "@ant-design/icons";
import { WorkbenchPanel } from "@lwmacct/260627-antd-workbench";
import { Button, Col, Empty, Row, Space, Statistic } from "antd";
import type { AppStats } from "../model/portsvcTypes";

type OverviewSectionProps = {
  canManage: boolean;
  onCreateGroup: () => void;
  stats: AppStats;
};

export function OverviewSection({
  canManage,
  onCreateGroup,
  stats,
}: OverviewSectionProps) {
  return (
    <Space direction="vertical" size={16} className="content-stack">
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="端口组" value={stats.groups} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="宿主机" value={stats.hosts} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="运行中" value={stats.runningGroups} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="服务组件" value={stats.slots} /></WorkbenchPanel>
        </Col>
      </Row>
      <Row gutter={[12, 12]}>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="空闲端口组" value={stats.freeGroups} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="依赖资产" value={stats.dependencyAssets} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="服务组" value={stats.serviceGroups} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={12} xl={6}>
          <WorkbenchPanel><Statistic title="资产关系" value={stats.assetLinks} /></WorkbenchPanel>
        </Col>
      </Row>
      {stats.groups === 0 ? (
        <Empty description="还没有端口组">
          <Button type="primary" icon={<PlusOutlined />} onClick={onCreateGroup} disabled={!canManage}>
            新建端口组
          </Button>
        </Empty>
      ) : (
        <WorkbenchPanel>
          <Statistic title="已管理端口" value={stats.groups * 10} suffix="个" />
        </WorkbenchPanel>
      )}
    </Space>
  );
}
