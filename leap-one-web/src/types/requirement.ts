/** 需求相关类型定义 */

import type { PaginationParams, SortParams } from './api';

/** 需求状态 */
export type RequirementStatus =
  | 'draft'
  | 'reviewing'
  | 'active'
  | 'developing'
  | 'testing'
  | 'completed'
  | 'closed'
  | 'rejected';

/** 需求优先级 */
export type Priority = 'P0' | 'P1' | 'P2' | 'P3';

/** 需求来源 */
export type RequirementSource = 'customer' | 'market' | 'internal' | 'competitive';

/** 需求列表查询参数 */
export interface RequirementListParams extends PaginationParams, SortParams {
  keyword?: string;
  productId?: number;
  projectId?: number;
  status?: RequirementStatus;
  priority?: Priority;
  assigneeId?: number;
  source?: RequirementSource;
}

/** 需求信息 */
export interface Requirement {
  id: number;
  title: string;
  code: string;
  description: string;
  status: RequirementStatus;
  priority: Priority;
  source: RequirementSource;
  productId: number;
  productName: string;
  projectId?: number;
  projectName?: string;
  moduleId?: number;
  moduleName?: string;
  storyPoints?: number;
  assigneeId?: number;
  assigneeName?: string;
  reviewerId?: number;
  reviewerName?: string;
  planRelease?: string;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
}
