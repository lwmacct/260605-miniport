import { Button, Divider, Drawer, Form, Input, InputNumber, Select, Space, Tabs } from "antd";
import { DeleteOutlined, PlusOutlined, ReloadOutlined, SaveOutlined } from "@ant-design/icons";
import type { FormInstance } from "antd";
import {
  componentTypeOptions,
  repositoryKindOptions,
  slotStatusOptions,
  statusOptions,
} from "../constants";
import type { GroupForm, Host, PortGroup } from "../types";
import { makeSlots } from "../utils";

type PortGroupDrawerProps = {
  open: boolean;
  saving: boolean;
  hosts: Host[];
  editingGroup: PortGroup | null;
  form: FormInstance<GroupForm>;
  onClose: () => void;
  onSave: (values: GroupForm) => void;
};

export function PortGroupDrawer({
  open,
  saving,
  hosts,
  editingGroup,
  form,
  onClose,
  onSave,
}: PortGroupDrawerProps) {
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
                <div className="form-grid">
                  <Form.Item name="hostId" label="IP 主机" rules={[{ required: true, message: "请选择主机" }]}>
                    <Select options={hosts.map((host) => ({ value: host.id, label: `${host.ip} ${host.name || ""}` }))} />
                  </Form.Item>
                  <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                    <Select options={statusOptions} />
                  </Form.Item>
                  <Form.Item name="portStart" label="起始端口" rules={[{ required: true, message: "请填写起始端口" }]}>
                    <InputNumber min={1} max={65535} className="full-width" />
                  </Form.Item>
                  <Form.Item name="portEnd" label="结束端口" rules={[{ required: true, message: "请填写结束端口" }]}>
                    <InputNumber min={1} max={65535} className="full-width" />
                  </Form.Item>
                  <Form.Item name="serviceName" label="服务名" rules={[{ required: true, message: "请填写服务名" }]}>
                    <Input placeholder="order-service" />
                  </Form.Item>
                  <Form.Item name="containerName" label="服务容器名">
                    <Input placeholder="order-service-dind" />
                  </Form.Item>
                  <Form.Item name="dindHost" label="DIND 宿主">
                    <Input placeholder="dind-01" />
                  </Form.Item>
                  <Form.Item name="owner" label="负责人">
                    <Input placeholder="platform" />
                  </Form.Item>
                  <Form.Item name="tags" label="标签">
                    <Input placeholder="api,redis,internal" />
                  </Form.Item>
                  <Form.Item name="notes" label="备注" className="form-wide">
                    <Input.TextArea rows={4} />
                  </Form.Item>
                </div>
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
                      const portEnd = form.getFieldValue("portEnd");
                      const slots = form.getFieldValue("slots");
                      form.setFieldValue("slots", makeSlots(portStart, portEnd, slots));
                    }}
                  >
                    生成槽位
                  </Button>
                  <Divider />
                  <Form.List name="slots">
                    {(fields) => (
                      <Space orientation="vertical" className="content-stack">
                        {fields.map((field) => (
                          <div key={field.key} className="slot-form-row">
                            <Form.Item name={[field.name, "port"]}>
                              <InputNumber disabled className="port-input" />
                            </Form.Item>
                            <Form.Item name={[field.name, "name"]}>
                              <Input placeholder="名称" />
                            </Form.Item>
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
                            <Form.Item name={[field.name, "status"]}>
                              <Select options={slotStatusOptions} />
                            </Form.Item>
                            <Form.Item name={[field.name, "purpose"]}>
                              <Input placeholder="用途" />
                            </Form.Item>
                          </div>
                        ))}
                      </Space>
                    )}
                  </Form.List>
                </>
              ),
            },
            {
              key: "components",
              label: "组件",
              children: (
                <Form.List name="components">
                  {(fields, { add, remove }) => (
                    <Space orientation="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ type: "opensource" })}>
                        添加组件
                      </Button>
                      {fields.map((field) => (
                        <div key={field.key} className="list-form-row">
                          <Form.Item name={[field.name, "name"]}>
                            <Input placeholder="kafka / nginx / redis" />
                          </Form.Item>
                          <Form.Item name={[field.name, "type"]}>
                            <Select options={componentTypeOptions} />
                          </Form.Item>
                          <Form.Item name={[field.name, "version"]}>
                            <Input placeholder="版本" />
                          </Form.Item>
                          <Form.Item name={[field.name, "url"]}>
                            <Input placeholder="项目地址" />
                          </Form.Item>
                          <Button danger icon={<DeleteOutlined />} onClick={() => remove(field.name)} />
                        </div>
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
                    <Space orientation="vertical" className="content-stack">
                      <Button icon={<PlusOutlined />} onClick={() => add({ kind: "source" })}>
                        添加仓库
                      </Button>
                      {fields.map((field) => (
                        <div key={field.key} className="list-form-row repo-row">
                          <Form.Item name={[field.name, "name"]}>
                            <Input placeholder="仓库名" />
                          </Form.Item>
                          <Form.Item name={[field.name, "kind"]}>
                            <Select options={repositoryKindOptions} />
                          </Form.Item>
                          <Form.Item name={[field.name, "url"]}>
                            <Input placeholder="Git URL" />
                          </Form.Item>
                          <Button danger icon={<DeleteOutlined />} onClick={() => remove(field.name)} />
                        </div>
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
