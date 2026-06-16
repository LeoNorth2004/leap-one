/**
 * 顶部导航栏组件
 *
 * 功能：
 * - 左侧：折叠按钮 + 面包屑导航
 * - 右侧：全局搜索框、主题切换、通知铃铛（带红点）、用户头像下拉菜单
 * - 用户下拉菜单包含：个人设置、退出登录
 */

import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Dropdown, Avatar, Badge, Input, Switch, Tooltip } from 'antd';
import type { MenuProps } from 'antd';
import {
  BellOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SunOutlined,
  MoonOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import BreadcrumbNav from './BreadcrumbNav';
import { useAuth } from '@/hooks/useAuth';
import { useAppStore } from '@/store/appStore';
import { useTheme } from '@/hooks/useTheme';
import styles from './HeaderBar.module.less';

const { Header } = Layout;

export default function HeaderBar() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const toggleSidebar = useAppStore((state) => state.toggleSidebar);
  const sidebarCollapsed = useAppStore((state) => state.sidebarCollapsed);
  const { isDark, toggleTheme } = useTheme();
  const [searchValue, setSearchValue] = useState('');
  const [notifCount] = useState(3); // TODO: 从通知 API 获取未读数

  /** 用户下拉菜单项 */
  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '系统设置',
      onClick: () => navigate('/settings'),
    },
    { type: 'divider' as const },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
      onClick: () => logout(),
    },
  ];

  /** 全局搜索处理 */
  const handleSearch = useCallback(
    (value: string) => {
      if (value.trim()) {
        navigate(`/search?q=${encodeURIComponent(value.trim())}`);
      }
    },
    [navigate]
  );

  /** 回车搜索 */
  const handleSearchPress = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter') {
        handleSearch(searchValue);
      }
    },
    [searchValue, handleSearch]
  );

  return (
    <Header className={styles.headerBar}>
      {/* ── 左侧区域 ──────────────────────────────────────── */}
      <div className={styles.headerLeft}>
        {/* 折叠 / 展开侧边栏按钮 */}
        <Tooltip title={sidebarCollapsed ? '展开菜单' : '收起菜单'}>
          <span className={styles.collapseBtn} onClick={toggleSidebar}>
            {sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          </span>
        </Tooltip>

        {/* 面包屑导航 */}
        <BreadcrumbNav />
      </div>

      {/* ── 右侧区域 ──────────────────────────────────────── */}
      <div className={styles.headerRight}>
        {/* 全局搜索框 */}
        <Input
          placeholder="全局搜索..."
          prefix={<SearchOutlined style={{ color: '#86909c' }} />}
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
          onKeyDown={handleSearchPress}
          allowClear
          style={{ width: 200 }}
          className={styles.searchInput}
        />

        {/* 主题切换开关 */}
        <Tooltip title={isDark ? '切换亮色模式' : '切换暗色模式'}>
          <Switch
            checkedChildren={<MoonOutlined />}
            unCheckedChildren={<SunOutlined />}
            checked={isDark}
            onChange={toggleTheme}
            size="small"
            className={styles.themeSwitch}
          />
        </Tooltip>

        {/* 通知铃铛（带未读数红点） */}
        <Tooltip title="通知中心">
          <Badge count={notifCount} size="small" offset={[-2, 2]}>
            <BellOutlined
              className={styles.headerIcon}
              onClick={() => navigate('/notification')}
            />
          </Badge>
        </Tooltip>

        {/* 用户头像 + 下拉菜单 */}
        <Dropdown menu={{ items: userMenuItems }} placement="bottomRight" arrow>
          <div className={styles.userArea}>
            <Avatar
              src={user?.avatar}
              icon={!user?.avatar ? <UserOutlined /> : undefined}
              size={32}
              style={{ flexShrink: 0 }}
            />
            {!sidebarCollapsed && (
              <span className={styles.userName}>
                {user?.realName ?? '用户'}
              </span>
            )}
          </div>
        </Dropdown>
      </div>
    </Header>
  );
}
