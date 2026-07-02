import { DownloadOutlined, PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, Modal, Select, Space, Spin, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useAuthStateQuery } from "@/modules/auth";
import { usePortsvcQuery } from "../model/portsvcQueries";
import {
  exportPortGroupsURL,
  removeDependencyAsset,
  removeHost,
  removePortGroup,
  removeServiceGroup,
  saveDependencyAsset,
  saveHost,
  savePortGroup,
  saveServiceGroup,
} from "../api/portsvcApi";
import {
  assetKindOptions,
  assetProviderOptions,
  assetTypeOptions,
  controllabilityOptions,
  serviceGroupKindOptions,
  serviceGroupStatusOptions,
  statusOptions,
  visibilityOptions,
} from "../model/portsvcConstants";
import type {
  DependencyAssetItem,
  HostForm,
  HostItem,
  PortGroupForm,
  PortGroupItem,
  PortsvcQuery,
  ServiceGroupForm,
  ServiceGroupItem,
} from "../model/portsvcTypes";
import { buildStats } from "../model/portsvcUtils";
import { OverviewSection } from "./OverviewSection";
import { ProjectServicesSection } from "./ProjectServicesSection";
import { DependenciesSection } from "./DependenciesSection";
import { HostsSection } from "./HostsSection";
import { PortGroupsUsageSection } from "./PortGroupsUsageSection";
import { ServiceGroupsSection } from "./ServiceGroupsSection";
import { ServiceGroupDetailDrawer } from "./ServiceGroupDetailDrawer";
import { ServiceDetailDrawer } from "./ServiceDetailDrawer";
import { ServiceDrawer } from "./ServiceDrawer";

type PortsvcView = "dependencies" | "hosts" | "overview" | "portGroups" | "projects" | "serviceGroups";

