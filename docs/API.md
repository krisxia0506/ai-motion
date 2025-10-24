# API 文档

## 概述

AI-Motion 提供 RESTful API 接口,用于小说解析、角色管理、场景管理、内容生成和导出功能。

**版本**: v0.1.0-alpha
**Base URL**:
- 开发环境:`http://localhost:8080`
- 生产环境:根据部署配置而定

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
- `POST /api/v1/auth/refresh`

### 认证接口 (需要 JWT)

所有其他接口均需要在 Header 中携带有效的 JWT token。

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
- `data`: 响应数据载体,成功时包含业务数据,失败时可为 `null`

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
| 10003 | 资源已存在              | 重复创建                                |
| 10004 | 资源状态不正确          | 操作不符合当前资源状态                   |
| 20001 | 认证失败                | Token 无效或过期                        |
| 20002 | 权限不足                | 无权限操作资源                          |
| 20003 | 用户名或密码错误        | 登录失败                                |
| 30001 | 文件上传失败            | 文件格式错误、大小超限                   |
| 30002 | 文件解析失败            | 文件格式不正确或内容无法解析             |
| 40001 | AI 服务调用失败         | Gemini/Sora API 错误                    |
| 40002 | AI 服务不可用           | AI 服务暂时不可用                       |
| 40003 | 生成任务失败            | 图像/视频生成失败                       |
| 50001 | 数据库错误              | 数据库操作失败                          |
| 50002 | 系统内部错误            | 未知错误                                |
| 50003 | 第三方服务错误          | 外部服务调用失败                        |

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
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "timestamp": "2024-01-01T12:00:00Z",
    "services": {
      "database": "connected",
      "geminiApi": "available",
      "soraApi": "available"
    }
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
  "code": 0,
  "message": "注册成功",
  "data": {
    "userId": "user_abc123",
    "username": "user123",
    "email": "user@example.com",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

**业务逻辑**
1. 验证用户名、邮箱、密码格式
2. 检查用户名和邮箱是否已存在
3. 使用 bcrypt 加密密码(cost=10)
4. 创建 User 实体
5. 保存到数据库
6. 返回用户基本信息

**错误示例**
```json
{
  "code": 10003,
  "message": "用户名已存在",
  "data": null
}
```

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
  "code": 0,
  "message": "登录成功",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyX2FiYzEyMyIsInVzZXJuYW1lIjoidXNlcjEyMyIsImV4cCI6MTcwNDE5NjgwMH0.xxx",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyX2FiYzEyMyIsInR5cGUiOiJyZWZyZXNoIiwiZXhwIjoxNzA2Nzg4ODAwfQ.yyy",
    "tokenType": "Bearer",
    "expiresIn": 604800,
    "user": {
      "userId": "user_abc123",
      "username": "user123",
      "email": "user@example.com"
    }
  }
}
```

**响应字段说明**
- `accessToken` - 访问令牌,用于 API 请求
- `refreshToken` - 刷新令牌,用于获取新的访问令牌
- `tokenType` - 令牌类型,固定为 "Bearer"
- `expiresIn` - 访问令牌过期时间(秒),默认 7 天

**业务逻辑**
1. 查找用户(支持用户名或邮箱登录)
2. 验证密码(使用 bcrypt.CompareHashAndPassword)
3. 生成 JWT accessToken:
   - Payload: `userId`, `username`, `exp`(过期时间)
   - 签名算法: HS256
   - 过期时间: 7天
4. 生成 JWT refreshToken:
   - Payload: `userId`, `type: "refresh"`, `exp`(过期时间)
   - 过期时间: 30天
5. 返回 token 和用户信息

**错误示例**
```json
{
  "code": 20003,
  "message": "用户名或密码错误",
  "data": null
}
```

---

### 2.3 POST /api/v1/auth/refresh

刷新访问令牌

**请求参数**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "Token 刷新成功",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_token",
    "tokenType": "Bearer",
    "expiresIn": 604800
  }
}
```

