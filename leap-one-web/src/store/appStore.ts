/** 应用全局状态管理 - Zustand Store */

import { create } from 'zustand';
import type { ThemeMode, LocaleType } from '@/types/common';
import { prefStorage } from '@/utils/storage';

interface AppState {
  /** 主题模式：亮色/暗色 */
  theme: ThemeMode;
  /** 侧边栏是否折叠 */
  sidebarCollapsed: boolean;
  /** 语言设置 */
  locale: LocaleType;

  /** 切换主题 */
  toggleTheme: () => void;
  /** 设置主题 */
  setTheme: (theme: ThemeMode) => void;
  /** 切换侧边栏折叠状态 */
  toggleSidebar: () => void;
  /** 设置侧边栏折叠状态 */
  setSidebarCollapsed: (collapsed: boolean) => void;
  /** 设置语言 */
  setLocale: (locale: LocaleType) => void;
}

export const useAppStore = create<AppState>((set) => ({
  theme: prefStorage.getTheme(),
  sidebarCollapsed: prefStorage.getSidebarCollapsed(),
  locale: 'zh-CN',

  toggleTheme: () =>
    set((state) => {
      const newTheme = state.theme === 'light' ? 'dark' : 'light';
      prefStorage.setTheme(newTheme);
      return { theme: newTheme };
    }),

  setTheme: (theme: ThemeMode) => {
    prefStorage.setTheme(theme);
    set({ theme });
  },

  toggleSidebar: () =>
    set((state) => {
      const newValue = !state.sidebarCollapsed;
      prefStorage.setSidebarCollapsed(newValue);
      return { sidebarCollapsed: newValue };
    }),

  setSidebarCollapsed: (collapsed: boolean) => {
    prefStorage.setSidebarCollapsed(collapsed);
    set({ sidebarCollapsed: collapsed });
  },

  setLocale: (locale: LocaleType) => set({ locale }),
}));