const portSegmentOptions = [1000, 2000, 3000, 4000, 5000] as const;

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
  const [portSegment, setPortSegment] = useState<number>();
  const [saving, setSaving] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<PortGroupItem | null>(null);
  const [editingGroup, setEditingGroup] = useState<PortGroupItem | null>(null);
  const [groupDrawerOpen, setGroupDrawerOpen] = useState(false);
  const [hostModalOpen, setHostModalOpen] = useState(false);
  const [editingHost, setEditingHost] = useState<HostItem | null>(null);
  const [assetModalOpen, setAssetModalOpen] = useState(false);
  const [editingAsset, setEditingAsset] = useState<DependencyAssetItem | null>(null);
  const [selectedServiceGroup, setSelectedServiceGroup] = useState<ServiceGroupItem | null>(null);
  const [serviceGroupModalOpen, setServiceGroupModalOpen] = useState(false);
  const [editingServiceGroup, setEditingServiceGroup] = useState<ServiceGroupItem | null>(null);
  const [groupForm] = Form.useForm<PortGroupForm>();
  const [hostForm] = Form.useForm<HostForm>();
  const [assetForm] = Form.useForm<Partial<DependencyAssetItem>>();
  const [serviceGroupForm] = Form.useForm<ServiceGroupForm>();

  const query = useMemo<PortsvcQuery>(() => {
    if (view === "hosts" || view === "overview" || view === "portGroups" || view === "serviceGroups") {
      return { sort: "port" };
    }

    return {
      query: search.trim(),
      sort,
      status,
    };
  }, [search, sort, status, view]);

  const portsvcQuery = usePortsvcQuery(query);
  const snapshot = portsvcQuery.data;
  const groups: PortGroupItem[] = snapshot?.portGroups ?? [];
  const hosts: HostItem[] = snapshot?.hosts ?? [];
  const dependencyAssets: DependencyAssetItem[] = snapshot?.dependencyAssets ?? [];
  const serviceGroups: ServiceGroupItem[] = snapshot?.serviceGroups ?? [];
  const meta = snapshot?.meta;
  const canManage = Boolean(authState.data?.session.authenticated);
  const stats = useMemo(
    () => buildStats(groups, hosts, dependencyAssets, serviceGroups),
    [dependencyAssets, groups, hosts, serviceGroups],
  );
  const projectGroups = useMemo(() => {
    if (view !== "projects" || !portSegment) {
      return groups;
    }

    const segmentEnd = portSegment + 999;
    return groups.filter((group) => group.portPrefix >= portSegment && group.portPrefix <= segmentEnd);
  }, [groups, portSegment, view]);
  const filteredServiceGroups = useMemo(() => {
    if (view !== "serviceGroups") {
      return serviceGroups;
    }
    const keyword = search.trim().toLowerCase();
    return serviceGroups.filter((group) => {
      if (status && group.status !== status) {
        return false;
      }
      if (!keyword) {
        return true;
      }
      return [group.name, group.kind, group.description, group.notes]
        .some((value) => value.toLowerCase().includes(keyword));
    });
  }, [search, serviceGroups, status, view]);

  useEffect(() => {
    setSelectedGroup(null);
    setSelectedServiceGroup(null);
  }, [query]);

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: ["portsvc", "snapshot"] });
  }

  function openCreateGroup() {
    setEditingGroup(null);
    groupForm.resetFields();
    groupForm.setFieldsValue({
      runtimeMode: "dind",
      status: "available",
      slots: [],
      assetLinks: [],
    });
    setGroupDrawerOpen(true);
  }

  function openEditGroup(group: PortGroupItem) {
    setEditingGroup(group);
    groupForm.setFieldsValue({
      ...group,
      assetLinks: group.assetLinks,
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

  function openCreateAsset() {
    setEditingAsset(null);
    assetForm.resetFields();
    assetForm.setFieldsValue({
      assetKind: "component",
      assetType: "middleware",
      controllability: "unknown",
      provider: "manual",
      status: "active",
      visibility: "unknown",
    });
    setAssetModalOpen(true);
  }

  function openEditAsset(asset: DependencyAssetItem) {
    setEditingAsset(asset);
    assetForm.setFieldsValue(asset);
    setAssetModalOpen(true);
  }

  async function handleSaveAsset(values: Partial<DependencyAssetItem>) {
    setSaving(true);
    try {
      await saveDependencyAsset(values, editingAsset);
      message.success("依赖资产已保存");
      setAssetModalOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteAsset(asset: DependencyAssetItem) {
    await removeDependencyAsset(asset);
    message.success("依赖资产已删除");
    setAssetModalOpen(false);
    await refresh();
  }

  function openCreateServiceGroup() {
    setEditingServiceGroup(null);
    serviceGroupForm.resetFields();
    serviceGroupForm.setFieldsValue({
      kind: "service",
      portGroups: [],
      status: "active",
    });
    setServiceGroupModalOpen(true);
  }

  function openEditServiceGroup(group: ServiceGroupItem) {
    setEditingServiceGroup(group);
    serviceGroupForm.setFieldsValue({
      ...group,
      portGroups: group.portGroups,
    });
    setServiceGroupModalOpen(true);
  }

  async function handleSaveServiceGroup(values: ServiceGroupForm) {
    setSaving(true);
    try {
      await saveServiceGroup(values, editingServiceGroup);
      message.success("服务组已保存");
      setServiceGroupModalOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteServiceGroup(group: ServiceGroupItem) {
    await removeServiceGroup(group);
    message.success("服务组已删除");
    setSelectedServiceGroup(null);
    setServiceGroupModalOpen(false);
    await refresh();
  }

  const exportURL = exportPortGroupsURL(query);
  const showProjectFilters = view === "projects" || view === "dependencies" || view === "serviceGroups";
  const sortOptions = [
    { label: "端口", value: "port" },
    { label: "运行环境", value: "environment" },
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
      case "hosts":
        content = (
          <HostsSection
            canManage={canManage}
            groups={groups}
            hosts={hosts}
            onEditHost={openEditHost}
          />
        );
        break;
      case "portGroups":
        content = <PortGroupsUsageSection canManage={canManage} groups={groups} onCreateGroup={openCreateGroup} />;
        break;
      case "projects":
        content = (
          <ProjectServicesSection
            canManage={canManage}
            groups={projectGroups}
            onDeleteGroup={(group) => void handleDeleteGroup(group)}
            onEditGroup={openEditGroup}
            onSelectGroup={setSelectedGroup}
            serviceGroups={serviceGroups}
          />
        );
        break;
      case "serviceGroups":
        content = (
          <ServiceGroupsSection
            canManage={canManage}
            onEditServiceGroup={openEditServiceGroup}
            onSelectServiceGroup={setSelectedServiceGroup}
            serviceGroups={filteredServiceGroups}
          />
        );
        break;
      case "dependencies":
        content = (
          <DependenciesSection
            canManage={canManage}
            dependencyAssets={dependencyAssets}
            groups={groups}
            onEditAsset={openEditAsset}
            stats={stats}
          />
        );
        break;
      default:
        content = (
          <OverviewSection
            canManage={canManage}
            onCreateGroup={openCreateGroup}
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
            {showProjectFilters ? (
              <>
                <Input
                  prefix={<SearchOutlined />}
                  className="page-search"
                  placeholder={view === "serviceGroups" ? "搜索服务组" : "搜索运行环境、服务 IP、宿主机、仓库"}
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                />
                <Select
                  allowClear
                  placeholder="状态"
                  style={{ width: 140 }}
                  value={status}
                  options={view === "serviceGroups" ? serviceGroupStatusOptions : statusOptions}
                  onChange={setStatus}
                />
                {view !== "serviceGroups" ? (
                  <Select style={{ width: 150 }} value={sort} options={sortOptions} onChange={setSort} />
                ) : null}
                {view === "projects" ? (
                  <Space.Compact>
                    {portSegmentOptions.map((segment) => (
                      <Button
                        key={segment}
                        type={portSegment === segment ? "primary" : "default"}
                        onClick={() => setPortSegment((current) => (current === segment ? undefined : segment))}
                      >
                        {segment}
                      </Button>
                    ))}
                  </Space.Compact>
                ) : null}
              </>
            ) : null}
            <Button icon={<ReloadOutlined />} onClick={() => void refresh()} loading={portsvcQuery.isFetching}>
              刷新
            </Button>
            <Button icon={<DownloadOutlined />} href={exportURL} target="_blank">
              导出 CSV
            </Button>
            {view === "hosts" ? (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreateHost} disabled={!canManage}>
                宿主机
              </Button>
            ) : null}
            {view === "dependencies" ? (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreateAsset} disabled={!canManage}>
                依赖资产
              </Button>
            ) : null}
            {view === "projects" ? (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreateGroup} disabled={!canManage}>
                运行环境
              </Button>
            ) : null}
            {view === "serviceGroups" ? (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreateServiceGroup} disabled={!canManage}>
                服务组
              </Button>
            ) : null}
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
        serviceGroups={serviceGroups}
      />
      <ServiceGroupDetailDrawer
        canManage={canManage}
        groups={groups}
        serviceGroup={selectedServiceGroup}
        onClose={() => setSelectedServiceGroup(null)}
        onDelete={handleDeleteServiceGroup}
        onEdit={(group) => {
          setSelectedServiceGroup(null);
          openEditServiceGroup(group);
        }}
      />
      <ServiceDrawer
        dependencyAssets={dependencyAssets}
        editingGroup={editingGroup}
        form={groupForm}
        hosts={hosts}
        onClose={() => setGroupDrawerOpen(false)}
        onSave={(values) => void handleSaveGroup(values)}
        open={groupDrawerOpen}
        saving={saving}
      />
      <Modal
        title={editingAsset ? "编辑依赖资产" : "新建依赖资产"}
        open={assetModalOpen}
        onCancel={() => setAssetModalOpen(false)}
        onOk={() => assetForm.submit()}
        confirmLoading={saving}
      >
        <Form form={assetForm} layout="vertical" onFinish={(values) => void handleSaveAsset(values)}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: "请填写资产名称" }]}>
            <Input placeholder="miniport / Redis / 闭源支付服务" />
          </Form.Item>
          <Space.Compact block>
            <Form.Item name="assetKind" label="资产类别" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={assetKindOptions} />
            </Form.Item>
            <Form.Item name="assetType" label="资产类型" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={assetTypeOptions} />
            </Form.Item>
          </Space.Compact>
          <Space.Compact block>
            <Form.Item name="provider" label="Provider" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={assetProviderOptions} />
            </Form.Item>
            <Form.Item name="visibility" label="可见性" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={visibilityOptions} />
            </Form.Item>
          </Space.Compact>
          <Form.Item name="controllability" label="可控性" rules={[{ required: true }]}>
            <Select options={controllabilityOptions} />
          </Form.Item>
          <Form.Item name="url" label="URL">
            <Input placeholder="Git URL / 服务地址 / 文档地址" />
          </Form.Item>
          <Form.Item name="fullName" label="仓库全名">
            <Input placeholder="org/repo" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          {editingAsset ? (
            <Button danger onClick={() => void handleDeleteAsset(editingAsset)}>
              删除依赖资产
            </Button>
          ) : null}
        </Form>
      </Modal>
      <Modal
        title={editingServiceGroup ? "编辑服务组" : "新建服务组"}
        open={serviceGroupModalOpen}
        onCancel={() => setServiceGroupModalOpen(false)}
        onOk={() => serviceGroupForm.submit()}
        confirmLoading={saving}
        width={860}
      >
        <Form form={serviceGroupForm} layout="vertical" onFinish={(values) => void handleSaveServiceGroup(values)}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: "请填写服务组名称" }]}>
            <Input placeholder="etcd 集群 / kafka 集群" />
          </Form.Item>
          <Space.Compact block>
            <Form.Item name="kind" label="类型" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={serviceGroupKindOptions} />
            </Form.Item>
            <Form.Item name="status" label="状态" rules={[{ required: true }]} style={{ width: "50%" }}>
              <Select options={serviceGroupStatusOptions} />
            </Form.Item>
          </Space.Compact>
          <Form.Item name="description" label="说明">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.List name="portGroups">
            {(fields, { add, remove }) => (
              <Space direction="vertical" className="content-stack">
                <Button icon={<PlusOutlined />} onClick={() => add({ role: "" })}>
                  添加运行环境
                </Button>
                {fields.map((field) => (
                  <Space.Compact key={field.key} block>
                    <Form.Item name={[field.name, "portGroupId"]} rules={[{ required: true }]} style={{ width: "48%" }}>
                      <Select
                        showSearch
                        placeholder="选择运行环境"
                        optionFilterProp="label"
                        options={groups.map((group) => ({
                          value: group.id,
                          label: `${group.portPrefix} · ${group.environmentName || group.serviceIp || group.runtimeMode}`,
                        }))}
                      />
                    </Form.Item>
                    <Form.Item name={[field.name, "role"]} style={{ width: "22%" }}>
                      <Input placeholder="角色" />
                    </Form.Item>
                    <Form.Item name={[field.name, "notes"]} style={{ width: "22%" }}>
                      <Input placeholder="备注" />
                    </Form.Item>
                    <Button danger onClick={() => remove(field.name)}>
                      删除
                    </Button>
                  </Space.Compact>
                ))}
              </Space>
            )}
          </Form.List>
          <Form.Item name="notes" label="备注" style={{ marginTop: 16 }}>
            <Input.TextArea rows={3} />
          </Form.Item>
          {editingServiceGroup ? (
            <Button danger onClick={() => void handleDeleteServiceGroup(editingServiceGroup)}>
              删除服务组
            </Button>
          ) : null}
        </Form>
      </Modal>
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
