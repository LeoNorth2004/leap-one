# ============================================
# Leap One AI服务 - Docker多阶段构建
# ============================================

# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-ai cmd/main.go

# 运行阶段
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /service-ai .
EXPOSE 8011
CMD ["./service-ai"]
