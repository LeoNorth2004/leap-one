/** 任务服务API */

import { apiClient } from './client';
import type { Task, TaskListParams } from '@/types/task';

const BASE_URL = '/task';

/** 获取任务列表 */
export function getTaskListApi(params?: TaskListParams) {
  return apiClient.getPage<Task>(`${BASE_URL}/list`, params as Record<string, unknown>);
}

/** 获取任务详情 */
export function getTaskDetailApi(id: number): Promise<Task> {
  return apiClient.get<Task>(`${BASE_URL}/${id}`).then((res) => res.data);
}

/** 创建任务 */
export function createTaskApi(data: Partial<Task>): Promise<Task> {
  return apiClient.post<Task>(BASE_URL, data).then((res) => res.data);
}

/** 更新任务 */
export function updateTaskApi(id: number, data: Partial<Task>): Promise<Task> {
  return apiClient.put<Task>(`${BASE_URL}/${id}`, data).then((res) => res.data);
}

/** 删除任务 */
export function deleteTaskApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 变更任务状态 */
export function changeTaskStatusApi(id: number, status: Task['status']): Promise<Task> {
  return apiClient.put<Task>(`${BASE_URL}/${id}/status`, { status }).then((res) => res.data);
}

/** 任务指派 */
export function assignTaskApi(id: number, assigneeId: number): Promise<Task> {
  return apiClient.put<Task>(`${BASE_URL}/${id}/assign`, { assigneeId }).then((res) => res.data);
}

/** 批量更新任务状态 */
export function batchUpdateTaskStatusApi(ids: number[], status: Task['status']): Promise<void> {
  return apiClient.put(`${BASE_URL}/batch/status`, { ids, status }).then(() => undefined);
}
