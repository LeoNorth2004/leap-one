/** 导出按钮组件 */

import { Button, Dropdown } from 'antd';
import { DownloadOutlined, FileExcelOutlined, FilePdfOutlined } from '@ant-design/icons';

type ExportFormat = 'csv' | 'excel' | 'pdf';

interface ExportButtonProps {
  /** 导出格式选项 */
  formats?: ExportFormat[];
  /** 导出回调 */
  onExport: (format: ExportFormat) => void;
  /** 加载状态 */
  loading?: boolean;
}

const formatLabels: Record<ExportFormat, { label: string; icon: React.ReactNode }> = {
  csv: { label: 'CSV格式', icon: <DownloadOutlined /> },
  excel: { label: 'Excel格式', icon: <FileExcelOutlined /> },
  pdf: { label: 'PDF格式', icon: <FilePdfOutlined /> },
};

export default function ExportButton({
  formats = ['excel'],
  onExport,
  loading = false,
}: ExportButtonProps) {
  if (formats.length === 1) {
    const fmt = formats[0];
    return (
      <Button icon={<DownloadOutlined />} loading={loading} onClick={() => onExport(fmt)}>
        导出{formatLabels[fmt].label}
      </Button>
    );
  }

  const items = formats.map((f: ExportFormat) => ({
    key: f,
    label: formatLabels[f].label,
    icon: formatLabels[f].icon,
    onClick: () => onExport(f),
  }));

  return (
    <Dropdown menu={{ items }} placement="bottomRight">
      <Button icon={<DownloadOutlined />} loading={loading}>
        导出
      </Button>
    </Dropdown>
  );
}
