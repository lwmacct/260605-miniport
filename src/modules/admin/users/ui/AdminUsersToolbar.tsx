import { ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { Button, Input, Select } from "antd";
import { adminUserRoleOptions, adminUserStatusOptions } from "../model/adminUsersOptions";
import styles from "./AdminUsersToolbar.module.css";

export interface AdminUsersFilters {
  keyword?: string;
  role?: string;
  status?: string;
}

interface AdminUsersToolbarProps {
  loading?: boolean;
  onFiltersChange(filters: Partial<AdminUsersFilters>): void;
  onRefresh(): void;
}

export function AdminUsersToolbar({
  loading,
  onFiltersChange,
  onRefresh,
}: AdminUsersToolbarProps) {
  return (
    <div className={styles.toolbar}>
      <Input.Search
        allowClear
        className={styles.search}
        enterButton={<SearchOutlined />}
        placeholder="搜索用户名、显示名称"
        onSearch={(keyword) => onFiltersChange({ keyword })}
      />
      <Select
        allowClear
        className={styles.select}
        options={[...adminUserRoleOptions]}
        placeholder="角色"
        onChange={(role?: string) => onFiltersChange({ role })}
      />
      <Select
        allowClear
        className={styles.select}
        options={[...adminUserStatusOptions]}
        placeholder="状态"
        onChange={(status?: string) => onFiltersChange({ status })}
      />
      <Button disabled={loading} icon={<ReloadOutlined />} onClick={onRefresh}>
        刷新
      </Button>
    </div>
  );
}
