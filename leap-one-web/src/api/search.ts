/** 搜索服务API */

import { apiClient } from './client';

const BASE_URL = '/search';

/** 全局搜索 */
export function globalSearchApi(keyword: string, params?: Record<string, unknown>) {
  return apiClient.get(`${BASE_URL}/global`, { keyword, ...params });
}

/** 搜索用户 */
export function searchUserApi(keyword: string) {
  return apiClient.get(`${BASE_URL}/user`, { keyword });
}

/** 搜索项目 */
export function searchProjectApi(keyword: string) {
  return apiClient.get(`${BASE_URL}/project`, { keyword });
}

/** 搜索需求 */
export function searchRequirementApi(keyword: string) {
  return apiClient.get(`${BASE_URL}/requirement`, { keyword });
}

/** 搜索任务 */
export function searchTaskApi(keyword: string) {
  return apiClient.get(`${BASE_URL}/task`, { keyword });
}

/** 搜索Bug */
export function searchBugApi(keyword: string) {
  return apiClient.get(`${BASE_URL}/bug`, { keyword });
}
