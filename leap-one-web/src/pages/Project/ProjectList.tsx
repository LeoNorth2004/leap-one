/** 项目列表页面 */

import { useState } from 'react';
import { Table, Button, Tag, Select, Space } from 'antd';
import { PlusOutlined, EyeOutlined, SettingOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import StatusTag from '@/components/Common/StatusTag';
import ProgressBadge from '@/components/Business/ProgressBadge';
import type { Project, ProjectStatus, ProjectType } from '@/types/project';

export default function ProjectList() {
  const [searchValue, setSearchValue] = useState('');
  const [statusFilter, setStatusFilter] = useState<ProjectStatus | undefined>();
  const [typeFilter, setTypeFilter] = useState<ProjectType | undefined>();

  /** 模拟项目数据 - TODO: 替换为真实API调用 */
  const mockProjects: Project[] = [
    { id: 1, name: '企业管理系统V2', code: 'PRJ-001', description: '企业综合管理平台升级', status: 'active', type: 'scrum', pmId: 1, pmName: '张三', productId: 1, productName: '企业管理系统', startDate: '2025-01-01', endDate: '2026-06-30', progress: 72, memberCount: 12, avatar: '', createdAt: '2025-01-01', updatedAt: '2026-06-01' },
    { id: 2, name: '移动办公App', code: 'PRJ-002', description: 'iOS/Android双端应用开发', status: 'active', type: 'scrum', pmId: 2, pmName: '李四', productId: 2, productName: '移动办公App', startDate: '2025-03-01', endDate: '2025-12-31', progress: 45, memberCount: 8, avatar: '', createdAt: '2025-03-01', updatedAt: '2026-05-28' },
    { id: 3, name: '数据中台建设', code: 'PRJ-003', description: '数据采集治理分析平台', status: 'active', type: 'waterfall', pmId: 3, pmName: '王五', productId: 3, productName: '数据分析平台', startDate: '2025-06-01', endDate: '2026-03-31', progress: 88, memberCount: 10, avatar: '', createdAt: '2025-06-01', updatedAt: '2026-06-05' },
    { id: 4, name: '旧系统迁移', code: 'PRJ-004', description: '遗留系统向新架构迁移', status: 'completed', type: 'waterfall', pmId: 4, pmName: '赵六', productId: 1, productName: '企业管理系统', startDate: '2024-01-01', endDate: '2025-01-31', progress: 100, memberCount: 15, avatar: '', createdAt: '2024-01-01', updatedAt: '2025-02-01' },
    { id: 5, name: '内部工具链', code: 'PRJ-005', description: '研发效能提升工具集', status: 'paused', type: 'kanban', pmId: 1, pmName: '张三', productId: 0, productName: '', startDate: '2025-09-01', endDate: '2026-09-01', progress: 30, memberCount: 5, avatar: '', createdAt: '2025-09-01', updatedAt: '2026-04-15' },
  ];

  const columns: ColumnsType<Project> = [
    {
      title: '项目名称',
      dataIndex: 'name',
      render: (name: string) => <strong>{name}</strong>,
    },
    { title: '代号', dataIndex: 'code', width: 100 },
    { title: '产品', dataIndex: 'productName' },
    { title: '项目经理', dataIndex: 'pmName' },
    {
      title: '类型',
      dataIndex: 'type',
      width: 90,
      render: (type: ProjectType) => {
        const map: Record<ProjectType, string> = { scrum: 'Scrum', waterfall: '瀑布', kanban: '看板', hybrid: '混合' };
        return <Tag>{map[type]}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (status: ProjectStatus) => <StatusTag status={status} />,
    },
    {
      title: '进度',
      dataIndex: 'progress',
      width: 140,
      render: (progress: number) => <ProgressBadge percent={progress} />,
    },
    { title: '成员', dataIndex: 'memberCount', width: 70, align: 'center', render: (n: number) => `${n}人` },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />}>
            详情
          </Button>
          <Button type="link" size="small">迭代</Button>
          <Button type="link" size="small" icon={<SettingOutlined />}>设置</Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="项目管理"
        subtitle="管理和跟踪所有项目的进度与状态"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            新建项目
          </Button>
        }
      />

      <div style={{ marginBottom: 16, display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'center' }}>
        <SearchBar value={searchValue} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索项目..." />
        <Select placeholder="状态" allowClear style={{ width: 120 }} value={statusFilter} onChange={setStatusFilter}
          options={[{ label: '规划中', value: 'planning' }, { label: '进行中', value: 'active' }, { label: '已暂停', value: 'paused' }, { label: '已完成', value: 'completed' }, { label: '已归档', value: 'archived' }]} />
        <Select placeholder="类型" allowClear style={{ width: 120 }} value={typeFilter} onChange={setTypeFilter}
          options={[{ label: 'Scrum', value: 'scrum' }, { label: '瀑布', value: 'waterfall' }, { label: '看板', value: 'kanban' }, { label: '混合', value: 'hybrid' }]} />
      </div>

      <Table rowKey="id" columns={columns} dataSource={mockProjects} pagination={{ pageSize: 10 }} scroll={{ x: 1100 }} />
    </div>
  );
}
