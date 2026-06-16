/** 需求详情页面 */

import { useParams, useNavigate } from 'react-router-dom';
import { Descriptions, Card, Tag, Tabs, Timeline, Divider, Button, Space, Typography } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';
import StatusTag from '@/components/Common/StatusTag';
import PriorityBadge from '@/components/Business/PriorityBadge';
import UserAvatar from '@/components/Business/UserAvatar';
import type { Requirement } from '@/types/requirement';

const { Paragraph } = Typography;

export default function RequirementDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  /** 模拟需求详情 - TODO: 替换为API */
  const requirement: Requirement = {
    id: Number(id) || 1,
    title: '用户权限RBAC模型实现',
    code: 'REQ-001',
    description: '实现基于角色(Role)的访问控制模型，支持角色继承、权限矩阵配置。包括：\n\n1. 角色(Role)定义和管理\n2. 权限(Permission)粒度到操作级别\n3. 用户-角色多对多关联\n4. 角色继承关系\n5. 前端菜单权限动态渲染\n6. API级别的鉴权拦截器',
    status: 'developing',
    priority: 'P0',
    source: 'internal',
    productId: 1,
    productName: '企业管理系统',
    moduleId: 1,
    moduleName: '权限管理',
    storyPoints: 8,
    assigneeId: 1,
    assigneeName: '张三',
    reviewerId: 2,
    reviewerName: '李四',
    planRelease: 'v2.0.0',
    createdAt: '2025-06-01T08:00:00Z',
    updatedAt: '2026-06-05T14:30:00Z',
    createdBy: '产品经理',
  };

  /** 模拟操作日志 */
  const activities = [
    { time: '2026-06-05 14:30', content: '张三 更新了需求描述', color: 'blue' },
    { time: '2026-06-03 10:15', content: '李四 通过了需求评审', color: 'green' },
    { time: '2026-06-01 09:00', content: '产品经理 将需求指派给张三', color: 'orange' },
    { time: '2026-05-28 16:45', content: '产品经理 提交了需求评审', color: 'purple' },
    { time: '2025-06-01 10:00', content: '产品经理 创建了此需求', color: 'green' },
  ];

  return (
    <div>
      <PageHeader
        title={`${requirement.code}: ${requirement.title}`}
        extra={
          <Space>
            <Button onClick={() => navigate(-1)}><ArrowLeftOutlined /> 返回</Button>
            <Button>编辑</Button>
            <Button type="primary">变更状态</Button>
          </Space>
        }
      />

      <Card>
        <Tabs
          defaultActiveKey="detail"
          items={[
            {
              key: 'detail',
              label: '基本信息',
              children: (
                <>
                  <Descriptions column={2} bordered size="middle">
                    <Descriptions.Item label="需求编号">{requirement.code}</Descriptions.Item>
                    <Descriptions.Item label="优先级"><PriorityBadge priority={requirement.priority} /></Descriptions.Item>
                    <Descriptions.Item label="状态"><StatusTag status={requirement.status} /></Descriptions.Item>
                    <Descriptions.Item label="来源">
                      <Tag>{requirement.source === 'internal' ? '内部' : requirement.source === 'customer' ? '客户' : requirement.source === 'market' ? '市场' : '竞品'}</Tag>
                    </Descriptions.Item>
                    <Descriptions.Item label="所属产品">{requirement.productName}</Descriptions.Item>
                    <Descriptions.Item label="所属模块">{requirement.moduleName || '-'}</Descriptions.Item>
                    <Descriptions.Item label="计划故事点">{requirement.storyPoints || '-'}</Descriptions.Item>
                    <Descriptions.Item label="计划版本">{requirement.planRelease || '-'}</Descriptions.Item>
                    <Descriptions.Item label="指派人">
                      {requirement.assigneeName && <Space><UserAvatar name={requirement.assigneeName} size={22} />{requirement.assigneeName}</Space>}
                    </Descriptions.Item>
                    <Descriptions.Item label="评审人">
                      {requirement.reviewerName && <Space><UserAvatar name={requirement.reviewerName} size={22} />{requirement.reviewerName}</Space>}
                    </Descriptions.Item>
                    <Descriptions.Item label="创建人">{requirement.createdBy}</Descriptions.Item>
                    <Descriptions.Item label="创建时间">{new Date(requirement.createdAt).toLocaleString('zh-CN')}</Descriptions.Item>
                    <Descriptions.Item label="更新时间">{new Date(requirement.updatedAt).toLocaleString('zh-CN')}</Descriptions.Item>
                  </Descriptions>

                  <Divider orientation="left">需求描述</Divider>
                  <Paragraph style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>{requirement.description}</Paragraph>
                </>
              ),
            },
            {
              key: 'tasks',
              label: '关联任务',
              children: <p style={{ color: '#999' }}>暂无关联任务。</p>,
            },
            {
              key: 'bugs',
              label: '关联Bug',
              children: <p style={{ color: '#999' }}>暂无关联Bug。</p>,
            },
            {
              key: 'activities',
              label: '操作日志',
              children: (
                <Timeline
                  items={activities.map((item) => ({
                    children: (
                      <div>
                        <span style={{ color: '#666', fontSize: 12, marginRight: 8 }}>{item.time}</span>
                        {item.content}
                      </div>
                    ),
                    color: item.color,
                  }))}
                />
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
}
