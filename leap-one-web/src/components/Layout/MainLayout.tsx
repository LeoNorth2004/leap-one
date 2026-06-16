/**
 * 主布局组件
 *
 * 功能：
 * - Ant Design Layout 构建经典后台布局（侧边栏 + 顶栏 + 内容区）
 * - 左侧 SiderMenu 导航菜单
 * - 顶部 HeaderBar 工具栏
 * - 中间内容区渲染子路由（Outlet）
 * - AI 悬浮球按钮
 * - 支持暗色 / 亮色主题切换
 */

import { Outlet } from 'react-router-dom';
import { Layout } from 'antd';
import SiderMenu from './SiderMenu';
import HeaderBar from './HeaderBar';
import AIFloatingButton from '../Business/AIFloatingButton';
import styles from './MainLayout.module.less';

const { Content } = Layout;

export default function MainLayout() {
  return (
    <Layout className={styles.mainLayout}>
      {/* 左侧导航菜单 */}
      <SiderMenu />

      {/* 右侧主体区域 */}
      <Layout className={styles.rightLayout}>
        {/* 顶部工具栏 */}
        <HeaderBar />

        {/* 页面内容区 */}
        <Content className={styles.contentArea}>
          <Outlet />
        </Content>
      </Layout>

      {/* AI 悬浮助手按钮 */}
      <AIFloatingButton />
    </Layout>
  );
}
