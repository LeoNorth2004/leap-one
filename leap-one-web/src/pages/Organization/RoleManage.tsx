/** 角色权限管理页面 */

import { Table, Button, Tag } from 'antd';
import { PlusOutlined, EditOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';

interface RoleItem {
  id: number;
  name: string;
  code: string;
  description: string;
  userCount: number;
  permissions: string[];
}

/** 模拟角色数据 - TODO: 替换为API */
const mockRoles: RoleItem[] = [
  { id: 1, name: '超级管理员', code: 'super_admin', description: '拥有系统全部权限', userCount: 2, permissions: ['*'] },
  { id: 2, name: '管理员', code: 'admin', description: '组织内管理权限', userCount: 5, permissions: ['user:*', 'project:*'] },
  { id: 3, name: '项目经理', code: 'pm', description: '项目管理相关权限', userCount: 12, permissions: ['project:*', 'task:*'] },
  { id: 4, name: '开发工程师', code: 'developer', description: '开发任务执行权限', userCount: 30, permissions: ['task:view', 'task:edit'] },
  { id: 5, name: '测试工程师', code: 'tester', description: '测试质量相关权限', userCount: 10, permissions: ['quality:*'] },
  { id: 6, name: '只读用户', code: 'viewer', description: '仅查看权限', userCount: 20, permissions: ['view:*'] },
];

export default function RoleManage() {
  return (
    <div>
      <PageHeader
        title="角色与权限管理"
        subtitle="配置系统角色及其对应的权限矩阵"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            新增角色
          </Button>
        }
      />

      <Table
        rowKey="id"
        dataSource={mockRoles}
        pagination={false}
        columns={[
          { title: '角色名称', dataIndex: 'name', render: (name: string) => <strong>{name}</strong> },
          { title: '角色标识', dataIndex: 'code', render: (code: string) => <Tag color="blue">{code}</Tag> },
          { title: '描述', dataIndex: 'description' },
          { title: '关联用户数', dataIndex: 'userCount', align: 'center' },
          {
            title: '操作',
            width: 120,
            render: () => (
              <Button type="link" icon={<EditOutlined />}>
                编辑权限
              </Button>
            ),
          },
        ]}
      />
    </div>
  );
}
