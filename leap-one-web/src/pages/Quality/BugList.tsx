/** Bug列表页面 */

import { useState } from 'react';
import { Table, Button, Tag, Space, Select } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import UserAvatar from '@/components/Business/UserAvatar';

interface BugItem {
  id: number;
  title: string;
  severity: string;
  status: string;
  project: string;
  reporter: string;
  assignee: string;
  createdAt: string;
}

const severityConfig: Record<string, { color: string; label: string }> = {
  fatal: { color: 'red', label: '致命' },
  serious: { color: 'orange', label: '严重' },
  normal: { color: 'blue', label: '一般' },
  slight: { color: 'default', label: '轻微' },
  suggest: { color: 'green', label: '建议' },
};

const mockBugs: BugItem[] = [
  { id: 1, title: '首页白屏 - 生产环境偶现', severity: 'fatal', status: 'open', project: '企业管理系统V2', reporter: '赵六', assignee: '张三', createdAt: '2026-06-07' },
  { id: 2, title: '导出Excel文件名乱码', severity: 'serious', status: 'fixed', project: '企业管理系统V2', reporter: '客户反馈', assignee: '李四', createdAt: '2026-06-05' },
  { id: 3, title: '移动端侧边栏无法收起', severity: 'normal', status: 'resolved', project: '移动办公App', reporter: '内部测试', assignee: '王五', createdAt: '2026-06-03' },
  { id: 4, title: '搜索结果分页不生效', severity: 'normal', status: 'open', project: '企业管理系统V2', reporter: '赵六', assignee: '张三', createdAt: '2026-06-02' },
  { id: 5, title: '暗色模式下表格边框不可见', severity: 'slight', status: 'closed', project: '企业管理系统V2', reporter: 'UI走查', assignee: '-', createdAt: '2026-05-30' },
];

export default function BugList() {
  const [searchValue, setSearchValue] = useState('');
  const [severityFilter, setSeverityFilter] = useState<string | undefined>();

  const columns: ColumnsType<BugItem> = [
    { title: 'Bug标题', dataIndex: 'title', render: (t: string) => <strong>{t}</strong>, ellipsis: true },
    {
      title: '严重程度', dataIndex: 'severity', width: 95,
      render: (s: string) => <Tag color={severityConfig[s]?.color}>{severityConfig[s]?.label || s}</Tag>,
    },
    {
      title: '状态', dataIndex: 'status', width: 90,
      render: (s: string) => <Tag color={
        s === 'open' ? 'red' : s === 'fixed' ? 'processing' :
        s === 'resolved' ? 'success' : 'default'
      }>{
        s === 'open' ? '待修复' : s === 'fixed' ? '已修复' :
        s === 'resolved' ? '已解决' : s === 'closed' ? '已关闭' : s
      }</Tag>,
    },
    { title: '项目', dataIndex: 'project', width: 140 },
    { title: '报告人', dataIndex: 'reporter', width: 90 },
    { title: '指派给', dataIndex: 'assignee', width: 90,
      render: (a: string) => a !== '-' ? <UserAvatar name={a} size={24} /> : a,
    },
    { title: '创建时间', dataIndex: 'createdAt', width: 110 },
    {
      title: '操作', key: 'action', width: 180,
      render: () => <Space><Button type="link" size="small">查看</Button><Button type="link" size="small">处理</Button></Space>,
    },
  ];

  return (
    <div>
      <PageHeader
        title="Bug管理"
        subtitle="跟踪和管理软件缺陷的全生命周期"
        extra={<Button type="primary" icon={<PlusOutlined />}>提交Bug</Button>}
      />
      <div style={{ marginBottom: 16, display: 'flex', gap: 12, alignItems: 'center' }}>
        <SearchBar value={searchValue || ''} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索Bug..." />
        <Select placeholder="严重程度" allowClear style={{ width: 130 }} value={severityFilter} onChange={setSeverityFilter}
          options={Object.entries(severityConfig).map(([k, v]) => ({ label: v.label, value: k }))} />
      </div>
      <Table rowKey="id" columns={columns} dataSource={mockBugs} pagination={{ pageSize: 10 }} scroll={{ x: 1000 }} />
    </div>
  );
}
