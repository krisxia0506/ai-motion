# 配置文档

## 配置文件

AI-Motion 使用环境变量进行配置管理。主要配置文件：

- `.env` - 主配置文件 (后端 + 数据库)
- `frontend/.env` - 前端配置文件

---

## 环境变量配置

### 创建配置文件

```bash
# 复制模板文件
cp .env.example .env

# 编辑配置
vim .env
```

---

## 后端配置 (.env)

### 数据库配置

```env
# MySQL 数据库配置
DATABASE_HOST=your-mysql-host        # 数据库主机地址
DATABASE_PORT=3306                   # 数据库端口
DATABASE_USER=ai_motion              # 数据库用户名
DATABASE_PASSWORD=your-password      # 数据库密码
DATABASE_NAME=ai_motion              # 数据库名称
```

#### 配置说明

- `DATABASE_HOST`:
  - 本地开发: `localhost` 或 `127.0.0.1`
  - Docker 环境: `mysql` (服务名)
  - 远程数据库: 实际 IP 或域名

- `DATABASE_PORT`: 默认 `3306`

- `DATABASE_USER`: 建议创建专用用户，不要使用 root

- `DATABASE_PASSWORD`: 使用强密码，至少 16 字符

- `DATABASE_NAME`: 数据库名称，默认 `ai_motion`

### AI 服务配置

```env
# Gemini API 配置 (文生图/图生图)
GEMINI_API_KEY=your-gemini-api-key
GEMINI_API_ENDPOINT=https://generativelanguage.googleapis.com/v1

# Sora2 API 配置 (文生视频/图生视频)
SORA_API_KEY=your-sora-api-key
SORA_API_ENDPOINT=https://api.openai.com/v1/sora

# OpenAI API 配置 (可选，用于文本处理)
OPENAI_API_KEY=your-openai-key
OPENAI_API_ENDPOINT=https://api.openai.com/v1
OPENAI_MODEL=gpt-4
```

#### API Key 获取方式

