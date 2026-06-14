import { SaveOutlined } from "@ant-design/icons";
import { Button, Drawer, Form, Input } from "antd";
import type { FormInstance } from "antd";
import type { Host, HostForm } from "../types";

type HostDrawerProps = {
  editingHost: Host | null;
  form: FormInstance<HostForm>;
  onClose: () => void;
  onSave: (values: HostForm) => void;
  open: boolean;
  saving: boolean;
};

export function HostDrawer({
  editingHost,
  form,
  onClose,
  onSave,
  open,
  saving,
}: HostDrawerProps) {
  return (
    <Drawer
      title={editingHost ? "编辑主机" : "新建主机"}
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
        <Form.Item name="ip" label="IP" rules={[{ required: true, message: "请填写 IP" }]}>
          <Input placeholder="172.22.11.12" />
        </Form.Item>
        <Form.Item name="name" label="名称">
          <Input placeholder="node-12" />
        </Form.Item>
        <Form.Item name="network" label="网段">
          <Input placeholder="172.22.11.0/24" />
        </Form.Item>
        <Form.Item name="environment" label="环境">
          <Input placeholder="dev / test / prod" />
        </Form.Item>
        <Form.Item name="notes" label="备注">
          <Input.TextArea rows={4} />
        </Form.Item>
      </Form>
    </Drawer>
  );
}
