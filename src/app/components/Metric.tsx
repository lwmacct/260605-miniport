import { Typography } from "antd";

export function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="metric">
      <Typography.Text type="secondary">{label}</Typography.Text>
      <strong>{value}</strong>
    </div>
  );
}
