# Miniport

Miniport 是一个端口服务资产管理应用，用来管理 IP 主机、固定 10 个端口一组的服务分配、服务容器、依赖组件和源码仓库。

项目已重构为清洁分层结构：

- `internal/appcmd/server`: 服务入口和运行时组装
- `internal/domain/inventory`: 领域模型、校验和业务服务
- `internal/adapter/inventoryhttp`: HTTP API 协议适配
- `internal/infra/*`: 数据库、schema、前端静态资源托管
- `src/app|src/pages|src/domains|src/shared`: 前端应用壳、页面、领域和共享层

## 数据模型

- 主机：一个 IP 节点，例如 `172.22.11.12`。
- 端口组：一个服务占用的 10 个连续端口，例如 `11120-11129`。
- 端口槽位：端口组里的每一个具体端口，记录协议、用途和状态。
- 组件：服务使用的开源项目或外部组件，例如 `kafka`、`nginx`、`etcd`、`redis`。
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
http://127.0.0.1:40239/#/hosts
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
- `GET /api/hosts`
- `POST /api/hosts`
- `PUT /api/hosts/{id}`
- `DELETE /api/hosts/{id}`
- `GET /api/port-groups`
- `POST /api/port-groups`
- `GET /api/port-groups/{id}`
- `PUT /api/port-groups/{id}`
- `DELETE /api/port-groups/{id}`

端口组必须正好包含 10 个端口。后端会拒绝同一主机下互相重叠的端口范围。