**业务逻辑**
1. 验证 refreshToken 是否有效
2. 解析 token 获取 userId
3. 检查用户是否存在
4. 生成新的 accessToken
5. 返回新 token

**错误示例**
```json
{
  "code": 20001,
  "message": "Token 无效或已过期",
  "data": null
}
```

---

### 2.4 POST /api/v1/auth/logout

**认证**: 需要 JWT

用户登出

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "登出成功",
  "data": null
}
```

**业务逻辑**
1. 解析 JWT token 获取 userId
2. (可选)将 token 加入黑名单(需要 Redis)
3. 返回成功消息

**注意**: 由于 JWT 是无状态的,客户端需要删除本地存储的 token

---

### 2.5 GET /api/v1/auth/me

**认证**: 需要 JWT

获取当前用户信息

**请求示例**
```bash
curl -X GET \
  http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userId": "user_abc123",
    "username": "user123",
    "email": "user@example.com",
    "createdAt": "2024-01-01T12:00:00Z",
    "usageStats": {
      "novelsCount": 5,
      "charactersCount": 30,
      "scenesGenerated": 150
    }
  }
}
```

**业务逻辑**
1. 从 JWT token 中解析 userId
2. 查询用户信息
3. 统计用户使用情况
4. 返回完整用户信息

---

### 2.6 PUT /api/v1/auth/password

**认证**: 需要 JWT

修改密码

**请求参数**
```json
{
  "oldPassword": "oldPassword123",
  "newPassword": "newPassword456"
}
```

**请求示例**
```bash
curl -X PUT \
  http://localhost:8080/api/v1/auth/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "oldPassword": "oldPassword123",
    "newPassword": "newPassword456"
  }'
```

**响应示例**
```json
{
  "code": 0,
  "message": "密码修改成功",
  "data": null
}
```

**业务逻辑**
1. 从 JWT 获取 userId
2. 验证旧密码是否正确
3. 验证新密码强度
4. 使用 bcrypt 加密新密码
5. 更新数据库
6. (可选)使所有旧 token 失效

**错误示例**
```json
{
  "code": 10001,
  "message": "旧密码错误",
  "data": null
}
```

---

### JWT Token 结构

**Access Token Payload**
```json
{
  "userId": "user_abc123",
  "username": "user123",
  "iat": 1704110400,
  "exp": 1704715200
}
```

**Refresh Token Payload**
```json
{
  "userId": "user_abc123",
  "type": "refresh",
  "iat": 1704110400,
  "exp": 1706788800
}
```

**签名密钥**: 使用环境变量 `JWT_SECRET_KEY`,至少 32 字符

---

## 3. 小说管理

### 3.1 POST /api/v1/novels/upload

**认证**: 需要 JWT

上传小说文件

**请求参数**
- `file` (multipart/form-data, required) - 小说文件,支持 TXT 格式
- `title` (form, optional) - 小说标题
- `author` (form, optional) - 作者名称

**幂等性**: 支持 `Idempotency-Key` 请求头,防止重复上传

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novels/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Idempotency-Key: uuid-generated-by-client" \
  -F "file=@novel.txt" \
  -F "title=修仙传" \
  -F "author=作者名"
```

