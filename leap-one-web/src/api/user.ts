/** 用户服务API */

import { apiClient } from './client';
import type { UserListParams, UserDetail, CreateUserParams, UpdateUserParams } from '@/types/user';

const BASE_URL = '/user';

/** 获取用户列表（分页） */
export function getUserListApi(params?: UserListParams) {
  return apiClient.getPage<UserDetail>(`${BASE_URL}/list`, params as Record<string, unknown>);
}

/** 获取用户详情 */
export function getUserDetailApi(id: number): Promise<UserDetail> {
  return apiClient.get<UserDetail>(`${BASE_URL}/${id}`).then((res) => res.data);
}

/** 创建用户 */
export function createUserApi(data: CreateUserParams): Promise<UserDetail> {
  return apiClient.post<UserDetail>(BASE_URL, data).then((res) => res.data);
}

/** 更新用户信息 */
export function updateUserApi(data: UpdateUserParams): Promise<UserDetail> {
  return apiClient.put<UserDetail>(`${BASE_URL}/${data.id}`, data).then((res) => res.data);
}

/** 删除用户 */
export function deleteUserApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 更新用户状态（启用/禁用） */
export function updateUserStatusApi(id: number, status: 'active' | 'disabled'): Promise<void> {
  return apiClient.put(`${BASE_URL}/${id}/status`, { status }).then(() => undefined);
}

/** 修改密码 */
export function changePasswordApi(oldPassword: string, newPassword: string): Promise<void> {
  return apiClient.put(`${BASE_URL}/password`, { oldPassword, newPassword }).then(() => undefined);
}
