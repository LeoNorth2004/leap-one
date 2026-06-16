/** 测试计划列表页面 */

import { Table, Button, Tag, Space, Progress } from 'antd';
import { PlusOutlined, PlayCircleOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';

interface TestPlanItem {
  id: number;
  name: string;
  project: string;
  status: string;
  totalCases: number;
  passedCases: number;
  failedCases: number;
  executor: string;
  startDate: string;
  endDate: string;
}

const mockPlans: TestPlanItem[] = [
  { id: 1, name: 'V2.0.0 冒烟测试计划', project: '企业管理系统V2', status: 'executing', totalCases: 50, passedCases: 35, failedCases: 3, executor: '赵六', startDate: '2026-06-01', endDate: '2026-06-15' },
  { id: 2, name: 'V1.5.0 回归测试计划', project: '移动办公App', status: 'pending', totalCases: 30, passedCases: 0, failedCases: 0, executor: '赵六', startDate: '2026-07-01', endDate: '2026-07-15' },
  { id: 3, name: '权限模块专项测试', project: '企业管理系统V2', status: 'completed', totalCases: 25, passedCases: 24, failedCases: 1, executor: '钱七', startDate: '2026-05-20', endDate: '2026-05-28' },
];

export default function TestPlanList() {
  const columns: ColumnsType<TestPlanItem> = [
    { title: '计划名称', dataIndex: 'name', render: (n: string) => <strong>{n}</strong> },
    { title: '关联项目', dataIndex: 'project', width: 160 },
    {
      title: '状态', dataIndex: 'status', width: 100,
      render: (s: string) => <Tag color={s === 'executing' ? 'processing' : s === 'completed' ? 'success' : 'default'}>{
        s === 'executing' ? '执行中' : s === 'completed' ? '已完成' : '待执行'
      }</Tag>,
    },
    { title: '负责人', dataIndex: 'executor', width: 90 },
    { title: '开始日期', dataIndex: 'startDate', width: 115 },
    { title: '结束日期', dataIndex: 'endDate', width: 115 },
    {
      title: '执行进度', width: 200,
      render: (_, r) => (
        <Progress
          percent={Math.round((r.passedCases + r.failedCases) / r.totalCases * 100)}
          success={{ percent: Math.round(r.passedCases / r.totalCases * 100) }}
          format={() => `${r.passedCases + r.failedCases}/${r.totalCases}`}
          size="small"
        />
      ),
    },
    {
      title: '操作', key: 'action', width: 200,
      render: (_, r) => (
        <Space size="small">
          <Button type="link" size="small" icon={<PlayCircleOutlined />} disabled={r.status === 'completed'}>执行</Button>
          <Button type="link" size="small">报告</Button>
          <Button type="link" size="small">详情</Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="测试计划"
        subtitle="规划和跟踪各轮次测试的执行情况"
        extra={<Button type="primary" icon={<PlusOutlined />}>创建计划</Button>}
      />
      <Table rowKey="id" columns={columns} dataSource={mockPlans} pagination={false} />
    </div>
  );
}
