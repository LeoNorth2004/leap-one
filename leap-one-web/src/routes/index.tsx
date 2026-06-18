/**
 * 路由主文件 - 使用 React Router v6 配置所有路由
 *
 * 功能：
 * - BrowserRouter 包裹整个应用
 * - MainLayout 作为父路由包裹所有需要认证的页面
 * - AuthGuard 保护需要认证的路由
 * - 登录页独立于 MainLayout
 * - React.lazy + Suspense 实现路由懒加载
 * - 404 兜底页面
 * - Ant Design ConfigProvider 全局配置（中文 + 主题）
 */

import { Suspense } from 'react';
import { useRoutes, Navigate, type RouteObject } from 'react-router-dom';
import { ConfigProvider, theme as antdTheme, Spin } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { routesConfig, type RouteMeta } from './routes.config';
import AuthGuard from './AuthGuard';
import useAppStore from '@/store/appStore';

/** 全局 Loading 组件（用于 Suspense fallback） */
function PageLoading() {
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: 'calc(100vh - 120px)',
    }}>
      <Spin size="large" tip="加载中..." />
    </div>
  );
}

/**
 * 渲染路由配置树
 * - 自动为需要认证的子路由包裹 AuthGuard
 * - 为懒加载组件包裹 Suspense
 */
function RenderRoutes() {
  const element = useRoutes(
    routesConfig.map((route): RouteObject => {
      const routeWithMeta = route as RouteObject & { meta?: RouteMeta };
      
      // 登录页直接返回，无需守卫和 Suspense（自身已处理）
      if (routeWithMeta.path === '/login') {
        return routeWithMeta;
      }

      // 主布局路由：为子路由添加 AuthGuard 和 Suspense
      if (routeWithMeta.meta?.requiresAuth !== false && routeWithMeta.children) {
        return {
          ...routeWithMeta,
          children: routeWithMeta.children.map((child): RouteObject => {
            const childWithMeta = child as RouteObject & { meta?: RouteMeta };
            const Element = childWithMeta.element as unknown as React.ComponentType;
            return {
              ...childWithMeta,
              element: Element ? (
                <Suspense fallback={<PageLoading />}>
                  <AuthGuard>
                    <Element />
                  </AuthGuard>
                </Suspense>
              ) : undefined,
            };
          }),
        };
      }

      return routeWithMeta;
    })
  );

  return element ?? <Navigate to="/" replace />;
}

/**
 * 应用根路由组件
 *
 * 包含：
 * - Ant Design ConfigProvider（中文国际化 + 主题配置）
 * - 暗色/亮色主题自动切换
 */
export default function AppRoutes() {
  const appTheme = useAppStore((state) => state.theme);

  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm:
          appTheme === 'dark'
            ? antdTheme.darkAlgorithm
            : antdTheme.defaultAlgorithm,
        token: {
          colorPrimary: '#1677ff',
          borderRadius: 6,
          fontFamily:
            '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif',
        },
        components: {
          Layout: {
            siderBg: appTheme === 'dark' ? '#1f1f1f' : '#fff',
            headerBg: appTheme === 'dark' ? '#1f1f1f' : '#fff',
            bodyBg: appTheme === 'dark' ? '#141414' : '#f5f7fa',
          },
          Menu: {
            itemBg: 'transparent',
          },
        },
      }}
    >
      <RenderRoutes />
    </ConfigProvider>
  );
}
