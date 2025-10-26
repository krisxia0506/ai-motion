# AI-Motion — 智能动漫生成系统

文本小说一键生成动漫，角色一致性保障 → 智能场景分割 → 图文配音输出，异步任务流程，生成记录、质量评估、角色复用与优化。

## 技术栈

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.10-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-7-646CFF?logo=vite&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-4-38B2AC?logo=tailwind-css&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-4169E1?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)

**关键代码参考：**
- 工作流服务：[backend/internal/application/service/manga_workflow_service.go](backend/internal/application/service/manga_workflow_service.go)
- 前端工作流：[frontend/src/pages/Dashboard.tsx](frontend/src/pages/Dashboard.tsx)
- 角色一致性设计：[docs/CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)
- AI 服务集成：`backend/internal/infrastructure/ai/`
- 数据模型：[backend/internal/domain/entity/](backend/internal/domain/entity/)
- API 路由：[backend/internal/interfaces/http/handler/](backend/internal/interfaces/http/handler/)

---

## 技术架构

### 架构概览

AI-Motion 是一个基于 **DDD（领域驱动设计）** 的全栈应用，采用 Go + React 技术栈，实现从文本小说到动漫内容的 AI 自动化生成服务。

```
┌─────────────────────────────────────────────────┐
│           前端层 (React 19 + TypeScript)         │
│    Vite 7 + TailwindCSS + Supabase Auth         │
│         工作流管理 + 实时任务状态追踪            │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│         接口层 (Gin HTTP Handlers)              │
│    RESTful API + JWT 中间件 + CORS              │
└─────┬───────────────────────────┬───────────────┘
      │                           │
┌─────▼──────────┐      ┌─────────▼────────────┐
│   应用层        │      │   领域层              │
│   工作流编排    │◄────►│   核心业务逻辑        │
│   Use Cases    │      │   Entity/Value Object│
└─────┬──────────┘      └──────────────────────┘
      │
┌─────▼──────────────────────────────────────────┐
│              基础设施层                         │
│  Repository (Supabase) + AI Services           │
│  Gemini 2.5 Flash (图像) + Sora2 (视频)        │
│  Supabase Storage (图片/视频/音频)              │
└─────────────────────────────────────────────────┘
```

### DDD 分层架构

#### **领域层 (Domain Layer)** - 核心业务逻辑
- **Novel Entity** - 小说实体：元数据、内容解析、状态管理
- **Character Entity** - 角色实体：信息提取、一致性追踪、参考图管理
- **Scene Entity** - 场景实体：分割算法、对话提取、提示词生成
- **Task Entity** - 任务实体：异步任务状态机、错误处理
- **无外部依赖** - 纯业务逻辑，不依赖框架和基础设施

#### **应用层 (Application Layer)** - 用例编排
- **MangaWorkflowService** - 动漫生成工作流
  - 小说解析 → 角色提取 → 场景分割 → 图像生成 → 配音合成
  - 异步任务管理、状态追踪、错误恢复
- **DTO (Data Transfer Objects)** - 数据传输对象
- **Use Cases** - 业务用例封装

#### **基础设施层 (Infrastructure Layer)** - 外部服务
- **Repository 实现** - Supabase PostgreSQL + PostgREST
  - NovelRepository, CharacterRepository, SceneRepository, TaskRepository
- **AI Services**
  - **Gemini 2.5 Flash Image** - 文生图/图生图，角色参考图生成
  - **Sora2** - 文生视频/图生视频，场景动画生成
- **Storage Service** - Supabase Storage，媒体文件管理

#### **接口层 (Interface Layer)** - HTTP API
- **Gin HTTP Handlers** - RESTful API 路由
- **Middleware** - JWT 认证、CORS、日志、错误处理
- **Request/Response** - 请求验证、响应序列化

### 项目结构

