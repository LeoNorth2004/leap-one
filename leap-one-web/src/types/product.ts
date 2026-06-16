/** 产品相关类型定义 */

import type { PaginationParams, SortParams } from './api';

/** 产品状态 */
export type ProductStatus = 'normal' | 'closed';

/** 产品列表查询参数 */
export interface ProductListParams extends PaginationParams, SortParams {
  keyword?: string;
  status?: ProductStatus;
  managerId?: number;
}

/** 产品信息 */
export interface Product {
  id: number;
  name: string;
  code: string;
  description: string;
  status: ProductStatus;
  managerId: number;
  managerName: string;
  projectId?: number;
  projectName?: string;
  line?: string;
  createdAt: string;
  updatedAt: string;
}

/** 产品路线图版本 */
export interface RoadmapVersion {
  id: number;
  productId: number;
  name: string;
  version: string;
  planDate: string;
  releaseDate?: string;
  status: 'planning' | 'developing' | 'released' | 'delayed';
  requirements: RoadmapRequirement[];
}

/** 路线图需求项 */
export interface RoadmapRequirement {
  id: number;
  title: string;
  priority: 'P0' | 'P1' | 'P2' | 'P3';
  status: string;
  assignee?: string;
}