**响应示例**
```json
{
  "code": 0,
  "message": "小说上传成功",
  "data": {
    "novelId": "novel_abc123",
    "title": "修仙传",
    "filename": "novel.txt",
    "fileSize": 1024000,
    "status": "uploaded",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

**业务逻辑**
1. 验证文件格式 (仅允许 .txt)
2. 验证文件大小限制 (最大 50MB)
3. 生成唯一 Novel ID
4. 保存文件到临时存储
5. 创建 Novel 实体,状态为 "uploaded"
6. 保存到数据库
7. 返回小说基本信息

**错误示例**
```json
{
  "code": 30001,
  "message": "文件格式不支持",
  "data": {
    "errorDetail": "仅支持 .txt 格式"
  }
}
```

---

### 3.2 POST /api/v1/novels/:novelId/parse

**认证**: 需要 JWT

解析小说内容,提取章节、角色、场景

**路径参数**
- `novelId` (required) - 小说 ID

**请求体**
```json
{
  "options": {
    "autoGenerateReferences": true,
    "sceneDivisionMode": "auto"
  }
}
```

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/novels/novel_abc123/parse \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "options": {
      "autoGenerateReferences": true,
      "sceneDivisionMode": "auto"
    }
  }'
```

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "解析任务已创建",
  "data": {
    "taskId": "task_abc123",
    "novelId": "novel_abc123",
    "status": "pending",
    "estimatedTime": 300
  }
}
```

**业务逻辑**
1. 读取小说原始文本
2. 创建异步解析任务
3. 调用 Gemini API 进行自然语言处理
4. 提取章节结构 → 创建 Chapter 实体
5. 提取角色信息 → 创建 Character 实体
6. 划分场景 → 创建 Scene 实体
7. 分析对话和场景描述
8. 更新 Novel 状态为 "parsed"
9. (可选) 异步触发角色参考图生成

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询解析进度

**注意**: 这是一个耗时操作,采用异步处理模式

---

### 3.3 GET /api/v1/novels/:novelId

**认证**: 需要 JWT

获取小说详细信息

**路径参数**
- `novelId` (required) - 小说 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/novels/novel_abc123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "novel_abc123",
    "title": "修仙传",
    "author": "作者名",
    "status": "parsed",
    "chaptersCount": 50,
    "charactersCount": 12,
    "scenesCount": 200,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:30:00Z"
  }
}
```

**状态值说明**
- `uploaded` - 已上传
- `parsing` - 解析中
- `parsed` - 解析完成
- `failed` - 解析失败

**错误示例**
```json
{
  "code": 10002,
  "message": "小说不存在",
  "data": null
}
```

---

### 3.4 GET /api/v1/novels

**认证**: 需要 JWT

获取小说列表

**查询参数**
- `page` (optional, default: 1) - 页码
- `pageSize` (optional, default: 20) - 每页数量
- `status` (optional) - 过滤状态

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels?page=1&pageSize=20&status=parsed" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "novel_abc123",
        "title": "修仙传",
        "author": "作者名",
        "status": "parsed",
        "createdAt": "2024-01-01T12:00:00Z"
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

### 3.5 DELETE /api/v1/novels/:novelId

**认证**: 需要 JWT

删除小说及关联数据

**路径参数**
- `novelId` (required) - 小说 ID

**请求示例**
```bash
curl -X DELETE http://localhost:8080/api/v1/novels/novel_abc123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "小说删除成功",
  "data": {
    "deletedItems": {
      "chapters": 50,
      "characters": 12,
      "scenes": 200,
      "media": 150
    }
  }
}
```

**业务逻辑**
- 级联删除所有章节、角色、场景、媒体文件
- 删除存储的文件 (参考图、场景图、视频)
- 返回删除统计

---

## 4. 角色管理

### 4.1 GET /api/v1/novels/:novelId/characters

**认证**: 需要 JWT

获取小说的所有角色

**路径参数**
- `novelId` (required) - 小说 ID

**查询参数**
- `includeReferences` (optional, default: true) - 是否包含参考图

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels/novel_abc123/characters?includeReferences=true" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
        "description": "女主角,18 岁,黑色长发,明亮的眼睛,身穿白色长裙",
        "role": "protagonist",
        "appearancesCount": 45,
        "referenceImages": [
          {
            "id": "ref_img_001",
            "url": "https://storage.example.com/ref_001.jpg",
            "state": "default",
            "createdAt": "2024-01-01T12:00:00Z"
          }
        ],
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

