/** 用户头像组件 */

import { Avatar, Tooltip } from 'antd';
import { UserOutlined } from '@ant-design/icons';

interface UserAvatarProps {
  /** 头像URL */
  src?: string;
  /** 用户名（用于显示tooltip和头像文字） */
  name?: string;
  /** 头像大小 */
  size?: number;
}

export default function UserAvatar({ src, name = '', size = 32 }: UserAvatarProps) {
  const initials = name
    ? name.length > 1
      ? name.slice(-2)
      : name.toUpperCase()
    : undefined;

  return (
    <Tooltip title={name}>
      <Avatar
        src={src}
        icon={!src ? <UserOutlined /> : undefined}
        size={size}
        style={{ backgroundColor: !src ? '#1677ff' : undefined }}
      >
        {!src && initials}
      </Avatar>
    </Tooltip>
  );
}
