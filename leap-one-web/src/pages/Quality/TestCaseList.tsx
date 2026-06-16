/** 测试用例列表页面 */

import { Table, Button, Tag, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';

interface TestCaseItem {
  id: number;
  title: string;
  module: string;
  priority: string;
  type: string;
  author: string;
  status: string;
  lastRunDate: string;
}

const mockTestCases: TestCaseItem[] = [
  { id: 1, title: '登录功能-正常账号登录', module: '认证模块', priority: 'P0', type: '功能测试', author: '赵六', status: '通过', lastRunDate: '2026-06-07' },
  { id: 2, title: '登录功能-密码错误提示', module: '认证模块', priority: 'P0', type: '功能测试', author: '赵六', status: '通过', lastRunDate: '2026-06-07' },
  { id: 3, title: '权限控制-无权限访问拦截', module: '权限模块', priority: 'P0', type: '安全测试', author: '钱七', status: '失败', lastRunDate: '2026-06-06' },
  { id: 4, title: '数据导出-Excel格式正确性', module: '报表模块', priority: 'P1', type: '功能测试', author: '赵六', status: '未执行', lastRunDate: '-' },
  { id: 5, title: '并发请求-接口压力测试', module: '性能测试', priority: 'P2', type: '性能测试', author: '孙八', status: '未执行', lastRunDate: '-' },
];

export default function TestCaseList() {
  const columns: ColumnsType<TestCaseItem> = [
    { title: '用例名称', dataIndex: 'title', render: (t: string) => <strong>{t}</strong>, ellipsis: true },
    { title: '所属模块', dataIndex: 'module', width: 110 },
    { title: '优先级', dataIndex: 'priority', width: 80, render: (p: string) => <Tag color={p === 'P0' ? 'red' : p === 'P1' ? 'orange' : 'blue'}>{p}</Tag> },
    { title: '类型', dataIndex: 'type', width: 100 },
    { title: '作者', dataIndex: 'author', width: 80 },
    {
      title: '状态', dataIndex: 'status', width: 85,
      render: (s: string) => <Tag color={s === '通过' ? 'success' : s === '失败' ? 'error' : 'default'}>{s}</Tag>,
    },
    { title: '最近执行', dataIndex: 'lastRunDate', width: 110 },
    {
      title: '操作', key: 'action', width: 120,
      render: () => <Space><Button type="link" size="small">执行</Button><Button type="link" size="small">编辑</Button></Space>,
    },
  ];

  return (
    <div>
      <PageHeader
        title="测试用例"
        subtitle="管理和维护所有测试用例"
        extra={<Button type="primary" icon={<PlusOutlined />}>新建用例</Button>}
      />
      <SearchBar value={''} onChange={() => {}} onSearch={() => {}} placeholder="搜索用例..." />
      <Table rowKey="id" columns={columns} dataSource={mockTestCases} pagination={{ pageSize: 10 }} />
    </div>
  );
}
