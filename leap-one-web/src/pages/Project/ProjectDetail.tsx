/** 项目详情页面 */

import { useParams, useNavigate } from 'react-router-dom';
import { Descriptions, Card, Tag, Tabs, Progress, Button, Space, Statistic, Row, Col } from 'antd';
import { ArrowLeftOutlined, TeamOutlined, FileTextOutlined, CarryOutOutlined, BugOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';
import StatusTag from '@/components/Common/StatusTag';
import ProgressBadge from '@/components/Business/ProgressBadge';
import UserAvatar from '@/components/Business/UserAvatar';
import type { Project } from '@/types/project';

export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  /** 模拟项目详情 - TODO: 替换为API */
  const project: Project = {
    id: Number(id) || 1,
    name: '企业管理系统V2',
    code: 'PRJ-001',
    description: '企业综合管理平台的全面升级改造，包含前端重构、后端微服务化、数据迁移等多个子项目。',
    status: 'active',
    type: 'scrum',
    pmId: 1,
    pmName: '张三',
    productId: 1,
    productName: '企业管理系统',
    startDate: '2025-01-01',
    endDate: '2026-06-30',
    progress: 72,
    memberCount: 12,
    avatar: '',
    createdAt: '2025-01-01',
    updatedAt: '2026-06-01',
  };

  /** 模拟成员列表 */
  const members = [
    { userId: 1, userName: '张三', avatar: '', role: 'pm', joinedAt: '2025-01-01' },
    { userId: 2, userName: '李四', avatar: '', role: 'developer', joinedAt: '2025-01-05' },
    { userId: 3, userName: '王五', avatar: '', role: 'developer', joinedAt: '2025-01-10' },
    { userId: 4, userName: '赵六', avatar: '', role: 'tester', joinedAt: '2025-01-15' },
    { userId: 5, userName: '钱七', avatar: '', role: 'developer', joinedAt: '2025-02-01' },
  ];

  const roleLabels: Record<string, string> = { pm: '项目经理', developer: '开发者', tester: '测试者', observer: '观察者' };

  return (
    <div>
      <PageHeader
        title={`${project.name} (${project.code})`}
        subtitle={project.description}
        extra={
          <Space>
            <Button onClick={() => navigate(-1)}><ArrowLeftOutlined /> 返回</Button>
            <Button type="primary">编辑项目</Button>
          </Space>
        }
      />

      {/* 统计概览 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card><Statistic title="总需求" value={28} prefix={<FileTextOutlined />} /></Card>
        </Col>
        <Col span={6}>
          <Card><Statistic title="总任务" value={156} prefix={<CarryOutOutlined />} /></Card>
        </Col>
        <Col span={6}>
          <Card><Statistic title="待修复Bug" value={7} prefix={<BugOutlined />} /></Card>
        </Col>
        <Col span={6}>
          <Card><Statistic title="团队成员" value={project.memberCount} prefix={<TeamOutlined />} /></Card>
        </Col>
      </Row>

      {/* 详情Tab */}
      <Card>
        <Tabs
          defaultActiveKey="info"
          items={[
            {
              key: 'info',
              label: '基本信息',
              children: (
                <Descriptions column={2} bordered size="middle">
                  <Descriptions.Item label="项目名称">{project.name}</Descriptions.Item>
                  <Descriptions.Item label="项目代号">{project.code}</Descriptions.Item>
                  <Descriptions.Item label="所属产品">{project.productName}</Descriptions.Item>
                  <Descriptions.Item label="项目经理">{project.pmName}</Descriptions.Item>
                  <Descriptions.Item label="项目类型"><Tag>{project.type.toUpperCase()}</Tag></Descriptions.Item>
                  <Descriptions.Item label="项目状态"><StatusTag status={project.status} /></Descriptions.Item>
                  <Descriptions.Item label="开始日期">{project.startDate}</Descriptions.Item>
                  <Descriptions.Item label="截止日期">{project.endDate}</Descriptions.Item>
                  <Descriptions.Item label="整体进度" span={2}><ProgressBadge percent={project.progress} size="default" /></Descriptions.Item>
                  <Descriptions.Item label="描述" span={2}>{project.description}</Descriptions.Item>
                </Descriptions>
              ),
            },
            {
              key: 'members',
              label: `团队成员 (${members.length})`,
              children: (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: 12 }}>
                  {members.map((m) => (
                    <Card key={m.userId} size="small" hoverable>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                        <UserAvatar src={m.avatar} name={m.userName} size={40} />
                        <div>
                          <strong>{m.userName}</strong>
                          <br />
                          <Tag color="blue" style={{ marginTop: 4 }}>{roleLabels[m.role]}</Tag>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              ),
            },
            {
              key: 'iterations',
              label: '迭代记录',
              children: <p style={{ color: '#999' }}>暂无迭代记录，请前往迭代管理创建。</p>,
            },
          ]}
        />
      </Card>
    </div>
  );
}
