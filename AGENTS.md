# Agent 规则

## 服务运行

- 工作区已启动 `npm run dev`, `go run ...` 相关热重载 不需要重复启动。

## 结构

- Go 服务模块: `github.com/lwmacct/260605-miniport`。
- 服务入口在 `main.go`; 后端启动、命令和基础设施代码放 `internal/app/`。
- 领域代码放 `internal/modules/<domain>`。
- `internal/modules/<domain>/` 内文件名使用通用职责语义, 例如 `model.go`, `dto.go`, `handler.go`, `service.go`, `repository.go`, `schema.go`, `validation.go`, `errors.go`, `utils.go`。
- 仅明确要对外复用的共享包放 `pkg/`。
- 前端是 React/Vite, 代码在 `src/`, 静态资源在 `public/`, 不提交构建产物。
- 测试文件放在被测代码旁边, 命名为 `*_test.go`。

## 常用命令

- `psql -c 'select 1;'`: 验证默认 PostgreSQL 连接。
- "sqlite3,litecli,sqlite-utils" 命令可用于操作 sqlite 数据库

## Language Rules

- Go 源码拆分命名规则:
  - 拆分优先依据职责边界和代码内聚性, 不以行数作为硬限制。
  - 仅为降低行数而做机械拆分时, 参考 300 行作为避免文件碎片化的下限。
  - 结构体方法按 `type.topic.go` 命名, 例如 `type.action.go`。
  - 纯包级函数逻辑统一放到 `utils.go`。
- Shell 命名规则:
  - 函数名使用双下划线前缀, 例如 `__func_name`。
  - 普通局部变量使用单下划线前缀, 例如 `_var_name`。
  - 全大写环境变量不受此规则约束。

## 后端规则

- 使用标准 Go 格式和简短小写包名。
- handler 保持简薄, 只做协议适配、参数接收和响应转换。
- 业务规则放 service, 数据访问放 repository, schema/迁移逻辑放 schema 文件, 不混进请求处理。
- DTO 命名保持明确, 例如 `CreateInput`, `UpdateInput`, `Response`; 现有接口命名可按当前 API 语义延续。
- 配置使用 `cfgm.MustLoadCmd`, 默认值从 `internal/config` 读取。

## 数据模型

- 本项目不要求历史兼容。需求或数据模型变化时, 直接修改当前代码、API、数据库字段和前端调用。
- 不为旧接口、旧字段、旧数据路径或旧行为写兼容层、回填分支、废弃提示或历史记录。
- 数据库 schema 表达当前真实模型。除非明确要求迁移已有生产数据, 否则不要为了兼容历史字段保留旧列、读旧列或双写旧列。

## 前端规则

- 前端使用 React + TypeScript 编写, 组件优先采用函数组件和 React Hooks。
- 前端 UI 优先使用 antd 生态组件, 包括 `antd` 和 `@ant-design/icons`。
- 优先使用 antd 组件自带的布局、表单、表格、抽屉、弹窗、标签、按钮、菜单、Tabs、Space、Grid 等能力。
- 优先通过 antd 的组件 props、Design Token、主题能力和组件组合完成样式调整。
- 尽可能避免编写大量自定义 CSS。只有在 antd 组件能力不足、需要稳定布局尺寸、响应式约束或少量业务视觉表达时才补充 CSS。
- 自定义 CSS 应保持小范围、语义化 class 命名, 不覆盖 antd 全局样式, 不写大面积页面级装饰样式。
