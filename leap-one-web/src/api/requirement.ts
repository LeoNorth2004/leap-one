/**
 * 需求服务 API
 *
 * 提供需求 CRUD、状态变更、评审等接口
 */

import { apiClient } from './client';
import type { Requirement, RequirementListParams } from '@/types/requirement';

const BASE = '/requirement';

// ── 查询接口 ─────────────────────────────────────────────────

/** 获取需求列表（分页） */
export const getRequirementListApi = (params?: RequirementListParams) =>
  apiClient.getPage<Requirement>(`${BASE}/list`, params as Record<string, unknown>);

/** 获取需求详情 */
export const getRequirementDetailApi = (id: number): Promise<Requirement> =>
  apiClient.get<Requirement>(`${BASE}/${id}`).then((res) => res.data);

// ── 写入接口 ─────────────────────────────────────────────────

/** 创建需求 */
export const createRequirementApi = (data: Partial<Requirement>): Promise<Requirement> =>
  apiClient.post<Requirement>(BASE, data).then((res) => res.data);

/** 更新需求 */
export const updateRequirementApi = (id: number, data: Partial<Requirement>): Promise<Requirement> =>
  apiClient.put<Requirement>(`${BASE}/${id}`, data).then((res) => res.data);

/** 删除需求 */
export const deleteRequirementApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);

// ── 状态 & 评审接口 ──────────────────────────────────────────

/** 变更需求状态 */
export const changeRequirementStatusApi = (
  id: number,
  status: Requirement['status']
): Promise<Requirement> =>
  apiClient
    .put<Requirement>(`${BASE}/${id}/status`, { status })
    .then((res) => res.data);

/** 需求评审（通过 / 拒绝） */
export const reviewRequirementApi = (
  id: number,
  action: 'approve' | 'reject',
  comment?: string
): Promise<Requirement> =>
  apiClient
    .put<Requirement>(`${BASE}/${id}/review`, { action, comment })
    .then((res) => res.data);
