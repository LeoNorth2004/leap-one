/**
 * 进度徽章组件
 *
 * 基于 Ant Design Progress 的封装，根据百分比自动设置状态：
 * - >= 100% → success（绿色）
 * - >= 80%  → normal（蓝色）
 * - < 80%   → active（动画进行中）
 */

import { Progress } from 'antd';
import type { ProgressProps } from 'antd';

// ── 类型定义 ─────────────────────────────────────────────────

interface ProgressBadgeProps {
  /** 进度百分比 (0-100) */
  percent: number;
  /** 进度条尺寸 */
  size?: ProgressProps['size'];
  /** 是否显示百分比文字 */
  showInfo?: boolean;
  /** 自定义进度条颜色 */
  statusColor?: string;
}

// ── 状态判定阈值 ─────────────────────────────────────────────

const SUCCESS_THRESHOLD = 100;
const NORMAL_THRESHOLD = 80;

// ── 工具函数 ─────────────────────────────────────────────────

/** 根据百分比计算进度状态 */
const resolveStatus = (percent: number): ProgressProps['status'] => {
  if (percent >= SUCCESS_THRESHOLD) return 'success';
  if (percent >= NORMAL_THRESHOLD) return 'normal';
  return 'active';
};

// ── 组件实现 ─────────────────────────────────────────────────

const ProgressBadge = ({
  percent,
  size = 'small',
  showInfo = true,
  statusColor,
}: ProgressBadgeProps) => {
  const status = resolveStatus(percent);

  return (
    <Progress
      percent={Math.round(percent)}
      size={size}
      status={status}
      showInfo={showInfo}
      strokeColor={statusColor ?? undefined}
      style={{ margin: 0 }}
    />
  );
};

export default ProgressBadge;
