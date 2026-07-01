export const statusOptions = [
  { value: "available", label: "可用" },
  { value: "planned", label: "规划中" },
  { value: "running", label: "运行中" },
  { value: "reserved", label: "保留" },
  { value: "stopped", label: "停用" },
];

export const serviceGroupStatusOptions = [
  { value: "active", label: "运行中" },
  { value: "planned", label: "规划中" },
  { value: "stopped", label: "停用" },
];

export const serviceGroupKindOptions = [
  { value: "service", label: "服务" },
  { value: "cluster", label: "集群" },
  { value: "stack", label: "服务栈" },
  { value: "other", label: "其他" },
];

export const runtimeModeOptions = [
  { value: "dind", label: "DIND" },
  { value: "host", label: "宿主机直跑" },
];

export const slotKindOptions = [
  { value: "app", label: "应用" },
  { value: "db", label: "数据库" },
  { value: "cache", label: "缓存" },
  { value: "mq", label: "消息队列" },
  { value: "middleware", label: "中间件" },
];

export const componentTypeOptions = [
  { value: "opensource", label: "开源项目" },
  { value: "external", label: "外部组件" },
  { value: "internal", label: "自有组件" },
];

export const assetKindOptions = [
  { value: "repository", label: "代码仓库" },
  { value: "service", label: "闭源/外部服务" },
  { value: "component", label: "组件依赖" },
  { value: "document", label: "文档" },
  { value: "other", label: "其他" },
];

export const assetTypeOptions = [
  { value: "owned", label: "自有可控" },
  { value: "opensource", label: "开源项目" },
  { value: "closed_source", label: "闭源服务" },
  { value: "saas", label: "SaaS" },
  { value: "third_party", label: "第三方" },
  { value: "internal_blackbox", label: "内部黑盒" },
  { value: "middleware", label: "中间件" },
];

export const assetProviderOptions = [
  { value: "manual", label: "手动" },
  { value: "github", label: "GitHub" },
  { value: "gitlab", label: "GitLab" },
  { value: "gitea", label: "Gitea" },
  { value: "vendor", label: "厂商" },
  { value: "internal", label: "内部" },
  { value: "other", label: "其他" },
];

export const visibilityOptions = [
  { value: "unknown", label: "未知" },
  { value: "private", label: "私有" },
  { value: "public", label: "公开" },
  { value: "internal", label: "内部" },
];

export const controllabilityOptions = [
  { value: "unknown", label: "未知" },
  { value: "full", label: "完全可控" },
  { value: "partial", label: "部分可控" },
  { value: "vendor", label: "厂商控制" },
  { value: "none", label: "不可控" },
];

export const relationTypeOptions = [
  { value: "source", label: "源码" },
  { value: "runtime", label: "运行依赖" },
  { value: "build", label: "构建依赖" },
  { value: "deploy", label: "部署" },
  { value: "infra", label: "基础设施" },
  { value: "api", label: "API" },
  { value: "docs", label: "文档" },
  { value: "observability", label: "可观测" },
];
