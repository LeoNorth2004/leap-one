/** 工单服务API */

import { apiClient } from './client';

const BASE_URL = '/issue';

/** 获取工单列表 */
export function getIssueListApi(params?: Record<string, unknown>) {
  return apiClient.getPage(`${BASE_URL}/list`, params);
}

/** 获取工单详情 */
export function getIssueDetailApi(id: number) {
  return apiClient.get(`${BASE_URL}/${id}`);
}

/** 创建工单 */
export function createIssueApi(data: Record<string, unknown>) {
  return apiClient.post(BASE_URL, data);
}

/** 更新工单 */
export function updateIssueApi(id: number, data: Record<string, unknown>) {
  return apiClient.put(`${BASE_URL}/${id}`, data);
}

/** 删除工单 */
export function deleteIssueApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 变更工单状态 */
export function changeIssueStatusApi(id: number, status: string) {
  return apiClient.put(`${BASE_URL}/${id}/status`, { status });
}