```
ai-motion/
├── backend/                      # Go 后端 (DDD 架构)
│   ├── cmd/main.go              # 应用入口
│   ├── internal/
│   │   ├── domain/              # 领域层 - 核心业务逻辑
│   │   │   ├── entity/         # 实体：Novel, Character, Scene, Task
│   │   │   ├── repository/     # 仓储接口定义
│   │   │   └── service/        # 领域服务
│   │   ├── application/         # 应用层 - 用例编排
│   │   │   ├── dto/            # 数据传输对象
│   │   │   └── service/        # 应用服务 (MangaWorkflowService)
│   │   ├── infrastructure/      # 基础设施层
│   │   │   ├── ai/             # AI 服务 (Gemini, Sora)
│   │   │   ├── database/       # 数据库连接 (Supabase)
│   │   │   ├── repository/     # 仓储实现
│   │   │   └── middleware/     # HTTP 中间件
│   │   └── interfaces/          # 接口层
│   │       └── http/           # HTTP API (Gin Handlers)
│   └── pkg/                     # 公共工具包
├── frontend/                     # React 前端
│   └── src/
│       ├── components/          # UI 组件
│       ├── pages/               # 页面 (Dashboard, WorkflowDetail)
│       ├── services/            # API 服务层
│       └── lib/                 # 工具函数 (Supabase Client)
├── docs/                        # 项目文档
│   ├── ARCHITECTURE.md          # 完整架构设计
│   ├── CHARACTER_CONSISTENCY.md # 角色一致性方案
│   ├── API.md                   # API 接口文档
│   └── DEVELOPMENT.md           # 开发指南
├── docker/                      # Docker 配置
│   ├── backend.Dockerfile
│   └── frontend.Dockerfile
└── scripts/                     # 工具脚本
    └── setup.sh
```

详细架构说明请查看 [ARCHITECTURE.md](docs/ARCHITECTURE.md)。

---

## 1) 面向用户与用户故事

### 目标用户
- **内容创作者/网文作者** - 希望将文字作品可视化，提升作品表现力和传播力
- **独立动漫工作室** - 需要快速产出分镜和预览内容，降低制作成本
- **教育/娱乐从业者** - 需要批量生成教学/娱乐动漫素材，无重型制作工具和团队

### 主要痛点
- **传统动漫制作成本高、周期长** - 需要专业团队和工具链，个人创作者门槛高
- **AI 生成角色不一致** - 同一角色在不同场景中外观差异大，影响观看体验
- **手动场景分割耗时** - 需要人工阅读小说、提取对话、设计分镜，效率低下
- **缺乏端到端解决方案** - 现有工具碎片化，需要在多个平台间切换

### 用户故事
- **作为网文作者**，我上传一部小说，系统自动提取角色、分割场景、生成图片和配音，30 分钟内得到完整动漫预览
- **作为独立制作人**，我先查看系统生成的角色参考图，确认视觉风格后再批量生成场景，确保角色一致性
- **作为教育工作者**，我给生成结果打分和标注问题，系统根据反馈优化生成参数，减少失败率

---

## 2) 功能与优先级

### 已实现 (P0) ✅
- [x] **异步工作流系统** - 小说上传 → 任务创建 → 状态追踪 → 结果查看
- [x] **Supabase 集成** - 认证、数据库、存储完整集成
- [x] **前端工作流界面** - 任务列表、详情页、状态实时更新
- [x] **DDD 架构搭建** - 领域层、应用层、基础设施层、接口层分离
- [x] **Docker 容器化** - 一键启动开发/生产环境

### 开发中 (P1) 🚧
- [ ] **小说解析引擎** - 自动分段、角色识别、对话提取
- [ ] **角色一致性方案** - 参考图生成 + 图生图工作流（详见 [CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)）
- [ ] **场景智能分割** - 基于情节变化、对话密度、时空转换的自动分镜
- [ ] **AI 服务集成** - Gemini 2.5 Flash (图像) + Sora2 (视频)

### 计划中 (P2) 📋
- [ ] **生成结果下载与分享** - 支持 MP4/GIF 导出和社交媒体分享链接
- [ ] **质量评估系统** - 用户评分 + 自动质量检测（CLIP 相似度、网格质量）
- [ ] **角色库与复用** - 提取常见角色风格，支持跨作品复用
- [ ] **多语言配音** - 支持中文、英文、日文等多语言 TTS

