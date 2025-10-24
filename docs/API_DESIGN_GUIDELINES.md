# AI-Motion API 设计规范

本文档定义 AI-Motion 项目的 API 设计规范,确保 API 的一致性、可维护性和可扩展性。

**版本**: v1.0  
**更新日期**: 2024-01-01

---

## RESTful 设计原则

- 使用标准 HTTP 方法: GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)
- 资源命名使用复数名词: `/api/v1/novels`, `/api/v1/characters`, `/api/v1/scenes`
- 路径参数用于资源标识: `/api/v1/novels/:novelId`, `/api/v1/characters/:characterId`
- 查询参数用于过滤和分页: `?page=1&pageSize=20&status=parsed`
- 嵌套资源表示从属关系: `/api/v1/novels/:novelId/characters`
- 动作用动词表示: `/api/v1/novels/:novelId/parse`, `/api/v1/generation/scene-image`
- 支持平级和嵌套两种查询方式:
  - 嵌套: `/api/v1/novels/:novelId/characters` (强调从属关系)
  - 平级: `/api/v1/characters?novelId=:novelId` (独立查询,支持多条件过滤)

---

## URL 结构规范

- **Base URL 格式**: `http(s)://domain/api/{version}/{resource}`
- **版本控制**: 在路径中使用 `/api/v1/` 前缀
- **资源路径层级**: 最多 3 层嵌套,避免过深层级
- **路径示例**:
  - 资源集合: `/api/v1/novels`
  - 单个资源: `/api/v1/novels/{novelId}`
  - 子资源: `/api/v1/novels/{novelId}/characters`
  - 资源操作: `/api/v1/novels/{novelId}/parse`
  - 批量操作: `/api/v1/generation/batch-scenes`
- **命名约定**: 路径使用小写字母,单词用连字符分隔: `/scene-image`, `/batch-scenes`

---

## 认证授权规范

- **认证方式**: JWT (JSON Web Token)
- **Token 传递**: HTTP Header `Authorization: Bearer {token}`
- **Token 类型**:
  - Access Token: 用于 API 请求,有效期 7 天
  - Refresh Token: 用于刷新访问令牌,有效期 30 天
- **Token Payload 结构**:
  ```json
  {
    "userId": "user_abc123",
    "username": "user123",
    "iat": 1704110400,
    "exp": 1704715200
  }
  ```
