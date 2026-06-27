import {
  DeleteOutlined,
  PlusOutlined,
  ReloadOutlined,
  SaveOutlined,
} from "@ant-design/icons";
import { Alert, Button, Col, Divider, Drawer, Form, Input, InputNumber, Row, Select, Space, Tabs } from "antd";
import type { FormInstance } from "antd";
import { componentTypeOptions, repositoryKindOptions, slotStatusOptions, statusOptions } from "../model/inventoryConstants";
import type { GroupForm, PortGroup } from "../model/inventoryTypes";
import { makeSlots } from "../model/inventoryUtils";

type PortGroupDrawerProps = {
  editingGroup: PortGroup | null;
  form: FormInstance<GroupForm>;
  onClose: () => void;
  onSave: (values: GroupForm) => void;
  open: boolean;
  saving: boolean;
};

export function PortGroupDrawer({
  editingGroup,
  form,
  onClose,
  onSave,
  open,
  saving,
}: PortGroupDrawerProps) {
  return (
    <Drawer
      title={editingGroup ? "编辑端口分配" : "新建端口分配"}
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
                  {!editingGroup ? (
                    <Col span={24}>
                      <Alert showIcon type="info" message="不填写起始端口时，系统会自动分配下一个可用的 10 端口组。" />
                    </Col>
                  ) : null}
                  <Col xs={24} md={12}>
                    <Form.Item name="name" label="分配名称" rules={[{ required: true, message: "请填写分配名称" }]}>
                      <Input placeholder="miniport-dev / order-service" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                      <Select options={statusOptions} />
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
                    <Form.Item name="portStart" label="起始端口">
                      <InputNumber min={10000} max={59990} step={10} style={{ width: "100%" }} placeholder="自动分配" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="owner" label="负责人">
                      <Input placeholder="platform" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="tags" label="标签">
                      <Input placeholder="api,redis,internal" />
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
                <>
                  <Button
                    icon={<ReloadOutlined />}
                    onClick={() => {
                      const portStart = form.getFieldValue("portStart");
                      const portEnd = portStart ? portStart + 9 : undefined;
                      const slots = form.getFieldValue("slots");
                      form.setFieldValue("slots", makeSlots(portStart, portEnd, slots));
                    }}
                  >
                    生成槽位
                  </Button>
                  <Divider />
                  <Form.List name="slots">
                    {(fields) => (
                      <Space direction="vertical" className="content-stack">
                        {fields.map((field) => (
                          <Row key={field.key} gutter={10} className="compact-form-row">
                            <Col xs={24} sm={6} lg={4}>
                              <Form.Item name={[field.name, "port"]}>
                                <InputNumber disabled style={{ width: "100%" }} />
                              </Form.Item>
                            </Col>
                            <Col xs={24} sm={18} lg={6}>
                              <Form.Item name={[field.name, "name"]}>
                                <Input placeholder="名称" />
                              </Form.Item>
                            </Col>
                            <Col xs={12} lg={4}>
                              <Form.Item name={[field.name, "protocol"]}>
                                <Select
                                  options={[
                                    { value: "tcp", label: "tcp" },
                                    { value: "udp", label: "udp" },
                                    { value: "http", label: "http" },
                                    { value: "https", label: "https" },
                                  ]}
                                />
                              </Form.Item>
                            </Col>
                            <Col xs={12} lg={4}>
                              <Form.Item name={[field.name, "status"]}>
                                <Select options={slotStatusOptions} />
                              </Form.Item>
                            </Col>
                            <Col xs={24} lg={6}>
                              <Form.Item name={[field.name, "purpose"]}>
                                <Input placeholder="用途" />
                              </Form.Item>
                            </Col>
                          </Row>
                        ))}
                      </Space>
                    )}
                  </Form.List>
                </>
              ),
            },
            {
              key: "projects",
              label: "项目",
              children: (
                <Form.List name="projects">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({})}>
                        添加项目
                      </Button>
                      {fields.map((field) => (
                        <Row key={field.key} gutter={10} className="compact-form-row">
                          <Col xs={24} md={7}>
                            <Form.Item name={[field.name, "name"]}>
                              <Input placeholder="项目名" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} md={10}>
                            <Form.Item name={[field.name, "description"]}>
                              <Input placeholder="说明" />
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
            {
              key: "components",
              label: "依赖",
              children: (
                <Form.List name="components">
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
                          <Col xs={24} sm={12} md={5}>
                            <Form.Item name={[field.name, "type"]}>
                              <Select options={componentTypeOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} sm={12} md={4}>
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
                          <Col xs={24} md={7}>
                            <Form.Item name={[field.name, "name"]}>
                              <Input placeholder="仓库名" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} sm={10} md={5}>
                            <Form.Item name={[field.name, "kind"]}>
                              <Select options={repositoryKindOptions} />
                            </Form.Item>
                          </Col>
                          <Col xs={24} sm={14} md={11}>
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
          ]}
        />
      </Form>
    </Drawer>
  );
}