### 展望 (P3) 🌟
- [ ] **实时协作编辑** - 多用户同时编辑场景、调整参数
- [ ] **提示词模板库** - 预设风格模板（少年漫、少女漫、奇幻等）
- [ ] **后处理工具链** - 视频剪辑、特效添加、字幕生成
- [ ] **管理后台** - 指标监控、用户管理、成本分析

**本次开发聚焦 P0 和 P1**：工作流系统、小说解析、角色一致性、AI 服务集成。

---

## 3) AI 服务选择

### 采用技术栈

#### **Gemini 2.5 Flash Image** - 图像生成
- **用途**：文生图、图生图
- **使用场景**：
  - 角色参考图生成（文本描述 → 角色外观）
  - 场景图像生成（带角色参考的图生图，保持一致性）
- **优势**：质量稳定、速度快、支持图像条件输入、API 简洁
- **集成位置**：`backend/internal/infrastructure/ai/gemini/`

#### **Sora2** - 视频生成
- **用途**：文生视频、图生视频
- **使用场景**：
  - 场景动画生成（静态图 → 短视频）
  - 动态内容生成（文本描述 → 动画片段）
- **优势**：视频质量高、支持时长控制、可与静态图结合
- **集成位置**：`backend/internal/infrastructure/ai/sora/`

### 技术选型对比

| 服务 | 文生图 | 图生图 | 文生视频 | 质量 | 速度 | 价格 | 角色一致性 |
|-----|--------|--------|----------|------|------|------|-----------|
| **Gemini 2.5 Flash** | ✅ | ✅ | ❌ | 高 | 快 | 中 | ⭐⭐⭐⭐⭐ |
| **Sora2** | ❌ | ❌ | ✅ | 极高 | 中 | 高 | ⭐⭐⭐⭐ |
| Midjourney | ✅ | ❌ | ❌ | 极高 | 慢 | 高 | ⭐⭐⭐ |
| Stable Diffusion | ✅ | ✅ | ❌ | 中 | 快 | 低 | ⭐⭐⭐ |

**选型理由**：
- **Gemini** 对图像条件输入支持好，适合角色一致性方案（参考图 + 提示词）
- **Sora2** 视频质量领先，适合高质量动漫输出
- 两者 API 稳定，有完善的 SDK，适合产品化

### 角色一致性实现

核心策略：**参考图生成 + 图生图转换**

```
文本描述 ──→ Gemini ──→ 角色参考图
                              ↓
场景文本 + 参考图 ──→ Gemini ──→ 场景图（角色一致）
```

详细设计请查看 [CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)。

---

## 4) 效果评估与持续迭代

### 关键指标

#### **生成质量指标**
- **成功率**：任务完成率（目标 >95%）
- **角色一致性得分**：CLIP 相似度对比（目标 >0.85）
- **用户满意度**：5 星评分系统（目标平均 >4.0）

#### **性能指标**
- **平均生成时长**：
  - 角色参考图：<30 秒
  - 单场景图像：<60 秒
  - 单场景视频：<120 秒
- **并发处理能力**：同时处理任务数（目标 >10）

#### **成本指标**
- **每章节成本**：API 调用费用（目标 <$0.50）
- **存储成本**：媒体文件存储费用

### 评估系统设计

```go
// backend/internal/application/service/quality_service.go
type QualityMetrics struct {
    ConsistencyScore  float64  // 角色一致性得分 (0-1)
    SceneRelevance    float64  // 场景相关性得分 (0-1)
    UserRating        float64  // 用户评分 (1-5)
    GenerationTime    int      // 生成时长（秒）
    ErrorRate         float64  // 错误率
}
```

#### **自动评估**
- **图像质量检测**：分辨率、清晰度、色彩饱和度
- **角色一致性检测**：CLIP Embedding 余弦相似度
- **场景相关性检测**：文本-图像匹配度（CLIP Score）

