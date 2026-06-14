import { PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, Space, Spin, Typography, message } from "antd";
import { useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useInventoryQuery, inventoryKeys } from "../queries";
import { removeHost, removePortGroup, saveHost, savePortGroup } from "../api";
import type { GroupForm, Host, HostForm, PortGroup } from "../types";
import { buildStats } from "../utils";
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
  const inventoryQuery = useInventoryQuery();
  const [search, setSearch] = useState("");
  const [saving, setSaving] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<PortGroup | null>(null);
  const [hostDrawerOpen, setHostDrawerOpen] = useState(false);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<Host | null>(null);
  const [editingGroup, setEditingGroup] = useState<PortGroup | null>(null);
  const [hostForm] = Form.useForm<HostForm>();
  const [groupForm] = Form.useForm<GroupForm>();

  const snapshot = inventoryQuery.data;
  const hosts: Host[] = snapshot?.hosts ?? [];
  const groups: PortGroup[] = snapshot?.groups ?? [];
  const meta = snapshot?.meta;

  const filteredGroups = useMemo(() => {
    const keyword = search.trim().toLowerCase();
    if (!keyword) {
      return groups;
    }
    return groups.filter((group) => {
      const text = [
        group.host?.ip,
        group.serviceName,
        group.containerName,
        group.dindHost,
        group.owner,
        group.tags,
        group.components.map((item) => item.name).join(" "),
        group.repositories.map((item) => item.url).join(" "),
      ]
        .join(" ")
        .toLowerCase();
      return text.includes(keyword);
    });
  }, [groups, search]);

  const stats = useMemo(() => buildStats(groups, hosts), [groups, hosts]);

  const groupsByHost = useMemo(
    () =>
      hosts.map((host) => ({
        host,
        groups: filteredGroups
          .filter((group) => group.hostId === host.id)
          .sort((left, right) => left.portStart - right.portStart),
      })),
    [filteredGroups, hosts],
  );

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: inventoryKeys.snapshot });
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
            groups={filteredGroups}
            onDeleteGroup={(group) => void handleDeleteGroup(group)}
            onEditGroup={openEditGroup}
            onSelectGroup={setSelectedGroup}
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
        content = <DependenciesSection groups={filteredGroups} stats={stats} />;
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
              placeholder="搜索服务、IP、组件、仓库"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
            <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={inventoryQuery.isFetching}>
              刷新
            </Button>
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
    </>
  );
}
