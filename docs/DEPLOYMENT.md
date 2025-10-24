# 部署指南

## 部署方式

AI-Motion 支持多种部署方式，推荐使用 Docker Compose 进行部署。

## 环境要求

### 开发环境
- Go 1.24+
- Node.js 20+
- MySQL 8.0+
- Docker 20+ (可选)
- Docker Compose 2.0+ (可选)

### 生产环境
- Docker 20+
- Docker Compose 2.0+
- 至少 2GB 内存
- 至少 10GB 磁盘空间

---

## Docker 部署 (推荐)

### 1. 克隆项目

```bash
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入必要的配置
vim .env
```

详细配置说明请参考 [配置文档](./CONFIGURATION.md)。

### 3. 启动服务

```bash
# 启动所有服务 (后台运行)
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 验证部署

```bash
# 检查后端服务
curl http://localhost:8080/health

# 检查前端服务
curl -I http://localhost:3000
```

访问地址：
- 前端: http://localhost:3000
- 后端: http://localhost:8080
- 健康检查: http://localhost:8080/health

### 5. 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷 (谨慎使用)
docker-compose down -v
```

---

## 本地开发部署

### 1. 准备数据库

确保 MySQL 8.0+ 已安装并运行：

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE ai_motion CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'ai_motion'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON ai_motion.* TO 'ai_motion'@'localhost';
FLUSH PRIVILEGES;
```

### 2. 后端服务

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 运行后端服务
go run cmd/main.go
```

后端服务将运行在 http://localhost:8080

### 3. 前端服务

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 运行前端服务
npm run dev
```

前端服务将运行在 http://localhost:3000

---

## 生产环境部署

### 使用 Docker Compose (推荐)

1. **配置生产环境变量**

```bash
cp .env.example .env.production
vim .env.production
```

确保配置：
- 生产数据库凭据
- AI 服务 API Keys
- 安全密钥

2. **使用生产配置启动**

```bash
docker-compose -f docker-compose.yml --env-file .env.production up -d
```

3. **配置 Nginx (可选)**

如果需要使用自定义域名或 HTTPS，可以在前面添加 Nginx 反向代理：

```nginx
# /etc/nginx/sites-available/ai-motion
server {
    listen 80;
    server_name your-domain.com;

    # 前端
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 后端 API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 使用 Make 命令

项目提供了便捷的 Make 命令：

```bash
# 查看所有可用命令
make help

# 构建项目
make build

# 启动 Docker 服务
make docker-up

# 停止 Docker 服务
make docker-down

# 查看日志
make logs
```

---

## 服务监控

### 健康检查

```bash
# 检查后端健康状态
curl http://localhost:8080/health

# 检查所有容器状态
docker-compose ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend

# 实时跟踪日志
docker-compose logs -f --tail=100
```

### 资源监控

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
docker system df
```

---

## 备份与恢复

### 数据库备份

```bash
# 备份数据库
docker exec ai-motion-mysql mysqldump -u ai_motion -p ai_motion > backup.sql

# 恢复数据库
docker exec -i ai-motion-mysql mysql -u ai_motion -p ai_motion < backup.sql
```

### 文件存储备份

```bash
# 备份存储目录
tar -czf storage-backup.tar.gz storage/

# 恢复存储目录
tar -xzf storage-backup.tar.gz
```

---

## 更新与升级

### Docker 部署更新

```bash
# 拉取最新代码
git pull origin main

# 重建并重启容器
docker-compose down
docker-compose up -d --build

# 查看更新后的状态
docker-compose ps
```

### 本地开发更新

```bash
# 更新后端
cd backend
go mod download
go build -o bin/server cmd/main.go

# 更新前端
cd frontend
npm install
npm run build
```

---

## 安全建议

1. **修改默认密码**: 确保修改 `.env` 中的所有默认密码
2. **使用 HTTPS**: 生产环境建议配置 SSL 证书
3. **限制端口访问**: 仅开放必要的端口（80, 443）
4. **定期备份**: 建立自动备份机制
5. **日志管理**: 配置日志轮转，避免磁盘空间耗尽
6. **API Key 安全**: 妥善保管 AI 服务的 API Keys

---

## 扩展部署

### 水平扩展

未来版本将支持：
- 多实例负载均衡
- Redis 缓存集群
- 对象存储 (S3/OSS)
- 消息队列 (RabbitMQ/Kafka)

### 云平台部署

项目可部署到以下云平台：
- AWS (ECS/EKS)
- 阿里云 (ACK)
- 腾讯云 (TKE)
- 自建 Kubernetes 集群

---

*部署文档版本: v0.1.0-alpha*
