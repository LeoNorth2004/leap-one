/**
 * 认证服务 API
 *
 * 提供登录、登出、Token 刷新、用户信息获取等认证相关接口
 */

import { post, get } from './client';
import type { LoginParams, LoginResult, TokenInfo } from '@/types/auth';

const BASE = '/auth';

// ── 接口实现 ─────────────────────────────────────────────────

/** 用户登录 */
export const loginApi = (data: LoginParams): Promise<LoginResult> =>
  post<LoginResult>(`${BASE}/login`, data).then((res) => res.data);

/** 用户登出 */
export const logoutApi = (): Promise<void> =>
  post<void>(`${BASE}/logout`).then(() => undefined);

/** 刷新访问令牌 */
export const refreshTokenApi = (refreshToken: string): Promise<TokenInfo> =>
  post<TokenInfo>(`${BASE}/refresh`, { refreshToken }).then((res) => res.data);

/** 获取当前用户信息 */
export const fetchUserProfileApi = (): Promise<LoginResult['user']> =>
  get<LoginResult['user']>(`${BASE}/profile`).then((res) => res.data);

/** 获取验证码（图片验证码） */
export const getCaptchaApi = (): Promise<{ captchaId: string; image: string }> =>
  get<{ captchaId: string; image: string }>(`${BASE}/captcha`).then((res) => res.data);
