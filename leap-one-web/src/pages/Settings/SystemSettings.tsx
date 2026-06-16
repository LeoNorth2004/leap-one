/** 系统设置页面 */

import { Card, Form, Input, Switch, Button, message, Tabs, Divider } from 'antd';
import { SaveOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';

export default function SystemSettings() {
  const [form] = Form.useForm();

  /** 保存基础设置 */
  const handleSaveBasic = () => {
    form.validateFields().then(() => {
      message.success('设置保存成功');
    });
  };

  return (
    <div>
      <PageHeader title="系统设置" subtitle="配置系统的全局参数和偏好" />

      <Card>
        <Tabs
          defaultActiveKey="basic"
          items={[
            {
              key: 'basic',
              label: '基础设置',
              children: (
                <Form form={form} layout="vertical" style={{ maxWidth: 600 }}>
                  <Form.Item label="系统名称" initialValue="Leap One 项目管理系统">
                    <Input />
                  </Form.Item>
                  <Form.Item label="系统Logo URL">
                    <Input placeholder="输入Logo图片地址" />
                  </Form.Item>
                  <Divider>安全设置</Divider>
                  <Form.Item label="密码最小长度" initialValue={8}>
                    <Input type="number" min={6} max={32} />
                  </Form.Item>
                  <Form.Item label="会话超时时间（分钟）" initialValue={30}>
                    <Input type="number" min={5} max={1440} />
                  </Form.Item>
                  <Form.Item label="启用验证码登录" valuePropName="checked" initialValue={false}>
                    <Switch />
                  </Form.Item>
                  <Form.Item label="允许注册新账号" valuePropName="checked" initialValue={false}>
                    <Switch />
                  </Form.Item>
                  <Button type="primary" icon={<SaveOutlined />} onClick={handleSaveBasic}>
                    保存设置
                  </Button>
                </Form>
              ),
            },
            {
              key: 'notification',
              label: '通知设置',
              children: (
                <Form layout="vertical" style={{ maxWidth: 600 }}>
                  <Form.Item label="邮件通知" valuePropName="checked" initialValue={true}>
                    <Switch />
                  </Form.Item>
                  <Form.Item label="站内消息" valuePropName="checked" initialValue={true}>
                    <Switch />
                  </Form.Item>
                  <Form.Item label="任务到期提醒（提前天数）" initialValue={1}>
                    <Input type="number" min={0} max={7} addonAfter="天" />
                  </Form.Item>
                  <Button type="primary" icon={<SaveOutlined />} onClick={() => message.success('通知设置已保存')}>
                    保存设置
                  </Button>
                </Form>
              ),
            },
            {
              key: 'log',
              label: '操作日志',
              children: (
                <div style={{ color: '#999' }}>
                  <p>操作日志记录了系统中所有关键操作，包括：</p>
                  <ul>
                    <li>用户登录/登出</li>
                    <li>数据的增删改操作</li>
                    <li>权限变更</li>
                    <li>系统配置修改</li>
                  </ul>
                  <p style={{ marginTop: 16 }}>日志保留期限：90天</p>
                  <Button>查看完整日志</Button>
                </div>
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
}
