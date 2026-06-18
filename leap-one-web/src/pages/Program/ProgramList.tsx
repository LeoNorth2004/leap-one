/** 项目集列表页面 */

import { useState } from 'react';
import { Row, Col, Card, Tag, Statistic, Button, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import PageHeader from '@/components/Common/PageHeader';
import SearchBar from '@/components/Common/SearchBar';

interface ProgramItem {
  id: number;
  name: string;
  description: string;
  status: 'active' | 'completed' | 'archived';
  projectCount: number;
  memberCount: number;
  pmName: string;
  progress: number;
}

/** 模拟数据 - TODO: 替换为API */
const mockPrograms: ProgramItem[] = [
  { id: 1, name: '企业数字化平台', description: '公司核心业务数字化转型项目集合', status: 'active', projectCount: 5, memberCount: 25, pmName: '张三', progress: 65 },
  { id: 2, name: '移动端产品线', description: 'iOS/Android/H5多端产品开发', status: 'active', projectCount: 3, memberCount: 18, pmName: '李四', progress: 40 },
  { id: 3, name: '数据中台建设', description: '数据采集、治理、分析平台', status: 'active', projectCount: 4, memberCount: 15, pmName: '王五', progress: 80 },
  { id: 4, name: '旧系统迁移', description: '遗留系统向新架构迁移', status: 'completed', projectCount: 6, memberCount: 30, pmName: '赵六', progress: 100 },
];

export default function ProgramList() {
  const [searchValue, setSearchValue] = useState('');
  return (
    <div>
      <PageHeader
        title="项目集管理"
        subtitle="管理和跟踪多个关联项目的整体进度"
        extra={
          <Button type="primary" icon={<PlusOutlined />}>
            创建项目集
          </Button>
        }
      />

      <SearchBar value={searchValue || ''} onChange={setSearchValue} onSearch={() => {}} placeholder="搜索项目集名称..." />

      <Row gutter={[16, 16]}>
        {mockPrograms.map((program) => (
          <Col xs={24} lg={12} key={program.id}>
            <Card hoverable className="program-card" style={{ borderRadius: 12 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div>
                  <h3 style={{ margin: 0 }}>{program.name}</h3>
                  <p style={{ color: '#666', fontSize: 13, marginTop: 4 }}>{program.description}</p>
                </div>
                <Tag color={program.status === 'active' ? 'processing' : program.status === 'completed' ? 'success' : 'default'}>
                  {program.status === 'active' ? '进行中' : program.status === 'completed' ? '已完成' : '已归档'}
                </Tag>
              </div>

              <Row gutter={16} style={{ marginTop: 16 }}>
                <Col span={8}>
                  <Statistic title="项目数" value={program.projectCount} suffix="个" />
                </Col>
                <Col span={8}>
                  <Statistic title="成员数" value={program.memberCount} suffix="人" />
                </Col>
                <Col span={8}>
                  <Statistic title="进度" value={program.progress} suffix="%" />
                </Col>
              </Row>

              <div style={{ marginTop: 12, paddingTop: 12, borderTop: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ color: '#999', fontSize: 13 }}>负责人：{program.pmName}</span>
                <Space>
                  <Button type="link" size="small">查看详情</Button>
                  <Button type="link" size="small">编辑</Button>
                </Space>
              </div>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
