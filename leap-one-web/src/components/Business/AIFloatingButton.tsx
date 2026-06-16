/** AI悬浮球组件 - 右下角固定定位 */

import { useState } from 'react';
import { FloatButton, Drawer, Button } from 'antd';
import {
  RobotOutlined,
  CloseOutlined,
  SendOutlined,
} from '@ant-design/icons';

export default function AIFloatingButton() {
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<Array<{ role: 'user' | 'ai'; content: string }>>([]);
  const [inputValue, setInputValue] = useState('');

  /** 发送消息 */
  const handleSend = () => {
    if (!inputValue.trim()) return;
    setMessages((prev) => [...prev, { role: 'user', content: inputValue }]);
    setInputValue('');

    // TODO: 调用AI对话API
    setTimeout(() => {
      setMessages((prev) => [
        ...prev,
        {
          role: 'ai',
          content: '我是Leap One AI助手，可以帮您分析需求、拆解任务、生成测试用例等。请问有什么可以帮助您的？',
        },
      ]);
    }, 800);
  };

  return (
    <>
      <FloatButton
        icon={<RobotOutlined />}
        tooltip="AI助手"
        onClick={() => setOpen(true)}
        style={{ right: 24, bottom: 24 }}
      />

      <Drawer
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <RobotOutlined style={{ color: '#1677ff' }} />
            <span>AI 智能助手</span>
          </div>
        }
        placement="right"
        width={380}
        open={open}
        onClose={() => setOpen(false)}
        closeIcon={<CloseOutlined />}
      >
        {/* 消息列表 */}
        <div className="ai-chat-messages" style={{ flex: 1, overflowY: 'auto', marginBottom: 16 }}>
          {messages.length === 0 && (
            <div style={{ textAlign: 'center', color: '#999', padding: '40px 0' }}>
              <RobotOutlined style={{ fontSize: 48, marginBottom: 16, display: 'block' }} />
              <p>您好！我是 Leap One AI 助手</p>
              <p style={{ fontSize: 12 }}>我可以帮您：</p>
              <ul style={{ fontSize: 12, textAlign: 'left', display: 'inline-block' }}>
                <li>智能拆解需求为任务</li>
                <li>自动生成测试用例</li>
                <li>分析项目风险</li>
                <li>推荐任务指派</li>
              </ul>
            </div>
          )}
          {messages.map((msg, idx) => (
            <div
              key={idx}
              style={{
                marginBottom: 12,
                display: 'flex',
                justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
              }}
            >
              <div
                style={{
                  maxWidth: '80%',
                  padding: '8px 14px',
                  borderRadius: msg.role === 'user'
                    ? '16px 4px 16px 16px'
                    : '4px 16px 16px 16px',
                  background: msg.role === 'user' ? '#1677ff' : '#f5f5f5',
                  color: msg.role === 'user' ? '#fff' : '#333',
                  fontSize: 14,
                  lineHeight: 1.6,
                }}
              >
                {msg.content}
              </div>
            </div>
          ))}
        </div>

        {/* 输入区域 */}
        <div className="ai-chat-input" style={{ display: 'flex', gap: 8 }}>
          <input
            className="ant-input"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') handleSend(); }}
            placeholder="输入您的问题..."
            style={{ flex: 1, borderRadius: 20, paddingLeft: 16 }}
          />
          <Button
            type="primary"
            shape="circle"
            icon={<SendOutlined />}
            onClick={handleSend}
            disabled={!inputValue?.trim()}
          />
        </div>
      </Drawer>
    </>
  );
}
