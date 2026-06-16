/** 状态标签组件 */

import { Tag } from 'antd';

interface StatusTagProps {
  /** 状态值 */
  status: string;
  /** 自定义状态映射，默认使用内置映射 */
  statusMap?: Record<string, { label: string; color: string }>;
}

/** 内置状态映射 */
const defaultStatusMaps: Record<string, Record<string, { label: string; color: string }>> = {
  task: {
    wait: { label: '待处理', color: 'default' },
    doing: { label: '进行中', color: 'processing' },
    done: { label: '已完成', color: 'success' },
    pause: { label: '已暂停', color: 'warning' },
    cancel: { label: '已取消', color: 'error' },
    closed: { label: '已关闭', color: 'default' },
  },
  requirement: {
    draft: { label: '草稿', color: 'default' },
    reviewing: { label: '评审中', color: 'processing' },
    active: { label: '激活', color: 'blue' },
    developing: { label: '开发中', color: 'cyan' },
    testing: { label: '测试中', color: 'orange' },
    completed: { label: '已完成', color: 'success' },
    closed: { label: '已关闭', color: 'default' },
    rejected: { label: '已拒绝', color: 'error' },
  },
  project: {
    planning: { label: '规划中', color: 'default' },
    active: { label: '进行中', color: 'processing' },
    paused: { label: '已暂停', color: 'warning' },
    completed: { label: '已完成', color: 'success' },
    archived: { label: '已归档', color: 'default' },
  },
};

export default function StatusTag({ status, statusMap }: StatusTagProps) {
  const map = statusMap || defaultStatusMaps.task;
  const config = map[status] || { label: status, color: 'default' };

  return <Tag color={config.color}>{config.label}</Tag>;
}
