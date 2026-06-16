/** 文档服务API */

import { apiClient } from './client';

const BASE_URL = '/document';

/** 获取文档列表/目录树 */
export function getDocumentListApi(params?: Record<string, unknown>) {
  return apiClient.getPage(`${BASE_URL}/list`, params);
}

/** 获取文档详情 */
export function getDocumentDetailApi(id: number) {
  return apiClient.get(`${BASE_URL}/${id}`);
}

/** 创建文档 */
export function createDocumentApi(data: Record<string, unknown>) {
  return apiClient.post(BASE_URL, data);
}

/** 更新文档 */
export function updateDocumentApi(id: number, data: Record<string, unknown>) {
  return apiClient.put(`${BASE_URL}/${id}`, data);
}

/** 删除文档 */
export function deleteDocumentApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 上传文档附件 */
export function uploadDocumentApi(file: FormData) {
  return apiClient.post(`${BASE_URL}/upload`, file);
}
