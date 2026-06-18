/**
 * 主题切换 Hook
 *
 * 支持亮色/暗色主题切换
 * 自动将主题 class 应用到 document 根元素
 */

import { useEffect, useMemo } from 'react';
import useAppStore from '@/store/appStore';

// ── 类型定义 ─────────────────────────────────────────────────

interface UseThemeReturn {
  /** 当前主题模式 */
  theme: 'light' | 'dark';
  /** 是否为暗色模式 */
  isDark: boolean;
  /** 切换主题 */
  toggleTheme: () => void;
  /** 设置指定主题 */
  setTheme: (theme: 'light' | 'dark') => void;
}

// ── Hook 实现 ────────────────────────────────────────────────

const useTheme = (): UseThemeReturn => {
  const { theme, toggleTheme, setTheme } = useAppStore();

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, [theme]);

  return useMemo(
    () => ({
      theme,
      isDark: theme === 'dark',
      toggleTheme,
      setTheme,
    }),
    [theme, toggleTheme, setTheme]
  );
};

export default useTheme;
