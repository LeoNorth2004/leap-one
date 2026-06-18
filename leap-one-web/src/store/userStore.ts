/**
 * 用户状态管理 - Zustand Store
 *
 * 管理当前用户详情信息和权限列表
 */

import { create } from 'zustand';
import type { UserDetail } from '@/types/user';

// ── 类型定义 ─────────────────────────────────────────────────

interface UserState {
  // ── State ──────────────────────────────────────────────────
  /** 当前用户详情（扩展信息） */
  currentUser: UserDetail | null;
  /** 用户权限列表 */
  permissions: string[];

  // ── Actions ────────────────────────────────────────────────
  /** 设置当前用户 */
  setCurrentUser: (user: UserDetail | null) => void;
  /** 设置权限列表 */
  setPermissions: (permissions: string[]) => void;
  /** 更新用户部分信息 */
  updateUserInfo: (data: Partial<UserDetail>) => void;
  /** 清除用户状态 */
  clearUser: () => void;
}

// ── Store 创建 ───────────────────────────────────────────────

const useUserStore = create<UserState>((set) => ({
  currentUser: null,
  permissions: [],

  setCurrentUser(user: UserDetail | null): void {
    set({ currentUser: user });
  },

  setPermissions(permissions: string[]): void {
    set({ permissions });
  },

  updateUserInfo(data: Partial<UserDetail>): void {
    set((state) => ({
      currentUser: state.currentUser
        ? { ...state.currentUser, ...data }
        : null,
    }));
  },

  clearUser(): void {
    set({ currentUser: null, permissions: [] });
  },
}));

export default useUserStore;
