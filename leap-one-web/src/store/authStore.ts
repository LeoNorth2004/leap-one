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

interface AuthState {
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

  /** 用户登录 */
  login: (params: LoginParams) => Promise<LoginResult>;
  /** 用户登出 */
  logout: () => Promise<void>;
  /** 设置 Tokens（手动设置时使用） */
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

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // ─── 初始状态 ─────────────────────────────────────────────
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      // ─── login ──────────────────────────────────────────────────
      login: async (params: LoginParams): Promise<LoginResult> => {
        set({ isLoading: true });
        try {
          const result = await loginApi(params);

          // 持久化 Token 到 localStorage
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

      // ─── logout ─────────────────────────────────────────────────
      logout: async (): Promise<void> => {
        try {
          await logoutApi();
        } finally {
          // 无论 API 是否成功，都清除本地状态
          get().clearAuth();
        }
      },

      // ─── setTokens ─────────────────────────────────────────────
      setTokens: (access: string, refresh: string): void => {
        tokenStorage.setToken(access);
        tokenStorage.setRefreshToken(refresh);
        set({
          token: access,
          refreshToken: refresh,
          isAuthenticated: true,
        });
      },

      // ─── fetchUserProfile ──────────────────────────────────────
      fetchUserProfile: async (): Promise<void> => {
        try {
          const user = await fetchUserProfileApi();
          set({ user });
        } catch (error) {
          console.error('[Auth] 获取用户信息失败:', error);
          throw error;
        }
      },

      // ─── refreshAccessToken ────────────────────────────────────
      refreshAccessToken: async (): Promise<void> => {
        const currentRefreshToken = get().refreshToken ?? tokenStorage.getRefreshToken();
        if (!currentRefreshToken) {
          get().clearAuth();
          return;
        }

        try {
          const result = await refreshTokenApi(currentRefreshToken);
          tokenStorage.setToken(result.accessToken);
          set({
            token: result.accessToken,
            isAuthenticated: true,
          });
        } catch {
          // 刷新失败，清除认证状态
          get().clearAuth();
        }
      },

      // ─── setUser ───────────────────────────────────────────────
      setUser: (user: UserInfo): void => {
        set({ user });
      },

      // ─── clearAuth ─────────────────────────────────────────────
      clearAuth: (): void => {
        tokenStorage.clearAll();
        set({
          user: null,
          token: null,
          refreshToken: null,
          isAuthenticated: false,
        });
      },
    }),
    {
      name: 'leap-one-auth', // localStorage key
      partialize: (state) => ({
        // 只持久化这些字段（不包含 isLoading）
        token: state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
        user: state.user,
      }),
    }
  )
);
