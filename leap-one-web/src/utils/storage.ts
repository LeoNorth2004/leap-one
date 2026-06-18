// 本地存储管理工具

// ── 存储键名常量 ──────────────────────────────────────────────
const STORAGE_KEYS = {
  TOKEN: 'leap_one_token',
  REFRESH_TOKEN: 'leap_one_refresh_token',
  USER: 'leap_one_user',
  THEME: 'leap_one_theme',
  SIDEBAR: 'leap_one_sidebar_collapsed',
} as const;

// ── 基础存储操作 ──────────────────────────────────────────────

const readFromStorage = <T>(key: string): T | null => {
  try {
    const rawValue = localStorage.getItem(key);
    return rawValue ? (JSON.parse(rawValue) as T) : null;
  } catch {
    return null;
  }
};

const writeToStorage = <T>(key: string, value: T): void => {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (err) {
    console.error('LocalStorage 写入失败:', err);
  }
};

const removeFromStorage = (key: string): void => {
  localStorage.removeItem(key);
};

export { readFromStorage as getStorage, writeToStorage as setStorage, removeFromStorage as removeStorage };

// ── Token 存储 ────────────────────────────────────────────────

export const tokenStorage = {
  getToken(): string | null {
    return readFromStorage<string>(STORAGE_KEYS.TOKEN);
  },

  setToken(token: string): void {
    writeToStorage(STORAGE_KEYS.TOKEN, token);
  },

  removeToken(): void {
    removeFromStorage(STORAGE_KEYS.TOKEN);
  },

  getRefreshToken(): string | null {
    return readFromStorage<string>(STORAGE_KEYS.REFRESH_TOKEN);
  },

  setRefreshToken(token: string): void {
    writeToStorage(STORAGE_KEYS.REFRESH_TOKEN, token);
  },

  removeRefreshToken(): void {
    removeFromStorage(STORAGE_KEYS.REFRESH_TOKEN);
  },

  clearAll(): void {
    removeFromStorage(STORAGE_KEYS.TOKEN);
    removeFromStorage(STORAGE_KEYS.REFRESH_TOKEN);
    removeFromStorage(STORAGE_KEYS.USER);
  },
};

// ── 用户信息存储 ──────────────────────────────────────────────

export const userStorage = {
  getUser(): Record<string, unknown> | null {
    return readFromStorage<Record<string, unknown>>(STORAGE_KEYS.USER);
  },

  setUser(user: Record<string, unknown>): void {
    writeToStorage(STORAGE_KEYS.USER, user);
  },

  removeUser(): void {
    removeFromStorage(STORAGE_KEYS.USER);
  },
};

// ── UI 偏好设置存储 ───────────────────────────────────────────

export const prefStorage = {
  getTheme(): 'light' | 'dark' {
    const stored = readFromStorage<'light' | 'dark'>(STORAGE_KEYS.THEME);
    return stored || 'light';
  },

  setTheme(theme: 'light' | 'dark'): void {
    writeToStorage(STORAGE_KEYS.THEME, theme);
  },

  getSidebarCollapsed(): boolean {
    const stored = readFromStorage<boolean>(STORAGE_KEYS.SIDEBAR);
    return stored || false;
  },

  setSidebarCollapsed(collapsed: boolean): void {
    writeToStorage(STORAGE_KEYS.SIDEBAR, collapsed);
  },
};
