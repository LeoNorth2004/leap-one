/**
 * 权限检查 Hook
 *
 * 基于当前用户角色和权限列表提供权限判断能力
 * 管理员默认拥有所有权限
 */

import { useMemo, useCallback } from 'react';
import { useAuthStore } from '@/store/authStore';
import { useUserStore } from '@/store/userStore';
import { hasPermission as checkPerm, hasAnyPermission as checkAnyPerm, hasAllPermissions as checkAllPerm, isAdmin } from '@/utils/permission';

// ── 类型定义 ─────────────────────────────────────────────────

interface UsePermissionReturn {
  /** 用户权限列表 */
  permissions: string[];
  /** 是否为管理员 */
  isAdmin: boolean;
  /** 检查单个权限 */
  hasPermission: (permission: string) => boolean;
  /** OR 逻辑：检查是否拥有任一权限 */
  hasAnyPermission: (requiredPermissions: string[]) => boolean;
  /** AND 逻辑：检查是否拥有全部权限 */
  hasAllPermissions: (requiredPermissions: string[]) => boolean;
}

// ── Hook 实现 ────────────────────────────────────────────────

const usePermission = (): UsePermissionReturn => {
  const { user } = useAuthStore();
  const { permissions } = useUserStore();

  // 单个权限检查
  const hasPermission = useCallback(
    (permission: string): boolean => {
      if (isAdmin(user)) return true;
      return checkPerm(permissions, permission);
    },
    [user, permissions]
  );

  // OR 权限检查
  const hasAnyPermission = useCallback(
    (requiredPermissions: string[]): boolean => {
      if (isAdmin(user)) return true;
      return checkAnyPerm(permissions, requiredPermissions);
    },
    [user, permissions]
  );

  // AND 权限检查
  const hasAllPermissions = useCallback(
    (requiredPermissions: string[]): boolean => {
      if (isAdmin(user)) return true;
      return checkAllPerm(permissions, requiredPermissions);
    },
    [user, permissions]
  );

  // 管理员标识
  const adminFlag = useMemo(() => isAdmin(user), [user]);

  return {
    permissions,
    isAdmin: adminFlag,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  };
};

export default usePermission;
