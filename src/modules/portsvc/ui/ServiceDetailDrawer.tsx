import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Descriptions, Drawer, Modal, Space, Table, Tag, Typography } from "antd";
import Card from "antd/es/card/Card";
import { statusOptions } from "../model/portsvcConstants";
import type { PortGroupItem, PortSlotItem, ServiceGroupItem } from "../model/portsvcTypes";
import { portGroupLabel, runtimeTag, splitTags } from "../model/portsvcUtils";

type ServiceDetailDrawerProps = {
  canManage: boolean;
  group: PortGroupItem | null;
  onClose: () => void;
  onDelete: (group: PortGroupItem) => Promise<void>;
  onEdit: (group: PortGroupItem) => void;
  serviceGroups: ServiceGroupItem[];
};

export function ServiceDetailDrawer({
  canManage,
  group,
  onClose,
  onDelete,
  onEdit,
  serviceGroups,
}: ServiceDetailDrawerProps) {
  const relatedServiceGroups = group
    ? serviceGroups
      .map((serviceGroup) => {
        const membership = serviceGroup.portGroups.find((member) => member.portGroupId === group.id);
        return membership ? { serviceGroup, membership } : null;
      })
      .filter((item): item is NonNullable<typeof item> => Boolean(item))
    : [];

  return (
    <Drawer
      title={group?.environmentName || group?.runtimeName || "端口组"}
      open={Boolean(group)}
      size="large"
      onClose={onClose}
      extra={
        group && canManage ? (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => onEdit(group)}>
              编辑
            </Button>
            <Button
              danger
              icon={<DeleteOutlined />}
              onClick={() => {
                Modal.confirm({
                  title: "删除端口组",
                  content: portGroupLabel(group),
                  onOk: () => onDelete(group),
                });
              }}
            >
              删除
            </Button>
          </Space>
        ) : null
      }
    >
      {group ? (
        <Space direction="vertical" size={16} className="content-stack">
          <Card>
            <Descriptions
              size="small"
              column={{ xs: 1, sm: 2 }}
              items={[
                { key: "user", label: "用户", children: group.ownerName || group.ownerSubject },
                { key: "project", label: "运行环境", children: group.environmentName || "-" },
                { key: "ports", label: "端口组", children: portGroupLabel(group) },
                { key: "runtime", label: "运行模式", children: runtimeTag(group.runtimeMode) },
                { key: "serviceIp", label: "服务 IP", children: group.serviceIp || "-" },
                { key: "runtimeName", label: "运行标识", children: group.runtimeName || "-" },
                { key: "host", label: "宿主机", children: group.host?.name || group.host?.ip || "-" },
                { key: "owner", label: "负责人", children: group.environmentOwner || "-" },
                {
                  key: "status",
                  label: "状态",
                  children: statusOptions.find((item) => item.value === group.status)?.label ?? group.status,
                },
              ]}
            />
            <Space size={[4, 4]} wrap className="detail-tags">
              {splitTags(group.tags).map((tag) => (
                <Tag key={tag}>{tag}</Tag>
              ))}
            </Space>
            {group.notes ? <Typography.Paragraph>{group.notes}</Typography.Paragraph> : null}
          </Card>
          <Card title="端口槽位">
            <Table<PortSlotItem>
              rowKey={(slot: PortSlotItem) => slot.id ?? String(slot.port)}
              size="small"
              pagination={false}
              dataSource={group.slots}
              columns={[
                { title: "端口", dataIndex: "port", width: 90 },
                { title: "名称", dataIndex: "name" },
                { title: "类型", dataIndex: "kind", width: 120 },
                { title: "协议", dataIndex: "protocol", width: 120 },
                { title: "容器", dataIndex: "containerName", render: (value?: string) => value || "-" },
              ]}
            />
          </Card>
          <Card title="所属服务组">
            <Space size={[6, 6]} wrap>
              {relatedServiceGroups.length === 0 ? <Typography.Text type="secondary">未加入服务组</Typography.Text> : null}
              {relatedServiceGroups.map(({ serviceGroup, membership }) => (
                <Tag key={serviceGroup.id ?? serviceGroup.name}>
                  {serviceGroup.name}
                  {membership.role ? ` · ${membership.role}` : ""}
                </Tag>
              ))}
            </Space>
          </Card>
          <Card title="依赖资产">
            <Space size={[6, 6]} wrap>
              {group.assetLinks.length === 0 ? <Typography.Text type="secondary">未记录依赖资产</Typography.Text> : null}
              {group.assetLinks.map((link) => (
                <Tag key={`${link.id ?? link.assetId}-${link.relationType}`}>
                  {link.relationType}: {link.asset?.name ?? link.assetId}
                </Tag>
              ))}
            </Space>
          </Card>
        </Space>
      ) : null}
    </Drawer>
  );
}
