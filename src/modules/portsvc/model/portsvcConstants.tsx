export const statusOptions = [
  { value: "available", label: "可用" },
  { value: "planned", label: "规划中" },
  { value: "running", label: "运行中" },
  { value: "reserved", label: "保留" },
  { value: "stopped", label: "停用" },
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

export const repositoryKindOptions = [
  { value: "source", label: "源码" },
  { value: "backend", label: "后端" },
  { value: "frontend", label: "前端" },
  { value: "deploy", label: "部署" },
  { value: "infra", label: "基础设施" },
];
