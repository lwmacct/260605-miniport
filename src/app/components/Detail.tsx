import { Typography } from "antd";

export function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="detail-item">
      <Typography.Text type="secondary">{label}</Typography.Text>
      <Typography.Text>{value}</Typography.Text>
    </div>
  );
}
