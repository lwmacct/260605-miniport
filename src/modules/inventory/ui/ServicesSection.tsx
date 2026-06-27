import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Card, Popconfirm, Space, Table, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { Key } from "react";
import type { PortGroup } from "../model/inventoryTypes";
import { statusTag } from "../model/inventoryUtils";

type ServicesSectionProps = {
  canManage: boolean;
  groups: PortGroup[];
  onChangeSelection: (ids: number[]) => void;
  onDeleteGroup: (group: PortGroup) => void;
  onEditGroup: (group: PortGroup) => void;
  onSelectGroup: (group: PortGroup) => void;
  selectedRowKeys: number[];
};

export function ServicesSection({
  canManage,
  groups,
  onChangeSelection,
  onDeleteGroup,
  onEditGroup,
  onSelectGroup,
  selectedRowKeys,
}: ServicesSectionProps) {
  const columns: ColumnsType<PortGroup> = [
    {
      title: "分配",
      dataIndex: "name",
      render: (_, group) => <Typography.Link onClick={() => onSelectGroup(group)}>{group.name}</Typography.Link>,
    },
    {
      title: "端口组",
      width: 140,
      render: (_, group) => `${group.portStart}-${group.portEnd}`,
    },
    {
      title: "DIND IP",
      dataIndex: "dindIp",
      width: 150,
      render: (value) => value || "-",
    },
    {
      title: "容器",
      dataIndex: "dindContainer",
      render: (value) => value || "-",
    },
    {
      title: "项目",
      render: (_, group) => (
        <Space size={[4, 4]} wrap>
          {group.projects.slice(0, 3).map((item) => (
            <Tag key={`${group.id}-${item.name}`}>{item.name}</Tag>
          ))}
          {group.projects.length > 3 ? <Tag>+{group.projects.length - 3}</Tag> : null}
        </Space>
      ),
    },
    {
      title: "用户",
      dataIndex: "username",
      width: 120,
    },
    {
      title: "状态",
      width: 110,
      render: (_, group) => statusTag(group.status),
    },
    {
      title: "操作",
      width: 120,
      render: (_, group) =>
        canManage ? (
          <Space>
            <Tooltip title="编辑">
              <Button icon={<EditOutlined />} onClick={() => onEditGroup(group)} />
            </Tooltip>
            <Popconfirm title="删除端口分配" onConfirm={() => onDeleteGroup(group)}>
              <Button danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        ) : null,
    },
  ];

  return (
    <Card>
      <Table<PortGroup>
        rowKey="id"
        columns={columns}
        dataSource={groups}
        pagination={{ pageSize: 12 }}
        rowSelection={
          canManage
            ? {
                selectedRowKeys,
                onChange: (keys: Key[]) => onChangeSelection(keys.map((key) => Number(key))),
              }
            : undefined
        }
        scroll={{ x: 1080 }}
      />
    </Card>
  );
}
