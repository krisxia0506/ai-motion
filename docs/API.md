# API 文档

## 概述

AI-Motion 提供 RESTful API 接口，用于小说解析、角色管理、场景管理、内容生成和导出功能。

**版本**: v0.1.0-alpha
**Base URL**:
- 开发环境：`http://localhost:8080`
- 生产环境：根据部署配置而定

## 认证

API 使用 **JWT (JSON Web Token)** 进行认证。

### 认证流程

1. 用户通过 `/api/v1/auth/register` 注册账号
2. 用户通过 `/api/v1/auth/login` 登录,获取 JWT token
3. 在后续请求的 HTTP Header 中携带 token: `Authorization: Bearer <token>`
4. Token 过期时间默认为 7 天,可通过 `/api/v1/auth/refresh` 刷新

### Header 格式

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 公开接口 (无需认证)

- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### 认证接口 (需要 JWT)

所有其他接口均需要在 Header 中携带有效的 JWT token。

---

## 接口分类

1. [系统健康检查](#1-系统健康检查)
2. [认证管理](#2-认证管理)
3. [小说管理](#3-小说管理)
4. [角色管理](#4-角色管理)
5. [场景管理](#5-场景管理)
6. [内容生成](#6-内容生成)
7. [项目管理](#7-项目管理)
8. [导出功能](#8-导出功能)

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
  "timestamp": "2024-01-01T12:00:00Z",
  "services": {
    "database": "connected",
    "gemini_api": "available",
    "sora_api": "available"
  }
}
```

**业务逻辑**
- 检查数据库连接状态 (可选)
- 检查 AI 服务可用性
- 返回系统当前时间戳

---

## 2. 认证管理

### 2.1 POST /api/v1/auth/register

用户注册

**请求参数**
```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**参数说明**
- `username` (required) - 用户名,3-20个字符,仅支持字母、数字、下划线
- `email` (required) - 邮箱地址
- `password` (required) - 密码,至少8个字符,需包含大小写字母和数字

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "email": "user@example.com",
    "password": "securePassword123"
  }'
```

**响应示例**
```json
{
  "user_id": "user_abc123",
  "username": "user123",
  "email": "user@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "message": "Registration successful"
}
```

**业务逻辑**
1. 验证用户名、邮箱、密码格式
2. 检查用户名和邮箱是否已存在
3. 使用 bcrypt 加密密码(cost=10)
4. 创建 User 实体
5. 保存到数据库
6. 返回用户基本信息

**错误码**
- `INVALID_USERNAME` - 用户名格式不正确
- `INVALID_EMAIL` - 邮箱格式不正确
- `WEAK_PASSWORD` - 密码强度不足
- `USERNAME_EXISTS` - 用户名已存在
- `EMAIL_EXISTS` - 邮箱已存在

---

### 2.2 POST /api/v1/auth/login

用户登录

**请求参数**
```json
{
  "username": "user123",
  "password": "securePassword123"
}
```

**参数说明**
- `username` (required) - 用户名或邮箱
- `password` (required) - 密码

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "password": "securePassword123"
  }'
```

**响应示例**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcl9hYmMxMjMiLCJ1c2VybmFtZSI6InVzZXIxMjMiLCJleHAiOjE3MDQxOTY4MDB9.xxx",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcl9hYmMxMjMiLCJ0eXBlIjoicmVmcmVzaCIsImV4cCI6MTcwNjc4ODgwMH0.yyy",
  "token_type": "Bearer",
  "expires_in": 604800,
  "user": {
    "user_id": "user_abc123",
    "username": "user123",
    "email": "user@example.com"
  }
}
```

**响应字段说明**
- `access_token` - 访问令牌,用于 API 请求
- `refresh_token` - 刷新令牌,用于获取新的访问令牌
- `token_type` - 令牌类型,固定为 "Bearer"
- `expires_in` - 访问令牌过期时间(秒),默认 7 天

**业务逻辑**
1. 查找用户(支持用户名或邮箱登录)
2. 验证密码(使用 bcrypt.CompareHashAndPassword)
3. 生成 JWT access_token:
   - Payload: `user_id`, `username`, `exp`(过期时间)
   - 签名算法: HS256
   - 过期时间: 7天
4. 生成 JWT refresh_token:
   - Payload: `user_id`, `type: "refresh"`, `exp`(过期时间)
   - 过期时间: 30天
5. 返回 token 和用户信息

**错误码**
- `INVALID_CREDENTIALS` - 用户名或密码错误
- `USER_NOT_FOUND` - 用户不存在
- `ACCOUNT_LOCKED` - 账号已锁定(可选,未来版本)

---

### 2.3 POST /api/v1/auth/refresh

刷新访问令牌

**请求参数**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**响应示例**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_token",
  "token_type": "Bearer",
  "expires_in": 604800
}
```

**业务逻辑**
1. 验证 refresh_token 是否有效
2. 解析 token 获取 user_id
3. 检查用户是否存在
4. 生成新的 access_token
5. 返回新 token

**错误码**
- `INVALID_TOKEN` - Token 无效
- `TOKEN_EXPIRED` - Token 已过期
- `USER_NOT_FOUND` - 用户不存在

---

### 2.4 POST /api/v1/auth/logout

**认证**: 需要 JWT

用户登出

**认证**: 需要 JWT

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "message": "Logout successful"
}
```

**业务逻辑**
1. 解析 JWT token 获取 user_id
2. (可选)将 token 加入黑名单(需要 Redis)
3. 返回成功消息

**注意**: 由于 JWT 是无状态的,客户端需要删除本地存储的 token

---

### 2.5 GET /api/v1/auth/me

**认证**: 需要 JWT

获取当前用户信息

**认证**: 需要 JWT

**请求示例**
```bash
curl -X GET \
  http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "user_id": "user_abc123",
  "username": "user123",
  "email": "user@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "usage_stats": {
    "novels_count": 5,
    "characters_count": 30,
    "scenes_generated": 150
  }
}
```

**业务逻辑**
1. 从 JWT token 中解析 user_id
2. 查询用户信息
3. 统计用户使用情况
4. 返回完整用户信息

---

### 2.6 PUT /api/v1/auth/password

**认证**: 需要 JWT

修改密码

**认证**: 需要 JWT

**请求参数**
```json
{
  "old_password": "oldPassword123",
  "new_password": "newPassword456"
}
```

**请求示例**
```bash
curl -X PUT \
  http://localhost:8080/api/v1/auth/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldPassword123",
    "new_password": "newPassword456"
  }'
```

**响应示例**
```json
{
  "message": "Password updated successfully"
}
```

**业务逻辑**
1. 从 JWT 获取 user_id
2. 验证旧密码是否正确
3. 验证新密码强度
4. 使用 bcrypt 加密新密码
5. 更新数据库
6. (可选)使所有旧 token 失效

**错误码**
- `INVALID_OLD_PASSWORD` - 旧密码错误
- `WEAK_PASSWORD` - 新密码强度不足
- `SAME_PASSWORD` - 新密码与旧密码相同

---

### JWT Token 结构

**Access Token Payload**
```json
{
  "user_id": "user_abc123",
  "username": "user123",
  "iat": 1704110400,
  "exp": 1704715200
}
```

**Refresh Token Payload**
```json
{
  "user_id": "user_abc123",
  "type": "refresh",
  "iat": 1704110400,
  "exp": 1706788800
}
```

**签名密钥**: 使用环境变量 `JWT_SECRET_KEY`,至少 32 字符

---

## 3. 小说管理

### 4.1 POST /api/v1/novels/upload

**认证**: 需要 JWT

上传小说文件

**认证**: 需要 JWT

**请求参数**
- `file` (multipart/form-data, required) - 小说文件，支持 TXT 格式
- `title` (form, optional) - 小说标题
- `author` (form, optional) - 作者名称

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novels/upload \
  -F "file=@novel.txt" \
  -F "title=修仙传" \
  -F "author=作者名"
```

**响应示例**
```json
{
  "novel_id": "novel_abc123",
  "title": "修仙传",
  "filename": "novel.txt",
  "file_size": 1024000,
  "status": "uploaded",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**业务逻辑**
1. 验证文件格式 (仅允许 .txt)
2. 验证文件大小限制 (如 50MB)
3. 生成唯一 Novel ID
4. 保存文件到临时存储
5. 创建 Novel 实体，状态为 "uploaded"
6. 保存到数据库
7. 返回小说基本信息

---

### 4.2 POST /api/v1/novels/:novel_id/parse

**认证**: 需要 JWT

解析小说内容，提取章节、角色、场景

**路径参数**
- `novel_id` (required) - 小说 ID

**请求体**
```json
{
  "options": {
    "auto_generate_references": true,
    "scene_division_mode": "auto"
  }
}
```

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novels/novel_abc123/parse \
  -H "Content-Type: application/json" \
  -d '{
    "options": {
      "auto_generate_references": true,
      "scene_division_mode": "auto"
    }
  }'
```

**响应示例**
```json
{
  "novel_id": "novel_abc123",
  "parse_result": {
    "chapters_count": 50,
    "characters_count": 12,
    "scenes_count": 200,
    "total_words": 500000
  },
  "status": "parsing",
  "estimated_time": 300
}
```

**业务逻辑**
1. 读取小说原始文本
2. 调用 Gemini API 进行自然语言处理
3. 提取章节结构 → 创建 Chapter 实体
4. 提取角色信息 → 创建 Character 实体
5. 划分场景 → 创建 Scene 实体
6. 分析对话和场景描述
7. 更新 Novel 状态为 "parsed"
8. (可选) 异步触发角色参考图生成
9. 返回解析统计信息

**注意**: 这是一个耗时操作，建议异步处理

---

### 4.3 GET /api/v1/novels/:novel_id

**认证**: 需要 JWT

获取小说详细信息

**路径参数**
- `novel_id` (required) - 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/novels/novel_abc123
```

**响应示例**
```json
{
  "id": "novel_abc123",
  "title": "修仙传",
  "author": "作者名",
  "status": "parsed",
  "chapters_count": 50,
  "characters_count": 12,
  "scenes_count": 200,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:30:00Z"
}
```

**状态值说明**
- `uploaded` - 已上传
- `parsing` - 解析中
- `parsed` - 解析完成
- `failed` - 解析失败

---

### 4.4 GET /api/v1/novels

**认证**: 需要 JWT

获取小说列表

**查询参数**
- `page` (optional, default: 1) - 页码
- `page_size` (optional, default: 20) - 每页数量
- `status` (optional) - 过滤状态

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels?page=1&page_size=20&status=parsed"
```

**响应示例**
```json
{
  "novels": [
    {
      "id": "novel_abc123",
      "title": "修仙传",
      "author": "作者名",
      "status": "parsed",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 5
  }
}
```

---

### 4.5 DELETE /api/v1/novels/:novel_id

**认证**: 需要 JWT

删除小说及关联数据

**路径参数**
- `novel_id` (required) - 小说 ID

**请求示例**
```bash
curl -X DELETE http://localhost:8080/api/v1/novels/novel_abc123
```

**响应示例**
```json
{
  "message": "Novel deleted successfully",
  "deleted_items": {
    "chapters": 50,
    "characters": 12,
    "scenes": 200,
    "media": 150
  }
}
```

**业务逻辑**
- 级联删除所有章节、角色、场景、媒体文件
- 删除存储的文件 (参考图、场景图、视频)
- 返回删除统计

---

## 4. 角色管理

### 4.1 GET /api/v1/novels/:novel_id/characters

**认证**: 需要 JWT

获取小说的所有角色

**路径参数**
- `novel_id` (required) - 小说 ID

**查询参数**
- `include_references` (optional, default: true) - 是否包含参考图

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels/novel_abc123/characters?include_references=true"
```

**响应示例**
```json
{
  "characters": [
    {
      "id": "char_001",
      "name": "李雪",
      "description": "女主角，18 岁，黑色长发，明亮的眼睛，身穿白色长裙",
      "role": "protagonist",
      "appearances_count": 45,
      "reference_images": [
        {
          "id": "ref_img_001",
          "url": "https://storage.example.com/ref_001.jpg",
          "state": "default",
          "created_at": "2024-01-01T12:00:00Z"
        }
      ],
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

**角色类型说明**
- `protagonist` - 主角
- `antagonist` - 反派
- `supporting` - 配角

---

### 4.2 GET /api/v1/characters/:character_id

**认证**: 需要 JWT

获取单个角色详情

**路径参数**
- `character_id` (required) - 角色 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/characters/char_001
```

**响应示例**
```json
{
  "id": "char_001",
  "novel_id": "novel_abc123",
  "name": "李雪",
  "description": "女主角，18 岁...",
  "appearance": {
    "age": "18",
    "gender": "female",
    "hair": "黑色长发",
    "eyes": "明亮的黑色眼睛",
    "clothing": "白色长裙"
  },
  "personality": {
    "traits": ["勇敢", "善良", "坚韧"],
    "description": "性格开朗，勇敢面对困难"
  },
  "reference_images": [
    {
      "id": "ref_img_001",
      "url": "https://storage.example.com/ref_001.jpg",
      "state": "default",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "scenes": ["scene_001", "scene_005"],
  "created_at": "2024-01-01T12:00:00Z"
}
```

---

### 4.3 POST /api/v1/characters/:character_id/references

**认证**: 需要 JWT

为角色生成参考图 (核心功能：角色一致性)

**路径参数**
- `character_id` (required) - 角色 ID

**请求体**
```json
{
  "state": "default",
  "custom_prompt": "",
  "style": "anime"
}
```

**状态类型说明**
- `default` - 默认外观
- `battle` - 战斗状态
- `formal` - 正式场合
- `custom` - 自定义状态

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/characters/char_001/references \
  -H "Content-Type: application/json" \
  -d '{
    "state": "default",
    "style": "anime"
  }'
```

**响应示例**
```json
{
  "reference_id": "ref_img_002",
  "character_id": "char_001",
  "status": "generating",
  "image_url": null,
  "estimated_time": 10
}
```

**业务逻辑**
1. 获取 Character 实体
2. 构建提示词 (使用角色 description + state + style)
3. 调用 Gemini TextToImage API
4. 保存图片到存储服务
5. 更新 Character.ReferenceImages
6. 返回生成结果

**角色一致性说明**: 此接口生成的参考图将用于后续所有场景生成，确保角色视觉一致性。

---

### 4.4 PUT /api/v1/characters/:character_id

**认证**: 需要 JWT

更新角色信息

**路径参数**
- `character_id` (required) - 角色 ID

**请求体**
```json
{
  "description": "更新后的描述",
  "appearance": {
    "age": "19",
    "hair": "银色长发"
  },
  "personality": {
    "traits": ["勇敢", "善良", "坚韧", "智慧"]
  }
}
```

**响应示例**
```json
{
  "id": "char_001",
  "name": "李雪",
  "description": "更新后的描述",
  "updated_at": "2024-01-01T13:00:00Z"
}
```

**注意**: 如果描述有重大变化，建议重新生成参考图

---

### 4.5 DELETE /api/v1/characters/:character_id/references/:reference_id

**认证**: 需要 JWT

删除角色参考图

**路径参数**
- `character_id` (required) - 角色 ID
- `reference_id` (required) - 参考图 ID

**响应示例**
```json
{
  "message": "Reference image deleted"
}
```

---

## 5. 场景管理

### 5.1 GET /api/v1/novels/:novel_id/scenes

**认证**: 需要 JWT

获取小说的所有场景

**路径参数**
- `novel_id` (required) - 小说 ID

**查询参数**
- `chapter_id` (optional) - 过滤章节
- `page` (optional, default: 1) - 页码
- `page_size` (optional, default: 20) - 每页数量

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels/novel_abc123/scenes?chapter_id=chapter_01&page=1"
```

**响应示例**
```json
{
  "scenes": [
    {
      "id": "scene_001",
      "novel_id": "novel_abc123",
      "chapter_id": "chapter_01",
      "sequence_num": 1,
      "description": "清晨的竹林，阳光透过竹叶洒下斑驳的光影",
      "location": "竹林",
      "time_of_day": "清晨",
      "characters": ["char_001", "char_003"],
      "dialogues": [
        {
          "character_id": "char_001",
          "text": "今天天气真好",
          "emotion": "happy"
        }
      ],
      "duration": 5.0,
      "media": {
        "image": "https://storage.example.com/scene_001.jpg",
        "video": null
      },
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 200
  }
}
```

---

### 5.2 GET /api/v1/scenes/:scene_id

**认证**: 需要 JWT

获取场景详情

**路径参数**
- `scene_id` (required) - 场景 ID

**响应示例**
```json
{
  "id": "scene_001",
  "novel_id": "novel_abc123",
  "chapter_id": "chapter_01",
  "sequence_num": 1,
  "description": "清晨的竹林，阳光透过竹叶洒下斑驳的光影",
  "location": "竹林",
  "time_of_day": "清晨",
  "characters": [
    {
      "id": "char_001",
      "name": "李雪",
      "description": "女主角..."
    }
  ],
  "dialogues": [
    {
      "character_id": "char_001",
      "text": "今天天气真好",
      "emotion": "happy"
    }
  ],
  "media": {
    "images": [
      {
        "id": "media_001",
        "url": "https://storage.example.com/scene_001.jpg",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "videos": []
  },
  "created_at": "2024-01-01T12:00:00Z"
}
```

---

### 5.3 PUT /api/v1/scenes/:scene_id

**认证**: 需要 JWT

更新场景信息

**路径参数**
- `scene_id` (required) - 场景 ID

**请求体**
```json
{
  "description": "更新后的场景描述",
  "location": "新地点",
  "time_of_day": "黄昏",
  "dialogues": [
    {
      "character_id": "char_001",
      "text": "太阳快要下山了",
      "emotion": "peaceful"
    }
  ]
}
```

**响应示例**
```json
{
  "id": "scene_001",
  "description": "更新后的场景描述",
  "updated_at": "2024-01-01T13:00:00Z"
}
```

---

### 5.4 POST /api/v1/scenes/:scene_id/prompt

**认证**: 需要 JWT

生成场景的 AI 提示词

**路径参数**
- `scene_id` (required) - 场景 ID

**响应示例**
```json
{
  "scene_id": "scene_001",
  "prompt": {
    "text_to_image": "清晨的竹林，阳光透过竹叶，一位黑发年轻女子站在竹林中，动漫风格，高质量，细节丰富",
    "image_to_video": "镜头缓慢推进，竹叶在微风中摇曳，女子转头看向远方"
  }
}
```

**业务逻辑**
1. 获取场景描述、位置、时间
2. 获取场景中的角色描述
3. 组合成结构化提示词
4. 返回文生图和图生视频两种提示词

---

## 6. 内容生成

### 6.1 POST /api/v1/generation/character-reference

**认证**: 需要 JWT

批量生成角色参考图

**请求体**
```json
{
  "character_ids": ["char_001", "char_002"],
  "options": {
    "style": "anime",
    "quality": "high"
  }
}
```

**响应示例**
```json
{
  "task_id": "gen_task_001",
  "status": "processing",
  "total": 2,
  "completed": 0,
  "estimated_time": 20
}
```

**业务逻辑**
- 创建异步任务
- 并发调用 Gemini TextToImage API
- 使用 Go 协程处理多个角色
- 保存结果到 Character 实体

---

### 6.2 POST /api/v1/generation/scene-image

**认证**: 需要 JWT

生成场景图片 (带角色一致性)

**请求体**
```json
{
  "scene_id": "scene_001",
  "options": {
    "style": "anime",
    "quality": "high",
    "consistency_strength": 0.8
  }
}
```

**参数说明**
- `consistency_strength` - 角色一致性强度，取值 0-1，越高越严格保持角色外观

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/generation/scene-image \
  -H "Content-Type: application/json" \
  -d '{
    "scene_id": "scene_001",
    "options": {
      "style": "anime",
      "quality": "high",
      "consistency_strength": 0.8
    }
  }'
```

**响应示例**
```json
{
  "media_id": "media_001",
  "scene_id": "scene_001",
  "status": "generating",
  "image_url": null,
  "estimated_time": 10
}
```

**业务逻辑**
1. 获取 Scene 实体
2. 获取场景中的所有角色
3. 收集角色参考图 (ReferenceImages[0])
4. 构建提示词 (场景描述 + 角色特征)
5. 调用 Gemini ImageToImage API:
   - prompt: 场景描述
   - reference_images: 角色参考图数组
   - consistency_strength: 一致性参数
6. 保存图片到存储
7. 创建 Media 实体，type="image"
8. 关联到 Scene
9. 返回结果

**角色一致性说明**: 此接口使用角色参考图进行图生图，确保场景中的角色与参考图保持视觉一致。

---

### 6.3 POST /api/v1/generation/scene-video

**认证**: 需要 JWT

生成场景视频 (图生视频)

**请求体**
```json
{
  "scene_id": "scene_001",
  "source_image": "media_001",
  "options": {
    "duration": 5,
    "fps": 30,
    "motion_type": "smooth"
  }
}
```

**参数说明**
- `source_image` - 可选，指定源图片 ID，默认使用场景已生成的图片
- `duration` - 视频时长 (秒)
- `fps` - 帧率
- `motion_type` - 动作类型：`smooth`(平滑) | `dynamic`(动态)

**响应示例**
```json
{
  "media_id": "media_002",
  "scene_id": "scene_001",
  "status": "generating",
  "video_url": null,
  "estimated_time": 60
}
```

**业务逻辑**
1. 获取 Scene 实体
2. 获取场景图片 (如果未指定，使用最新生成的图片)
3. 调用 Sora2 ImageToVideo API
4. 保存视频到存储
5. 创建 Media 实体，type="video"
6. 返回结果

---

### 6.4 POST /api/v1/generation/voice

**认证**: 需要 JWT

生成角色配音

**请求体**
```json
{
  "character_id": "char_001",
  "text": "今天天气真好",
  "options": {
    "emotion": "happy",
    "voice_profile": "female_young",
    "speed": 1.0
  }
}
```

**参数说明**
- `emotion` - 情感：`neutral` | `happy` | `sad` | `angry` | `surprised`
- `voice_profile` - 声音配置：`female_young` | `male_young` | `female_mature` | `male_mature`
- `speed` - 语速倍率，默认 1.0

**响应示例**
```json
{
  "audio_id": "audio_001",
  "character_id": "char_001",
  "audio_url": "https://storage.example.com/audio_001.mp3",
  "duration": 2.5,
  "status": "completed"
}
```

**注意**: TTS 服务待集成 (如 Google TTS、Azure TTS)

---

### 6.5 POST /api/v1/generation/batch-scenes

**认证**: 需要 JWT

批量生成场景内容

**请求体**
```json
{
  "scene_ids": ["scene_001", "scene_002", "scene_003"],
  "content_types": ["image", "video"],
  "options": {
    "parallel": true,
    "max_concurrent": 3
  }
}
```

**参数说明**
- `content_types` - 生成内容类型数组
- `parallel` - 是否并行处理
- `max_concurrent` - 最大并发数

**响应示例**
```json
{
  "batch_task_id": "batch_001",
  "total_scenes": 3,
  "status": "processing",
  "progress": {
    "completed": 0,
    "failed": 0,
    "pending": 3
  },
  "estimated_time": 120
}
```

**业务逻辑**
- 创建批量任务
- 使用 Go 协程并发处理
- 使用 semaphore 控制并发数
- 依次调用场景图片、视频生成
- 实时更新进度

---

### 6.6 GET /api/v1/generation/tasks/:task_id

**认证**: 需要 JWT

查询生成任务状态

**路径参数**
- `task_id` (required) - 任务 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/generation/tasks/batch_001
```

**响应示例**
```json
{
  "task_id": "batch_001",
  "type": "batch_scenes",
  "status": "processing",
  "progress": {
    "total": 3,
    "completed": 1,
    "failed": 0,
    "current": "scene_002"
  },
  "results": [
    {
      "scene_id": "scene_001",
      "status": "completed",
      "media_ids": ["media_001", "media_002"]
    },
    {
      "scene_id": "scene_002",
      "status": "processing",
      "media_ids": []
    }
  ],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:05:00Z"
}
```

**状态说明**
- `pending` - 等待中
- `processing` - 处理中
- `completed` - 已完成
- `failed` - 失败

---

## 7. 项目管理

### 7.1 POST /api/v1/projects

**认证**: 需要 JWT

创建动漫项目

**请求体**
```json
{
  "novel_id": "novel_abc123",
  "title": "修仙传动画版",
  "chapters": [1, 2, 3],
  "style": "anime",
  "settings": {
    "image_quality": "high",
    "video_duration": 5,
    "voice_enabled": true
  }
}
```

**响应示例**
```json
{
  "project_id": "proj_001",
  "novel_id": "novel_abc123",
  "title": "修仙传动画版",
  "status": "created",
  "created_at": "2024-01-01T12:00:00Z"
}
```

---

### 7.2 GET /api/v1/projects/:project_id

**认证**: 需要 JWT

获取项目详情

**路径参数**
- `project_id` (required) - 项目 ID

**响应示例**
```json
{
  "id": "proj_001",
  "novel_id": "novel_abc123",
  "title": "修仙传动画版",
  "status": "in_progress",
  "progress": {
    "total_scenes": 150,
    "generated_images": 100,
    "generated_videos": 50,
    "generated_voices": 200
  },
  "chapters": [1, 2, 3],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:30:00Z"
}
```

**状态说明**
- `created` - 已创建
- `in_progress` - 进行中
- `completed` - 已完成

---

### 7.3 POST /api/v1/projects/:project_id/generate-all

**认证**: 需要 JWT

一键生成项目所有内容

**路径参数**
- `project_id` (required) - 项目 ID

**响应示例**
```json
{
  "task_id": "gen_all_001",
  "status": "started",
  "estimated_time": 3600
}
```

**业务逻辑**
1. 生成所有角色参考图
2. 为所有场景生成图片
3. 为所有场景生成视频
4. 为所有对话生成配音
5. 分步异步处理，可中断恢复

---

## 8. 导出功能

### 8.1 POST /api/v1/export/video

**认证**: 需要 JWT

导出完整视频

**请求体**
```json
{
  "project_id": "proj_001",
  "chapters": [1, 2, 3],
  "options": {
    "format": "mp4",
    "resolution": "1920x1080",
    "quality": "high",
    "include_subtitles": true,
    "include_voice": true
  }
}
```

**响应示例**
```json
{
  "export_id": "exp_001",
  "status": "processing",
  "estimated_time": 300,
  "download_url": null
}
```

**业务逻辑**
1. 按章节顺序收集所有场景
2. 合并场景视频
3. 叠加配音轨道
4. 添加字幕 (SRT 格式)
5. 使用 FFmpeg 渲染最终视频
6. 保存到存储
7. 返回下载链接

---

### 8.2 GET /api/v1/export/:export_id

**认证**: 需要 JWT

查询导出任务状态

**路径参数**
- `export_id` (required) - 导出任务 ID

**响应示例**
```json
{
  "export_id": "exp_001",
  "status": "completed",
  "progress": 100,
  "download_url": "https://storage.example.com/exports/exp_001.mp4",
  "file_size": 1024000000,
  "duration": 1800,
  "created_at": "2024-01-01T12:00:00Z",
  "expires_at": "2024-01-08T12:00:00Z"
}
```

**状态说明**
- `processing` - 处理中
- `completed` - 已完成
- `failed` - 失败

---

### 8.3 POST /api/v1/export/scenes

**认证**: 需要 JWT

批量导出场景素材

**请求体**
```json
{
  "scene_ids": ["scene_001", "scene_002"],
  "types": ["image", "video", "audio"],
  "format": "zip"
}
```

**响应示例**
```json
{
  "export_id": "exp_002",
  "download_url": "https://storage.example.com/exports/scenes_exp_002.zip",
  "file_size": 50000000
}
```

---

## 通用规范

### 错误响应格式

所有错误响应遵循统一格式：

```json
{
  "error": {
    "code": "INVALID_NOVEL_ID",
    "message": "小说 ID 不存在",
    "details": {
      "novel_id": "invalid_id"
    },
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### HTTP 状态码

- `200` - 请求成功
- `201` - 资源创建成功
- `400` - 请求参数错误
- `401` - 未认证
- `403` - 无权限
- `404` - 资源不存在
- `409` - 资源冲突 (如重复上传)
- `422` - 业务逻辑错误 (如角色未生成参考图)
- `429` - 请求过于频繁
- `500` - 服务器内部错误
- `503` - AI 服务不可用

### 常见错误码

**认证相关**
- `UNAUTHORIZED` - 未认证,缺少或无效的 JWT token
- `TOKEN_EXPIRED` - Token 已过期
- `INVALID_TOKEN` - Token 无效或格式错误
- `INVALID_CREDENTIALS` - 用户名或密码错误
- `USER_NOT_FOUND` - 用户不存在
- `USERNAME_EXISTS` - 用户名已存在
- `EMAIL_EXISTS` - 邮箱已存在
- `INVALID_USERNAME` - 用户名格式不正确
- `INVALID_EMAIL` - 邮箱格式不正确
- `WEAK_PASSWORD` - 密码强度不足
- `INVALID_OLD_PASSWORD` - 旧密码错误

**业务相关**
- `INVALID_PARAMETER` - 请求参数无效
- `RESOURCE_NOT_FOUND` - 资源不存在
- `DUPLICATE_RESOURCE` - 资源已存在
- `INVALID_FILE_FORMAT` - 文件格式不支持
- `FILE_TOO_LARGE` - 文件过大
- `PARSE_FAILED` - 解析失败
- `GENERATION_FAILED` - 生成失败
- `AI_SERVICE_UNAVAILABLE` - AI 服务不可用
- `INSUFFICIENT_CREDITS` - 额度不足 (未来版本)

---

## 异步任务模式

长时间任务 (如解析、生成、导出) 采用异步模式：

1. POST 请求返回 `task_id`
2. 客户端轮询 GET `/api/v1/generation/tasks/:task_id`
3. 状态流转：`pending` → `processing` → `completed` / `failed`
4. 建议轮询间隔：2-5 秒

**未来版本**: 支持 WebSocket 实时推送进度

---

## 接口开发优先级

### P0 (核心功能，MVP 必需)
1. ✅ GET /health
2. POST /api/v1/auth/register
3. POST /api/v1/auth/login
4. POST /api/v1/novels/upload
5. POST /api/v1/novels/:novel_id/parse
6. GET /api/v1/novels/:novel_id/characters
7. POST /api/v1/characters/:character_id/references
8. POST /api/v1/generation/scene-image
9. GET /api/v1/novels/:novel_id/scenes

### P1 (重要功能)
10. POST /api/v1/generation/scene-video
11. POST /api/v1/generation/batch-scenes
12. GET /api/v1/generation/tasks/:task_id
13. GET /api/v1/novels/:novel_id
14. GET /api/v1/scenes/:scene_id
15. GET /api/v1/auth/me
16. POST /api/v1/auth/refresh

### P2 (增强功能)
17. POST /api/v1/projects
18. POST /api/v1/export/video
19. POST /api/v1/generation/voice
20. PUT /api/v1/characters/:character_id
21. PUT /api/v1/scenes/:scene_id
22. POST /api/v1/auth/logout
23. PUT /api/v1/auth/password

### P3 (优化功能)
24. GET /api/v1/novels (列表)
25. DELETE /api/v1/novels/:novel_id
26. POST /api/v1/export/scenes

---

## 相关文档

- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构设计
- [CHARACTER_CONSISTENCY.md](CHARACTER_CONSISTENCY.md) - 角色一致性设计
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南
- [QUICKSTART.md](../QUICKSTART.md) - 快速开始

---

*API 文档版本：v0.1.0-alpha*
*最后更新：2024-01-01*
