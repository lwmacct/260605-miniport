import { Button, Checkbox, Form, Input, Modal, Select } from "antd";
import type { FormInstance } from "antd";
import { statusOptions } from "../model/inventoryConstants";

export type BatchPortGroupForm = {
  applyOwner?: boolean;
  applyStatus?: boolean;
  applyTags?: boolean;
  owner?: string;
  status?: string;
  tags?: string;
};

type BatchPortGroupModalProps = {
  form: FormInstance<BatchPortGroupForm>;
  onCancel: () => void;
  onSubmit: (values: BatchPortGroupForm) => void;
  open: boolean;
  selectedCount: number;
  saving: boolean;
};

export function BatchPortGroupModal({
  form,
  onCancel,
  onSubmit,
  open,
  selectedCount,
  saving,
}: BatchPortGroupModalProps) {
  return (
    <Modal
      title={`批量更新端口组 (${selectedCount})`}
      open={open}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          取消
        </Button>,
        <Button key="submit" type="primary" loading={saving} onClick={() => form.submit()}>
          保存
        </Button>,
      ]}
    >
      <Form form={form} layout="vertical" onFinish={onSubmit}>
        <Form.Item name="applyStatus" valuePropName="checked">
          <Checkbox>更新状态</Checkbox>
        </Form.Item>
        <Form.Item shouldUpdate noStyle>
          {() => (
            <Form.Item
              name="status"
              label="状态"
              rules={form.getFieldValue("applyStatus") ? [{ required: true, message: "请选择状态" }] : []}
            >
              <Select disabled={!form.getFieldValue("applyStatus")} options={statusOptions} />
            </Form.Item>
          )}
        </Form.Item>

        <Form.Item name="applyOwner" valuePropName="checked">
          <Checkbox>更新负责人</Checkbox>
        </Form.Item>
        <Form.Item name="owner" label="负责人">
          <Input disabled={!form.getFieldValue("applyOwner")} placeholder="留空可清空负责人" />
        </Form.Item>

        <Form.Item name="applyTags" valuePropName="checked">
          <Checkbox>更新标签</Checkbox>
        </Form.Item>
        <Form.Item name="tags" label="标签">
          <Input disabled={!form.getFieldValue("applyTags")} placeholder="留空可清空标签" />
        </Form.Item>
      </Form>
    </Modal>
  );
}
