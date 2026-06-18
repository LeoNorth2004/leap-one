/**
 * 应用入口
 *
 * - React 18 createRoot 并发渲染模式
 * - StrictMode 开发模式检查
 * - 全局样式引入
 */

import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './assets/styles/global.css';

const rootElement = document.getElementById('root');

if (rootElement) {
  ReactDOM.createRoot(rootElement).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
}
