/**
 * 通知服务 API
 *
 * 提供通知列表查询、已读标记、删除等接口
 */

import { apiClient } from './client';

const BASE = '/notification';

// ── 查询接口 ─────────────────────────────────────────────────

/** 获取通知列表（分页） */
export const getNotificationListApi = (params?: Record<string, unknown>) =>
  apiClient.getPage(`${BASE}/list`, params);

/** 获取未读通知数量 */
export const getUnreadCountApi = (): Promise<number> =>
  apiClient
    .get<{ count: number }>(`${BASE}/unread/count`)
    .then((res) => res.data.count);

// ── 操作接口 ─────────────────────────────────────────────────

/** 标记指定通知为已读 */
export const markAsReadApi = (ids: number[]): Promise<void> =>
  apiClient.put(`${BASE}/read`, { ids }).then(() => undefined);

/** 全部标记为已读 */
export const markAllReadApi = (): Promise<void> =>
  apiClient.put(`${BASE}/read/all`).then(() => undefined);

/** 删除通知 */
export const deleteNotificationApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);
