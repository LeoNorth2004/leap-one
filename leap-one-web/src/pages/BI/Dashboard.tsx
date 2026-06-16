/** BI数据大屏页面 */

import { Row, Col, Card, Statistic, Progress, Table } from 'antd';
import {
  ArrowUpOutlined,
  ArrowDownOutlined,
  ProjectOutlined,
  FileTextOutlined,
  CarryOutOutlined,
  BugOutlined,
} from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';

export default function BIDashboard() {
  /** 模拟统计数据 */
  const overviewStats = [
    { title: '活跃项目', value: 5, icon: <ProjectOutlined />, color: '#1677ff', trend: 2 },
    { title: '总需求数', value: 128, icon: <FileTextOutlined />, color: '#13c2c2', trend: 15 },
    { title: '进行中任务', value: 67, icon: <CarryOutOutlined />, color: '#faad14', trend: -3 },
    { title: '待修复Bug', value: 23, icon: <BugOutlined />, color: '#ff4d4f', trend: -5 },
  ];

  /** 项目进度数据 */
  const projectProgress = [
    { name: '企业管理系统V2', progress: 72, tasks: 156, completed: 112 },
    { name: '移动办公App', progress: 45, tasks: 89, completed: 40 },
    { name: '数据中台建设', progress: 88, tasks: 67, completed: 59 },
    { name: '内部工具链', progress: 30, tasks: 34, completed: 10 },
  ];

  /** 团队效能排行 */
  const teamRanking = [
    { rank: 1, name: '前端组', members: 6, completedTasks: 45, avgHours: 6.2 },
    { rank: 2, name: '后端组', members: 8, completedTasks: 52, avgHours: 7.8 },
    { rank: 3, name: '测试组', members: 4, completedTasks: 38, avgHours: 5.5 },
    { rank: 4, name: '产品组', members: 3, completedTasks: 20, avgHours: 4.0 },
  ];

  const rankingColumns = [
    { title: '排名', dataIndex: 'rank', width: 60, align: 'center' as const },
    { title: '团队', dataIndex: 'name' },
    { title: '人数', dataIndex: 'members', width: 70, align: 'center' as const },
    { title: '完成任务', dataIndex: 'completedTasks', width: 90, align: 'center' as const },
    { title: '平均工时(h)', dataIndex: 'avgHours', width: 110, align: 'center' as const },
  ];

  return (
    <div className="bi-dashboard">
      <PageHeader title="BI 数据概览" subtitle="多维度数据分析与可视化展示" />

      {/* 核心指标卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {overviewStats.map((stat) => (
          <Col xs={24} sm={12} lg={6} key={stat.title}>
            <Card className="bi-stat-card" style={{ borderRadius: 12 }}>
              <Statistic
                title={<span><span style={{ marginRight: 8 }}>{stat.icon}</span>{stat.title}</span>}
                value={stat.value}
                prefix={
                  stat.trend > 0 ? (
                    <ArrowUpOutlined style={{ color: '#52c41a', fontSize: 14 }} />
                  ) : (
                    <ArrowDownOutlined style={{ color: '#ff4d4f', fontSize: 14 }} />
                  )
                }
                suffix={
                  <span style={{
                    fontSize: 13,
                    marginLeft: 8,
                    color: stat.trend > 0 ? '#52c41a' : '#ff4d4f',
                  }}>
                    {stat.trend > 0 ? '+' : ''}{stat.trend}
                  </span>
                }
                valueStyle={{ fontSize: 32, fontWeight: 700 }}
              />
            </Card>
          </Col>
        ))}
      </Row>

      {/* 项目进度 + 团队效能 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={14}>
          <Card title="项目进度总览" style={{ borderRadius: 12 }} size="small">
            {projectProgress.map((proj) => (
              <div key={proj.name} style={{ marginBottom: 16 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                  <strong>{proj.name}</strong>
                  <span>{proj.completed}/{proj.tasks} 任务 ({proj.progress}%)</span>
                </div>
                <Progress percent={proj.progress} strokeColor="#1677ff" size="small" />
              </div>
            ))}
          </Card>
        </Col>
        <Col xs={24} lg={10}>
          <Card title="团队效能排行" style={{ borderRadius: 12 }} size="small">
            <Table
              rowKey="name"
              dataSource={teamRanking}
              columns={rankingColumns}
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
