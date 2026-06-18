/** 个人中心页面 */

import { Card, Form, Input, Button, Avatar, Upload, message, Tabs, Tag } from 'antd';
import { UserOutlined, MailOutlined, CameraOutlined, LockOutlined } from '@ant-design/icons';
import useAuth from '@/hooks/useAuth';
import PageHeader from '@/components/Common/PageHeader';

export default function Profile() {
  const { user } = useAuth();
  const [profileForm] = Form.useForm();
  const [pwdForm] = Form.useForm();

  /** 更新个人信息 */
  const handleUpdateProfile = () => {
    profileForm.validateFields().then(() => {
      message.success('个人信息更新成功');
    });
  };

  /** 修改密码 */
  const handleChangePassword = () => {
    pwdForm.validateFields().then((values) => {
      // TODO: 调用修改密码API
      console.log('修改密码:', values);
      message.success('密码修改成功');
      pwdForm.resetFields();
    });
  };

  return (
    <div>
      <PageHeader title="个人中心" subtitle="管理您的个人信息和账户设置" />

      <Card>
        <Tabs
          defaultActiveKey="profile"
          items={[
            {
              key: 'profile',
              label: '基本信息',
              children: (
                <div style={{ maxWidth: 600 }}>
                  {/* 头像区域 */}
                  <div style={{ textAlign: 'center', marginBottom: 32 }}>
                    <Upload
                      showUploadList={false}
                      beforeUpload={() => false}
                      onChange={(info) => {
                        if (info.file.originFileObj) {
                          message.success('头像上传成功');
                        }
                      }}
                    >
                      <Avatar
                        src={user?.avatar}
                        icon={!user?.avatar ? <UserOutlined /> : undefined}
                        size={100}
                        style={{ cursor: 'pointer', border: '2px dashed #d9d9d9' }}
                      >
                        <CameraOutlined style={{ position: 'absolute', bottom: 0, right: 0, background: '#1677ff', color: '#fff', borderRadius: '50%', padding: 4, fontSize: 12 }} />
                      </Avatar>
                    </Upload>
                    <p style={{ marginTop: 8, color: '#999' }}>点击更换头像</p>
                  </div>

                  <Form
                    form={profileForm}
                    layout="vertical"
                    initialValues={{
                      realName: user?.realName || '',
                      email: user?.email || '',
                      phone: user?.phone || '',
                    }}
                  >
                    <Form.Item label="用户名">
                      <Input value={user?.username} disabled />
                    </Form.Item>
                    <Form.Item label="真实姓名" name="realName" rules={[{ required: true, message: '请输入姓名' }]}>
                      <Input prefix={<UserOutlined />} placeholder="请输入真实姓名" />
                    </Form.Item>
                    <Form.Item label="邮箱地址" name="email" rules={[{ required: true, type: 'email', message: '请输入有效邮箱' }]}>
                      <Input prefix={<MailOutlined />} placeholder="请输入邮箱地址" />
                    </Form.Item>
                    <Form.Item label="手机号码" name="phone">
                      <Input placeholder="请输入手机号码" />
                    </Form.Item>
                    <Form.Item label="所属部门">
                      <Input value={user?.department || '-'} disabled />
                    </Form.Item>
                    <Form.Item label="角色">
                      {(user?.roles || []).map((role: string) => (
                        <Tag key={role} color="blue">{role}</Tag>
                      ))}
                    </Form.Item>
                    <Button type="primary" icon={<UserOutlined />} onClick={handleUpdateProfile}>
                      保存信息
                    </Button>
                  </Form>
                </div>
              ),
            },
            {
              key: 'security',
              label: '账户安全',
              children: (
                <div style={{ maxWidth: 500 }}>
                  <h4 style={{ marginBottom: 16 }}><LockOutlined /> 修改密码</h4>
                  <Form form={pwdForm} layout="vertical">
                    <Form.Item label="当前密码" name="oldPassword"
                      rules={[{ required: true, message: '请输入当前密码' }]}>
                      <Input.Password placeholder="请输入当前密码" />
                    </Form.Item>
                    <Form.Item label="新密码" name="newPassword"
                      rules={[
                        { required: true, message: '请输入新密码' },
                        { min: 8, message: '密码至少8位' },
                      ]}>
                      <Input.Password placeholder="请输入新密码（至少8位）" />
                    </Form.Item>
                    <Form.Item label="确认新密码" name="confirmPassword"
                      rules={[
                        { required: true, message: '请确认新密码' },
                        ({ getFieldValue }) => ({
                          validator(_, value) {
                            if (value && getFieldValue('newPassword') !== value) {
                              return Promise.reject(new Error('两次密码不一致'));
                            }
                            return Promise.resolve();
                          },
                        }),
                      ]}>
                      <Input.Password placeholder="请再次输入新密码" />
                    </Form.Item>
                    <Button type="primary" icon={<LockOutlined />} onClick={handleChangePassword}>
                      确认修改
                    </Button>
                  </Form>
                </div>
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
}
