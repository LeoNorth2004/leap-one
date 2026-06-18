/**
 * 认证状态管理 - Zustand Store
 *
 * 功能：
 * - 使用 Zustand + persist 中间件持久化到 localStorage
 * - login 方法调用认证 API，成功后存储 token 和用户信息
 * - logout 清除所有状态并跳转登录页
 * - fetchUserProfile 获取当前用户信息
 * - Token 自动刷新机制
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { UserInfo, LoginParams, LoginResult } from '@/types/auth';
import { loginApi, logoutApi, refreshTokenApi, fetchUserProfileApi } from '@/api/auth';
import { tokenStorage } from '@/utils/storage';

// ── 类型定义 ─────────────────────────────────────────────────

interface AuthState {
  // ── State ──────────────────────────────────────────────────
  /** 当前用户信息 */
  user: UserInfo | null;
  /** 访问令牌 */
  token: string | null;
  /** 刷新令牌 */
  refreshToken: string | null;
  /** 是否已认证 */
  isAuthenticated: boolean;
  /** 登录中 / 请求中状态 */
  isLoading: boolean;

  // ── Actions ────────────────────────────────────────────────
  /** 用户登录 */
  login: (params: LoginParams) => Promise<LoginResult>;
  /** 用户登出 */
  logout: () => Promise<void>;
  /** 手动设置 Tokens */
  setTokens: (access: string, refresh: string) => void;
  /** 获取当前用户信息 */
  fetchUserProfile: () => Promise<void>;
  /** 刷新访问令牌 */
  refreshAccessToken: () => Promise<void>;
  /** 设置用户信息 */
  setUser: (user: UserInfo) => void;
  /** 清除所有认证状态（内部使用） */
  clearAuth: () => void;
}

// ── Persist 配置 ─────────────────────────────────────────────

const PERSIST_CONFIG = {
  name: 'leap-one-auth',
  partialize: (state: AuthState) => ({
    token: state.token,
    refreshToken: state.refreshToken,
    isAuthenticated: state.isAuthenticated,
    user: state.user,
  }),
};

// ── Store 创建 ───────────────────────────────────────────────

const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // ═══ 初始状态 ══════════════════════════════════════════
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      // ═══ login ══════════════════════════════════════════════
      login: async (params: LoginParams): Promise<LoginResult> => {
        set({ isLoading: true });
        try {
          const result = await loginApi(params);

          tokenStorage.setToken(result.token);
          tokenStorage.setRefreshToken(result.refreshToken);

          set({
            user: result.user,
            token: result.token,
            refreshToken: result.refreshToken,
            isAuthenticated: true,
            isLoading: false,
          });

          return result;
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      // ═══ logout ═════════════════════════════════════════════
      logout: async (): Promise<void> => {
        try {
          await logoutApi();
        } finally {
          get().clearAuth();
        }
      },

      // ═══ setTokens ══════════════════════════════════════════
      setTokens(access: string, refresh: string): void {
        tokenStorage.setToken(access);
        tokenStorage.setRefreshToken(refresh);
        set({
          token: access,
          refreshToken: refresh,
          isAuthenticated: true,
        });
      },

      // ═══ fetchUserProfile ═══════════════════════════════════
      fetchUserProfile: async (): Promise<void> => {
        try {
          const userInfo = await fetchUserProfileApi();
          set({ user: userInfo });
        } catch (error) {
          console.error('[Auth] 获取用户信息失败:', error);
          throw error;
        }
      },

      // ═══ refreshAccessToken ═════════════════════════════════
      refreshAccessToken: async (): Promise<void> => {
        const currentRefreshToken =
          get().refreshToken ?? tokenStorage.getRefreshToken();

        if (!currentRefreshToken) {
          get().clearAuth();
          return;
        }

        try {
          const tokenInfo = await refreshTokenApi(currentRefreshToken);
          tokenStorage.setToken(tokenInfo.accessToken);
          set({
            token: tokenInfo.accessToken,
            isAuthenticated: true,
          });
        } catch {
          get().clearAuth();
        }
      },

      // ═══ setUser ════════════════════════════════════════════
      setUser(user: UserInfo): void {
        set({ user });
      },

      // ═══ clearAuth ═════════════════════════════════════════
      clearAuth(): void {
        tokenStorage.clearAll();
        set({
          user: null,
          token: null,
          refreshToken: null,
          isAuthenticated: false,
        });
      },
    }),
    PERSIST_CONFIG
  )
);

export default useAuthStore;
