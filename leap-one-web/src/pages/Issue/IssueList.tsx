/** 工单列表页面 */

import { useState } from 'react';
import { Table, Button, Tag, Space, Select } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';

interface IssueItem {
  id: number;
  title: string;
  type: string;
  priority: string;
  status: string;
  requester: string;
  assignee: string;
  createdAt: string;
}

const mockIssues: IssueItem[] = [
  { id: 1, title: '申请开通生产环境SSH权限', type: '运维支持', priority: 'high', status: '处理中', requester: '张三', assignee: '运维组', createdAt: '2026-06-07' },
  { id: 2, title: '需要新增数据库只读账号', type: '资源申请', priority: 'medium', status: '待处理', requester: '李四', assignee: '-', createdAt: '2026-06-06' },
  { id: 3, title: '服务器磁盘空间告警', type: '故障报告', priority: 'urgent', status: '已解决', requester: '监控系统', assignee: '运维组', createdAt: '2026-06-05' },
  { id: 4, title: 'VPN访问速度慢问题', type: '问题咨询', priority: 'low', status: '已关闭', requester: '王五', assignee: 'IT支持', createdAt: '2026-04-20' },
  { id: 5, title: '申请采购MacBook开发机', type: '资源申请', priority: 'medium', status: '待审批', requester: '赵六', assignee: '行政部', createdAt: '2026-06-08' },
];

export default function IssueList() {
  const [searchValue, setSearchValue] = useState('');
  const [typeFilter, setTypeFilter] = useState<string | undefined>();

  const columns: ColumnsType<IssueItem> = [
    { title: '工单标题', dataIndex: 'title', render: (t: string) => <strong>{t}</strong>, ellipsis: true },
    { title: '类型', dataIndex: 'type', width: 100, render: (t: string) => <Tag>{t}</Tag> },
    {
      title: '优先级', dataIndex: 'priority', width: 90,
      render: (p: string) => <Tag color={p === 'urgent' ? 'red' : p === 'high' ? 'orange' : p === 'medium' ? 'blue' : 'default'}>{p}</Tag>,
    },
    {
      title: '状态', dataIndex: 'status', width: 95,
      render: (s: string) => <Tag color={s === '处理中' || s === '待审批' ? 'processing' : s === '已解决' ? 'success' : s === '待处理' ? 'warning' : 'default'}>{s}</Tag>,
    },
    { title: '申请人', dataIndex: 'requester', width: 90 },
    { title: '处理人', dataIndex: 'assignee', width: 90 },
    { title: '创建时间', dataIndex: 'createdAt', width: 110 },
    {
      title: '操作', key: 'action', width: 150,
      render: () => <Space><Button type="link" size="small">处理</Button><Button type="link" size="small">详情</Button></Space>,
    },
  ];

  return (
    <div>
      <PageHeader
        title="工单管理"
        subtitle="跟踪和处理各类服务请求与问题"
        extra={<Button type="primary" icon={<PlusOutlined />}>提交工单</Button>}
      />
      <div style={{ marginBottom: 16, display: 'flex', gap: 12, alignItems: 'center' }}>
        <SearchBar value={searchValue || ''} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索工单..." />
        <Select placeholder="类型筛选" allowClear style={{ width: 130 }} value={typeFilter} onChange={setTypeFilter}
          options={[
            { label: '运维支持', value: '运维支持' }, { label: '资源申请', value: '资源申请' },
            { label: '故障报告', value: '故障报告' }, { label: '问题咨询', value: '问题咨询' },
          ]} />
      </div>
      <Table rowKey="id" columns={columns} dataSource={mockIssues} pagination={{ pageSize: 10 }} />
    </div>
  );
}
