/** 用户状态管理 - Zustand Store */

import { create } from 'zustand';
import type { UserDetail } from '@/types/user';

interface UserState {
  /** 当前用户详情（扩展信息） */
  currentUser: UserDetail | null;
  /** 用户权限列表 */
  permissions: string[];

  /** 设置当前用户 */
  setCurrentUser: (user: UserDetail | null) => void;
  /** 设置权限列表 */
  setPermissions: (permissions: string[]) => void;
  /** 更新用户部分信息 */
  updateUserInfo: (data: Partial<UserDetail>) => void;
  /** 清除用户状态 */
  clearUser: () => void;
}

export const useUserStore = create<UserState>((set) => ({
  currentUser: null,
  permissions: [],

  setCurrentUser: (user: UserDetail | null) =>
    set({ currentUser: user }),

  setPermissions: (permissions: string[]) =>
    set({ permissions }),

  updateUserInfo: (data: Partial<UserDetail>) =>
    set((state) => ({
      currentUser: state.currentUser ? { ...state.currentUser, ...data } : null,
    })),

  clearUser: () =>
    set({ currentUser: null, permissions: [] }),
}));