- **签名算法**: HS256
- **密钥管理**: 使用环境变量 `JWT_SECRET_KEY`,最少 32 字符
- **公开端点** (无需认证):
  - `GET /health`
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/refresh`
- **保护端点**: 所有其他端点需要有效 JWT Token

---

## 请求格式规范

### Content-Type

- `application/json` (JSON 请求体)
- `multipart/form-data` (文件上传)

### 请求头要求

- `Content-Type`: 指定请求体格式
- `Authorization`: Bearer Token (除公开端点外)
- `Idempotency-Key`: 幂等性键 (可选,用于创建类操作)

### 请求体结构

- **字段命名**: 统一使用驼峰命名 (camelCase): `novelId`, `characterId`, `createdAt`
- **必填字段验证**: 在应用层验证
- **字段类型严格**: 字符串、数字、布尔值、对象、数组

### 文件上传规范

- **字段名**: 
  - `file` (单文件上传)
  - `files` (多文件上传)
- **支持格式**: 
  - `.txt` (小说文本)
  - `.jpg`, `.png` (角色参考图)
- **大小限制**: 
  - 文本文件: 50MB
  - 图片文件: 10MB
- **文件名**: 使用 UTF-8 编码,支持中文
- **附加字段**: `title`, `author` (可选)
- **安全检查**: 
  - 文件类型验证 (MIME type)
  - 文件内容扫描 (防病毒)
  - 文件名过滤 (防路径穿越)

---

## 响应格式规范

### 统一响应结构

所有响应必须包含 `code`、`message`、`data` 三个顶层字段:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 字段说明

- `code`: 业务状态码,0 表示成功,非 0 表示失败 (详见业务状态码定义)
- `message`: 操作结果描述信息
- `data`: 响应数据载体,成功时包含业务数据,失败时可为 `null`

### 字段命名规范

- **统一使用驼峰命名 (camelCase)**: `novelId`, `characterId`, `createdAt`
- **布尔字段**: 使用 `is`, `has`, `can` 前缀: `isActive`, `hasPermission`
- **时间字段**: 使用 `At` 后缀: `createdAt`, `updatedAt`
- **ID 字段**: 使用 `Id` 后缀: `novelId`, `userId`
- **数量字段**: 使用 `Count` 后缀: `chaptersCount`
- **URL 字段**: 使用 `Url` 后缀: `imageUrl`, `downloadUrl`

### 成功响应示例

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "novelId": "novel_abc123",
    "title": "修仙传",
    "status": "uploaded",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

### 列表响应结构

分页信息包含在 `data` 字段内:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 100,
      "totalPages": 5,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

### 异步任务响应结构

```json
{
  "code": 0,
  "message": "任务已创建",
  "data": {
    "taskId": "task_abc123",
    "status": "pending",
    "estimatedTime": 300
  }
}
```

### 失败响应结构

`code` 非 0,`data` 可包含错误详情或为 `null`:

```json
{
  "code": 10002,
  "message": "小说 ID 不存在",
  "data": {
    "errorDetail": "novelId: novel_invalid not found"
  }
}
```

### 数据类型规范

- **时间格式**: ISO 8601 `2024-01-01T12:00:00Z` (UTC 时间)
  - **服务端**: 统一使用 UTC 时间,后缀 `Z`
  - **前端**: 根据用户时区转换显示
  - **数据库**: 存储 UTC 时间
- **布尔值**: `true` / `false`
- **空值处理**: 使用 `null`,避免空字符串

---

## 业务状态码定义

业务状态码 (`code`) 用于表示业务逻辑层面的结果,与 HTTP 状态码分离。

| Code  | 说明                     | 场景                                    |
|-------|------------------------|-----------------------------------------|
| 0     | 成功                    | 操作成功                                |
| 10001 | 参数错误                | 必填参数缺失、格式错误、类型不匹配        |
| 10002 | 资源不存在              | Novel/Character/Scene 不存在            |
| 10003 | 资源已存在              | 重复创建 (如重复上传同一小说)            |
| 10004 | 资源状态不正确          | 操作不符合当前资源状态                   |
| 20001 | 认证失败                | Token 无效或过期                        |
| 20002 | 权限不足                | 无权限操作资源                          |
| 20003 | 用户名或密码错误        | 登录失败                                |
| 30001 | 文件上传失败            | 文件格式错误、大小超限                   |
| 30002 | 文件解析失败            | 小说文件格式不正确或内容无法解析         |
| 40001 | AI 服务调用失败         | Gemini/Sora API 错误                    |
| 40002 | AI 服务不可用           | AI 服务暂时不可用                       |
| 40003 | 生成任务失败            | 图像/视频生成失败                       |
| 50001 | 数据库错误              | 数据库操作失败                          |
| 50002 | 系统内部错误            | 未知错误                                |
| 50003 | 第三方服务错误          | 外部服务调用失败                        |

---

## HTTP 状态码规范

HTTP 状态码表示 HTTP 协议层面的结果,与业务状态码分离使用。

### 职责分工

| 层级 | 作用 | 何时使用 |
|------|------|---------|
| **HTTP 状态码** | 表示 HTTP 协议层面的结果 | 请求是否到达、路由是否匹配、服务是否可用 |
| **业务状态码** | 表示业务逻辑层面的结果 | 业务操作是否成功、具体失败原因 |

### 推荐实践

- **业务逻辑错误**: 统一返回 `HTTP 200`,通过 `code` 区分
- **协议层错误**: 使用对应的 HTTP 状态码

### HTTP 状态码列表

- `200 OK`: 请求成功 (GET, PUT, DELETE 成功,以及业务逻辑错误)
- `201 Created`: 资源创建成功 (POST 成功)
- `204 No Content`: 请求成功但无返回内容
- `400 Bad Request`: 请求格式错误 (JSON 解析失败等)
- `401 Unauthorized`: 未认证,缺少或无效 Token
- `403 Forbidden`: 已认证但无权限
- `404 Not Found`: 路由不存在
- `409 Conflict`: 资源冲突 (可选,也可用 HTTP 200 + code 10003)
- `422 Unprocessable Entity`: 请求格式正确但语义错误 (可选)
- `429 Too Many Requests`: 请求频率超限
- `500 Internal Server Error`: 服务器崩溃
- `503 Service Unavailable`: 服务不可用

### 使用示例

#### 业务错误 - HTTP 200, code 非 0

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 10002,
  "message": "小说不存在",
  "data": null
}
```

