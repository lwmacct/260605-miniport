import { apiGet } from "../../shared/api/client";

export interface AdminUser {
  id: number;
  username: string;
  displayName: string;
  status: string;
  admin: boolean;
  disabledAt?: string;
}

export function fetchAdminUsers() {
  return apiGet<AdminUser[]>("/api/admin/users");
}
