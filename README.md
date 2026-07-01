# Miniport

Miniport 是一个端口服务资产管理应用，用来管理宿主机、固定 10 个端口一组的项目端口组、端口槽位、依赖组件和源码仓库。

项目使用主项目组装模块模式：

- `internal/appcmd/server`: 最终服务入口、模块组装和运行时依赖
- `github.com/lwmacct/260630-go-hsr-auth/pkg/auth`: 身份认证、用户、会话、管理员用户 API
- `internal/handler`: portsvc 业务 HTTP API 协议适配、DTO 和响应转换
- `internal/service`: portsvc 业务规则和应用服务
- `internal/repository`: portsvc 数据访问、record 和 schema；只保存业务表和稳定身份主体
- `internal/infra/*`: portsvc schema 辅助和前端静态资源托管
- `src/app|src/pages|src/domains|src/shared`: 前端应用壳、页面、领域和共享层

## 数据模型

- 宿主机：承载端口组的物理机、云主机或小规格设备，例如 `4h4g`。
- 端口组：一个身份主体占用的 10 个连续端口，例如 `11120-11129`，通常对应一个项目和一个 DIND/直跑运行单元。
- 端口槽位：端口组内的具体服务组件，例如 `redis:11120`、`mysql:11121`、`kafka:11122`。
- 依赖资产：端口组项目使用的代码仓库、闭源服务、SaaS、组件或文档，例如 GitHub 仓库、闭源 API、`kafka`、`nginx`、`redis`。
- 资产关系：端口组到依赖资产的关系，例如源码、运行依赖、构建依赖、部署、基础设施、API 或文档。

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
- `GET /api/hosts`
- `POST /api/hosts`
- `PUT /api/hosts/{id}`
- `DELETE /api/hosts/{id}`
- `GET /api/dependency-assets`
- `POST /api/dependency-assets`
- `PUT /api/dependency-assets/{id}`
- `DELETE /api/dependency-assets/{id}`
- `GET /api/port-groups`
- `POST /api/port-groups`
- `GET /api/port-groups/{id}`
- `PUT /api/port-groups/{id}`
- `DELETE /api/port-groups/{id}`
- `POST /api/port-groups/{id}/slots`
- `PUT /api/port-slots/{id}`
- `DELETE /api/port-slots/{id}`
- `GET /api/port-groups/export.csv`

端口组必须正好包含 10 个端口。同一身份主体下端口起点不能重复。端口槽位的端口必须落在端口组范围内，且同一端口组内端口不能重复。依赖资产是全局资产，端口组通过资产关系引用它们。身份认证和 `users` 表由 auth 模块拥有，Miniport 只保存稳定 UUID7 `owner_subject` 作为业务归属语义，并通过 shared `identity.Directory` 解析主体展示信息。
