import { Card, Descriptions, Typography } from "antd";
import styles from "@/shared/ui/SectionPage.module.css";

export function SystemRoute() {
  return (
    <section className={styles.section}>
      <div className={styles.sectionHeader}>
        <Typography.Title level={2}>系统设置</Typography.Title>
        <Typography.Paragraph type="secondary">
          预留服务运行参数、注册开关、会话策略和数据库状态入口。
        </Typography.Paragraph>
      </div>
      <Card className={styles.panel} title="功能占位">
        <Descriptions
          bordered
          column={1}
          size="small"
          items={[
            {
              key: "registration",
              label: "注册开关",
              children: "由 server.auth.local.registration-enabled 控制",
            },
            {
              key: "admins",
              label: "管理员",
              children: "由 server.auth.admins 运行时配置控制",
            },
            {
              key: "database",
              label: "数据库",
              children: "sqlite 默认，可选 pgsql",
            },
          ]}
        />
      </Card>
    </section>
  );
}
