/** 分页管理Hook */

import { useState, useCallback } from 'react';

interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

interface UsePaginationOptions {
  /** 默认每页条数 */
  defaultPageSize?: number;
}

/**
 * 分页状态管理Hook
 * 提供统一的分页状态和处理函数
 */
export function usePagination(options: UsePaginationOptions = {}) {
  const { defaultPageSize = 10 } = options;

  const [pagination, setPagination] = useState<PaginationState>({
    page: 1,
    pageSize: defaultPageSize,
    total: 0,
  });

  /** 切换页码 */
  const changePage = useCallback((page: number) => {
    setPagination((prev) => ({ ...prev, page }));
  }, []);

  /** 切换每页条数 */
  const changePageSize = useCallback((_current: number, size: number) => {
    setPagination({ page: 1, pageSize: size, total: pagination.total });
  }, [pagination.total]);

  /** 设置总记录数（通常在API响应后调用） */
  const setTotal = useCallback((total: number) => {
    setPagination((prev) => ({ ...prev, total }));
  }, []);

  /** 重置分页到第一页 */
  const reset = useCallback(() => {
    setPagination({ page: 1, pageSize: defaultPageSize, total: 0 });
  }, [defaultPageSize]);

  /** 获取Ant Design Table的pagination属性 */
  const tablePagination = {
    current: pagination.page,
    pageSize: pagination.pageSize,
    total: pagination.total,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number) => `共 ${total} 条`,
    onChange: changePage,
    onShowSizeChange: changePageSize,
  };

  return {
    page: pagination.page,
    pageSize: pagination.pageSize,
    total: pagination.total,
    changePage,
    changePageSize,
    setTotal,
    reset,
    tablePagination,
  };
}
