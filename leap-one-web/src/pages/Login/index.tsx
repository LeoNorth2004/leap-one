/**
 * 登录页面
 *
 * 设计：
 * - 居中卡片布局，左侧品牌展示区（Leap One logo + 标语），右侧登录表单
 * - 表单字段：用户名、密码、记住我
 * - Ant Design Form 表单验证 + 回车提交
 * - 响应式设计（移动端隐藏左侧品牌区）
 */

import { useCallback } from 'react';
import { Form, Input, Button, Checkbox, Card } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import styles from './index.module.less';

// ── 表单验证规则 ─────────────────────────────────────────────

const USERNAME_RULES = [
  { required: true, message: '请输入用户名' },
  { min: 2, message: '用户名至少 2 个字符' },
];

const PASSWORD_RULES = [
  { required: true, message: '请输入密码' },
  { min: 6, message: '密码至少 6 个字符' },
];

// ── 粒子数量 ─────────────────────────────────────────────────

const PARTICLE_COUNT = 8;

// ── 页面组件 ─────────────────────────────────────────────────

const LoginPage = () => {
  const { login, isLoading } = useAuth();
  const [form] = Form.useForm();

  // 提交登录
  const handleSubmit = useCallback(
    async (values: { username: string; password: string; remember?: boolean }) => {
      try {
        await login({
          username: values.username,
          password: values.password,
          remember: values.remember,
        });
      } catch {
        // 错误由 Axios 拦截器统一处理
      }
    },
    [login]
  );

  // 回车键提交
  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Enter') {
        form.submit();
      }
    },
    [form]
  );

  return (
    <div className={styles.loginPage}>
      {/* 背景装饰 */}
      <div className={styles.loginBg} />

      {/* 粒子效果 */}
      <div className={styles.loginParticles}>
        {Array.from({ length: PARTICLE_COUNT }, (_, idx) => (
          <span key={idx} className={styles.particle} />
        ))}
      </div>

      {/* 主卡片容器 */}
      <div className={styles.loginContainer}>
        {/* ═══ 左侧品牌区 ════════════════════════════════════ */}
        <div className={styles.brandArea}>
          <div className={styles.brandContent}>
            <span className={styles.brandIcon}>🚀</span>
            <h1 className={styles.brandTitle}>Leap One</h1>
            <p className={styles.brandSubtitle}>企业级项目管理系统</p>
            <p className={styles.brandDesc}>
              高效协作 · 敏捷交付 · 数据驱动
            </p>

            {/* 特性列表 */}
            <div className={styles.featureList}>
              {[
                '需求全生命周期管理',
                '敏捷 / 看板双模式支持',
                '自动化测试与质量追踪',
                'BI 数据大屏与智能分析',
              ].map((feature) => (
                <div key={feature} className={styles.featureItem}>
                  <span className={styles.featureDot} />
                  <span>{feature}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* ═══ 右侧登录表单区 ══════════════════════════════ */}
        <Card className={styles.loginCard}>
          <div className={styles.formHeader}>
            <h2 className={styles.formTitle}>欢迎回来</h2>
            <p className={styles.formDesc}>请登录您的账户以继续</p>
          </div>

          <Form
            form={form}
            onFinish={handleSubmit}
            size="large"
            autoComplete="off"
            initialValues={{ remember: true }}
            onKeyDown={handleKeyDown}
            requiredMark={false}
          >
            <Form.Item name="username" rules={USERNAME_RULES}>
              <Input
                prefix={<UserOutlined className={styles.inputIcon} />}
                placeholder="请输入用户名"
              />
            </Form.Item>

            <Form.Item name="password" rules={PASSWORD_RULES}>
              <Input.Password
                prefix={<LockOutlined className={styles.inputIcon} />}
                placeholder="请输入密码"
              />
            </Form.Item>

            <Form.Item name="remember" valuePropName="checked">
              <Checkbox className={styles.rememberCheckbox}>记住我</Checkbox>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                block
                loading={isLoading}
                size="large"
                className={styles.submitBtn}
              >
                {isLoading ? '登录中...' : '登 录'}
              </Button>
            </Form.Item>
          </Form>

          <p className={styles.footerText}>
            Leap One &copy; {new Date().getFullYear()} — 企业级项目管理平台
          </p>
        </Card>
      </div>
    </div>
  );
};

export default LoginPage;
