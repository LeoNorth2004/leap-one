/** 格式化工具函数 */

import dayjs from 'dayjs';

/**
 * 格式化日期
 * @param date - 日期字符串或Date对象
 * @param format - 格式模板，默认 YYYY-MM-DD HH:mm:ss
 */
export function formatDate(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!date) return '';
  return dayjs(date).format(format);
}

/**
 * 格式化相对时间（如：3分钟前）
 */
export function formatRelativeTime(date: string | Date): string {
  if (!date) return '';
  return dayjs(date).fromNow();
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const k = 1024;
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${units[i]}`;
}

/**
 * 格式化数字（千分位）
 */
export function formatNumber(num: number): string {
  return num.toLocaleString('zh-CN');
}

/**
 * 格式化百分比
 */
export function formatPercent(value: number, decimals = 1): string {
  return `${(value * 100).toFixed(decimals)}%`;
}

/**
 * 截断文本
 */
export function truncateText(text: string, maxLength: number): string {
  if (!text || text.length <= maxLength) return text;
  return text.slice(0, maxLength) + '...';
}
