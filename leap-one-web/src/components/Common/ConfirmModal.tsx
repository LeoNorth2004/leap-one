/** 确认弹窗组件 */

import { Modal } from 'antd';

interface ConfirmModalProps {
  /** 是否可见 */
  open: boolean;
  /** 标题 */
  title: string;
  /** 描述内容 */
  content: string;
  /** 确认文字 */
  confirmText?: string;
  /** 取消文字 */
  cancelText?: string;
  /** 是否危险操作（红色确认按钮） */
  danger?: boolean;
  /** 确认回调 */
  onConfirm: () => void | Promise<void>;
  /** 取消/关闭回调 */
  onCancel: () => void;
  /** 是否加载中 */
  loading?: boolean;
}

export default function ConfirmModal({
  open,
  title,
  content,
  confirmText = '确定',
  cancelText = '取消',
  danger = false,
  onConfirm,
  onCancel,
  loading = false,
}: ConfirmModalProps) {
  return (
    <Modal
      open={open}
      title={title}
      okText={confirmText}
      cancelText={cancelText}
      okType={danger ? 'danger' : undefined}
      confirmLoading={loading}
      onOk={onConfirm}
      onCancel={onCancel}
    >
      <p>{content}</p>
    </Modal>
  );
}
