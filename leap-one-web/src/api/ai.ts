/** AI智能助手服务API */

import { apiClient } from './client';

const BASE_URL = '/ai';

/** AI对话 - 发送消息 */
export function chatAiApi(message: string, context?: string) {
  return apiClient.post(`${BASE_URL}/chat`, { message, context });
}

/** AI智能拆解需求为任务 */
export function aiSplitRequirementApi(requirementId: number) {
  return apiClient.post(`${BASE_URL}/split-requirement`, { requirementId });
}

/** AI生成测试用例 */
export function aiGenerateTestCaseApi(requirementId: number) {
  return apiClient.post(`${BASE_URL}/generate-testcase`, { requirementId });
}

/** AI分析项目风险 */
export function aiAnalyzeRiskApi(projectId: number) {
  return apiClient.post(`${BASE_URL}/analyze-risk`, { projectId });
}

/** AI生成周报/日报摘要 */
export function aiGenerateReportApi(type: 'weekly' | 'daily', userId?: number) {
  return apiClient.post(`${BASE_URL}/generate-report`, { type, userId });
}

/** AI智能推荐任务指派 */
export function aiSuggestAssigneeApi(taskId: number) {
  return apiClient.post(`${BASE_URL}/suggest-assignee`, { taskId });
}
