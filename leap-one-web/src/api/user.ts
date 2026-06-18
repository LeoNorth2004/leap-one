/**
 * 用户服务 API
 *
 * 提供用户 CRUD、状态管理、密码修改等接口
 */

import { apiClient } from './client';
import type { UserListParams, UserDetail, CreateUserParams, UpdateUserParams } from '@/types/user';

const BASE = '/user';

// ── 查询接口 ─────────────────────────────────────────────────

/** 获取用户列表（分页） */
export const getUserListApi = (params?: UserListParams) =>
  apiClient.getPage<{ list: UserDetail[]; total: number }>(
    `${BASE}/list`,
    params as Record<string, unknown>
  );

/** 获取用户详情 */
export const getUserDetailApi = (id: number): Promise<UserDetail> =>
  apiClient.get<UserDetail>(`${BASE}/${id}`).then((res) => res.data);

// ── 写入接口 ─────────────────────────────────────────────────

/** 创建用户 */
export const createUserApi = (data: CreateUserParams): Promise<UserDetail> =>
  apiClient.post<UserDetail>(BASE, data).then((res) => res.data);

/** 更新用户信息 */
export const updateUserApi = (data: UpdateUserParams): Promise<UserDetail> =>
  apiClient.put<UserDetail>(`${BASE}/${data.id}`, data).then((res) => res.data);

/** 删除用户 */
export const deleteUserApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);

// ── 状态 & 安全接口 ──────────────────────────────────────────

/** 更新用户状态（启用/禁用） */
export const updateUserStatusApi = (
  id: number,
  status: 'active' | 'disabled'
): Promise<void> =>
  apiClient.put(`${BASE}/${id}/status`, { status }).then(() => undefined);

/** 修改密码 */
export const changePasswordApi = (
  oldPassword: string,
  newPassword: string
): Promise<void> =>
  apiClient.put(`${BASE}/password`, { oldPassword, newPassword }).then(() => undefined);
