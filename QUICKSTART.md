# 快速启动指南

本指南将帮助你快速启动 AI-Motion 项目。根据你的需求选择合适的启动方式。

## 前置要求

根据不同的启动方式,你需要准备:

### 方式一: Docker 部署 (推荐)
- Docker 20.0+
- Docker Compose 2.0+

### 方式二: 本地开发
- Go 1.24+
- Node.js 20+
- Supabase 账号 (用于 PostgreSQL 数据库)

## 方式一: Docker 部署 (推荐)

这是最简单快速的启动方式,适合快速体验和生产环境部署。

### 步骤 1: 克隆项目

```bash
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion
```

### 步骤 2: 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件,配置数据库连接信息
# 如果暂时不需要数据库,可以使用默认配置
nano .env  # 或使用其他编辑器
```

.env 配置示例:
```env
DATABASE_HOST=your-mysql-host.com
DATABASE_PORT=3306
DATABASE_USER=ai_motion
DATABASE_PASSWORD=your-secure-password
DATABASE_NAME=ai_motion
```

### 步骤 3: 启动服务

```bash
# 启动所有服务(后端 + 前端)
docker-compose up -d

# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 步骤 4: 验证服务

```bash
# 测试后端健康检查
curl http://localhost:8080/health

# 预期返回:
# {"service":"ai-motion","status":"ok"}
```

### 步骤 5: 访问应用

- **前端界面**: http://localhost:3000
- **后端 API**: http://localhost:8080
- **健康检查**: http://localhost:8080/health

### 常用 Docker 命令

```bash
# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 重新构建并启动
docker-compose up -d --build

# 查看容器日志
docker-compose logs backend   # 后端日志
docker-compose logs frontend  # 前端日志
docker-compose logs -f        # 实时查看所有日志

# 进入容器
docker exec -it ai-motion-backend sh
docker exec -it ai-motion-frontend sh
```

## 方式二: 本地开发

适合需要修改代码、调试或开发新功能的场景。

### 步骤 1: 克隆项目

```bash
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion
```

### 步骤 2: 使用 Make 快速启动

```bash
# 查看所有可用命令
make help

# 安装所有依赖(后端 + 前端)
make install

# 启动开发环境(会同时启动后端和前端)
make dev
```

### 步骤 3: 手动启动 (可选)

如果你想分别启动后端和前端:

#### 启动后端

```bash
# 进入后端目录
cd backend

# 安装 Go 依赖
go mod download

# 运行后端服务
go run cmd/main.go

# 后端将在 http://localhost:8080 启动
```

#### 启动前端

在新的终端窗口中:

```bash
# 进入前端目录
cd frontend

# 安装 npm 依赖
npm install

# 启动前端开发服务器
npm run dev

# 前端将在 http://localhost:5173 启动
```

### 步骤 4: 验证服务

```bash
# 测试后端 API
curl http://localhost:8080/health

# 访问前端(开发模式)
open http://localhost:5173
```

## 方式三: 使用初始化脚本

项目提供了自动化初始化脚本。

```bash
# 添加执行权限
chmod +x scripts/setup.sh

# 运行初始化脚本
./scripts/setup.sh

# 脚本会自动:
# 1. 检查环境依赖
# 2. 安装后端依赖
# 3. 安装前端依赖
# 4. 创建必要的目录
```

## 端口说明

默认端口配置:

| 服务 | 开发模式 | Docker 模式 |
|------|---------|------------|
| 后端 API | 8080 | 8080 |
| 前端 | 5173 | 3000 |

如需修改端口:
- **Docker 模式**: 编辑 [docker-compose.yml](docker-compose.yml) 中的 `ports` 配置
- **开发模式**:
  - 后端: 修改 `backend/cmd/main.go` 中的端口
  - 前端: 修改 `frontend/vite.config.ts` 中的 server 配置

## 验证安装

### 1. 检查后端服务

```bash
# 健康检查
curl http://localhost:8080/health

# 测试 API 路由
curl http://localhost:8080/api/v1/novel/upload
```

### 2. 检查前端服务

访问浏览器:
- 开发模式: http://localhost:5173
- Docker 模式: http://localhost:3000

## 配置说明

### 后端配置

后端目前通过环境变量配置,主要配置项在 `.env` 文件中。

**当前可用配置:**
- `DATABASE_HOST`: MySQL 主机地址
- `DATABASE_PORT`: MySQL 端口 (默认 3306)
- `DATABASE_USER`: 数据库用户名
- `DATABASE_PASSWORD`: 数据库密码
- `DATABASE_NAME`: 数据库名称

