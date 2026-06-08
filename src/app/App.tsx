import { useEffect, useMemo, useState } from "react";
import { Alert, Button, Form, Input, Layout, Menu, Space, Spin, Typography, message } from "antd";
import { PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { loadInventory, removeHost, removePortGroup, saveHost, savePortGroup } from "./api";
import { navItems } from "./constants";
import { GroupDetailDrawer } from "./components/GroupDetailDrawer";
import { HostDrawer } from "./components/HostDrawer";
import { PortGroupDrawer } from "./components/PortGroupDrawer";
import { DependenciesSection } from "./sections/DependenciesSection";
import { HostsSection } from "./sections/HostsSection";
import { OverviewSection } from "./sections/OverviewSection";
import { ServicesSection } from "./sections/ServicesSection";
import type { GroupForm, Host, HostForm, Meta, PortGroup, SectionKey } from "./types";

export default function App() {
  const [section, setSection] = useState<SectionKey>("overview");
  const [meta, setMeta] = useState<Meta | null>(null);
  const [hosts, setHosts] = useState<Host[]>([]);
  const [groups, setGroups] = useState<PortGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [selectedGroup, setSelectedGroup] = useState<PortGroup | null>(null);
  const [hostDrawerOpen, setHostDrawerOpen] = useState(false);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<Host | null>(null);
  const [editingGroup, setEditingGroup] = useState<PortGroup | null>(null);
  const [hostForm] = Form.useForm<HostForm>();
  const [groupForm] = Form.useForm<GroupForm>();

  const load = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await loadInventory();
      setMeta(data.meta);
      setHosts(data.hosts);
      setGroups(data.groups);
      if (selectedGroup) {
        setSelectedGroup(data.groups.find((item) => item.id === selectedGroup.id) ?? null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载数据失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const filteredGroups = useMemo(() => {
    const keyword = search.trim().toLowerCase();
    if (!keyword) return groups;
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

  const stats = useMemo(() => {
    const slots = groups.flatMap((group) => group.slots);
    return {
      hosts: hosts.length,
      groups: groups.length,
      usedSlots: slots.filter((slot) => slot.status !== "empty").length,
      emptySlots: slots.filter((slot) => slot.status === "empty").length,
      components: groups.reduce((sum, group) => sum + group.components.length, 0),
      repositories: groups.reduce((sum, group) => sum + group.repositories.length, 0),
    };
  }, [groups, hosts]);

  const groupsByHost = useMemo(() => {
    return hosts.map((host) => ({
      host,
      groups: filteredGroups
        .filter((group) => group.hostId === host.id)
        .sort((left, right) => left.portStart - right.portStart),
    }));
  }, [filteredGroups, hosts]);

  const openCreateHost = () => {
    setEditingHost(null);
    hostForm.resetFields();
    setHostDrawerOpen(true);
  };

  const openEditHost = (host: Host) => {
    setEditingHost(host);
    hostForm.setFieldsValue(host);
    setHostDrawerOpen(true);
  };

  const openCreateGroup = () => {
    setEditingGroup(null);
    groupForm.resetFields();
    groupForm.setFieldsValue({
      hostId: hosts[0]?.id,
      status: "planned",
      slots: [],
      components: [],
      repositories: [],
    } as GroupForm);
    setGroupDrawerOpen(true);
  };

  const openEditGroup = (group: PortGroup) => {
    setEditingGroup(group);
    groupForm.setFieldsValue({
      ...group,
      components: group.components,
      repositories: group.repositories,
      slots: group.slots,
    });
    setGroupDrawerOpen(true);
  };

  const handleSaveHost = async (values: HostForm) => {
    setSaving(true);
    try {
      await saveHost(values, editingHost);
      message.success("主机已保存");
      setHostDrawerOpen(false);
      await load();
    } catch (err) {
      message.error(err instanceof Error ? err.message : "保存失败");
    } finally {
      setSaving(false);
    }
  };

  const handleSaveGroup = async (values: GroupForm) => {
    setSaving(true);
    try {
      await savePortGroup(values, editingGroup);
      message.success("端口组已保存");
      setGroupDrawerOpen(false);
      await load();
    } catch (err) {
      message.error(err instanceof Error ? err.message : "保存失败");
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteGroup = async (group: PortGroup) => {
    await removePortGroup(group);
    message.success("端口组已删除");
    setSelectedGroup(null);
    await load();
  };

  const handleDeleteHost = async (host: Host) => {
    await removeHost(host);
    message.success("主机已删除");
    await load();
  };

  return (
    <Layout className="app-shell">
      <Layout.Sider className="app-sider" width={244}>
        <div className="brand">
          <Typography.Text className="brand-label">Miniport</Typography.Text>
          <Typography.Text className="brand-subtitle">端口服务资产管理</Typography.Text>
        </div>
        <Menu
          className="nav"
          selectedKeys={[section]}
          mode="inline"
          items={navItems}
          onClick={({ key }) => setSection(key as SectionKey)}
        />
        <div className="runtime-card">
          <Typography.Text className="runtime-title">{meta?.version ?? "-"}</Typography.Text>
          <Typography.Text className="runtime-line">{meta?.listen ?? "-"}</Typography.Text>
          <Typography.Text className="runtime-line">{meta?.database ?? "-"}</Typography.Text>
        </div>
      </Layout.Sider>

      <Layout>
        <Layout.Header className="app-header">
          <div>
            <Typography.Title level={4} className="header-title">
              {navItems.find((item) => item.key === section)?.label}
            </Typography.Title>
            <Typography.Text type="secondary">按 IP 和 10 端口组维护服务、容器、依赖和仓库。</Typography.Text>
          </div>
          <Space>
            <Input
              prefix={<SearchOutlined />}
              className="header-search"
              placeholder="搜索服务、IP、组件、仓库"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
            <Button icon={<ReloadOutlined />} onClick={() => void load()} loading={loading}>
              刷新
            </Button>
            <Button icon={<PlusOutlined />} onClick={openCreateHost}>
              主机
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreateGroup} disabled={hosts.length === 0}>
              端口组
            </Button>
          </Space>
        </Layout.Header>

        <Layout.Content className="app-content">
          {error ? <Alert type="error" showIcon message="后端不可用" description={error} /> : null}

          {loading && groups.length === 0 && hosts.length === 0 ? (
            <div className="loading-state">
              <Spin size="large" />
            </div>
          ) : (
            <>
              {section === "overview" ? (
                <OverviewSection
                  stats={stats}
                  hosts={hosts}
                  groupsByHost={groupsByHost}
                  onCreateHost={openCreateHost}
                  onEditHost={openEditHost}
                  onSelectGroup={setSelectedGroup}
                />
              ) : null}

              {section === "services" ? (
                <ServicesSection
                  groups={filteredGroups}
                  onSelectGroup={setSelectedGroup}
                  onEditGroup={openEditGroup}
                  onDeleteGroup={(group) => void handleDeleteGroup(group)}
                />
              ) : null}

              {section === "hosts" ? (
                <HostsSection
                  hosts={hosts}
                  groups={groups}
                  onEditHost={openEditHost}
                  onDeleteHost={(host) => void handleDeleteHost(host)}
                />
              ) : null}

              {section === "dependencies" ? <DependenciesSection groups={groups} stats={stats} /> : null}
            </>
          )}
        </Layout.Content>
      </Layout>

      <HostDrawer
        open={hostDrawerOpen}
        saving={saving}
        editingHost={editingHost}
        form={hostForm}
        onClose={() => setHostDrawerOpen(false)}
        onSave={(values) => void handleSaveHost(values)}
      />

      <PortGroupDrawer
        open={groupDrawerOpen}
        saving={saving}
        hosts={hosts}
        editingGroup={editingGroup}
        form={groupForm}
        onClose={() => setGroupDrawerOpen(false)}
        onSave={(values) => void handleSaveGroup(values)}
      />

      <GroupDetailDrawer
        group={selectedGroup}
        onClose={() => setSelectedGroup(null)}
        onEdit={openEditGroup}
        onDelete={handleDeleteGroup}
      />
    </Layout>
  );
}
