# 故障排查指南

本文档提供常见问题的解决方案和故障排查步骤。

---

## Docker 相关问题

### 容器无法启动

**问题描述**: 执行 `docker-compose up -d` 后容器无法正常启动

**排查步骤**:

1. 查看容器状态
```bash
docker-compose ps
```

2. 查看详细日志
```bash
# 查看所有容器日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mysql

# 实时跟踪日志
docker-compose logs -f --tail=100
```

3. 常见原因及解决方案

**原因 1: 端口被占用**
```bash
# 检查端口占用
lsof -i :8080  # 后端端口
lsof -i :3000  # 前端端口
lsof -i :3306  # MySQL 端口

# 解决方案 1: 停止占用端口的进程
kill -9 <PID>

# 解决方案 2: 修改端口映射
# 编辑 docker-compose.yml
services:
  backend:
    ports:
      - "8081:8080"  # 改用 8081 端口
```

**原因 2: 配置文件缺失或错误**
```bash
# 检查 .env 文件是否存在
ls -la .env

# 如果不存在，从模板创建
cp .env.example .env

# 验证配置格式
cat .env
```

**原因 3: 镜像构建失败**
```bash
# 重新构建镜像
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### 容器频繁重启

**问题描述**: 容器启动后立即退出并不断重启

**排查步骤**:

1. 查看容器退出原因
```bash
docker-compose ps
docker-compose logs backend --tail=50
```

2. 常见原因

**数据库连接失败**
```bash
# 检查数据库容器是否正常运行
docker-compose ps mysql

# 检查数据库日志
docker-compose logs mysql

# 测试数据库连接
docker-compose exec mysql mysql -u ai_motion -p
```

**环境变量配置错误**
```bash
# 检查容器内的环境变量
docker-compose exec backend env | grep DATABASE

# 确认 .env 文件配置正确
cat .env | grep DATABASE
```

### 数据持久化问题

**问题描述**: 容器重启后数据丢失

**解决方案**:

确保 docker-compose.yml 中配置了数据卷：

```yaml
services:
  mysql:
    volumes:
      - mysql-data:/var/lib/mysql

volumes:
  mysql-data:
```

**备份现有数据**:
```bash
# 导出数据库
docker-compose exec mysql mysqldump -u ai_motion -p ai_motion > backup.sql

# 备份存储文件
tar -czf storage-backup.tar.gz storage/
```

---

## 数据库问题

### 无法连接数据库

**问题描述**: 应用提示数据库连接失败

**排查步骤**:

1. 检查数据库服务状态
```bash
# Docker 环境
docker-compose ps mysql

# 本地 MySQL
systemctl status mysql  # Linux
brew services list | grep mysql  # macOS
```

2. 验证连接参数
```bash
# 测试数据库连接
mysql -h DATABASE_HOST -P 3306 -u ai_motion -p

# Docker 环境测试
docker-compose exec mysql mysql -u ai_motion -p
```

3. 常见错误

**错误: Access denied for user**
```bash
# 原因: 用户名或密码错误
# 解决: 重置数据库密码

# 进入 MySQL 容器
docker-compose exec mysql mysql -u root -p

# 重置用户密码
ALTER USER 'ai_motion'@'%' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;

# 更新 .env 文件中的密码
```

**错误: Unknown database**
```bash
# 原因: 数据库不存在
# 解决: 创建数据库

docker-compose exec mysql mysql -u root -p
CREATE DATABASE ai_motion CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**错误: Can't connect to MySQL server**
```bash
# 原因 1: MySQL 服务未启动
docker-compose up -d mysql

# 原因 2: 网络问题 (Docker 环境)
# 检查 Docker 网络
docker network ls
docker network inspect ai-motion_default

# 原因 3: 主机地址错误
# Docker 环境应使用服务名: DATABASE_HOST=mysql
# 本地开发应使用: DATABASE_HOST=localhost
```

### 数据库性能问题

**问题描述**: 查询缓慢或超时

**解决方案**:

1. 检查慢查询日志
```bash
docker-compose exec mysql mysql -u root -p
SHOW VARIABLES LIKE 'slow_query%';
SET GLOBAL slow_query_log = 'ON';
```

2. 优化查询
```bash
# 查看表索引
SHOW INDEX FROM table_name;

# 分析查询计划
EXPLAIN SELECT * FROM table_name WHERE ...;
```

---

## 后端服务问题

### API 请求失败

**问题描述**: 前端无法访问后端 API

**排查步骤**:

1. 检查后端服务状态
```bash
# 健康检查
curl http://localhost:8080/health

# 检查端口监听
lsof -i :8080
netstat -an | grep 8080
```

2. 查看后端日志
```bash
# Docker 环境
docker-compose logs backend --tail=100

# 本地开发
# 查看应用控制台输出
```

3. 常见错误

**错误: Connection refused**
```bash
# 原因: 后端服务未启动或监听地址错误
# 检查 SERVER_HOST 配置
# Docker 环境应使用: SERVER_HOST=0.0.0.0
# 本地开发可使用: SERVER_HOST=localhost
```