**角色类型说明**
- `protagonist` - 主角
- `antagonist` - 反派
- `supporting` - 配角

---

### 4.2 GET /api/v1/characters/:characterId

**认证**: 需要 JWT

获取单个角色详情

**路径参数**
- `characterId` (required) - 角色 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/characters/char_001 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "char_001",
    "novelId": "novel_abc123",
    "name": "李雪",
    "description": "女主角,18 岁...",
    "appearance": {
      "age": "18",
      "gender": "female",
      "hair": "黑色长发",
      "eyes": "明亮的黑色眼睛",
      "clothing": "白色长裙"
    },
    "personality": {
      "traits": ["勇敢", "善良", "坚韧"],
      "description": "性格开朗,勇敢面对困难"
    },
    "referenceImages": [
      {
        "id": "ref_img_001",
        "url": "https://storage.example.com/ref_001.jpg",
        "state": "default",
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ],
    "scenes": ["scene_001", "scene_005"],
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

---

### 4.3 POST /api/v1/characters/:characterId/references

**认证**: 需要 JWT

为角色生成参考图 (核心功能:角色一致性)

**路径参数**
- `characterId` (required) - 角色 ID

**请求体**
```json
{
  "state": "default",
  "customPrompt": "",
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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "state": "default",
    "style": "anime"
  }'
```

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "参考图生成任务已创建",
  "data": {
    "taskId": "task_ref_001",
    "referenceId": "ref_img_002",
    "characterId": "char_001",
    "status": "pending",
    "estimatedTime": 10
  }
}
```

**业务逻辑**
1. 获取 Character 实体
2. 构建提示词 (使用角色 description + state + style)
3. 调用 Gemini TextToImage API
4. 保存图片到存储服务
5. 更新 Character.ReferenceImages
6. 返回生成结果

**角色一致性说明**: 此接口生成的参考图将用于后续所有场景生成,确保角色视觉一致性。

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询生成进度

---

### 4.4 PUT /api/v1/characters/:characterId

**认证**: 需要 JWT

更新角色信息

**路径参数**
- `characterId` (required) - 角色 ID

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
  "code": 0,
  "message": "角色信息更新成功",
  "data": {
    "id": "char_001",
    "name": "李雪",
    "description": "更新后的描述",
    "updatedAt": "2024-01-01T13:00:00Z"
  }
}
```

**注意**: 如果描述有重大变化,建议重新生成参考图

---

### 4.5 DELETE /api/v1/characters/:characterId/references/:referenceId

**认证**: 需要 JWT

删除角色参考图

**路径参数**
- `characterId` (required) - 角色 ID
- `referenceId` (required) - 参考图 ID

**响应示例**
```json
{
  "code": 0,
  "message": "参考图删除成功",
  "data": null
}
```

---

## 5. 场景管理

### 5.1 GET /api/v1/novels/:novelId/scenes

**认证**: 需要 JWT

获取小说的所有场景

**路径参数**
- `novelId` (required) - 小说 ID

**查询参数**
- `chapterId` (optional) - 过滤章节
- `page` (optional, default: 1) - 页码
- `pageSize` (optional, default: 20) - 每页数量

