/** 全局看板视图页面 */

import { useState } from 'react';
import { Card, Tag, Badge, Select, Space } from 'antd';
import PageHeader from '@/components/Common/PageHeader';
import PriorityBadge from '@/components/Business/PriorityBadge';
import UserAvatar from '@/components/Business/UserAvatar';

interface KanbanItem {
  id: string;
  title: string;
  type: 'requirement' | 'task' | 'bug';
  priority: string;
  status: string;
  project: string;
  assignee?: string;
}

const initialItems: KanbanItem[] = [
  { id: 'k1', title: 'RBAC权限模型实现', type: 'requirement', priority: 'P0', status: 'developing', project: '企业管理系统', assignee: '张三' },
  { id: 'k2', title: '用户认证中间件开发', type: 'task', priority: 'urgent', status: 'doing', project: '企业管理系统', assignee: '张三' },
  { id: 'k3', title: '首页白屏偶现问题', type: 'bug', priority: 'fatal', status: 'open', project: '企业管理系统', assignee: '李四' },
  { id: 'k4', title: '数据导出Excel功能', type: 'requirement', priority: 'P1', status: 'reviewing', project: '企业管理系统' },
  { id: 'k5', title: '移动端适配优化', type: 'requirement', priority: 'P1', status: 'active', project: '移动办公App', assignee: '王五' },
  { id: 'k6', title: '编写权限API单元测试', type: 'task', priority: 'high', status: 'wait', project: '企业管理系统', assignee: '赵六' },
];

const columnConfig = [
  { key: 'todo', title: '待处理', color: '#d9d9d9' },
  { key: 'doing', title: '进行中', color: '#bae7ff' },
  { key: 'review', title: '评审/测试', color: '#fff1b8' },
  { key: 'done', title: '已完成', color: '#b7eb8f' },
];

const typeColors: Record<string, string> = {
  requirement: 'blue',
  task: 'cyan',
  bug: 'red',
};

export default function KanbanView() {
  const [projectFilter, setProjectFilter] = useState<string | undefined>();

  /** 按列分组 */
  const groupedByColumn = columnConfig.map((col) => ({
    ...col,
    items: initialItems.filter((item) => {
      if (col.key === 'todo') return ['draft', 'wait', 'open', 'pending'].includes(item.status);
      if (col.key === 'doing') return ['doing', 'developing', 'active'].includes(item.status);
      if (col.key === 'review') return ['reviewing', 'testing'].includes(item.status);
      if (col.key === 'done') return ['done', 'completed', 'closed', 'resolved'].includes(item.status);
      return false;
    }),
  }));

  return (
    <div>
      <PageHeader
        title="全局看板"
        subtitle="跨项目、跨类型的统一任务视图"
        extra={
          <Select placeholder="筛选项目" allowClear style={{ width: 180 }} value={projectFilter} onChange={setProjectFilter}
            options={[{ label: '企业管理系统', value: '1' }, { label: '移动办公App', value: '2' }]} />
        }
      />

      <div style={{ display: 'flex', gap: 16, overflowX: 'auto', padding: '4px 0' }}>
        {groupedByColumn.map((col) => (
          <div key={col.key} style={{
            minWidth: 300, flex: 1, background: col.color, borderRadius: 12, padding: 12,
            display: 'flex', flexDirection: 'column',
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
              <strong>{col.title}</strong>
              <Badge count={col.items.length} style={{ backgroundColor: '#666' }} />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {col.items.length === 0 ? (
                <div style={{ textAlign: 'center', padding: 20, color: '#999' }}>暂无</div>
              ) : (
                col.items.map((item) => (
                  <Card key={item.id} size="small" hoverable style={{ borderRadius: 8 }}>
                    <Space direction="vertical" size={4} style={{ width: '100%' }}>
                      <div>
                        <Tag color={typeColors[item.type]}>{item.type === 'requirement' ? '需求' : item.type === 'task' ? '任务' : 'Bug'}</Tag>
                        <PriorityBadge priority={item.priority as 'P0' | 'P1' | 'P2' | 'P3'} />
                      </div>
                    <span style={{ fontWeight: 500 }}>{item.title}</span>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <span style={{ fontSize: 12, color: '#999' }}>{item.project}</span>
                      {item.assignee && <UserAvatar name={item.assignee} size={22} />}
                    </div>
                  </Space>
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
