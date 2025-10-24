# AI-Motion API 设计规范

本文档总结 AI-Motion 项目的 API 设计规范,以结构化格式呈现供 AI 模型理解。

---

## RESTful 设计原则

- 使用标准 HTTP 方法: GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)
- 资源命名使用复数名词: `/api/v1/novels`, `/api/v1/characters`, `/api/v1/scenes`
- 路径参数用于资源标识: `/api/v1/novels/:novel_id`, `/api/v1/characters/:character_id`
- 查询参数用于过滤和分页: `?page=1&page_size=20&status=parsed`
- 嵌套资源表示从属关系: `/api/v1/novels/:novel_id/characters`
- 动作用动词表示: `/api/v1/novels/:novel_id/parse`, `/api/v1/generation/scene-image`

---

## URL 结构规范

- **Base URL 格式**: `http(s)://domain/api/{version}/{resource}`
- **版本控制**: 在路径中使用 `/api/v1/` 前缀
- **资源路径层级**: 最多 3 层嵌套,避免过深层级
- **路径示例**:
  - 资源集合: `/api/v1/novels`
  - 单个资源: `/api/v1/novels/{novel_id}`
  - 子资源: `/api/v1/novels/{novel_id}/characters`
  - 资源操作: `/api/v1/novels/{novel_id}/parse`
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
    "user_id": "user_abc123",
    "username": "user123",
    "iat": 1704110400,
    "exp": 1704715200
  }
  ```
- **签名算法**: HS256
- **密钥管理**: 使用环境变量 `JWT_SECRET_KEY`,最少 32 字符
- **公开端点**(无需认证):
  - `GET /health`
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
- **保护端点**: 所有其他端点需要有效 JWT Token

---

## 请求格式规范

- **Content-Type**:
  - `application/json` (JSON 请求体)
  - `multipart/form-data` (文件上传)
- **请求头要求**:
  - `Content-Type`: 指定请求体格式
  - `Authorization`: Bearer Token (除公开端点外)
- **请求体结构**:
  - 使用驼峰命名: `novelId`, `characterId`, `createdAt`
  - 必填字段验证: 在应用层验证
  - 字段类型严格: 字符串、数字、布尔值、对象、数组
- **文件上传格式**:
  - 字段名: `file`
  - 支持格式: `.txt` (小说文件)
  - 大小限制: 50MB
  - 附加字段: `title`, `author` (可选)

---

## 响应格式规范

- **Content-Type**: `application/json`
- **统一响应结构**: 所有响应必须包含 `code`、`message`、`data` 三个顶层字段
  ```json
  {
    "code": 0,
    "message": "success",
    "data": { }
  }
  ```
- **字段说明**:
  - `code`: 业务状态码,0 表示成功,非 0 表示失败
  - `message`: 操作结果描述信息
  - `data`: 响应数据载体,成功时包含业务数据,失败时可为 `null`

- **成功响应示例**:
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "novel_id": "novel_abc123",
      "title": "修仙传",
      "status": "uploaded"
    }
  }
  ```

