/**
 * Axios 实例配置 - API 客户端
 *
 * 功能：
 * - baseURL 从环境变量 VITE_API_BASE_URL 读取
 * - 请求拦截器：自动注入 Authorization Bearer Token
 * - 响应拦截器：统一错误处理（401/403/其他）
 * - 支持请求取消（AbortController / CancelToken）
 * - 重复请求自动取消
 * - 导出 client 实例及独立的 get/post/put/del 方法
 */

import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
  type AxiosResponse,
  type Canceler,
} from 'axios';
import { message } from 'antd';
import { API_CODE } from '@/utils/constants';
import { tokenStorage } from '@/utils/storage';
import type { ApiResponse } from '@/types/api';

// ─── 重复请求管理 ───────────────────────────────────────────────

const pendingRequests = new Map<string, Canceler>();

/** 根据请求配置生成唯一标识 */
function generateReqKey(config: InternalAxiosRequestConfig): string {
  const { method, url, params, data } = config;
  return [method ?? 'get', url, JSON.stringify(params), JSON.stringify(data)].join('&');
}

/** 添加待取消请求到 Map */
function addPendingRequest(config: InternalAxiosRequestConfig): void {
  const key = generateReqKey(config);
  config.cancelToken = new axios.CancelToken((cancel) => {
    if (!pendingRequests.has(key)) {
      pendingRequests.set(key, cancel);
    }
  });
}

/** 移除并取消已存在的重复请求 */
function removePendingRequest(config: InternalAxiosRequestConfig): void {
  const key = generateReqKey(config);
  if (pendingRequests.has(key)) {
    const cancel = pendingRequests.get(key)!;
    cancel(`重复请求被取消: ${config.url}`);
    pendingRequests.delete(key);
  }
}

// ─── 创建 Axios 实例 ─────────────────────────────────────────────

const client: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api',
  timeout: 30_000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
});

// ─── 请求拦截器 ──────────────────────────────────────────────────

client.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 取消重复请求
    removePendingRequest(config);
    addPendingRequest(config);

    // 从 localStorage 读取 token 并注入 Authorization header
    const token = tokenStorage.getToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// ─── 响应拦截器 ──────────────────────────────────────────────────

client.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // 请求完成后从 pendingMap 中移除
    removePendingRequest(response.config);

    const resData = response.data;

    // 业务层错误码处理（后端返回的 code !== 0）
    if (resData.code !== undefined && resData.code !== API_CODE.SUCCESS) {
      switch (resData.code) {
        case API_CODE.UNAUTHORIZED:
          tokenStorage.clearAll();
          window.location.href = '/login';
          message.error('登录已过期，请重新登录');
          break;
        case API_CODE.FORBIDDEN:
          message.error('抱歉，您没有权限执行此操作');
          break;
        default:
          message.warning(resData.message ?? '请求失败');
      }
      return Promise.reject(new Error(resData.message ?? '请求失败'));
    }

    // 成功响应：直接返回完整 response（调用方通过 .data 获取数据）
    return response;
  },

  (error) => {
    // 移除失败的请求记录
    if (error.config) {
      removePendingRequest(error.config as InternalAxiosRequestConfig);
    }

    // 请求被取消时不显示错误提示
    if (axios.isCancel(error)) {
      console.warn('[API] 请求被取消:', error.message);
      return Promise.reject(error);
    }

    // HTTP 状态码错误处理
    const status = error.response?.status;
    switch (status) {
      case 401:
        tokenStorage.clearAll();
        window.location.href = '/login';
        message.error('登录已过期，请重新登录');
        break;
      case 403:
        message.error('抱歉，您没有权限访问此资源');
        break;
      case 404:
        message.error('请求的资源不存在');
        break;
      case 500:
        message.error('服务器内部错误，请稍后重试');
        break;
      case 502:
      case 503:
      case 504:
        message.error('服务暂时不可用，请稍后重试');
        break;
      default:
        message.error(error.message ?? '网络异常，请检查网络连接');
    }

    return Promise.reject(error);
  }
);

// ─── 便捷方法封装 ───────────────────────────────────────────────

/** 带 AbortController 的请求配置扩展 */
interface RequestConfig extends AxiosRequestConfig {
  /** 是否启用AbortController，默认 true */
  abortable?: boolean;
}

/**
 * GET 请求
 */
export function get<T = unknown>(
  url: string,
  params?: Record<string, unknown>,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  const controller = config?.abortable !== false ? new AbortController() : undefined;
  return client
    .get<ApiResponse<T>>(url, { params, ...config, signal: controller?.signal })
    .then((res) => res.data);
}

/**
 * POST 请求
 */
export function post<T = unknown>(
  url: string,
  data?: unknown,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  const controller = config?.abortable !== false ? new AbortController() : undefined;
  return client
    .post<ApiResponse<T>>(url, data, { ...config, signal: controller?.signal })
    .then((res) => res.data);
}

/**
 * PUT 请求
 */
export function put<T = unknown>(
  url: string,
  data?: unknown,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  const controller = config?.abortable !== false ? new AbortController() : undefined;
  return client
    .put<ApiResponse<T>>(url, data, { ...config, signal: controller?.signal })
    .then((res) => res.data);
}

/**
 * DELETE 请求
 */
export function del<T = unknown>(
  url: string,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  const controller = config?.abortable !== false ? new AbortController() : undefined;
  return client
    .delete<ApiResponse<T>>(url, { ...config, signal: controller?.signal })
    .then((res) => res.data);
}

// ─── 导出 ─────────────────────────────────────────────────────────

export default client;
