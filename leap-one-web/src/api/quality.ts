/**
 * 质量服务 API（测试用例、Bug、测试计划）
 *
 * 以对象形式组织三类资源的 CRUD 接口
 */

import { apiClient } from './client';

const BASE = '/quality';

// ── 测试用例 API ─────────────────────────────────────────────

export const testCaseApi = Object.freeze({
  /** 获取用例列表（分页） */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE}/testcase/list`, params),

  /** 获取用例详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE}/testcase/${id}`),

  /** 创建用例 */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE}/testcase`, data),

  /** 更新用例 */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE}/testcase/${id}`, data),

  /** 删除用例 */
  delete: (id: number) =>
    apiClient.delete(`${BASE}/testcase/${id}`),
});

// ── Bug API ───────────────────────────────────────────────────

export const bugApi = Object.freeze({
  /** 获取 Bug 列表（分页） */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE}/bug/list`, params),

  /** 获取 Bug 详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE}/bug/${id}`),

  /** 提交 Bug */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE}/bug`, data),

  /** 更新 Bug */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE}/bug/${id}`, data),

  /** 删除 Bug */
  delete: (id: number) =>
    apiClient.delete(`${BASE}/bug/${id}`),

  /** 变更 Bug 状态（确认/修复/关闭/激活） */
  changeStatus: (id: number, status: string) =>
    apiClient.put(`${BASE}/bug/${id}/status`, { status }),
});

// ── 测试计划 API ─────────────────────────────────────────────

export const testPlanApi = Object.freeze({
  /** 获取测试计划列表（分页） */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE}/testplan/list`, params),

  /** 获取测试计划详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE}/testplan/${id}`),

  /** 创建测试计划 */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE}/testplan`, data),

  /** 更新测试计划 */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE}/testplan/${id}`, data),

  /** 删除测试计划 */
  delete: (id: number) =>
    apiClient.delete(`${BASE}/testplan/${id}`),
});
