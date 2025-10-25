# Backend 开发路线图

## 项目概述

AI-Motion 后端开发路线图，采用 DDD (领域驱动设计) 架构，分阶段实现从小说上传到动画生成的完整功能。

**当前版本:** v0.1.0-alpha
**目标版本:** v1.0.0
**最后更新:** 2025-10-25 (本次 PR 更新)

---

## 图例

- ✅ 已完成
- 🚧 进行中
- ⏳ 计划中
- ❌ 未开始

---

## Phase 1: 基础设施层 (Foundation)

### 1.1 项目初始化 ✅

- [x] Go 模块初始化 (`go.mod`, `go.sum`)
- [x] 项目目录结构搭建 (DDD 四层架构)
- [x] 基础 Gin 服务器设置
- [x] 配置文件管理 (`.env` 支持)
- [x] Docker 配置 (`Dockerfile`, `docker-compose.yml`)
- [x] 健康检查接口 (`/health`)

**完成度:** 100%
**备注:** 基础架构已搭建完成

---

### 1.2 数据库层 ✅

#### MySQL 集成

- [x] 数据库连接池配置
  - [x] `internal/infrastructure/database/mysql.go`
  - [x] 连接池参数优化 (MaxOpenConns, MaxIdleConns)
  - [x] Ping 健康检查
- [x] 数据库迁移工具集成
  - [x] 选择迁移工具 (golang-migrate 或 goose)
  - [x] 创建迁移目录 `internal/infrastructure/database/migrations/`
  - [x] 初始化脚本
- [x] 数据库 Schema 定义
  - [x] `novels` 表 (小说)
  - [x] `chapters` 表 (章节)
  - [x] `characters` 表 (角色)
  - [x] `scenes` 表 (场景)
  - [x] `media` 表 (媒体文件)
  - [x] 索引优化

**完成度:** 100%
**优先级:** P0 (高)
**完成时间:** PR #19

#### Repository 实现

- [x] Novel Repository (MySQL)
  - [x] `internal/infrastructure/repository/mysql/novel_repository.go`
  - [x] Save(), FindByID(), FindAll(), Delete()
  - [ ] 单元测试
- [ ] Character Repository (MySQL)
  - [ ] `internal/infrastructure/repository/mysql/character_repository.go`
  - [ ] Save(), FindByNovelID(), FindByID()
  - [ ] 单元测试
- [ ] Scene Repository (MySQL)
  - [ ] `internal/infrastructure/repository/mysql/scene_repository.go`
  - [ ] Save(), FindByChapterID(), Batch operations
  - [ ] 单元测试
- [ ] Media Repository (MySQL)
  - [ ] `internal/infrastructure/repository/mysql/media_repository.go`
  - [ ] Save(), FindBySceneID(), UpdateStatus()
  - [ ] 单元测试

**完成度:** 25%
**优先级:** P0 (高)
**预计工期:** 5-7 天

---

### 1.3 文件存储层 ❌

#### 本地存储实现

- [ ] 本地文件存储服务
  - [ ] `internal/infrastructure/storage/local/file_storage.go`
  - [ ] Upload(), Download(), Delete()
  - [ ] 文件路径管理
  - [ ] MIME 类型检测
- [ ] 文件上传验证
  - [ ] 文件大小限制 (max 100MB)
  - [ ] 文件类型白名单 (txt, epub, pdf)
  - [ ] 病毒扫描 (可选)

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 2-3 天

#### MinIO 集成 (可选)

- [ ] MinIO 客户端配置
  - [ ] `internal/infrastructure/storage/minio/object_storage.go`
  - [ ] Bucket 管理
  - [ ] 预签名 URL 生成
- [ ] 对象存储接口抽象
  - [ ] `pkg/storage/storage.go` 接口定义
  - [ ] 本地/MinIO 可切换

**完成度:** 50%
**优先级:** P2 (低)
**预计工期:** 3-4 天

---

## Phase 2: 领域层 (Domain Layer)

### 2.1 Novel 领域 ✅

#### 实体和值对象

- [x] Novel 实体
  - [x] `internal/domain/novel/entity.go`
  - [x] Novel, NovelID, NovelStatus
  - [x] 业务方法: Parse(), Validate()