#### **用户反馈**
- **5 星评分系统**
- **问题标签**：角色不一致、场景不符、质量差、其他
- **文字评论**：详细问题描述

#### **策略优化**
- **参数调优**：根据评分调整提示词模板、生成参数
- **失败重试**：低分任务自动重新生成（最多 3 次）
- **黑名单机制**：记录失败案例，避免重复错误

### 迭代计划

1. **V1.0（当前）**：基础工作流 + 手动评估
2. **V1.1**：自动质量检测 + 用户评分系统
3. **V1.2**：CLIP 一致性检测 + 自动重试
4. **V2.0**：多供应商 AB 测试 + 智能路由

---

## 5) 如何降低 AI 成本

### 成本优化策略

#### **1. 参考图复用** 💰
- **策略**：同一角色只生成一次参考图，后续场景复用
- **节省**：每角色节省 N-1 次图像生成（N = 场景数）
- **实现**：CharacterRepository 缓存参考图 URL

#### **2. 场景批处理** 🚀
- **策略**：相似场景合并生成（同一地点、相同角色组合）
- **节省**：减少重复背景生成
- **实现**：SceneService 场景聚类算法

#### **3. 质量分级生成** 📊
- **策略**：
  - 预览模式：低分辨率（512x512），快速验证
  - 精修模式：高分辨率（1024x1024），最终输出
- **节省**：试错成本降低 50%
- **实现**：GenerationService 支持质量档位

#### **4. 失败早停** ⚠️
- **策略**：检测到低质量结果立即停止后续步骤
- **节省**：避免在错误基础上继续浪费 API 调用
- **实现**：QualityService 实时质量检查

#### **5. 提示词优化** ✍️
- **策略**：维护高质量提示词模板库，减少重试
- **节省**：提高首次成功率，减少重试次数
- **实现**：PromptService 模板管理

### 成本监控

```go
// backend/internal/domain/entity/cost_tracking.go
type CostMetrics struct {
    APICallCount      int       // API 调用次数
    TotalCost         float64   // 总成本（美元）
    CostPerChapter    float64   // 每章节成本
    CacheHitRate      float64   // 缓存命中率
    RetryRate         float64   // 重试率
}
```

---

## 快速开始

### 前置依赖
- **Node.js 18+** - 前端开发
- **Go 1.24+** - 后端开发
- **Docker & Docker Compose** - 容器化部署
- **Supabase 账号** - 数据库和存储服务
- **Gemini API Key** - AI 图像生成（开发期可先不配置）

### 环境变量配置

创建 `.env` 文件（根目录）：

```env
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=your-service-key
SUPABASE_JWT_SECRET=your-jwt-secret

# AI Services (开发期可选)
GEMINI_API_KEY=your-gemini-key
SORA_API_KEY=your-sora-key

# Backend
GO_ENV=development
PORT=8080

# Frontend (frontend/.env)
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
VITE_API_BASE_URL=http://localhost:8080
```

### Docker 部署（推荐）⭐

```bash
# 1. 克隆项目
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion

# 2. 配置环境变量
cp .env.example .env
vim .env  # 填入 Supabase 配置

# 3. 启动服务
docker-compose up -d

# 4. 访问应用
# 前端: http://localhost:3000
# 后端 API: http://localhost:8080
# 健康检查: http://localhost:8080/health

# 5. 查看日志
docker-compose logs -f

# 6. 停止服务
docker-compose down
```

### 本地开发

```bash
# 安装依赖
pnpm install             # 前端依赖
cd backend && go mod download  # 后端依赖

# 启动开发环境
make dev                 # 同时启动前后端
# 或分别启动：
make dev-backend         # 后端 :8080
make dev-frontend        # 前端 :3000

# 运行测试
make test                # 全部测试
make test-backend        # 后端测试
make test-frontend       # 前端测试

# 代码检查
make lint                # 代码规范检查
make format              # 代码格式化
```

详细指南请查看 [QUICKSTART.md](QUICKSTART.md)。

---

## 产品与技术要点

### 核心特性总结

