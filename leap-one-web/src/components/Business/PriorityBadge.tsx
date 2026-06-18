/**
 * 优先级徽章组件
 *
 * 根据优先级值渲染对应颜色和标签的 Tag
 * 支持 P0-P3 和 urgent/high/medium/low 两套编码
 */

import { Tag } from 'antd';

// ── 类型定义 ─────────────────────────────────────────────────

type PriorityValue = 'P0' | 'P1' | 'P2' | 'P3' | 'urgent' | 'high' | 'medium' | 'low';

interface PriorityBadgeProps {
  /** 优先级值 */
  priority: PriorityValue;
}

// ── 优先级配置映射表 ─────────────────────────────────────────

const PRIORITY_CONFIG: Record<PriorityValue, { label: string; color: string }> = Object.freeze({
  P0:     { label: 'P0-紧急', color: 'red' },
  P1:     { label: 'P1-高',   color: 'orange' },
  P2:     { label: 'P2-中',   color: 'blue' },
  P3:     { label: 'P3-低',   color: 'default' },
  urgent: { label: '紧急',    color: 'red' },
  high:   { label: '高',      color: 'orange' },
  medium: { label: '中',      color: 'blue' },
  low:    { label: '低',      color: 'default' },
});

/** 未匹配时的回退配置 */
const FALLBACK_PRIORITY = Object.freeze({ label: '', color: 'default' });

// ── 组件实现 ─────────────────────────────────────────────────

const PriorityBadge = ({ priority }: PriorityBadgeProps) => {
  const config = PRIORITY_CONFIG[priority] ?? { ...FALLBACK_PRIORITY, label: priority };

  return <Tag color={config.color}>{config.label}</Tag>;
};

export default PriorityBadge;
