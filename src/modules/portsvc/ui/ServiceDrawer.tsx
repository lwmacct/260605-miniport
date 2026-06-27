import { DeleteOutlined, PlusOutlined, SaveOutlined } from "@ant-design/icons";
import { Button, Col, Drawer, Form, Input, Row, Select, Space, Tabs } from "antd";
import type { FormInstance } from "antd";
import { componentTypeOptions, repositoryKindOptions, statusOptions } from "../model/portsvcConstants";
import type { PortAllocation, ServiceForm, ServiceItem } from "../model/portsvcTypes";

type ServiceDrawerProps = {
  editingService: ServiceItem | null;
  form: FormInstance<ServiceForm>;
  onClose: () => void;
  onSave: (values: ServiceForm) => void;
  open: boolean;
  ports: PortAllocation[];
  saving: boolean;
};

export function ServiceDrawer({
  editingService,
  form,
  onClose,
  onSave,
  open,
  ports,
  saving,
}: ServiceDrawerProps) {
  return (
    <Drawer
      title={editingService ? "编辑服务" : "新建服务"}
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
                    <Form.Item name="name" label="服务名称" rules={[{ required: true, message: "请填写服务名称" }]}>
                      <Input placeholder="miniport-api / order-service" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                      <Select options={statusOptions} />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="projectName" label="项目">
                      <Input placeholder="miniport" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="portAllocationId" label="端口组">
                      <Select
                        allowClear
                        options={ports.map((port) => ({
                          value: port.id,
                          label: `${port.portStart}-${port.portEnd}`,
                        }))}
                      />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="dindIp" label="DIND 内网 IP">
                      <Input placeholder="172.20.0.12" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="dindContainer" label="DIND 容器">
                      <Input placeholder="miniport-dind-01" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="owner" label="负责人">
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
              key: "repositories",
              label: "仓库",
              children: (
                <Form.List name="repositories">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ kind: "source" })}>
                        添加仓库
                      </Button>
                      {fields.map((field) => (
                        <Row key={field.key} gutter={10} className="compact-form-row">
                          <Col xs={24} md={6}>
                            <Form.Item name={[field.name, "name"]}>
                              <Input placeholder="仓库名" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={5}>
                            <Form.Item name={[field.name, "kind"]}>
                              <Select options={repositoryKindOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={12}>
                            <Form.Item name={[field.name, "url"]}>
                              <Input placeholder="Git URL" />
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
              key: "dependencies",
              label: "依赖",
              children: (
                <Form.List name="dependencies">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ type: "opensource" })}>
                        添加依赖
                      </Button>
                      {fields.map((field) => (
                        <Row key={field.key} gutter={10} className="compact-form-row">
                          <Col xs={24} md={7}>
                            <Form.Item name={[field.name, "name"]}>
                              <Input placeholder="kafka / etcd / postgresql" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={5}>
                            <Form.Item name={[field.name, "type"]}>
                              <Select options={componentTypeOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={4}>
                            <Form.Item name={[field.name, "version"]}>
                              <Input placeholder="版本" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={7}>
                            <Form.Item name={[field.name, "url"]}>
                              <Input placeholder="项目地址" />
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
