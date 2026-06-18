/**
 * 搜索栏组件
 *
 * 提供搜索输入框、重置按钮、额外操作区的组合布局
 */

import { Input, Button, Space } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';

// ── 类型定义 ─────────────────────────────────────────────────

interface SearchBarProps {
  /** 搜索关键词值 */
  value: string;
  /** 搜索回调（按回车或点击搜索图标触发） */
  onSearch: (value: string) => void;
  /** 输入变化回调 */
  onChange?: (value: string) => void;
  /** 输入框占位文本 */
  placeholder?: string;
  /** 是否显示重置按钮 */
  showReset?: boolean;
  /** 重置按钮点击回调 */
  onReset?: () => void;
  /** 额外操作区域（显示在搜索框右侧） */
  extra?: React.ReactNode;
}

// ── 默认配置 ─────────────────────────────────────────────────

const DEFAULT_PLACEHOLDER = '请输入关键词搜索...';

// ── 组件实现 ─────────────────────────────────────────────────

const SearchBar = ({
  value,
  onSearch,
  onChange,
  placeholder = DEFAULT_PLACEHOLDER,
  showReset,
  onReset,
  extra,
}: SearchBarProps) => (
  <div className="search-bar" style={{ marginBottom: 16 }}>
    <Space wrap>
      <Input.Search
        allowClear
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        onSearch={onSearch}
        style={{ width: 300 }}
      />
      {showReset && (
        <Button icon={<ReloadOutlined />} onClick={onReset}>
          重置
        </Button>
      )}
      {extra}
    </Space>
  </div>
);

export default SearchBar;