**错误: CORS 错误**
```bash
# 原因: 跨域配置不正确
# 解决: 更新 .env 配置

CORS_ALLOWED_ORIGINS=http://localhost:3000,https://your-domain.com
```

### AI 服务调用失败

**问题描述**: 图片生成或视频生成失败

**排查步骤**:

1. 验证 API Key
```bash
# 检查环境变量
docker-compose exec backend env | grep API_KEY

# 测试 Gemini API
curl -H "x-goog-api-key: YOUR_API_KEY" \
  https://generativelanguage.googleapis.com/v1/models
```

2. 常见错误

**错误: API key not valid**
```bash
# 原因: API Key 无效或过期
# 解决: 重新生成 API Key 并更新配置
```

**错误: Quota exceeded**
```bash
# 原因: API 配额用尽
# 解决:
# 1. 检查账户配额
# 2. 升级账户套餐
# 3. 等待配额重置
```

**错误: Request timeout**
```bash
# 原因: 网络问题或 API 响应慢
# 解决:
# 1. 检查网络连接
# 2. 增加超时时间配置
# 3. 使用代理 (如果在国内)
```

---

## 前端问题

### 前端无法访问

**问题描述**: 无法打开 http://localhost:3000

**排查步骤**:

1. 检查前端服务
```bash
# Docker 环境
docker-compose ps frontend
docker-compose logs frontend

# 本地开发
# 检查 npm run dev 是否正常运行
```

2. 检查端口
```bash
# 检查 3000 端口是否被占用
lsof -i :3000

# 如果被占用，修改 Vite 配置
# 编辑 frontend/vite.config.ts
server: {
  port: 3001  // 使用其他端口
}
```

### API 请求失败

**问题描述**: 前端调用后端 API 失败

**排查步骤**:

1. 检查 API 配置
```bash
# 检查前端环境变量
cat frontend/.env

# 确认 VITE_API_BASE_URL 正确
VITE_API_BASE_URL=http://localhost:8080
```

2. 浏览器控制台检查
- 打开浏览器开发者工具 (F12)
- 查看 Network 标签
- 检查请求 URL 和响应状态

3. 常见错误

**错误: net::ERR_CONNECTION_REFUSED**
```bash
# 原因: 后端服务未启动
# 解决: 启动后端服务
docker-compose up -d backend
```

**错误: CORS policy blocked**
```bash
# 原因: CORS 配置问题
# 解决: 更新后端 CORS 配置
# 在 .env 中添加前端地址
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

---

## 开发环境问题

### Go 依赖下载失败

**问题描述**: `go mod download` 失败

**解决方案**:

```bash
# 设置 Go 代理 (国内用户)
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 如果仍然失败，尝试 tidy
go mod tidy
```

### Node.js 依赖安装失败

**问题描述**: `npm install` 失败

**解决方案**:

```bash
# 清理缓存
npm cache clean --force

# 删除 node_modules
rm -rf frontend/node_modules
rm -rf frontend/package-lock.json

# 使用淘宝镜像 (国内用户)
npm config set registry https://registry.npmmirror.com

# 重新安装
cd frontend
npm install

# 或使用 pnpm/yarn
pnpm install
# 或
yarn install
```

---

## 性能问题

### 响应速度慢

**排查步骤**:

1. 检查资源使用
```bash
# 查看容器资源占用
docker stats

# 查看系统资源
top  # Linux/macOS
```

2. 优化建议

**数据库优化**:
- 添加适当索引
- 优化慢查询
- 调整连接池大小

**后端优化**:
- 启用响应缓存
- 优化日志级别 (生产环境使用 info 或 warn)
- 调整并发处理数

**前端优化**:
- 启用生产模式构建
- 使用 CDN 加载静态资源
- 实施代码分割

---

## 日志分析

### 查看日志

```bash
# Docker 环境 - 查看所有日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mysql

# 实时跟踪日志
docker-compose logs -f

# 查看最近 N 行日志
docker-compose logs --tail=50 backend

# 查看特定时间范围的日志
docker-compose logs --since 2024-01-01T00:00:00 backend
```

### 调整日志级别

```env
# 开发环境 - 详细日志
LOG_LEVEL=debug
LOG_FORMAT=text

# 生产环境 - 简洁日志
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 获取帮助

如果上述方法无法解决问题:

1. 查看项目 [Issues](https://github.com/xiajiayi/ai-motion/issues)
2. 提交新的 Issue，包含:
   - 问题描述
   - 重现步骤
   - 错误日志
   - 环境信息 (OS, Docker 版本等)
3. 加入社区讨论

---

## 诊断信息收集

提交 Issue 时，请提供以下信息:

```bash
# 系统信息
uname -a

# Docker 版本
docker --version
docker-compose --version

# 容器状态
docker-compose ps

# 近期日志
docker-compose logs --tail=100 > logs.txt

# 配置文件 (移除敏感信息)
cat .env.example
```

---

*故障排查文档版本: v0.1.0-alpha*
