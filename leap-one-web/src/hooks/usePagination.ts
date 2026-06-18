/**
 * 分页状态管理 Hook
 *
 * 提供统一的分页状态和处理函数，兼容 Ant Design Table
 */

import { useState, useCallback, useMemo } from 'react';

// ── 类型定义 ─────────────────────────────────────────────────

interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

interface UsePaginationOptions {
  /** 默认每页条数 */
  defaultPageSize?: number;
}

interface UsePaginationReturn {
  /** 当前页码 */
  page: number;
  /** 每页条数 */
  pageSize: number;
  /** 总记录数 */
  total: number;
  /** 切换页码 */
  changePage: (page: number) => void;
  /** 切换每页条数 */
  changePageSize: (_current: number, size: number) => void;
  /** 设置总记录数 */
  setTotal: (total: number) => void;
  /** 重置分页到第一页 */
  reset: () => void;
  /** Ant Design Table pagination 配置 */
  tablePagination: {
    current: number;
    pageSize: number;
    total: number;
    showSizeChanger: boolean;
    showQuickJumper: boolean;
    showTotal: (total: number) => string;
    onChange: (page: number) => void;
    onShowSizeChange: (_current: number, size: number) => void;
  };
}

// ── 默认配置 ─────────────────────────────────────────────────

const DEFAULT_PAGE_SIZE = 10;

// ── Hook 实现 ────────────────────────────────────────────────

const usePagination = (options: UsePaginationOptions = {}): UsePaginationReturn => {
  const { defaultPageSize = DEFAULT_PAGE_SIZE } = options;

  const [state, setState] = useState<PaginationState>({
    page: 1,
    pageSize: defaultPageSize,
    total: 0,
  });

  // 页码切换
  const changePage = useCallback((targetPage: number): void => {
    setState((prev) => ({ ...prev, page: targetPage }));
  }, []);

  // 每页条数切换
  const changePageSize = useCallback(
    (_current: number, newSize: number): void => {
      setState({ page: 1, pageSize: newSize, total: state.total });
    },
    [state.total]
  );

  // 设置总记录数
  const setTotal = useCallback((newTotal: number): void => {
    setState((prev) => ({ ...prev, total: newTotal }));
  }, []);

  // 重置分页
  const reset = useCallback((): void => {
    setState({ page: 1, pageSize: defaultPageSize, total: 0 });
  }, [defaultPageSize]);

  // Ant Design Table 兼容的 pagination 配置
  const tablePagination = useMemo(
    () => ({
      current: state.page,
      pageSize: state.pageSize,
      total: state.total,
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: (totalCount: number) => `共 ${totalCount} 条`,
      onChange: changePage,
      onShowSizeChange: changePageSize,
    }),
    [state.page, state.pageSize, state.total, changePage, changePageSize]
  );

  return {
    page: state.page,
    pageSize: state.pageSize,
    total: state.total,
    changePage,
    changePageSize,
    setTotal,
    reset,
    tablePagination,
  };
};

export default usePagination;
