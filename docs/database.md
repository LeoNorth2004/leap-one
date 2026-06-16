# Leap One 数据库设计文档

> **版本**: v1.0.0
> **状态**: 草稿
> **最后更新**: 2026-06-08

---

## 1. 数据库架构概览

### 1.1 分库策略
Leap One 采用**每个微服务独占一个数据库**的策略，共15个PostgreSQL实例：

| 服务名 | 数据库名 | 端口 | 主要用途 |
|--------|----------|------|----------|
| service-user-org | user_db | 5432 | 用户、组织、角色、权限 |
| service-portfolio | portfolio_db | 5433 | 投资组合、战略对齐 |
| service-project | project_db | 5434 | 项目信息、里程碑、阶段 |
| service-task | task_db | 5435 | 任务、子任务、工时 |
| service-requirement | requirement_db | 5436 | 需求、需求条目、变更记录 |
| service-quality | quality_db | 5437 | 缺陷、测试用例、审计报告 |
| service-devops | devops_db | 5438 | 流水线、部署记录、环境 |
| service-document | document_db | 5439 | 文档元数据、版本记录 |
| service-kanban | kanban_db | 5440 | 看板、列、卡片、WIP限制 |
| service-bi | bi_db | 5441 | 报表配置、仪表盘、数据快照 |
| service-ai | ai_db | 5442 | AI会话、知识库、向量索引 |
| service-notification | notification_db | 5443 | 消息、订阅、推送记录 |
| service-search | search_db | 5444 | 搜索索引配置、搜索历史 |
| service-config | config_db | 5445 | 全局配置、功能开关、租户设置 |

### 1.2 连接池配置
<!-- TODO: 各服务的数据库连接池参数（最大连接数、最小空闲、超时时间） -->

### 1.3 备份策略
<!-- TODO: pg_dump定时备份、Point-in-Time Recovery配置 -->

---

## 2. 核心数据模型

### 2.1 用户与组织 (user_db)

#### users 表
```sql
-- TODO: 定义users表结构
-- 字段：id, username, email, password_hash, avatar, status, created_at, updated_at
```

#### organizations 表
```sql
-- TODO: 定义organizations表结构
```

#### roles 表 & permissions 表
```sql
-- TODO: RBAC相关表结构
```

### 2.2 项目管理 (project_db)

#### projects 表
```sql
-- TODO: 定义projects表结构
```

#### milestones 表
```sql
-- TODO: 定义milestones表结构
```

### 2.3 任务管理 (task_db)

#### tasks 表
```sql
-- TODO: 定义tasks表结构，包含优先级、状态、关联关系
```

#### task_dependencies 表
```sql
-- TODO: 任务依赖关系表
```

#### time_entries 表
```sql
-- TODO: 工时记录表
```

### 2.4 需求管理 (requirement_db)

#### requirements 表
```sql
-- TODO: 需求主表，包含类型、优先级、状态
```

#### requirement_items 表
```sql
-- TODO: 需求条目表
```

### 2.5 质量管理 (quality_db)

#### defects 表
```sql
-- TODO: 缺陷表，包含严重程度、复现步骤
```

#### test_cases 表
```sql
-- TODO: 测试用例表
```

### 2.6 其他服务数据模型
<!-- TODO: DevOps/Document/Kanban/BI/AI/Notification/Search/Config 的表结构 -->

---

## 3. 索引设计

### 3.1 主键与外键
<!-- TODO: 各表的主键策略（UUID自增）、外键约束说明 -->

### 3.2 查询优化索引
<!-- TODO: 高频查询场景的复合索引设计 -->

### 3.3 全文搜索索引
<!-- TODO: PostgreSQL tsvector全文索引配置 -->

---

## 4. 数据迁移

### 4.1 迁移工具
<!-- TODO: golang-migrate / Goose / Flyway 选型及使用方法 -->

### 4.2 迁移脚本规范
<!-- TODO: 命名规范、回滚脚本要求 -->

### 4.3 已有迁移列表
<!-- TODO: 记录已执行的迁移文件 -->

---

## 5. 数据安全

### 5.1 敏感数据加密
<!-- TODO: 密码哈希(bcrypt)、PII数据加密存储 -->

### 5.2 行级安全策略 (RLS)
<!-- TODO: PostgreSQL RLS实现多租户数据隔离 -->

### 5.3 SQL注入防护
<!-- TODO: 参数化查询、ORM使用规范 -->

---

## 附录

- A. 完整ER图（Mermaid）
- B. 所有DDL语句汇总
- C. 种子数据SQL
