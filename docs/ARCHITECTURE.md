# AI-Motion 系统架构文档

## 架构概述

AI-Motion 采用 **DDD (领域驱动设计)** 架构模式,将系统按照业务领域划分,实现高内聚、低耦合的设计。

### 核心理念

- **以领域为中心**: 围绕小说、角色、场景、媒体等核心业务概念建模
- **分层架构**: 清晰的层次划分,职责分离
- **依赖倒置**: 外层依赖内层,核心业务逻辑不依赖外部技术实现

## 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         用户界面层                            │
│                    (React + TypeScript)                     │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP/REST
┌───────────────────────────▼─────────────────────────────────┐
│                      接口层 (Interfaces)                      │
│              HTTP Handlers + Middleware                      │
│              - 请求验证                                        │
│              - 响应格式化                                      │
│              - 错误处理                                        │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                     应用层 (Application)                      │
│                      应用服务 + DTO                           │
│              - 用例编排                                        │
│              - 事务管理                                        │
│              - 业务流程协调                                    │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      领域层 (Domain)                          │
│                   领域模型 + 业务规则                          │
│   ┌──────────┬──────────┬──────────┬──────────┐            │
│   │Novel领域 │Character │Scene领域 │Media领域 │            │
│   │          │领域      │          │          │            │
│   └──────────┴──────────┴──────────┴──────────┘            │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                  基础设施层 (Infrastructure)                  │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │Repository│  │AI Services│  │Storage   │  │External  │   │
│  │ Supabase │  │ Gemini   │  │  MinIO   │  │  APIs    │   │
│  │(PostgREST)│  │  Sora2   │  │          │  │          │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## DDD 分层架构

### 1. 领域层 (Domain Layer)

**职责**: 核心业务逻辑和业务规则

**目录结构**:
```
internal/domain/
├── novel/              # 小说领域
│   ├── entity.go      # 实体: Novel
│   ├── value_object.go # 值对象: Chapter, Paragraph
│   ├── repository.go  # 仓储接口
│   └── service.go     # 领域服务
├── character/         # 角色领域
│   ├── entity.go      # 实体: Character
│   ├── value_object.go # 值对象: Appearance, Personality
│   ├── repository.go  # 仓储接口
│   └── service.go     # 领域服务
├── scene/             # 场景领域
│   ├── entity.go      # 实体: Scene
│   ├── value_object.go # 值对象: Dialogue, Description
│   ├── repository.go  # 仓储接口
│   └── service.go     # 领域服务
└── media/             # 媒体领域
    ├── entity.go      # 实体: Image, Video
    ├── value_object.go # 值对象: MediaMetadata, GenerationParams
    ├── repository.go  # 仓储接口
    └── service.go     # 领域服务
```

**核心实体**:

#### Novel (小说实体)
```go
type Novel struct {
    ID          NovelID
    Title       string
    Author      string
    Content     string
    Chapters    []Chapter
    Status      NovelStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 业务方法
func (n *Novel) Parse() error
func (n *Novel) ExtractCharacters() []Character
func (n *Novel) DivideIntoScenes() []Scene
```

#### Character (角色实体)
```go
type Character struct {
    ID              CharacterID
    NovelID         NovelID
    Name            string
    Appearance      Appearance      // 值对象
    Personality     Personality     // 值对象
    ReferenceImages []ImageID
    CreatedAt       time.Time
}

// 业务方法
func (c *Character) GenerateReferenceImage(prompt string) (*Image, error)
func (c *Character) ValidateConsistency(newImage *Image) (bool, float64)
```

#### Scene (场景实体)
```go
type Scene struct {
    ID          SceneID
    NovelID     NovelID
    ChapterID   ChapterID
    SequenceNum int
    Description Description  // 值对象
    Dialogues   []Dialogue   // 值对象
    Characters  []CharacterID
    Location    string
    TimeOfDay   string
}

// 业务方法
func (s *Scene) GeneratePrompt() string
func (s *Scene) CalculateDuration() time.Duration
```

#### Media (媒体实体)
```go
type Media struct {
    ID          MediaID
    Type        MediaType    // Image, Video
    SceneID     SceneID
    URL         string
    Metadata    MediaMetadata // 值对象
    Status      MediaStatus
    CreatedAt   time.Time
}

// 业务方法
func (m *Media) IsReady() bool
func (m *Media) GetDimensions() (width, height int)
```

### 2. 应用层 (Application Layer)

**职责**: 用例编排,协调领域对象完成业务流程