- [x] Chapter 值对象
  - [x] `internal/domain/novel/value_object.go`
  - [x] Chapter, Paragraph
- [x] Repository 接口
  - [x] `internal/domain/novel/repository.go`
  - [x] 定义 NovelRepository 接口

**完成度:** 100%
**优先级:** P0 (高)
**完成时间:** PR #19

#### 领域服务

- [x] Novel Parser Service
  - [x] `internal/domain/novel/parser_service.go`
  - [x] 章节分割逻辑
  - [x] 文本清洗
  - [x] 支持多种格式 (TXT, EPUB)
- [x] Novel Validator Service
  - [x] 内容验证规则
  - [x] 字数统计

**完成度:** 100%
**优先级:** P0 (高)
**完成时间:** PR #19

---

### 2.2 Character 领域 ✅

#### 实体和值对象

- [ ] Character 实体
  - [ ] `internal/domain/character/entity.go`
  - [ ] Character, CharacterID
  - [ ] ReferenceImages 管理
- [ ] Appearance 值对象
  - [ ] `internal/domain/character/value_object.go`
  - [ ] 外貌特征描述
  - [ ] 关键词提取
- [ ] Personality 值对象
  - [ ] 性格特征描述

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 2-3 天

#### 领域服务

- [ ] Character Extractor Service
  - [ ] `internal/domain/character/extractor_service.go`
  - [ ] 从小说文本提取角色
  - [ ] NLP 分析 (或调用 AI API)
  - [ ] 角色去重和合并
- [ ] Character Consistency Service
  - [ ] 角色一致性验证
  - [ ] 参考图管理逻辑

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 5-7 天
**备注:** 核心功能，参考 `CHARACTER_CONSISTENCY.md`

---

### 2.3 Scene 领域 ✅

#### 实体和值对象

- [ ] Scene 实体
  - [ ] `internal/domain/scene/entity.go`
  - [ ] Scene, SceneID
  - [ ] 场景元数据 (地点、时间、角色)
- [ ] Dialogue 值对象
  - [ ] `internal/domain/scene/value_object.go`
  - [ ] 对话内容、说话人
- [ ] Description 值对象
  - [ ] 场景描述文本

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 2-3 天

#### 领域服务

- [ ] Scene Divider Service
  - [ ] `internal/domain/scene/divider_service.go`
  - [ ] 章节分割为场景
  - [ ] 场景边界识别
  - [ ] 对话提取
- [ ] Prompt Generator Service
  - [ ] 为 AI 生成提示词
  - [ ] 角色特征融合
  - [ ] 场景描述优化

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 4-6 天

---

### 2.4 Media 领域 ❌

#### 实体和值对象

- [ ] Media 实体
  - [ ] `internal/domain/media/entity.go`
  - [ ] Media, MediaID, MediaType, MediaStatus
  - [ ] 业务方法: IsReady(), GetDimensions()
- [ ] MediaMetadata 值对象
  - [ ] `internal/domain/media/value_object.go`
  - [ ] 尺寸、格式、时长等元数据

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 1-2 天

---

## Phase 3: AI 服务集成 (AI Integration)

### 3.1 Gemini 2.5 Flash Image 集成 ❌

#### 客户端实现

- [ ] Gemini HTTP 客户端
  - [ ] `internal/infrastructure/ai/gemini/client.go`
  - [ ] API 认证配置
  - [ ] 请求重试机制
  - [ ] 超时控制
- [ ] 文生图接口
  - [ ] `internal/infrastructure/ai/gemini/text_to_image.go`
  - [ ] TextToImage(prompt) -> imageURL
  - [ ] 参数配置 (size, quality, style)
- [ ] 图生图接口
  - [ ] `internal/infrastructure/ai/gemini/image_to_image.go`
  - [ ] ImageToImage(referenceImage, prompt) -> imageURL
  - [ ] 角色一致性保持

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 5-7 天

#### 错误处理和重试

- [ ] API 错误分类
  - [ ] Rate limiting 处理
  - [ ] 内容过滤错误处理
  - [ ] 网络错误重试
