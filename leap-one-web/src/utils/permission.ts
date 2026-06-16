/** 权限工具函数 */

import type { UserInfo } from '@/types/auth';

/** 权限码前缀 */
const PERMISSION_PREFIXES = {
  VIEW: 'view:',
  CREATE: 'create:',
  EDIT: 'edit:',
  DELETE: 'delete:',
  EXPORT: 'export:',
  MANAGE: 'manage:',
} as const;

/**
 * 检查用户是否拥有指定权限
 * @param permissions - 用户权限列表
 * @param permission - 需要检查的权限码
 */
export function hasPermission(permissions: string[], permission: string): boolean {
  return permissions.includes(permission);
}

/**
 * 检查是否拥有任一权限（OR逻辑）
 */
export function hasAnyPermission(permissions: string[], requiredPermissions: string[]): boolean {
  return requiredPermissions.some((p) => permissions.includes(p));
}

/**
 * 检查是否拥有全部权限（AND逻辑）
 */
export function hasAllPermissions(permissions: string[], requiredPermissions: string[]): boolean {
  return requiredPermissions.every((p) => permissions.includes(p));
}

/**
 * 检查是否为管理员
 */
export function isAdmin(user: UserInfo | null): boolean {
  if (!user) return false;
  return user.roles.includes('admin') || user.roles.includes('super_admin');
}

/**
 * 检查是否为项目成员
 */
export function isProjectMember(projectMembers: number[], userId: number): boolean {
  return projectMembers.includes(userId);
}

/**
 * 根据资源生成查看权限码
 */
export function viewPermission(resource: string): string {
  return `${PERMISSION_PREFIXES.VIEW}${resource}`;
}

/**
 * 根据资源生成编辑权限码
 */
export function editPermission(resource: string): string {
  return `${PERMISSION_PREFIXES.EDIT}${resource}`;
}