**目录结构**:
```
internal/application/
├── service/
│   ├── novel_service.go      # 小说相关用例
│   ├── character_service.go  # 角色相关用例
│   ├── generation_service.go # 内容生成用例
│   └── export_service.go     # 导出用例
└── dto/
    ├── novel_dto.go
    ├── character_dto.go
    ├── scene_dto.go
    └── media_dto.go
```

**应用服务示例**:

```go
// NovelService - 小说应用服务
type NovelService struct {
    novelRepo      domain.NovelRepository
    characterRepo  domain.CharacterRepository
    sceneRepo      domain.SceneRepository
    parserService  domain.NovelParserService
}

func (s *NovelService) UploadAndParse(ctx context.Context, req *UploadNovelRequest) (*NovelDTO, error) {
    // 1. 创建小说实体
    novel := domain.NewNovel(req.Title, req.Author, req.Content)

    // 2. 调用领域服务解析小说
    if err := s.parserService.Parse(novel); err != nil {
        return nil, err
    }

    // 3. 提取角色
    characters := novel.ExtractCharacters()

    // 4. 持久化
    if err := s.novelRepo.Save(ctx, novel); err != nil {
        return nil, err
    }

    for _, char := range characters {
        if err := s.characterRepo.Save(ctx, char); err != nil {
            return nil, err
        }
    }

    // 5. 返回 DTO
    return s.toDTO(novel), nil
}
```

```go
// GenerationService - 内容生成应用服务
type GenerationService struct {
    sceneRepo     domain.SceneRepository
    mediaRepo     domain.MediaRepository
    characterRepo domain.CharacterRepository
    geminiClient  *gemini.Client
    soraClient    *sora.Client
}

func (s *GenerationService) GenerateSceneImage(ctx context.Context, sceneID string) (*MediaDTO, error) {
    // 1. 获取场景实体
    scene, err := s.sceneRepo.FindByID(ctx, domain.SceneID(sceneID))
    if err != nil {
        return nil, err
    }

    // 2. 获取场景中的角色
    characters, err := s.characterRepo.FindByIDs(ctx, scene.Characters)
    if err != nil {
        return nil, err
    }

    // 3. 生成提示词
    prompt := s.buildPrompt(scene, characters)

    // 4. 调用 Gemini API 生成图片
    imageURL, err := s.geminiClient.TextToImage(ctx, prompt)
    if err != nil {
        return nil, err
    }

    // 5. 创建媒体实体
    media := domain.NewMedia(domain.MediaTypeImage, scene.ID, imageURL)

    // 6. 持久化
    if err := s.mediaRepo.Save(ctx, media); err != nil {
        return nil, err
    }

    return s.toMediaDTO(media), nil
}

func (s *GenerationService) GenerateSceneVideo(ctx context.Context, sceneID string) (*MediaDTO, error) {
    // 1. 获取场景图片
    images, err := s.mediaRepo.FindBySceneID(ctx, domain.SceneID(sceneID))
    if err != nil {
        return nil, err
    }

    // 2. 调用 Sora2 API 图生视频
    videoURL, err := s.soraClient.ImageToVideo(ctx, images[0].URL)
    if err != nil {
        return nil, err
    }

    // 3. 创建视频媒体实体
    media := domain.NewMedia(domain.MediaTypeVideo, domain.SceneID(sceneID), videoURL)

    // 4. 持久化
    if err := s.mediaRepo.Save(ctx, media); err != nil {
        return nil, err
    }

    return s.toMediaDTO(media), nil
}
```

### 3. 接口层 (Interfaces Layer)

**职责**: 处理外部请求,调用应用服务

**目录结构**:
```
internal/interfaces/
├── http/
│   ├── handler/
│   │   ├── novel_handler.go
│   │   ├── character_handler.go
│   │   ├── scene_handler.go
│   │   └── generation_handler.go
│   ├── request/
│   │   └── validators.go
│   └── response/
│       └── formatters.go
└── middleware/
    ├── auth.go
    ├── cors.go
    ├── logger.go
    └── error_handler.go
```

**HTTP Handler 示例**:

```go
type NovelHandler struct {
    novelService *application.NovelService
}

func (h *NovelHandler) Upload(c *gin.Context) {
    var req dto.UploadNovelRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    novel, err := h.novelService.UploadAndParse(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": novel})
}
```

### 4. 基础设施层 (Infrastructure Layer)

**职责**: 提供技术实现,如数据库访问、外部 API 调用等

