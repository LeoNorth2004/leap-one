/** 主题切换Hook */

import { useEffect } from 'react';
import { useAppStore } from '@/store/appStore';

/**
 * 主题Hook - 支持亮色/暗色主题切换
 * 自动将主题class应用到document根元素
 */
export function useTheme() {
  const { theme, toggleTheme, setTheme } = useAppStore();

  useEffect(() => {
    // 将主题class应用到html根元素，供Ant Design和全局样式使用
    document.documentElement.setAttribute('data-theme', theme);
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, [theme]);

  return {
    theme,
    isDark: theme === 'dark',
    toggleTheme,
    setTheme,
  };
}
