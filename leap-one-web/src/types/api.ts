/** API通用响应类型定义 */

/** 通用API响应结构 */
export interface ApiResponse<T = unknown> {
  /** 业务状态码，0表示成功 */
  code: number;
  /** 响应消息 */
  message: string;
  /** 响应数据 */
  data: T;
}

/** 分页响应数据结构 */
export interface PaginatedData<T = unknown> {
  /** 数据列表 */
  list: T[];
  /** 总记录数 */
  total: number;
  /** 当前页码（从1开始） */
  page: number;
  /** 每页条数 */
  pageSize: number;
}

/** 分页API响应结构 */
export type PaginatedResponse<T = unknown> = ApiResponse<PaginatedData<T>>;

/** 分页查询参数 */
export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

/** 排序参数 */
export interface SortParams {
  sortField?: string;
  sortOrder?: 'ascend' | 'descend';
}
