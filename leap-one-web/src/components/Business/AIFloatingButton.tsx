import { useState, useRef, useEffect } from 'react';
import { FloatButton, Drawer } from 'antd';
import {
  RobotOutlined,
  CloseOutlined,
  SendOutlined,
} from '@ant-design/icons';
import styles from './AIFloatingButton.module.less';

export default function AIFloatingButton() {
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<Array<{ role: 'user' | 'ai'; content: string }>>([]);
  const [inputValue, setInputValue] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isTyping]);

  const handleSend = () => {
    if (!inputValue.trim()) return;
    setMessages((prev) => [...prev, { role: 'user', content: inputValue }]);
    setInputValue('');
    setIsTyping(true);

    setTimeout(() => {
      setMessages((prev) => [
        ...prev,
        {
          role: 'ai',
          content: '我是Leap One AI助手，可以帮您分析需求、拆解任务、生成测试用例等。请问有什么可以帮助您的？',
        },
      ]);
      setIsTyping(false);
    }, 1200);
  };

  return (
    <>
      <FloatButton
        icon={<RobotOutlined />}
        tooltip="AI助手"
        onClick={() => setOpen(true)}
        style={{
          right: 24,
          bottom: 24,
          width: 56,
          height: 56,
          borderRadius: '50%',
          background: 'linear-gradient(135deg, #1677ff 0%, #13c2c2 100%)',
          boxShadow: '0 4px 16px rgba(22, 119, 255, 0.4)',
        }}
      />

      <Drawer
        className={styles.aiDrawer}
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <RobotOutlined style={{ color: '#fff' }} />
            <span>AI 智能助手</span>
          </div>
        }
        placement="right"
        width={380}
        open={open}
        onClose={() => setOpen(false)}
        closeIcon={<CloseOutlined />}
      >
        <div className={styles.aiChatMessages}>
          {messages.length === 0 && (
            <div className={styles.aiWelcome}>
              <RobotOutlined className={styles.welcomeIcon} />
              <div className={styles.welcomeTitle}>您好！我是 Leap One AI 助手</div>
              <div className={styles.welcomeSubtitle}>我可以帮您：</div>
              <ul className={styles.featureList}>
                <li className={styles.featureItem}>智能拆解需求为任务</li>
                <li className={styles.featureItem}>自动生成测试用例</li>
                <li className={styles.featureItem}>分析项目风险</li>
                <li className={styles.featureItem}>推荐任务指派</li>
              </ul>
            </div>
          )}

          {messages.map((msg, idx) => (
            <div
              key={idx}
              className={`${styles.messageItem} ${styles[msg.role]}`}
            >
              <div className={`${styles.messageBubble} ${styles[msg.role]}`}>
                {msg.content}
              </div>
            </div>
          ))}

          {isTyping && (
            <div className={`${styles.messageItem} ${styles.ai}`}>
              <div className={styles.typingIndicator}>
                <span className={styles.typingDot} />
                <span className={styles.typingDot} />
                <span className={styles.typingDot} />
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        <div className={styles.aiChatInput}>
          <div className={styles.inputWrapper}>
            <input
              className={styles.chatInput}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleSend();
              }}
              placeholder="输入您的问题..."
              disabled={isTyping}
            />
          </div>
          <button
            className={styles.sendBtn}
            onClick={handleSend}
            disabled={!inputValue?.trim() || isTyping}
          >
            <SendOutlined />
          </button>
        </div>
      </Drawer>
    </>
  );
}