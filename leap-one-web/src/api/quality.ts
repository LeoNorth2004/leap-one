/** 质量服务API（测试用例、Bug、测试计划） */

import { apiClient } from './client';

const BASE_URL = '/quality';

/** 测试用例相关 */
export const testCaseApi = {
  /** 获取用例列表 */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE_URL}/testcase/list`, params),

  /** 获取用例详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE_URL}/testcase/${id}`),

  /** 创建用例 */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE_URL}/testcase`, data),

  /** 更新用例 */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE_URL}/testcase/${id}`, data),

  /** 删除用例 */
  delete: (id: number) =>
    apiClient.delete(`${BASE_URL}/testcase/${id}`),
};

/** Bug相关 */
export const bugApi = {
  /** 获取Bug列表 */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE_URL}/bug/list`, params),

  /** 获取Bug详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE_URL}/bug/${id}`),

  /** 提交Bug */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE_URL}/bug`, data),

  /** 更新Bug */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE_URL}/bug/${id}`, data),

  /** 删除Bug */
  delete: (id: number) =>
    apiClient.delete(`${BASE_URL}/bug/${id}`),

  /** 变更Bug状态（确认/修复/关闭/激活） */
  changeStatus: (id: number, status: string) =>
    apiClient.put(`${BASE_URL}/bug/${id}/status`, { status }),
};

/** 测试计划相关 */
export const testPlanApi = {
  /** 获取测试计划列表 */
  getList: (params?: Record<string, unknown>) =>
    apiClient.getPage(`${BASE_URL}/testplan/list`, params),

  /** 获取测试计划详情 */
  getDetail: (id: number) =>
    apiClient.get(`${BASE_URL}/testplan/${id}`),

  /** 创建测试计划 */
  create: (data: Record<string, unknown>) =>
    apiClient.post(`${BASE_URL}/testplan`, data),

  /** 更新测试计划 */
  update: (id: number, data: Record<string, unknown>) =>
    apiClient.put(`${BASE_URL}/testplan/${id}`, data),

  /** 删除测试计划 */
  delete: (id: number) =>
    apiClient.delete(`${BASE_URL}/testplan/${id}`),
};
