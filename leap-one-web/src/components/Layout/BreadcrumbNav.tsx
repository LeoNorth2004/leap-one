/** 面包屑导航组件 */

import { useLocation } from 'react-router-dom';
import { Breadcrumb } from 'antd';
import { HomeOutlined } from '@ant-design/icons';

/** 路由到面包文名称映射 */
const routeLabelMap: Record<string, string> = {
  '/': '工作台',
  '/user': '用户管理',
  '/user/list': '用户列表',
  '/org/department': '部门管理',
  '/org/role': '角色权限',
  '/program': '项目集',
  '/program/list': '项目集列表',
  '/product': '产品',
  '/product/list': '产品列表',
  '/product/roadmap': '产品路线图',
  '/project': '项目',
  '/project/list': '项目列表',
  '/requirement': '需求',
  '/requirement/list': '需求列表',
  '/task': '任务',
  '/task/list': '任务列表',
  '/quality': '质量',
  '/quality/testcase': '测试用例',
  '/quality/bug': 'Bug列表',
  '/quality/testplan': '测试计划',
  '/issue': '工单',
  '/issue/list': '工单列表',
  '/document': '文档',
  '/document/list': '文档列表',
  '/kanban': '看板',
  '/bi': 'BI大屏',
  '/bi/dashboard': '数据概览',
  '/settings': '系统设置',
  '/profile': '个人中心',
};

export default function BreadcrumbNav() {
  const location = useLocation();

  /** 解析路径生成面包屑项 */
  const breadcrumbItems = (() => {
    const paths = location.pathname.split('/').filter(Boolean);
    const items: Array<{ title: React.ReactNode }> = [{ title: <HomeOutlined /> }];

    let currentPath = '';
    for (const segment of paths) {
      currentPath += `/${segment}`;
      const label = routeLabelMap[currentPath] || segment;
      items.push({ title: label });
    }

    return items;
  })();

  return (
    <Breadcrumb items={breadcrumbItems} style={{ fontSize: 14 }} />
  );
}