**请求示例**
```bash
curl "http://localhost:8080/api/v1/novels/novel_abc123/scenes?chapterId=chapter_01&page=1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
        "novelId": "novel_abc123",
        "chapterId": "chapter_01",
        "sequenceNum": 1,
        "description": "清晨的竹林,阳光透过竹叶洒下斑驳的光影",
        "location": "竹林",
        "timeOfDay": "清晨",
        "characters": ["char_001", "char_003"],
        "dialogues": [
          {
            "characterId": "char_001",
            "text": "今天天气真好",
            "emotion": "happy"
          }
        ],
        "duration": 5.0,
        "media": {
          "image": "https://storage.example.com/scene_001.jpg",
          "video": null
        },
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 200,
      "totalPages": 10,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

---

### 5.2 GET /api/v1/scenes/:sceneId

**认证**: 需要 JWT

获取场景详情

**路径参数**
- `sceneId` (required) - 场景 ID

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "scene_001",
    "novelId": "novel_abc123",
    "chapterId": "chapter_01",
    "sequenceNum": 1,
    "description": "清晨的竹林,阳光透过竹叶洒下斑驳的光影",
    "location": "竹林",
    "timeOfDay": "清晨",
    "characters": [
      {
        "id": "char_001",
        "name": "李雪",
        "description": "女主角..."
      }
    ],
    "dialogues": [
      {
        "characterId": "char_001",
        "text": "今天天气真好",
        "emotion": "happy"
      }
    ],
    "media": {
      "images": [
        {
          "id": "media_001",
          "url": "https://storage.example.com/scene_001.jpg",
          "createdAt": "2024-01-01T12:00:00Z"
        }
      ],
      "videos": []
    },
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

---

### 5.3 PUT /api/v1/scenes/:sceneId

**认证**: 需要 JWT

更新场景信息

**路径参数**
- `sceneId` (required) - 场景 ID

**请求体**
```json
{
  "description": "更新后的场景描述",
  "location": "新地点",
  "timeOfDay": "黄昏",
  "dialogues": [
    {
      "characterId": "char_001",
      "text": "太阳快要下山了",
      "emotion": "peaceful"
    }
  ]
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "场景信息更新成功",
  "data": {
    "id": "scene_001",
    "description": "更新后的场景描述",
    "updatedAt": "2024-01-01T13:00:00Z"
  }
}
```

---

### 5.4 POST /api/v1/scenes/:sceneId/prompt

**认证**: 需要 JWT

生成场景的 AI 提示词

**路径参数**
- `sceneId` (required) - 场景 ID

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sceneId": "scene_001",
    "prompt": {
      "textToImage": "清晨的竹林,阳光透过竹叶,一位黑发年轻女子站在竹林中,动漫风格,高质量,细节丰富",
      "imageToVideo": "镜头缓慢推进,竹叶在微风中摇曳,女子转头看向远方"
    }
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
  "characterIds": ["char_001", "char_002"],
  "options": {
    "style": "anime",
    "quality": "high"
  }
}
```

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/generation/character-reference \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "characterIds": ["char_001", "char_002"],
    "options": {
      "style": "anime",
      "quality": "high"
    }
  }'
```

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "批量生成任务已创建",
  "data": {
    "taskId": "gen_task_001",
    "status": "pending",
    "total": 2,
    "completed": 0,
    "estimatedTime": 20
  }
}
```

**业务逻辑**
- 创建异步任务
- 并发调用 Gemini TextToImage API
- 使用 Go 协程处理多个角色
- 保存结果到 Character 实体

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询生成进度

---

### 6.2 POST /api/v1/generation/scene-image

**认证**: 需要 JWT

生成场景图片 (带角色一致性)

**请求体**
```json
{
  "sceneId": "scene_001",
  "options": {
    "style": "anime",
    "quality": "high",
    "consistencyStrength": 0.8
  }
}
```

**参数说明**
- `consistencyStrength` - 角色一致性强度,取值 0-1,越高越严格保持角色外观

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/generation/scene-image \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "sceneId": "scene_001",
    "options": {
      "style": "anime",
      "quality": "high",
      "consistencyStrength": 0.8
    }
  }'
```

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "场景图片生成任务已创建",
  "data": {
    "taskId": "task_img_001",
    "mediaId": "media_001",
    "sceneId": "scene_001",
    "status": "pending",
    "estimatedTime": 10
  }
}
```

**业务逻辑**
1. 获取 Scene 实体
2. 获取场景中的所有角色
3. 收集角色参考图 (ReferenceImages[0])
4. 构建提示词 (场景描述 + 角色特征)
5. 调用 Gemini ImageToImage API:
   - prompt: 场景描述
   - referenceImages: 角色参考图数组
   - consistencyStrength: 一致性参数