- **列表响应结构**: 分页信息包含在 `data` 字段内
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "items": [ ],
      "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 100,
        "total_pages": 5
      }
    }
  }
  ```

- **异步任务响应结构**:
  ```json
  {
    "code": 0,
    "message": "任务已创建",
    "data": {
      "task_id": "task_abc123",
      "status": "pending",
      "estimated_time": 300
    }
  }
  ```

- **失败响应结构**: `code` 非 0,`data` 可包含错误详情或为 `null`
  ```json
  {
    "code": 40001,
    "message": "小说 ID 不存在",
    "data": {
      "error_detail": "novel_id: novel_invalid not found"
    }
  }
  ```

- **字段命名**: 使用 snake_case (与请求体区分)
- **时间格式**: ISO 8601 `2024-01-01T12:00:00Z`
- **布尔值**: `true` / `false`
- **空值处理**: 使用 `null` 或省略字段

---

## HTTP 状态码规范

- `200 OK`: 请求成功 (GET, PUT, DELETE 成功)
- `201 Created`: 资源创建成功 (POST 成功)
- `204 No Content`: 请求成功但无返回内容
- `400 Bad Request`: 请求参数错误,格式不正确
- `401 Unauthorized`: 未认证,缺少或无效 Token
- `403 Forbidden`: 无权限访问资源
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突 (如重复创建)
- `422 Unprocessable Entity`: 业务逻辑错误
- `429 Too Many Requests`: 请求频率超限
- `500 Internal Server Error`: 服务器内部错误
- `503 Service Unavailable`: 外部服务不可用 (如 AI API)

---

## 错误响应规范

- **统一错误响应格式**: 遵循 `code`、`message`、`data` 结构
  ```json
  {
    "code": 40001,
    "message": "请求参数无效",
    "data": {
      "field": "novel_id",
      "value": "invalid_id",
      "timestamp": "2024-01-01T12:00:00Z"
    }
  }
  ```

- **业务错误码体系**: 5 位数字,按模块划分
  - **0**: 成功
  - **10xxx**: 认证授权错误
  - **20xxx**: 小说管理错误
  - **30xxx**: 角色管理错误
  - **40xxx**: 场景管理错误
  - **50xxx**: 内容生成错误
  - **60xxx**: 导出功能错误
  - **90xxx**: 系统级错误

- **常见错误码定义**:

  **认证授权 (10xxx)**:
  - `10001`: Token 无效
  - `10002`: Token 已过期
  - `10003`: 用户名或密码错误
  - `10004`: 用户未登录
  - `10005`: 权限不足
  - `10006`: 用户名已存在
  - `10007`: 邮箱已存在
  - `10008`: 密码强度不足

  **小说管理 (20xxx)**:
  - `20001`: 小说 ID 不存在
  - `20002`: 文件格式不支持
  - `20003`: 文件过大
  - `20004`: 小说解析失败
  - `20005`: 小说已存在

  **角色管理 (30xxx)**:
  - `30001`: 角色 ID 不存在
  - `30002`: 参考图生成失败
  - `30003`: 角色描述不完整

  **场景管理 (40xxx)**:
  - `40001`: 场景 ID 不存在
  - `40002`: 场景数据无效

  **内容生成 (50xxx)**:
  - `50001`: 图片生成失败
  - `50002`: 视频生成失败
  - `50003`: 配音生成失败
  - `50004`: AI 服务不可用
  - `50005`: 任务 ID 不存在
  - `50006`: 任务已取消

  **系统级错误 (90xxx)**:
  - `90001`: 请求参数无效
  - `90002`: 资源冲突
  - `90003`: 请求频率超限
  - `90004`: 服务器内部错误
  - `90005`: 数据库错误
  - `90006`: 外部服务调用失败

- **错误响应示例**:
  ```json
  {
    "code": 10003,
    "message": "用户名或密码错误",
    "data": null
  }
  ```

  ```json
  {
    "code": 20002,
    "message": "文件格式不支持",
    "data": {
      "allowed_formats": [".txt"],
      "received_format": ".pdf"
    }
  }
  ```

- **错误消息语言**: 中文
- **敏感信息保护**: 不暴露内部实现细节、堆栈信息和数据库结构

---

## 分页规范

- **查询参数**:
  - `page`: 页码,从 1 开始
  - `page_size`: 每页数量,默认 20,最大 100
- **响应格式**: 分页信息包含在 `data` 对象内
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "items": [
        {
          "novel_id": "novel_001",
          "title": "修仙传"
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 100,
        "total_pages": 5
      }
    }
  }
  ```
- **默认值**: `page=1`, `page_size=20`
- **边界处理**: 页码超出范围返回空 `items` 数组,`total` 保持实际值

---

## 过滤和查询规范

- 过滤参数直接作为查询参数: `?status=parsed&author=作者名`
- 多值过滤使用逗号分隔: `?character_ids=char_001,char_002`
- 布尔过滤: `?include_references=true`
- 时间范围过滤: `?created_after=2024-01-01&created_before=2024-12-31`
- 模糊搜索: `?search=关键词`
- 排序: `?sort_by=created_at&order=desc`

---

## 异步任务模式

