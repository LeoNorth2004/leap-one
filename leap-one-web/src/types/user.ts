/** 用户相关类型定义 */

import type { PaginationParams, SortParams } from './api';

/** 用户列表查询参数 */
export interface UserListParams extends PaginationParams, SortParams {
  keyword?: string;
  departmentId?: number;
  status?: 'active' | 'disabled';
  roleId?: number;
}

/** 用户信息（扩展版，含部门角色详情） */
export interface UserDetail {
  id: number;
  username: string;
  realName: string;
  avatar: string;
  email: string;
  phone: string;
  gender: 'male' | 'female' | 'unknown';
  departmentId: number;
  departmentName: string;
  position: string;
  roles: UserRole[];
  status: 'active' | 'disabled';
  createdAt: string;
  updatedAt: string;
  lastLoginTime: string;
}

/** 用户角色关联 */
export interface UserRole {
  id: number;
  name: string;
  code: string;
}

/** 创建用户参数 */
export interface CreateUserParams {
  username: string;
  password: string;
  realName: string;
  email: string;
  phone?: string;
  departmentId: number;
  roleIds: number[];
}

/** 更新用户参数 */
export interface UpdateUserParams extends Partial<CreateUserParams> {
  id: number;
  status?: 'active' | 'disabled';
}
