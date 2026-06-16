/** 认证相关类型定义 */

/** 登录请求参数 */
export interface LoginParams {
  username: string;
  password: string;
  remember?: boolean;
  captcha?: string;
  captchaId?: string;
}

/** 登录响应数据 */
export interface LoginResult {
  token: string;
  refreshToken: string;
  expiresIn: number;
  user: UserInfo;
}

/** Token信息 */
export interface TokenInfo {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

/** 用户基本信息 */
export interface UserInfo {
  id: number;
  username: string;
  realName: string;
  avatar: string;
  email: string;
  phone: string;
  department: string;
  departmentId: number;
  roles: string[];
  permissions: string[];
  status: 'active' | 'disabled';
  lastLoginTime: string;
}
