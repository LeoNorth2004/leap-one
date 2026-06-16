/** 认证相关 API */

import { post, get } from './client';
import type { LoginParams, LoginResult, TokenInfo } from '@/types/auth';

const BASE_URL = '/auth';

/** 用户登录 */
export function loginApi(data: LoginParams): Promise<LoginResult> {
  return post<LoginResult>(`${BASE_URL}/login`, data);
}

/** 用户登出 */
export function logoutApi(): Promise<void> {
  return post<void>(`${BASE_URL}/logout`).then(() => undefined);
}

/** 刷新 Token */
export function refreshTokenApi(refreshToken: string): Promise<TokenInfo> {
  return post<TokenInfo>(`${BASE_URL}/refresh`, { refreshToken });
}

/** 获取当前用户信息 */
export function fetchUserProfileApi(): Promise<LoginResult['user']> {
  return get<LoginResult['user']>(`${BASE_URL}/profile`);
}

/** 获取验证码 */
export function getCaptchaApi(): Promise<{ captchaId: string; image: string }> {
  return get<{ captchaId: string; image: string }>(`${BASE_URL}/captcha`);
}
