/**
 * 常量定义
 *
 * 包含：API 状态码、分页默认值、任务/Bug/需求/项目枚举映射等
 */

// ─── API 状态码 ───────────────────────────────────────────────────

export const API_CODE = {
  SUCCESS: 0,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  SERVER_ERROR: 500,
} as const;

// ─── 分页默认值 ───────────────────────────────────────────────────

export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
} as const;

// ─── 任务状态映射 ─────────────────────────────────────────────────

export const TASK_STATUS_MAP: Record<string, { label: string; color: string }> = {
  wait: { label: '待处理', color: 'default' },
  doing: { label: '进行中', color: 'processing' },
  done: { label: '已完成', color: 'success' },
  pause: { label: '已暂停', color: 'warning' },
  cancel: { label: '已取消', color: 'error' },
  closed: { label: '已关闭', color: 'default' },
};

/** 任务状态选项列表（用于筛选器） */
export const TASK_STATUS_OPTIONS = Object.entries(TASK_STATUS_MAP).map(([value, { label }]) => ({
  label,
  value,
}));

// ─── 任务优先级映射 ───────────────────────────────────────────────

export const PRIORITY_MAP: Record<string, { label: string; color: string }> = {
  P0: { label: '紧急', color: 'red' },
  urgent: { label: '紧急', color: 'red' },
  P1: { label: '高', color: 'orange' },
  high: { label: '高', color: 'orange' },
  P2: { label: '中', color: 'blue' },
  medium: { label: '中', color: 'blue' },
  P3: { label: '低', color: 'default' },
  low: { label: '低', color: 'default' },
};

/** 优先级选项列表 */
export const PRIORITY_OPTIONS = [
  { label: '紧急', value: 'P0' },
  { label: '高', value: 'P1' },
  { label: '中', value: 'P2' },
  { label: '低', value: 'P3' },
];

// ─── 需求状态映射 ─────────────────────────────────────────────────

export const REQUIREMENT_STATUS_MAP: Record<string, { label: string; color: string }> = {
  draft: { label: '草稿', color: 'default' },
  reviewing: { label: '评审中', color: 'processing' },
  active: { label: '激活', color: 'blue' },
  developing: { label: '开发中', color: 'cyan' },
  testing: { label: '测试中', color: 'orange' },
  completed: { label: '已完成', color: 'success' },
  closed: { label: '已关闭', color: 'default' },
  rejected: { label: '已拒绝', color: 'error' },
};

/** 需求类型映射 */
export const REQUIREMENT_TYPE_MAP: Record<string, { label: string; color: string }> = {
  user: { label: '用户需求', color: 'blue' },
  business: { label: '业务需求', color: 'green' },
  technical: { label: '技术需求', color: 'orange' },
  defect: { label: '缺陷修复', color: 'red' },
};

// ─── Bug 严重程度映射 ─────────────────────────────────────────────

export const BUG_SEVERITY_MAP: Record<string, { label: string; color: string }> = {
  fatal: { label: '致命', color: 'red' },
  serious: { label: '严重', color: 'orange' },
  normal: { label: '一般', color: 'blue' },
  slight: { label: '轻微', color: 'default' },
  suggest: { label: '建议', color: 'green' },
};

/** Bug 状态映射 */
export const BUG_STATUS_MAP: Record<string, { label: string; color: string }> = {
  active: { label: '激活', color: 'red' },
  resolved: { label: '已解决', color: 'success' },
  closed: { label: '已关闭', color: 'default' },
  postponed: { label: '已延期', color: 'warning' },
};

// ─── 项目状态映射 ─────────────────────────────────────────────────

export const PROJECT_STATUS_MAP: Record<string, { label: string; color: string }> = {
  planning: { label: '规划中', color: 'default' },
  active: { label: '进行中', color: 'processing' },
  paused: { label: '已暂停', color: 'warning' },
  completed: { label: '已完成', color: 'success' },
  archived: { label: '已归档', color: 'default' },
};

/** 项目类型映射 */
export const PROJECT_TYPE_MAP: Record<string, { label: string; color: string }> = {
  scrum: { label: 'Scrum', color: 'blue' },
  waterfall: { label: '瀑布', color: 'green' },
  kanban: { label: '看板', color: 'orange' },
  hybrid: { label: '混合', color: 'purple' },
};

// ─── 用户角色映射 ─────────────────────────────────────────────────

export const USER_ROLE_MAP: Record<string, { label: string; color: string }> = {
  admin: { label: '管理员', color: 'red' },
  project_manager: { label: '项目经理', color: 'blue' },
  developer: { label: '开发人员', color: 'green' },
  tester: { label: '测试人员', color: 'orange' },
  guest: { label: '访客', color: 'default' },
  viewer: { label: '只读用户', color: 'default' },
};
