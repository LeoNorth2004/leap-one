/** 用户详情/个人资料页面 */

import { useState } from 'react';
import { Descriptions, Card, Tag, Divider, Button, Space } from 'antd';
import { MailOutlined, PhoneOutlined, TeamOutlined } from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import PageHeader from '@/components/Common/PageHeader';
import UserAvatar from '@/components/Business/UserAvatar';
import { getUserDetailApi } from '@/api/user';
import type { UserDetail } from '@/types/user';

export default function UserProfile() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [user, setUser] = useState<UserDetail | null>(null);
  const [loading, setLoading] = useState(true);

  /** 加载用户详情 */
  useState(() => {
    if (id) {
      getUserDetailApi(Number(id))
        .then(setUser)
        .finally(() => setLoading(false));
    }
  });

  return (
    <div>
      <PageHeader
        title="用户详情"
        extra={
          <Space>
            <Button onClick={() => navigate(-1)}>返回</Button>
            <Button type="primary">编辑信息</Button>
          </Space>
        }
      />

      {/* 用户基本信息卡片 */}
      <Card loading={loading}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <UserAvatar src={user?.avatar} name={user?.realName} size={80} />
          <h2 style={{ marginTop: 12, marginBottom: 4 }}>{user?.realName}</h2>
          <p style={{ color: '#999', marginBottom: 0 }}>@{user?.username}</p>
        </div>

        <Descriptions column={{ xxl: 3, xl: 2, lg: 1 }} bordered size="middle">
          <Descriptions.Item label="用户ID">{user?.id}</Descriptions.Item>
          <Descriptions.Item label="性别">
            {user?.gender === 'male' ? '男' : user?.gender === 'female' ? '女' : '未知'}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={user?.status === 'active' ? 'success' : 'error'}>
              {user?.status === 'active' ? '正常' : '已禁用'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label={<><MailOutlined /> 邮箱</>}>
            {user?.email}
          </Descriptions.Item>
          <Descriptions.Item label={<><PhoneOutlined /> 手机</>}>
            {user?.phone || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={<><TeamOutlined /> 部门</>}>
            {user?.departmentName}
          </Descriptions.Item>
          <Descriptions.Item label="职位">{user?.position || '-'}</Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {user?.createdAt ? new Date(user.createdAt).toLocaleString('zh-CN') : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="最后登录">
            {user?.lastLoginTime ? new Date(user.lastLoginTime).toLocaleString('zh-CN') : '-'}
          </Descriptions.Item>
        </Descriptions>

        <Divider orientation="left">角色与权限</Divider>
        <div>
          {(user?.roles || []).map((role) => (
            <Tag key={role.id} color="blue" style={{ marginBottom: 8, marginRight: 8 }}>
              {role.name} ({role.code})
            </Tag>
          ))}
        </div>
      </Card>
    </div>
  );
}