✅ **异步工作流系统** - 小说上传 → 任务创建 → 状态追踪 → 结果查看
✅ **DDD 分层架构** - 领域层/应用层/基础设施层/接口层清晰分离
✅ **角色一致性方案** - 参考图生成 + 图生图，确保视觉统一
✅ **质量评估闭环** - 用户评分 + 自动检测，驱动策略优化
✅ **成本优化策略** - 参考图复用、场景批处理、质量分级
✅ **Docker 容器化** - 一键启动开发/生产环境

### 技术亮点

- **Go DDD 架构**：严格依赖倒置，领域层零外部依赖，易测试易扩展
- **Supabase 全栈**：认证、数据库、存储统一管理，开发效率高
- **AI 服务解耦**：基础设施层实现，可灵活切换供应商
- **前端 React 19**：并发特性提升 UI 响应速度
- **异步任务系统**：支持长时间 AI 生成，用户体验流畅

---

## 开发者参考（关键文件）

### 后端核心
- **工作流服务**：[backend/internal/application/service/manga_workflow_service.go](backend/internal/application/service/manga_workflow_service.go)
- **领域实体**：[backend/internal/domain/entity/](backend/internal/domain/entity/)
- **仓储接口**：[backend/internal/domain/repository/](backend/internal/domain/repository/)
- **Supabase 实现**：[backend/internal/infrastructure/repository/supabase/](backend/internal/infrastructure/repository/supabase/)
- **HTTP 处理器**：[backend/internal/interfaces/http/handler/manga_workflow_handler.go](backend/internal/interfaces/http/handler/manga_workflow_handler.go)
- **认证中间件**：[backend/internal/infrastructure/middleware/auth_middleware.go](backend/internal/infrastructure/middleware/auth_middleware.go)

### 前端核心
- **工作流页面**：[frontend/src/pages/Dashboard.tsx](frontend/src/pages/Dashboard.tsx)
- **任务详情**：[frontend/src/pages/WorkflowDetail.tsx](frontend/src/pages/WorkflowDetail.tsx)
- **API 服务**：[frontend/src/services/api.ts](frontend/src/services/api.ts)
- **Supabase 客户端**：[frontend/src/lib/supabase.ts](frontend/src/lib/supabase.ts)

