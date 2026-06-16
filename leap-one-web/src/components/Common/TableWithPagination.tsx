/** 带分页表格组件 - 封装Ant Design Table常用配置 */

import type { TableProps } from 'antd';
import { Table } from 'antd';

interface TableWithPaginationProps<T> extends Omit<TableProps<T>, 'onChange' | 'pagination'> {
  /** 数据源 */
  dataSource: T[];
  /** 总记录数 */
  total: number;
  /** 当前页码 */
  current: number;
  /** 每页条数 */
  pageSize: number;
  /** 切换页码 */
  onChange: (page: number, pageSize: number) => void;
  /** 是否加载中 */
  loading?: boolean;
  /** 行key */
  rowKey?: string | ((record: T) => string);
}

export default function TableWithPagination<T extends Record<string, unknown>>({
  columns,
  dataSource,
  total,
  current,
  pageSize,
  onChange,
  loading = false,
  rowKey = 'id',
  ...restProps
}: TableWithPaginationProps<T>) {
  const paginationConfig = {
    current,
    pageSize,
    total,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (t: number) => `共 ${t} 条`,
    onChange,
  };

  return (
    <Table<T>
      columns={columns}
      dataSource={dataSource}
      rowKey={rowKey}
      pagination={paginationConfig}
      loading={loading}
      size="middle"
      scroll={{ x: 'max-content' }}
      {...restProps}
    />
  );
}
