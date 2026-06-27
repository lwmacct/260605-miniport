import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Descriptions, Drawer, Modal, Space, Tag, Typography } from "antd";
import { statusOptions } from "../model/portsvcConstants";
import type { ServiceItem } from "../model/portsvcTypes";
import { portRange, splitTags } from "../model/portsvcUtils";

type ServiceDetailDrawerProps = {
  canManage: boolean;
  onClose: () => void;
  onDelete: (service: ServiceItem) => Promise<void>;
  onEdit: (service: ServiceItem) => void;
  service: ServiceItem | null;
};

export function ServiceDetailDrawer({
  canManage,
  onClose,
  onDelete,
  onEdit,
  service,
}: ServiceDetailDrawerProps) {
  return (
    <Drawer
      title={service?.name}
      open={Boolean(service)}
      size="large"
      onClose={onClose}
      extra={
        service && canManage ? (
          <Space>
            <Button icon={<EditOutlined />} onClick={() => onEdit(service)}>
              编辑
            </Button>
            <Button
              danger
              icon={<DeleteOutlined />}
              onClick={() => {
                Modal.confirm({
                  title: "删除服务",
                  content: service.name,
                  onOk: () => onDelete(service),
                });
              }}
            >
              删除
            </Button>
          </Space>
        ) : null
      }
    >
      {service ? (
        <Space direction="vertical" size={16} className="content-stack">
          <Card>
            <Descriptions
              size="small"
              column={{ xs: 1, sm: 2 }}
              items={[
                { key: "user", label: "用户", children: service.username || service.userId },
                { key: "project", label: "项目", children: service.projectName || "-" },
                { key: "ports", label: "端口组", children: portRange(service.portAllocation) },
                { key: "dindIp", label: "DIND IP", children: service.dindIp || "-" },
                { key: "container", label: "容器", children: service.dindContainer || "-" },
                { key: "owner", label: "负责人", children: service.owner || "-" },
                {
                  key: "status",
                  label: "状态",
                  children: statusOptions.find((item) => item.value === service.status)?.label ?? service.status,
                },
              ]}
            />
            <Space size={[4, 4]} wrap className="detail-tags">
              {splitTags(service.tags).map((tag) => (
                <Tag key={tag}>{tag}</Tag>
              ))}
            </Space>
            {service.notes ? <Typography.Paragraph>{service.notes}</Typography.Paragraph> : null}
          </Card>
          <Card title="仓库">
            <Space direction="vertical" className="content-stack">
              {service.repositories.length === 0 ? <Typography.Text type="secondary">未记录仓库</Typography.Text> : null}
              {service.repositories.map((repo) => (
                <Typography.Text key={`${repo.id ?? repo.url}`} copyable>
                  {repo.name || repo.kind}: {repo.url || "-"}
                </Typography.Text>
              ))}
            </Space>
          </Card>
          <Card title="依赖">
            <Space size={[6, 6]} wrap>
              {service.dependencies.length === 0 ? <Typography.Text type="secondary">未记录依赖</Typography.Text> : null}
              {service.dependencies.map((item) => (
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
