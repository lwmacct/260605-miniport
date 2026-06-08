# Agent 规则

## 服务运行

- 工作区已启动 `npm run dev`, `go run ...` 相关热重载 不需要重复启动。

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

## 前端规则

- 前端使用 React + TypeScript 编写, 组件优先采用函数组件和 React Hooks。
- 前端 UI 优先使用 antd 生态组件, 包括 `antd` 和 `@ant-design/icons`。
- 优先使用 antd 组件自带的布局、表单、表格、抽屉、弹窗、标签、按钮、菜单、Tabs、Space、Grid 等能力。
- 优先通过 antd 的组件 props、Design Token、主题能力和组件组合完成样式调整。
- 尽可能避免编写大量自定义 CSS。只有在 antd 组件能力不足、需要稳定布局尺寸、响应式约束或少量业务视觉表达时才补充 CSS。
- 自定义 CSS 应保持小范围、语义化 class 命名, 不覆盖 antd 全局样式, 不写大面积页面级装饰样式。