**目录结构**:
```
internal/infrastructure/
├── repository/
│   ├── supabase/
│   │   ├── novel_repository.go
│   │   ├── character_repository.go
│   │   ├── scene_repository.go
│   │   └── media_repository.go
│   └── memory/          # 内存实现(用于测试)
├── ai/
│   ├── gemini/
│   │   ├── client.go
│   │   ├── text_to_image.go
│   │   └── image_to_image.go
│   └── sora/
│       ├── client.go
│       ├── text_to_video.go
│       └── image_to_video.go
├── storage/
│   ├── local/
│   │   └── file_storage.go
│   └── minio/
│       └── object_storage.go
└── config/
    └── config.go
```

**仓储实现示例**:

```go
type SupabaseNovelRepository struct {
    client *postgrest.Client
}

func (r *SupabaseNovelRepository) Save(ctx context.Context, novel *domain.Novel) error {
    _, _, err := r.client.From("aimotion_novel").
        Upsert(novel, "", "", "*").
        Execute()
    return err
}

func (r *SupabaseNovelRepository) FindByID(ctx context.Context, id domain.NovelID) (*domain.Novel, error) {
    var novels []domain.Novel
    _, _, err := r.client.From("aimotion_novel").
        Select("*", "exact", false).
        Eq("id", string(id)).
        Single().
        Execute(&novels)
    
    if err != nil {
        return nil, err
    }
    
    if len(novels) == 0 {
        return nil, domain.ErrNovelNotFound
    }

    return &novels[0], nil
}
```

**AI 服务客户端示例**:

```go
// Gemini 客户端
type GeminiClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func (c *GeminiClient) TextToImage(ctx context.Context, prompt string) (string, error) {
    req := &TextToImageRequest{
        Prompt: prompt,
        Model:  "gemini-2.5-flash-image",
        Size:   "1024x768",
    }

    resp, err := c.post(ctx, "/v1/images/generations", req)
    if err != nil {
        return "", err
    }

    return resp.ImageURL, nil
}

func (c *GeminiClient) ImageToImage(ctx context.Context, imageURL string, prompt string) (string, error) {
    req := &ImageToImageRequest{
        ImageURL: imageURL,
        Prompt:   prompt,
        Model:    "gemini-2.5-flash-image",
    }

    resp, err := c.post(ctx, "/v1/images/variations", req)
    if err != nil {
        return "", err
    }

    return resp.ImageURL, nil
}
```

```go
// Sora2 客户端
type Sora2Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func (c *Sora2Client) TextToVideo(ctx context.Context, prompt string) (string, error) {
    req := &TextToVideoRequest{
        Prompt:   prompt,
        Duration: 5,  // 5秒视频
        FPS:      30,
    }

    resp, err := c.post(ctx, "/v1/videos/generations", req)
    if err != nil {
        return "", err
    }

    return resp.VideoURL, nil
}

func (c *Sora2Client) ImageToVideo(ctx context.Context, imageURL string) (string, error) {
    req := &ImageToVideoRequest{
        ImageURL: imageURL,
        Duration: 5,
        FPS:      30,
    }

    resp, err := c.post(ctx, "/v1/videos/animations", req)
    if err != nil {
        return "", err
    }

    return resp.VideoURL, nil
}
```

## 核心业务流程

### 1. 小说上传与解析流程

```
用户上传小说
    ↓
NovelHandler.Upload()
    ↓
NovelService.UploadAndParse()
    ↓
创建 Novel 实体
    ↓
调用 NovelParserService 解析
    ↓
提取角色 → 创建 Character 实体
    ↓
划分场景 → 创建 Scene 实体
    ↓
持久化到数据库
    ↓
返回解析结果
```

### 2. 角色一致性生成流程

```
创建角色
    ↓
生成角色参考图 (文生图)
    ↓
GenerationService.GenerateCharacterReference()
    ↓
Gemini TextToImage API
    ↓
保存参考图到 Character.ReferenceImages
    ↓
后续场景生成时
    ↓
使用参考图 + 场景描述
    ↓
Gemini ImageToImage API (图生图)
    ↓
保持角色一致性
```

### 3. 场景内容生成流程

```
选择场景
    ↓
GenerationService.GenerateSceneImage()
    ↓
获取场景中的角色
    ↓
获取角色参考图
    ↓
构建提示词 (场景描述 + 角色特征)
    ↓
Gemini ImageToImage API
    ↓
生成场景图片
    ↓
保存 Media 实体
    ↓
(可选) GenerateSceneVideo()
    ↓
Sora2 ImageToVideo API
    ↓
生成场景视频
```