### 配置与文档
- **架构设计**：[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- **角色一致性**：[docs/CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md)
- **API 文档**：[docs/API.md](docs/API.md)
- **开发指南**：[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
- **Docker 配置**：[docker-compose.yml](docker-compose.yml)

---

## 文档导航

📖 **新手指南**
- [QUICKSTART.md](QUICKSTART.md) - 5 分钟快速开始
- [README.md](README.md) - 项目概览（本文档）

🏗️ **架构与设计**
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - 完整架构设计和 DDD 分层说明
- [CHARACTER_CONSISTENCY.md](docs/CHARACTER_CONSISTENCY.md) - 角色一致性技术方案
- [API.md](docs/API.md) - RESTful API 接口规范

👨‍💻 **开发者文档**
- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - 开发规范、编码标准、测试策略
- [CONFIGURATION.md](docs/CONFIGURATION.md) - 环境变量和配置说明
- [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) - 常见问题与解决方案

🚀 **部署运维**
- [DEPLOYMENT.md](docs/DEPLOYMENT.md) - 生产环境部署指南
- [docker-compose.yml](docker-compose.yml) - Docker 编排配置

---

## 常用命令

```bash
# 查看帮助
make help

# 开发相关
make install         # 安装依赖
make dev             # 启动开发环境（前端 + 后端）
make dev-backend     # 仅启动后端 (:8080)
make dev-frontend    # 仅启动前端 (:3000)

# Docker 相关
make docker-up       # 启动 Docker 服务
make docker-down     # 停止 Docker 服务
make docker-rebuild  # 重新构建并启动
make docker-logs     # 查看容器日志

# 构建与测试
make build           # 编译前后端
make test            # 运行所有测试
make test-backend    # 后端测试
make test-frontend   # 前端测试
make lint            # 代码检查
make format          # 代码格式化

# 清理
make clean           # 清理构建文件
```

---

## 项目状态

**当前版本**: v0.1.0-alpha (开发中)

### 完成情况

✅ **已完成 (P0)**
- 项目初始化与容器化部署
- DDD 架构搭建（领域层/应用层/基础设施层/接口层）
- Supabase 集成（认证/数据库/存储）
- 异步工作流系统（任务创建/状态追踪）
- 前端工作流界面（Dashboard + 任务详情）

🚧 **开发中 (P1)**
- 小说解析引擎（分段/角色识别/对话提取）
- 角色一致性方案（参考图生成 + 图生图）
- AI 服务集成（Gemini 2.5 Flash + Sora2）
- 场景智能分割算法

📋 **计划中 (P2)**
- 质量评估系统（用户评分 + 自动检测）
- 生成结果下载与分享
- 角色库与复用机制
- 多语言配音支持

### 开发路线图

**Phase 1: MVP 核心功能（当前阶段）**
- [x] 基础架构与 Docker 环境
- [x] Supabase 认证与数据库
- [x] 异步任务工作流
- [ ] 小说解析引擎
- [ ] AI 服务集成（Gemini + Sora）
- [ ] 角色一致性实现

**Phase 2: 产品化与优化**
- [ ] 质量评估与反馈系统
- [ ] 成本优化策略（缓存/批处理）
- [ ] 前端完整用户界面
- [ ] 用户权限与配额管理

**Phase 3: 规模化与商业化**
- [ ] 多供应商智能路由
- [ ] 实时协作编辑
- [ ] 管理后台与数据分析
- [ ] 生产环境性能优化

---

## 贡献指南

欢迎贡献代码、报告问题或提出建议！

### 参与方式

1. **报告 Bug** - 提交 [Issue](https://github.com/xiajiayi/ai-motion/issues)
2. **功能建议** - 在 [讨论区](https://github.com/xiajiayi/ai-motion/discussions) 讨论
3. **贡献代码** - 提交 Pull Request

### 贡献流程

```bash
# 1. Fork 项目
git clone https://github.com/your-username/ai-motion.git
cd ai-motion

# 2. 创建特性分支
git checkout -b feature/amazing-feature

# 3. 开发并测试
make dev
make test
make lint

# 4. 提交更改（遵循 Conventional Commits）
git commit -m 'feat: add amazing feature'

# 5. 推送分支
git push origin feature/amazing-feature

# 6. 创建 Pull Request
```

### Commit 规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 构建/工具链更新

详细贡献指南请查看 [DEVELOPMENT.md](docs/DEVELOPMENT.md)。

---

## 安全与合规

### 安全最佳实践

- ✅ **JWT 认证**：所有 API 请求需要有效 JWT Token
- ✅ **环境变量**：敏感信息通过环境变量注入，不提交代码
- ✅ **CORS 配置**：生产环境限制允许的域名
- ✅ **输入验证**：所有用户输入进行验证和清理
- ⚠️ **API Key 管理**：建议使用密钥管理服务（AWS Secrets Manager 等）

### 生产环境建议

- 启用 HTTPS（通过反向代理如 Nginx）
- 配置速率限制防止 API 滥用
- 定期备份 Supabase 数据库
- 监控 AI 服务调用成本
- 设置日志告警（错误率/响应时间）

---

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

## 联系与支持

### 获取帮助

- 📖 **查看文档** - 先查看 [docs/](docs/) 目录下的文档
- 💬 **讨论区** - [GitHub Discussions](https://github.com/xiajiayi/ai-motion/discussions)
- 🐛 **报告 Bug** - [GitHub Issues](https://github.com/xiajiayi/ai-motion/issues)

### 社区

- ⭐ **Star 项目** - 如果觉得有帮助，欢迎 Star 支持！
- 🔄 **分享** - 将项目分享给需要的朋友
- 💡 **反馈** - 您的反馈是我们改进的动力

---

**Built with ❤️ by AI-Motion Team**
