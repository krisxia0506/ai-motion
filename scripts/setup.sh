#!/bin/bash

# AI-Motion 项目初始化脚本

set -e

echo "================================"
echo "AI-Motion 项目初始化"
echo "================================"
echo ""

# 检查依赖
echo "检查系统依赖..."

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi
echo "✓ Go 版本: $(go version)"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 18+"
    exit 1
fi
echo "✓ Node.js 版本: $(node --version)"

# 检查 npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm 未安装"
    exit 1
fi
echo "✓ npm 版本: $(npm --version)"

echo ""
echo "================================"
echo "安装后端依赖"
echo "================================"
cd backend
go mod download
go mod tidy
echo "✓ 后端依赖安装完成"

echo ""
echo "================================"
echo "安装前端依赖"
echo "================================"
cd ../frontend
npm install
echo "✓ 前端依赖安装完成"

echo ""
echo "================================"
echo "创建配置文件"
echo "================================"
cd ..

# 复制后端配置示例
if [ ! -f "backend/config/config.yaml" ]; then
    cp backend/config/config.example.yaml backend/config/config.yaml
    echo "✓ 已创建 backend/config/config.yaml"
    echo "  请编辑此文件填入必要的 API 密钥"
else
    echo "⚠ backend/config/config.yaml 已存在，跳过"
fi

# 复制前端环境变量
if [ ! -f "frontend/.env" ]; then
    cp frontend/.env.example frontend/.env
    echo "✓ 已创建 frontend/.env"
else
    echo "⚠ frontend/.env 已存在，跳过"
fi

# 创建存储目录
mkdir -p storage/{uploads,temp,models,images,audio}
echo "✓ 已创建存储目录"

echo ""
echo "================================"
echo "初始化完成！"
echo "================================"
echo ""
echo "下一步:"
echo "1. 编辑 backend/config/config.yaml 填入 API 密钥"
echo "2. 启动后端: cd backend && go run cmd/main.go"
echo "3. 启动前端: cd frontend && npm run dev"
echo ""
echo "或使用 Docker:"
echo "  docker-compose up -d"
echo ""
echo "访问应用: http://localhost:5173"
echo "访问 API: http://localhost:8080"
echo ""
