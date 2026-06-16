/** 项目相关类型定义 */

import type { PaginationParams, SortParams } from './api';

/** 项目状态 */
export type ProjectStatus = 'planning' | 'active' | 'paused' | 'completed' | 'archived';

/** 项目类型 */
export type ProjectType = 'scrum' | 'waterfall' | 'kanban' | 'hybrid';

/** 项目列表查询参数 */
export interface ProjectListParams extends PaginationParams, SortParams {
  keyword?: string;
  status?: ProjectStatus;
  type?: ProjectType;
  pmId?: number;
  productId?: number;
  programId?: number;
}

/** 项目信息 */
export interface Project {
  id: number;
  name: string;
  code: string;
  description: string;
  status: ProjectStatus;
  type: ProjectType;
  pmId: number;
  pmName: string;
  productId: number;
  productName: string;
  programId?: number;
  programName?: string;
  startDate: string;
  endDate: string;
  progress: number;
  memberCount: number;
  avatar: string;
  createdAt: string;
  updatedAt: string;
}

/** 项目成员 */
export interface ProjectMember {
  userId: number;
  userName: string;
  avatar: string;
  role: 'pm' | 'developer' | 'tester' | 'observer';
  joinedAt: string;
}

/** 迭代/Sprint */
export interface Iteration {
  id: number;
  projectId: number;
  name: string;
  status: 'pending' | 'active' | 'completed';
  startDate: string;
  endDate: string;
  goal: string;
  taskCount: number;
  completedTaskCount: number;
  progress: number;
}
