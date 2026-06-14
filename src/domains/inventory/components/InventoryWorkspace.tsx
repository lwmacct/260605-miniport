import { DeleteOutlined, DownloadOutlined, EditOutlined, PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, Modal, Select, Space, Spin, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useInventoryQuery, inventoryKeys } from "../queries";
import {
  batchDeletePortGroups,
  batchUpdatePortGroups,
  exportPortGroupsURL,
  removeHost,
  removePortGroup,
  saveHost,
  savePortGroup,
} from "../api";
import { statusOptions } from "../constants";
import type { BatchPortGroupUpdate, GroupForm, Host, HostForm, InventoryQuery, PortGroup } from "../types";
import { buildStats } from "../utils";
import { BatchPortGroupModal, type BatchPortGroupForm } from "./BatchPortGroupModal";
import { OverviewSection } from "./OverviewSection";
import { ServicesSection } from "./ServicesSection";
import { HostsSection } from "./HostsSection";
import { DependenciesSection } from "./DependenciesSection";
import { GroupDetailDrawer } from "./GroupDetailDrawer";
import { HostDrawer } from "./HostDrawer";
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
  const [search, setSearch] = useState("");
  const [environment, setEnvironment] = useState<string>();
  const [status, setStatus] = useState<string>();
  const [sort, setSort] = useState<string>(view === "hosts" ? "ip" : "host_port");
  const [saving, setSaving] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<PortGroup | null>(null);
  const [hostDrawerOpen, setHostDrawerOpen] = useState(false);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<Host | null>(null);
  const [editingGroup, setEditingGroup] = useState<PortGroup | null>(null);
  const [selectedGroupIDs, setSelectedGroupIDs] = useState<number[]>([]);
  const [batchModalOpen, setBatchModalOpen] = useState(false);
  const [hostForm] = Form.useForm<HostForm>();
  const [groupForm] = Form.useForm<GroupForm>();
  const [batchForm] = Form.useForm<BatchPortGroupForm>();

  const query = useMemo<InventoryQuery>(() => {
    if (view === "hosts") {
      return {
        environment,
        hostQuery: search.trim(),
        hostSort: sort,
      };
    }

    return {
      hostQuery: "",
      portGroupQuery: search.trim(),
      portGroupSort: sort,
      status,
    };
  }, [environment, search, sort, status, view]);

  const inventoryQuery = useInventoryQuery(query);

  const snapshot = inventoryQuery.data;
  const hosts: Host[] = snapshot?.hosts ?? [];
  const groups: PortGroup[] = snapshot?.groups ?? [];
  const meta = snapshot?.meta;

  useEffect(() => {
    setSelectedGroupIDs([]);
  }, [query]);

  const stats = useMemo(() => buildStats(groups, hosts), [groups, hosts]);

  const groupsByHost = useMemo(
    () =>
      hosts.map((host) => ({
        host,
        groups: groups
          .filter((group) => group.hostId === host.id)
          .sort((left, right) => left.portStart - right.portStart),
      })),
    [groups, hosts],
  );

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: inventoryKeys.snapshot(query) });
  }

  function openCreateHost() {
    setEditingHost(null);
    hostForm.resetFields();
    setHostDrawerOpen(true);
  }

  function openEditHost(host: Host) {
    setEditingHost(host);
    hostForm.setFieldsValue(host);
    setHostDrawerOpen(true);
  }

  function openCreateGroup() {
    const firstHost = hosts[0];
    setEditingGroup(null);
    groupForm.resetFields();
    groupForm.setFieldsValue({
      hostId: firstHost?.id,
      status: "planned",
      slots: [],
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
      repositories: group.repositories,
      slots: group.slots,
    });
    setGroupDrawerOpen(true);
  }

  async function handleSaveHost(values: HostForm) {
    setSaving(true);
    try {
      await saveHost(values, editingHost);
      message.success("主机已保存");
      setHostDrawerOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleSaveGroup(values: GroupForm) {
    setSaving(true);
    try {
      await savePortGroup(values, editingGroup);
      message.success("端口组已保存");
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
    message.success("端口组已删除");
    setSelectedGroup(null);
    await refresh();
  }

  async function handleDeleteHost(host: Host) {
    await removeHost(host);
    message.success("主机已删除");
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
      message.success(`已批量更新 ${selectedGroupIDs.length} 个端口组`);
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
      title: "批量删除端口组",
      content: `确认删除已选择的 ${selectedGroupIDs.length} 个端口组？`,
      okButtonProps: { danger: true },
      onOk: async () => {
        await batchDeletePortGroups(selectedGroupIDs);
        message.success(`已删除 ${selectedGroupIDs.length} 个端口组`);
        setSelectedGroupIDs([]);
        await refresh();
      },
    });
  }

  const canBatchOperate = view === "services" && selectedGroupIDs.length > 0;
  const canExport = view !== "hosts";
  const exportURL = exportPortGroupsURL(query);
  const environmentOptions = useMemo(
    () =>
      Array.from(new Set((snapshot?.hosts ?? []).map((host) => host.environment).filter(Boolean))).map((value) => ({
        label: value,
        value,
      })),
    [snapshot?.hosts],
  );

  const sortOptions =
    view === "hosts"
      ? [
          { label: "IP", value: "ip" },
          { label: "环境", value: "environment" },
          { label: "最近更新", value: "updated_desc" },
        ]
      : [
          { label: "主机 + 端口", value: "host_port" },
          { label: "服务名", value: "service" },
          { label: "状态", value: "status" },
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
        content = (
          <HostsSection
            groups={groups}
            hosts={hosts}
            onDeleteHost={(host) => void handleDeleteHost(host)}
            onEditHost={openEditHost}
          />
        );
        break;
      case "dependencies":
        content = <DependenciesSection groups={groups} stats={stats} />;
        break;
      default:
        content = (
          <OverviewSection
            groupsByHost={groupsByHost}
            hosts={hosts}
            onCreateHost={openCreateHost}
            onEditHost={openEditHost}
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
              placeholder={view === "hosts" ? "搜索 IP、名称、网段、备注" : "搜索服务、IP、组件、仓库"}
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
            {view === "hosts" ? (
              <Select
                allowClear
                placeholder="环境"
                style={{ width: 140 }}
                value={environment}
                options={environmentOptions}
                onChange={setEnvironment}
              />
            ) : (
              <Select
                allowClear
                placeholder="状态"
                style={{ width: 140 }}
                value={status}
                options={statusOptions}
                onChange={setStatus}
              />
            )}
            <Select style={{ width: 150 }} value={sort} options={sortOptions} onChange={setSort} />
            <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={inventoryQuery.isFetching}>
              刷新
            </Button>
            {canExport ? (
              <Button icon={<DownloadOutlined />} href={exportURL} target="_blank">
                导出 CSV
              </Button>
            ) : null}
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
            <Button icon={<PlusOutlined />} onClick={openCreateHost}>
              主机
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreateGroup} disabled={hosts.length === 0}>
              端口组
            </Button>
          </Space>
        </Flex>
      </section>

      <section className="app-content">{content}</section>

      <GroupDetailDrawer
        group={selectedGroup}
        onClose={() => setSelectedGroup(null)}
        onDelete={handleDeleteGroup}
        onEdit={(group) => {
          setSelectedGroup(null);
          openEditGroup(group);
        }}
      />
      <HostDrawer
        editingHost={editingHost}
        form={hostForm}
        onClose={() => setHostDrawerOpen(false)}
        onSave={(values) => void handleSaveHost(values)}
        open={hostDrawerOpen}
        saving={saving}
      />
      <PortGroupDrawer
        editingGroup={editingGroup}
        form={groupForm}
        hosts={hosts}
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
