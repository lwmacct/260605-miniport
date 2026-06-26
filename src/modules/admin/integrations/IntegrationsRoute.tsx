import { Card, Descriptions, Typography } from "antd";
import styles from "@/shared/ui/SectionPage.module.css";

export function IntegrationsRoute() {
  return (
    <section className={styles.section}>
      <div className={styles.sectionHeader}>
        <Typography.Title level={2}>集成配置</Typography.Title>
        <Typography.Paragraph type="secondary">
          预留验证码、Webhook、第三方服务和通知渠道入口。
        </Typography.Paragraph>
      </div>
      <Card className={styles.panel} title="功能占位">
        <Descriptions
          bordered
          column={1}
          size="small"
          items={[
            {
              key: "challenge",
              label: "验证码",
              children: "image / hcaptcha / turnstile",
            },
            {
              key: "webhook",
              label: "Webhook",
              children: "待接入",
            },
            {
              key: "notification",
              label: "通知",
              children: "待接入",
            },
          ]}
        />
      </Card>
    </section>
  );
}