- **适用场景**: 解析、生成、导出等长时间操作
- **初始请求**: POST 返回任务信息
  ```json
  {
    "code": 0,
    "message": "任务已创建",
    "data": {
      "task_id": "task_abc123",
      "status": "pending",
      "estimated_time": 300
    }
  }
  ```
- **状态查询**: GET `/api/v1/generation/tasks/{task_id}`
- **任务状态流转**: `pending` → `processing` → `completed` / `failed`
- **轮询建议**: 2-5 秒间隔
- **进度跟踪响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "task_id": "task_abc123",
      "status": "processing",
      "progress": {
        "total": 100,
        "completed": 45,
        "current": "scene_046"
      }
    }
  }
  ```
- **任务完成响应**:
  ```json
  {
    "code": 0,
    "message": "任务已完成",
    "data": {
      "task_id": "task_abc123",
      "status": "completed",
      "result": {
        "media_ids": ["media_001", "media_002"]
      }
    }
  }
  ```
- **任务失败响应**:
  ```json
  {
    "code": 50001,
    "message": "图片生成失败",
    "data": {
      "task_id": "task_abc123",
      "status": "failed",
      "error": "AI service timeout"
    }
  }
  ```
- **未来扩展**: 支持 WebSocket 实时推送

---

## 批量操作规范

- **批量请求格式**:
  ```json
  {
    "items": ["id1", "id2", "id3"],
    "options": {
      "parallel": true,
      "max_concurrent": 3
    }
  }
  ```
- **批量响应格式**:
  ```json
  {
    "code": 0,
    "message": "批量操作已提交",
    "data": {
      "batch_id": "batch_001",
      "total": 3,
      "status": "processing",
      "results": [
        {"id": "id1", "status": "success"},
        {"id": "id2", "status": "failed", "error": "资源不存在"},
        {"id": "id3", "status": "pending"}
      ]
    }
  }
  ```
- **部分成功响应**: 当部分项失败时,`code` 仍为 0,但在 `data.results` 中标记每项状态
  ```json
  {
    "code": 0,
    "message": "批量操作完成,部分失败",
    "data": {
      "batch_id": "batch_001",
      "total": 3,
      "success_count": 2,
      "failed_count": 1,
      "results": [
        {"id": "id1", "status": "success"},
        {"id": "id2", "status": "failed", "error": "资源不存在"},
        {"id": "id3", "status": "success"}
      ]
    }
  }
  ```
- **并发控制**: 使用 `max_concurrent` 限制并发数
- **部分失败处理**: 返回每个项的独立状态,不中断整体操作

---

## 文件上传规范

- **编码方式**: `multipart/form-data`
- **文件字段名**: `file`
- **支持格式**: `.txt` (小说), `.jpg/.png` (参考图)
- **大小限制**: 50MB
- **验证要求**:
  - 文件格式检查
  - 文件大小检查
  - MIME 类型验证
- **响应包含**:
  - `file_size`: 文件大小(字节)
  - `filename`: 原始文件名
  - `url`: 访问 URL (如适用)

---

## 资源标识符规范

- **ID 格式**: 前缀 + 下划线 + 唯一标识符
- **ID 示例**:
  - Novel: `novel_abc123`
  - Character: `char_001`
  - Scene: `scene_001`
  - Media: `media_001`
  - User: `user_abc123`
  - Task: `task_abc123`
- **ID 生成**: UUID 或雪花算法
- **ID 长度**: 建议 20-32 字符
- **URL 编码**: 对特殊字符进行编码

---

## CORS 规范

- **允许的源**: 配置环境变量 `ALLOWED_ORIGINS`
- **允许的方法**: `GET, POST, PUT, DELETE, OPTIONS`
- **允许的头**: `Content-Type, Authorization, X-Request-ID`
- **暴露的头**: `X-Total-Count, X-Page, X-Page-Size`
- **凭证支持**: `Access-Control-Allow-Credentials: true`
- **预检请求缓存**: 86400 秒(24 小时)

---

## 请求幂等性规范

- **GET**: 幂等,多次调用结果相同
- **PUT**: 幂等,多次更新结果相同
- **DELETE**: 幂等,多次删除结果相同
- **POST**: 非幂等,可能创建多个资源
- **幂等键支持**: 使用 `Idempotency-Key` 头 (可选)
- **幂等键有效期**: 24 小时

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

## 字段命名规范

- **请求参数**: camelCase (前端友好)
- **响应字段**: snake_case (后端惯例)
- **布尔字段**: 使用 `is_`, `has_`, `can_` 前缀
- **时间字段**: 使用 `_at` 后缀 (`created_at`, `updated_at`)
- **ID 字段**: 使用 `_id` 后缀 (`novel_id`, `user_id`)
- **数量字段**: 使用 `_count` 后缀 (`chapters_count`)
- **URL 字段**: 使用 `_url` 后缀 (`image_url`, `download_url`)

---

## 数据类型规范

- **字符串**: 用于文本、ID、枚举值
- **整数**: 用于计数、序号
- **浮点数**: 用于百分比、评分
- **布尔值**: `true` / `false`
- **日期时间**: ISO 8601 格式字符串
- **数组**: 用于列表数据
- **对象**: 用于嵌套结构
- **枚举值**: 字符串形式,小写下划线分隔
- **空值**: 使用 `null`,避免空字符串表示空值

---

## 枚举值规范

- **小说状态**: `uploaded`, `parsing`, `parsed`, `failed`
- **媒体状态**: `pending`, `processing`, `completed`, `failed`
- **角色类型**: `protagonist`, `antagonist`, `supporting`
- **媒体类型**: `image`, `video`, `audio`
- **任务状态**: `pending`, `processing`, `completed`, `failed`
- **情感类型**: `neutral`, `happy`, `sad`, `angry`, `surprised`
- **命名约定**: 一致性使用小写,单词用下划线分隔

---

## 安全规范

- **密码要求**:
  - 最少 8 个字符
  - 包含大小写字母和数字
  - 使用 bcrypt 加密,cost=10
- **Token 安全**:
  - 使用 HTTPS 传输
  - 不在 URL 中传递 Token
  - 设置合理过期时间
- **输入验证**:
  - 验证所有用户输入
  - 防止 SQL 注入
  - 防止 XSS 攻击
- **敏感信息**:
  - 不返回密码哈希
  - 不暴露内部错误详情
  - 日志中脱敏敏感数据

---

## 日志和监控

- **请求日志包含**:
  - 请求 ID: `X-Request-ID` 头
  - 时间戳
  - 方法和路径
  - 状态码
  - 响应时间
  - 用户 ID (如已认证)
- **错误日志包含**:
  - 错误类型
  - 错误消息
  - 堆栈跟踪
  - 请求上下文
- **性能监控**:
  - API 响应时间
  - 数据库查询时间
  - 外部 API 调用时间

---

## 文档规范

- **每个端点包含**:
  - 端点描述
  - HTTP 方法
  - 完整路径
  - 认证要求
  - 请求参数 (路径/查询/请求体)
  - 响应示例 (成功和失败)
  - 错误码说明
  - 业务逻辑说明
- **推荐格式**: OpenAPI/Swagger 3.0 规范
- **提供交互式 API 文档**
- **包含完整的 curl 示例**

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
  - 异步任务: <100ms (返回 task_id)
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

## 业务领域端点分类

- **认证管理**: `/api/v1/auth/*`
- **小说管理**: `/api/v1/novels/*`
- **角色管理**: `/api/v1/characters/*`
- **场景管理**: `/api/v1/scenes/*`
- **内容生成**: `/api/v1/generation/*`
- **项目管理**: `/api/v1/projects/*`
- **导出功能**: `/api/v1/export/*`

---

## 角色一致性规范

- **参考图生成**: POST `/api/v1/characters/:character_id/references`
- **参考图用途**: 确保角色在所有场景中视觉一致
- **一致性参数**: `consistency_strength` (0-1,默认 0.8)
- **实现方式**: 使用参考图进行图生图转换
- **状态管理**: 支持多种状态 (default, battle, formal, custom)

---

## 总结

本规范文档涵盖了 AI-Motion 项目 API 设计的所有核心方面,从 RESTful 原则到 DDD 架构实践,从安全认证到性能优化。遵循这些规范可以确保 API 的一致性、可维护性和可扩展性。
