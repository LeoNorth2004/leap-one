/**
 * 文档服务 API
 *
 * 提供文档 CRUD、附件上传等接口
 */

import { apiClient } from './client';

const BASE = '/document';

// ── 查询接口 ─────────────────────────────────────────────────

/** 获取文档列表 / 目录树（分页） */
export const getDocumentListApi = (params?: Record<string, unknown>) =>
  apiClient.getPage(`${BASE}/list`, params);

/** 获取文档详情 */
export const getDocumentDetailApi = (id: number) =>
  apiClient.get(`${BASE}/${id}`);

// ── 写入接口 ─────────────────────────────────────────────────

/** 创建文档 */
export const createDocumentApi = (data: Record<string, unknown>) =>
  apiClient.post(BASE, data);

/** 更新文档 */
export const updateDocumentApi = (id: number, data: Record<string, unknown>) =>
  apiClient.put(`${BASE}/${id}`, data);

/** 删除文档 */
export const deleteDocumentApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);

// ── 文件操作 ─────────────────────────────────────────────────

/** 上传文档附件 */
export const uploadDocumentApi = (file: FormData) =>
  apiClient.post(`${BASE}/upload`, file);
