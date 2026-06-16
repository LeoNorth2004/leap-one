# Leap One 安装与部署指南

> **版本**: v1.0.0
> **状态**: 草稿
> **最后更新**: 2026-06-08

---

## 1. 部署架构概览

### 1.1 环境分层
| 环境 | 用途 | 规模建议 |
|------|------|----------|
| development | 本地开发调试 | 单机Docker Compose |
| staging | 预发布验证 | 小规模K8s集群 |
| production | 生产运行 | 高可用K8s集群 |

### 1.2 硬件要求
<!-- TODO: 开发环境/生产环境的CPU、内存、磁盘最低/推荐配置 -->

### 1.3 软件依赖
| 软件 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Docker | 24.0+ | 最新稳定版 |
| Docker Compose | v2.20+ | 最新稳定版 |
| Go | 1.23+ | 1.23.x |
| Node.js | 20+ | 20 LTS |
| PostgreSQL | 16+ | 16.x |
| Redis | 7.0+ | 7.0.x |
| Nginx | 1.24+ | 1.24.x |

---

## 2. 快速开始（开发环境）

### 2.1 克隆仓库
```bash
git clone <repository-url>
cd leap-one
```

### 2.2 环境配置
```bash
cp .env.example .env
# 编辑 .env 文件，修改密码等敏感信息
```

### 2.3 启动基础设施
```bash
make run-dev
# 或手动执行：
docker compose up -d postgres-user redis
```

### 2.4 启动应用服务
<!-- TODO: 各服务的本地启动方式（Air热重载或直接运行） -->

### 2.5 启动前端
```bash
make frontend-install
make frontend-dev
```

### 2.6 验证安装
<!-- TODO: 健康检查端点、默认账号登录验证 -->

---

## 3. Docker Compose 部署

### 3.1 完整启动所有服务
```bash
docker compose up -d
```

### 3.2 分步启动（推荐）
```bash
# 第一步：启动基础设施
docker compose up -d postgres-* redis

# 第二步：等待数据库健康后启动应用
docker compose up -d gateway-service service-*

# 第三步：启动Nginx
docker compose up -d nginx
```

### 3.3 常用运维命令
```bash
# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f [service-name]

# 重启单个服务
docker compose restart service-user-org

# 进入容器调试
docker compose exec postgres-user psql -U leapone -d user_db
```

---

## 4. Kubernetes 部署

### 4.1 前置准备
<!-- TODO: K8s集群要求、kubectl配置、镜像仓库配置 -->

### 4.2 Helm Chart 部署
```bash
# TODO: Helm chart 安装命令
helm install leapone ./helm/leapone -f values-production.yaml
```

### 4.3 K8s 资源清单
<!-- TODO: Deployment/Service/ConfigMap/Secret/PVC/Ingress YAML示例 -->

### 4.4 滚动更新策略
<!-- TODO: RollingUpdate配置、健康探针配置 -->

---

## 5. 生产环境配置

### 5.1 安全加固
<!-- TODO: TLS证书配置、防火墙规则、网络隔离 -->

### 5.2 性能调优
<!-- TODO: PostgreSQL连接池、Redis内存、Gin worker数量、Nginx worker_processes -->

### 5.3 高可用配置
<!-- TODO: PostgreSQL主从复制、Redis Sentinel/Cluster、多副本部署 -->

### 5.4 备份与恢复
<!-- TODO: 自动化备份脚本、恢复演练流程 -->

---

## 6. 监控与运维

### 6.1 健康检查端点
| 端点 | 说明 |
|------|------|
| `/healthz` | 网关健康检查 |
| `/readyz` | 就绪探针 |
| `/livez` | 存活探针 |

### 6.2 日志收集
<!-- TODO: 日志格式、收集方案(Loki/ELK)、日志级别配置 -->

### 6.3 告警配置
<!-- TODO: 关键告警规则、通知渠道(PagerDuty/钉钉/企微) -->

---

## 7. 故障排查

### 7.1 常见问题
<!-- TODO: FAQ：连接失败、端口冲突、内存不足等常见问题及解决方案 -->

### 7.2 诊断工具
<!-- TODO: 诊断命令速查表 -->

---

## 附录

- A. 环境变量完整参考
- B. 端口映射一览表
- C. 升级指南
- D. 回滚操作手册
