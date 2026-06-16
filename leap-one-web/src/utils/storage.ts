/** 本地存储封装工具 */

const TOKEN_KEY = 'leap_one_token';
const REFRESH_TOKEN_KEY = 'leap_one_refresh_token';
const USER_KEY = 'leap_one_user';
const THEME_KEY = 'leap_one_theme';
const SIDEBAR_KEY = 'leap_one_sidebar_collapsed';

/**
 * 获取存储项
 */
export function getStorage<T>(key: string): T | null {
  try {
    const item = localStorage.getItem(key);
    return item ? (JSON.parse(item) as T) : null;
  } catch {
    return null;
  }
}

/**
 * 设置存储项
 */
export function setStorage<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error('LocalStorage 写入失败:', error);
  }
}

/**
 * 移除存储项
 */
export function removeStorage(key: string): void {
  localStorage.removeItem(key);
}

/** Token 相关操作 */
export const tokenStorage = {
  getToken: () => getStorage<string>(TOKEN_KEY),
  setToken: (token: string) => setStorage(TOKEN_KEY, token),
  removeToken: () => removeStorage(TOKEN_KEY),

  getRefreshToken: () => getStorage<string>(REFRESH_TOKEN_KEY),
  setRefreshToken: (token: string) => setStorage(REFRESH_TOKEN_KEY, token),
  removeRefreshToken: () => removeStorage(REFRESH_TOKEN_KEY),

  clearAll: () => {
    removeStorage(TOKEN_KEY);
    removeStorage(REFRESH_TOKEN_KEY);
    removeStorage(USER_KEY);
  },
};

/** 用户信息存储 */
export const userStorage = {
  getUser: () => getStorage<Record<string, unknown>>(USER_KEY),
  setUser: (user: Record<string, unknown>) => setStorage(USER_KEY, user),
  removeUser: () => removeStorage(USER_KEY),
};

/** UI偏好存储 */
export const prefStorage = {
  getTheme: (): 'light' | 'dark' => getStorage<'light' | 'dark'>(THEME_KEY) || 'light',
  setTheme: (theme: 'light' | 'dark') => setStorage(THEME_KEY, theme),

  getSidebarCollapsed: (): boolean => getStorage<boolean>(SIDEBAR_KEY) || false,
  setSidebarCollapsed: (collapsed: boolean) => setStorage(SIDEBAR_KEY, collapsed),
};
