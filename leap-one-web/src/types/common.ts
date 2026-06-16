/** 通用类型定义 */

// ─── 基础通用类型 ─────────────────────────────────────────────────

/** 选项类型（用于下拉选择等） */
export interface OptionItem {
  label: string;
  value: string | number;
  disabled?: boolean;
}

/** 树形节点 */
export interface TreeNode<T = unknown> {
  key: string;
  title: string;
  children?: TreeNode<T>[];
  data?: T;
}

/** 时间范围 */
export interface DateRange {
  start: string;
  end: string;
}

// ─── 统计与展示类型 ───────────────────────────────────────────────

/** 统计卡片数据 */
export interface StatCard {
  title: string;
  value: number;
  icon: string;
  color: string;
  trend?: number;
  trendLabel?: string;
}

/** 动态/操作日志条目 */
export interface ActivityItem {
  id: number;
  action: string;
  target: string;
  targetType: string;
  user: string;
  userAvatar: string;
  timestamp: string;
  detail?: string;
}

/** 待办事项条目 */
export interface TodoItem {
  id: number;
  title: string;
  type: 'task' | 'requirement' | 'bug' | 'issue';
  priority: string;
  status: string;
  dueDate?: string;
  projectName?: string;
}

/** 日历日程项 */
export interface CalendarEvent {
  id: number;
  title: string;
  date: string;
  type: 'meeting' | 'deadline' | 'review' | 'other';
  color?: string;
}

// ─── 文件相关 ─────────────────────────────────────────────────────

/** 文件上传结果 */
export interface UploadFileResult {
  uid: string;
  name: string;
  url: string;
  size: number;
  type: string;
}

// ─── 系统配置类型 ─────────────────────────────────────────────────

/** 主题模式 */
export type ThemeMode = 'light' | 'dark';

/** 语言 */
export type LocaleType = 'zh-CN' | 'en-US';

// ─── 枚举类型 ─────────────────────────────────────────────────────

/** 任务优先级 */
export type TaskPriority = 'P0' | 'P1' | 'P2' | 'P3';

/** 任务状态 */
export type TaskStatus = 'wait' | 'doing' | 'done' | 'pause' | 'cancel' | 'closed';

/** Bug 严重程度 */
export type BugSeverity = 'fatal' | 'serious' | 'normal' | 'slight' | 'suggest';

/** Bug 状态 */
export type BugStatus = 'active' | 'resolved' | 'closed' | 'postponed';

/** 需求类型 */
export type RequirementType = 'user' | 'business' | 'technical' | 'defect';

/** 需求状态 */
export type RequirementStatus =
  | 'draft'
  | 'reviewing'
  | 'active'
  | 'developing'
  | 'testing'
  | 'completed'
  | 'closed'
  | 'rejected';

/** 项目类型 */
export type ProjectType = 'scrum' | 'waterfall' | 'kanban' | 'hybrid';

/** 项目状态 */
export type ProjectStatus = 'planning' | 'active' | 'paused' | 'completed' | 'archived';

/** 用户角色 */
export type UserRole = 'admin' | 'project_manager' | 'developer' | 'tester' | 'guest' | 'viewer';
