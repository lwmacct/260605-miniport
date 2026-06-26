import { Card, Descriptions, Typography } from "antd";
import styles from "@/shared/ui/SectionPage.module.css";

export function AuditRoute() {
  return (
    <section className={styles.section}>
      <div className={styles.sectionHeader}>
        <Typography.Title level={2}>安全审计</Typography.Title>
        <Typography.Paragraph type="secondary">
          预留登录、注册、权限变更和系统事件审计入口。
        </Typography.Paragraph>
      </div>
      <Card className={styles.panel} title="功能占位">
        <Descriptions
          bordered
          column={1}
          size="small"
          items={[
            {
              key: "sources",
              label: "事件来源",
              children: "认证、管理操作、系统任务",
            },
            {
              key: "retention",
              label: "保留策略",
              children: "待接入",
            },
            {
              key: "export",
              label: "导出",
              children: "待接入",
            },
          ]}
        />
      </Card>
    </section>
  );
}
