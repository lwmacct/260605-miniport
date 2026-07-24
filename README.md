# Miniport

Miniport 是一个端口服务资产管理应用，用来管理宿主机、固定 10 个端口一组的项目端口组、端口槽位、依赖组件和源码仓库。

项目使用主项目组装模块模式：

- `internal/appcmd/server`: 最终服务入口、模块组装和运行时依赖
- `github.com/lwmacct/260711-go-pkg-authme/pkg/authme`: 单操作员认证、加密浏览器会话和 HTTP 访问边界
- `internal/handler`: portsvc 业务 HTTP API 协议适配、DTO 和响应转换
- `internal/service`: portsvc 业务规则和应用服务
- `internal/repository`: portsvc 数据访问、record 和当前 schema；业务表不保存登录身份
- `internal/infra/*`: portsvc schema 辅助和前端静态资源托管
- `src/app|src/pages|src/domains|src/shared`: 前端应用壳、页面、领域和共享层

<!--TOC-->

## Table of Contents

- [数据模型](#数据模型) `:27+10`
- [GitHub App](#github-app) `:37+12`
- [单操作员认证](#单操作员认证) `:49+8`
- [本地运行](#本地运行) `:57+28`
- [API](#api) `:85+41`

<!--TOC-->

## 数据模型

- 宿主机：承载端口组的物理机、云主机或小规格设备，例如 `4h4g`。
- 端口组：系统全局分配的 10 个连续端口，例如 `11120-11129`，通常对应一个项目和一个 DIND/直跑运行单元。
- 端口槽位：端口组内的具体服务组件，例如 `redis:11120`、`mysql:11121`、`kafka:11122`。
- 依赖资产：端口组项目使用的代码仓库、闭源服务、SaaS、组件或文档，例如 GitHub 仓库、闭源 API、`kafka`、`nginx`、`redis`。
- 资产关系：端口组到依赖资产的关系，例如源码、运行依赖、构建依赖、部署、基础设施、API 或文档。
- GitHub 仓库：通过 GitHub App 安装同步的公开或私有仓库，以 GitHub repository ID 作为稳定身份。
- 仓库关系：端口组或端口槽位到 GitHub 仓库的源码、构建、部署或文档关系。

## GitHub App

仓库同步使用 GitHub App Installation Token，不接收用户 Personal Access Token。GitHub App 注册时配置：

- Repository permissions：`Metadata: Read-only`。
- Setup URL：`https://<miniport-host>/api/github/setup`。
- Webhook URL：`https://<miniport-host>/api/github/webhooks`。
- Webhook events：`installation`、`installation_repositories`、`repository`。
- 安装范围：`Any account`。

在 `server.github` 中设置 `enabled`、`app-id`、`app-slug`、`private-key-file` 和 `webhook-secret`。操作员需要分别在个人账号和每个 GitHub Organization 安装 App；系统全局同步安装时授权的仓库。

## 单操作员认证

认证只负责保护 HTTP 访问，不参与端口组、服务组、依赖资产或 GitHub 安装的数据归属。系统不提供用户表、注册、角色、管理员或按用户过滤的数据路径。

默认启用一个静态 Access Token。设置 `AUTHME_ACCESS_TOKEN`，并为浏览器会话设置一个 base64url 编码的 32 字节 `AUTHME_SESSION_KEY`。Dex GitHub OIDC 默认启用，部署时需配置 `server.http.authme.dexgithub.client-secret`，并可用 `allowed-github-user` 调整唯一允许登录的 GitHub 用户名。

本次模型不提供历史数据库迁移。部署新版本前必须备份旧库并创建空数据库，让服务按当前 schema 初始化；不要把新二进制直接指向仍含 `users`、`auth_sessions`、`owner_subject` 或 `github_installation_subjects` 的旧库。

## 本地运行

```bash
go run . server
corepack enable
pnpm install
pnpm run dev
```

后端默认监听 `:40238`。Vite 默认监听 `:40239`，并把 `/api` 和 `/authme` 代理到后端。

前端使用 hash 路由。主要页面入口：

```text
http://127.0.0.1:40239/#/console/overview
http://127.0.0.1:40239/#/console/projects
http://127.0.0.1:40239/#/console/service-groups
http://127.0.0.1:40239/#/console/dependencies
```

界面支持明亮/暗色主题切换。主题状态保存在浏览器 `localStorage` 中，刷新后会保持上次选择。

如果需要让 Vite 代理到其他后端地址：

```bash
API_PROXY_TARGET=http://127.0.0.1:40240 pnpm run dev
```

## API

主要接口挂载在 `/api` 下：

- `GET /authme/session`
- `POST /authme/login/token`
- `GET /authme/login/github`
- `DELETE /authme/session`
- `GET /api/health`
- `GET /api/meta`
- `GET /api/github/status`
- `POST /api/github/connections`
- `GET /api/github/setup`
- `GET /api/github/installations`
- `POST /api/github/installations/{id}/sync`
- `GET /api/github/repositories`
- `POST /api/github/webhooks`
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
- `GET /api/service-groups`
- `POST /api/service-groups`
- `GET /api/service-groups/{id}`
- `PUT /api/service-groups/{id}`
- `DELETE /api/service-groups/{id}`

端口组必须正好包含 10 个端口，端口起点在系统内全局唯一。端口槽位的端口必须落在端口组范围内，且同一端口组内端口不能重复。依赖资产、服务组和 GitHub 安装均为全局数据，端口组通过关系引用它们。认证身份不会写入业务表或业务 API。
