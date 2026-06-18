/**
 * 应用全局状态 - Zustand Store
 *
 * 管理主题模式、侧边栏折叠状态、语言设置等全局 UI 偏好
 */

import { create } from 'zustand';
import type { ThemeMode, LocaleType } from '@/types/common';
import { prefStorage } from '@/utils/storage';

// ── 类型定义 ─────────────────────────────────────────────────

interface AppState {
  /** 主题模式：亮色/暗色 */
  theme: ThemeMode;
  /** 侧边栏是否折叠 */
  sidebarCollapsed: boolean;
  /** 语言设置 */
  locale: LocaleType;

  // ── Actions ────────────────────────────────────────────────

  /** 切换主题（亮色 <-> 暗色） */
  toggleTheme: () => void;
  /** 设置指定主题 */
  setTheme: (theme: ThemeMode) => void;
  /** 切换侧边栏折叠状态 */
  toggleSidebar: () => void;
  /** 设置侧边栏折叠状态 */
  setSidebarCollapsed: (collapsed: boolean) => void;
  /** 设置语言 */
  setLocale: (locale: LocaleType) => void;
}

// ── 初始状态工厂 ─────────────────────────────────────────────

const createInitialState = (): Pick<AppState, 'theme' | 'sidebarCollapsed' | 'locale'> => ({
  theme: prefStorage.getTheme(),
  sidebarCollapsed: prefStorage.getSidebarCollapsed(),
  locale: 'zh-CN',
});

// ── Store 创建 ───────────────────────────────────────────────

const useAppStore = create<AppState>((set) => ({
  ...createInitialState(),

  toggleTheme() {
    set((state) => {
      const nextTheme = state.theme === 'light' ? 'dark' : 'light';
      prefStorage.setTheme(nextTheme);
      return { theme: nextTheme };
    });
  },

  setTheme(theme: ThemeMode) {
    prefStorage.setTheme(theme);
    set({ theme });
  },

  toggleSidebar() {
    set((state) => {
      const nextCollapsed = !state.sidebarCollapsed;
      prefStorage.setSidebarCollapsed(nextCollapsed);
      return { sidebarCollapsed: nextCollapsed };
    });
  },

  setSidebarCollapsed(collapsed: boolean) {
    prefStorage.setSidebarCollapsed(collapsed);
    set({ sidebarCollapsed: collapsed });
  },

  setLocale(locale: LocaleType) {
    set({ locale });
  },
}));

export default useAppStore;
