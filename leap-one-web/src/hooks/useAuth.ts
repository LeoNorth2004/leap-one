/**
 * 认证相关 Hook
 *
 * 从 authStore 获取状态和方法，封装登录/登出逻辑并自动处理路由跳转
 */

import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store/authStore';
import type { LoginParams } from '@/types/auth';

interface UseAuthReturn {
  /** 当前用户信息 */
  user: ReturnType<typeof useAuthStore>['user'];
  /** 访问令牌 */
  token: ReturnType<typeof useAuthStore>['token'];
  /** 是否已认证 */
  isAuthenticated: ReturnType<typeof useAuthStore>['isAuthenticated'];
  /** 加载中状态 */
  isLoading: ReturnType<typeof useAuthStore>['isLoading'];
  /** 执行登录（成功后自动跳转到工作台） */
  login: (params: LoginParams) => Promise<void>;
  /** 执行登出（成功后自动跳转到登录页） */
  logout: () => Promise<void>;
  /** 刷新访问令牌 */
  refreshToken: ReturnType<typeof useAuthStore>['refreshAccessToken'];
  /** 设置用户信息 */
  setUser: ReturnType<typeof useAuthStore>['setUser'];
  /** 手动设置 Tokens */
  setTokens: ReturnType<typeof useAuthStore>['setTokens'];
  /** 获取当前用户信息 */
  fetchUserProfile: ReturnType<typeof useAuthStore>['fetchUserProfile'];
}

export function useAuth(): UseAuthReturn {
  const navigate = useNavigate();
  const {
    user,
    token,
    isAuthenticated,
    isLoading,
    login: storeLogin,
    logout: storeLogout,
    refreshAccessToken,
    setUser,
    setTokens,
    fetchUserProfile,
  } = useAuthStore();

  /** 执行登录并跳转到工作台 */
  const login = useCallback(
    async (params: LoginParams) => {
      await storeLogin(params);
      navigate('/', { replace: true });
    },
    [storeLogin, navigate]
  );

  /** 执行登出并跳转到登录页 */
  const logout = useCallback(async () => {
    await storeLogout();
    navigate('/login', { replace: true });
  }, [storeLogout, navigate]);

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    logout,
    refreshToken: refreshAccessToken,
    setUser,
    setTokens,
    fetchUserProfile,
  };
}
