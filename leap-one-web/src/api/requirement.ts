/** 需求服务API */

import { apiClient } from './client';
import type { Requirement, RequirementListParams } from '@/types/requirement';

const BASE_URL = '/requirement';

/** 获取需求列表 */
export function getRequirementListApi(params?: RequirementListParams) {
  return apiClient.getPage<Requirement>(`${BASE_URL}/list`, params as Record<string, unknown>);
}

/** 获取需求详情 */
export function getRequirementDetailApi(id: number): Promise<Requirement> {
  return apiClient.get<Requirement>(`${BASE_URL}/${id}`).then((res) => res.data);
}

/** 创建需求 */
export function createRequirementApi(data: Partial<Requirement>): Promise<Requirement> {
  return apiClient.post<Requirement>(BASE_URL, data).then((res) => res.data);
}

/** 更新需求 */
export function updateRequirementApi(id: number, data: Partial<Requirement>): Promise<Requirement> {
  return apiClient.put<Requirement>(`${BASE_URL}/${id}`, data).then((res) => res.data);
}

/** 删除需求 */
export function deleteRequirementApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 变更需求状态 */
export function changeRequirementStatusApi(
  id: number,
  status: Requirement['status']
): Promise<Requirement> {
  return apiClient.put<Requirement>(`${BASE_URL}/${id}/status`, { status }).then((res) => res.data);
}

/** 需求评审 */
export function reviewRequirementApi(
  id: number,
  action: 'approve' | 'reject',
  comment?: string
): Promise<Requirement> {
  return apiClient
    .put<Requirement>(`${BASE_URL}/${id}/review`, { action, comment })
    .then((res) => res.data);
}
