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

import { useMemo } from 'react';
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
} from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import PageHeader from '@/components/Common/PageHeader';
import { PRIORITY_MAP, TASK_STATUS_MAP } from '@/utils/constants';
import type { TodoItem, CalendarEvent, ActivityItem } from '@/types/common';
import styles from './index.module.less';

const { Text, Title } = Typography;

// ─── Mock 数据函数（替换为真实 API 即可）────────────────────

/** 获取统计数据 */
async function fetchDashboardStats(): Promise<{
  taskCount: number;
  requirementCount: number;
  bugCount: number;
  issueCount: number;
}> {
  // TODO: 替换为真实 API 调用
  // return get<DashboardStats>('/dashboard/stats');
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        taskCount: 12,
        requirementCount: 5,
        bugCount: 3,
        issueCount: 7,
      });
    }, 600);
  });
}

/** 获取最近动态 */
async function fetchRecentActivities(): Promise<ActivityItem[]> {
  // TODO: 替换为真实 API 调用
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

/** 获取待办事项 Top 5 */
async function fetchTodoItems(): Promise<TodoItem[]> {
  // TODO: 替换为真实 API 调用
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

/** 获取日历日程 */
async function fetchCalendarEvents(): Promise<CalendarEvent[]> {
  // TODO: 替换为真实 API 调用
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        { id: 1, title: 'Sprint 迭代评审会', date: '2026-06-09', type: 'meeting', color: '#1677ff' },
        { id: 2, title: 'V1.2 版本发布', date: '2026-06-12', type: 'deadline', color: '#ff4d4f' },
        { id: 3, title: '需求评审会议', date: '2026-06-10', type: 'review', color: '#faad14' },
        { id: 4, title: '技术方案讨论', date: '2026-06-11', type: 'meeting', color: '#1677ff' },
        { id: 5, title: '代码审查截止', date: '2026-06-13', type: 'deadline', color: '#ff4d4f' },
      ]);
    }, 300);
  });
}

// ─── 组件 ──────────────────────────────────────────────────────

