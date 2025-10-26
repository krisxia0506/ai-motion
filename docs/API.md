# API 文档

## 概述

AI-Motion 提供 RESTful API 接口,用于小说解析、角色管理、场景管理和内容生成功能。

**版本**: v0.1.0-alpha  
**Base URL**:
- 开发环境: `http://localhost:8080`
- 生产环境: 根据部署配置而定

---

## 统一响应格式

所有 API 响应遵循统一的三段式结构:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**字段说明:**
- `code`: 业务状态码,0 表示成功,非 0 表示失败
- `message`: 操作结果描述信息
- `data`: 响应数据载体,成功时包含业务数据,失败时为 `null`

**命名规范:**
- 所有字段使用驼峰命名 (camelCase): `novelId`, `characterId`, `createdAt`
- 时间格式: ISO 8601 `2024-01-01T12:00:00Z` (UTC)

---

## 业务状态码

| Code  | 说明                     | 场景                                    |
|-------|------------------------|-----------------------------------------|
| 0     | 成功                    | 操作成功                                |
| 10001 | 参数错误                | 必填参数缺失、格式错误、类型不匹配        |
| 10002 | 资源不存在              | Novel/Character/Scene 不存在            |
| 30002 | 文件解析失败            | 小说解析失败                            |
| 40001 | AI 服务调用失败         | Gemini/Sora API 错误                    |
| 40003 | 生成任务失败            | 图像/视频生成失败                       |
| 50001 | 数据库错误              | 数据库操作失败                          |
| 50002 | 系统内部错误            | 未知错误                                |

---

## 接口分类

