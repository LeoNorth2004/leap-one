/**
 * 状态标签组件
 *
 * 根据状态值自动渲染对应颜色和文字的 Tag
 * 支持自定义状态映射，内置任务/需求/项目三种映射
 */

import { Tag } from 'antd';
import type { ReactNode } from 'react';

// ── 类型定义 ─────────────────────────────────────────────────

interface StatusTagProps {
  /** 状态值 */
  status: string;
  /** 自定义状态映射（默认使用任务状态映射） */
  statusMap?: Record<string, { label: string; color: string }>;
}

// ── 内置状态映射表 ─────────────────────────────────────────────

const STATUS_MAPS = Object.freeze({
  task: {
    wait:     { label: '待处理', color: 'default' },
    doing:    { label: '进行中', color: 'processing' },
    done:     { label: '已完成', color: 'success' },
    pause:    { label: '已暂停', color: 'warning' },
    cancel:   { label: '已取消', color: 'error' },
    closed:   { label: '已关闭', color: 'default' },
  },
  requirement: {
    draft:       { label: '草稿',   color: 'default' },
    reviewing:   { label: '评审中', color: 'processing' },
    active:      { label: '激活',   color: 'blue' },
    developing:  { label: '开发中', color: 'cyan' },
    testing:     { label: '测试中', color: 'orange' },
    completed:   { label: '已完成', color: 'success' },
    closed:      { label: '已关闭', color: 'default' },
    rejected:    { label: '已拒绝', color: 'error' },
  },
  project: {
    planning:  { label: '规划中', color: 'default' },
    active:    { label: '进行中', color: 'processing' },
    paused:    { label: '已暂停', color: 'warning' },
    completed: { label: '已完成', color: 'success' },
    archived:  { label: '已归档', color: 'default' },
  },
}) as Record<string, Record<string, { label: string; color: string }>>;

/** 未匹配状态的默认配置 */
const FALLBACK_CONFIG = Object.freeze({ label: '', color: 'default' });

// ── 组件实现 ─────────────────────────────────────────────────

const StatusTag = ({ status, statusMap }: StatusTagProps): ReactNode => {
  const activeMap = statusMap ?? STATUS_MAPS.task;
  const config = (activeMap as Record<string, { label: string; color: string }>)[status]
    ?? { ...FALLBACK_CONFIG, label: status };

  return <Tag color={config.color}>{config.label}</Tag>;
};

export default StatusTag;
