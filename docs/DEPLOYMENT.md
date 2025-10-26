# AI-Motion 部署指南

本文档提供 AI-Motion 的完整部署方案,涵盖从快速开发环境到生产级云原生架构的多种部署策略。

## 目录

- [部署方案概览](#部署方案概览)
- [方案一: Docker Compose 快速部署](#方案一-docker-compose-快速部署)
- [方案二: 云原生 Kubernetes 部署](#方案二-云原生-kubernetes-部署)
- [方案三: Serverless 部署](#方案三-serverless-部署)
- [环境配置详解](#环境配置详解)
- [数据库部署方案](#数据库部署方案)
- [媒体存储方案](#媒体存储方案)
- [AI 服务集成](#ai-服务集成)
- [性能优化建议](#性能优化建议)
- [监控与日志](#监控与日志)
- [安全加固](#安全加固)
- [CI/CD 流程](#cicd-流程)
- [分阶段部署路线图](#分阶段部署路线图)

---

## 部署方案概览

### 方案对比

| 方案 | 适用场景 | 成本估算 | 优势 | 劣势 |
|------|---------|---------|------|------|
| Docker Compose | 开发环境、MVP 验证、中小规模 | ~¥350/月 | 部署简单、成本低、易调试 | 单点故障、扩展性受限 |
| Kubernetes | 生产环境、高可用需求 | ~¥1000/月起 | 高可用、自动伸缩、零停机更新 | 复杂度高、学习曲线陡 |
| Serverless | 流量波动大、按需付费 | 按使用量 | 自动伸缩、无需运维 | 冷启动延迟、需代码改造 |

---

## 方案一: Docker Compose 快速部署

**推荐用于**: 开发环境、MVP 快速验证、中小规模应用

### 架构图

```
┌──────────────────────────────────────┐
│         单台云服务器 (4核8G)          │
│  ┌────────────────────────────────┐  │
│  │   Nginx (反向代理 + SSL)       │  │
│  └───────┬────────────────────────┘  │
│          │                           │
│  ┌───────▼────────┐  ┌──────────┐   │
│  │  Frontend      │  │ Backend  │   │
│  │  容器 (3000)   │  │容器(8080)│   │
│  └────────────────┘  └─────┬────┘   │
│                            │         │
│                      ┌─────▼────┐   │
│                      │  MySQL   │   │
│                      │  容器     │   │
│                      └──────────┘   │
└──────────────────────────────────────┘
```

### 快速开始

#### 1. 克隆项目并配置

```bash
# 克隆项目
git clone https://github.com/krisxia0506/ai-motion.git
cd ai-motion

# 配置环境变量
cp .env.example .env
vim .env  # 填入真实配置
```

#### 2. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 实时查看日志
docker-compose logs -f
```

#### 3. 验证部署

```bash
# 检查后端健康状态
curl http://localhost:8080/health

# 检查前端服务
curl -I http://localhost:3000
```

访问地址:
- 前端: http://localhost:3000
- 后端 API: http://localhost:8080
- 健康检查: http://localhost:8080/health

### 生产环境配置

#### Nginx 反向代理配置

```nginx
# /etc/nginx/conf.d/ai-motion.conf
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 前端静态资源
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 后端 API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        
        # 大文件上传支持
        client_max_body_size 100M;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }
}
```

### 服务器要求

**最低配置**:
- CPU: 2核
- 内存: 4GB
- 磁盘: 40GB SSD
- 带宽: 5Mbps

**推荐配置**:
- CPU: 4核
- 内存: 8GB
- 磁盘: 100GB SSD
- 带宽: 10Mbps

### 成本估算 (阿里云)

- 4核8G ECS 服务器: ~¥300/月
- 100GB 云盘: ~¥50/月
- **总计**: ~¥350/月

---

## 方案二: 云原生 Kubernetes 部署

**推荐用于**: 生产环境、需要高可用性和弹性伸缩

### 架构图

```
┌─────────────────────────────────────────────────────┐
│              负载均衡器 (ALB/SLB)                     │
│              + SSL/TLS 终止                          │
└────────────┬────────────────────────┬────────────────┘
             │                        │
    ┌────────▼────────┐      ┌───────▼────────┐
    │  前端容器集群    │      │  后端容器集群   │
    │  (React + Nginx)│      │  (Go Backend)  │
    │  副本数: 2-5    │      │  副本数: 3-10   │
    └────────┬────────┘      └───────┬────────┘
             │                        │
             │                ┌───────▼────────┐
             │                │  云数据库 RDS   │
             │                │  (MySQL 8.0)   │
             │                │  主从双机       │
             │                └───────┬────────┘
             │                        │
    ┌────────▼────────────────────────▼────────┐
    │           对象存储 (OSS/S3)              │
    │      (媒体文件、角色参考图)               │
    └────────────────────────────────────────┘
             │
    ┌────────▼────────┐
    │   AI 服务集成    │
    │  Gemini + Sora2 │
    └─────────────────┘
```

### 技术栈

- **容器编排**: Kubernetes (K8s) / 阿里云 ACK / AWS EKS
- **镜像仓库**: 阿里云 ACR / Docker Hub / AWS ECR
- **数据库**: 云托管 MySQL 8.0+ (阿里云 RDS / AWS RDS)
- **对象存储**: 阿里云 OSS / AWS S3 / MinIO
- **负载均衡**: 阿里云 SLB / AWS ALB
- **CDN**: 前端静态资源加速
- **监控**: Prometheus + Grafana

### 部署步骤

#### 1. 构建并推送镜像

```bash
# 登录镜像仓库
docker login your-registry.com

# 构建后端镜像
docker build -f docker/Dockerfile.backend \
  -t your-registry/ai-motion-backend:v0.1.0 .
docker push your-registry/ai-motion-backend:v0.1.0

# 构建前端镜像
docker build -f docker/Dockerfile.frontend \
  -t your-registry/ai-motion-frontend:v0.1.0 .
docker push your-registry/ai-motion-frontend:v0.1.0
```

#### 2. Kubernetes 配置文件

**后端 Deployment** (`k8s/backend-deployment.yaml`):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-motion-backend
  namespace: ai-motion
spec:
  replicas: 3  # 3个副本实现高可用
  selector:
    matchLabels:
      app: ai-motion-backend
  template:
    metadata:
      labels:
        app: ai-motion-backend
    spec:
      containers:
      - name: backend
        image: your-registry/ai-motion-backend:v0.1.0
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: DATABASE_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
        - name: DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        - name: GEMINI_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-api-keys
              key: gemini
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: ai-motion-backend
  namespace: ai-motion
spec:
  selector:
    app: ai-motion-backend
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
  type: ClusterIP
```

**前端 Deployment** (`k8s/frontend-deployment.yaml`):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-motion-frontend
  namespace: ai-motion
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ai-motion-frontend
  template:
    metadata:
      labels:
        app: ai-motion-frontend
    spec:
      containers:
      - name: frontend
        image: your-registry/ai-motion-frontend:v0.1.0
        ports:
        - containerPort: 80
          name: http
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: ai-motion-frontend
  namespace: ai-motion
spec:
  selector:
    app: ai-motion-frontend
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
  type: ClusterIP
```

**Ingress 配置** (`k8s/ingress.yaml`):

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ai-motion-ingress
  namespace: ai-motion
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - your-domain.com
    secretName: ai-motion-tls
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: ai-motion-backend
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: ai-motion-frontend
            port:
              number: 80
```

**Secrets 配置** (`k8s/secrets.yaml`):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
  namespace: ai-motion
type: Opaque
stringData:
  host: "your-rds-endpoint.mysql.rds.aliyuncs.com"
  user: "ai_motion"
  password: "your-strong-password"
  database: "ai_motion"
---
apiVersion: v1
kind: Secret
metadata:
  name: ai-api-keys
  namespace: ai-motion
type: Opaque
stringData:
  gemini: "your-gemini-api-key"
  sora: "your-sora-api-key"
```

#### 3. 部署到 Kubernetes

```bash
# 创建命名空间
kubectl create namespace ai-motion

# 应用 Secrets
kubectl apply -f k8s/secrets.yaml

# 部署后端
kubectl apply -f k8s/backend-deployment.yaml

# 部署前端
kubectl apply -f k8s/frontend-deployment.yaml

# 配置 Ingress
kubectl apply -f k8s/ingress.yaml

# 查看部署状态
kubectl get pods -n ai-motion
kubectl get svc -n ai-motion
kubectl get ingress -n ai-motion
```

#### 4. 配置水平自动伸缩 (HPA)

```yaml
# k8s/backend-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ai-motion-backend-hpa
  namespace: ai-motion
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ai-motion-backend
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

```bash
kubectl apply -f k8s/backend-hpa.yaml
```

### 成本估算 (阿里云 ACK)

- 3 × 2核4G ECS 节点: ~¥600/月
- RDS MySQL 双机高可用: ~¥400/月
- OSS 存储 100GB: ~¥12/月
- SLB 负载均衡: ~¥50/月
- **总计**: ~¥1062/月 起

---

## 方案三: Serverless 部署

**推荐用于**: 流量波动大、按需付费、降低运维成本

### 架构方案

**前端**: 
- Vercel / Netlify / 阿里云 OSS + CDN
- 全球 CDN 加速
- 自动 HTTPS

**后端**: 
- 阿里云函数计算 (FC) / AWS Lambda + API Gateway
- 按请求计费
- 自动伸缩

**数据库**: 
- 阿里云 RDS / AWS RDS for MySQL
- Serverless 版本 (按实际使用计费)

### 前端部署 (Vercel)

```bash
# 安装 Vercel CLI
npm install -g vercel

# 部署前端
cd frontend
vercel deploy --prod
```

### 后端改造建议

1. **将 Gin 路由拆分为独立函数**
2. **使用 API Gateway 统一路由**
3. **优化冷启动时间** (减少依赖、使用预留实例)

### 优势与劣势

**优势**:
- ✅ 自动伸缩,应对流量高峰
- ✅ 按实际使用付费,成本可控
- ✅ 无需运维服务器
- ✅ 前端全球 CDN 加速

**劣势**:
- ⚠️ 冷启动延迟 (首次请求慢)
- ⚠️ 需要改造现有代码
- ⚠️ 调试相对复杂

---

## 环境配置详解

### 后端环境变量

在项目根目录的 `.env` 文件中配置：

```bash
# 数据库配置
DATABASE_HOST=your-mysql-host.com
DATABASE_PORT=3306
DATABASE_USER=ai_motion
DATABASE_PASSWORD=strong-password-here
DATABASE_NAME=ai_motion

# AI 服务 API Keys
GEMINI_API_KEY=your-gemini-api-key
SORA_API_KEY=your-sora-api-key

# 对象存储配置 (如使用 OSS/S3)
STORAGE_TYPE=oss  # 或 s3/local
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
OSS_BUCKET=ai-motion-media
OSS_ACCESS_KEY=your-access-key
OSS_SECRET_KEY=your-secret-key

# 应用配置
APP_ENV=production
APP_PORT=8080
LOG_LEVEL=info

# JWT 密钥 (建议使用强随机字符串)
JWT_SECRET=your-jwt-secret-key-here

# CORS 配置
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

### 前端环境变量

在 `frontend/.env` 文件中配置：

```bash
# API 配置
# 本地开发环境
VITE_API_BASE_URL=http://localhost:8080
# 生产环境示例
# VITE_API_BASE_URL=https://api.your-domain.com

VITE_API_TIMEOUT=30000

# 上传配置
VITE_UPLOAD_MAX_SIZE=52428800  # 50MB

# 功能开关
VITE_ENABLE_PREVIEW=true
VITE_ENABLE_EXPORT=true

# Supabase 配置
VITE_SUPABASE_URL=your_supabase_url
VITE_SUPABASE_ANON_KEY=your_supabase_anon_key
```

**重要说明：**

- **前后端分离部署**: 前端通过 `VITE_API_BASE_URL` 环境变量配置后端 API 地址
- **本地开发**: 使用 `http://localhost:8080`
- **生产环境**: 使用实际的后端 API 域名，例如 `https://api.your-domain.com`
- **Docker 部署**: 如果前后端在同一容器网络，可以使用容器服务名
- **Vercel/Netlify 部署**: 在部署平台的环境变量设置中配置 `VITE_API_BASE_URL`

### 配置文件示例

详细配置说明请参考 [配置文档](./CONFIGURATION.md)。

---

## 数据库部署方案

### 推荐: 云托管数据库

**优势**:
- 自动备份与恢复
- 主从自动切换
- 高可用保障
- 性能监控完善

**推荐服务**:
- 阿里云 RDS MySQL 8.0
- AWS RDS for MySQL
- 腾讯云 TencentDB

### 配置建议

**最小规格**:
- 2核4G (可平滑升级)
- 100GB SSD 存储 (支持自动扩容)
- 主从双机部署
- 每日自动备份,保留 7 天

### 数据库初始化

```bash
# 连接到数据库
mysql -h your-rds-host -u root -p

# 创建数据库
CREATE DATABASE ai_motion CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户并授权
CREATE USER 'ai_motion'@'%' IDENTIFIED BY 'your-strong-password';
GRANT ALL PRIVILEGES ON ai_motion.* TO 'ai_motion'@'%';
FLUSH PRIVILEGES;

# 运行数据库迁移 (如有)
# 当前项目尚未包含迁移脚本,建议使用 golang-migrate 等工具
```

---

## 媒体存储方案

### 对象存储对比

| 方案 | 优势 | 适用场景 | 成本 |
|------|------|---------|------|
| 阿里云 OSS | 国内访问快,CDN 集成好 | 主要用户在中国 | ¥0.12/GB/月 |
| AWS S3 | 全球分布,生态完善 | 国际用户 | $0.023/GB/月 |
| MinIO | 自建,成本可控,私有化 | 私有化部署,大容量 | 仅服务器成本 |

### 推荐存储结构

```
ai-motion-bucket/
├── characters/              # 角色数据
│   └── {character_id}/
│       ├── reference_1.jpg  # 角色参考图 1
│       ├── reference_2.jpg  # 角色参考图 2
│       └── metadata.json    # 角色元数据
├── scenes/                  # 场景数据
│   └── {scene_id}/
│       ├── image.jpg        # 场景图片
│       ├── video.mp4        # 场景视频
│       └── metadata.json    # 场景元数据
└── novels/                  # 小说原文
    └── {novel_id}/
        └── original.txt     # 原始文本
```

### 阿里云 OSS 配置示例

```bash
# 环境变量配置
STORAGE_TYPE=oss
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
OSS_BUCKET=ai-motion-media
OSS_ACCESS_KEY=your-access-key-id
OSS_SECRET_KEY=your-access-key-secret

# 启用 CDN 加速 (可选)
OSS_CDN_DOMAIN=cdn.your-domain.com
```

---

## AI 服务集成

### API 调用优化建议

#### 1. 并发控制

```go
// 限制并发请求数量,避免 API 限流
semaphore := make(chan struct{}, 5) // 最多 5 个并发请求

for _, task := range tasks {
    semaphore <- struct{}{}
    go func(t Task) {
        defer func() { <-semaphore }()
        processTask(t)
    }(task)
}
```

#### 2. 重试机制

```go
maxRetries := 3
backoff := time.Second * 2

for i := 0; i < maxRetries; i++ {
    err := callAIAPI()
    if err == nil {
        break
    }
    time.Sleep(backoff * time.Duration(i+1))
}
```

#### 3. 结果缓存

- 缓存角色参考图生成结果 (避免重复生成)
- 缓存场景描述解析结果
- 使用 Redis 或本地缓存

#### 4. 成本监控

- 记录 API 调用次数和成本
- 设置每日/每月调用上限
- 异常情况告警

---

## 性能优化建议

### 后端优化

- [ ] **添加 Redis 缓存层** (角色参考图、API 响应)
- [ ] **使用 Goroutine Pool** 控制并发
- [ ] **数据库连接池优化** (最大连接数、空闲连接数)
- [ ] **API 请求限流** (防止滥用)
- [ ] **异步任务队列** (长时间 AI 生成任务)

### 前端优化

- [ ] **启用 CDN 加速** 静态资源
- [ ] **图片懒加载和压缩** (WebP 格式)
- [ ] **代码分割** (Code Splitting)
- [ ] **使用 Service Worker** 缓存
- [ ] **React 组件优化** (React.memo, useMemo, useCallback)

### 数据库优化

- [ ] **添加适当索引** (查询频繁的字段)
- [ ] **定期分析查询计划** (EXPLAIN)
- [ ] **读写分离** (主从复制)
- [ ] **分表分库** (数据量大时)

---

## 监控与日志

### 监控方案

**推荐技术栈**:

| 组件 | 用途 | 方案 |
|------|------|------|
| 应用监控 | 性能指标、错误追踪 | Prometheus + Grafana / 云监控 |
| 日志收集 | 集中式日志管理 | ELK Stack / 阿里云 SLS |
| 链路追踪 | 分布式追踪 | Jaeger / SkyWalking |
| 告警 | 异常通知 | AlertManager / 钉钉/企业微信 |

### Prometheus 监控配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'ai-motion-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### 关键监控指标

- **系统资源**: CPU 使用率、内存使用率、磁盘 I/O
- **应用性能**: API 响应时间、QPS、错误率
- **数据库**: 连接数、慢查询、锁等待
- **AI API**: 调用次数、成功率、平均延迟、成本

### 告警规则建议

```yaml
# 告警规则示例
groups:
  - name: ai-motion-alerts
    rules:
      - alert: HighCPUUsage
        expr: cpu_usage > 85
        for: 5m
        annotations:
          summary: "CPU 使用率过高"
      
      - alert: HighMemoryUsage
        expr: memory_usage > 90
        for: 5m
        annotations:
          summary: "内存使用率过高"
      
      - alert: HighAPIErrorRate
        expr: api_error_rate > 5
        for: 2m
        annotations:
          summary: "API 错误率超过 5%"
      
      - alert: SlowAPIResponse
        expr: api_response_time_p95 > 2000
        for: 5m
        annotations:
          summary: "API 响应时间过慢 (P95 > 2s)"
```

### 日志管理

```bash
# 查看后端日志 (Docker Compose)
docker-compose logs -f backend

# 查看特定时间段日志
docker-compose logs --since 2024-01-01T00:00:00 backend

# Kubernetes 日志
kubectl logs -f deployment/ai-motion-backend -n ai-motion

# 查看最近 100 行日志
kubectl logs --tail=100 deployment/ai-motion-backend -n ai-motion
```

---

## 安全加固

### 基础安全措施

- [ ] **HTTPS 强制** (Let's Encrypt 免费证书)
- [ ] **API 限流和防刷** (基于 IP 或用户)
- [ ] **文件上传白名单** (只允许 .txt, .jpg, .png)
- [ ] **SQL 注入防护** (使用参数化查询)
- [ ] **敏感信息加密存储** (密码、API Keys)
- [ ] **定期安全扫描** (依赖漏洞扫描)

### 网络安全

```yaml
# Kubernetes NetworkPolicy 示例
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ai-motion-backend-policy
  namespace: ai-motion
spec:
  podSelector:
    matchLabels:
      app: ai-motion-backend
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: ai-motion-frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: mysql
    ports:
    - protocol: TCP
      port: 3306
```

### 密钥管理

**不要将密钥写入代码或配置文件**,使用以下方式管理:

- **Kubernetes Secrets** (K8s 部署)
- **环境变量注入**
- **云服务密钥管理** (阿里云 KMS / AWS Secrets Manager)
- **HashiCorp Vault** (企业级)

### 定期安全检查

```bash
# Go 依赖漏洞扫描
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# npm 依赖漏洞扫描
npm audit
npm audit fix

# Docker 镜像扫描
docker scan your-image:tag
```

---

## CI/CD 流程

### GitHub Actions 示例

创建 `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  REGISTRY: your-registry.com
  BACKEND_IMAGE: ai-motion-backend
  FRONTEND_IMAGE: ai-motion-frontend

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}
      
      - name: Build and push backend
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./docker/Dockerfile.backend
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.BACKEND_IMAGE }}:${{ github.sha }}
            ${{ env.REGISTRY }}/${{ env.BACKEND_IMAGE }}:latest
      
      - name: Build and push frontend
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./docker/Dockerfile.frontend
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.FRONTEND_IMAGE }}:${{ github.sha }}
            ${{ env.REGISTRY }}/${{ env.FRONTEND_IMAGE }}:latest
  
  deploy-to-k8s:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Set up kubectl
        uses: azure/setup-kubectl@v3
      
      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig.yaml
          export KUBECONFIG=kubeconfig.yaml
      
      - name: Deploy backend
        run: |
          kubectl set image deployment/ai-motion-backend \
            backend=${{ env.REGISTRY }}/${{ env.BACKEND_IMAGE }}:${{ github.sha }} \
            -n ai-motion
          kubectl rollout status deployment/ai-motion-backend -n ai-motion
      
      - name: Deploy frontend
        run: |
          kubectl set image deployment/ai-motion-frontend \
            frontend=${{ env.REGISTRY }}/${{ env.FRONTEND_IMAGE }}:${{ github.sha }} \
            -n ai-motion
          kubectl rollout status deployment/ai-motion-frontend -n ai-motion
```

### 部署流程

1. **代码推送到 main 分支**
2. **自动触发 CI/CD**
3. **运行测试** (单元测试、集成测试)
4. **构建 Docker 镜像**
5. **推送镜像到仓库**
6. **部署到 Kubernetes**
7. **健康检查验证**
8. **发送部署通知** (钉钉/企业微信)

---

## 分阶段部署路线图

### Phase 1: MVP 快速验证 (当前阶段)

**目标**: 快速上线,验证核心功能

**方案**: Docker Compose 单服务器部署

**配置**:
- 1 台 4核8G 云服务器
- Docker Compose 编排
- 本地 MySQL 容器

**成本**: ~¥350/月

**时间**: 1-2 周

---

### Phase 2: 正式上线 (1-3个月)

**目标**: 稳定服务,支持中等流量

**方案**: 云原生 Kubernetes 部署

**配置**:
- 3 × 2核4G 容器节点
- 云数据库 RDS (主从)
- 对象存储 OSS
- 负载均衡 SLB

**成本**: ~¥1000/月

**新增功能**:
- 自动伸缩
- 滚动更新
- 健康检查
- 监控告警

**时间**: 2-4 周

---

### Phase 3: 规模化运营 (3个月后)

**目标**: 高可用,应对流量高峰

**方案**: 混合架构 (K8s + Serverless + CDN)

**配置**:
- 弹性伸缩 (5-20 副本)
- 全球 CDN 加速
- Redis 集群缓存
- 消息队列 (异步任务)

**成本**: 按实际使用 ~¥2000-5000/月

**新增功能**:
- 多地域部署
- 读写分离
- 分布式追踪
- AI 成本优化

---

## 备份与恢复

### 数据库备份

```bash
# 自动备份脚本 (每日凌晨 2 点)
0 2 * * * docker exec ai-motion-mysql \
  mysqldump -u ai_motion -p${DB_PASSWORD} ai_motion \
  | gzip > /backup/ai_motion_$(date +\%Y\%m\%d).sql.gz

# 保留最近 7 天备份
find /backup -name "ai_motion_*.sql.gz" -mtime +7 -delete
```

### 文件存储备份

```bash
# OSS 跨区域复制 (阿里云)
# 在 OSS 控制台配置跨区域复制规则

# 或使用 ossutil 同步
ossutil sync oss://ai-motion-media/ oss://ai-motion-backup/
```

### 恢复演练

**建议每季度进行一次恢复演练**,验证备份有效性。

---

## 更新与升级

### 滚动更新 (Kubernetes)

```bash
# 更新后端镜像
kubectl set image deployment/ai-motion-backend \
  backend=your-registry/ai-motion-backend:v0.2.0 \
  -n ai-motion

# 查看滚动更新状态
kubectl rollout status deployment/ai-motion-backend -n ai-motion

# 回滚到上一版本 (如有问题)
kubectl rollout undo deployment/ai-motion-backend -n ai-motion
```

### Docker Compose 更新

```bash
# 拉取最新代码
git pull origin main

# 重建并重启容器 (会有短暂停机)
docker-compose down
docker-compose up -d --build

# 查看更新后的状态
docker-compose ps
```

---

## 故障排查

### 常见问题

详细故障排查指南请参考 [故障排查文档](./TROUBLESHOOTING.md)。

**快速诊断**:

```bash
# 检查容器状态
docker-compose ps
kubectl get pods -n ai-motion

# 查看日志
docker-compose logs --tail=100 backend
kubectl logs -f deployment/ai-motion-backend -n ai-motion

# 检查资源使用
docker stats
kubectl top pods -n ai-motion

# 检查网络连通性
curl http://localhost:8080/health
kubectl exec -it <pod-name> -n ai-motion -- curl http://backend:8080/health
```

---

## 总结

### 推荐部署路径

1. **短期 (1-3个月)**: 使用 **Docker Compose 方案** 快速部署到单台云服务器,验证产品可行性
2. **中期 (3-6个月)**: 迁移到 **Kubernetes 云原生架构**,实现高可用和弹性伸缩
3. **长期 (6个月+)**: 引入 **Serverless** 和 **CDN**,优化成本和性能

### 立即可以开始的工作

- [x] 完善 `.env.example`,补充 AI API Keys 等配置项
- [ ] 添加数据库迁移脚本 (推荐使用 golang-migrate)
- [ ] 创建 K8s 部署配置文件 (参考本文档)
- [ ] 配置 GitHub Actions CI/CD
- [ ] 添加 Prometheus 监控端点
- [ ] 编写备份恢复脚本

---

**文档版本**: v0.2.0  
**最后更新**: 2024-10-24  
**维护者**: AI-Motion Team

如有部署问题,请参考 [故障排查文档](./TROUBLESHOOTING.md) 或提交 Issue。
