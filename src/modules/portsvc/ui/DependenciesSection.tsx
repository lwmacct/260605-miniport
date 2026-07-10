import { EditOutlined } from "@ant-design/icons";
import { WorkbenchPanel } from "@lwmacct/260627-antd-workbench";
import { Button, Col, Row, Space, Statistic, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { AppStats, DependencyAssetItem, PortGroupItem } from "../model/portsvcTypes";

type DependenciesSectionProps = {
  canManage: boolean;
  dependencyAssets: DependencyAssetItem[];
  groups: PortGroupItem[];
  onEditAsset: (asset: DependencyAssetItem) => void;
  stats: AppStats;
};

export function DependenciesSection({ canManage, dependencyAssets, groups, onEditAsset, stats }: DependenciesSectionProps) {
  const assetColumns: ColumnsType<DependencyAssetItem> = [
    { title: "名称", dataIndex: "name", render: (value, item) => value || item.fullName || item.url },
    { title: "类别", width: 120, render: (_, item) => <Tag>{item.assetKind}</Tag> },
    { title: "类型", width: 150, render: (_, item) => <Tag>{item.assetType}</Tag> },
    { title: "Provider", dataIndex: "provider", width: 120 },
    {
      title: "地址",
      dataIndex: "url",
      render: (url) => (url ? <Typography.Text copyable ellipsis>{url}</Typography.Text> : "-"),
    },
    { title: "可控性", dataIndex: "controllability", width: 120 },
    { title: "状态", dataIndex: "status", width: 100 },
    {
      title: "操作",
      width: 80,
      render: (_, item) =>
        canManage ? (
          <Tooltip title="编辑资产">
            <Button size="small" icon={<EditOutlined />} onClick={() => onEditAsset(item)} />
          </Tooltip>
        ) : null,
    },
  ];

  const groupColumns: ColumnsType<PortGroupItem> = [
    { title: "运行环境", dataIndex: "environmentName", render: (value) => value || "-" },
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
      title: "依赖资产",
      render: (_, item) => (
        <Space size={[4, 4]} wrap>
          {item.assetLinks.length === 0 ? <Typography.Text type="secondary">-</Typography.Text> : null}
          {item.assetLinks.map((link) => (
            <Tag key={link.id ?? `${link.assetId}-${link.relationType}`}>
              {link.relationType}: {link.asset?.name ?? link.assetId}
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
          <WorkbenchPanel><Statistic title="服务组件" value={stats.slots} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={8}>
          <WorkbenchPanel><Statistic title="依赖资产" value={stats.dependencyAssets} /></WorkbenchPanel>
        </Col>
        <Col xs={24} sm={8}>
          <WorkbenchPanel><Statistic title="资产关系" value={stats.assetLinks} /></WorkbenchPanel>
        </Col>
      </Row>
      <WorkbenchPanel title="依赖资产">
        <Table<DependencyAssetItem>
          rowKey={(item: DependencyAssetItem) => item.id ?? item.name}
          columns={assetColumns}
          dataSource={dependencyAssets}
          pagination={{ pageSize: 10 }}
          scroll={{ x: 1000 }}
        />
      </WorkbenchPanel>
      <WorkbenchPanel title="端口组依赖关系">
        <Table<PortGroupItem>
          rowKey="id"
          columns={groupColumns}
          dataSource={groups}
          pagination={{ pageSize: 12 }}
          scroll={{ x: 900 }}
        />
      </WorkbenchPanel>
    </Space>
  );
}
