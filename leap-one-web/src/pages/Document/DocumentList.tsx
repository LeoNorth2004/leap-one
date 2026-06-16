/** 文档列表页面 */

import { Tree, Button, Space, Card, Modal, Input, message, Form, Select } from 'antd';
import { PlusOutlined, FileTextOutlined, FolderOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useState } from 'react';
import PageHeader from '@/components/Common/PageHeader';

interface DocNode {
  key: string;
  title: string;
  isFolder: boolean;
  children?: DocNode[];
}

/** 模拟文档树 */
const mockDocs: DocNode[] = [
  {
    key: 'root-1', title: '产品文档', isFolder: true, children: [
      { key: 'doc-1', title: '产品需求文档(PRD).md', isFolder: false },
      { key: 'doc-2', title: '产品路线图.md', isFolder: false },
      { key: 'folder-1', title: '版本发布说明', isFolder: true, children: [
        { key: 'doc-3', title: 'v2.0.0 发布说明.md', isFolder: false },
        { key: 'doc-4', title: 'v1.5.0 发布说明.md', isFolder: false },
      ]},
    ],
  },
  {
    key: 'root-2', title: '技术文档', isFolder: true, children: [
      { key: 'doc-5', title: 'API接口文档.md', isFolder: false },
      { key: 'doc-6', title: '数据库设计文档.md', isFolder: false },
      { key: 'doc-7', title: '部署运维手册.md', isFolder: false },
      { key: 'folder-2', title: '架构设计', isFolder: true, children: [
        { key: 'doc-8', title: '系统架构图.drawio', isFolder: false },
        { key: 'doc-9', title: '微服务拆分方案.md', isFolder: false },
      ]},
    ],
  },
  {
    key: 'root-3', title: '会议纪要', isFolder: true, children: [
      { key: 'doc-10', title: '2026-06-03 周会记录.md', isFolder: false },
      { key: 'doc-11', title: '2026-05-27 需求评审记录.md', isFolder: false },
    ],
  },
];

export default function DocumentList() {
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [createModalOpen, setCreateModalOpen] = useState(false);

  return (
    <div>
      <PageHeader
        title="文档中心"
        subtitle="项目知识库与文档协作管理"
        extra={
          <Space>
            <Button icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>新建文档</Button>
            <Button icon={<FolderOutlined />}>新建文件夹</Button>
          </Space>
        }
      />

      <Card>
        <Tree
          treeData={mockDocs.map((node) => ({
            ...node,
            icon: node.isFolder ? <FolderOutlined style={{ color: '#faad14' }} /> : <FileTextOutlined style={{ color: '#1677ff' }} />,
            title: (
              <span style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', minWidth: 300 }}>
                <span>{node.title as string}</span>
                {!node.isFolder && (
                  <Space size={4}>
                    <Button type="text" size="small" icon={<EyeOutlined />} onClick={(e) => { e.stopPropagation(); message.info('预览文档'); }} />
                    <Button type="text" size="small" icon={<EditOutlined />} onClick={(e) => e.stopPropagation()} />
                    <Button type="text" size="small" danger icon={<DeleteOutlined />} onClick={(e) => e.stopPropagation()} />
                  </Space>
                )}
              </span>
            ),
          }))}
          selectedKeys={selectedKeys}
          onSelect={(keys) => setSelectedKeys(keys as string[])}
          defaultExpandAll
          showIcon
          blockNode
        />
      </Card>

      <Modal open={createModalOpen} title="新建文档" onCancel={() => setCreateModalOpen(false)}
        onOk={() => { setCreateModalOpen(false); message.success('文档创建成功'); }}>
        <Form layout="vertical">
          <Form.Item label="文档名称" required><Input placeholder="请输入文档名称" /></Form.Item>
          <Form.Item label="父目录"><Input placeholder="选择父目录（可选）" disabled /></Form.Item>
          <Form.Item label="文档模板">
            <Select options={[{ label: '空白文档', value: 'blank' }, { label: 'Markdown模板', value: 'md' }, { label: 'API文档模板', value: 'api' }]} defaultValue="blank" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
