import { useEffect, useMemo, useState } from "react";
import { Alert, Button, ConfigProvider, Flex, Form, Input, Layout, Space, Spin, Switch, Tooltip, Typography, message } from "antd";
import { MoonOutlined, PlusOutlined, ReloadOutlined, SearchOutlined, SunOutlined } from "@ant-design/icons";
import { loadInventory, removeHost, removePortGroup, saveHost, savePortGroup } from "./api";
import { navItems } from "./constants";
import { AppHeader } from "./components/AppHeader";
import { GroupDetailDrawer } from "./components/GroupDetailDrawer";
import { HostDrawer } from "./components/HostDrawer";
import { PortGroupDrawer } from "./components/PortGroupDrawer";
import { DependenciesSection } from "./sections/DependenciesSection";
import { HostsSection } from "./sections/HostsSection";
import { OverviewSection } from "./sections/OverviewSection";
import { ServicesSection } from "./sections/ServicesSection";
import { navigateSectionHash, readSectionFromHash, replaceSectionHash } from "./routing";
import { appTheme, applyTheme, readInitialTheme } from "./theme";
import type { GroupForm, Host, HostForm, Meta, PortGroup, SectionKey } from "./types";
import type { ThemeMode } from "./theme";

export default function App() {
  const [section, setSection] = useState<SectionKey>(() => readSectionFromHash());
  const [themeMode, setThemeMode] = useState<ThemeMode>(() => readInitialTheme());
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
  const themeConfig = useMemo(() => appTheme(themeMode), [themeMode]);

  const load = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await loadInventory();
      setMeta(data.meta);
      setHosts(data.hosts ?? []);
      setGroups(data.groups ?? []);
      if (selectedGroup) {
        setSelectedGroup((data.groups ?? []).find((item) => item.id === selectedGroup.id) ?? null);
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

  useEffect(() => {
    replaceSectionHash(readSectionFromHash());

    const handleHashChange = () => {
      const nextSection = readSectionFromHash();
      setSection(nextSection);
      replaceSectionHash(nextSection);
    };

    window.addEventListener("hashchange", handleHashChange);
    return () => window.removeEventListener("hashchange", handleHashChange);
  }, []);

  const filteredGroups = useMemo(() => {
    const keyword = search.trim().toLowerCase();
    const safeGroups = groups ?? [];
    if (!keyword) return safeGroups;
    return safeGroups.filter((group) => {
      const components = group.components ?? [];
      const repositories = group.repositories ?? [];
      const text = [
        group.host?.ip,
        group.serviceName,
        group.containerName,
        group.dindHost,
        group.owner,
        group.tags,
        components.map((item) => item.name).join(" "),
        repositories.map((item) => item.url).join(" "),
      ]
        .join(" ")
        .toLowerCase();
      return text.includes(keyword);
    });
  }, [groups, search]);

  const stats = useMemo(() => {
    const safeHosts = hosts ?? [];
    const safeGroups = groups ?? [];
    const slots = safeGroups.flatMap((group) => group.slots ?? []);
    return {
      hosts: safeHosts.length,
      groups: safeGroups.length,
      usedSlots: slots.filter((slot) => slot.status !== "empty").length,
      emptySlots: slots.filter((slot) => slot.status === "empty").length,
      components: safeGroups.reduce((sum, group) => sum + (group.components ?? []).length, 0),
      repositories: safeGroups.reduce((sum, group) => sum + (group.repositories ?? []).length, 0),
    };
  }, [groups, hosts]);

  const groupsByHost = useMemo(() => {
    const safeHosts = hosts ?? [];
    return safeHosts.map((host) => ({
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
    const firstHost = (hosts ?? [])[0];
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

  const handleNavigate = (key: string) => {
    const nextSection = key as SectionKey;
    setSection(nextSection);
    navigateSectionHash(nextSection);
  };

  const currentNavItem = navItems.find((item) => item.key === section);
  const visibleHosts = hosts ?? [];
  const visibleGroups = groups ?? [];

  return (
    <ConfigProvider theme={themeConfig}>
      <Layout className="app-shell">
        <AppHeader
          activeKey={section}
          actions={
            <Tooltip title={themeMode === "dark" ? "切换到明亮模式" : "切换到暗色模式"}>
              <Switch
                checked={themeMode === "dark"}
                checkedChildren={<MoonOutlined />}
                unCheckedChildren={<SunOutlined />}
                onChange={(checked) => {
                  const nextTheme = checked ? "dark" : "light";
                  applyTheme(nextTheme);
                  setThemeMode(nextTheme);
                }}
              />
            </Tooltip>
          }
          brandName="Miniport"
          navItems={navItems}
          onNavigate={handleNavigate}
          version={meta?.version}
        />

        <Layout.Content className="app-page">
          <section className="app-page-head">
            <Flex align="flex-start" justify="space-between" gap={24} wrap="wrap">
              <div>
                <Typography.Title level={2} className="page-title">
                  {currentNavItem?.label}
                </Typography.Title>
                <Typography.Text type="secondary">按 IP 和 10 端口组维护服务、容器、依赖和仓库。</Typography.Text>
              </div>
              <Space wrap>
                <Input
                  prefix={<SearchOutlined />}
                  className="page-search"
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
                <Button type="primary" icon={<PlusOutlined />} onClick={openCreateGroup} disabled={visibleHosts.length === 0}>
                  端口组
                </Button>
              </Space>
            </Flex>
          </section>

          <div className="app-content">
            {error ? <Alert type="error" showIcon message="后端不可用" description={error} /> : null}

            {loading && visibleGroups.length === 0 && visibleHosts.length === 0 ? (
              <div className="loading-state">
                <Spin size="large" />
              </div>
            ) : (
              <>
                {section === "overview" ? (
                  <OverviewSection
                    stats={stats}
                    hosts={visibleHosts}
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
                    hosts={visibleHosts}
                    groups={visibleGroups}
                    onEditHost={openEditHost}
                    onDeleteHost={(host) => void handleDeleteHost(host)}
                  />
                ) : null}

                {section === "dependencies" ? <DependenciesSection groups={visibleGroups} stats={stats} /> : null}
              </>
            )}
          </div>
        </Layout.Content>

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
          hosts={visibleHosts}
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
    </ConfigProvider>
  );
}
