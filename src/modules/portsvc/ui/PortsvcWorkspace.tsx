import { DownloadOutlined, PlusOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Alert, Button, Flex, Form, Input, InputNumber, Modal, Select, Space, Spin, Typography, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useAuthStateQuery } from "@/modules/auth";
import { usePortsvcQuery, portsvcKeys } from "../model/portsvcQueries";
import {
  exportServicesURL,
  removePortAllocation,
  removeService,
  savePortAllocation,
  saveService,
} from "../api/portsvcApi";
import { statusOptions } from "../model/portsvcConstants";
import type { PortsvcQuery, PortAllocation, PortAllocationForm, ServiceForm, ServiceItem } from "../model/portsvcTypes";
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
  const [sort, setSort] = useState<string>("name");
  const [saving, setSaving] = useState(false);
  const [selectedService, setSelectedService] = useState<ServiceItem | null>(null);
  const [editingService, setEditingService] = useState<ServiceItem | null>(null);
  const [serviceDrawerOpen, setServiceDrawerOpen] = useState(false);
  const [portModalOpen, setPortModalOpen] = useState(false);
  const [editingPort, setEditingPort] = useState<PortAllocation | null>(null);
  const [serviceForm] = Form.useForm<ServiceForm>();
  const [portForm] = Form.useForm<PortAllocationForm>();

  const query = useMemo<PortsvcQuery>(() => {
    return {
      serviceQuery: search.trim(),
      serviceSort: sort,
      status,
    };
  }, [search, sort, status]);

  const portsvcQuery = usePortsvcQuery(query);
  const snapshot = portsvcQuery.data;
  const services: ServiceItem[] = snapshot?.services ?? [];
  const ports: PortAllocation[] = snapshot?.ports ?? [];
  const meta = snapshot?.meta;
  const canManage = Boolean(authState.data?.session.authenticated);
  const stats = useMemo(() => buildStats(services, ports), [ports, services]);

  useEffect(() => {
    setSelectedService(null);
  }, [query]);

  async function refresh() {
    await queryClient.invalidateQueries({ queryKey: portsvcKeys.snapshot(query) });
  }

  function openCreateService() {
    setEditingService(null);
    serviceForm.resetFields();
    serviceForm.setFieldsValue({
      status: "planned",
      repositories: [],
      dependencies: [],
    });
    setServiceDrawerOpen(true);
  }

  function openEditService(service: ServiceItem) {
    setEditingService(service);
    serviceForm.setFieldsValue({
      ...service,
      dependencies: service.dependencies,
      repositories: service.repositories,
    });
    setServiceDrawerOpen(true);
  }

  async function handleSaveService(values: ServiceForm) {
    setSaving(true);
    try {
      await saveService(values, editingService);
      message.success("服务已保存");
      setServiceDrawerOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteService(service: ServiceItem) {
    await removeService(service);
    message.success("服务已删除");
    setSelectedService(null);
    await refresh();
  }

  function openCreatePort() {
    setEditingPort(null);
    portForm.resetFields();
    portForm.setFieldsValue({ status: "available" });
    setPortModalOpen(true);
  }

  function openEditPort(port: PortAllocation) {
    setEditingPort(port);
    portForm.setFieldsValue(port);
    setPortModalOpen(true);
  }

  async function handleSavePort(values: PortAllocationForm) {
    setSaving(true);
    try {
      await savePortAllocation(values, editingPort);
      message.success("端口组已保存");
      setPortModalOpen(false);
      await refresh();
    } catch (error) {
      message.error(error instanceof Error ? error.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  async function handleDeletePort(port: PortAllocation) {
    await removePortAllocation(port);
    message.success("端口组已删除");
    await refresh();
  }

  const exportURL = exportServicesURL(query);
  const sortOptions = [
    { label: "名称", value: "name" },
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
            onDeleteService={(service) => void handleDeleteService(service)}
            onEditService={openEditService}
            onSelectService={setSelectedService}
            services={services}
          />
        );
        break;
      case "projects":
        content = <ProjectServicesSection onSelectService={setSelectedService} services={services} />;
        break;
      case "dependencies":
        content = <DependenciesSection services={services} stats={stats} />;
        break;
      default:
        content = (
          <OverviewSection
            canManage={canManage}
            onCreateService={openCreateService}
            onSelectService={setSelectedService}
            services={services}
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
              placeholder="搜索服务、项目、DIND IP、仓库"
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
            <Button icon={<PlusOutlined />} onClick={openCreatePort} disabled={!canManage}>
              端口组
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreateService} disabled={!canManage}>
              服务
            </Button>
          </Space>
        </Flex>
      </section>

      <section className="app-content">{content}</section>

      <ServiceDetailDrawer
        canManage={canManage}
        service={selectedService}
        onClose={() => setSelectedService(null)}
        onDelete={handleDeleteService}
        onEdit={(service) => {
          setSelectedService(null);
          openEditService(service);
        }}
      />
      <ServiceDrawer
        editingService={editingService}
        form={serviceForm}
        onClose={() => setServiceDrawerOpen(false)}
        onSave={(values) => void handleSaveService(values)}
        open={serviceDrawerOpen}
        ports={ports}
        saving={saving}
      />
      <Modal
        title={editingPort ? "编辑端口组" : "新建端口组"}
        open={portModalOpen}
        onCancel={() => setPortModalOpen(false)}
        onOk={() => portForm.submit()}
        confirmLoading={saving}
      >
        <Form form={portForm} layout="vertical" onFinish={(values) => void handleSavePort(values)}>
          <Form.Item name="portStart" label="起始端口">
            <InputNumber min={10000} max={59990} step={10} style={{ width: "100%" }} placeholder="自动分配" />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select options={[{ value: "available", label: "可用" }, ...statusOptions]} />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          {editingPort ? (
            <Button danger onClick={() => void handleDeletePort(editingPort)}>
              删除端口组
            </Button>
          ) : null}
        </Form>
      </Modal>
    </>
  );
}