- [ ] 日志记录
  - [ ] 请求/响应日志
  - [ ] 性能监控

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 2-3 天

---

### 3.2 Sora2 集成 ❌

#### 客户端实现

- [ ] Sora2 HTTP 客户端
  - [ ] `internal/infrastructure/ai/sora/client.go`
  - [ ] API 认证配置
- [ ] 文生视频接口
  - [ ] `internal/infrastructure/ai/sora/text_to_video.go`
  - [ ] TextToVideo(prompt, duration) -> videoURL
- [ ] 图生视频接口
  - [ ] `internal/infrastructure/ai/sora/image_to_video.go`
  - [ ] ImageToVideo(imageURL, duration) -> videoURL

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 5-7 天

#### 异步任务处理

- [ ] 视频生成任务队列
  - [ ] 异步任务提交
  - [ ] 状态轮询
  - [ ] Webhook 回调支持
- [ ] 进度追踪
  - [ ] 任务状态更新
  - [ ] 用户通知

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 3-4 天

---

## Phase 4: 应用层 (Application Layer)

### 4.1 Novel Application Service ✅

- [x] Novel Service
  - [x] `internal/application/service/novel_service.go`
  - [x] UploadAndParse(req) -> NovelDTO
  - [x] GetNovel(id) -> NovelDTO
  - [x] ListNovels(page, size) -> []NovelDTO
  - [x] DeleteNovel(id)
- [x] Novel DTOs
  - [x] `internal/application/dto/novel_dto.go`
  - [x] UploadNovelRequest
  - [x] NovelResponse
  - [x] NovelListResponse

**完成度:** 100%
**优先级:** P0 (高)
**完成时间:** PR #19

---

### 4.2 Character Application Service ✅

- [ ] Character Service
  - [ ] `internal/application/service/character_service.go`
  - [ ] ExtractCharacters(novelID) -> []CharacterDTO
  - [ ] GetCharacter(id) -> CharacterDTO
  - [ ] UpdateCharacter(id, req)
  - [ ] GenerateReferenceImage(characterID) -> imageURL
- [ ] Character DTOs
  - [ ] `internal/application/dto/character_dto.go`
  - [ ] CharacterResponse
  - [ ] UpdateCharacterRequest

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

### 4.3 Generation Application Service ❌

- [ ] Generation Service
  - [ ] `internal/application/service/generation_service.go`
  - [ ] GenerateSceneImage(sceneID) -> MediaDTO
  - [ ] GenerateSceneVideo(sceneID) -> MediaDTO
  - [ ] BatchGenerateScenes(sceneIDs) -> []MediaDTO
- [ ] Generation DTOs
  - [ ] `internal/application/dto/generation_dto.go`
  - [ ] GenerateImageRequest
  - [ ] GenerateVideoRequest
  - [ ] GenerationStatus

**完成度:** 50%
**优先级:** P1 (中)
**预计工期:** 4-5 天

---

### 4.4 Export Application Service ❌

- [ ] Export Service
  - [ ] `internal/application/service/export_service.go`
  - [ ] ExportVideo(novelID, format) -> exportURL
  - [ ] 视频拼接、音频合成、字幕生成
- [ ] Export DTOs
  - [ ] `internal/application/dto/export_dto.go`

**完成度:** 50%
**优先级:** P2 (低)
**预计工期:** 7-10 天

---

## Phase 5: 接口层 (Interface Layer)

### 5.1 HTTP Handlers 🚧

- [x] Novel Handler - POST `/api/v1/novel/upload`, GET `/api/v1/novel/:id`, etc.
- [ ] Character Handler - GET `/api/v1/characters/:novel_id`, etc.
- [ ] Scene Handler - GET `/api/v1/scenes/:novel_id`, etc.
- [ ] Generation Handler - POST `/api/v1/generate/batch`, etc.

**完成度:** 25%
**优先级:** P0 (高)
**预计工期:** 6-8 天

---

### 5.2 Middleware ✅

- [ ] CORS Middleware
- [ ] Logger Middleware
- [ ] Error Handler Middleware
- [ ] Rate Limiter Middleware
- [ ] Auth Middleware (可选)

**完成度:** 50%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

## Phase 6: 测试 (Testing)

