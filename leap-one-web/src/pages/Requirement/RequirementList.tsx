/** 需求列表页面 */

import { useState } from 'react';
import { Table, Button, Select, Space } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import StatusTag from '@/components/Common/StatusTag';
import PriorityBadge from '@/components/Business/PriorityBadge';
import type { Requirement, RequirementStatus, Priority } from '@/types/requirement';

export default function RequirementList() {
  const [searchValue, setSearchValue] = useState('');
  const [statusFilter, setStatusFilter] = useState<RequirementStatus | undefined>();
  const [priorityFilter, setPriorityFilter] = useState<Priority | undefined>();

  /** 模拟需求数据 - TODO: 替换为真实API */
  const mockRequirements: Requirement[] = [
    { id: 1, title: '用户权限RBAC模型实现', code: 'REQ-001', description: '基于角色的访问控制系统', status: 'developing', priority: 'P0', source: 'internal', productId: 1, productName: '企业管理系统', moduleId: 1, moduleName: '权限管理', storyPoints: 8, assigneeId: 1, assigneeName: '张三', planRelease: 'v2.0.0', createdAt: '2025-06-01', updatedAt: '2026-06-05', createdBy: '产品经理' },
    { id: 2, title: '数据导出Excel功能', code: 'REQ-002', description: '支持表格数据导出为Excel格式', status: 'reviewing', priority: 'P1', source: 'customer', productId: 1, productName: '企业管理系统', planRelease: 'v2.0.0', createdAt: '2025-06-10', updatedAt: '2026-06-03', createdBy: '产品经理' },
    { id: 3, title: '移动端适配优化', code: 'REQ-003', description: '响应式布局支持移动设备访问', status: 'active', priority: 'P1', source: 'market', productId: 2, productName: '移动办公App', assigneeId: 2, assigneeName: '李四', planRelease: 'v1.5.0', createdAt: '2025-07-01', updatedAt: '2026-05-28', createdBy: '产品经理' },
    { id: 4, title: '实时消息通知推送', code: 'REQ-004', description: 'WebSocket实时消息推送机制', status: 'testing', priority: 'P2', source: 'internal', productId: 1, productName: '企业管理系统', assigneeId: 3, assigneeName: '王五', planRelease: 'v2.1.0', createdAt: '2025-08-15', updatedAt: '2026-06-06', createdBy: '技术负责人' },
    { id: 5, title: '多语言国际化支持', code: 'REQ-005', description: '支持中文/英文切换显示', status: 'draft', priority: 'P3', source: 'competitive', productId: 1, productName: '企业管理系统', createdAt: '2026-01-10', updatedAt: '2026-06-01', createdBy: '产品经理' },
  ];

  const columns: ColumnsType<Requirement> = [
    { title: '需求标题', dataIndex: 'title', render: (title: string) => <strong>{title}</strong>, ellipsis: true },
    { title: '编号', dataIndex: 'code', width: 95 },
    { title: '产品', dataIndex: 'productName', width: 110 },
    {
      title: '优先级', dataIndex: 'priority', width: 90,
      render: (priority: Priority) => <PriorityBadge priority={priority} />,
    },
    {
      title: '状态', dataIndex: 'status', width: 100,
      render: (status: RequirementStatus) => <StatusTag status={status} />,
    },
    { title: '指派人', dataIndex: 'assigneeName', width: 90 },
    { title: '计划版本', dataIndex: 'planRelease', width: 100 },
    { title: '更新时间', dataIndex: 'updatedAt', width: 110 },
    {
      title: '操作', key: 'action', width: 150,
      render: () => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />}>查看</Button>
          <Button type="link" size="small">编辑</Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="需求管理"
        subtitle="管理所有产品的需求全生命周期"
        extra={<Button type="primary" icon={<PlusOutlined />}>新建需求</Button>}
      />

      <div style={{ marginBottom: 16, display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'center' }}>
        <SearchBar value={searchValue} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索需求..." />
        <Select placeholder="状态" allowClear style={{ width: 120 }} value={statusFilter} onChange={setStatusFilter}
          options={[
            { label: '草稿', value: 'draft' }, { label: '评审中', value: 'reviewing' },
            { label: '激活', value: 'active' }, { label: '开发中', value: 'developing' },
            { label: '测试中', value: 'testing' }, { label: '已完成', value: 'completed' },
          ]} />
        <Select placeholder="优先级" allowClear style={{ width: 110 }} value={priorityFilter} onChange={setPriorityFilter}
          options={[{ label: 'P0-紧急', value: 'P0' }, { label: 'P1-高', value: 'P1' }, { label: 'P2-中', value: 'P2' }, { label: 'P3-低', value: 'P3' }]} />
      </div>

      <Table rowKey="id" columns={columns} dataSource={mockRequirements} pagination={{ pageSize: 10 }} scroll={{ x: 1100 }} />
    </div>
  );
}
