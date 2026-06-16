/** 用户列表页面 */

import { useState } from 'react';
import { Table, Button, Tag, Space, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, StopOutlined, CheckCircleOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';
import UserAvatar from '@/components/Business/UserAvatar';
import StatusTag from '@/components/Common/StatusTag';
import ConfirmModal from '@/components/Common/ConfirmModal';
import { getUserListApi, deleteUserApi, updateUserStatusApi } from '@/api/user';
import type { UserDetail } from '@/types/user';

export default function UserList() {
  const [searchValue, setSearchValue] = useState('');
  const [users, setUsers] = useState<UserDetail[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);

  /** 删除确认弹窗状态 */
  const [deleteModal, setDeleteModal] = useState<{ open: boolean; user: UserDetail | null }>({
    open: false,
    user: null,
  });

  /** 加载用户列表 */
  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await getUserListApi({
        keyword: searchValue || undefined,
        page,
        pageSize,
      });
      if (res.code === 0 && res.data) {
        setUsers(res.data.list);
        setTotal(res.data.total);
      }
    } finally {
      setLoading(false);
    }
  };

  /** 初始化加载和搜索时重新获取数据 */
  useState(() => { fetchUsers(); });
  // 注意：实际项目中应使用 useEffect 或 TanStack Query

  /** 表格列定义 */
  const columns: ColumnsType<UserDetail> = [
    {
      title: '用户',
      dataIndex: 'realName',
      key: 'realName',
      render: (_: unknown, record) => (
        <Space>
          <UserAvatar src={record.avatar} name={record.realName} size={28} />
          <span>{record.realName}</span>
          <span style={{ color: '#999', fontSize: 12 }}>@{record.username}</span>
        </Space>
      ),
    },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '部门', dataIndex: 'departmentName', key: 'departmentName' },
    { title: '职位', dataIndex: 'position', key: 'position' },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
      render: (roles: UserDetail['roles']) => (
        <>
          {roles.map((role) => (
            <Tag key={role.id} color="blue">{role.name}</Tag>
          ))}
        </>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (status: UserDetail['status']) => (
        <StatusTag status={status} statusMap={{
          active: { label: '正常', color: 'success' },
          disabled: { label: '已禁用', color: 'error' },
        }} />
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />}>
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            icon={record.status === 'active' ? <StopOutlined /> : <CheckCircleOutlined />}
            onClick={() => handleToggleStatus(record)}
          >
            {record.status === 'active' ? '禁用' : '启用'}
          </Button>
          <Button
            type="link"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => setDeleteModal({ open: true, user: record })}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  /** 切换用户状态 */
  const handleToggleStatus = async (user: UserDetail) => {
    const newStatus = user.status === 'active' ? 'disabled' : 'active';
    try {
      await updateUserStatusApi(user.id, newStatus);
      message.success(`用户已${newStatus === 'active' ? '启用' : '禁用'}`);
      fetchUsers();
    } catch {
      // 错误已在拦截器中处理
    }
  };

  /** 确认删除用户 */
  const handleConfirmDelete = async () => {
    if (!deleteModal.user) return;
    try {
      await deleteUserApi(deleteModal.user.id);
      message.success('用户已删除');
      setDeleteModal({ open: false, user: null });
      fetchUsers();
    } catch {
      // 错误已在拦截器中处理
    }
  };

  return (
    <div>
      <PageHeader
        title="用户管理"
        subtitle="管理系统中的所有用户账号"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            新增用户
          </Button>
        }
      />

      <SearchBar
        value={searchValue}
        onChange={setSearchValue}
        onSearch={() => { setPage(1); fetchUsers(); }}
        placeholder="搜索用户名、姓名或邮箱..."
      />

      <Table
        rowKey="id"
        columns={columns}
        dataSource={users}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); fetchUsers(); },
        }}
      />

      <ConfirmModal
        open={deleteModal.open}
        title="确认删除用户"
        content={`确定要删除用户 "${deleteModal.user?.realName}" 吗？此操作不可恢复。`}
        danger
        onConfirm={handleConfirmDelete}
        onCancel={() => setDeleteModal({ open: false, user: null })}
      />
    </div>
  );
}
