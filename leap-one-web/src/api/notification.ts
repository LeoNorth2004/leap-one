/** 通知服务API */

import { apiClient } from './client';

const BASE_URL = '/notification';

/** 获取通知列表 */
export function getNotificationListApi(params?: Record<string, unknown>) {
  return apiClient.getPage(`${BASE_URL}/list`, params);
}

/** 获取未读通知数量 */
export function getUnreadCountApi(): Promise<number> {
  return apiClient
    .get<{ count: number }>(`${BASE_URL}/unread/count`)
    .then((res) => res.data.count);
}

/** 标记通知为已读 */
export function markAsReadApi(ids: number[]): Promise<void> {
  return apiClient.put(`${BASE_URL}/read`, { ids }).then(() => undefined);
}

/** 全部标记为已读 */
export function markAllReadApi(): Promise<void> {
  return apiClient.put(`${BASE_URL}/read/all`).then(() => undefined);
}

/** 删除通知 */
export function deleteNotificationApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}
