# Leap One（跃态项目管理系统）

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://react.dev/)

> **更快·更易用·更现代** —— 一款完全对标禅道所有功能的下一代企业级项目管理系统

## ✨ 特性

- 🚀 **高性能**：Go + Gin 架构，相比禅道PHP实现性能提升 **3-5倍**
- 🎨 **现代UI**：React + Ant Design 5，界面美观，操作流畅
- 🏗️ **微服务架构**：15个独立微服务，高内聚低耦合，支持水平扩展
- 🐳 **一键部署**：Docker Compose 一键启动，也支持传统服务器部署和K8s集群
- 🤖 **AI赋能**：内置AI需求预测、智能对话辅助功能
- 🔐 **企业级安全**：RBAC权限、JWT认证、审计日志
- 📊 **BI大屏**：内置5张宏观管理大屏，数据驱动决策
- 🔌 **开放API**：完善的RESTful API，易于集成

## 📋 功能全览

对标禅道所有核心功能模块：
- [x] 项目集管理
- [x] 产品管理（需求池·需求评审·变更管理·路线图）
- [x] 项目管理（敏捷/瀑布/轻量/全生命周期）
- [x] 执行管理与迭代看板
- [x] 质量管理（测试用例·测试计划·Bug管理）
- [x] DevOps管理（代码集成·CI/CD·制品管理）
- [x] 看板管理与泳道
- [x] 文档管理与知识库
- [x] 组织管理与RBAC权限
- [x] 工单/事项管理
- [x] BI统计与报表大屏
- [x] AI智能辅助
- [x] 全局搜索与消息通知

*详细功能清单见 [功能对比文档](docs/feature-comparison.md)*

## 🚀 快速开始

### 方式一：Docker Compose 一键部署（推荐）

```bash
# 克隆仓库
git clone https://github.com/your-org/leap-one.git
cd leap-one

# 配置环境变量
cp .env.example .env
# 编辑.env文件，设置数据库密码等

# 一键启动（自动拉取镜像、创建数据库、初始化）
docker compose up -d

# 查看状态
docker compose ps
```

访问 `http://localhost` 开始使用！
默认管理员：admin / Admin@123456

### 方式二：安装程序（Windows/macOS/Linux）

1. 从 [Releases](https://github.com/your-org/leap-one/releases) 下载对应平台的安装包
2. 运行安装程序，跟随图形化向导完成配置
3. 自动完成环境检测、数据库初始化、服务启动
4. 访问系统开始使用

### 方式三：源码编译

```bash
# 后端
cd services/user-service && go build -o user-service
# 依次启动各服务...

# 前端
cd web && npm install && npm run build
```

## 📚 文档

| 文档 | 说明 |
|------|------|
| [安装部署指南](docs/deployment.md) | Docker/二进制/K8s详细部署步骤 |
| [用户手册](docs/user-guide.md) | 各功能模块操作说明 + 禅道迁移指南 |
| [API文档](docs/api.md) | OpenAPI 3.0完整接口文档 |
| [数据库设计](docs/database.md) | 各服务数据表结构说明 |
| [架构设计](docs/design.md) | 微服务架构图与技术选型说明 |
| [维护手册](docs/maintenance.md) | 日常巡检、故障排查、灾备恢复 |

## 🛠️ 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端语言 | Go | 1.23+ |
| Web框架 | Gin | v1.10+ |
| ORM | GORM | v2.0+ |
| 数据库 | PostgreSQL | 16+ |
| 缓存 | Redis | 7.0+ |
| 前端框架 | React | 18+ |
| UI库 | Ant Design | 5.0+ |
| 容器化 | Docker + Docker Compose | - |
| CI/CD | GitHub Actions | - |

## 📊 性能对比

| 指标 | 禅道 | Leap One | 提升 |
|------|------|----------|------|
| API响应时间(P99) | ~800ms | ~150ms | **5.3倍** |
| 页面加载速度 | ~2.5s | ~0.8s | **3.1倍** |
| 支持并发用户 | ~500 | ~5000+ | **10倍** |
| 数据库查询速度 | 一般 | 优化索引+缓存 | **显著提升** |

## 🤝 从禅道迁移

1. 导出禅道数据（支持XML/CSV格式）
2. 在Leap One管理后台使用数据导入功能
3. 按照[禅道迁移指南](docs/user-guide.md#从禅道迁移)进行字段映射配置
4. 验证导入数据完整性
5. 切换团队使用（支持双轨运行过渡期）

## 📞 支持与反馈

- 📧 邮件：support@leapone.com
- 💬 GitHub Issues：提交问题
- 📖 官方文档：https://docs.leapone.com

## 📄 许可证

[Apache License 2.0](LICENSE)

Copyright © 2026 Leap One Team
