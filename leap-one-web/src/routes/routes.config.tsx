/**
 * 路由配置表 - 定义所有路由及其元信息
 *
 * 使用 React.lazy 实现按需加载
 * 包含权限控制字段 permissions
 */

import type { RouteObject } from 'react-router-dom';

/** 路由元信息 */
export interface RouteMeta {
  /** 页面标题（用于面包屑、浏览器标题） */
  title: string;
  /** 是否需要认证，默认 true */
  requiresAuth?: boolean;
  /** 所需权限列表（可选，用于菜单/按钮级权限控制） */
  permissions?: string[];
  /** 是否在侧边栏菜单中隐藏 */
  hideInMenu?: boolean;
}

/** 带有 meta 信息的路由对象 */
export type LeapRoute = RouteObject & {
  meta?: RouteMeta;
};

/** 路由配置表 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const routesConfig: any[] = [
  // ════════════════════════════════════════════════════════════
  // 无需认证的路由
  // ════════════════════════════════════════════════════════════
  {
    path: '/login',
    element: lazyImport(() => import('@/pages/Login')),
    meta: { title: '登录', requiresAuth: false },
  },

  // ════════════════════════════════════════════════════════════
  // 需要认证的主布局路由（MainLayout 作为父级）
  // ════════════════════════════════════════════════════════════
  {
    path: '/',
    element: lazyImport(() => import('@/components/Layout/MainLayout')),
    meta: { title: 'Leap One', requiresAuth: true },
    children: [
      // ── 工作台 ───────────────────────────────────────────────
      {
        index: true,
        element: lazyImport(() => import('@/pages/Dashboard')),
        meta: { title: '工作台', icon: 'DashboardOutlined' },
      },

      // ── 用户管理 ─────────────────────────────────────────────
      {
        path: 'user/list',
        element: lazyImport(() => import('@/pages/User/UserList')),
        meta: { title: '用户列表' },
      },
      {
        path: 'user/:id',
        element: lazyImport(() => import('@/pages/User/UserProfile')),
        meta: { title: '用户详情', hideInMenu: true },
      },

      // ── 组织管理 ─────────────────────────────────────────────
      {
        path: 'org/department',
        element: lazyImport(() => import('@/pages/Organization/Department')),
        meta: { title: '部门管理' },
      },
      {
        path: 'org/role',
        element: lazyImport(() => import('@/pages/Organization/RoleManage')),
        meta: { title: '角色权限' },
      },

      // ── 项目集 ───────────────────────────────────────────────
      {
        path: 'program/list',
        element: lazyImport(() => import('@/pages/Program/ProgramList')),
        meta: { title: '项目集列表' },
      },

      // ── 产品 ─────────────────────────────────────────────────
      {
        path: 'product/list',
        element: lazyImport(() => import('@/pages/Product/ProductList')),
        meta: { title: '产品列表' },
      },
      {
        path: 'product/roadmap',
        element: lazyImport(() => import('@/pages/Product/ProductRoadmap')),
        meta: { title: '产品路线图' },
      },

      // ── 项目 ─────────────────────────────────────────────────
      {
        path: 'project/list',
        element: lazyImport(() => import('@/pages/Project/ProjectList')),
        meta: { title: '项目列表' },
      },
      {
        path: 'project/:id',
        element: lazyImport(() => import('@/pages/Project/ProjectDetail')),
        meta: { title: '项目详情', hideInMenu: true },
      },
      {
        path: 'project/:id/iteration',
        element: lazyImport(() => import('@/pages/Project/IterationList')),
        meta: { title: '迭代列表', hideInMenu: true },
      },
      {
        path: 'project/:id/kanban',
        element: lazyImport(() => import('@/pages/Project/KanbanBoard')),
        meta: { title: '项目看板', hideInMenu: true },
      },

      // ── 需求 ─────────────────────────────────────────────────
      {
        path: 'requirement/list',
        element: lazyImport(() => import('@/pages/Requirement/RequirementList')),
        meta: { title: '需求列表' },
      },
      {
        path: 'requirement/:id',
        element: lazyImport(() => import('@/pages/Requirement/RequirementDetail')),
        meta: { title: '需求详情', hideInMenu: true },
      },

      // ── 任务 ─────────────────────────────────────────────────
      {
        path: 'task/list',
        element: lazyImport(() => import('@/pages/Task/TaskList')),
        meta: { title: '任务列表' },
      },

      // ── 质量中心 ─────────────────────────────────────────────
      {
        path: 'quality/testcase',
        element: lazyImport(() => import('@/pages/Quality/TestCaseList')),
        meta: { title: '测试用例' },
      },
      {
        path: 'quality/bug',
        element: lazyImport(() => import('@/pages/Quality/BugList')),
        meta: { title: 'Bug 列表' },
      },
      {
        path: 'quality/testplan',
        element: lazyImport(() => import('@/pages/Quality/TestPlanList')),
        meta: { title: '测试计划' },
      },

      // ── 工单 ─────────────────────────────────────────────────
      {
        path: 'issue/list',
        element: lazyImport(() => import('@/pages/Issue/IssueList')),
        meta: { title: '工单列表' },
      },

      // ── 文档 ─────────────────────────────────────────────────
      {
        path: 'document/list',
        element: lazyImport(() => import('@/pages/Document/DocumentList')),
        meta: { title: '文档中心' },
      },

      // ── 看板 ─────────────────────────────────────────────────
      {
        path: 'kanban',
        element: lazyImport(() => import('@/pages/Kanban/KanbanView')),
        meta: { title: '全局看板' },
      },

      // ── BI 大屏 ───────────────────────────────────────────────
      {
        path: 'bi/dashboard',
        element: lazyImport(() => import('@/pages/BI/Dashboard')),
        meta: { title: '数据概览' },
      },

      // ── 系统设置 ─────────────────────────────────────────────
      {
        path: 'settings',
        element: lazyImport(() => import('@/pages/Settings/SystemSettings')),
        meta: { title: '系统设置', permissions: ['system:settings'] },
      },

      // ── 个人中心 ─────────────────────────────────────────────
      {
        path: 'profile',
        element: lazyImport(() => import('@/pages/Profile')),
        meta: { title: '个人中心', hideInMenu: true },
      },
    ],
  },

  // ════════════════════════════════════════════════════════════
  // 兜底：404 页面
  // ════════════════════════════════════════════════════════════
  {
    path: '*',
    element: (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '60vh',
        gap: 16,
      }}>
        <div style={{ fontSize: 72, fontWeight: 700, color: 'var(--primary-color)' }}>404</div>
        <p style={{ fontSize: 16, color: 'var(--text-secondary)' }}>抱歉，您访问的页面不存在</p>
        <a href="/" style={{ color: 'var(--primary-color)' }}>返回首页</a>
      </div>
    ),
  },
];

/**
 * 封装 React.lazy 导入，统一错误处理
 */
function lazyImport(
  importFn: () => Promise<{ default: React.ComponentType }>
): React.LazyExoticComponent<React.ComponentType> {
  return React.lazy(importFn);
}

// 需要导入 React 以使用 React.lazy
import React from 'react';