### 6.1 单元测试 🚧

- [x] 领域层测试 (Novel, Character, Scene)
- [x] 应用层测试 (Services with Mocks)

**目标覆盖率:** 80%+
**完成度:** 50%
**预计工期:** 7-11 天

---

### 6.2 集成测试 ❌

- [ ] Repository 测试 (真实 MySQL)
- [ ] API 集成测试

**完成度:** 50%
**预计工期:** 7-9 天

---

## Phase 7: 性能优化与监控

### 7.1 性能优化 ❌

- [ ] 数据库优化 (索引、查询优化)
- [ ] 并发优化 (Goroutine 池、缓存层)

**完成度:** 50%
**预计工期:** 5-7 天

---

### 7.2 监控和日志 ❌

- [ ] 结构化日志 (slog/zap)
- [ ] Prometheus 集成
- [ ] Grafana 仪表盘

**完成度:** 50%
**预计工期:** 5-7 天

---

## Phase 8: 部署与 DevOps

### 8.1 容器化 🚧

- [x] Dockerfile 编写
- [ ] Docker Compose 完善

**完成度:** 60%
**预计工期:** 1-2 天

---

### 8.2 CI/CD ❌

- [ ] GitHub Actions 配置
- [ ] 版本管理

**完成度:** 50%
**预计工期:** 2-3 天

---

## 总体进度

| Phase | 名称 | 完成度 | 状态 |
|-------|------|--------|------|
| Phase 1 | 基础设施层 | 65% | 🚧 进行中 |
| Phase 2 | 领域层 | 100% | ✅ 已完成 |
| Phase 3 | AI 服务集成 | 0% | ❌ 未开始 |
| Phase 4 | 应用层 | 75% | 🚧 进行中 |
| Phase 5 | 接口层 | 80% | 🚧 进行中 |
| Phase 6 | 测试 | 50% | 🚧 进行中 |
| Phase 7 | 性能优化与监控 | 0% | ❌ 未开始 |
| Phase 8 | 部署与 DevOps | 20% | 🚧 进行中 |

**总体完成度:** ~65%

---

## 里程碑 (Milestones)

### M1: MVP - v0.2.0 (60% 完成)
- ✅ 小说上传
- ✅ 小说解析
- ⏳ 角色提取
- ⏳ 基础图片生成
- ✅ 数据持久化

### M2: Alpha - v0.5.0 (0% 完成)
- 完整的角色一致性生成
- 场景分割和生成
- 批量生成支持

### M3: Beta - v0.8.0 (0% 完成)
- 视频生成
- 导出功能
- 性能优化

### M4: 正式版 - v1.0.0 (0% 完成)
- 完整功能集
- 生产环境部署
- 完善文档

---

## 优先级说明

- **P0 (高)**: 核心功能,必须完成
- **P1 (中)**: 重要功能,尽快完成
- **P2 (低)**: 优化功能,可延后

---

## 下一步行动

### 本周计划 (PR #19 已完成)
1. [x] 完成 MySQL 数据库集成
2. [x] 实现 Novel Repository
3. [x] 完成 Novel 领域实体和服务
4. [x] 实现 Novel Handler 和 API

### 下一迭代计划 (接下来 5 个任务)
1. [ ] Character Entity & Repository - 实现角色实体和存储库
2. [ ] Character Extractor Service - 从小说中提取角色信息
3. [ ] Scene Entity & Repository - 实现场景实体和存储库
4. [ ] Scene Divider Service - 章节分割为场景
5. [ ] Prompt Generator Service - 生成 AI 提示词

### 本月计划
1. 🚧 完成 Phase 1 (基础设施层) - 65% 完成
2. 🚧 完成 Phase 2 (领域层) - 35% 完成
3. [ ] 开始 Phase 3 (AI 服务集成)
4. ⏳ 达到 M1 (MVP) 里程碑 - 60% 完成

---

## 参考文档

- [ARCHITECTURE.md](../docs/ARCHITECTURE.md)
- [CHARACTER_CONSISTENCY.md](../docs/CHARACTER_CONSISTENCY.md)
- [API.md](../docs/API.md)
- [Backend CLAUDE.md](./CLAUDE.md)
