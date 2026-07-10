import { WorkbenchPanel } from "@lwmacct/260627-antd-workbench";
import { Alert, Space, Typography } from "antd";
import { useMemo, useState } from "react";
import type { AdminUser } from "../api/adminUsersApi";
import { useAdminUsersQuery } from "../model/adminUsersQueries";
import styles from "@/shared/ui/SectionPage.module.css";
import { AdminUsersTable } from "./AdminUsersTable";
import { AdminUsersToolbar, type AdminUsersFilters } from "./AdminUsersToolbar";

export function AdminUsersPage() {
  const [filters, setFilters] = useState<AdminUsersFilters>({});
  const users = useAdminUsersQuery();
  const filteredUsers = useMemo(
    () => filterUsers(users.data ?? [], filters),
    [filters, users.data],
  );

  function updateFilters(nextFilters: Partial<AdminUsersFilters>) {
    setFilters((current) => ({ ...current, ...nextFilters }));
  }

  return (
    <section className={styles.section}>
      <div className={styles.sectionHeader}>
        <Typography.Title level={2}>用户管理</Typography.Title>
        <Typography.Paragraph type="secondary">
          查看用户状态与管理员权限，管理员权限来自 server.auth.admins 运行时配置。
        </Typography.Paragraph>
      </div>

      {users.isError ? <Alert showIcon type="error" message={users.error.message} /> : null}

      <WorkbenchPanel>
        <Space orientation="vertical" size="middle" style={{ width: "100%" }}>
          <AdminUsersToolbar
            loading={users.isFetching}
            onFiltersChange={updateFilters}
            onRefresh={() => void users.refetch()}
          />
          <AdminUsersTable data={filteredUsers} loading={users.isPending} />
        </Space>
      </WorkbenchPanel>
    </section>
  );
}

function filterUsers(users: AdminUser[], filters: AdminUsersFilters) {
  const keyword = filters.keyword?.trim().toLowerCase();

  return users.filter((user) => {
    const role = user.admin ? "admin" : "user";
    if (filters.role && role !== filters.role) {
      return false;
    }
    if (filters.status && user.status !== filters.status) {
      return false;
    }
    if (!keyword) {
      return true;
    }
    return (
      user.username.toLowerCase().includes(keyword) ||
      user.displayName.toLowerCase().includes(keyword)
    );
  });
}