**注意**: 目前数据库配置是为将来功能预留,后端服务可以在没有数据库的情况下启动。

### 前端配置

前端配置在 `frontend/.env` 中:

```env
# API 基础地址
VITE_API_BASE_URL=http://localhost:8080
```

在 Docker 模式下,前端会自动使用正确的 API 地址。

## 故障排查

### Docker 相关问题

#### 问题: 容器无法启动

```bash
# 查看详细错误日志
docker-compose logs backend
docker-compose logs frontend

# 清理并重建
docker-compose down
docker-compose up -d --build
```

#### 问题: 端口被占用

```bash
# macOS/Linux: 检查端口占用
lsof -i :8080  # 检查后端端口
lsof -i :3000  # 检查前端端口

# Windows: 检查端口占用
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# 解决方案:
# 1. 停止占用端口的进程
# 2. 或修改 docker-compose.yml 中的端口映射
```

#### 问题: 镜像构建失败

```bash
# 清理 Docker 缓存
docker system prune -a

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

### 本地开发问题

#### 问题: Go 依赖下载失败

```bash
# 设置 Go 代理(中国大陆)
export GOPROXY=https://goproxy.cn,direct

# 或使用官方代理
export GOPROXY=https://proxy.golang.org,direct

# 重新下载依赖
cd backend
go mod download
```

#### 问题: Node.js 依赖安装失败

```bash
# 清理缓存
cd frontend
rm -rf node_modules package-lock.json

# 使用淘宝镜像
npm install --registry=https://registry.npmmirror.com

# 或使用 cnpm
npm install -g cnpm --registry=https://registry.npmmirror.com
cnpm install
```

#### 问题: 前端连接不到后端

1. 确认后端服务正在运行:
   ```bash
   curl http://localhost:8080/health
   ```

2. 检查前端 `.env` 配置:
   ```bash
   cat frontend/.env
   # 确保 VITE_API_BASE_URL=http://localhost:8080
   ```

3. 检查浏览器控制台是否有 CORS 错误

### 权限问题

#### macOS/Linux: scripts/setup.sh 无法执行

```bash
# 添加执行权限
chmod +x scripts/setup.sh
```

#### Docker: 权限拒绝

```bash
# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或刷新用户组
newgrp docker
```

## 性能优化

### Docker 构建加速

1. 使用 BuildKit:
   ```bash
   export DOCKER_BUILDKIT=1
   docker-compose build
   ```

2. 使用国内镜像源(中国大陆用户):
   编辑 `/etc/docker/daemon.json`:
   ```json
   {
     "registry-mirrors": [
       "https://docker.mirrors.ustc.edu.cn"
     ]
   }
   ```

### 开发模式优化

1. 启用 Go 模块缓存:
   ```bash
   export GOMODCACHE=$HOME/go/pkg/mod
   ```

2. 使用 npm 缓存:
   ```bash
   npm config set cache ~/.npm-cache
   ```

## 下一步

成功启动后,你可以:

1. **查看 API 文档**:
   - 查看 [README.md](README.md) 中的 API 接口说明

2. **开始开发**:
   - 后端代码在 `backend/` 目录
   - 前端代码在 `frontend/` 目录

3. **查看架构文档**:
   - [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

4. **准备测试数据**:
   - 准备一个 TXT 格式的小说文件
   - 通过 API 上传测试

## 开发工作流

### 推荐的开发流程

1. **代码修改**
   - 修改后端代码后,后端会自动重启(使用 `air` 或手动重启)
   - 修改前端代码后,Vite 会自动热更新

2. **测试更改**
   ```bash
   # 后端测试
   cd backend && go test ./...

   # 前端测试
   cd frontend && npm test
   ```

3. **构建生产版本**
   ```bash
   make build

   # 或分别构建
   cd backend && go build -o bin/ai-motion cmd/main.go
   cd frontend && npm run build
   ```

4. **Docker 部署测试**
   ```bash
   docker-compose up -d --build
   ```

## 获取帮助

如果遇到问题:

1. 查看 [README.md](README.md) 中的详细文档
2. 查看项目 [Issues](https://github.com/xiajiayi/ai-motion/issues)
3. 提交新的 Issue 描述你的问题
4. 查看 [docs/](docs/) 目录中的其他文档

## 有用的资源

- [Go 官方文档](https://go.dev/doc/)
- [React 官方文档](https://react.dev/)
- [Vite 文档](https://vitejs.dev/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [Docker 文档](https://docs.docker.com/)

---

**提示**: 首次启动推荐使用 Docker 方式,可以避免环境配置问题,快速体验项目功能。
