import { DeleteOutlined, PlusOutlined, SaveOutlined } from "@ant-design/icons";
import { Button, Col, Drawer, Form, Input, InputNumber, Row, Select, Space, Tabs } from "antd";
import type { FormInstance } from "antd";
import {
  relationTypeOptions,
  runtimeModeOptions,
  slotKindOptions,
  statusOptions,
} from "../model/portsvcConstants";
import type { DependencyAssetItem, HostItem, PortGroupForm, PortGroupItem } from "../model/portsvcTypes";

type ServiceDrawerProps = {
  editingGroup: PortGroupItem | null;
  dependencyAssets: DependencyAssetItem[];
  form: FormInstance<PortGroupForm>;
  hosts: HostItem[];
  onClose: () => void;
  onSave: (values: PortGroupForm) => void;
  open: boolean;
  saving: boolean;
};

export function ServiceDrawer({
  editingGroup,
  dependencyAssets,
  form,
  hosts,
  onClose,
  onSave,
  open,
  saving,
}: ServiceDrawerProps) {
  const portPrefix = Form.useWatch("portPrefix", form);
  const start = typeof portPrefix === "number" && portPrefix > 0 ? portPrefix * 10 : 10000;
  const portOptions = Array.from({ length: 10 }, (_, idx) => {
    const port = start + idx;
    return { value: port, label: String(port) };
  });

  return (
    <Drawer
      title={editingGroup ? "编辑端口组" : "新建端口组"}
      open={open}
      size="large"
      onClose={onClose}
      extra={
        <Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={() => form.submit()}>
          保存
        </Button>
      }
    >
      <Form form={form} layout="vertical" onFinish={onSave}>
        <Tabs
          items={[
            {
              key: "base",
              label: "基础信息",
              children: (
                <Row gutter={14}>
                  <Col xs={24} md={12}>
                    <Form.Item name="environmentName" label="运行环境">
                      <Input placeholder="miniport-dind-01 / 4h4g-host-env" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                      <Select options={statusOptions} />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="portPrefix" label="端口组">
                      <InputNumber min={1000} max={5999} step={1} style={{ width: "100%" }} placeholder="自动分配" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="runtimeMode" label="运行模式" rules={[{ required: true }]}>
                      <Select options={runtimeModeOptions} />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="hostId" label="宿主机">
                      <Select
                        allowClear
                        options={hosts.map((host) => ({
                          value: host.id,
                          label: `${host.name}${host.ip ? ` · ${host.ip}` : ""}`,
                        }))}
                      />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="serviceIp" label="服务 IP">
                      <Input placeholder="172.20.0.12" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="runtimeName" label="运行标识">
                      <Input placeholder="miniport-dind-01 / systemd unit" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="environmentOwner" label="负责人">
                      <Input placeholder="platform" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="tags" label="标签">
                      <Input placeholder="api,internal" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="notes" label="备注">
                      <Input.TextArea rows={4} />
                    </Form.Item>
                  </Col>
                </Row>
              ),
            },
            {
              key: "slots",
              label: "端口槽位",
              children: (
                <Form.List name="slots">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ kind: "app", protocol: "tcp", status: "planned" })}>
                        添加槽位
                      </Button>
                      {fields.map((field) => (
                        <Row key={field.key} gutter={10} className="compact-form-row">
                          <Col xs={24} md={4}>
                            <Form.Item name={[field.name, "port"]} rules={[{ required: true }]}>
                              <Select placeholder="端口" options={portOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={5}>
                            <Form.Item name={[field.name, "name"]} rules={[{ required: true }]}>
                              <Input placeholder="redis / mysql / api" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={4}>
                            <Form.Item name={[field.name, "kind"]}>
                              <Select options={slotKindOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={4}>
                            <Form.Item name={[field.name, "protocol"]}>
                              <Input placeholder="tcp / http" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={6}>
                            <Form.Item name={[field.name, "containerName"]}>
                              <Input placeholder="容器名" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={1}>
                            <Button danger icon={<DeleteOutlined />} onClick={() => remove(field.name)} />
                          </Col>
                        </Row>
                      ))}
                    </Space>
                  )}
                </Form.List>
              ),
            },
            {
              key: "assetLinks",
              label: "依赖资产",
              children: (
                <Form.List name="assetLinks">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ relationType: "runtime", required: true })}>
                        添加资产关系
                      </Button>
                      {fields.map((field) => (
                        <Row key={field.key} gutter={10} className="compact-form-row">
                          <Col xs={24} md={9}>
                            <Form.Item name={[field.name, "assetId"]} rules={[{ required: true }]}>
                              <Select
                                showSearch
                                placeholder="选择依赖资产"
                                optionFilterProp="label"
                                options={dependencyAssets.map((asset) => ({
                                  value: asset.id,
                                  label: `${asset.name} · ${asset.assetKind}/${asset.assetType}`,
                                }))}
                              />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={5}>
                            <Form.Item name={[field.name, "relationType"]}>
                              <Select options={relationTypeOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={3}>
                            <Form.Item name={[field.name, "required"]}>
                              <Select options={[{ value: true, label: "必需" }, { value: false, label: "可选" }]} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={6}>
                            <Form.Item name={[field.name, "notes"]}>
                              <Input placeholder="备注" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={1}>
                            <Button danger icon={<DeleteOutlined />} onClick={() => remove(field.name)} />
                          </Col>
                        </Row>
                      ))}
                    </Space>
                  )}
                </Form.List>
              ),
            },
          ]}
        />
      </Form>
    </Drawer>
  );
}
