/**
 * 确认弹窗组件
 *
 * 基于 Ant Design Modal 的二次封装
 * 支持自定义按钮文字、危险操作模式、加载状态等
 */

import { Modal } from 'antd';

// ── 类型定义 ─────────────────────────────────────────────────

interface ConfirmModalProps {
  /** 弹窗是否可见 */
  open: boolean;
  /** 弹窗标题 */
  title: string;
  /** 弹窗描述内容 */
  content: string;
  /** 确认按钮文字 */
  confirmText?: string;
  /** 取消按钮文字 */
  cancelText?: string;
  /** 是否为危险操作（确认按钮变红） */
  danger?: boolean;
  /** 确认回调 */
  onConfirm: () => void | Promise<void>;
  /** 取消/关闭回调 */
  onCancel: () => void;
  /** 确认按钮加载状态 */
  loading?: boolean;
}

// ── 默认值 ───────────────────────────────────────────────────

const DEFAULT_CONFIRM_TEXT = '确定';
const DEFAULT_CANCEL_TEXT = '取消';

// ── 组件实现 ─────────────────────────────────────────────────

const ConfirmModal = ({
  open,
  title,
  content,
  confirmText = DEFAULT_CONFIRM_TEXT,
  cancelText = DEFAULT_CANCEL_TEXT,
  danger = false,
  onConfirm,
  onCancel,
  loading = false,
}: ConfirmModalProps) => (
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

export default ConfirmModal;
