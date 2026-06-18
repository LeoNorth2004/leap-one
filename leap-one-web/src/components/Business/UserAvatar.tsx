/**
 * 用户头像组件
 *
 * 支持图片头像 / 文字头像（取用户名后两位）/ 默认图标三种模式
 * 悬停显示 Tooltip 用户名
 */

import { Avatar, Tooltip } from 'antd';
import { UserOutlined } from '@ant-design/icons';

// ── 类型定义 ─────────────────────────────────────────────────

interface UserAvatarProps {
  /** 头像图片地址 */
  src?: string;
  /** 用户名（用于 Tooltip 和文字头像） */
  name?: string;
  /** 头像尺寸 */
  size?: number;
}

// ── 默认值 ───────────────────────────────────────────────────

const DEFAULT_SIZE = 32;
const FALLBACK_COLOR = '#1677ff';

// ── 工具函数 ─────────────────────────────────────────────────

/** 从用户名提取头像文字 */
const extractInitials = (name: string): string | undefined => {
  if (!name) return undefined;
  return name.length > 1 ? name.slice(-2) : name.toUpperCase();
};

// ── 组件实现 ─────────────────────────────────────────────────

const UserAvatar = ({ src, name = '', size = DEFAULT_SIZE }: UserAvatarProps) => {
  const initials = src ? undefined : extractInitials(name);
  const bgColor = src ? undefined : FALLBACK_COLOR;

  return (
    <Tooltip title={name}>
      <Avatar
        src={src}
        icon={!src ? <UserOutlined /> : undefined}
        size={size}
        style={{ backgroundColor: bgColor }}
      >
        {!src && initials}
      </Avatar>
    </Tooltip>
  );
};

export default UserAvatar;
