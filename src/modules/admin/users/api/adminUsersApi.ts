import { apiGet } from "@/shared/api/client";

export interface AdminUser {
  id: number;
  username: string;
  displayName: string;
  role: string;
  status: string;
  admin: boolean;
  disabledAt?: string;
}

interface AdminUserList {
  items: AdminUser[];
}

export async function fetchAdminUsers() {
  const result = await apiGet<AdminUserList>("/api/admin/users?pageSize=100");
  return result?.items ?? [];
}
