# Miniport

Miniport 是一个端口服务资产管理应用，用来管理 IP 主机、固定 10 个端口一组的服务分配、服务容器、依赖组件和源码仓库。

项目使用主项目组装模块模式：

- `internal/appcmd/server`: 最终服务入口、模块组装和运行时依赖
- `github.com/lwmacct/260630-go-hsr-auth/pkg/auth`: 身份认证、用户、会话、管理员用户 API
- `internal/handler`: portsvc 业务 HTTP API 协议适配、DTO 和响应转换
- `internal/service`: portsvc 业务规则和应用服务
- `internal/repository`: portsvc 数据访问、record 和 schema；只通过只读 SQL 查询引用 auth-owned `users`
- `internal/infra/*`: portsvc schema 辅助和前端静态资源托管
- `src/app|src/pages|src/domains|src/shared`: 前端应用壳、页面、领域和共享层

## 数据模型

- 端口分配：一个身份主体占用的 10 个连续端口，例如 `11120-11129`。
- 服务：绑定到身份主体的服务资产，可关联端口分配、项目、DIND 信息、负责人、标签和备注。
- 依赖：服务使用的开源项目或外部组件，例如 `kafka`、`nginx`、`etcd`、`redis`。
- 仓库：服务关联的源码、部署、前端、后端或基础设施仓库。

## 本地运行

```bash
go run . server
npm install
npm run dev
```

后端默认监听 `:40238`。Vite 默认监听 `:40239`，并把 `/api` 代理到后端。

前端使用 hash 路由。页面拆成四个独立入口：

```text
http://127.0.0.1:40239/#/overview
http://127.0.0.1:40239/#/services
http://127.0.0.1:40239/#/projects
http://127.0.0.1:40239/#/dependencies
```

界面支持明亮/暗色主题切换。主题状态保存在浏览器 `localStorage` 中，刷新后会保持上次选择。

如果需要让 Vite 代理到其他后端地址：

```bash
API_PROXY_TARGET=http://127.0.0.1:40240 npm run dev
```

## API

主要接口挂载在 `/api` 下：

- `GET /api/health`
- `GET /api/meta`
- `GET /api/auth/config`
- `POST /api/auth/challenges`
- `POST /api/auth/password/login`
- `POST /api/auth/password/register`
- `GET /api/auth/me`
- `POST /api/auth/logout`
- `GET /api/admin/users`
- `GET /api/services`
- `POST /api/services`
- `GET /api/services/{id}`
- `PUT /api/services/{id}`
- `DELETE /api/services/{id}`
- `POST /api/services/batch-delete`
- `GET /api/services/export.csv`
- `GET /api/port-allocations`
- `POST /api/port-allocations`
- `PUT /api/port-allocations/{id}`
- `DELETE /api/port-allocations/{id}`

端口分配必须正好包含 10 个端口。同一身份主体下端口起点不能重复。身份认证和 `users` 表由 auth 模块拥有，Miniport 只保留 `user_id` 作为业务归属语义。
