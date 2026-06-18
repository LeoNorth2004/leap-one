/**
 * 带分页表格组件
 *
 * 封装 Ant Design Table 的常用配置：
 * - 内置分页（支持切换每页条数、快速跳转）
 * - 默认中等尺寸、横向滚动
 */

import type { TableProps } from 'antd';
import { Table } from 'antd';

// ── 类型定义 ─────────────────────────────────────────────────

interface TableWithPaginationProps<T> extends Omit<TableProps<T>, 'onChange' | 'pagination'> {
  /** 数据源 */
  dataSource: T[];
  /** 总记录数 */
  total: number;
  /** 当前页码 */
  current: number;
  /** 每页条数 */
  pageSize: number;
  /** 分页变化回调 */
  onChange: (page: number, pageSize: number) => void;
  /** 加载状态 */
  loading?: boolean;
  /** 行唯一标识字段或函数 */
  rowKey?: string | ((record: T) => string);
}

// ── 默认配置 ─────────────────────────────────────────────────

const DEFAULT_ROW_KEY = 'id';

const buildPaginationConfig = (
  current: number,
  pageSize: number,
  total: number,
  onChange: (page: number, size: number) => void
) => ({
  current,
  pageSize,
  total,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (t: number) => `共 ${t} 条`,
  onChange,
});

// ── 组件实现 ─────────────────────────────────────────────────

function TableWithPagination<T extends Record<string, unknown>>({
  columns,
  dataSource,
  total,
  current,
  pageSize,
  onChange,
  loading = false,
  rowKey = DEFAULT_ROW_KEY,
  ...restProps
}: TableWithPaginationProps<T>) {
  const pagination = buildPaginationConfig(current, pageSize, total, onChange);

  return (
    <Table<T>
      columns={columns}
      dataSource={dataSource}
      rowKey={rowKey}
      pagination={pagination}
      loading={loading}
      size="middle"
      scroll={{ x: 'max-content' }}
      {...restProps}
    />
  );
}

export default TableWithPagination;
