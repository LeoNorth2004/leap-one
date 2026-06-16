/**
 * 应用入口文件
 *
 * - ReactDOM.createRoot（React 18 并发模式）
 * - 渲染 App 根组件
 * - StrictMode（开发模式额外检查）
 * - 导入全局样式
 */

import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './assets/styles/global.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
