/**
 * 导出按钮组件
 *
 * 支持单格式直接导出 / 多格式下拉选择导出
 * 支持 CSV / Excel / PDF 三种导出格式
 */

import { Button, Dropdown } from 'antd';
import { DownloadOutlined, FileExcelOutlined, FilePdfOutlined } from '@ant-design/icons';

// ── 类型定义 ─────────────────────────────────────────────────

type ExportFormat = 'csv' | 'excel' | 'pdf';

interface ExportButtonProps {
  /** 可选导出格式列表 */
  formats?: ExportFormat[];
  /** 导出回调 */
  onExport: (format: ExportFormat) => void;
  /** 加载状态 */
  loading?: boolean;
}

// ── 格式配置映射表 ─────────────────────────────────────────────

const FORMAT_CONFIG: Record<ExportFormat, { label: string; icon: React.ReactNode }> = Object.freeze({
  csv:   { label: 'CSV格式',  icon: <DownloadOutlined /> },
  excel: { label: 'Excel格式', icon: <FileExcelOutlined /> },
  pdf:   { label: 'PDF格式',   icon: <FilePdfOutlined /> },
});

const DEFAULT_FORMATS: ExportFormat[] = ['excel'];

// ── 组件实现 ─────────────────────────────────────────────────

const ExportButton = ({
  formats = DEFAULT_FORMATS,
  onExport,
  loading = false,
}: ExportButtonProps) => {
  // 单格式模式：直接渲染按钮
  if (formats.length === 1) {
    const singleFmt = formats[0] as ExportFormat;
    const { label, icon } = FORMAT_CONFIG[singleFmt];

    return (
      <Button icon={icon} loading={loading} onClick={() => onExport(singleFmt)}>
        导出{label}
      </Button>
    );
  }

  // 多格式模式：下拉菜单选择
  const menuItems = formats.map((fmt) => ({
    key: fmt,
    label: FORMAT_CONFIG[fmt].label,
    icon: FORMAT_CONFIG[fmt].icon,
    onClick: () => onExport(fmt),
  }));

  return (
    <Dropdown menu={{ items: menuItems }} placement="bottomRight">
      <Button icon={<DownloadOutlined />} loading={loading}>
        导出
      </Button>
    </Dropdown>
  );
};

export default ExportButton;
