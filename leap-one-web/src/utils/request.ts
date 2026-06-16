/** Axios请求封装 */

import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
  type AxiosResponse,
  type Canceler,
} from 'axios';
import { message } from 'antd';
import { API_CODE } from './constants';
import { tokenStorage } from './storage';
import type { ApiResponse } from '@/types/api';

/** 存储待取消的请求Map */
const pendingRequests = new Map<string, Canceler>();

/**
 * 生成请求唯一标识
 */
function generateReqKey(config: InternalAxiosRequestConfig): string {
  const { method, url, params, data } = config;
  return [method, url, JSON.stringify(params), JSON.stringify(data)].join('&');
}

/**
 * 添加待取消请求
 */
function addPendingRequest(config: InternalAxiosRequestConfig): void {
  const key = generateReqKey(config);
  config.cancelToken = new axios.CancelToken((cancel) => {
    if (!pendingRequests.has(key)) {
      pendingRequests.set(key, cancel);
    }
  });
}

/**
 * 移除并取消已存在的重复请求
 */
function removePendingRequest(config: InternalAxiosRequestConfig): void {
  const key = generateReqKey(config);
  if (pendingRequests.has(key)) {
    const cancel = pendingRequests.get(key)!;
    cancel(`重复请求被取消: ${config.url}`);
    pendingRequests.delete(key);
  }
}

/** 创建Axios实例 */
const client: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
});

/** 请求拦截器：注入Token、处理重复请求 */
client.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 取消重复请求
    removePendingRequest(config);
    addPendingRequest(config);

    // 注入Authorization Bearer Token
    const token = tokenStorage.getToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

/** 响应拦截器：统一错误处理、401跳转登录、403提示无权限 */
client.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // 请求完成后移除该请求
    removePendingRequest(response.config);

    const resData = response.data;

    // 业务层错误处理
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
          message.warning(resData.message || '请求失败');
      }
      return Promise.reject(new Error(resData.message || '请求失败'));
    }

    return response;
  },
  (error) => {
    // 移除失败的请求
    if (error.config) {
      removePendingRequest(error.config);
    }

    if (axios.isCancel(error)) {
      console.warn('请求被取消:', error.message);
      return Promise.reject(error);
    }

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
      default:
        message.error(error.message || '网络异常，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

export default client;
