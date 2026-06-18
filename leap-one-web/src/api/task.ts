/**
 * 任务服务 API
 *
 * 提供任务 CRUD、状态变更、指派、批量操作等接口
 */

import { apiClient } from './client';
import type { Task, TaskListParams } from '@/types/task';

const BASE = '/task';

// ── 查询接口 ─────────────────────────────────────────────────

/** 获取任务列表（分页） */
export const getTaskListApi = (params?: TaskListParams) =>
  apiClient.getPage<Task>(`${BASE}/list`, params as Record<string, unknown>);

/** 获取任务详情 */
export const getTaskDetailApi = (id: number): Promise<Task> =>
  apiClient.get<Task>(`${BASE}/${id}`).then((res) => res.data);

// ── 写入接口 ─────────────────────────────────────────────────

/** 创建任务 */
export const createTaskApi = (data: Partial<Task>): Promise<Task> =>
  apiClient.post<Task>(BASE, data).then((res) => res.data);

/** 更新任务 */
export const updateTaskApi = (id: number, data: Partial<Task>): Promise<Task> =>
  apiClient.put<Task>(`${BASE}/${id}`, data).then((res) => res.data);

/** 删除任务 */
export const deleteTaskApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);

// ── 状态 & 指派接口 ──────────────────────────────────────────

/** 变更任务状态 */
export const changeTaskStatusApi = (id: number, status: Task['status']): Promise<Task> =>
  apiClient.put<Task>(`${BASE}/${id}/status`, { status }).then((res) => res.data);

/** 任务指派 */
export const assignTaskApi = (id: number, assigneeId: number): Promise<Task> =>
  apiClient.put<Task>(`${BASE}/${id}/assign`, { assigneeId }).then((res) => res.data);

/** 批量更新任务状态 */
export const batchUpdateTaskStatusApi = (ids: number[], status: Task['status']): Promise<void> =>
  apiClient.put(`${BASE}/batch/status`, { ids, status }).then(() => undefined);
