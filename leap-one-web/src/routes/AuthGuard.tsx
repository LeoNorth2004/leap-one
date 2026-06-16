/** 权限守卫组件 - 检查用户是否已登录 */

import { Navigate, useLocation } from 'react-router-dom';
import React from 'react';
import { useAuthStore } from '@/store/authStore';

interface AuthGuardProps {
  children: React.ReactNode;
}

/**
 * 认证守卫：
 * - 已登录 → 渲染子组件
 * - 未登录 → 重定向到登录页，并记录原始路径用于登录后跳回
 */
export default function AuthGuard({ children }: AuthGuardProps) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const location = useLocation();

  if (!isAuthenticated) {
    // 将当前路径保存到location state，以便登录后可以跳回来
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}
