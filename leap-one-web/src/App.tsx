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

// 创建 TanStack Query 客户端实例
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 数据 5 分钟内视为新鲜
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </QueryClientProvider>
  );
}
