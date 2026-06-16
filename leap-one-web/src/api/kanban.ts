/** 看板服务API */

import { apiClient } from './client';

const BASE_URL = '/kanban';

/** 获取看板数据（含列和卡片） */
export function getKanbanDataApi(projectId?: number) {
  const params = projectId ? { projectId } : {};
  return apiClient.get(`${BASE_URL}/data`, params);
}

/** 移动看板卡片（拖拽排序） */
export function moveKanbanCardApi(
  cardId: number,
  targetColumnId: string,
  targetIndex: number
) {
  return apiClient.put(`${BASE_URL}/card/move`, {
    cardId,
    targetColumnId,
    targetIndex,
  });
}

/** 更新看板卡片 */
export function updateKanbanCardApi(cardId: number, data: Record<string, unknown>) {
  return apiClient.put(`${BASE_URL}/card/${cardId}`, data);
}
