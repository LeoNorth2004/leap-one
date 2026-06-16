/** 项目看板页面（Kanban Board） */

import { useState } from 'react';
import { Card, Tag, Badge, Avatar, Space, Empty } from 'antd';
import PageHeader from '@/components/Common/PageHeader';
import PriorityBadge from '@/components/Business/PriorityBadge';
import UserAvatar from '@/components/Business/UserAvatar';

interface KanbanCard {
  id: string;
  title: string;
  priority: string;
  assignee?: string;
  assigneeAvatar?: string;
  tags?: string[];
}

interface KanbanColumn {
  id: string;
  title: string;
  cards: KanbanCard[];
}

/** 模拟看板数据 */
const initialColumns: KanbanColumn[] = [
  {
    id: 'todo',
    title: '待处理',
    cards: [
      { id: 't1', title: '设计数据库表结构', priority: 'P0', assignee: '李四', tags: ['设计'] },
      { id: 't2', title: '编写API接口文档', priority: 'P1', assignee: '王五', tags: ['文档'] },
      { id: 't3', title: '配置CI/CD流水线', priority: 'P2', tags: ['运维'] },
    ],
  },
  {
    id: 'doing',
    title: '进行中',
    cards: [
      { id: 't4', title: '实现用户认证模块', priority: 'P0', assignee: '张三', tags: ['开发'] },
      { id: 't5', title: '优化首页加载性能', priority: 'P1', assignee: '赵六', tags: ['性能优化'] },
    ],
  },
  {
    id: 'review',
    title: '评审中',
    cards: [
      { id: 't6', title: '代码审查：权限模块', priority: 'P1', assignee: '钱七', tags: ['审查'] },
    ],
  },
  {
    id: 'done',
    title: '已完成',
    cards: [
      { id: 't7', title: '项目初始化和脚手架搭建', priority: 'P2', assignee: '张三', tags: ['基础设施'] },
      { id: 't8', title: 'Git工作流规范制定', priority: 'P3', tags: ['流程'] },
    ],
  },
];

// 简化的看板实现（不依赖外部拖拽库）
export default function KanbanBoard() {
  const [columns] = useState<KanbanColumn[]>(initialColumns);

  return (
    <div>
      <PageHeader
        title="项目看板"
        subtitle="以可视化方式管理任务的流转"
      />

      <div className="kanban-board" style={{ display: 'flex', gap: 16, overflowX: 'auto', padding: '4px 0' }}>
        {columns.map((col) => (
          <div key={col.id} className="kanban-column" style={{
            minWidth: 280,
            flex: 1,
            background: '#f5f5f5',
            borderRadius: 12,
            padding: 12,
            display: 'flex',
            flexDirection: 'column',
          }}>
            {/* 列头 */}
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: 12,
              padding: '0 4px',
            }}>
              <strong>{col.title}</strong>
              <Badge count={col.cards.length} showZero style={{ backgroundColor: '#1677ff' }} />
            </div>

            {/* 卡片列表 */}
            <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 8 }}>
              {col.cards.length === 0 ? (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无任务" style={{ margin: 'auto' }} />
              ) : (
                col.cards.map((card) => (
                  <Card
                    key={card.id}
                    size="small"
                    hoverable
                    className="kanban-card"
                    style={{ borderRadius: 8, cursor: 'grab' }}
                  >
                    <div style={{ marginBottom: 8 }}>
                      <PriorityBadge priority={card.priority as 'P0' | 'P1' | 'P2' | 'P3'} />
                      <span style={{ fontWeight: 500, marginLeft: 8 }}>{card.title}</span>
                    </div>
                    {(card.tags || []).map((tag) => (
                      <Tag key={tag} style={{ fontSize: 11, marginBottom: 4 }}>{tag}</Tag>
                    ))}
                    {card.assignee && (
                      <div style={{ marginTop: 8, display: 'flex', justifyContent: 'flex-end' }}>
                        <UserAvatar name={card.assignee} size={24} />
                      </div>
                    )}
                  </Card>
                ))
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
