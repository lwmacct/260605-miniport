import { useQuery } from "@tanstack/react-query";
import { fetchAdminUsers } from "../api/adminUsersApi";

export const adminUsersKeys = {
  list: ["admin-users", "list"] as const,
};

export function useAdminUsersQuery() {
  return useQuery({
    queryKey: adminUsersKeys.list,
    queryFn: fetchAdminUsers,
  });
}
