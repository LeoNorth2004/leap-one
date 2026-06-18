// 格式化工具集

import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);

const DEFAULT_DATE_FORMAT = 'YYYY-MM-DD HH:mm:ss';
const FILE_SIZE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB'];
const BASE_UNIT = 1024;

/** 日期格式化 */
export const formatDate = (date: string | Date, format: string = DEFAULT_DATE_FORMAT): string => {
  if (!date) {
    return '';
  }
  return dayjs(date).format(format);
};

/** 相对时间格式化（如：3分钟前） */
export const formatRelativeTime = (date: string | Date): string => {
  if (!date) {
    return '';
  }
  return dayjs(date).fromNow();
};

/** 文件大小格式化 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) {
    return '0 B';
  }

  const unitIndex = Math.floor(Math.log(bytes) / Math.log(BASE_UNIT));
  const size = parseFloat((bytes / Math.pow(BASE_UNIT, unitIndex)).toFixed(2));
  const unit = FILE_SIZE_UNITS[unitIndex];

  return `${size} ${unit}`;
};

/** 数字千分位格式化 */
export const formatNumber = (num: number): string => num.toLocaleString('zh-CN');

/** 百分比格式化 */
export const formatPercent = (value: number, decimals: number = 1): string => {
  const percent = (value * 100).toFixed(decimals);
  return `${percent}%`;
};

/** 文本截断（超出长度显示省略号） */
export const truncateText = (text: string, maxLength: number): string => {
  if (!text || text.length <= maxLength) {
    return text;
  }
  return `${text.slice(0, maxLength)}...`;
};
