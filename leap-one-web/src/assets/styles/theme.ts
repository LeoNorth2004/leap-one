/** 主题配置 - 支持亮色/暗色模式切换 */

export const lightTheme = {
  name: 'light' as const,
  token: {
    colorPrimary: '#1677ff',
    colorBgContainer: '#ffffff',
    colorBgLayout: '#f5f7fa',
    colorText: '#1d2129',
    colorTextSecondary: '#4e5969',
    colorBorder: '#e5e6eb',
    borderRadius: 6,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif',
  },
};

export const darkTheme = {
  name: 'dark' as const,
  token: {
    colorPrimary: '#1677ff',
    colorBgContainer: '#1f1f1f',
    colorBgLayout: '#141414',
    colorText: '#e8e8e8',
    colorTextSecondary: '#a6a6a6',
    colorBorder: '#303030',
    borderRadius: 6,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif',
  },
};

/** 根据主题名称获取配置 */
export function getThemeConfig(themeName: 'light' | 'dark') {
  return themeName === 'dark' ? darkTheme : lightTheme;
}
