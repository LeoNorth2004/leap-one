/** 进度徽章组件 */

import { Progress } from 'antd';
import type { ProgressProps } from 'antd';

interface ProgressBadgeProps {
  /** 进度百分比 (0-100) */
  percent: number;
  /** 进度条尺寸 */
  size?: ProgressProps['size'];
  /** 是否显示文字 */
  showInfo?: boolean;
  /** 状态颜色 */
  statusColor?: string;
}

export default function ProgressBadge({
  percent,
  size = 'small',
  showInfo = true,
  statusColor,
}: ProgressBadgeProps) {
  let status: ProgressProps['status'] = 'active';
  if (percent >= 100) status = 'success';
  else if (percent >= 80) status = 'normal';

  return (
    <Progress
      percent={Math.round(percent)}
      size={size}
      status={status}
      showInfo={showInfo}
      strokeColor={statusColor || undefined}
      style={{ margin: 0 }}
    />
  );
}
