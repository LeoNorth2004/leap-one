# Leap One API 接口文档

> **版本**: v1.0.0
> **Base URL**: `http://localhost:8080/api/v1`
> **认证方式**: Bearer Token (JWT)
> **状态**: 草稿
> **最后更新**: 2026-06-08

---

## 1. API 概述

### 1.1 通用约定
- **协议**: HTTPS (生产环境) / HTTP (开发环境)
- **内容类型**: `application/json`
- **字符编码**: UTF-8
- **日期格式**: ISO 8601 (`YYYY-MM-DDTHH:mm:ssZ`)

### 1.2 统一响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1717843200
}
```

### 1.3 错误码规范
| 错误码 | 含义 | HTTP状态码 |
|--------|------|------------|
| 0 | 成功 | 200 |
| 10001 | 参数校验失败 | 400 |
| 10002 | 未授权 | 401 |
| 10003 | 权限不足 | 403 |
| 10004 | 资源不存在 | 404 |
| 10005 | 重复操作 | 409 |
| 20001 | 服务器内部错误 | 500 |
| 20002 | 服务不可用 | 503 |

### 1.4 分页规范
```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

---

## 2. 认证接口 (`/auth`)

### 2.1 用户登录
```
POST /api/v1/auth/login
```
<!-- TODO: 请求体、响应示例 -->

### 2.2 用户登出
```
POST /api/v1/auth/logout
```

### 2.3 刷新Token
```
POST /api/v1/auth/refresh
```

### 2.4 修改密码
```
PUT /api/v1/auth/password
```

---

## 3. 用户与组织接口 (`/users`, `/orgs`)

### 3.1 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users | 用户列表（分页） |
| POST | /api/v1/users | 创建用户 |
| GET | /api/v1/users/:id | 用户详情 |
| PUT | /api/v1/users/:id | 更新用户 |
| DELETE | /api/v1/users/:id | 删除用户 |

### 3.2 组织管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/orgs | 组织列表 |
| POST | /api/v1/orgs | 创建组织 |
| GET | /api/v1/orgs/:id | 组织详情 |
| PUT | /api/v1/orgs/:id | 更新组织 |

### 3.3 角色权限
<!-- TODO: 角色 CRUD、权限分配、权限检查接口 -->

---

## 4. 项目管理接口 (`/projects`)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/projects | 项目列表 |
| POST | /api/v1/projects | 创建项目 |
| GET | /api/v1/projects/:id | 项目详情 |
| PUT | /api/v1/projects/:id | 更新项目 |
| DELETE | /api/v1/projects/:id | 归档项目 |
| GET | /api/v1/projects/:id/milestones | 里程碑列表 |
| POST | /api/v1/projects/:id/milestones | 创建里程碑 |

<!-- TODO: 各接口的请求/响应示例 -->

---

## 5. 任务管理接口 (`/tasks`)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/tasks | 任务列表（支持筛选/排序） |
| POST | /api/v1/tasks | 创建任务 |
| GET | /api/v1/tasks/:id | 任务详情 |
| PUT | /api/v1/tasks/:id | 更新任务 |
| DELETE | /api/v1/tasks/:id | 删除任务 |
| POST | /api/v1/tasks/:id/assign | 分配任务 |
| POST | /api/v1/tasks/:id/status | 更改状态 |
| GET | /api/v1/tasks/:id/time-entries | 工时记录 |
| POST | /api/v1/tasks/:id/time-entries | 添加工时 |

---

## 6. 需求管理接口 (`/requirements`)
<!-- TODO: 需求CRUD、需求条目、变更记录接口 -->

---

## 7. 质量管理接口 (`/quality`)
<!-- TODO: 缺陷CRUD、测试用例、质量报告接口 -->

---

## 8. DevOps 接口 (`/devops`)
<!-- TODO: 流水线、部署、环境管理接口 -->

---

## 9. 文档管理接口 (`/documents`)
<!-- TODO: 文档上传/下载/删除、版本管理接口 -->
<!-- 注意：大文件上传可能需要分片上传接口 -->

---

## 10. 看板接口 (`/kanban`)
<!-- TODO: 看板/列/卡片/WIP CRUD接口 -->

---

## 11. BI 接口 (`/bi`)
<!-- TODO: 仪表盘、报表、数据导出接口 -->

---

## 12. AI 接口 (`/ai`)

### 12.1 对话接口
```
POST /api/v1/ai/chat
Content-Type: application/json

{
  "message": "帮我分析当前项目的风险",
  "session_id": "uuid",
  "stream": true
}
```
<!-- TODO: SSE流式响应格式说明 -->

### 12.2 智能推荐
<!-- TODO: 推荐、预测类接口 -->

---

## 13. 通知接口 (`/notifications`)
<!-- TODO: 消息列表、已读/未读、订阅管理接口 -->

---

## 14. 搜索接口 (`/search`)

### 14.1 全局搜索
```
GET /api/v1/search?q=keyword&type=all&page=1&page_size=20
```
<!-- TODO: 搜索结果聚合、高亮、过滤参数说明 -->

### 14.2 搜索建议
```
GET /api/v1/search/suggest?q=prefix
```

---

## 15. 配置接口 (`/config`)
<!-- TODO: 全局配置读取/更新、功能开关接口 -->

---

## 16. WebSocket 接口 (`/ws`)

### 16.1 连接
```
WS /ws?token=jwt_token
```

### 16.2 消息格式
<!-- TODO: 心跳、实时通知、AI流式输出的WebSocket消息协议 -->

---

## 附录

- A. Postman/Swagger 导入文件
- B. SDK 示例代码（Go/Python/JavaScript）
- C. 接口变更日志
