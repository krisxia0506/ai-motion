# 角色一致性设计方案

## 概述

角色一致性是 AI-Motion 系统的**核心功能**。它确保同一角色在所有生成的场景中保持一致的视觉外观,这对于从小说创建连贯的动画内容至关重要。

## 设计目标

1. **视觉一致性**: 同一角色在不同场景中外观完全一致
2. **可扩展性**: 支持每部小说中的多个角色
3. **高质量**: 生成高质量的参考图以确保准确复现
4. **高性能**: 高效的图像生成,避免冗余 API 调用

## 架构设计

### 1. 参考图生成

**触发时机**: 从小说中首次提取角色时

**处理流程**:
```
小说文本 → 角色提取 → 生成参考图 → 存储
```

**实现细节**:
- **服务**: Gemini 2.5 Flash Image (文本到图像)
- **输入**: 从小说中提取的角色描述
- **输出**: 高质量参考图
- **存储**: 保存到数据库 `Character.ReferenceImages` 字段

**示例代码**:
```go
// Domain 实体
type Character struct {
    ID              string
    NovelID         string
    Name            string
    Description     string
    ReferenceImages []string  // 参考图 URL 列表
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 2. 带一致性的场景生成

**触发时机**: 为场景生成图像时

**处理流程**:
```
场景描述 + 角色参考图 → 图像生成 → 场景图片
```

**实现细节**:
- **服务**: Gemini 2.5 Flash Image (图像到图像)
- **输入**:
  - 场景描述文本
  - 角色参考图
  - 场景上下文(地点、时间、氛围)
- **输出**: 具有一致角色外观的场景图片

**工作流示例**:
```go
// 伪代码
func GenerateSceneImage(scene Scene, characters []Character) (string, error) {
    // 1. 收集角色参考图
    var refImages []string
    for _, char := range scene.Characters {
        refImages = append(refImages, char.ReferenceImages[0])
    }

    // 2. 生成带角色一致性的场景图
    sceneImage, err := geminiClient.ImageToImage(
        prompt: scene.Description,
        referenceImages: refImages,
        style: scene.Style,
    )

    return sceneImage, err
}
```

## API 集成

### Gemini 2.5 Flash Image API

#### 文本到图像(参考图生成)
```http
POST /v1/images:generate
Content-Type: application/json

{
  "prompt": "一位年轻女子,长发黑发,明亮的眼睛,穿着中国传统服饰",
  "model": "gemini-2.5-flash-image",
  "style": "anime",
  "quality": "high"
}
```

#### 图像到图像(场景生成)
```http
POST /v1/images:transform
Content-Type: application/json

{
  "prompt": "角色站在日落时分的竹林中",
  "reference_images": [
    "https://storage.example.com/characters/char_001_ref.jpg"
  ],
  "model": "gemini-2.5-flash-image",
  "style": "anime",
  "consistency_strength": 0.8
}
```

## 数据流

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 小说上传与角色提取                                        │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. 生成参考图 (Gemini 文本到图像)                            │
│    - 输入: 角色描述                                          │
│    - 输出: 参考图 URL                                        │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. 将参考图存储到角色实体                                    │
│    - 保存到: Character.ReferenceImages                      │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. 场景划分与分析                                            │
│    - 识别每个场景中的角色                                    │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. 生成场景图片 (Gemini 图像到图像)                          │
│    - 输入: 场景描述 + 角色参考图                             │
│    - 输出: 具有一致性的场景图片                              │
└─────────────────────────────────────────────────────────────┘
```

## DDD 分层实现

### Domain 层
```go
// internal/domain/character/entity.go
type Character struct {
    ID              string
    ReferenceImages []string
    // ...
}

// internal/domain/character/repository.go
type Repository interface {
    FindByID(ctx context.Context, id string) (*Character, error)
    Save(ctx context.Context, char *Character) error
}
```

