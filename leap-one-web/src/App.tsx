/**
 * 应用根组件
 *
 * 功能：
 * - QueryClientProvider（TanStack Query 配置）
 * - BrowserRouter（路由支持）
 * - AppRoutes（路由 + Ant Design ConfigProvider）
 */

import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import AppRoutes from '@/routes';

// ── TanStack Query 配置 ───────────────────────────────────────

const QUERY_DEFAULTS = Object.freeze({
  queries: {
    staleTime: 5 * 60 * 1000,
    retry: 1,
    refetchOnWindowFocus: false,
  },
  mutations: {
    retry: 0,
  },
});

const queryClient = new QueryClient({ defaultOptions: QUERY_DEFAULTS });

// ── 根组件 ───────────────────────────────────────────────────

const App = () => (
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <AppRoutes />
    </BrowserRouter>
  </QueryClientProvider>
);

export default App;
