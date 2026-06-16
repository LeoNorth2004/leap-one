.PHONY: help run-dev build test clean docker-build docker-up docker-down docker-logs lint fmt vet frontend-install frontend-dev frontend-build frontend-lint

help: ## 显示帮助信息
	@echo "Leap One 构建命令"
	@echo "=================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run-dev: ## 启动开发环境（Docker Compose + Air热重载）
	docker compose up -d postgres-user postgres-portfolio postgres-project postgres-task \
		postgres-requirement postgres-quality postgres-devops postgres-document \
		postgres-kanban postgres-bi postgres-ai postgres-notification \
		postgres-search postgres-config redis
	@echo "开发环境基础设施已启动"

build: ## 编译所有服务
	@for dir in $$(ls -d service-*/); do \
		echo "Building $$dir..."; \
		cd $$dir && go build -o ../bin/$$(basename $$dir) cmd/main.go && cd ..; \
	done

test: ## 运行所有测试
	@for dir in $$(ls -d service-*/); do \
		echo "Testing $$dir..."; \
		cd $$dir && go test -v -race ./... && cd ..; \
	done

clean: ## 清理构建产物
	rm -rf bin/

docker-build: ## 构建所有Docker镜像
	docker compose build

docker-up: ## 启动所有服务
	docker compose up -d

docker-down: ## 停止所有服务
	docker compose down

docker-logs: ## 查看日志
	docker compose logs -f

lint: ## 代码检查
	@golangci-lint run ./...

fmt: ## 格式化代码
	@gofmt -l -w .

vet: ## 静态分析
	@go vet ./...

frontend-install: ## 安装前端依赖
	cd leap-one-web && npm install

frontend-dev: ## 启动前端开发服务器
	cd leap-one-web && npm run dev

frontend-build: ## 构建前端
	cd leap-one-web && npm run build

frontend-lint: ## 前端代码检查
	cd leap-one-web && npm run lint
