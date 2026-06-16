/** 优先级徽章组件 */

import { Tag } from 'antd';

type PriorityValue = 'P0' | 'P1' | 'P2' | 'P3' | 'urgent' | 'high' | 'medium' | 'low';

interface PriorityBadgeProps {
  /** 优先级值 */
  priority: PriorityValue;
}

const priorityConfig: Record<PriorityValue, { label: string; color: string }> = {
  P0: { label: 'P0-紧急', color: 'red' },
  P1: { label: 'P1-高', color: 'orange' },
  P2: { label: 'P2-中', color: 'blue' },
  P3: { label: 'P3-低', color: 'default' },
  urgent: { label: '紧急', color: 'red' },
  high: { label: '高', color: 'orange' },
  medium: { label: '中', color: 'blue' },
  low: { label: '低', color: 'default' },
};

export default function PriorityBadge({ priority }: PriorityBadgeProps) {
  const config = priorityConfig[priority] || { label: priority, color: 'default' };
  return <Tag color={config.color}>{config.label}</Tag>;
}