**Gemini API**
1. 访问 [Google AI Studio](https://makersuite.google.com/app/apikey)
2. 创建新的 API Key
3. 复制 API Key 到配置文件

**Sora2 API**
1. 访问 [OpenAI Platform](https://platform.openai.com)
2. 进入 API Keys 页面
3. 创建新的 API Key
4. 确保账户有 Sora API 访问权限

**OpenAI API**
1. 访问 [OpenAI Platform](https://platform.openai.com)
2. 创建 API Key
3. 根据需求选择模型 (gpt-4, gpt-3.5-turbo 等)

### 应用配置

```env
# 服务器配置
SERVER_PORT=8080                     # 后端服务端口
SERVER_HOST=0.0.0.0                  # 监听地址
ENV=development                      # 环境: development, production

# 日志配置
LOG_LEVEL=info                       # 日志级别: debug, info, warn, error
LOG_FORMAT=json                      # 日志格式: json, text

# 文件存储配置
STORAGE_PATH=./storage               # 文件存储路径
STORAGE_MAX_SIZE=104857600           # 最大文件大小 (100MB)
```

#### 配置说明

- `SERVER_PORT`: 后端服务监听端口，默认 8080
- `ENV`:
  - `development` - 开发环境，启用调试日志
  - `production` - 生产环境，优化性能
- `LOG_LEVEL`:
  - `debug` - 详细调试信息
  - `info` - 一般信息
  - `warn` - 警告信息
  - `error` - 仅错误信息
- `STORAGE_PATH`: 上传文件和生成内容的存储位置

### CORS 配置

```env
# CORS 跨域配置
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://your-domain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

---

## 前端配置 (frontend/.env)

### 创建前端配置

```bash
# 进入前端目录
cd frontend

# 创建配置文件
cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=AI-Motion
VITE_MAX_UPLOAD_SIZE=104857600
EOF
```

### 配置项说明

```env
# API 配置
VITE_API_BASE_URL=http://localhost:8080    # 后端 API 地址

# 应用配置
VITE_APP_TITLE=AI-Motion                   # 应用标题

# 文件上传配置
VITE_MAX_UPLOAD_SIZE=104857600             # 最大上传大小 (100MB)

# 功能开关
VITE_ENABLE_VIDEO_GENERATION=true          # 启用视频生成功能
VITE_ENABLE_VOICE_GENERATION=true          # 启用语音生成功能
```

#### 环境特定配置

Vite 支持环境特定配置文件：

- `.env` - 所有环境
- `.env.local` - 本地环境 (不提交到 git)
- `.env.development` - 开发环境
- `.env.production` - 生产环境

**生产环境配置示例** (`.env.production`):

```env
VITE_API_BASE_URL=https://api.your-domain.com
VITE_APP_TITLE=AI-Motion
VITE_MAX_UPLOAD_SIZE=104857600
```

---

## Docker Compose 配置

### docker-compose.yml 环境变量

Docker Compose 会自动读取项目根目录的 `.env` 文件。

**端口映射配置**:

```yaml
services:
  backend:
    ports:
      - "${BACKEND_PORT:-8080}:8080"

  frontend:
    ports:
      - "${FRONTEND_PORT:-3000}:3000"
```

可以在 `.env` 中自定义端口：

```env
BACKEND_PORT=8080
FRONTEND_PORT=3000
```

---

## 配置最佳实践

### 1. 安全性

- 永远不要提交包含真实凭据的 `.env` 文件到版本控制
- 使用 `.env.example` 作为模板，不包含敏感信息
- 生产环境使用强密码和随机密钥
- 定期轮换 API Keys

### 2. 环境分离

```bash
# 开发环境
.env.development

# 测试环境
.env.test

# 生产环境
.env.production
```

### 3. 配置验证

启动时系统会验证必需的配置项：

必需配置：
- `DATABASE_HOST`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `DATABASE_NAME`
- `GEMINI_API_KEY` 或 `SORA_API_KEY` (至少一个)

### 4. 使用环境变量

在 Docker 部署时，可以直接传递环境变量：

```bash
docker-compose up -d \
  -e DATABASE_PASSWORD=new_password \
  -e GEMINI_API_KEY=your_key
```

---

## 配置示例

### 开发环境完整配置

```env
# 数据库配置
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=ai_motion
DATABASE_PASSWORD=dev_password_123
DATABASE_NAME=ai_motion_dev

# AI 服务配置
GEMINI_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SORA_API_KEY=sk-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# 应用配置
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENV=development
LOG_LEVEL=debug
LOG_FORMAT=text

# 存储配置
STORAGE_PATH=./storage
STORAGE_MAX_SIZE=104857600

# CORS 配置
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 生产环境配置示例

```env
# 数据库配置
DATABASE_HOST=production-mysql.example.com
DATABASE_PORT=3306
DATABASE_USER=ai_motion_prod
DATABASE_PASSWORD=STRONG_RANDOM_PASSWORD_HERE
DATABASE_NAME=ai_motion

# AI 服务配置
GEMINI_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SORA_API_KEY=sk-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# 应用配置
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENV=production
LOG_LEVEL=info
LOG_FORMAT=json

# 存储配置
STORAGE_PATH=/var/lib/ai-motion/storage
STORAGE_MAX_SIZE=104857600

# CORS 配置
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

---

## 故障排查

### 配置未生效

1. 确认 `.env` 文件位置正确
2. 检查环境变量格式 (无空格，无引号)
3. 重启服务使配置生效

```bash
# Docker 环境
docker-compose down
docker-compose up -d

# 本地开发
# 重启后端和前端服务
```

### 数据库连接失败

检查配置项：
- `DATABASE_HOST` 是否正确
- `DATABASE_PORT` 是否开放
- `DATABASE_USER` 和 `DATABASE_PASSWORD` 是否正确
- 数据库是否已创建

### API Key 无效

- 确认 API Key 格式正确
- 检查 API Key 是否已激活
- 验证 API 额度是否充足

---

*配置文档版本: v0.1.0-alpha*
