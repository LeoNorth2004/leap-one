/**
 * 工作台 / 个人仪表盘页面
 *
 * 功能：
 * - 顶部统计卡片区（4 个 Statistic Card）：任务数、需求数、Bug 数、工单数
 * - 最近动态时间线（Timeline 组件）
 * - 待办事项快捷列表（我的待办 Top 5）
 * - 日历视图（简要日程）
 * - 使用 TanStack Query 获取数据（含 Skeleton 加载态）
 * - 响应式网格布局
 * - Empty 占位（无数据时）
 */

import { useMemo, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Row,
  Col,
  Card,
  Statistic,
  Timeline,
  Typography,
  Calendar,
  Badge,
  Empty,
  List,
  Tag,
  Spin,
} from 'antd';
import type { Dayjs } from 'dayjs';
import dayjs from 'dayjs';
import {
  CarryOutOutlined,
  FileTextOutlined,
  BugOutlined,
  CustomerServiceOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  WarningOutlined,
  PlusOutlined,
  RocketOutlined,
  AimOutlined,
  BarChartOutlined,
} from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import PageHeader from '@/components/Common/PageHeader';
import { PRIORITY_MAP } from '@/utils/constants';
import type { TodoItem, CalendarEvent, ActivityItem } from '@/types/common';
import styles from './index.module.less';

const { Text } = Typography;

async function fetchDashboardStats(): Promise<{
  taskCount: number;
  requirementCount: number;
  bugCount: number;
  issueCount: number;
  completedTaskCount: number;
  totalTaskCount: number;
  sprintProgress: number;
}> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        taskCount: 12,
        requirementCount: 5,
        bugCount: 3,
        issueCount: 7,
        completedTaskCount: 28,
        totalTaskCount: 40,
        sprintProgress: 70,
      });
    }, 600);
  });
}

async function fetchRecentActivities(): Promise<ActivityItem[]> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        { id: 1, action: '完成了', target: '【前端登录页开发】', targetType: 'task', user: '张三', userAvatar: '', timestamp: '10分钟前', detail: '' },
        { id: 2, action: '提交了 Bug', target: '【首页样式错乱】', targetType: 'bug', user: '李四', userAvatar: '', timestamp: '30分钟前', detail: '' },
        { id: 3, action: '创建了需求', target: '【用户权限管理优化】', targetType: 'requirement', user: '王五', userAvatar: '', timestamp: '1小时前', detail: '' },
        { id: 4, action: '评论了任务', target: '【API 接口文档编写】', targetType: 'task', user: '赵六', userAvatar: '', timestamp: '2小时前', detail: '' },
        { id: 5, action: '标记迭代为已完成', target: '【Sprint-2026-W23】', targetType: 'iteration', user: '系统', userAvatar: '', timestamp: '3小时前', detail: '' },
      ]);
    }, 500);
  });
}

async function fetchTodoItems(): Promise<TodoItem[]> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        { id: 101, title: '完成用户认证模块开发', type: 'task', priority: 'P1', status: 'doing', dueDate: '2026-06-10', projectName: 'Leap One' },
        { id: 102, title: '修复首页数据加载慢的问题', type: 'bug', priority: 'P0', status: 'active', dueDate: '2026-06-09', projectName: 'Leap One' },
        { id: 103, title: '编写 API 接口文档', type: 'task', priority: 'P2', status: 'wait', dueDate: '2026-06-12', projectName: 'Leap One' },
        { id: 104, title: '评审产品路线图 PRD', type: 'requirement', priority: 'P1', status: 'reviewing', dueDate: '2026-06-11', projectName: 'Leap One' },
        { id: 105, title: '处理客户反馈工单 #2031', type: 'issue', priority: 'P2', status: 'active', dueDate: '2026-06-15', projectName: '客服系统' },
      ]);
    }, 400);
  });
}

async function fetchCalendarEvents(): Promise<CalendarEvent[]> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        { id: 1, title: 'Sprint 迭代评审会', date: dayjs().add(1, 'day').format('YYYY-MM-DD'), type: 'meeting', color: '#1677ff' },
        { id: 2, title: 'V1.2 版本发布', date: dayjs().add(4, 'day').format('YYYY-MM-DD'), type: 'deadline', color: '#ff4d4f' },
        { id: 3, title: '需求评审会议', date: dayjs().add(2, 'day').format('YYYY-MM-DD'), type: 'review', color: '#faad14' },
        { id: 4, title: '技术方案讨论', date: dayjs().add(3, 'day').format('YYYY-MM-DD'), type: 'meeting', color: '#1677ff' },
        { id: 5, title: '代码审查截止', date: dayjs().add(5, 'day').format('YYYY-MM-DD'), type: 'deadline', color: '#ff4d4f' },
      ]);
    }, 300);
  });
}

