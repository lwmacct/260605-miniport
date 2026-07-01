import { DownloadOutlined, PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, InputNumber, Modal, Select, Space, Spin, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useAuthStateQuery } from "@/modules/auth";
import { usePortsvcQuery, portsvcKeys } from "../model/portsvcQueries";
import {
  exportPortGroupsURL,
  removeHost,
  removePortGroup,
  saveHost,
  savePortGroup,
} from "../api/portsvcApi";
import { statusOptions } from "../model/portsvcConstants";
import type { HostForm, HostItem, PortGroupForm, PortGroupItem, PortsvcQuery } from "../model/portsvcTypes";
import { buildStats } from "../model/portsvcUtils";
import { OverviewSection } from "./OverviewSection";
import { ServicesSection } from "./ServicesSection";
import { ProjectServicesSection } from "./ProjectServicesSection";
import { DependenciesSection } from "./DependenciesSection";
import { ServiceDetailDrawer } from "./ServiceDetailDrawer";
import { ServiceDrawer } from "./ServiceDrawer";

type PortsvcView = "dependencies" | "projects" | "overview" | "services";

type PortsvcWorkspaceProps = {
  description: string;
  title: string;
  view: PortsvcView;
};

export function PortsvcWorkspace({ description, title, view }: PortsvcWorkspaceProps) {
  const queryClient = useQueryClient();
  const authState = useAuthStateQuery();
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState<string>();
  const [sort, setSort] = useState<string>("port");
  const [saving, setSaving] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<PortGroupItem | null>(null);
  const [editingGroup, setEditingGroup] = useState<PortGroupItem | null>(null);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [hostModalOpen, setHostModalOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<HostItem | null>(null);
  const [groupForm] = Form.useForm<PortGroupForm>();
  const [hostForm] = Form.useForm<HostForm>();

  const query = useMemo<PortsvcQuery>(() => {
    return {
      query: search.trim(),
      sort,
      status,
    };
  }, [search, sort, status]);

  const portsvcQuery = usePortsvcQuery(query);
  const snapshot = portsvcQuery.data;
  const groups: PortGroupItem[] = snapshot?.portGroups ?? [];
  const hosts: HostItem[] = snapshot?.hosts ?? [];
  const meta = snapshot?.meta;
  const canManage = Boolean(authState.data?.session.authenticated);
  const stats = useMemo(() => buildStats(groups, hosts), [groups, hosts]);

  useEffect(() => {
    setSelectedGroup(null);
  }, [query]);

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: portsvcKeys.snapshot(query) });
  }

  function openCreateGroup() {
    setEditingGroup(null);
    groupForm.resetFields();
    groupForm.setFieldsValue({
      runtimeMode: "dind",
      status: "available",
      slots: [],
      repositories: [],
      dependencies: [],
    });
    setGroupDrawerOpen(true);
  }

  function openEditGroup(group: PortGroupItem) {
    setEditingGroup(group);
    groupForm.setFieldsValue({
      ...group,
      dependencies: group.dependencies,
      repositories: group.repositories,
      slots: group.slots,
    });
    setGroupDrawerOpen(true);
  }

  async function handleSaveGroup(values: PortGroupForm) {
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

  async function handleDeleteGroup(group: PortGroupItem) {
    await removePortGroup(group);
    message.success("端口组已删除");
    setSelectedGroup(null);
    await refresh();
  }

  function openCreateHost() {
    setEditingHost(null);
    hostForm.resetFields();
    hostForm.setFieldsValue({ status: "active" });
    setHostModalOpen(true);
  }

  function openEditHost(host: HostItem) {
    setEditingHost(host);
    hostForm.setFieldsValue(host);
    setHostModalOpen(true);
  }

  async function handleSaveHost(values: HostForm) {
    setSaving(true);
    try {
      await saveHost(values, editingHost);
      message.success("宿主机已保存");
      setHostModalOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteHost(host: HostItem) {
    await removeHost(host);
    message.success("宿主机已删除");
    setHostModalOpen(false);
    await refresh();
  }

  const exportURL = exportPortGroupsURL(query);
  const sortOptions = [
    { label: "端口", value: "port" },
    { label: "项目", value: "project" },
    { label: "状态", value: "status" },
    { label: "最近更新", value: "updated_desc" },
  ];

  let content: ReactNode;
  if (portsvcQuery.isPending) {
    content = (
      <div className="loading-state">
        <Spin />
      </div>
    );
  } else if (portsvcQuery.isError) {
    content = <Alert showIcon type="error" message="加载数据失败" description={portsvcQuery.error.message} />;
  } else {
    switch (view) {
      case "services":
        content = (
          <ServicesSection
            canManage={canManage}
            groups={groups}
            onDeleteGroup={(group) => void handleDeleteGroup(group)}
            onEditGroup={openEditGroup}
            onSelectGroup={setSelectedGroup}
          />
        );
        break;
      case "projects":
        content = <ProjectServicesSection groups={groups} onSelectGroup={setSelectedGroup} />;
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
            onEditGroup={openEditGroup}
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
              placeholder="搜索项目、服务 IP、宿主机、仓库"
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
            <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={portsvcQuery.isFetching}>
              刷新
            </Button>
            <Button icon={<DownloadOutlined />} href={exportURL} target="_blank">
              导出 CSV
            </Button>
            <Button icon={<PlusOutlined />} onClick={openCreateHost} disabled={!canManage}>
              宿主机
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreateGroup} disabled={!canManage}>
              端口组
            </Button>
          </Space>
        </Flex>
      </section>

      <section className="app-content">{content}</section>

      <ServiceDetailDrawer
        canManage={canManage}
        group={selectedGroup}
        onClose={() => setSelectedGroup(null)}
        onDelete={handleDeleteGroup}
        onEdit={(group) => {
          setSelectedGroup(null);
          openEditGroup(group);
        }}
      />
      <ServiceDrawer
        editingGroup={editingGroup}
        form={groupForm}
        hosts={hosts}
        onClose={() => setGroupDrawerOpen(false)}
        onSave={(values) => void handleSaveGroup(values)}
        open={groupDrawerOpen}
        saving={saving}
      />
      <Modal
        title={editingHost ? "编辑宿主机" : "新建宿主机"}
        open={hostModalOpen}
        onCancel={() => setHostModalOpen(false)}
        onOk={() => hostForm.submit()}
        confirmLoading={saving}
      >
        <Form form={hostForm} layout="vertical" onFinish={(values) => void handleSaveHost(values)}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: "请填写宿主机名称" }]}>
            <Input placeholder="miniport-host-01" />
          </Form.Item>
          <Form.Item name="ip" label="IP">
            <Input placeholder="10.0.0.12" />
          </Form.Item>
          <Form.Item name="spec" label="规格">
            <Input placeholder="4h4g" />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select options={[{ value: "active", label: "运行中" }, { value: "stopped", label: "停用" }]} />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          {editingHost ? (
            <Space>
              <Button danger onClick={() => void handleDeleteHost(editingHost)}>
                删除宿主机
              </Button>
              {hosts.length > 1 ? (
                <Button onClick={() => openEditHost(hosts[(hosts.findIndex((item) => item.id === editingHost.id) + 1) % hosts.length])}>
                  下一个
                </Button>
              ) : null}
            </Space>
          ) : null}
        </Form>
      </Modal>
    </>
  );
}