1. [系统健康检查](#1-系统健康检查)
2. [小说管理](#2-小说管理)
3. [角色管理](#3-角色管理)
4. [场景管理](#4-场景管理)
5. [提示词生成](#5-提示词生成)
6. [内容生成](#6-内容生成)
7. [漫画生成](#7-漫画生成)

---

## 1. 系统健康检查

### 1.1 GET /health

检查服务健康状态

**请求示例**
```bash
curl http://localhost:8080/health
```

**响应示例**
```json
{
  "status": "ok",
  "service": "ai-motion"
}
```

**说明**: 此接口不使用统一响应格式,直接返回健康状态

---

## 2. 小说管理

### 2.1 POST /api/v1/novel/upload

上传小说内容

**请求参数**
```json
{
  "title": "小说标题",
  "author": "作者名",
  "content": "小说内容..."
}
```

**参数说明**
- `title` (required) - 小说标题
- `author` (required) - 作者名称
- `content` (required) - 小说内容,100-5000 字

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novel/upload \
  -H "Content-Type: application/json" \
  -d '{
    "title": "修仙传",
    "author": "作者名",
    "content": "从前有座山..."
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "修仙传",
    "author": "作者名",
    "status": "uploaded",
    "word_count": 1500,
    "chapter_count": 0,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**业务逻辑**
1. 验证字数限制 (100-5000 字)
2. 创建 Novel 实体
3. 自动解析章节
4. 保存到数据库
5. 返回小说信息

**错误示例**
```json
{
  "code": 10001,
  "message": "小说内容不能少于100字",
  "data": null
}
```

---

### 2.2 GET /api/v1/novel/:id

获取小说详细信息

**路径参数**
- `id` (required) - 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/novel/550e8400-e29b-41d4-a716-446655440000
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "修仙传",
    "author": "作者名",
    "status": "uploaded",
    "word_count": 1500,
    "chapter_count": 3,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**错误示例**
```json
{
  "code": 10002,
  "message": "小说不存在",
  "data": null
}
```

---

### 2.3 GET /api/v1/novel

获取小说列表

**查询参数**
- `offset` (optional, default: 0) - 偏移量
- `limit` (optional, default: 20) - 每页数量

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novel?offset=0&limit=20"
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "修仙传",
        "author": "作者名",
        "status": "uploaded",
        "word_count": 1500,
        "chapter_count": 3,
        "created_at": "2024-01-01T12:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 5,
      "totalPages": 1,
      "hasNext": false,
      "hasPrev": false
    }
  }
}
```

---

### 2.4 DELETE /api/v1/novel/:id

删除小说及关联数据

**路径参数**
- `id` (required) - 小说 ID

**请求示例**
```bash
curl -X DELETE http://localhost:8080/api/v1/novel/550e8400-e29b-41d4-a716-446655440000
```

**响应示例**
```json
{
  "code": 0,
  "message": "小说删除成功",
  "data": null
}
```

**业务逻辑**: 级联删除所有章节、角色、场景、媒体文件

---

### 2.5 GET /api/v1/novel/:id/chapters

获取小说的所有章节

**路径参数**
- `id` (required) - 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/novel/550e8400-e29b-41d4-a716-446655440000/chapters
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "chapters": [
      {
        "id": "chapter_001",
        "chapter_number": 1,
        "title": "第一章 初入江湖",
        "word_count": 500,
        "created_at": "2024-01-01T12:00:00Z"
      },
      {
        "id": "chapter_002",
        "chapter_number": 2,
        "title": "第二章 奇遇",
        "word_count": 500,
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

## 3. 角色管理

### 3.1 POST /api/v1/characters/novel/:novel_id/extract

从小说中提取角色信息

**路径参数**
- `novel_id` (required) - 小说 ID

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/characters/novel/550e8400-e29b-41d4-a716-446655440000/extract
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "characters": [
      {
        "id": "char_001",
        "novel_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "李雪",
        "description": "女主角,18 岁,黑色长发,明亮的眼睛",
        "appearance_count": 45,
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

**业务逻辑**
1. 读取小说内容
2. 使用正则表达式识别中文角色名
3. 提取角色对话和外貌描述
4. 创建 Character 实体
5. 保存到数据库

---

### 3.2 GET /api/v1/characters/:id

获取单个角色详情

**路径参数**
- `id` (required) - 角色 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/characters/char_001
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "char_001",
    "novel_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "李雪",
    "description": "女主角,18 岁,黑色长发,明亮的眼睛",
    "appearance_count": 45,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 3.3 GET /api/v1/characters/novel/:novel_id

获取小说的所有角色

**路径参数**
- `novel_id` (required) - 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/characters/novel/550e8400-e29b-41d4-a716-446655440000
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "characters": [
      {
        "id": "char_001",
        "name": "李雪",
        "description": "女主角,18 岁...",
        "appearance_count": 45
      },
      {
        "id": "char_002",
        "name": "张伟",
        "description": "男主角,20 岁...",
        "appearance_count": 42
      }
    ]
  }
}
```

---

### 3.4 PUT /api/v1/characters/:id

更新角色信息

**路径参数**
- `id` (required) - 角色 ID

**请求体**
```json
{
  "name": "李雪",
  "description": "更新后的角色描述"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "角色信息更新成功",
  "data": {
    "id": "char_001",
    "name": "李雪",
    "description": "更新后的角色描述",
    "updated_at": "2024-01-01T13:00:00Z"
  }
}
```

---

### 3.5 DELETE /api/v1/characters/:id

删除角色

**路径参数**
- `id` (required) - 角色 ID

**响应示例**
```json
{
  "code": 0,
  "message": "角色删除成功",
  "data": null
}
```

---

### 3.6 POST /api/v1/characters/merge

合并重复的角色

**请求体**
```json
{
  "source_id": "char_002",
  "target_id": "char_001"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "角色合并成功",
  "data": {
    "id": "char_001",
    "name": "李雪",
    "appearance_count": 87
  }
}
```

**业务逻辑**: 将 `source_id` 的角色信息合并到 `target_id`,并删除 `source_id`

---

## 4. 场景管理

### 4.1 POST /api/v1/scenes/chapter/:chapter_id/divide

将章节划分为场景

**路径参数**
- `chapter_id` (required) - 章节 ID

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/scenes/chapter/chapter_001/divide
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "scenes": [
      {
        "id": "scene_001",
        "chapter_id": "chapter_001",
        "sequence_num": 1,
        "description": "清晨的竹林,阳光透过竹叶洒下",
        "location": "竹林",
        "time_of_day": "清晨",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

**业务逻辑**
1. 读取章节内容
2. 根据地点、时间等标记自动划分场景
3. 提取场景描述和对话
4. 创建 Scene 实体
5. 保存到数据库

---

### 4.2 GET /api/v1/scenes/:id

获取场景详情

**路径参数**
- `id` (required) - 场景 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/scenes/scene_001
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "scene_001",
    "chapter_id": "chapter_001",
    "sequence_num": 1,
    "description": "清晨的竹林,阳光透过竹叶洒下",
    "location": "竹林",
    "time_of_day": "清晨",
    "dialogues": [
      {
        "character_id": "char_001",
        "character_name": "李雪",
        "text": "今天天气真好"
      }
    ],
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

### 4.3 GET /api/v1/scenes/chapter/:chapter_id

获取章节的所有场景

**路径参数**
- `chapter_id` (required) - 章节 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/scenes/chapter/chapter_001
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "scenes": [
      {
        "id": "scene_001",
        "sequence_num": 1,
        "description": "清晨的竹林",
        "location": "竹林"
      },
      {
        "id": "scene_002",
        "sequence_num": 2,
        "description": "山间小路",
        "location": "山路"
      }
    ]
  }
}
```

---

### 4.4 GET /api/v1/scenes/novel/:novel_id

获取小说的所有场景

**路径参数**
- `novel_id` (required) - 小说 ID

**查询参数**
- `offset` (optional, default: 0) - 偏移量
- `limit` (optional, default: 20) - 每页数量

**请求示例**
```bash
curl "http://localhost:8080/api/v1/scenes/novel/550e8400-e29b-41d4-a716-446655440000?offset=0&limit=20"
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "scene_001",
        "chapter_id": "chapter_001",
        "sequence_num": 1,
        "description": "清晨的竹林",
        "location": "竹林"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 50,
      "totalPages": 3,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

---

### 4.5 DELETE /api/v1/scenes/:id

删除场景

**路径参数**
- `id` (required) - 场景 ID

**响应示例**
```json
{
  "code": 0,
  "message": "场景删除成功",
  "data": null
}
```

---

## 5. 提示词生成

### 5.1 POST /api/v1/prompts/generate

为场景生成 AI 提示词

**请求体**
```json
{
  "scene_id": "scene_001",
  "type": "image"
}
```

**参数说明**
- `scene_id` (required) - 场景 ID
- `type` (required) - 提示词类型: `image` | `video`

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/prompts/generate \
  -H "Content-Type: application/json" \
  -d '{
    "scene_id": "scene_001",
    "type": "image"
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "scene_id": "scene_001",
    "type": "image",
    "prompt": "清晨的竹林,阳光透过竹叶,一位黑发年轻女子站在竹林中,anime style,高质量,细节丰富"
  }
}
```

**业务逻辑**
1. 获取场景描述、位置、时间
2. 获取场景中的角色描述
3. 组合成结构化提示词
4. 返回优化后的提示词

---

### 5.2 POST /api/v1/prompts/generate/batch

批量生成场景提示词

**请求体**
```json
{
  "scene_ids": ["scene_001", "scene_002", "scene_003"],
  "type": "image"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "prompts": [
      {
        "scene_id": "scene_001",
        "prompt": "清晨的竹林,阳光透过竹叶..."
      },
      {
        "scene_id": "scene_002",
        "prompt": "山间小路,夕阳西下..."
      }
    ]
  }
}
```

---

## 6. 内容生成

### 6.1 POST /api/v1/generate/image

生成场景图片

**请求体**
```json
{
  "scene_id": "scene_001",
  "style": "anime",
  "use_character_reference": true
}
```

**参数说明**
- `scene_id` (required) - 场景 ID
- `style` (optional, default: "anime") - 图片风格
- `use_character_reference` (optional, default: true) - 是否使用角色参考图保证一致性

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/generate/image \
  -H "Content-Type: application/json" \
  -d '{
    "scene_id": "scene_001",
    "style": "anime",
    "use_character_reference": true
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "图片生成成功",
  "data": {
    "scene_id": "scene_001",
    "media_id": "media_001",
    "url": "https://storage.example.com/scene_001.jpg",
    "status": "completed"
  }
}
```

**业务逻辑**
1. 获取场景信息和角色参考图
2. 构建提示词
3. 调用 Gemini API 生成图片 (Image-to-Image 或 Text-to-Image)
4. 保存图片到存储
5. 创建 Media 实体
6. 返回结果

**角色一致性**: 如果 `use_character_reference=true`,会使用角色参考图进行 Image-to-Image 生成

---

### 6.2 POST /api/v1/generate/video

生成场景视频

**请求体**
```json
{
  "scene_id": "scene_001",
  "source_image_id": "media_001",
  "duration": 5
}
```

**参数说明**
- `scene_id` (required) - 场景 ID
- `source_image_id` (optional) - 源图片 ID,默认使用场景已生成的图片
- `duration` (optional, default: 5) - 视频时长(秒)

**响应示例**
```json
{
  "code": 0,
  "message": "视频生成成功",
  "data": {
    "scene_id": "scene_001",
    "media_id": "media_002",
    "url": "https://storage.example.com/scene_001.mp4",
    "duration": 5,
    "status": "completed"
  }
}
```

**业务逻辑**
1. 获取场景图片
2. 调用 Sora2 ImageToVideo API
3. 保存视频到存储
4. 创建 Media 实体
5. 返回结果

---

### 6.3 POST /api/v1/generate/batch

批量生成场景内容

**请求体**
```json
{
  "scene_ids": ["scene_001", "scene_002", "scene_003"],
  "content_type": "image",
  "options": {
    "style": "anime",
    "use_character_reference": true
  }
}
```

**参数说明**
- `scene_ids` (required) - 场景 ID 数组
- `content_type` (required) - 生成类型: `image` | `video`
- `options` (optional) - 生成选项

**响应示例**
```json
{
  "code": 0,
  "message": "批量生成成功",
  "data": {
    "results": [
      {
        "scene_id": "scene_001",
        "status": "completed",
        "media_id": "media_001"
      },
      {
        "scene_id": "scene_002",
        "status": "completed",
        "media_id": "media_002"
      },
      {
        "scene_id": "scene_003",
        "status": "failed",
        "error": "AI 服务调用失败"
      }
    ],
    "summary": {
      "total": 3,
      "completed": 2,
      "failed": 1
    }
  }
}
```

**业务逻辑**: 使用 Go 协程并发处理多个场景生成任务

---

### 6.4 GET /api/v1/generate/status/:scene_id

查询场景生成状态

**路径参数**
- `scene_id` (required) - 场景 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/generate/status/scene_001
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "scene_id": "scene_001",
    "has_image": true,
    "has_video": false,
    "media": [
      {
        "id": "media_001",
        "type": "image",
        "url": "https://storage.example.com/scene_001.jpg",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

## 7. 漫画生成

### 7.1 POST /api/v1/manga/generate

一键生成漫画 (端到端流程)

**请求体**
```json
{
  "title": "小红帽",
  "author": "格林兄弟",
  "content": "从前有个可爱的小姑娘..."
}
```

**参数说明**
- `title` (required) - 小说标题
- `author` (required) - 作者名称
- `content` (required) - 小说内容,100-5000 字

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/manga/generate \
  -H "Content-Type: application/json" \
  -d '{
    "title": "小红帽",
    "author": "格林兄弟",
    "content": "从前有个可爱的小姑娘,谁见了都喜欢..."
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "novel_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "小红帽",
    "character_count": 3,
    "scene_count": 8,
    "status": "completed",
    "message": "Successfully generated manga with 3 characters and 8 scenes"
  }
}
```

**业务逻辑**
1. 上传并解析小说
2. 自动提取角色
3. 为每个角色生成参考图 (Gemini Text-to-Image)
4. 将章节划分为场景
5. 匹配场景与角色
6. 为每个场景生成图片 (使用角色参考图保证一致性)
7. 返回生成结果

**核心特性**: 
- **字数限制**: 100-5000 字,确保大模型能一次性处理
- **角色一致性**: 通过参考图 + Image-to-Image 保证角色外观统一
- **端到端自动化**: 一个 API 调用完成从文本到漫画的全流程

**工作流程**:
```
上传小说 → 解析小说 → 提取角色 → 生成角色参考图 → 
划分场景 → 匹配场景与角色 → 生成场景图片 → 完成
```

**错误示例**
```json
{
  "code": 10001,
  "message": "小说内容不能超过5000字",
  "data": null
}
```

---

## HTTP 状态码

- `200 OK` - 请求成功 (包括业务逻辑错误,通过 code 区分)
- `400 Bad Request` - 请求格式错误
- `404 Not Found` - 路由不存在
- `500 Internal Server Error` - 服务器内部错误
- `503 Service Unavailable` - AI 服务不可用

**注意**: 业务逻辑错误统一返回 HTTP 200,通过 `code` 字段区分具体错误

---

## 响应结构设计

所有接口使用统一的响应结构,定义在 `backend/internal/interfaces/http/response/response.go`:

### 基础响应
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

### 分页响应
```go
type PaginationData struct {
    Items      interface{}     `json:"items"`
    Pagination *PaginationInfo `json:"pagination,omitempty"`
}

type PaginationInfo struct {
    Page       int  `json:"page"`
    PageSize   int  `json:"pageSize"`
    Total      int  `json:"total"`
    TotalPages int  `json:"totalPages"`
    HasNext    bool `json:"hasNext"`
    HasPrev    bool `json:"hasPrev"`
}
```

### 辅助函数

- `Success(c, data)` - 成功响应
- `SuccessWithMessage(c, message, data)` - 自定义消息的成功响应
- `SuccessList(c, items, page, pageSize, total)` - 分页列表响应
- `Error(c, code, message)` - 错误响应
- `InvalidParams(c, message)` - 参数错误 (10001)
- `ResourceNotFound(c, message)` - 资源不存在 (10002)
- `FileParseError(c, message)` - 文件解析失败 (30002)
- `AIServiceError(c, message)` - AI 服务错误 (40001)
- `GenerationError(c, message)` - 生成任务失败 (40003)
- `DatabaseError(c, message)` - 数据库错误 (50001)
- `InternalError(c, message)` - 系统内部错误 (50002)

---

## 接口实现状态

| 功能模块 | 状态 | 说明 |
|---------|------|------|
| 系统健康检查 | ✅ 已实现 | 基础健康检查 |
| 小说管理 | ✅ 已实现 | 上传、查询、删除、章节列表 |
| 角色管理 | ✅ 已实现 | 提取、查询、更新、删除、合并 |
| 场景管理 | ✅ 已实现 | 划分、查询、删除 |
| 提示词生成 | ✅ 已实现 | 单个和批量生成 |
| 内容生成 | ✅ 已实现 | 图片、视频、批量生成、状态查询 |
| 漫画生成 | ✅ 已实现 | 端到端自动化生成流程 |
| 用户认证 | ⏳ 待实现 | JWT 认证、注册、登录 |
| 项目管理 | ⏳ 待实现 | 项目创建、管理 |
| 导出功能 | ⏳ 待实现 | 视频导出、素材打包 |

---

## 相关文档

- [API_DESIGN_GUIDELINES.md](API_DESIGN_GUIDELINES.md) - API 设计规范
- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构设计
- [CHARACTER_CONSISTENCY.md](CHARACTER_CONSISTENCY.md) - 角色一致性设计
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南
- [backend/CLAUDE.md](../backend/CLAUDE.md) - 后端开发指南
- [QUICKSTART.md](../QUICKSTART.md) - 快速开始

---

*API 文档版本: v0.1.0-alpha*  
*最后更新: 2024-01-26*  
*符合 API 设计规范 v1.0*  
*基于实际代码实现编写*
