# ============================================
# Leap One Web - Dockerfile
# 多阶段构建：Node构建 → Nginx托管
# ============================================

# 阶段1: Node.js 构建环境
FROM node:20-alpine AS builder

WORKDIR /app

# 先复制package文件，利用Docker缓存层
COPY package.json package-lock.json* ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# TypeScript类型检查 + Vite生产构建
RUN npm run build

# ============================================
# 阶段2: Nginx 生产环境
FROM nginx:alpine

# 移除默认Nginx配置
RUN rm /etc/nginx/conf.d/default.conf

# 复制自定义Nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制产物到Nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]
