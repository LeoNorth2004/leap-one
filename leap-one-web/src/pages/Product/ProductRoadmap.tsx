/** 产品路线图页面 */

import { Timeline, Card, Tag, Badge, Space } from 'antd';
import PageHeader from '@/components/Common/PageHeader';

interface VersionItem {
  id: number;
  version: string;
  name: string;
  planDate: string;
  releaseDate?: string;
  status: 'planning' | 'developing' | 'released' | 'delayed';
  requirements: Array<{ id: number; title: string; priority: string; status: string }>;
}

/** 模拟路线图数据 */
const mockRoadmap: VersionItem[] = [
  {
    id: 1, version: 'v2.0.0', name: '智能分析模块',
    planDate: '2026-07-01', releaseDate: undefined,
    status: 'developing',
    requirements: [
      { id: 1, title: 'AI数据分析引擎', priority: 'P0', status: 'developing' },
      { id: 2, title: '自定义报表生成', priority: 'P1', status: 'active' },
      { id: 3, title: '数据导出优化', priority: 'P2', status: 'draft' },
    ],
  },
  {
    id: 2, version: 'v1.5.0', name: '移动端适配',
    planDate: '2026-08-15', releaseDate: undefined,
    status: 'planning',
    requirements: [
      { id: 4, title: '响应式布局重构', priority: 'P0', status: 'draft' },
      { id: 5, title: '移动端导航优化', priority: 'P1', status: 'draft' },
    ],
  },
  {
    id: 3, version: 'v1.4.0', name: '协作功能增强',
    planDate: '2026-05-01', releaseDate: '2026-05-28',
    status: 'released',
    requirements: [
      { id: 6, title: '实时评论通知', priority: 'P1', status: 'completed' },
      { id: 7, title: '文件共享功能', priority: 'P2', status: 'completed' },
    ],
  },
  {
    id: 4, version: 'v1.3.0', name: '基础权限体系',
    planDate: '2026-03-01', releaseDate: '2026-03-25',
    status: 'released',
    requirements: [
      { id: 8, title: 'RBAC权限模型', priority: 'P0', status: 'completed' },
    ],
  },
];

const statusConfig: Record<string, { color: string; label: string }> = {
  planning: { color: 'default', label: '规划中' },
  developing: { color: 'processing', label: '开发中' },
  released: { color: 'success', label: '已发布' },
  delayed: { color: 'error', label: '延期' },
};

const priorityColors: Record<string, string> = {
  P0: 'red', P1: 'orange', P2: 'blue', P3: 'default',
};

export default function ProductRoadmap() {
  return (
    <div>
      <PageHeader title="产品路线图" subtitle="查看版本规划与里程碑时间线" />

      <Timeline
        mode="left"
        items={mockRoadmap.map((version) => {
          const status = statusConfig[version.status] || { color: 'default', label: version.status };
          return {
            color: status.color,
            children: (
              <Card
                size="small"
                title={
                  <Space>
                    <strong>{version.name}</strong>
                    <Tag>{version.version}</Tag>
                    <Badge
                      status={status.color as 'success' | 'processing' | 'error' | 'default'}
                      text={status.label}
                    />
                  </Space>
                }
                style={{ marginBottom: 8 }}
              >
                <div style={{ marginBottom: 8, color: '#666', fontSize: 13 }}>
                  计划日期：{version.planDate}
                  {version.releaseDate && ` → 发布日期：${version.releaseDate}`}
                </div>
                <div>
                  {version.requirements.map((req) => (
                    <Tag key={req.id} color={priorityColors[req.priority]} style={{ marginBottom: 4 }}>
                      [{req.priority}] {req.title}
                    </Tag>
                  ))}
                </div>
              </Card>
            ),
          };
        })}
      />
    </div>
  );
}