export default function Dashboard() {
  const navigate = useNavigate();

  // TanStack Query 数据获取
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

  // ─── 统计卡片配置 ────────────────────────────────────────────
  const statCards = useMemo(
    () =>
      statsData
        ? [
            {
              title: '进行中任务',
              value: statsData.taskCount,
              icon: <CarryOutOutlined />,
              color: '#1677ff',
              trend: 8,
              suffix: '个',
            },
            {
              title: '我的需求',
              value: statsData.requirementCount,
              icon: <FileTextOutlined />,
              color: '#13c2c2',
              trend: -2,
              suffix: '个',
            },
            {
              title: '待修复 Bug',
              value: statsData.bugCount,
              icon: <BugOutlined />,
              color: '#ff4d4f',
              trend: 1,
              suffix: '个',
            },
            {
              title: '我的工单',
              value: statsData.issueCount,
              icon: <CustomerServiceOutlined />,
              color: '#faad14',
              trend: 3,
              suffix: '个',
            },
          ]
        : [],
    [statsData]
  );

  // ─── 动态时间线项 ────────────────────────────────────────────
  const timelineItems = useMemo(
    () =>
      activities?.map((item) => ({
        color: item.targetType === 'bug'
          ? 'red'
          : item.targetType === 'requirement'
            ? 'blue'
            : item.targetType === 'task'
              ? 'green'
              : 'orange',
        children: (
          <div>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {item.timestamp}
            </Text>
            <br />
            <Text>
              {item.user} {item.action} {item.target}
            </Text>
          </div>
        ),
      })) ?? [],
    [activities]
  );

  // ─── 日历单元格渲染 ──────────────────────────────────────────
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

  // ─── 待办优先级图标 ──────────────────────────────────────────
  const getPriorityTag = (priority: string) => {
    const info = PRIORITY_MAP[priority];
    return info ? <Tag color={info.color}>{info.label}</Tag> : null;
  };

  /** 获取待办状态图标 */
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'doing':
      case 'active':
        return <ClockCircleOutlined style={{ color: '#1677ff' }} />;
      case 'done':
      case 'resolved':
        return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
      case 'wait':
        return <WarningOutlined style={{ color: '#faad14' }} />;
      default:
        return <ClockCircleOutlined style={{ color: '#86909c' }} />;
    }
  };

  // ─── 快捷操作入口 ────────────────────────────────────────────
  const quickActions = [
    { title: '创建任务', desc: '新建一个任务', icon: <CarryOutOutlined />, path: '/task/list' },
    { title: '提交 Bug', desc: '报告一个问题', icon: <BugOutlined />, path: '/quality/bug' },
    { title: '创建工单', desc: '提交服务请求', icon: <CustomerServiceOutlined />, path: '/issue/list' },
    { title: '写文档', desc: '编辑项目文档', icon: <FileTextOutlined />, path: '/document/list' },
  ];

  // 全局加载状态
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

      {/* ═══ 统计卡片区域 ═══ */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {statCards.map((stat) => (
          <Col xs={24} sm={12} lg={6} key={stat.title}>
            <Card className={`${styles.statCard} card-hover`}>
              <Statistic
                title={stat.title}
                value={stat.value}
                suffix={stat.suffix}
                prefix={
                  <span className={styles.statIcon} style={{ backgroundColor: stat.color }}>
                    {stat.icon}
                  </span>
                }
                valueStyle={{ fontSize: 28, fontWeight: 700 }}
              />
              <div className={styles.trend}>
                {stat.trend > 0 ? (
                  <span style={{ color: '#52c41a' }}>
                    <ArrowUpOutlined /> +{Math.abs(stat.trend)}
                  </span>
                ) : (
                  <span style={{ color: '#ff4d4f' }}>
                    <ArrowDownOutlined /> {stat.trend}
                  </span>
                )}
                <Text type="secondary" style={{ marginLeft: 8 }}>
                  较昨日
                </Text>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      {/* ═══ 主要内容区：动态 + 待办 + 快捷入口 + 日历 ═══ */}
      <Row gutter={[16, 16]}>
        {/* 最近动态时间线 */}
        <Col xs={24} lg={12}>
          <Card
            title="最近动态"
            className={`${styles.activityCard} card-hover`}
            loading={activitiesLoading}
          >
            {timelineItems.length > 0 ? (
              <Timeline items={timelineItems} />
            ) : (
              <Empty description="暂无动态" image={Empty.PRESENTED_IMAGE_SIMPLE} />
            )}
          </Card>
        </Col>

        {/* 待办事项 Top 5 */}
        <Col xs={24} lg={12}>
          <Card
            title="我的待办"
            className={`${styles.todoCard} card-hover`}
            loading={todosLoading}
            extra={<a onClick={() => navigate('/task/list')}>查看全部</a>}
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
                            <Tag>{todo.projectName}</Tag>
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

        {/* 日历视图 */}
        <Col xs={24} lg={12}>
          <Card
            title="日程安排"
            className={`${styles.calendarCard} card-hover`}
            loading={calendarLoading}
          >
            <Calendar
              fullscreen={false}
              cellRender={(date) => cellRender(date as unknown as Dayjs)}
              style={{ borderRadius: 8 }}
            />
          </Card>
        </Col>

        {/* 快捷入口 */}
        <Col xs={24} lg={12}>
          <Card title="快捷入口" className={`${styles.quickCard} card-hover`}>
            <div className={styles.quickGrid}>
              {quickActions.map((action) => (
                <div
                  key={action.title}
                  className={styles.quickItem}
                  onClick={() => navigate(action.path)}
                  role="button"
                  tabIndex={0}
                  onKeyDown={(e) => e.key === 'Enter' && navigate(action.path)}
                >
                  <span className={styles.quickIcon}>{action.icon}</span>
                  <div>
                    <div className={styles.quickTitle}>{action.title}</div>
                    <div className={styles.quickDesc}>{action.desc}</div>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