### Application 层
```go
// internal/application/service/character_service.go
type CharacterService struct {
    repo      domain.Repository
    aiService AIService
}

func (s *CharacterService) GenerateReferenceImage(
    ctx context.Context,
    characterID string,
) error {
    char, _ := s.repo.FindByID(ctx, characterID)
    refImage, _ := s.aiService.TextToImage(char.Description)
    char.ReferenceImages = append(char.ReferenceImages, refImage)
    return s.repo.Save(ctx, char)
}
```

### Infrastructure 层
```go
// internal/infrastructure/ai/gemini/client.go
type GeminiClient struct {
    apiKey string
    client *http.Client
}

func (c *GeminiClient) TextToImage(prompt string) (string, error) {
    // 调用 Gemini API
}

func (c *GeminiClient) ImageToImage(
    prompt string,
    refImages []string,
) (string, error) {
    // 调用 Gemini API
}
```

## 数据库设计

```sql
CREATE TABLE characters (
    id VARCHAR(36) PRIMARY KEY,
    novel_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    reference_images JSON,  -- 图片 URL 数组
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_novel_id (novel_id),
    FOREIGN KEY (novel_id) REFERENCES novels(id) ON DELETE CASCADE
);
```

**数据示例**:
```json
{
  "id": "char_001",
  "novel_id": "novel_001",
  "name": "李雪",
  "description": "一位年轻女子,长发黑发...",
  "reference_images": [
    "https://storage.example.com/characters/char_001_ref_01.jpg",
    "https://storage.example.com/characters/char_001_ref_02.jpg"
  ]
}
```

## 边缘场景与注意事项

### 一个场景中的多个角色
- **挑战**: 场景包含 2 个或更多角色
- **解决方案**: 将所有角色参考图传递给图像到图像 API
- **实现**:
  ```go
  refImages := []string{
      char1.ReferenceImages[0],
      char2.ReferenceImages[0],
  }
  ```

### 角色外观变化
- **挑战**: 角色更换服装或发型
- **解决方案**: 为不同状态生成新的参考图
- **实现**: 存储带元数据的多个参考图
  ```go
  type CharacterReference struct {
      URL         string
      State       string  // "default", "battle", "formal"
      Description string
  }
  ```

### 参考图质量问题
- **挑战**: 生成的参考图与描述不符
- **解决方案**:
  1. 允许使用调整后的提示词重新生成
  2. 支持手动上传参考图
  3. 实现质量评分机制

### 性能优化
- **缓存**: 在 CDN 中缓存参考图
- **懒加载**: 仅在需要时生成参考图
- **批处理**: 并行生成多个参考图

## 质量保证

### 验证检查清单
- [ ] 参考图与角色描述匹配
- [ ] 场景图片保持角色特征(面部、发型、服装)
- [ ] 多个场景显示一致的外观
- [ ] 不同角度/姿势保持可识别性

### 测试策略
```go
// 测试参考图生成
func TestGenerateReferenceImage(t *testing.T) {
    service := NewCharacterService(mockRepo, mockAI)
    char := &Character{
        Description: "黑发年轻女子",
    }

    err := service.GenerateReferenceImage(ctx, char.ID)
    assert.NoError(t, err)
    assert.NotEmpty(t, char.ReferenceImages)
}

// 测试场景生成一致性
func TestSceneGenerationConsistency(t *testing.T) {
    // 使用同一角色生成 3 个场景
    // 使用图像比较验证视觉相似性
}
```

## 未来增强

1. **多参考支持**: 使用多个参考图以获得更好的一致性
2. **风格迁移**: 对所有角色应用小说特定的艺术风格
3. **角色演变**: 跟踪故事时间线上的角色外观变化
4. **AI 驱动验证**: 使用计算机视觉自动验证一致性
5. **用户反馈循环**: 允许用户评分并改进参考图质量

## 相关文档

- [ARCHITECTURE.md](ARCHITECTURE.md) - 整体系统架构
- [API.md](API.md) - API 端点规范
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南

## 参考资料

- Gemini 2.5 Flash Image API 文档
- 图像到图像生成最佳实践
- 角色设计一致性指南
