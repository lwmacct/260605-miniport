import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Descriptions, Drawer, Modal, Space, Table, Tag, Typography } from "antd";
import { statusOptions } from "../model/portsvcConstants";
import type { PortGroupItem, PortSlotItem } from "../model/portsvcTypes";
import { portRange, runtimeTag, splitTags } from "../model/portsvcUtils";

type ServiceDetailDrawerProps = {
  canManage: boolean;
  group: PortGroupItem | null;
  onClose: () => void;
  onDelete: (group: PortGroupItem) => Promise<void>;
  onEdit: (group: PortGroupItem) => void;
};

export function ServiceDetailDrawer({
  canManage,
  group,
  onClose,
  onDelete,
  onEdit,
}: ServiceDetailDrawerProps) {
  return (
    <Drawer
      title={group?.projectName || group?.runtimeName || "端口组"}
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
                  content: portRange(group),
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
                { key: "project", label: "项目", children: group.projectName || "-" },
                { key: "ports", label: "端口组", children: portRange(group) },
                { key: "runtime", label: "运行模式", children: runtimeTag(group.runtimeMode) },
                { key: "serviceIp", label: "服务 IP", children: group.serviceIp || "-" },
                { key: "runtimeName", label: "运行标识", children: group.runtimeName || "-" },
                { key: "host", label: "宿主机", children: group.host?.name || group.host?.ip || "-" },
                { key: "owner", label: "项目负责人", children: group.projectOwner || "-" },
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
              rowKey={(slot) => slot.id ?? String(slot.port)}
              size="small"
              pagination={false}
              dataSource={group.slots}
              columns={[
                { title: "端口", dataIndex: "port", width: 90 },
                { title: "名称", dataIndex: "name" },
                { title: "类型", dataIndex: "kind", width: 120 },
                { title: "协议", dataIndex: "protocol", width: 120 },
                { title: "容器", dataIndex: "containerName", render: (value) => value || "-" },
              ]}
            />
          </Card>
          <Card title="仓库">
            <Space direction="vertical" className="content-stack">
              {group.repositories.length === 0 ? <Typography.Text type="secondary">未记录仓库</Typography.Text> : null}
              {group.repositories.map((repo) => (
                <Typography.Text key={`${repo.id ?? repo.url}`} copyable>
                  {repo.name || repo.kind}: {repo.url || "-"}
                </Typography.Text>
              ))}
            </Space>
          </Card>
          <Card title="依赖">
            <Space size={[6, 6]} wrap>
              {group.dependencies.length === 0 ? <Typography.Text type="secondary">未记录依赖</Typography.Text> : null}
              {group.dependencies.map((item) => (
                <Tag key={`${item.id ?? item.name}-${item.version}`}>
                  {item.name}
                  {item.version ? ` ${item.version}` : ""}
                </Tag>
              ))}
            </Space>
          </Card>
        </Space>
      ) : null}
    </Drawer>
  );
}
