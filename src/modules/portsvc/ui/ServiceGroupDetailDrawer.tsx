import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Descriptions, Drawer, Modal, Space, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { PortGroupItem, ServiceGroupItem, ServiceGroupPortGroupItem } from "../model/portsvcTypes";
import { statusTag } from "../model/portsvcUtils";

type ServiceGroupDetailDrawerProps = {
  canManage: boolean;
  groups: PortGroupItem[];
  onClose: () => void;
  onDelete: (group: ServiceGroupItem) => Promise<void>;
  onEdit: (group: ServiceGroupItem) => void;
  serviceGroup: ServiceGroupItem | null;
};

type ServiceGroupMember = ServiceGroupPortGroupItem & {
  resolvedPortGroup?: PortGroupItem;
};

export function ServiceGroupDetailDrawer({
  canManage,
  groups,
  onClose,
  onDelete,
  onEdit,
  serviceGroup,
}: ServiceGroupDetailDrawerProps) {
  const groupByID = new Map(groups.map((group) => [group.id, group]));
  const members: ServiceGroupMember[] = serviceGroup
    ? serviceGroup.portGroups.map((member) => ({
      ...member,
      resolvedPortGroup: groupByID.get(member.portGroupId) ?? member.portGroup,
    }))
    : [];

  const columns: ColumnsType<ServiceGroupMember> = [
    {
      title: "端口组",
      width: 100,
      render: (_, item) => item.resolvedPortGroup?.portStart ?? item.portGroup?.portStart ?? "-",
    },
    {
      title: "运行环境",
      render: (_, item) => item.resolvedPortGroup?.environmentName || item.portGroup?.environmentName || "-",
    },
    { title: "角色", dataIndex: "role", width: 120, render: (value) => value || "-" },
    {
      title: "服务 IP",
      width: 150,
      render: (_, item) => item.resolvedPortGroup?.serviceIp || item.portGroup?.serviceIp || "-",
    },
    {
      title: "宿主机",
      width: 160,
      render: (_, item) => {
        const host = item.resolvedPortGroup?.host ?? item.portGroup?.host;
        return host?.name || host?.ip || "-";
      },
    },
    {
      title: "槽位",
      render: (_, item) => {
        const slots = item.resolvedPortGroup?.slots ?? [];
        return slots.length ? (
          <Space size={[4, 4]} wrap>
            {slots.map((slot) => (
              <Tag key={slot.id ?? slot.port}>{slot.name}:{slot.port}</Tag>
            ))}
          </Space>
        ) : "-";
      },
    },
    {
      title: "依赖资产",
      render: (_, item) => {
        const links = item.resolvedPortGroup?.assetLinks ?? [];
        return links.length ? (
          <Space size={[4, 4]} wrap>
            {links.map((link) => (
              <Tag key={link.id ?? `${link.assetId}-${link.relationType}`}>
                {link.relationType}: {link.asset?.name ?? link.assetId}
              </Tag>
            ))}
          </Space>
        ) : "-";
      },
    },
  ];

  return (
    <Drawer
      title={serviceGroup?.name || "服务组"}
      open={Boolean(serviceGroup)}
      size="large"
      onClose={onClose}
      extra={
        serviceGroup && canManage ? (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => onEdit(serviceGroup)}>
              编辑
            </Button>
            <Button
              danger
              icon={<DeleteOutlined />}
              onClick={() => {
                Modal.confirm({
                  title: "删除服务组",
                  content: serviceGroup.name,
                  onOk: () => onDelete(serviceGroup),
                });
              }}
            >
              删除
            </Button>
          </Space>
        ) : null
      }
    >
      {serviceGroup ? (
        <Space direction="vertical" size={16} className="content-stack">
          <Card>
            <Descriptions
              size="small"
              column={{ xs: 1, sm: 2 }}
              items={[
                { key: "owner", label: "用户", children: serviceGroup.ownerName || serviceGroup.ownerSubject || "-" },
                { key: "kind", label: "类型", children: <Tag>{serviceGroup.kind}</Tag> },
                { key: "status", label: "状态", children: statusTag(serviceGroup.status) },
                { key: "members", label: "运行环境", children: serviceGroup.portGroups.length },
              ]}
            />
            {serviceGroup.description ? <Typography.Paragraph style={{ marginTop: 12 }}>{serviceGroup.description}</Typography.Paragraph> : null}
            {serviceGroup.notes ? <Typography.Paragraph type="secondary">{serviceGroup.notes}</Typography.Paragraph> : null}
          </Card>
          <Card title="运行环境">
            <Table<ServiceGroupMember>
              rowKey={(item) => item.id ?? item.portGroupId}
              size="small"
              columns={columns}
              dataSource={members}
              pagination={false}
              scroll={{ x: 1200 }}
            />
          </Card>
        </Space>
      ) : null}
    </Drawer>
  );
}
