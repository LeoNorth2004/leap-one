# ============================================
# Leap One 项目管理服务 - Docker多阶段构建
# ============================================

# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-project cmd/main.go

# 运行阶段
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /service-project .
EXPOSE 8003
CMD ["./service-project"]