## 数据模型

### 数据库 Schema

> 数据库设计遵循《阿里巴巴Java开发手册》规范,详细规范请参考 [DATABASE_DESIGN_STANDARDS.md](DATABASE_DESIGN_STANDARDS.md)

数据库表结构已拆分为独立的 SQL 文件,每个表一个文件,便于管理和维护。

**Schema 文件位置**: [`database/schema/`](../database/schema/)

**核心表结构**:

| 表名 | 说明 | Schema 文件 |
|------|------|-------------|
| `aimotion_novel` | 小说表 | [01_aimotion_novel.sql](../database/schema/01_aimotion_novel.sql) |
| `aimotion_novel_content` | 小说内容表(大字段独立) | [02_aimotion_novel_content.sql](../database/schema/02_aimotion_novel_content.sql) |
| `aimotion_chapter` | 章节表 | [03_aimotion_chapter.sql](../database/schema/03_aimotion_chapter.sql) |
| `aimotion_character` | 角色表 | [04_aimotion_character.sql](../database/schema/04_aimotion_character.sql) |
| `aimotion_character_image` | 角色图片表 | [05_aimotion_character_image.sql](../database/schema/05_aimotion_character_image.sql) |
| `aimotion_scene` | 场景表 | [06_aimotion_scene.sql](../database/schema/06_aimotion_scene.sql) |
| `aimotion_scene_character` | 场景角色关联表 | [07_aimotion_scene_character.sql](../database/schema/07_aimotion_scene_character.sql) |
| `aimotion_media` | 媒体表 | [08_aimotion_media.sql](../database/schema/08_aimotion_media.sql) |

**快速初始化**: 使用 [`database/schema/init.sql`](../database/schema/init.sql) 可一次性创建所有表。

**关键设计规范**:

1. **表名规范**: 统一使用 `aimotion_` 前缀
2. **ID 类型**: 使用 `BIGINT UNSIGNED AUTO_INCREMENT`
3. **时间字段**: 使用 `gmt_create` 和 `gmt_modified`
4. **逻辑删除**: 新增 `is_deleted` 字段
5. **索引命名**: 使用 `uk_` (唯一索引) 和 `idx_` (普通索引) 前缀
6. **移除外键**: 不使用 `FOREIGN KEY`,应用层维护关联
7. **垂直拆分**: 将小说内容独立为 `aimotion_novel_content` 表
8. **新增关联表**: `aimotion_character_image` (角色图片) 和 `aimotion_scene_character` (场景角色关联)

完整的表关系图和使用说明请参考 [`database/schema/README.md`](../database/schema/README.md)

## AI 服务集成策略

### Gemini 2.5 Flash Image

**用途**:
- 文生图: 生成角色参考图、场景初始图
- 图生图: 基于参考图生成一致性角色场景

**API 调用示例**:
```go
// 文生图
imageURL := geminiClient.TextToImage(ctx, "一位穿着白色长袍的年轻女子,黑色长发,站在古代庭院中")

// 图生图 (保持角色一致性)
consistentImageURL := geminiClient.ImageToImage(ctx, referenceImageURL, "同样的女子,现在在竹林中")
```

### Sora2

**用途**:
- 文生视频: 直接从文本生成动态场景
- 图生视频: 将静态场景图片动态化

**API 调用示例**:
```go
// 文生视频
videoURL := soraClient.TextToVideo(ctx, "年轻女子在竹林中行走,微风吹动竹叶")

// 图生视频
animatedVideoURL := soraClient.ImageToVideo(ctx, sceneImageURL)
```

## 性能优化

### 1. 异步处理
- 使用 Go 协程处理耗时的 AI 生成任务
- 任务队列管理 (可选 Redis + Worker)

### 2. 缓存策略
- 角色参考图缓存
- 生成结果缓存
- API 响应缓存

### 3. 并发控制
- 限制并发 AI API 调用数量
- 使用 semaphore 控制并发

## 扩展性设计

### 1. 策略模式
- 支持多个 AI 服务提供商
- 可插拔的生成策略

### 2. 事件驱动
- 领域事件发布/订阅
- 异步解耦

### 3. 微服务化
- 各领域可独立部署
- API Gateway 统一入口

## 安全性

- API Key 安全管理
- JWT 认证
- 请求限流
- 文件上传验证
- SQL 注入防护

## 监控与日志

- 结构化日志 (使用 zap/logrus)
- 错误追踪
- 性能监控
- AI API 调用统计
