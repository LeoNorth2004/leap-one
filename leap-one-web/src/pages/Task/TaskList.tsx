/** 任务列表页面 */

import { useState } from 'react';
import { Table, Button, Select, Space } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import StatusTag from '@/components/Common/StatusTag';
import PriorityBadge from '@/components/Business/PriorityBadge';
import UserAvatar from '@/components/Business/UserAvatar';
import type { Task, TaskStatus, TaskPriority } from '@/types/task';

export default function TaskList() {
  const [searchValue, setSearchValue] = useState('');
  const [statusFilter, setStatusFilter] = useState<TaskStatus | undefined>();
  const [priorityFilter, setPriorityFilter] = useState<TaskPriority | undefined>();

  /** 模拟任务数据 - TODO: 替换为真实API */
  const mockTasks: Task[] = [
    { id: 1, title: '实现用户认证中间件', description: 'JWT Token校验', status: 'doing', type: 'dev', priority: 'urgent', projectId: 1, projectName: '企业管理系统V2', iterationId: 1, iterationName: 'Sprint-W23', requirementId: 1, requirementTitle: 'RBAC模型实现', assigneeId: 1, assigneeName: '张三', createdBy: 2, createdByName: '李四', estimatedHours: 8, consumedHours: 5, leftHours: 3, startDate: '2026-06-02', dueDate: '2026-06-08', createdAt: '2026-06-02', updatedAt: '2026-06-07' },
    { id: 2, title: '编写权限API单元测试', description: '', status: 'wait', type: 'test', priority: 'high', projectId: 1, projectName: '企业管理系统V2', iterationId: 1, iterationName: 'Sprint-W23', assigneeId: 4, assigneeName: '赵六', createdBy: 1, createdByName: '张三', estimatedHours: 4, startDate: '2026-06-10', dueDate: '2026-06-12', createdAt: '2026-06-05', updatedAt: '2026-06-05' },
    { id: 3, title: '优化数据库查询性能', description: '慢SQL优化', status: 'done', type: 'dev', priority: 'medium', projectId: 1, projectName: '企业管理系统V2', iterationId: 2, iterationName: 'Sprint-W22', assigneeId: 3, assigneeName: '王五', createdBy: 1, createdByName: '张三', estimatedHours: 6, consumedHours: 7, finishedDate: '2026-05-30', createdAt: '2026-05-20', updatedAt: '2026-05-30' },
    { id: 5, title: '设计系统UI规范文档', description: '', status: 'pause', type: 'design', priority: 'low', projectId: 1, projectName: '企业管理系统V2', iterationId: 1, iterationName: 'Sprint-W23', assigneeId: 5, assigneeName: '钱七', createdBy: 1, createdByName: '张三', estimatedHours: 4, createdAt: '2026-06-01', updatedAt: '2026-06-04' },
    { id: 5, title: '代码审查：权限模块PR', description: '', status: 'doing', type: 'review', priority: 'high', projectId: 1, projectName: '企业管理系统V2', iterationId: 1, iterationName: 'Sprint-W23', assigneeId: 2, assigneeName: '李四', createdBy: 1, createdByName: '张三', estimatedHours: 2, consumedHours: 1, createdAt: '2026-06-06', updatedAt: '2026-06-07' },
  ];

  const columns: ColumnsType<Task> = [
    { title: '任务标题', dataIndex: 'title', render: (title: string) => <strong>{title}</strong>, ellipsis: true },
    { title: '所属项目', dataIndex: 'projectName', width: 130 },
    { title: '迭代', dataIndex: 'iterationName', width: 110 },
    {
      title: '优先级', dataIndex: 'priority', width: 90,
      render: (p: TaskPriority) => <PriorityBadge priority={p} />,
    },
    {
      title: '状态', dataIndex: 'status', width: 95,
      render: (s: TaskStatus) => <StatusTag status={s} />,
    },
    {
      title: '指派人', dataIndex: 'assigneeName', width: 95,
      render: (name?: string) => name ? <UserAvatar name={name} size={24} /> : '-',
    },
    { title: '截止日期', dataIndex: 'dueDate', width: 110 },
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
        title="任务管理"
        subtitle="跟踪和分配所有开发、测试等任务"
        extra={<Button type="primary" icon={<PlusOutlined />}>新建任务</Button>}
      />

      <div style={{ marginBottom: 16, display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'center' }}>
        <SearchBar value={searchValue} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索任务..." />
        <Select placeholder="状态" allowClear style={{ width: 120 }} value={statusFilter} onChange={setStatusFilter}
          options={[
            { label: '待处理', value: 'wait' }, { label: '进行中', value: 'doing' },
            { label: '已完成', value: 'done' }, { label: '已暂停', value: 'pause' },
          ]} />
        <Select placeholder="优先级" allowClear style={{ width: 110 }} value={priorityFilter} onChange={setPriorityFilter}
          options={[{ label: '紧急', value: 'urgent' }, { label: '高', value: 'high' }, { label: '中', value: 'medium' }, { label: '低', value: 'low' }]} />
      </div>

      <Table rowKey="id" columns={columns} dataSource={mockTasks} pagination={{ pageSize: 10 }} scroll={{ x: 1000 }} />
    </div>
  );
}