#### 协议层错误 - HTTP 401

```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "code": 20001,
  "message": "Token 无效或已过期",
  "data": null
}
```

---

## 分页规范

### 查询参数

- `page`: 页码,从 1 开始 (默认: 1)
- `pageSize`: 每页数量 (默认: 20,最大: 100)

**可选的 RESTful 风格**:
- `limit`: 每页数量
- `offset`: 偏移量

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 100,
      "totalPages": 5,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

### 边界处理

- 页码超出范围返回空数组
- `pageSize` 超过最大值时自动限制为 100

---

## 幂等性设计

### HTTP 方法幂等性

- **GET**: 幂等,多次调用结果相同
- **PUT**: 幂等,多次更新结果相同
- **DELETE**: 幂等,多次删除结果相同
- **POST**: 非幂等,可能创建多个资源

### 幂等性 Token

对于创建类 POST 请求,支持使用 `Idempotency-Key` 请求头确保幂等性。

- **幂等性 Token**: 客户端在创建类请求中携带 `Idempotency-Key` 请求头
- **服务端**: 在一定时间窗口内 (24 小时),相同 Key 的请求返回相同结果
- **适用场景**: 小说上传、场景生成等耗时操作

**示例**:

```bash
POST /api/v1/novels
Headers:
  Authorization: Bearer {token}
  Idempotency-Key: uuid-client-generated
  Content-Type: application/json

Body:
{
  "title": "修仙传",
  "author": "作者名",
  "content": "..."
}
```

**响应**: 如果重复请求,返回首次结果,避免重复创建

---

## 异步任务模式

### 适用场景

解析、生成、导出等长时间操作

### 初始请求

POST 请求返回 `taskId` 和状态:

```json
{
  "code": 0,
  "message": "任务已创建",
  "data": {
    "taskId": "task_abc123",
    "status": "pending",
    "estimatedTime": 300
  }
}
```

### 任务状态查询

#### GET /api/v1/tasks/:taskId

