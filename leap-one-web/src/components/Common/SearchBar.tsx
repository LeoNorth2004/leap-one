/** 搜索栏组件 */

import { Input, Button, Space } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';

interface SearchBarProps {
  /** 搜索值 */
  value: string;
  /** 搜索回调 */
  onSearch: (value: string) => void;
  /** 值变化回调 */
  onChange?: (value: string) => void;
  /** 占位文本 */
  placeholder?: string;
  /** 是否显示重置按钮 */
  showReset?: boolean;
  /** 重置回调 */
  onReset?: () => void;
  /** 额外操作按钮 */
  extra?: React.ReactNode;
}

export default function SearchBar({
  value,
  onSearch,
  onChange,
  placeholder = '请输入关键词搜索...',
  showReset,
  onReset,
  extra,
}: SearchBarProps) {
  return (
    <div className="search-bar" style={{ marginBottom: 16 }}>
      <Space wrap>
        <Input.Search
          allowClear
          placeholder={placeholder}
          value={value}
          onChange={(e) => {
            onChange?.(e.target.value);
          }}
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
}
