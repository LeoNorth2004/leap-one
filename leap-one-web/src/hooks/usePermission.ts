/** 权限检查Hook */

import { useMemo, useCallback } from 'react';
import { useAuthStore } from '@/store/authStore';
import { useUserStore } from '@/store/userStore';
import { hasPermission, hasAnyPermission, hasAllPermissions, isAdmin } from '@/utils/permission';

export function usePermission() {
  const { user } = useAuthStore();
  const { permissions } = useUserStore();

  /** 检查单个权限 */
  const checkPermission = useCallback(
    (permission: string): boolean => {
      // 管理员拥有所有权限
      if (isAdmin(user)) return true;
      return hasPermission(permissions, permission);
    },
    [user, permissions]
  );

  /** 检查是否拥有任一权限（OR） */
  const checkAnyPermission = useCallback(
    (requiredPermissions: string[]): boolean => {
      if (isAdmin(user)) return true;
      return hasAnyPermission(permissions, requiredPermissions);
    },
    [user, permissions]
  );

  /** 检查是否拥有全部权限（AND） */
  const checkAllPermissions = useCallback(
    (requiredPermissions: string[]): boolean => {
      if (isAdmin(user)) return true;
      return hasAllPermissions(permissions, requiredPermissions);
    },
    [user, permissions]
  );

  /** 是否为管理员 */
  const admin = useMemo(() => isAdmin(user), [user]);

  return {
    permissions,
    isAdmin: admin,
    hasPermission: checkPermission,
    hasAnyPermission: checkAnyPermission,
    hasAllPermissions: checkAllPermissions,
  };
}
