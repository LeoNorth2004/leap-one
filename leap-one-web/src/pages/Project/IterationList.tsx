/** 迭代列表页面 */

import { useState } from 'react';
import { Table, Button, Space } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import StatusTag from '@/components/Common/StatusTag';
import ProgressBadge from '@/components/Business/ProgressBadge';
import type { Iteration } from '@/types/project';

export default function IterationList() {
  /** 模拟迭代数据 - TODO: 替换为真实API */
  const [iterations] = useState<Iteration[]>([
    { id: 1, projectId: 1, name: 'Sprint-2026-W23', status: 'active', startDate: '2026-06-02', endDate: '2026-06-15', goal: '完成用户管理模块和权限体系开发', taskCount: 18, completedTaskCount: 14, progress: 78 },
    { id: 2, projectId: 1, name: 'Sprint-2026-W22', status: 'completed', startDate: '2026-05-19', endDate: '2026-06-01', goal: '完成报表模块和数据分析功能', taskCount: 22, completedTaskCount: 22, progress: 100 },
    { id: 3, projectId: 1, name: 'Sprint-2026-W21', status: 'completed', startDate: '2026-05-05', endDate: '2026-05-18', goal: '完成基础UI组件库搭建', taskCount: 15, completedTaskCount: 15, progress: 100 },
    { id: 4, projectId: 1, name: 'Sprint-2026-W20', status: 'completed', startDate: '2026-04-21', endDate: '2026-05-04', goal: '完成登录认证和路由框架', taskCount: 12, completedTaskCount: 11, progress: 92 },
    { id: 5, projectId: 1, name: 'Sprint-2026-W24', status: 'pending', startDate: '2026-06-16', endDate: '2026-06-29', goal: '完成消息通知和工单系统', taskCount: 0, completedTaskCount: 0, progress: 0 },
  ]);

  const columns: ColumnsType<Iteration> = [
    { title: '迭代名称', dataIndex: 'name', render: (name: string) => <strong>{name}</strong> },
    { title: '目标', dataIndex: 'goal', ellipsis: true },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: Iteration['status']) => (
        <StatusTag status={status} statusMap={{
          pending: { label: '未开始', color: 'default' },
          active: { label: '进行中', color: 'processing' },
          completed: { label: '已完成', color: 'success' },
        }} />
      ),
    },
    { title: '开始时间', dataIndex: 'startDate', width: 120 },
    { title: '结束时间', dataIndex: 'endDate', width: 120 },
    {
      title: '进度',
      dataIndex: 'progress',
      width: 160,
      render: (progress: number, record) => (
        <span>
          <ProgressBadge percent={progress} />
          <span style={{ marginLeft: 8, fontSize: 13 }}>{record.completedTaskCount}/{record.taskCount}</span>
        </span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
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
        title="迭代管理"
        subtitle="规划和跟踪每个Sprint的执行情况"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            创建迭代
          </Button>
        }
      />

      <Table rowKey="id" columns={columns} dataSource={iterations} pagination={false} />
    </div>
  );
}
