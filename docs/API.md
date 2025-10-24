# API 文档

## 概述

AI-Motion 提供 RESTful API 接口，用于小说解析、角色管理、内容生成和导出功能。

## Base URL

- 开发环境: `http://localhost:8080`
- 生产环境: 根据部署配置而定

## 认证

目前处于开发阶段，暂未实现认证机制。未来版本将支持 API Key 或 JWT 认证。

## 接口列表

### 健康检查

#### GET /health

检查服务健康状态

**请求示例**
```bash
curl http://localhost:8080/health
```

**响应示例**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

### 小说管理

#### POST /api/v1/novel/upload

上传小说文件

**请求参数**
- `file` (multipart/form-data): 小说文件 (支持 TXT 格式)

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novel/upload \
  -F "file=@novel.txt"
```

**响应示例**
```json
{
  "novel_id": "uuid-string",
  "filename": "novel.txt",
  "upload_time": "2024-01-01T00:00:00Z"
}
```

#### POST /api/v1/novel/:id/parse

解析小说内容，提取章节、角色和场景信息

**路径参数**
- `id`: 小说 ID

**请求示例**
```bash
curl -X POST http://localhost:8080/api/v1/novel/abc123/parse
```

**响应示例**
```json
{
  "novel_id": "abc123",
  "chapters": 50,
  "characters": 12,
  "scenes": 200,
  "status": "completed"
}
```

---

### 角色管理

#### GET /api/v1/characters/:novel_id

获取指定小说的角色列表

**路径参数**
- `novel_id`: 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/characters/abc123
```

**响应示例**
```json
{
  "characters": [
    {
      "id": "char001",
      "name": "张三",
      "description": "主角，年轻的剑客",
      "appearances": 45,
      "reference_image": "https://storage/char001.jpg"
    },
    {
      "id": "char002",
      "name": "李四",
      "description": "反派角色",
      "appearances": 30,
      "reference_image": "https://storage/char002.jpg"
    }
  ]
}
```

---

### 内容生成

#### POST /api/v1/generate/scene

生成场景图片

**请求参数**
```json
{
  "novel_id": "abc123",
  "scene_id": "scene001",
  "description": "清晨的竹林，阳光透过竹叶",
  "characters": ["char001"],
  "style": "anime"
}
```

**响应示例**
```json
{
  "scene_id": "scene001",
  "image_url": "https://storage/scene001.jpg",
  "generation_time": 5.2,
  "status": "success"
}
```

#### POST /api/v1/generate/voice

生成角色配音

**请求参数**
```json
{
  "character_id": "char001",
  "text": "你好，我是张三",
  "emotion": "neutral",
  "voice_profile": "male_young"
}
```

**响应示例**
```json
{
  "audio_url": "https://storage/voice001.mp3",
  "duration": 2.5,
  "status": "success"
}
```

---

### 导出功能

#### POST /api/v1/anime/:id/export

导出生成的动漫内容

**路径参数**
- `id`: 动漫项目 ID

**请求参数**
```json
{
  "format": "mp4",
  "quality": "high",
  "chapters": [1, 2, 3]
}
```

**响应示例**
```json
{
  "export_id": "export001",
  "status": "processing",
  "estimated_time": 120
}
```

---

## 错误处理

所有错误响应遵循统一格式：

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述",
    "details": {}
  }
}
```

### 常见错误码

- `400` - 请求参数错误
- `404` - 资源不存在
- `500` - 服务器内部错误
- `503` - 服务暂时不可用

## 速率限制

目前暂未实施速率限制。未来版本将根据用户类型设置不同的请求限制。

## WebSocket API

未来版本将支持 WebSocket 接口，用于实时推送生成进度。

---

*API 文档版本: v0.1.0-alpha*
