/** 产品服务API */

import { apiClient } from './client';
import type { Product, ProductListParams, RoadmapVersion } from '@/types/product';

const BASE_URL = '/product';

/** 获取产品列表 */
export function getProductListApi(params?: ProductListParams) {
  return apiClient.getPage<Product>(`${BASE_URL}/list`, params as Record<string, unknown>);
}

/** 获取产品详情 */
export function getProductDetailApi(id: number): Promise<Product> {
  return apiClient.get<Product>(`${BASE_URL}/${id}`).then((res) => res.data);
}

/** 创建产品 */
export function createProductApi(data: Partial<Product>): Promise<Product> {
  return apiClient.post<Product>(BASE_URL, data).then((res) => res.data);
}

/** 更新产品 */
export function updateProductApi(id: number, data: Partial<Product>): Promise<Product> {
  return apiClient.put<Product>(`${BASE_URL}/${id}`, data).then((res) => res.data);
}

/** 删除产品 */
export function deleteProductApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 获取产品路线图 */
export function getRoadmapApi(productId: number): Promise<RoadmapVersion[]> {
  return apiClient
    .get<RoadmapVersion[]>(`${BASE_URL}/${productId}/roadmap`)
    .then((res) => res.data);
}
