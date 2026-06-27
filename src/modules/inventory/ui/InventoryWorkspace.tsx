import { DeleteOutlined, DownloadOutlined, EditOutlined, PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, Modal, Select, Space, Spin, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useAuthStateQuery } from "@/modules/auth";
import { useInventoryQuery, inventoryKeys } from "../model/inventoryQueries";
import {
  batchDeletePortGroups,
  batchUpdatePortGroups,
  exportPortGroupsURL,
  removePortGroup,
  savePortGroup,
} from "../api/inventoryApi";
import { statusOptions } from "../model/inventoryConstants";
import type { BatchPortGroupUpdate, GroupForm, InventoryQuery, PortGroup } from "../model/inventoryTypes";
import { buildStats } from "../model/inventoryUtils";
import { BatchPortGroupModal, type BatchPortGroupForm } from "./BatchPortGroupModal";
import { OverviewSection } from "./OverviewSection";
import { ServicesSection } from "./ServicesSection";
import { HostsSection } from "./HostsSection";
import { DependenciesSection } from "./DependenciesSection";
import { GroupDetailDrawer } from "./GroupDetailDrawer";
import { PortGroupDrawer } from "./PortGroupDrawer";

type InventoryView = "dependencies" | "hosts" | "overview" | "services";

type InventoryWorkspaceProps = {
  description: string;
  title: string;
  view: InventoryView;
};