export default function Dashboard() {
  const navigate = useNavigate();

  const { data: statsData, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: fetchDashboardStats,
    staleTime: 5 * 60 * 1000,
  });

  const { data: activities, isLoading: activitiesLoading } = useQuery({
    queryKey: ['dashboard-activities'],
    queryFn: fetchRecentActivities,
    staleTime: 5 * 60 * 1000,
  });

  const { data: todoList, isLoading: todosLoading } = useQuery({
    queryKey: ['dashboard-todos'],
    queryFn: fetchTodoItems,
    staleTime: 5 * 60 * 1000,
  });

  const { data: calendarEvents, isLoading: calendarLoading } = useQuery({
    queryKey: ['dashboard-calendar'],
    queryFn: fetchCalendarEvents,
    staleTime: 10 * 60 * 1000,
  });

  const statCards = useMemo(
    () =>
      statsData
        ? [
            {
              title: '进行中任务',
              value: statsData.taskCount,
              icon: <CarryOutOutlined />,
              color: '#1677ff',
              gradient: 'linear-gradient(135deg, #1677ff 0%, #4096ff 100%)',
              trend: 8,
              suffix: '个',
              delay: 0,
            },
            {
              title: '我的需求',
              value: statsData.requirementCount,
              icon: <FileTextOutlined />,
              color: '#13c2c2',
              gradient: 'linear-gradient(135deg, #13c2c2 0%, #36cfc9 100%)',
              trend: -2,
              suffix: '个',
              delay: 0.1,
            },
            {
              title: '待修复 Bug',
              value: statsData.bugCount,
              icon: <BugOutlined />,
              color: '#ff4d4f',
              gradient: 'linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%)',
              trend: 1,
              suffix: '个',
              delay: 0.2,
            },
            {
              title: '我的工单',
              value: statsData.issueCount,
              icon: <CustomerServiceOutlined />,
              color: '#faad14',
              gradient: 'linear-gradient(135deg, #faad14 0%, #ffc53d 100%)',
              trend: 3,
              suffix: '个',
              delay: 0.3,
            },
          ]
        : [],
    [statsData]
  );

  const timelineItems = useMemo(
    () =>
      activities?.map((item, index) => ({
        color: item.targetType === 'bug'
          ? '#ff4d4f'
          : item.targetType === 'requirement'
            ? '#1677ff'
            : item.targetType === 'task'
              ? '#52c41a'
              : '#faad14',
        children: (
          <div className={styles.timelineContent} style={{ animationDelay: `${index * 0.1}s` }}>
            <div className={styles.timelineMeta}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                {item.timestamp}
              </Text>
            </div>
            <div className={styles.timelineText}>
              <Text strong>{item.user}</Text>
              <Text type="secondary" style={{ margin: '0 4px' }}>{item.action}</Text>
              <Text className={styles.timelineTarget}>{item.target}</Text>
            </div>
          </div>
        ),
      })) ?? [],
    [activities]
  );

  const cellRender = useCallback(
    (date: Dayjs) => {
      const dateStr = date.format('YYYY-MM-DD');
      const events = calendarEvents?.filter((e) => e.date === dateStr) ?? [];

      if (events.length > 0) {
        return (
          <div className={styles.calendarCell}>
            {events.slice(0, 3).map((event) => (
              <Badge
                key={event.id}
                color={event.color ?? '#1677ff'}
                text={
                  <span className={styles.eventText} title={event.title}>
                    {event.title}
                  </span>
                }
                style={{ fontSize: 11 }}
              />
            ))}
            {events.length > 3 && (
              <Text type="secondary" style={{ fontSize: 10 }}>
                +{events.length - 3}
              </Text>
            )}
          </div>
        );
      }
      return null;
    },
    [calendarEvents]
  );

  const getPriorityTag = (priority: string) => {
    const info = PRIORITY_MAP[priority];
    return info ? <Tag color={info.color}>{info.label}</Tag> : null;
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'doing':
      case 'active':
        return <ClockCircleOutlined className={styles.statusIcon} style={{ color: '#1677ff' }} />;
      case 'done':
      case 'resolved':
        return <CheckCircleOutlined className={styles.statusIcon} style={{ color: '#52c41a' }} />;
      case 'wait':
        return <WarningOutlined className={styles.statusIcon} style={{ color: '#faad14' }} />;
      default:
        return <ClockCircleOutlined className={styles.statusIcon} style={{ color: '#86909c' }} />;
    }
  };

  const quickActions = [
    { title: '创建任务', desc: '新建一个任务', icon: <CarryOutOutlined />, path: '/task/list', color: '#1677ff' },
    { title: '提交 Bug', desc: '报告一个问题', icon: <BugOutlined />, path: '/quality/bug', color: '#ff4d4f' },
    { title: '创建工单', desc: '提交服务请求', icon: <CustomerServiceOutlined />, path: '/issue/list', color: '#faad14' },
    { title: '写文档', desc: '编辑项目文档', icon: <FileTextOutlined />, path: '/document/list', color: '#13c2c2' },
  ];

  if (statsLoading && !statsData) {
    return (
      <div className={styles.dashboard}>
        <PageHeader title="工作台" subtitle="加载中..." />
        <Spin size="large" tip="正在加载数据..." style={{ display: 'block', marginTop: 100 }} />
      </div>
    );
  }

  return (
    <div className={styles.dashboard}>
      <PageHeader
        title="工作台"
        subtitle={`欢迎回来，今天是 ${new Date().toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          weekday: 'long',
        })}`}
      />

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {statCards.map((stat) => (
          <Col xs={24} sm={12} lg={6} key={stat.title}>
            <Card
              className={`${styles.statCard} card-hover`}
              style={{ animationDelay: `${stat.delay}s` }}
              bordered={false}
            >
              <div className={styles.statCardInner}>
                <div className={styles.statIconWrapper} style={{ background: stat.gradient }}>
                  <span className={styles.statIcon}>{stat.icon}</span>
                </div>
                <div className={styles.statContent}>
                  <Statistic
                    title={stat.title}
                    value={stat.value}
                    suffix={stat.suffix}
                    valueStyle={{ fontSize: 28, fontWeight: 700, color: stat.color }}
                  />
                  <div className={styles.trend}>
                    {stat.trend > 0 ? (
                      <span className={styles.trendUp}>
                        <ArrowUpOutlined /> +{Math.abs(stat.trend)}
                      </span>
                    ) : (
                      <span className={styles.trendDown}>
                        <ArrowDownOutlined /> {stat.trend}
                      </span>
                    )}
                    <Text type="secondary" style={{ marginLeft: 8, fontSize: 12 }}>
                      较昨日
                    </Text>
                  </div>
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      {statsData && (
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col xs={24} lg={8}>
            <Card className={`${styles.progressCard} card-hover`} bordered={false}>
              <div className={styles.progressHeader}>
                <AimOutlined style={{ color: '#1677ff', fontSize: 20 }} />
                <span className={styles.progressTitle}>迭代进度</span>
              </div>
              <div className={styles.progressContent}>
                <div className={styles.progressRing}>
                  <svg className={styles.progressSvg} viewBox="0 0 100 100">
                    <circle
                      className={styles.progressBg}
                      cx="50"
                      cy="50"
                      r="40"
                      fill="none"
                      stroke="#e5e6eb"
                      strokeWidth="8"
                    />
                    <circle
                      className={styles.progressCircle}
                      cx="50"
                      cy="50"
                      r="40"
                      fill="none"
                      stroke="url(#progressGradient)"
                      strokeWidth="8"
                      strokeLinecap="round"
                      strokeDasharray={`${statsData.sprintProgress * 2.51} 251`}
                      transform="rotate(-90 50 50)"
                    />
                    <defs>
                      <linearGradient id="progressGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                        <stop offset="0%" stopColor="#1677ff" />
                        <stop offset="100%" stopColor="#13c2c2" />
                      </linearGradient>
                    </defs>
                  </svg>
                  <div className={styles.progressText}>
                    <span className={styles.progressValue}>{statsData.sprintProgress}%</span>
                    <span className={styles.progressLabel}>完成</span>
                  </div>
                </div>
                <div className={styles.progressStats}>
                  <div className={styles.progressStatItem}>
                    <span className={styles.progressStatValue}>{statsData.completedTaskCount}</span>
                    <span className={styles.progressStatLabel}>已完成</span>
                  </div>
                  <div className={styles.progressStatDivider} />
                  <div className={styles.progressStatItem}>
                    <span className={styles.progressStatValue}>{statsData.totalTaskCount}</span>
                    <span className={styles.progressStatLabel}>总任务</span>
                  </div>
                </div>
              </div>
            </Card>
          </Col>
          <Col xs={24} lg={16}>
            <Card className={`${styles.chartCard} card-hover`} bordered={false}>
              <div className={styles.chartHeader}>
                <BarChartOutlined style={{ color: '#13c2c2', fontSize: 20 }} />
                <span className={styles.chartTitle}>本周工作趋势</span>
              </div>
              <div className={styles.chartContent}>
                <div className={styles.weekChart}>
                  {['周一', '周二', '周三', '周四', '周五', '周六', '周日'].map((day, idx) => {
                    const heights = [60, 80, 45, 90, 70, 30, 55];
                    return (
                      <div key={day} className={styles.chartBarItem}>
                        <div className={styles.chartBarWrapper}>
                          <div
                            className={styles.chartBar}
                            style={{
                              height: `${heights[idx]}%`,
                              background: `linear-gradient(180deg, #1677ff 0%, #13c2c2 100%)`,
                              animationDelay: `${idx * 0.1}s`,
                            }}
                          />
                        </div>
                        <span className={styles.chartBarLabel}>{day}</span>
                      </div>
                    );
                  })}
                </div>
                <div className={styles.chartLegend}>
                  <span className={styles.legendItem}>
                    <span className={styles.legendDot} style={{ background: '#1677ff' }} />
                    任务完成
                  </span>
                  <span className={styles.legendItem}>
                    <span className={styles.legendDot} style={{ background: '#13c2c2' }} />
                    需求处理
                  </span>
                </div>
              </div>
            </Card>
          </Col>
        </Row>
      )}

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card
            title="最近动态"
            className={`${styles.activityCard} card-hover`}
            loading={activitiesLoading}
            bordered={false}
            extra={
              <button className={styles.viewAllBtn} onClick={() => navigate('/activity')}>
                查看全部 <RocketOutlined style={{ fontSize: 12 }} />
              </button>
            }
          >
            {timelineItems.length > 0 ? (
              <Timeline items={timelineItems} className={styles.timeline} />
            ) : (
              <Empty description="暂无动态" image={Empty.PRESENTED_IMAGE_SIMPLE} />
            )}
          </Card>
        </Col>

        <Col xs={24} lg={12}>
          <Card
            title="我的待办"
            className={`${styles.todoCard} card-hover`}
            loading={todosLoading}
            bordered={false}
            extra={
              <button className={styles.viewAllBtn} onClick={() => navigate('/task/list')}>
                查看全部 <RocketOutlined style={{ fontSize: 12 }} />
              </button>
            }
          >
            {todoList && todoList.length > 0 ? (
              <List
                dataSource={todoList}
                renderItem={(todo) => (
                  <List.Item
                    className={styles.todoItem}
                    actions={[getPriorityTag(todo.priority)]}
                  >
                    <List.Item.Meta
                      avatar={getStatusIcon(todo.status)}
                      title={
                        <a
                          className={styles.todoTitle}
                          onClick={() => {
                            if (todo.type === 'task') navigate('/task/list');
                            else if (todo.type === 'bug') navigate('/quality/bug');
                            else if (todo.type === 'requirement') navigate('/requirement/list');
                            else if (todo.type === 'issue') navigate('/issue/list');
                          }}
                        >
                          {todo.title}
                        </a>
                      }
                      description={
                        <span className={styles.todoMeta}>
                          {todo.projectName && (
                            <Tag color="#1677ff">{todo.projectName}</Tag>
                          )}
                          {todo.dueDate && (
                            <Text type="secondary" style={{ fontSize: 12 }}>
                              截止: {todo.dueDate}
                            </Text>
                          )}
                        </span>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无待办" image={Empty.PRESENTED_IMAGE_SIMPLE} />
            )}
          </Card>
        </Col>

        <Col xs={24} lg={12}>
          <Card
            title="日程安排"
            className={`${styles.calendarCard} card-hover`}
            loading={calendarLoading}
            bordered={false}
          >
            <Calendar
              fullscreen={false}
              cellRender={(date) => cellRender(date as unknown as Dayjs)}
              className={styles.calendar}
            />
          </Card>
        </Col>

        <Col xs={24} lg={12}>
          <Card title="快捷入口" className={`${styles.quickCard} card-hover`} bordered={false}>
            <div className={styles.quickGrid}>
              {quickActions.map((action, index) => (
                <div
                  key={action.title}
                  className={styles.quickItem}
                  onClick={() => navigate(action.path)}
                  role="button"
                  tabIndex={0}
                  onKeyDown={(e) => e.key === 'Enter' && navigate(action.path)}
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <span className={styles.quickIcon} style={{ background: `${action.color}15`, color: action.color }}>
                    {action.icon}
                  </span>
                  <div>
                    <div className={styles.quickTitle}>{action.title}</div>
                    <div className={styles.quickDesc}>{action.desc}</div>
                  </div>
                  <PlusOutlined className={styles.quickPlus} style={{ color: action.color }} />
                </div>
              ))}
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
}