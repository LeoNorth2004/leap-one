/** 任务相关类型定义 */

import type { PaginationParams, SortParams } from './api';

/** 任务状态 */
export type TaskStatus =
  | 'wait'
  | 'doing'
  | 'done'
  | 'pause'
  | 'cancel'
  | 'closed';

/** 任务类型 */
export type TaskType = 'dev' | 'test' | 'design' | 'review' | 'meeting' | 'other';

/** 任务优先级 */
export type TaskPriority = 'urgent' | 'high' | 'medium' | 'low';

/** 任务列表查询参数 */
export interface TaskListParams extends PaginationParams, SortParams {
  keyword?: string;
  projectId?: number;
  iterationId?: number;
  requirementId?: number;
  status?: TaskStatus;
  type?: TaskType;
  priority?: TaskPriority;
  assigneeId?: number;
  createdBy?: number;
}

/** 任务信息 */
export interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  type: TaskType;
  priority: TaskPriority;
  projectId: number;
  projectName: string;
  iterationId?: number;
  iterationName?: string;
  requirementId?: number;
  requirementTitle?: string;
  parentId?: number;
  assigneeId?: number;
  assigneeName?: string;
  createdBy: number;
  createdByName: string;
  estimatedHours?: number;
  consumedHours?: number;
  leftHours?: number;
  startDate?: string;
  dueDate?: string;
  finishedDate?: string;
  createdAt: string;
  updatedAt: string;
}
