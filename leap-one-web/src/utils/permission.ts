// 权限校验工具

import type { UserInfo } from '@/types/auth';

// ── 权限码前缀定义 ───────────────────────────────────────────
const PERMISSION_PREFIXES = Object.freeze({
  VIEW: 'view:',
  CREATE: 'create:',
  EDIT: 'edit:',
  DELETE: 'delete:',
  EXPORT: 'export:',
  MANAGE: 'manage:',
});

// ── 权限检查核心方法 ─────────────────────────────────────────

/** 校验用户是否拥有指定权限 */
export const hasPermission = (permissions: string[], permission: string): boolean =>
  permissions.includes(permission);

/** OR 逻辑：校验是否拥有任一权限 */
export const hasAnyPermission = (permissions: string[], requiredPermissions: string[]): boolean =>
  requiredPermissions.some((p) => permissions.includes(p));

/** AND 逻辑：校验是否拥有全部权限 */
export const hasAllPermissions = (permissions: string[], requiredPermissions: string[]): boolean =>
  requiredPermissions.every((p) => permissions.includes(p));

/** 判断当前用户是否为管理员 */
export const isAdmin = (user: UserInfo | null): boolean => {
  if (!user) {
    return false;
  }
  const adminRoles = ['admin', 'super_admin'];
  return adminRoles.some((role) => user.roles.includes(role));
};

/** 判断用户是否为项目成员 */
export const isProjectMember = (projectMembers: number[], userId: number): boolean =>
  projectMembers.includes(userId);

// ── 权限码生成器 ─────────────────────────────────────────────

/** 生成查看权限码 */
export const viewPermission = (resource: string): string =>
  `${PERMISSION_PREFIXES.VIEW}${resource}`;

/** 生成编辑权限码 */
export const editPermission = (resource: string): string =>
  `${PERMISSION_PREFIXES.EDIT}${resource}`;
