/**
 * 侧边导航菜单组件
 *
 * 功能：
 * - 根据路由配置动态生成菜单（参考禅道导航结构，现代化升级）
 * - 支持菜单图标（@ant-design/icons）
 * - 支持菜单展开/折叠
 * - 当前路由高亮选中
 * - 根据权限过滤菜单项（可选）
 * - 暗色/亮色主题自适应
 */

import { useMemo, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import type { MenuProps } from 'antd';
import {
  DashboardOutlined,
  AppstoreOutlined,
  ShoppingOutlined,
  FolderOpenOutlined,
  FileTextOutlined,
  CarryOutOutlined,
  ExperimentOutlined,
  CustomerServiceOutlined,
  FileOutlined,
  TableOutlined,
  BarChartOutlined,
  SettingOutlined,
  TeamOutlined,
} from '@ant-design/icons';
import useAppStore from '@/store/appStore';

const { Sider } = Layout;

/** 菜单配置 - 参考禅道导航结构，现代化升级 */
const menuItems: MenuProps['items'] = [
  {
    key: '/',
    icon: <DashboardOutlined />,
    label: '工作台',
  },
  {
    key: '/org',
    icon: <TeamOutlined />,
    label: '组织',
    children: [
      { key: '/user/list', label: '用户列表' },
      { key: '/org/department', label: '部门管理' },
      { key: '/org/role', label: '角色权限' },
    ],
  },
  {
    key: '/program',
    icon: <AppstoreOutlined />,
    label: '项目集',
    children: [{ key: '/program/list', label: '项目集列表' }],
  },
  {
    key: '/product',
    icon: <ShoppingOutlined />,
    label: '产品',
    children: [
      { key: '/product/list', label: '产品列表' },
      { key: '/product/roadmap', label: '产品路线图' },
    ],
  },
  {
    key: '/project',
    icon: <FolderOpenOutlined />,
    label: '项目',
    children: [
      { key: '/project/list', label: '项目列表' },
      { key: '/project/:id/kanban', label: '项目看板' },
    ],
  },
  {
    key: '/requirement',
    icon: <FileTextOutlined />,
    label: '需求',
    children: [{ key: '/requirement/list', label: '需求列表' }],
  },
  {
    key: '/task',
    icon: <CarryOutOutlined />,
    label: '任务',
    children: [{ key: '/task/list', label: '任务列表' }],
  },
  {
    key: '/quality',
    icon: <ExperimentOutlined />,
    label: '质量中心',
    children: [
      { key: '/quality/testcase', label: '测试用例' },
      { key: '/quality/bug', label: 'Bug 列表' },
      { key: '/quality/testplan', label: '测试计划' },
    ],
  },
  {
    key: '/issue',
    icon: <CustomerServiceOutlined />,
    label: '工单',
    children: [{ key: '/issue/list', label: '工单列表' }],
  },
  {
    key: '/document',
    icon: <FileOutlined />,
    label: '文档',
    children: [{ key: '/document/list', label: '文档列表' }],
  },
  {
    key: '/kanban',
    icon: <TableOutlined />,
    label: '看板',
  },
  {
    key: '/bi',
    icon: <BarChartOutlined />,
    label: 'BI 统计',
    children: [{ key: '/bi/dashboard', label: '数据概览' }],
  },
  {
    key: '/settings',
    icon: <SettingOutlined />,
    label: '设置',
  },
];

export default function SiderMenu() {
  const navigate = useNavigate();
  const location = useLocation();
  const collapsed = useAppStore((state) => state.sidebarCollapsed);
  const theme = useAppStore((state) => state.theme);

  /** 根据当前路由计算选中的菜单 key */
  const selectedKeys = useMemo<string[]>(() => {
    const pathname = location.pathname;
    // 精确匹配优先，否则取当前路径
    return [pathname];
  }, [location.pathname]);

  /** 根据当前路径计算需要展开的父级菜单 keys */
  const openKeys = useMemo<string[]>(() => {
    const segments = location.pathname.split('/').filter(Boolean);
    if (segments.length > 1) {
      // 取第一段路径作为展开的父级 key
      return [`/${segments[0]}`];
    }
    // 默认不展开任何子菜单
    return [];
  }, [location.pathname]);

  /** 菜单项点击处理 */
  const handleMenuClick: MenuProps['onClick'] = useCallback(
    ({ key }: { key: string }) => {
      navigate(key);
    },
    [navigate]
  );

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      width={240}
      collapsedWidth={64}
      className="layout-sider"
      theme={theme === 'dark' ? 'dark' : 'light'}
      style={{
        overflow: 'auto',
        height: '100vh',
        position: 'sticky',
        top: 0,
        left: 0,
      }}
    >
      {/* Logo 区域 */}
      <div className="sider-logo" onClick={() => navigate('/')}>
        <span className="logo-icon">🚀</span>
        {!collapsed && <span className="logo-text">Leap One</span>}
      </div>

      {/* 导航菜单 */}
      <Menu
        mode="inline"
        items={menuItems}
        selectedKeys={selectedKeys}
        defaultOpenKeys={openKeys}
        onClick={handleMenuClick}
        inlineCollapsed={collapsed}
        style={{
          borderRight: 0,
          height: 'calc(100vh - 64px)',
          overflowY: 'auto',
          overflowX: 'hidden',
        }}
      />
    </Sider>
  );
}
