import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Descriptions, Drawer, Modal, Space, Table, Tag, Typography } from "antd";
import { statusOptions } from "../constants";
import type { PortGroup, PortSlot } from "../types";
import { splitTags } from "../utils";

type GroupDetailDrawerProps = {
  group: PortGroup | null;
  onClose: () => void;
  onDelete: (group: PortGroup) => Promise<void>;
  onEdit: (group: PortGroup) => void;
};

export function GroupDetailDrawer({
  group,
  onClose,
  onDelete,
  onEdit,
}: GroupDetailDrawerProps) {
  return (
    <Drawer
      title={group?.serviceName}
      open={Boolean(group)}
      size="large"
      onClose={onClose}
      extra={
        group ? (
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
                  content: `${group.host?.ip ?? ""} ${group.portStart}-${group.portEnd}`,
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
                { key: "ip", label: "IP", children: group.host?.ip ?? "-" },
                { key: "ports", label: "端口组", children: `${group.portStart}-${group.portEnd}` },
                { key: "container", label: "容器", children: group.containerName || "-" },
                { key: "dind", label: "DIND", children: group.dindHost || "-" },
                { key: "owner", label: "负责人", children: group.owner || "-" },
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
            <Table<PortSlot>
              rowKey="port"
              size="small"
              pagination={false}
              dataSource={group.slots}
              columns={[
                { title: "端口", dataIndex: "port", width: 90 },
                { title: "名称", dataIndex: "name" },
                { title: "协议", dataIndex: "protocol", width: 90 },
                { title: "状态", dataIndex: "status", width: 90 },
                { title: "用途", dataIndex: "purpose" },
              ]}
            />
          </Card>
          <Card title="组件">
            <Space size={[6, 6]} wrap>
              {group.components.length === 0 ? <Typography.Text type="secondary">未记录组件</Typography.Text> : null}
              {group.components.map((item) => (
                <Tag key={item.name}>
                  {item.name}
                  {item.version ? ` ${item.version}` : ""}
                </Tag>
              ))}
            </Space>
          </Card>
          <Card title="仓库">
            <Space direction="vertical" className="content-stack">
              {group.repositories.length === 0 ? <Typography.Text type="secondary">未记录仓库</Typography.Text> : null}
              {group.repositories.map((repo) => (
                <Typography.Text key={repo.url} copyable>
                  {repo.name}: {repo.url}
                </Typography.Text>
              ))}
            </Space>
          </Card>
        </Space>
      ) : null}
    </Drawer>
  );
}
