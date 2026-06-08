# Miniport

Miniport 是一个端口服务资产管理应用，用来管理 IP 主机、固定 10 个端口一组的服务分配、服务容器、依赖组件和源码仓库。

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

后端默认监听 `:40238`。Vite 默认监听 `:40237`，并把 `/api` 代理到后端。

如果需要让 Vite 代理到其他后端地址：

```bash
VITE_API_TARGET=http://127.0.0.1:40240 npm run dev
```

## API

主要接口挂载在 `/api` 下：

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
