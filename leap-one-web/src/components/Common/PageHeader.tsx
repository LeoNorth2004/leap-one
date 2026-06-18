/**
 * 页面标题组件
 *
 * 渲染页面标题、可选副标题，右侧支持放置操作按钮等额外内容
 */

import { Typography } from 'antd';

const { Title } = Typography;

// ── 类型定义 ─────────────────────────────────────────────────

interface PageHeaderProps {
  /** 页面主标题 */
  title: string;
  /** 副标题 / 页面描述文字 */
  subtitle?: string;
  /** 右侧额外操作区 */
  extra?: React.ReactNode;
}

// ── 样式常量 ─────────────────────────────────────────────────

const HEADER_STYLE = { marginBottom: 16 };
const LAYOUT_STYLE = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
};
const SUBTITLE_STYLE = { color: '#999', margin: '4px 0 0', fontSize: 13 };

// ── 组件实现 ─────────────────────────────────────────────────

const PageHeader = ({ title, subtitle, extra }: PageHeaderProps) => (
  <div className="page-header" style={HEADER_STYLE}>
    <div style={LAYOUT_STYLE}>
      <div>
        <Title level={4} style={{ margin: 0 }}>
          {title}
        </Title>
        {subtitle && <p style={SUBTITLE_STYLE}>{subtitle}</p>}
      </div>
      {extra && <div>{extra}</div>}
    </div>
  </div>
);

export default PageHeader;