6. 保存图片到存储
7. 创建 Media 实体,type="image"
8. 关联到 Scene
9. 返回结果

**角色一致性说明**: 此接口使用角色参考图进行图生图,确保场景中的角色与参考图保持视觉一致。

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询生成进度

---

### 6.3 POST /api/v1/generation/scene-video

**认证**: 需要 JWT

生成场景视频 (图生视频)

**请求体**
```json
{
  "sceneId": "scene_001",
  "sourceImage": "media_001",
  "options": {
    "duration": 5,
    "fps": 30,
    "motionType": "smooth"
  }
}
```

**参数说明**
- `sourceImage` - 可选,指定源图片 ID,默认使用场景已生成的图片
- `duration` - 视频时长 (秒)
- `fps` - 帧率
- `motionType` - 动作类型:`smooth`(平滑) | `dynamic`(动态)

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "场景视频生成任务已创建",
  "data": {
    "taskId": "task_vid_001",
    "mediaId": "media_002",
    "sceneId": "scene_001",
    "status": "pending",
    "estimatedTime": 60
  }
}
```

**业务逻辑**
1. 获取 Scene 实体
2. 获取场景图片 (如果未指定,使用最新生成的图片)
3. 调用 Sora2 ImageToVideo API
4. 保存视频到存储
5. 创建 Media 实体,type="video"
6. 返回结果

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询生成进度

---

### 6.4 POST /api/v1/generation/voice

**认证**: 需要 JWT

生成角色配音

**请求体**
```json
{
  "characterId": "char_001",
  "text": "今天天气真好",
  "options": {
    "emotion": "happy",
    "voiceProfile": "female_young",
    "speed": 1.0
  }
}
```

**参数说明**
- `emotion` - 情感:`neutral` | `happy` | `sad` | `angry` | `surprised`
- `voiceProfile` - 声音配置:`female_young` | `male_young` | `female_mature` | `male_mature`
- `speed` - 语速倍率,默认 1.0

**响应示例**
```json
{
  "code": 0,
  "message": "配音生成成功",
  "data": {
    "audioId": "audio_001",
    "characterId": "char_001",
    "audioUrl": "https://storage.example.com/audio_001.mp3",
    "duration": 2.5,
    "status": "completed"
  }
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
  "sceneIds": ["scene_001", "scene_002", "scene_003"],
  "contentTypes": ["image", "video"],
  "options": {
    "parallel": true,
    "maxConcurrent": 3
  }
}
```

**参数说明**
- `contentTypes` - 生成内容类型数组
- `parallel` - 是否并行处理
- `maxConcurrent` - 最大并发数

**请求示例**
```bash
curl -X POST \
  http://localhost:8080/api/v1/generation/batch-scenes \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "sceneIds": ["scene_001", "scene_002", "scene_003"],
    "contentTypes": ["image", "video"],
    "options": {
      "parallel": true,
      "maxConcurrent": 3
    }
  }'
