/** 主题配置 - 支持亮色/暗色模式切换 */

export const lightTheme = {
  name: 'light' as const,
  token: {
    colorPrimary: '#1677ff',
    colorPrimaryHover: '#4096ff',
    colorPrimaryActive: '#0958d9',
    colorBgContainer: '#ffffff',
    colorBgLayout: '#f0f2f5',
    colorBgPage: '#f5f7fa',
    colorText: '#1d2129',
    colorTextSecondary: '#4e5969',
    colorTextPlaceholder: '#86909c',
    colorTextDisabled: '#c9cdd4',
    colorBorder: '#e5e6eb',
    colorBorderHover: '#d0d3d9',
    colorSplit: '#f0f0f0',
    colorBgElevated: '#ffffff',
    colorBgHover: '#f5f7fa',
    colorBgActive: '#e8eaed',
    borderRadius: 8,
    borderRadiusLG: 12,
    borderRadiusSM: 4,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif',
    fontWeightStrong: 600,
    shadowColor: 'rgba(0, 0, 0, 0.08)',
    shadowColorSecondary: 'rgba(0, 0, 0, 0.12)',
    shadowColorTertiary: 'rgba(0, 0, 0, 0.16)',
    motionDurationSlow: 300,
    motionDurationMid: 200,
    motionDurationFast: 100,
    motionEaseInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    motionEaseOut: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
  },
};

export const darkTheme = {
  name: 'dark' as const,
  token: {
    colorPrimary: '#4096ff',
    colorPrimaryHover: '#69b1ff',
    colorPrimaryActive: '#1677ff',
    colorBgContainer: '#1f1f1f',
    colorBgLayout: '#0a0a0a',
    colorBgPage: '#141414',
    colorText: '#f5f5f5',
    colorTextSecondary: '#bfbfbf',
    colorTextPlaceholder: '#666666',
    colorTextDisabled: '#434343',
    colorBorder: '#303030',
    colorBorderHover: '#434343',
    colorSplit: '#262626',
    colorBgElevated: '#292929',
    colorBgHover: '#262626',
    colorBgActive: '#303030',
    borderRadius: 8,
    borderRadiusLG: 12,
    borderRadiusSM: 4,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif',
    fontWeightStrong: 600,
    shadowColor: 'rgba(0, 0, 0, 0.3)',
    shadowColorSecondary: 'rgba(0, 0, 0, 0.4)',
    shadowColorTertiary: 'rgba(0, 0, 0, 0.5)',
    motionDurationSlow: 300,
    motionDurationMid: 200,
    motionDurationFast: 100,
    motionEaseInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    motionEaseOut: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
  },
};

export const themeColors = {
  primary: '#1677ff',
  success: '#52c41a',
  warning: '#faad14',
  error: '#ff4d4f',
  info: '#1677ff',
  cyan: '#13c2c2',
  purple: '#722ed1',
  orange: '#fa8c16',
  pink: '#eb2f96',
  red: '#ff4d4f',
  green: '#52c41a',
  blue: '#1677ff',
};

export const gradientColors = {
  primary: 'linear-gradient(135deg, #1677ff 0%, #13c2c2 100%)',
  success: 'linear-gradient(135deg, #52c41a 0%, #73d13d 100%)',
  warning: 'linear-gradient(135deg, #faad14 0%, #ffc53d 100%)',
  error: 'linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%)',
  dark: 'linear-gradient(180deg, #1f1f1f 0%, #141414 100%)',
  light: 'linear-gradient(180deg, #ffffff 0%, #f5f7fa 100%)',
  brand: 'linear-gradient(135deg, #1677ff 0%, #0958d9 50%, #13c2c2 100%)',
};

export function getThemeConfig(themeName: 'light' | 'dark') {
  return themeName === 'dark' ? darkTheme : lightTheme;
}
