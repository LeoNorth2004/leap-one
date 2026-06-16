/** 部门管理页面 */

import { useState } from 'react';
import { Tree, Button, Modal, Form, Input, message, Space, Card } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';

interface DepartmentNode {
  key: string;
  title: string;
  children?: DepartmentNode[];
}

/** 模拟部门树数据 - TODO: 替换为API */
const mockDepartments: DepartmentNode[] = [
  {
    key: '1',
    title: '总公司',
    children: [
      {
        key: '1-1',
        title: '研发中心',
        children: [
          { key: '1-1-1', title: '前端组' },
          { key: '1-1-2', title: '后端组' },
          { key: '1-1-3', title: '测试组' },
        ],
      },
      { key: '1-2', title: '产品部' },
      { key: '1-3', title: '设计部' },
      { key: '1-4', title: '运维部' },
    ],
  },
];

export default function Department() {
  const [treeData] = useState<DepartmentNode[]>(mockDepartments);
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <div>
      <PageHeader
        title="部门管理"
        subtitle="管理组织架构和部门层级"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新增部门
          </Button>
        }
      />

      <Card>
        <Tree
          treeData={treeData}
          defaultExpandAll
          showIcon
          blockNode
          titleRender={(node) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingRight: 16 }}>
              <span>{node.title as string}</span>
              <Space size="small">
                <Button type="text" size="small" icon={<EditOutlined />} />
                <Button type="text" size="small" danger icon={<DeleteOutlined />} />
              </Space>
            </div>
          )}
        />
      </Card>

      <Modal
        open={modalOpen}
        title="新增部门"
        onCancel={() => setModalOpen(false)}
        onOk={() => { setModalOpen(false); message.success('部门创建成功'); }}
      >
        <Form layout="vertical">
          <Form.Item label="上级部门" required>
            <Input placeholder="选择上级部门（可选）" disabled />
          </Form.Item>
          <Form.Item label="部门名称" required>
            <Input placeholder="请输入部门名称" />
          </Form.Item>
          <Form.Item label="负责人">
            <Input placeholder="请输入负责人" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