查询异步任务状态

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "taskId": "task_abc123",
    "status": "processing",
    "progress": 65,
    "result": null,
    "error": null,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:05:00Z"
  }
}
```

**任务状态**: `pending`, `processing`, `completed`, `failed`

### 轮询建议

- 轮询间隔: 2-5 秒
- 超时时间: 根据 `estimatedTime` 设置

### 未来扩展

支持 WebSocket 实时推送任务进度

---

## 批量操作规范

### 批量请求格式

```json
{
  "sceneIds": ["scene_1", "scene_2", "scene_3"],
  "options": {
    "parallel": true,
    "maxConcurrent": 3
  }
}
```

### 批量响应格式

```json
{
  "code": 0,
  "message": "批量任务已创建",
  "data": {
    "batchId": "batch_abc123",
    "tasks": [
      {"sceneId": "scene_1", "taskId": "task_1", "status": "pending"},
      {"sceneId": "scene_2", "taskId": "task_2", "status": "pending"},
      {"sceneId": "scene_3", "taskId": "task_3", "status": "pending"}
    ]
  }
}
```

### 并发控制

- 使用 `maxConcurrent` 限制并发数
- 部分失败不影响其他任务执行

---

## 领域资源映射

将 DDD 领域映射到 API 资源路径:

| 领域 (Domain)   | API 资源路径              | 说明                    |
|----------------|--------------------------|------------------------|
| Novel 领域      | `/api/v1/novels`         | 小说管理               |
| Character 领域  | `/api/v1/characters`     | 角色管理               |
| Scene 领域      | `/api/v1/scenes`         | 场景管理               |
| Media 领域      | `/api/v1/media`          | 媒体资源管理           |
| Generation 服务 | `/api/v1/generation/*`   | 内容生成服务           |
| Auth 服务       | `/api/v1/auth/*`         | 认证授权               |
| Task 服务       | `/api/v1/tasks/*`        | 异步任务管理           |

---

## 资源标识符规范

### ID 格式

前缀 + 下划线 + 唯一标识符

### ID 示例

- Novel: `novel_abc123`
- Character: `char_001`
- Scene: `scene_001`
- Media: `media_001`
- User: `user_abc123`
- Task: `task_abc123`

### ID 生成

- UUID 或雪花算法
- 建议长度: 20-32 字符

---

## 过滤和查询规范

- 过滤参数直接作为查询参数: `?status=parsed&author=作者名`
- 多值过滤使用逗号分隔: `?characterIds=char_001,char_002`
- 布尔过滤: `?includeReferences=true`
- 时间范围过滤: `?createdAfter=2024-01-01&createdBefore=2024-12-31`
- 模糊搜索: `?search=关键词`
- 排序: `?sortBy=createdAt&order=desc`

---

## CORS 规范

- **允许的源**: 配置环境变量 `ALLOWED_ORIGINS`
- **允许的方法**: `GET, POST, PUT, DELETE, OPTIONS`
- **允许的头**: `Content-Type, Authorization, Idempotency-Key, X-Request-ID`
- **暴露的头**: `X-Total-Count, X-Page, X-Page-Size`
- **凭证支持**: `Access-Control-Allow-Credentials: true`
- **预检请求缓存**: 86400 秒 (24 小时)

---

## 速率限制规范

- **限制策略**: 按用户 ID 或 IP 限制
- **限制级别**:
  - 全局限制: 1000 请求/小时
  - 生成接口: 100 请求/小时
  - 登录接口: 10 请求/15 分钟
- **响应头**:
  - `X-RateLimit-Limit`: 限制总数
  - `X-RateLimit-Remaining`: 剩余次数
  - `X-RateLimit-Reset`: 重置时间戳
- **超限响应**: HTTP 429 + `Retry-After` 头

---

## 缓存策略

- **缓存头**:
  - `Cache-Control: no-cache` (动态内容)
  - `Cache-Control: public, max-age=3600` (静态资源)
- **ETag 支持**: 资源版本控制
- **条件请求**:
  - `If-None-Match`: ETag 验证
  - `If-Modified-Since`: 时间验证
- **304 Not Modified**: 资源未变化时返回

---

## API 版本控制

- **版本方式**: URL 路径版本 `/api/v1/`, `/api/v2/`
- **版本兼容**: 主版本不兼容,次版本向后兼容
- **废弃策略**:
  - 提前 6 个月通知
  - 响应头添加 `X-API-Deprecated: true`
  - 文档标注 `[Deprecated]`
- **版本生命周期**: 至少支持 2 个主版本

---

## 枚举值规范

- **小说状态**: `uploaded`, `parsing`, `parsed`, `failed`
- **媒体状态**: `pending`, `processing`, `completed`, `failed`
- **角色类型**: `protagonist`, `antagonist`, `supporting`
- **媒体类型**: `image`, `video`, `audio`
- **任务状态**: `pending`, `processing`, `completed`, `failed`
- **情感类型**: `neutral`, `happy`, `sad`, `angry`, `surprised`
- **命名约定**: 使用小写,单词用下划线分隔

---

## 安全规范

### 密码要求

- 最少 8 个字符
- 包含大小写字母和数字
- 使用 bcrypt 加密,cost=10

### Token 安全

- 使用 HTTPS 传输
- 不在 URL 中传递 Token
- 设置合理过期时间

### 输入验证

- 验证所有用户输入
- 防止 SQL 注入
- 防止 XSS 攻击

### 敏感信息

- 不返回密码哈希
- 不暴露内部错误详情
- 日志中脱敏敏感数据

---

## 日志和监控

### 请求日志包含

- 请求 ID: `X-Request-ID` 头
- 时间戳
- 方法和路径
- 状态码
- 响应时间
- 用户 ID (如已认证)

### 错误日志包含

- 错误类型
- 错误消息
- 堆栈跟踪
- 请求上下文

### 性能监控

- API 响应时间
- 数据库查询时间
- 外部 API 调用时间

---

## 文档规范

### 每个端点包含

- 端点描述
- HTTP 方法
- 完整路径
- 认证要求
- 请求参数 (路径/查询/请求体)
- 响应示例 (成功和失败)
- 错误码说明
- 业务逻辑说明

### 推荐格式

- OpenAPI/Swagger 3.0 规范
- 提供交互式 API 文档
- 包含完整的 curl 示例

---

## 测试规范

- **单元测试**: 覆盖业务逻辑
- **集成测试**: 覆盖 API 端点
- **测试用例包含**:
  - 正常情况
  - 边界情况
  - 错误情况
  - 认证失败
  - 权限不足
- **使用测试数据库**
- **Mock 外部服务** (AI API)

---

## 性能要求

- **API 响应时间**:
  - 查询接口: <200ms
  - 创建/更新: <500ms
  - 异步任务: <100ms (返回 taskId)
- **并发支持**: 100 并发请求
- **数据库连接池**: 最大 25 连接
- **超时设置**:
  - HTTP 请求: 30 秒
  - 数据库查询: 10 秒
  - AI API 调用: 60 秒

---

## DDD 架构依赖原则

- **接口层** 依赖 **应用层**
- **应用层** 依赖 **领域层**
- **基础设施层** 实现 **领域层接口**
- **领域层** 零外部依赖
- **依赖注入** 管理层间依赖关系

---

## 角色一致性规范

- **参考图生成**: POST `/api/v1/characters/:characterId/references`
- **参考图用途**: 确保角色在所有场景中视觉一致
- **一致性参数**: `consistencyStrength` (0-1,默认 0.8)
- **实现方式**: 使用参考图进行图生图转换 (Gemini ImageToImage API)
- **状态管理**: 支持多种状态 (default, battle, formal, custom)

详见: [CHARACTER_CONSISTENCY.md](CHARACTER_CONSISTENCY.md)

---

## 总结

本规范文档涵盖了 AI-Motion 项目 API 设计的所有核心方面,从 RESTful 原则到 DDD 架构实践,从安全认证到性能优化。

### 核心要点

1. **统一字段命名**: 请求和响应统一使用驼峰命名 (camelCase)
2. **双层状态码**: HTTP 状态码处理协议层,业务状态码处理业务逻辑层
3. **完整的业务状态码**: 定义了 0-50002 的完整业务错误码体系
4. **幂等性保证**: 支持 `Idempotency-Key` 确保创建操作幂等
5. **异步任务支持**: 完整的异步任务创建、查询、进度跟踪机制
6. **领域驱动设计**: API 路径与 DDD 领域清晰映射
7. **安全第一**: JWT 认证、输入验证、速率限制、HTTPS 传输

遵循这些规范可以确保 API 的一致性、可维护性和可扩展性。

---

**文档版本**: v1.0  
**最后更新**: 2024-01-01  
**维护者**: AI-Motion 开发团队
