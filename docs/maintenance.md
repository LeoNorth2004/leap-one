# Leap One 运维维护手册

> **版本**: v1.0.0
> **状态**: 草稿
> **最后更新**: 2026-06-08

---

## 1. 日常运维

### 1.1 每日巡检清单
- [ ] 检查所有服务运行状态 (`docker compose ps`)
- [ ] 检查各服务健康检查端点
- [ ] 检查磁盘空间使用率
- [ ] 检查数据库连接数
- [ ] 检查Redis内存使用情况
- [ ] 检查错误日志是否有异常

### 1.2 每周维护任务
- [ ] 清理过期日志文件
- [ ] 检查数据库慢查询日志
- [ ] 更新依赖包安全补丁
- [ ] 备份验证（恢复演练）

### 1.3 每月维护任务
- [ ] 性能基线回顾
- [ ] 容量规划评估
- [ ] 安全漏洞扫描
- [ ] 备份完整性校验

---

## 2. 服务管理

### 2.1 服务启停
```bash
# 启动所有服务
docker compose up -d

# 停止所有服务
docker compose down

# 重启指定服务
docker compose restart service-user-org

# 扩容服务实例
docker compose up -d --scale service-task=3
```

### 2.2 服务健康检查
```bash
# 检查网关健康
curl http://localhost:8080/healthz

# 检查各微服务健康
curl http://localhost:8001/healthz
curl http://localhost:8002/healthz
# ... 以此类推
```

### 2.3 日志管理
```bash
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f service-project

# 查看最近100行错误日志
docker compose logs --tail=100 service-task 2>&1 | grep -i error

# 导出日志到文件
docker compose logs > logs_$(date +%Y%m%d).log
```

---

## 3. 数据库运维

### 3.1 PostgreSQL 维护

#### 连接数据库
```bash
# 通过Docker进入
docker compose exec postgres-user psql -U leapone -d user_db

# 从宿主机连接
psql -h localhost -p 5432 -U leapone -d user_db
```

#### 常用运维SQL
```sql
-- 查看活跃连接
SELECT * FROM pg_stat_activity WHERE state = 'active';

-- 查看表大小
SELECT relname, pg_size_pretty(pg_total_relation_size(relid)) 
FROM pg_stat_user_tables ORDER BY pg_total_relation_size(relid) DESC;

-- 查看慢查询（需开启pg_stat_statements）
SELECT * FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 20;

-- VACUUM分析
VACUUM ANALYZE;
```

#### 备份与恢复
```bash
# 全量备份
docker compose exec postgres-user pg_dump -U leapone -d user_db > backup_user_$(date +%Y%m%d).sql

# 仅备份数据（无schema）
docker compose exec postgres-user pg_dump -U leapone -d user_data-only > backup_data.sql

# 恢复
cat backup.sql | docker compose exec -T postgres-user psql -U leapone -d user_db
```

### 3.2 Redis 维护

#### 连接Redis
```bash
docker compose exec redis redis-cli -a ${REDIS_PASSWORD}
```

#### 常用运维命令
```bash
# 查看内存使用
INFO memory

# 查看Key数量
DBSIZE

# 查看大Key
redis-cli --bigkeys

# 清理过期Key（谨慎操作）
# FLUSHDB  # 清空当前库
# FLUSHALL # 清空所有库
```

---

## 4. Nginx 运维

### 4.1 配置重载
```bash
# 测试配置语法
docker compose exec nginx nginx -t

# 平滑重载（不中断连接）
docker compose exec nginx nginx -s reload
```

### 4.2 SSL证书更新
<!-- TODO: Let's Encrypt自动续期或手动更新流程 -->

### 4.3 访问日志分析
```bash
# Top 10 访问IP
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head -10

# Top 10 慢请求
awk '$NF > 5' /var/log/nginx/access.log | awk '{print $7}' | sort | uniq -c | sort -rn | head -10
```

---

## 5. 性能调优

### 5.1 PostgreSQL 调优参数
<!-- TODO: shared_buffers, work_mem, effective_cache_size等参数建议值 -->

### 5.2 Redis 调优参数
<!-- TODO: maxmemory-policy, tcp-backlog, timeout等参数建议值 -->

### 5.3 Go/Gin 调优
<!-- TODO: GOMAXPROCS、连接池大小、GC参数调整 -->

### 5.4 Nginx 调优
<!-- TODO: worker_connections、keepalive、gzip配置优化 -->

---

## 6. 故障处理

### 6.1 故障分级
| 等级 | 定义 | 响应时间 | 恢复时间 |
|------|------|----------|----------|
| P0 | 系统完全不可用 | 5分钟 | 1小时 |
| P1 | 核心功能受损 | 15分钟 | 4小时 |
| P2 | 非核心功能受损 | 2小时 | 24小时 |
| P3 | 优化建议 | 下个工作日 | - |

### 6.2 常见故障处理
<!-- TODO: 数据库连接耗尽、Redis宕机、OOM、磁盘满等常见故障的处理SOP -->

### 6.3 应急联系
<!-- TODO: 值班人员、厂商技术支持联系方式 -->

---

## 7. 安全运维

### 7.1 定期安全检查
- [ ] 密码轮换（每90天）
- [ ] 证书有效期检查
- [ ] 依赖漏洞扫描（每周）
- [ ] 访问权限审计（每月）

### 7.2 安全事件响应
<!-- TODO: 安全事件分类、处置流程、取证保留 -->

---

## 8. 升级维护

### 8.1 滚动升级流程
<!-- TODO: 零停机升级的操作步骤 -->

### 8.2 数据库Schema升级
<!-- TODO: 迁移脚本执行顺序、回滚方案 -->

### 8.3 版本兼容矩阵
<!-- TODO: 各服务之间的版本兼容性要求 -->

---

## 附录

- A. 运维命令速查表
- B. 监控面板地址
- C. 应急联系人名单
- D. 变更记录
