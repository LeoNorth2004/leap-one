/** 产品列表页面 */

import { useState } from 'react';
import { Table, Button, Space, Select } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import StatusTag from '@/components/Common/StatusTag';
import type { Product } from '@/types/product';

export default function ProductList() {
  const [searchValue, setSearchValue] = useState('');
  const [statusFilter, setStatusFilter] = useState<string | undefined>();

  /** 模拟产品数据 - TODO: 替换为真实API调用 */
  const mockProducts: Product[] = [
    { id: 1, name: '企业管理系统', code: 'PRD-001', description: '面向企业的综合管理平台', status: 'normal', managerId: 1, managerName: '张三', createdAt: '2025-01-15', updatedAt: '2025-06-01' },
    { id: 2, name: '移动办公App', code: 'PRD-002', description: '移动端协同办公应用', status: 'normal', managerId: 2, managerName: '李四', createdAt: '2025-02-20', updatedAt: '2025-05-28' },
    { id: 3, name: '数据分析平台', code: 'PRD-003', description: 'BI数据分析与可视化工具', status: 'normal', managerId: 3, managerName: '王五', createdAt: '2025-03-10', updatedAt: '2025-06-05' },
    { id: 4, name: '旧版CRM', code: 'PRD-004', description: '已停用的客户关系管理系统', status: 'closed', managerId: 1, managerName: '张三', createdAt: '2024-06-01', updatedAt: '2025-01-01' },
  ];

  const columns: ColumnsType<Product> = [
    { title: '产品名称', dataIndex: 'name', render: (name: string) => <strong>{name}</strong> },
    { title: '产品代号', dataIndex: 'code' },
    { title: '描述', dataIndex: 'description', ellipsis: true },
    { title: '产品经理', dataIndex: 'managerName' },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (status: Product['status']) => (
        <StatusTag status={status} statusMap={{
          normal: { label: '正常', color: 'success' },
          closed: { label: '已关闭', color: 'default' },
        }} />
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, _record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />}>
            查看
          </Button>
          <Button type="link" size="small">路线图</Button>
          <Button type="link" size="small">编辑</Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="产品管理"
        subtitle="管理所有产品的需求与路线图"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            新建产品
          </Button>
        }
      />

      <div style={{ marginBottom: 16, display: 'flex', gap: 12, alignItems: 'center' }}>
        <SearchBar value={searchValue} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索产品..." />
        <Select
          placeholder="状态筛选"
          allowClear
          style={{ width: 130 }}
          value={statusFilter}
          onChange={setStatusFilter}
          options={[
            { label: '正常', value: 'normal' },
            { label: '已关闭', value: 'closed' },
          ]}
        />
      </div>

      <Table rowKey="id" columns={columns} dataSource={mockProducts} pagination={{ pageSize: 10 }} />
    </div>
  );
}