```

**响应示例 (批量任务)**
```json
{
  "code": 0,
  "message": "批量生成任务已创建",
  "data": {
    "batchId": "batch_abc123",
    "tasks": [
      {"sceneId": "scene_001", "taskId": "task_1", "status": "pending"},
      {"sceneId": "scene_002", "taskId": "task_2", "status": "pending"},
      {"sceneId": "scene_003", "taskId": "task_3", "status": "pending"}
    ],
    "totalScenes": 3,
    "status": "pending",
    "progress": {
      "completed": 0,
      "failed": 0,
      "pending": 3
    },
    "estimatedTime": 120
  }
}
```

**业务逻辑**
- 创建批量任务
- 使用 Go 协程并发处理
- 使用 semaphore 控制并发数
- 依次调用场景图片、视频生成
- 实时更新进度

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询批量任务进度

---

### 6.6 GET /api/v1/tasks/:taskId

**认证**: 需要 JWT

查询异步任务状态 (通用任务查询接口)

**路径参数**
- `taskId` (required) - 任务 ID

**请求示例**
```bash
curl http://localhost:8080/api/v1/tasks/task_abc123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例 (单个任务)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "taskId": "task_abc123",
    "type": "scene_image",
    "status": "processing",
    "progress": 65,
    "result": null,
    "error": null,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:05:00Z"
  }
}
```

**响应示例 (批量任务)**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "taskId": "batch_001",
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
        "sceneId": "scene_001",
        "status": "completed",
        "mediaIds": ["media_001", "media_002"]
      },
      {
        "sceneId": "scene_002",
        "status": "processing",
        "mediaIds": []
      },
      {
        "sceneId": "scene_003",
        "status": "pending",
        "mediaIds": []
      }
    ],
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:05:00Z"
  }
}
```

**状态说明**
- `pending` - 等待中
- `processing` - 处理中
- `completed` - 已完成
- `failed` - 失败

**轮询建议**: 轮询间隔 2-5 秒

---

## 7. 项目管理

### 7.1 POST /api/v1/projects

**认证**: 需要 JWT

创建动漫项目

**请求体**
```json
{
  "novelId": "novel_abc123",
  "title": "修仙传动画版",
  "chapters": [1, 2, 3],
  "style": "anime",
  "settings": {
    "imageQuality": "high",
    "videoDuration": 5,
    "voiceEnabled": true
  }
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "项目创建成功",
  "data": {
    "projectId": "proj_001",
    "novelId": "novel_abc123",
    "title": "修仙传动画版",
    "status": "created",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

---

### 7.2 GET /api/v1/projects/:projectId

**认证**: 需要 JWT

获取项目详情

**路径参数**
- `projectId` (required) - 项目 ID

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "proj_001",
    "novelId": "novel_abc123",
    "title": "修仙传动画版",
    "status": "inProgress",
    "progress": {
      "totalScenes": 150,
      "generatedImages": 100,
      "generatedVideos": 50,
      "generatedVoices": 200
    },
    "chapters": [1, 2, 3],
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:30:00Z"
  }
}
```

**状态说明**
- `created` - 已创建
- `inProgress` - 进行中
- `completed` - 已完成

---

### 7.3 POST /api/v1/projects/:projectId/generate-all

**认证**: 需要 JWT

一键生成项目所有内容

**路径参数**
- `projectId` (required) - 项目 ID

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "全量生成任务已创建",
  "data": {
    "taskId": "gen_all_001",
    "status": "pending",
    "estimatedTime": 3600
  }
}
```

**业务逻辑**
1. 生成所有角色参考图
2. 为所有场景生成图片
3. 为所有场景生成视频
4. 为所有对话生成配音
5. 分步异步处理,可中断恢复

**任务查询**: 使用 `GET /api/v1/tasks/:taskId` 查询生成进度

---

## 8. 导出功能

### 8.1 POST /api/v1/export/video

**认证**: 需要 JWT

导出完整视频

**请求体**
```json
{
  "projectId": "proj_001",
  "chapters": [1, 2, 3],
  "options": {
    "format": "mp4",
    "resolution": "1920x1080",
    "quality": "high",
    "includeSubtitles": true,
    "includeVoice": true
  }
}
```

**响应示例 (异步任务)**
```json
{
  "code": 0,
  "message": "导出任务已创建",
  "data": {
    "exportId": "exp_001",
    "status": "pending",
    "estimatedTime": 300,
    "downloadUrl": null
  }
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

### 8.2 GET /api/v1/export/:exportId

**认证**: 需要 JWT

查询导出任务状态

**路径参数**
- `exportId` (required) - 导出任务 ID

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "exportId": "exp_001",
    "status": "completed",
    "progress": 100,
    "downloadUrl": "https://storage.example.com/exports/exp_001.mp4",
    "fileSize": 1024000000,
    "duration": 1800,
    "createdAt": "2024-01-01T12:00:00Z",
    "expiresAt": "2024-01-08T12:00:00Z"
  }
}
```

**状态说明**
- `pending` - 等待中
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
  "sceneIds": ["scene_001", "scene_002"],
  "types": ["image", "video", "audio"],
  "format": "zip"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "场景素材导出成功",
  "data": {
    "exportId": "exp_002",
    "downloadUrl": "https://storage.example.com/exports/scenes_exp_002.zip",
    "fileSize": 50000000
  }
}
```

