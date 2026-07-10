import {
  DisconnectOutlined,
  GithubOutlined,
  LinkOutlined,
  ReloadOutlined,
  SearchOutlined,
  SettingOutlined,
} from "@ant-design/icons";
import { WorkbenchPage, WorkbenchPanel } from "@lwmacct/260627-antd-workbench";
import { useQueryClient } from "@tanstack/react-query";
import { Alert, Avatar, Button, Input, Popconfirm, Space, Table, Tag, Tooltip, Typography, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import {
  beginGitHubConnection,
  disconnectGitHubInstallation,
  syncGitHubInstallation,
  type GitHubInstallation,
  type GitHubRepository,
} from "@/modules/github/api/githubApi";
import {
  githubKeys,
  useGitHubInstallationsQuery,
  useGitHubRepositoriesQuery,
  useGitHubStatusQuery,
} from "@/modules/github/model/githubQueries";

function formatTime(value?: string) {
  return value ? new Date(value).toLocaleString() : "-";
}

export function GitHubSettingsRoute() {
  const location = useLocation();
  const queryClient = useQueryClient();
  const [query, setQuery] = useState("");
  const [busyId, setBusyId] = useState<string>();
  const statusQuery = useGitHubStatusQuery();
  const installationsQuery = useGitHubInstallationsQuery();
  const repositoriesQuery = useGitHubRepositoriesQuery(query);
  const enabled = Boolean(statusQuery.data?.enabled);
  const loadError = statusQuery.error ?? installationsQuery.error ?? repositoriesQuery.error;

  useEffect(() => {
    if (new URLSearchParams(location.search).get("github") === "connected") {
      message.success("GitHub 已连接，仓库同步完成");
      void queryClient.invalidateQueries({ queryKey: githubKeys.all });
    }
  }, [location.search, queryClient]);

  async function connect() {
    try {
      const result = await beginGitHubConnection();
      window.location.assign(result.url);
    } catch (error) {
      message.error(error instanceof Error ? error.message : "无法连接 GitHub");
    }
  }

  async function sync(item: GitHubInstallation) {
    setBusyId(item.id);
    try {
      await syncGitHubInstallation(item.id);
      await queryClient.invalidateQueries({ queryKey: githubKeys.all });
      await queryClient.invalidateQueries({ queryKey: ["portsvc", "snapshot"] });
      message.success(`${item.accountLogin} 同步完成`);
    } catch (error) {
	  await queryClient.invalidateQueries({ queryKey: githubKeys.all });
      message.error(error instanceof Error ? error.message : "同步失败");
    } finally {
      setBusyId(undefined);
    }
  }

  async function disconnect(item: GitHubInstallation) {
    setBusyId(item.id);
    try {
      await disconnectGitHubInstallation(item.id);
      await queryClient.invalidateQueries({ queryKey: githubKeys.all });
      await queryClient.invalidateQueries({ queryKey: ["portsvc", "snapshot"] });
      message.success(`${item.accountLogin} 已断开`);
    } catch (error) {
      message.error(error instanceof Error ? error.message : "断开失败");
    } finally {
      setBusyId(undefined);
    }
  }

  const installationColumns: ColumnsType<GitHubInstallation> = [
    {
      title: "账号",
      render: (_, item) => (
        <Space>
          <Avatar size="small" src={item.avatarUrl} icon={<GithubOutlined />} />
          <Typography.Text strong>{item.accountLogin}</Typography.Text>
          <Tag>{item.accountType}</Tag>
        </Space>
      ),
    },
    {
      title: "仓库范围",
      dataIndex: "repositorySelection",
      width: 130,
      render: (value) => <Tag color={value === "all" ? "green" : "gold"}>{value === "all" ? "全部" : "部分"}</Tag>,
    },
    {
      title: "状态",
      dataIndex: "status",
      width: 110,
      render: (value) => <Tag color={value === "active" ? "green" : "red"}>{value}</Tag>,
    },
    { title: "最近同步", dataIndex: "lastSyncedAt", width: 190, render: formatTime },
    {
      title: "操作",
      width: 150,
      render: (_, item) => (
        <Space size={4}>
          <Tooltip title="同步仓库">
            <Button disabled={item.status !== "active"} icon={<ReloadOutlined />} loading={busyId === item.id} onClick={() => void sync(item)} />
          </Tooltip>
          <Tooltip title="管理 GitHub 安装">
            <Button
              icon={<SettingOutlined />}
              href={`https://github.com/settings/installations/${item.githubInstallationId}`}
              target="_blank"
            />
          </Tooltip>
          <Popconfirm title="断开此账号？" onConfirm={() => void disconnect(item)}>
            <Tooltip title="断开连接">
              <Button danger icon={<DisconnectOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const repositoryColumns: ColumnsType<GitHubRepository> = [
    {
      title: "仓库",
      dataIndex: "fullName",
      render: (value, item) => (
        <Space>
          <Typography.Link href={item.htmlUrl} target="_blank">{value}</Typography.Link>
          {item.fork ? <Tag>Fork</Tag> : null}
          {item.archived ? <Tag>Archived</Tag> : null}
        </Space>
      ),
    },
    { title: "可见性", dataIndex: "visibility", width: 110, render: (value) => <Tag>{value}</Tag> },
    { title: "默认分支", dataIndex: "defaultBranch", width: 150, render: (value) => value || "-" },
    {
      title: "状态",
      dataIndex: "state",
      width: 120,
      render: (value) => <Tag color={value === "active" ? "green" : "red"}>{value}</Tag>,
    },
    { title: "GitHub 更新", dataIndex: "remoteUpdatedAt", width: 190, render: formatTime },
  ];

  return (
    <WorkbenchPage
      title="GitHub 仓库"
      extra={
        <Button type="primary" icon={<LinkOutlined />} disabled={!enabled} onClick={() => void connect()}>
          连接 GitHub
        </Button>
      }
    >
      {!enabled ? <Alert showIcon type="warning" message="GitHub App 未启用" /> : null}
      {loadError ? <Alert showIcon type="error" message="GitHub 数据加载失败" description={loadError.message} /> : null}
      <WorkbenchPanel title="已连接账号">
        <Table
          rowKey="id"
          columns={installationColumns}
          dataSource={installationsQuery.data ?? []}
          loading={installationsQuery.isPending}
          pagination={false}
          locale={{ emptyText: "尚未连接 GitHub 账号" }}
          scroll={{ x: 800 }}
        />
      </WorkbenchPanel>
      <WorkbenchPanel
        title="同步仓库"
        extra={
          <Input.Search
            allowClear
            prefix={<SearchOutlined />}
            placeholder="搜索仓库"
            onSearch={setQuery}
            style={{ width: 260 }}
          />
        }
      >
        <Table
          rowKey="id"
          columns={repositoryColumns}
          dataSource={repositoriesQuery.data ?? []}
          loading={repositoriesQuery.isPending}
          pagination={{ pageSize: 20 }}
          scroll={{ x: 900 }}
        />
      </WorkbenchPanel>
    </WorkbenchPage>
  );
}
