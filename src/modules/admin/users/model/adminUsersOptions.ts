export const adminUserRoleOptions = [
  { label: "管理员", value: "admin" },
  { label: "普通用户", value: "user" },
] as const;

export const adminUserStatusOptions = [
  { label: "启用", value: "active" },
  { label: "禁用", value: "disabled" },
] as const;

export type AdminUserRole = (typeof adminUserRoleOptions)[number]["value"];
export type AdminUserStatus = (typeof adminUserStatusOptions)[number]["value"];
