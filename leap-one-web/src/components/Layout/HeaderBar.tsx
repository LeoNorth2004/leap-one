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
import { Layout, Dropdown, Avatar, Input, Switch, Tooltip } from 'antd';
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
  DownOutlined,
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
  const [notifCount] = useState(3);

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

  const handleSearch = useCallback(
    (value: string) => {
      if (value.trim()) {
        navigate(`/search?q=${encodeURIComponent(value.trim())}`);
      }
    },
    [navigate]
  );

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
      <div className={styles.headerLeft}>
        <Tooltip title={sidebarCollapsed ? '展开菜单' : '收起菜单'}>
          <span className={styles.collapseBtn} onClick={toggleSidebar}>
            {sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          </span>
        </Tooltip>

        <BreadcrumbNav />
      </div>

      <div className={styles.headerRight}>
        <div className={styles.searchWrapper}>
          <Input
            placeholder="全局搜索..."
            prefix={<SearchOutlined style={{ color: '#86909c' }} />}
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            onKeyDown={handleSearchPress}
            allowClear
            className={styles.searchInput}
          />
        </div>

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

        <div className={styles.divider} />

        <Tooltip title="通知中心">
          <span className={styles.headerIcon} onClick={() => navigate('/notification')}>
            <BellOutlined />
            {notifCount > 0 && (
              <span className={styles.notificationBadge}>
                {notifCount > 99 ? '99+' : notifCount}
              </span>
            )}
          </span>
        </Tooltip>

        <Dropdown menu={{ items: userMenuItems }} placement="bottomRight" arrow>
          <div className={styles.userArea}>
            <Avatar
              src={user?.avatar}
              icon={!user?.avatar ? <UserOutlined /> : undefined}
              size={32}
              className={styles.userAvatar}
            />
            {!sidebarCollapsed && (
              <>
                <span className={styles.userName}>
                  {user?.realName ?? '用户'}
                </span>
                <DownOutlined style={{ fontSize: 12, color: '#86909c' }} />
              </>
            )}
          </div>
        </Dropdown>
      </div>
    </Header>
  );
}
