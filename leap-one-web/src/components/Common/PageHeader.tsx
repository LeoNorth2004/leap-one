/** 页面标题组件 */

import React from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

interface PageHeaderProps {
  /** 页面标题 */
  title: string;
  /** 副标题/描述 */
  subtitle?: string;
  /** 额外操作区域（右侧） */
  extra?: React.ReactNode;
}

export default function PageHeader({ title, subtitle, extra }: PageHeaderProps) {
  return (
    <div className="page-header" style={{ marginBottom: 16 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Title level={4} style={{ margin: 0 }}>
            {title}
          </Title>
          {subtitle && <p style={{ color: '#999', margin: '4px 0 0', fontSize: 13 }}>{subtitle}</p>}
        </div>
        {extra && <div>{extra}</div>}
      </div>
    </div>
  );
}
