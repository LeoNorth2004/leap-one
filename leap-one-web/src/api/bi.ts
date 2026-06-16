/** BI统计服务API */

import { apiClient } from './client';

const BASE_URL = '/bi';

/** 获取工作台概览统计数据 */
export function getDashboardStatsApi() {
  return apiClient.get(`${BASE_URL}/dashboard/stats`);
}

/** 获取项目进度趋势数据 */
export function getProjectTrendApi(projectId: number, days = 30) {
  return apiClient.get(`${BASE_URL}/project/trend`, { projectId, days });
}

/** 获取需求分布数据 */
export function getRequirementDistributionApi(params?: Record<string, unknown>) {
  return apiClient.get(`${BASE_URL}/requirement/distribution`, params);
}

/** 获取质量指标数据 */
export function getQualityMetricsApi(projectId?: number) {
  const params = projectId ? { projectId } : {};
  return apiClient.get(`${BASE_URL}/quality/metrics`, params);
}

/** 获取团队效能数据 */
export function getTeamEfficiencyApi(teamId?: number) {
  const params = teamId ? { teamId } : {};
  return apiClient.get(`${BASE_URL}/team/efficiency`, params);
}

/** 获取工时统计数据 */
export function getWorkHourStatsApi(params?: Record<string, unknown>) {
  return apiClient.get(`${BASE_URL}/workhour/stats`, params);
}
