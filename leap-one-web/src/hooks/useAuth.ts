/**
 * 认证 Hook
 *
 * 封装登录/登出逻辑，自动处理路由跳转
 */

import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store/authStore';
import type { LoginParams, UserInfo } from '@/types/auth';

// ── 类型定义 ─────────────────────────────────────────────────

interface UseAuthReturn {
  /** 当前用户信息 */
  user: UserInfo | null;
  /** 访问令牌 */
  token: string | null;
  /** 是否已认证 */
  isAuthenticated: boolean;
  /** 加载中状态 */
  isLoading: boolean;
  /** 执行登录（成功后自动跳转到工作台） */
  login: (params: LoginParams) => Promise<void>;
  /** 执行登出（成功后自动跳转到登录页） */
  logout: () => Promise<void>;
  /** 刷新访问令牌 */
  refreshToken: () => Promise<void>;
  /** 设置用户信息 */
  setUser: (user: UserInfo) => void;
  /** 手动设置 Tokens */
  setTokens: (access: string, refresh: string) => void;
  /** 获取当前用户信息 */
  fetchUserProfile: () => Promise<void>;
}

// ── Hook 实现 ────────────────────────────────────────────────

const useAuth = (): UseAuthReturn => {
  const navigate = useNavigate();
  const store = useAuthStore();

  const loginHandler = useCallback(
    async (params: LoginParams): Promise<void> => {
      await store.login(params);
      navigate('/', { replace: true });
    },
    [store.login, navigate]
  );

  const logoutHandler = useCallback(async (): Promise<void> => {
    await store.logout();
    navigate('/login', { replace: true });
  }, [store.logout, navigate]);

  return {
    user: store.user,
    token: store.token,
    isAuthenticated: store.isAuthenticated,
    isLoading: store.isLoading,
    login: loginHandler,
    logout: logoutHandler,
    refreshToken: store.refreshAccessToken,
    setUser: store.setUser,
    setTokens: store.setTokens,
    fetchUserProfile: store.fetchUserProfile,
  };
};

export default useAuth;
