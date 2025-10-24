.PHONY: help install dev build clean docker-up docker-down test

help:
	@echo "AI-Motion 项目命令"
	@echo "===================="
	@echo "install         - 安装所有依赖"
	@echo "dev             - 启动开发环境"
	@echo "build           - 编译项目"
	@echo "clean           - 清理编译产物"
	@echo "docker-up       - 启动 Docker 容器"
	@echo "docker-down     - 停止 Docker 容器"
	@echo "test            - 运行测试"

install:
	@echo "安装后端依赖..."
	cd backend && go mod download
	@echo "安装前端依赖..."
	cd frontend && npm install

dev:
	@echo "启动后端服务..."
	cd backend && go run cmd/main.go &
	@echo "启动前端服务..."
	cd frontend && npm run dev

build:
	@echo "编译后端..."
	cd backend && go build -o bin/ai-motion cmd/main.go
	@echo "编译前端..."
	cd frontend && npm run build

clean:
	@echo "清理编译产物..."
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf storage/temp

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

test:
	@echo "运行后端测试..."
	cd backend && go test ./...
	@echo "运行前端测试..."
	cd frontend && npm test