---

## 幂等性设计

对于创建类操作,支持使用 `Idempotency-Key` 请求头确保幂等性:

**示例**:
```bash
POST /api/v1/novels/upload
Headers:
  Authorization: Bearer {token}
  Idempotency-Key: uuid-client-generated
```

**行为**: 在 24 小时内,相同 `Idempotency-Key` 的请求返回相同结果,避免重复创建

**适用场景**:
- 小说上传
- 项目创建
- 生成任务创建

---

## 异步任务模式

长时间任务 (如解析、生成、导出) 采用异步模式:

1. POST 请求返回 `taskId`
2. 客户端轮询 GET `/api/v1/tasks/:taskId`
3. 状态流转:`pending` → `processing` → `completed` / `failed`
4. 建议轮询间隔:2-5 秒

**未来版本**: 支持 WebSocket 实时推送进度

---

## HTTP 状态码

- `200 OK` - 请求成功 (包括业务逻辑错误,通过 code 区分)
- `201 Created` - 资源创建成功
- `204 No Content` - 请求成功但无返回内容
- `400 Bad Request` - 请求格式错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 无权限
- `404 Not Found` - 路由不存在
- `429 Too Many Requests` - 请求过于频繁
- `500 Internal Server Error` - 服务器内部错误
- `503 Service Unavailable` - AI 服务不可用

**注意**: 业务逻辑错误统一返回 HTTP 200,通过 `code` 字段区分具体错误

---

## 接口开发优先级

### P0 (核心功能,MVP 必需)
1. ✅ GET /health
2. POST /api/v1/auth/register
3. POST /api/v1/auth/login
4. POST /api/v1/novels/upload
5. POST /api/v1/novels/:novelId/parse
6. GET /api/v1/novels/:novelId/characters
7. POST /api/v1/characters/:characterId/references
8. POST /api/v1/generation/scene-image
9. GET /api/v1/novels/:novelId/scenes
10. GET /api/v1/tasks/:taskId

### P1 (重要功能)
11. POST /api/v1/generation/scene-video
12. POST /api/v1/generation/batch-scenes
13. GET /api/v1/novels/:novelId
14. GET /api/v1/scenes/:sceneId
15. GET /api/v1/auth/me
16. POST /api/v1/auth/refresh

### P2 (增强功能)
17. POST /api/v1/projects
18. POST /api/v1/export/video
19. POST /api/v1/generation/voice
20. PUT /api/v1/characters/:characterId
21. PUT /api/v1/scenes/:sceneId
22. POST /api/v1/auth/logout
23. PUT /api/v1/auth/password

### P3 (优化功能)
24. GET /api/v1/novels (列表)
25. DELETE /api/v1/novels/:novelId
26. POST /api/v1/export/scenes

---

## 相关文档

- [API_DESIGN_GUIDELINES.md](API_DESIGN_GUIDELINES.md) - API 设计规范
- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构设计
- [CHARACTER_CONSISTENCY.md](CHARACTER_CONSISTENCY.md) - 角色一致性设计
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南
- [QUICKSTART.md](../QUICKSTART.md) - 快速开始

---

*API 文档版本:v0.1.0-alpha*
*最后更新:2024-01-01*
*符合 API 设计规范 v1.0*