export function InventoryWorkspace({
  description,
  title,
  view,
}: InventoryWorkspaceProps) {
  const queryClient = useQueryClient();
  const authState = useAuthStateQuery();
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState<string>();
  const [sort, setSort] = useState<string>("port");
  const [saving, setSaving] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<PortGroup | null>(null);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<PortGroup | null>(null);
  const [selectedGroupIDs, setSelectedGroupIDs] = useState<number[]>([]);
  const [batchModalOpen, setBatchModalOpen] = useState(false);
  const [groupForm] = Form.useForm<GroupForm>();
  const [batchForm] = Form.useForm<BatchPortGroupForm>();

  const query = useMemo<InventoryQuery>(() => {
    return {
      portGroupQuery: search.trim(),
      portGroupSort: sort,
      status,
    };
  }, [search, sort, status]);

  const inventoryQuery = useInventoryQuery(query);

  const snapshot = inventoryQuery.data;
  const groups: PortGroup[] = snapshot?.groups ?? [];
  const meta = snapshot?.meta;
  const canManage = Boolean(authState.data?.session.authenticated);

  useEffect(() => {
    setSelectedGroupIDs([]);
  }, [query]);

  const stats = useMemo(() => buildStats(groups), [groups]);

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: inventoryKeys.snapshot(query) });
  }

  function openCreateGroup() {
    setEditingGroup(null);
    groupForm.resetFields();
    groupForm.setFieldsValue({
      status: "planned",
      slots: [],
      projects: [],
      components: [],
      repositories: [],
    } as GroupForm);
    setGroupDrawerOpen(true);
  }

  function openEditGroup(group: PortGroup) {
    setEditingGroup(group);
    groupForm.setFieldsValue({
      ...group,
      components: group.components,
      projects: group.projects,
      repositories: group.repositories,
      slots: group.slots,
    });
    setGroupDrawerOpen(true);
  }

  async function handleSaveGroup(values: GroupForm) {
    setSaving(true);
    try {
      await savePortGroup(values, editingGroup);
      message.success("端口分配已保存");
      setGroupDrawerOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteGroup(group: PortGroup) {
    await removePortGroup(group);
    message.success("端口分配已删除");
    setSelectedGroup(null);
    await refresh();
  }

  function openBatchModal() {
    batchForm.resetFields();
    setBatchModalOpen(true);
  }

  async function handleBatchUpdate(values: BatchPortGroupForm) {
    const changes: BatchPortGroupUpdate = {};
    if (values.applyStatus) {
      changes.status = values.status;
    }
    if (values.applyOwner) {
      changes.owner = values.owner ?? "";
    }
    if (values.applyTags) {
      changes.tags = values.tags ?? "";
    }

    setSaving(true);
    try {
      await batchUpdatePortGroups(selectedGroupIDs, changes);
      message.success(`已批量更新 ${selectedGroupIDs.length} 个端口分配`);
      setBatchModalOpen(false);
      setSelectedGroupIDs([]);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "批量更新失败");
    } finally {
      setSaving(false);
    }
  }

  function confirmBatchDelete() {
    Modal.confirm({
      title: "批量删除端口分配",
      content: `确认删除已选择的 ${selectedGroupIDs.length} 个端口分配？`,
      okButtonProps: { danger: true },
      onOk: async () => {
        await batchDeletePortGroups(selectedGroupIDs);
        message.success(`已删除 ${selectedGroupIDs.length} 个端口分配`);
        setSelectedGroupIDs([]);
        await refresh();
      },
    });
  }

  const canBatchOperate = canManage && view === "services" && selectedGroupIDs.length > 0;
  const exportURL = exportPortGroupsURL(query);
  const sortOptions = [
    { label: "端口", value: "port" },
    { label: "名称", value: "name" },
    { label: "状态", value: "status" },
    { label: "用户", value: "user" },
    { label: "最近更新", value: "updated_desc" },
  ];

  let content: ReactNode;
  if (inventoryQuery.isPending) {
    content = (
      <div className="loading-state">
        <Spin />
      </div>
    );
  } else if (inventoryQuery.isError) {
    content = <Alert showIcon type="error" message="加载数据失败" description={inventoryQuery.error.message} />;
  } else {
    switch (view) {
      case "services":
        content = (
          <ServicesSection
            canManage={canManage}
            groups={groups}
            onChangeSelection={setSelectedGroupIDs}
            onDeleteGroup={(group) => void handleDeleteGroup(group)}
            onEditGroup={openEditGroup}
            onSelectGroup={setSelectedGroup}
            selectedRowKeys={selectedGroupIDs}
          />
        );
        break;
      case "hosts":
        content = <HostsSection groups={groups} onSelectGroup={setSelectedGroup} />;
        break;
      case "dependencies":
        content = <DependenciesSection groups={groups} stats={stats} />;
        break;
      default:
        content = (
          <OverviewSection
            canManage={canManage}
            groups={groups}
            onCreateGroup={openCreateGroup}
            onSelectGroup={setSelectedGroup}
            stats={stats}
          />
        );
        break;
    }
  }

  return (
    <>
      <section className="app-page-head">
        <Flex align="flex-start" justify="space-between" gap={24} wrap="wrap">
          <div>
            <Typography.Title level={2} className="page-title">
              {title}
            </Typography.Title>
            <Typography.Text type="secondary">{description}</Typography.Text>
            {meta?.version ? (
              <Typography.Paragraph type="secondary" style={{ marginBottom: 0, marginTop: 8 }}>
                {meta.version} · {meta.database}
              </Typography.Paragraph>
            ) : null}
          </div>
          <Space wrap>
            <Input
              prefix={<SearchOutlined />}
              className="page-search"
              placeholder="搜索分配、DIND IP、项目、依赖、仓库"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
            <Select
              allowClear
              placeholder="状态"
              style={{ width: 140 }}
              value={status}
              options={statusOptions}
              onChange={setStatus}
            />
            <Select style={{ width: 150 }} value={sort} options={sortOptions} onChange={setSort} />
            <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={inventoryQuery.isFetching}>
              刷新
            </Button>
            <Button icon={<DownloadOutlined />} href={exportURL} target="_blank">
              导出 CSV
            </Button>
            {canBatchOperate ? (
              <>
                <Typography.Text type="secondary">{`已选 ${selectedGroupIDs.length} 项`}</Typography.Text>
                <Button icon={<EditOutlined />} onClick={openBatchModal}>
                  批量编辑
                </Button>
                <Button danger icon={<DeleteOutlined />} onClick={confirmBatchDelete}>
                  批量删除
                </Button>
              </>
            ) : null}
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={openCreateGroup}
              disabled={!canManage}
            >
              端口分配
            </Button>
          </Space>
        </Flex>
      </section>

      <section className="app-content">{content}</section>

      <GroupDetailDrawer
        canManage={canManage}
        group={selectedGroup}
        onClose={() => setSelectedGroup(null)}
        onDelete={handleDeleteGroup}
        onEdit={(group) => {
          setSelectedGroup(null);
          openEditGroup(group);
        }}
      />
      <PortGroupDrawer
        editingGroup={editingGroup}
        form={groupForm}
        onClose={() => setGroupDrawerOpen(false)}
        onSave={(values) => void handleSaveGroup(values)}
        open={groupDrawerOpen}
        saving={saving}
      />
      <BatchPortGroupModal
        form={batchForm}
        onCancel={() => setBatchModalOpen(false)}
        onSubmit={(values) => void handleBatchUpdate(values)}
        open={batchModalOpen}
        selectedCount={selectedGroupIDs.length}
        saving={saving}
      />
    </>
  );
}
