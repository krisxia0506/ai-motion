# 简化版漫画生成业务流程设计

## 文档概述

本文档详细描述了简化版漫画生成功能的业务流程、后端接口设计和前端交互方案。

**版本**: v1.0
**日期**: 2025-01-26
**状态**: 设计文档

---

## 1. 业务流程概述

### 1.1 简化前后对比

**简化前** (复杂流程):
1. 用户上传小说
2. 用户手动提取角色
3. 用户手动划分场景
4. 用户为每个场景生成图片
5. 用户导出结果

**简化后** (一键流程):
1. 用户输入小说内容 (文件上传或文本输入)
2. 点击"开始生成"按钮
3. 系统返回任务ID，跳转到任务详情页
4. 前端轮询任务状态，实时显示进度
5. 生成完成后查看结果

### 1.2 核心业务场景

```
┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  用户输入   │ ───> │  创建异步任务    │ ───> │  返回任务ID     │
│  小说内容   │      │  POST /generate  │      │  task_id        │
└─────────────┘      └──────────────────┘      └─────────────────┘
                                                         │
                                                         ▼
┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  显示结果   │ <─── │  轮询任务状态    │ <─── │  跳转详情页     │
│  或重试     │      │  GET /status     │      │  /task/:id      │
└─────────────┘      └──────────────────┘      └─────────────────┘
```

---

## 2. 用户认证机制

### 2.1 认证方案: Supabase Auth

本系统使用 **Supabase Authentication** 进行用户认证和授权管理。

**核心特性**:
- JWT (JSON Web Token) 基于 Token 的认证
- 支持邮箱/密码登录
- 自动刷新 Token (Refresh Token)
- 用户会话管理

### 2.2 认证流程

```
┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  用户登录   │ ───> │  Supabase Auth   │ ───> │  返回 JWT Token │
│  邮箱/密码  │      │  验证凭证        │      │  + Refresh Token│
└─────────────┘      └──────────────────┘      └─────────────────┘
                                                         │
                                                         ▼
┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐
│  访问API    │ <─── │  存储到本地      │ <─── │  前端保存Token  │
│  带Token    │      │  localStorage    │      │                 │
└─────────────┘      └──────────────────┘      └─────────────────┘
```

### 2.3 JWT Token 结构

Supabase 返回的 JWT Token 包含以下信息:

```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",  // 用户ID
  "email": "user@example.com",
  "role": "authenticated",
  "iat": 1706270400,
  "exp": 1706274000
}
```

**字段说明**:
- `sub`: 用户唯一标识 (User ID)
- `email`: 用户邮箱
- `role`: 用户角色 (authenticated/anon)
- `iat`: Token 签发时间
- `exp`: Token 过期时间

### 2.4 后端认证中间件

所有 API 接口需要验证 JWT Token:

**验证流程**:
1. 从请求头 `Authorization: Bearer <token>` 中提取 Token
2. 使用 Supabase JWT Secret 验证 Token 签名
3. 检查 Token 是否过期
4. 提取 User ID 并注入到请求上下文
5. 继续执行业务逻辑

**认证失败响应**:
```json
{
  "code": 20001,
  "message": "未授权: Token无效或已过期",
  "data": null
}
```

**HTTP状态码**: `401 Unauthorized`

### 2.5 业务状态码扩展

| Code  | 说明 | 场景 |
|-------|------|------|
| 20001 | 未授权 | Token缺失、无效或过期 |
| 20002 | Token已过期 | 需要刷新Token |
| 20003 | 权限不足 | 访问其他用户的资源 |

### 2.6 数据隔离

**用户数据隔离原则**:
- 每个任务关联一个 `user_id` 字段
- 用户只能访问自己创建的任务
- 查询时自动添加 `user_id` 过滤条件
- 防止横向越权访问

**示例**:
```sql
-- 查询任务时自动过滤用户ID
SELECT * FROM tasks WHERE user_id = '当前用户ID' AND id = '任务ID';
```

---

## 3. 后端接口设计

### 3.1 接口列表

| 接口路径 | 方法 | 功能 | 认证 | 说明 |
|---------|------|------|------|------|
| `/api/v1/manga/generate` | POST | 创建漫画生成任务 | 需要 | 异步任务，立即返回任务ID |
| `/api/v1/manga/task/:task_id` | GET | 获取任务状态 | 需要 | 轮询接口，返回任务进度 |
| `/api/v1/manga/tasks` | GET | 获取任务列表 | 需要 | 分页查询当前用户的任务列表 |
| `/api/v1/manga/task/:task_id/cancel` | POST | 取消任务 | 需要 | 可选功能，取消正在执行的任务 |

**认证说明**: 所有接口需要在请求头中携带 Supabase JWT Token

---

### 3.2 接口详细设计

#### 3.2.1 创建漫画生成任务

**接口**: `POST /api/v1/manga/generate`

**认证**: 需要 JWT Token

**功能**: 接收用户提交的小说内容，创建异步生成任务，立即返回任务ID

**请求头**:
```
Authorization: Bearer <supabase_jwt_token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "title": "小红帽",
  "author": "格林兄弟",
  "content": "从前有个可爱的小姑娘，谁见了都喜欢..."
}
```

**参数说明**:
- `title` (required, string): 小说标题，最大长度200字符
- `author` (optional, string): 作者名称，默认为"Unknown"
- `content` (required, string): 小说内容，100-5000字

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/v1/manga/generate \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "title": "小红帽",
    "author": "格林兄弟",
    "content": "从前有个可爱的小姑娘，谁见了都喜欢..."
  }'
```

**响应示例** (成功):
```json
{
  "code": 0,
  "message": "任务已创建",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "created_at": "2024-01-26T12:00:00Z"
  }
}
```

**响应示例** (失败):
```json
{
  "code": 10001,
  "message": "小说内容不能少于100字",
  "data": null
}
```

**业务逻辑**:
1. 验证请求参数 (标题、内容长度)
2. 创建Novel实体 (status: "processing")
3. 创建Task实体保存到数据库
4. 启动异步Goroutine执行以下流程:
   - 解析小说 (章节划分)
   - 提取角色信息
   - 为每个角色生成参考图
   - 划分场景
   - 匹配场景与角色
   - 为每个场景生成图片
   - 更新任务状态为"completed"或"failed"
5. 立即返回任务ID给前端

**状态码**:
- `0`: 成功
- `10001`: 参数错误 (标题为空、内容长度不符)
- `50001`: 数据库错误
- `50002`: 系统内部错误

---

#### 3.2.2 获取任务状态

**接口**: `GET /api/v1/manga/task/:task_id`

**认证**: 需要 JWT Token

**功能**: 查询任务的执行状态、进度和结果

**请求头**:
```
Authorization: Bearer <supabase_jwt_token>
```

**路径参数**:
- `task_id` (required): 任务ID

**请求示例**:
```bash
curl http://localhost:8080/api/v1/manga/task/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例** (任务进行中):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing",
    "progress": {
      "current_step": "生成角色参考图",
      "current_step_index": 3,
      "total_steps": 6,
      "percentage": 50,
      "details": {
        "characters_extracted": 3,
        "characters_generated": 2,
        "scenes_divided": 8,
        "scenes_generated": 0
      }
    },
    "created_at": "2024-01-26T12:00:00Z",
    "updated_at": "2024-01-26T12:05:32Z"
  }
}
```

**响应示例** (任务完成):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "completed",
    "progress": {
      "current_step": "完成",
      "current_step_index": 6,
      "total_steps": 6,
      "percentage": 100,
      "details": {
        "characters_extracted": 3,
        "characters_generated": 3,
        "scenes_divided": 8,
        "scenes_generated": 8
      }
    },
    "result": {
      "novel_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "小红帽",
      "character_count": 3,
      "scene_count": 8,
      "characters": [
        {
          "id": "char_001",
          "name": "小红帽",
          "reference_image_url": "https://storage.example.com/char_001.jpg"
        },
        {
          "id": "char_002",
          "name": "大灰狼",
          "reference_image_url": "https://storage.example.com/char_002.jpg"
        },
        {
          "id": "char_003",
          "name": "奶奶",
          "reference_image_url": "https://storage.example.com/char_003.jpg"
        }
      ],
      "scenes": [
        {
          "id": "scene_001",
          "sequence_num": 1,
          "description": "森林小路",
          "image_url": "https://storage.example.com/scene_001.jpg"
        },
        {
          "id": "scene_002",
          "sequence_num": 2,
          "description": "奶奶的房子",
          "image_url": "https://storage.example.com/scene_002.jpg"
        }
      ]
    },
    "created_at": "2024-01-26T12:00:00Z",
    "updated_at": "2024-01-26T12:10:45Z",
    "completed_at": "2024-01-26T12:10:45Z"
  }
}
```

**响应示例** (任务失败):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "failed",
    "progress": {
      "current_step": "生成场景图片",
      "current_step_index": 5,
      "total_steps": 6,
      "percentage": 83
    },
    "error": {
      "code": 40001,
      "message": "AI服务调用失败: Gemini API rate limit exceeded",
      "retry_able": true
    },
    "created_at": "2024-01-26T12:00:00Z",
    "updated_at": "2024-01-26T12:08:23Z",
    "failed_at": "2024-01-26T12:08:23Z"
  }
}
```

**响应示例** (任务不存在):
```json
{
  "code": 10002,
  "message": "任务不存在",
  "data": null
}
```

**任务状态说明**:

| 状态 | 说明 | 前端行为 |
|------|------|---------|
| `pending` | 任务已创建，等待执行 | 显示"等待中"，继续轮询 |
| `processing` | 任务执行中 | 显示进度条，继续轮询 |
| `completed` | 任务成功完成 | 停止轮询，显示结果 |
| `failed` | 任务失败 | 停止轮询，显示错误信息 |
| `cancelled` | 任务已取消 | 停止轮询，显示取消提示 |

**进度步骤说明**:

| 步骤索引 | 步骤名称 | 说明 |
|---------|---------|------|
| 1 | 解析小说 | 解析章节结构 |
| 2 | 提取角色 | 识别角色和外貌 |
| 3 | 生成角色参考图 | 为每个角色生成参考图 |
| 4 | 划分场景 | 将章节划分为场景 |
| 5 | 生成场景图片 | 为每个场景生成图片 |
| 6 | 完成 | 所有步骤完成 |

**状态码**:
- `0`: 成功
- `10002`: 任务不存在

**轮询建议**:
- 轮询间隔: 2秒
- 超时时间: 15分钟 (停止轮询并提示用户)
- 错误重试: 3次

---

#### 3.2.3 获取任务列表

**接口**: `GET /api/v1/manga/tasks`

**认证**: 需要 JWT Token

**功能**: 分页查询当前用户的所有任务列表，支持按状态筛选

**请求头**:
```
Authorization: Bearer <supabase_jwt_token>
```

**查询参数**:
- `page` (optional, number, default: 1): 页码，从1开始
- `page_size` (optional, number, default: 20): 每页数量，最大100
- `status` (optional, string): 按状态筛选，可选值: `pending`, `processing`, `completed`, `failed`, `cancelled`
- `sort` (optional, string, default: `created_at_desc`): 排序方式，可选值: `created_at_desc`, `created_at_asc`, `updated_at_desc`

**请求示例**:
```bash
# 获取所有任务
curl http://localhost:8080/api/v1/manga/tasks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 只查看已完成的任务
curl "http://localhost:8080/api/v1/manga/tasks?status=completed&page=1&page_size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例** (成功):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "task_id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "小红帽",
        "status": "completed",
        "progress": {
          "percentage": 100,
          "current_step": "完成"
        },
        "character_count": 3,
        "scene_count": 8,
        "created_at": "2024-01-26T12:00:00Z",
        "updated_at": "2024-01-26T12:10:45Z",
        "completed_at": "2024-01-26T12:10:45Z"
      },
      {
        "task_id": "660e8400-e29b-41d4-a716-446655440001",
        "title": "白雪公主",
        "status": "processing",
        "progress": {
          "percentage": 50,
          "current_step": "生成角色参考图"
        },
        "created_at": "2024-01-26T13:00:00Z",
        "updated_at": "2024-01-26T13:05:32Z"
      },
      {
        "task_id": "770e8400-e29b-41d4-a716-446655440002",
        "title": "灰姑娘",
        "status": "failed",
        "progress": {
          "percentage": 83,
          "current_step": "生成场景图片"
        },
        "error": {
          "message": "AI服务调用失败"
        },
        "created_at": "2024-01-26T11:00:00Z",
        "updated_at": "2024-01-26T11:08:23Z",
        "failed_at": "2024-01-26T11:08:23Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 3,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

**响应示例** (空列表):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 0,
      "total_pages": 0,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

**状态码**:
- `0`: 成功
- `10001`: 参数错误 (page、page_size 无效)
- `20001`: 未授权 (Token无效或缺失)

**业务逻辑**:
1. 验证JWT Token，提取用户ID
2. 根据用户ID查询任务列表
3. 支持按状态筛选和分页
4. 返回任务基本信息 (不包含完整的结果数据)

**注意事项**:
- 任务列表只返回当前用户创建的任务
- 不同用户的任务数据相互隔离
- 任务按创建时间倒序排列

---

#### 3.2.4 取消任务 (可选)

**接口**: `POST /api/v1/manga/task/:task_id/cancel`

**认证**: 需要 JWT Token

**功能**: 取消正在执行的任务

**请求头**:
```
Authorization: Bearer <supabase_jwt_token>
```

**路径参数**:
- `task_id` (required): 任务ID

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/v1/manga/task/550e8400-e29b-41d4-a716-446655440000/cancel \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应示例** (成功):
```json
{
  "code": 0,
  "message": "任务已取消",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "cancelled"
  }
}
```

**响应示例** (失败 - 任务已完成):
```json
{
  "code": 10001,
  "message": "任务已完成，无法取消",
  "data": null
}
```

**业务逻辑**:
1. 检查任务是否存在
2. 检查任务状态是否为"pending"或"processing"
3. 设置取消标志位
4. 等待Goroutine检测到取消标志并退出
5. 更新任务状态为"cancelled"

**状态码**:
- `0`: 成功
- `10001`: 参数错误 (任务已完成或已取消)
- `10002`: 任务不存在

---

## 4. 前端交互流程

### 4.1 认证集成

#### 4.1.1 Supabase Client 初始化

```typescript
// src/lib/supabase.ts
import { createClient } from '@supabase/supabase-js'

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

export const supabase = createClient(supabaseUrl, supabaseAnonKey)
```

#### 4.1.2 登录/注册页面

```typescript
// src/pages/LoginPage.tsx
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { supabase } from '../lib/supabase'

const LoginPage = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const { data, error } = await supabase.auth.signInWithPassword({
        email,
        password,
      })

      if (error) throw error

      // 登录成功，跳转到首页
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSignUp = async () => {
    setLoading(true)
    setError(null)

    try {
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
      })

      if (error) throw error

      alert('注册成功！请查收邮箱验证邮件')
    } catch (err) {
      setError(err instanceof Error ? err.message : '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <form onSubmit={handleLogin}>
        <h1>登录 AI-Motion</h1>

        <input
          type="email"
          placeholder="邮箱"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />

        <input
          type="password"
          placeholder="密码"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />

        {error && <div className="error">{error}</div>}

        <button type="submit" disabled={loading}>
          {loading ? '登录中...' : '登录'}
        </button>

        <button type="button" onClick={handleSignUp} disabled={loading}>
          注册新账号
        </button>
      </form>
    </div>
  )
}

export default LoginPage
```

#### 4.1.3 认证状态管理

```typescript
// src/hooks/useAuth.ts
import { useEffect, useState } from 'react'
import { User, Session } from '@supabase/supabase-js'
import { supabase } from '../lib/supabase'

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null)
  const [session, setSession] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // 获取当前会话
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session)
      setUser(session?.user ?? null)
      setLoading(false)
    })

    // 监听认证状态变化
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
      setUser(session?.user ?? null)
    })

    return () => subscription.unsubscribe()
  }, [])

  const signOut = async () => {
    await supabase.auth.signOut()
  }

  return {
    user,
    session,
    loading,
    signOut,
  }
}
```

#### 4.1.4 路由守卫 (Protected Route)

```typescript
// src/components/ProtectedRoute.tsx
import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

interface ProtectedRouteProps {
  children: React.ReactNode
}

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { user, loading } = useAuth()

  if (loading) {
    return <div>加载中...</div>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export default ProtectedRoute
```

#### 4.1.5 API Client 认证拦截器

```typescript
// src/services/api.ts
import axios from 'axios'
import { supabase } from '../lib/supabase'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：自动添加 JWT Token
apiClient.interceptors.request.use(
  async (config) => {
    const { data: { session } } = await supabase.auth.getSession()

    if (session?.access_token) {
      config.headers.Authorization = `Bearer ${session.access_token}`
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：处理认证错误
apiClient.interceptors.response.use(
  (response) => {
    // 统一处理三段式响应
    if (response.data.code !== 0) {
      return Promise.reject(new Error(response.data.message))
    }
    return response.data
  },
  async (error) => {
    // 处理 401 认证错误
    if (error.response?.status === 401) {
      // Token 过期，尝试刷新
      const { data, error: refreshError } = await supabase.auth.refreshSession()

      if (refreshError || !data.session) {
        // 刷新失败，跳转到登录页
        await supabase.auth.signOut()
        window.location.href = '/login'
        return Promise.reject(new Error('登录已过期，请重新登录'))
      }

      // 刷新成功，重试原请求
      const originalRequest = error.config
      originalRequest.headers.Authorization = `Bearer ${data.session.access_token}`
      return apiClient(originalRequest)
    }

    return Promise.reject(error)
  }
)

export { apiClient }
```

#### 4.1.6 环境变量配置

```bash
# frontend/.env
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 4.2 页面结构

```
App (根组件)
  │
  ├─ LoginPage (登录页) ────────────> 未登录用户访问
  │
  ├─ ProtectedRoute (需要认证) ──┐
  │                              │
  │  ┌───────────────────────────┘
  │  │
  │  ├─ HomePage (首页)
  │  │   ├─ 输入表单
  │  │   │   ├─ 模式切换 (文件上传 / 文本输入)
  │  │   │   ├─ 标题输入框
  │  │   │   ├─ 作者输入框 (可选)
  │  │   │   └─ "开始生成" 按钮
  │  │   └─ 导航链接 → 任务列表
  │  │
  │  ├─ TaskListPage (任务列表页)
  │  │   ├─ 筛选器 (按状态)
  │  │   ├─ 任务卡片列表
  │  │   └─ 分页器
  │  │
  │  └─ TaskDetailPage (任务详情页)
  │      ├─ 任务信息卡片
  │      ├─ 进度条 + 当前步骤
  │      ├─ 详细进度信息
  │      ├─ 结果展示区域
  │      └─ 操作按钮 (取消/重试/查看详情)
```

### 4.3 路由配置

```typescript
// src/App.tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { useAuth } from './hooks/useAuth'
import ProtectedRoute from './components/ProtectedRoute'
import LoginPage from './pages/LoginPage'
import HomePage from './pages/HomePage'
import TaskListPage from './pages/TaskListPage'
import TaskDetailPage from './pages/TaskDetailPage'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route
          path="/"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/tasks"
          element={
            <ProtectedRoute>
              <TaskListPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/task/:taskId"
          element={
            <ProtectedRoute>
              <TaskDetailPage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
```

### 4.4 任务列表页面实现

```typescript
// src/pages/TaskListPage.tsx
import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { apiClient } from '../services/api'
import { TaskStatus } from '../types/task'

interface TaskItem {
  task_id: string
  title: string
  status: TaskStatus
  progress: {
    percentage: number
    current_step: string
  }
  character_count?: number
  scene_count?: number
  created_at: string
  updated_at: string
  completed_at?: string
  failed_at?: string
  error?: {
    message: string
  }
}

const TaskListPage = () => {
  const [tasks, setTasks] = useState<TaskItem[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [statusFilter, setStatusFilter] = useState<TaskStatus | 'all'>('all')
  const navigate = useNavigate()

  useEffect(() => {
    fetchTasks()
  }, [page, statusFilter])

  const fetchTasks = async () => {
    setLoading(true)
    try {
      const params: any = { page, page_size: 20 }
      if (statusFilter !== 'all') {
        params.status = statusFilter
      }

      const response = await apiClient.get('/manga/tasks', { params })
      setTasks(response.data.items)
      setTotalPages(response.data.pagination.total_pages)
    } catch (err) {
      console.error('Failed to fetch tasks:', err)
    } finally {
      setLoading(false)
    }
  }

  const getStatusColor = (status: TaskStatus) => {
    const colors = {
      pending: 'gray',
      processing: 'blue',
      completed: 'green',
      failed: 'red',
      cancelled: 'orange',
    }
    return colors[status]
  }

  const getStatusLabel = (status: TaskStatus) => {
    const labels = {
      pending: '等待中',
      processing: '处理中',
      completed: '已完成',
      failed: '失败',
      cancelled: '已取消',
    }
    return labels[status]
  }

  return (
    <div className="task-list-page">
      <header>
        <h1>我的任务</h1>
        <button onClick={() => navigate('/')}>创建新任务</button>
      </header>

      {/* 筛选器 */}
      <div className="filters">
        <label>状态筛选:</label>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as TaskStatus | 'all')}
        >
          <option value="all">全部</option>
          <option value="completed">已完成</option>
          <option value="processing">处理中</option>
          <option value="failed">失败</option>
          <option value="pending">等待中</option>
          <option value="cancelled">已取消</option>
        </select>
      </div>

      {/* 任务列表 */}
      {loading ? (
        <div className="loading">加载中...</div>
      ) : tasks.length === 0 ? (
        <div className="empty">
          <p>暂无任务</p>
          <button onClick={() => navigate('/')}>创建第一个任务</button>
        </div>
      ) : (
        <div className="task-grid">
          {tasks.map((task) => (
            <div
              key={task.task_id}
              className="task-card"
              onClick={() => navigate(`/task/${task.task_id}`)}
            >
              <div className="task-header">
                <h3>{task.title}</h3>
                <span className={`status-badge status-${getStatusColor(task.status)}`}>
                  {getStatusLabel(task.status)}
                </span>
              </div>

              <div className="task-progress">
                <div className="progress-bar">
                  <div
                    className="progress-fill"
                    style={{ width: `${task.progress.percentage}%` }}
                  />
                </div>
                <span className="progress-text">
                  {task.progress.percentage}% - {task.progress.current_step}
                </span>
              </div>

              {task.status === 'completed' && (
                <div className="task-stats">
                  <span>角色: {task.character_count}</span>
                  <span>场景: {task.scene_count}</span>
                </div>
              )}

              {task.status === 'failed' && task.error && (
                <div className="task-error">{task.error.message}</div>
              )}

              <div className="task-footer">
                <span>创建于: {new Date(task.created_at).toLocaleString()}</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 分页器 */}
      {totalPages > 1 && (
        <div className="pagination">
          <button
            disabled={page === 1}
            onClick={() => setPage(page - 1)}
          >
            上一页
          </button>
          <span>第 {page} / {totalPages} 页</span>
          <button
            disabled={page === totalPages}
            onClick={() => setPage(page + 1)}
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}

export default TaskListPage
```

### 4.5 前端状态管理

**任务状态类型定义** (TypeScript):

```typescript
// 任务状态枚举
type TaskStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';

// 任务进度详情
interface TaskProgressDetails {
  characters_extracted: number;
  characters_generated: number;
  scenes_divided: number;
  scenes_generated: number;
}

// 任务进度信息
interface TaskProgress {
  current_step: string;
  current_step_index: number;
  total_steps: number;
  percentage: number;
  details?: TaskProgressDetails;
}

// 角色信息
interface Character {
  id: string;
  name: string;
  reference_image_url: string;
}

// 场景信息
interface Scene {
  id: string;
  sequence_num: number;
  description: string;
  image_url: string;
}

// 任务结果
interface TaskResult {
  novel_id: string;
  title: string;
  character_count: number;
  scene_count: number;
  characters: Character[];
  scenes: Scene[];
}

// 任务错误信息
interface TaskError {
  code: number;
  message: string;
  retry_able: boolean;
}

// 任务数据
interface TaskData {
  task_id: string;
  status: TaskStatus;
  progress: TaskProgress;
  result?: TaskResult;
  error?: TaskError;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  failed_at?: string;
}
```

### 4.6 前端交互流程详解

#### 4.6.1 创建任务流程

```typescript
// 1. 用户点击"开始生成"按钮
const handleGenerate = async () => {
  // 2. 验证表单
  if (!title) {
    setError('请输入标题');
    return;
  }

  if (inputMode === 'file' && !selectedFile) {
    setError('请选择文件');
    return;
  }

  if (inputMode === 'text' && !textContent.trim()) {
    setError('请输入小说内容');
    return;
  }

  try {
    setUploading(true);
    setError(null);

    // 3. 读取文件内容或使用文本输入
    let content = '';
    if (inputMode === 'file' && selectedFile) {
      content = await readFileContent(selectedFile);
    } else {
      content = textContent;
    }

    // 4. 调用后端接口创建任务
    const response = await apiClient.post<{ task_id: string }>('/manga/generate', {
      title,
      author: author || 'Unknown',
      content,
    });

    // 5. 获取任务ID，跳转到任务详情页
    const taskId = response.data.task_id;
    navigate(`/task/${taskId}`);
  } catch (err) {
    // 6. 处理错误
    const error = err instanceof Error ? err : new Error('生成失败');
    setError(error.message);
  } finally {
    setUploading(false);
  }
};
```

#### 4.6.2 任务详情页轮询流程

```typescript
// TaskDetailPage.tsx

const TaskDetailPage = () => {
  const { taskId } = useParams<{ taskId: string }>();
  const [taskData, setTaskData] = useState<TaskData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 轮询函数
  const fetchTaskStatus = async () => {
    try {
      const response = await apiClient.get<TaskData>(`/manga/task/${taskId}`);
      setTaskData(response.data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取任务状态失败');
    } finally {
      setLoading(false);
    }
  };

  // 启动轮询
  useEffect(() => {
    if (!taskId) return;

    // 立即获取一次
    fetchTaskStatus();

    // 设置轮询定时器 (每2秒)
    const pollInterval = setInterval(() => {
      // 只在任务未完成时继续轮询
      if (taskData?.status === 'processing' || taskData?.status === 'pending') {
        fetchTaskStatus();
      } else {
        clearInterval(pollInterval);
      }
    }, 2000);

    // 清理定时器
    return () => clearInterval(pollInterval);
  }, [taskId, taskData?.status]);

  // 15分钟超时检测
  useEffect(() => {
    const timeout = setTimeout(() => {
      if (taskData?.status === 'processing' || taskData?.status === 'pending') {
        setError('任务执行超时，请稍后刷新查看或重试');
      }
    }, 15 * 60 * 1000); // 15分钟

    return () => clearTimeout(timeout);
  }, [taskData?.status]);

  // 渲染UI...
};
```

#### 4.6.3 UI组件示例

```typescript
// 进度条组件
const ProgressBar = ({ percentage }: { percentage: number }) => (
  <div className="progress-bar">
    <div className="progress-fill" style={{ width: `${percentage}%` }} />
    <span className="progress-text">{percentage}%</span>
  </div>
);

// 步骤指示器
const StepIndicator = ({ progress }: { progress: TaskProgress }) => (
  <div className="step-indicator">
    <div className="step-label">当前步骤: {progress.current_step}</div>
    <div className="step-progress">
      {progress.current_step_index} / {progress.total_steps}
    </div>
  </div>
);

// 任务状态徽章
const StatusBadge = ({ status }: { status: TaskStatus }) => {
  const statusConfig = {
    pending: { label: '等待中', color: 'gray' },
    processing: { label: '处理中', color: 'blue' },
    completed: { label: '已完成', color: 'green' },
    failed: { label: '失败', color: 'red' },
    cancelled: { label: '已取消', color: 'orange' },
  };

  const config = statusConfig[status];
  return <span className={`badge badge-${config.color}`}>{config.label}</span>;
};

// 结果展示区域
const ResultSection = ({ result }: { result: TaskResult }) => (
  <div className="result-section">
    <h2>生成结果</h2>

    {/* 角色列表 */}
    <section>
      <h3>角色 ({result.character_count})</h3>
      <div className="character-grid">
        {result.characters.map(char => (
          <div key={char.id} className="character-card">
            <img src={char.reference_image_url} alt={char.name} />
            <p>{char.name}</p>
          </div>
        ))}
      </div>
    </section>

    {/* 场景列表 */}
    <section>
      <h3>场景 ({result.scene_count})</h3>
      <div className="scene-grid">
        {result.scenes.map(scene => (
          <div key={scene.id} className="scene-card">
            <img src={scene.image_url} alt={scene.description} />
            <p>第{scene.sequence_num}幕: {scene.description}</p>
          </div>
        ))}
      </div>
    </section>

    {/* 操作按钮 */}
    <div className="actions">
      <Button onClick={() => navigate(`/novels/${result.novel_id}`)}>
        查看详情
      </Button>
      <Button variant="secondary" onClick={() => navigate('/export')}>
        导出视频
      </Button>
    </div>
  </div>
);
```

### 4.7 错误处理与重试

```typescript
// 错误处理组件
const ErrorDisplay = ({ error, onRetry }: { error: TaskError; onRetry: () => void }) => (
  <div className="error-display">
    <div className="error-icon">❌</div>
    <h3>生成失败</h3>
    <p className="error-message">{error.message}</p>
    <p className="error-code">错误代码: {error.code}</p>

    {error.retry_able && (
      <Button onClick={onRetry}>重新生成</Button>
    )}

    <Button variant="secondary" onClick={() => navigate('/')}>
      返回首页
    </Button>
  </div>
);

// 重试逻辑
const handleRetry = async () => {
  // 使用相同参数重新创建任务
  const response = await apiClient.post('/manga/generate', {
    title: taskData?.result?.title,
    author: 'Unknown',
    content: savedContent, // 需要保存原始内容
  });

  navigate(`/task/${response.data.task_id}`);
};
```

### 4.8 取消任务功能

```typescript
const handleCancel = async () => {
  if (!confirm('确定要取消任务吗？')) return;

  try {
    await apiClient.post(`/manga/task/${taskId}/cancel`);
    setTaskData(prev => prev ? { ...prev, status: 'cancelled' } : null);
  } catch (err) {
    setError('取消任务失败');
  }
};

// UI中的取消按钮
{(taskData?.status === 'pending' || taskData?.status === 'processing') && (
  <Button variant="danger" onClick={handleCancel}>
    取消任务
  </Button>
)}
```

---

## 5. 后端实现要点

### 5.1 数据库表设计

#### Task表 (新增)

```sql
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,  -- 新增：关联用户ID
    novel_id VARCHAR(36),
    status ENUM('pending', 'processing', 'completed', 'failed', 'cancelled') NOT NULL,
    progress_step VARCHAR(100),
    progress_step_index INT DEFAULT 0,
    progress_percentage INT DEFAULT 0,
    progress_details JSON,
    error_code INT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    failed_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),        -- 新增：用户ID索引
    INDEX idx_novel_id (novel_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_user_status (user_id, status)  -- 新增：组合索引，用于按用户和状态查询
);
```

#### 进度详情JSON结构

```json
{
  "characters_extracted": 3,
  "characters_generated": 2,
  "scenes_divided": 8,
  "scenes_generated": 5
}
```

### 5.2 Supabase 认证中间件实现

#### Go Supabase JWT 验证

```go
// backend/internal/infrastructure/middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// SupabaseAuth 验证 Supabase JWT Token
func (m *AuthMiddleware) SupabaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: 缺少 Authorization 头",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 2. 解析 Bearer Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Authorization 格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 3. 验证 JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Token 无效或已过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 4. 提取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Token Claims 解析失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: 用户ID不存在",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 5. 将用户ID存入上下文
		c.Set("user_id", userID)
		c.Set("user_email", claims["email"])

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}
```

#### 应用认证中间件

```go
// backend/cmd/main.go
func main() {
	// ... 初始化配置和依赖 ...

	// 初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.Supabase.JWTSecret)

	r := gin.Default()
	r.Use(cors.Default())

	// 公开路由（无需认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ai-motion",
		})
	})

	// API 路由（需要认证）
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware.SupabaseAuth()) // 应用认证中间件
	{
		// 漫画生成路由
		mangaGroup := v1.Group("/manga")
		{
			mangaGroup.POST("/generate", mangaWorkflowHandler.GenerateManga)
			mangaGroup.GET("/task/:task_id", mangaWorkflowHandler.GetTaskStatus)
			mangaGroup.GET("/tasks", mangaWorkflowHandler.GetTaskList)
			mangaGroup.POST("/task/:task_id/cancel", mangaWorkflowHandler.CancelTask)
		}
	}

	// 启动服务器
	r.Run(":" + cfg.Server.Port)
}
```

### 5.3 Go实现示例

#### Task Domain Entity

```go
// backend/internal/domain/task/task.go
package task

import (
    "time"
    "github.com/google/uuid"
)

type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusProcessing TaskStatus = "processing"
    TaskStatusCompleted  TaskStatus = "completed"
    TaskStatusFailed     TaskStatus = "failed"
    TaskStatusCancelled  TaskStatus = "cancelled"
)

type ProgressDetails struct {
    CharactersExtracted  int `json:"characters_extracted"`
    CharactersGenerated  int `json:"characters_generated"`
    ScenesDivided        int `json:"scenes_divided"`
    ScenesGenerated      int `json:"scenes_generated"`
}

type Task struct {
    ID                 string          `json:"id"`
    UserID             string          `json:"user_id"`  // 新增：用户ID
    NovelID            string          `json:"novel_id"`
    Status             TaskStatus      `json:"status"`
    ProgressStep       string          `json:"progress_step"`
    ProgressStepIndex  int             `json:"progress_step_index"`
    ProgressPercentage int             `json:"progress_percentage"`
    ProgressDetails    ProgressDetails `json:"progress_details"`
    ErrorCode          int             `json:"error_code,omitempty"`
    ErrorMessage       string          `json:"error_message,omitempty"`
    CreatedAt          time.Time       `json:"created_at"`
    UpdatedAt          time.Time       `json:"updated_at"`
    CompletedAt        *time.Time      `json:"completed_at,omitempty"`
    FailedAt           *time.Time      `json:"failed_at,omitempty"`
    cancelChan         chan struct{}   // 取消信号通道
}

func NewTask(novelID string) *Task {
    return &Task{
        ID:                 uuid.New().String(),
        NovelID:            novelID,
        Status:             TaskStatusPending,
        ProgressStep:       "等待中",
        ProgressStepIndex:  0,
        ProgressPercentage: 0,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
        cancelChan:         make(chan struct{}),
    }
}

func (t *Task) UpdateProgress(step string, stepIndex int, percentage int, details ProgressDetails) {
    t.Status = TaskStatusProcessing
    t.ProgressStep = step
    t.ProgressStepIndex = stepIndex
    t.ProgressPercentage = percentage
    t.ProgressDetails = details
    t.UpdatedAt = time.Now()
}

func (t *Task) MarkCompleted() {
    t.Status = TaskStatusCompleted
    t.ProgressStep = "完成"
    t.ProgressStepIndex = 6
    t.ProgressPercentage = 100
    now := time.Now()
    t.CompletedAt = &now
    t.UpdatedAt = now
}

func (t *Task) MarkFailed(errCode int, errMsg string) {
    t.Status = TaskStatusFailed
    t.ErrorCode = errCode
    t.ErrorMessage = errMsg
    now := time.Now()
    t.FailedAt = &now
    t.UpdatedAt = now
}

func (t *Task) Cancel() {
    t.Status = TaskStatusCancelled
    t.UpdatedAt = time.Now()
    close(t.cancelChan)
}

func (t *Task) IsCancelled() bool {
    select {
    case <-t.cancelChan:
        return true
    default:
        return false
    }
}
```

#### Handler实现（带用户认证）

```go
// backend/internal/interfaces/http/handler/manga_workflow_handler.go
func (h *MangaWorkflowHandler) GenerateManga(c *gin.Context) {
    var req dto.GenerateMangaRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.InvalidParams(c, err.Error())
        return
    }

    // 获取当前用户ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Error(c, 20001, "未授权")
        return
    }

    // 创建任务（传入用户ID）
    task, err := h.workflowService.CreateTask(c.Request.Context(), userID, &req)
    if err != nil {
        response.InternalError(c, "创建任务失败")
        return
    }

    // 异步执行任务
    go h.workflowService.ExecuteTask(context.Background(), task.ID)

    // 立即返回任务ID
    response.Success(c, gin.H{
        "task_id":    task.ID,
        "status":     task.Status,
        "created_at": task.CreatedAt,
    })
}

func (h *MangaWorkflowHandler) GetTaskStatus(c *gin.Context) {
    taskID := c.Param("task_id")

    // 获取当前用户ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Error(c, 20001, "未授权")
        return
    }

    // 获取任务状态（验证用户权限）
    taskStatus, err := h.workflowService.GetTaskStatus(c.Request.Context(), userID, taskID)
    if err != nil {
        response.ResourceNotFound(c, "任务不存在")
        return
    }

    response.Success(c, taskStatus)
}

func (h *MangaWorkflowHandler) GetTaskList(c *gin.Context) {
    // 获取当前用户ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Error(c, 20001, "未授权")
        return
    }

    // 解析查询参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    status := c.Query("status")

    // 获取任务列表
    tasks, pagination, err := h.workflowService.GetTaskList(c.Request.Context(), userID, page, pageSize, status)
    if err != nil {
        response.InternalError(c, "获取任务列表失败")
        return
    }

    response.SuccessList(c, tasks, pagination)
}

func (h *MangaWorkflowHandler) CancelTask(c *gin.Context) {
    taskID := c.Param("task_id")

    // 获取当前用户ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Error(c, 20001, "未授权")
        return
    }

    // 取消任务（验证用户权限）
    err := h.workflowService.CancelTask(c.Request.Context(), userID, taskID)
    if err != nil {
        response.Error(c, 10001, err.Error())
        return
    }

    response.SuccessWithMessage(c, "任务已取消", gin.H{
        "task_id": taskID,
        "status":  "cancelled",
    })
}
```

#### Service实现 (异步执行)

```go
// backend/internal/application/service/manga_workflow_service.go
func (s *MangaWorkflowService) ExecuteTask(ctx context.Context, taskID string) {
    task, err := s.taskRepo.FindByID(ctx, taskID)
    if err != nil {
        return
    }

    // 步骤1: 解析小说
    task.UpdateProgress("解析小说", 1, 17, task.ProgressDetails)
    s.taskRepo.Save(ctx, task)
    if task.IsCancelled() return

    // ... 执行解析逻辑 ...

    // 步骤2: 提取角色
    task.UpdateProgress("提取角色", 2, 33, ProgressDetails{
        CharactersExtracted: len(characters),
    })
    s.taskRepo.Save(ctx, task)
    if task.IsCancelled() return

    // ... 提取角色逻辑 ...

    // 步骤3: 生成角色参考图
    for i, char := range characters {
        if task.IsCancelled() return

        task.UpdateProgress("生成角色参考图", 3, 50, ProgressDetails{
            CharactersExtracted: len(characters),
            CharactersGenerated: i,
        })
        s.taskRepo.Save(ctx, task)

        // ... 生成参考图逻辑 ...
    }

    // 步骤4: 划分场景
    task.UpdateProgress("划分场景", 4, 67, ProgressDetails{
        CharactersExtracted: len(characters),
        CharactersGenerated: len(characters),
        ScenesDivided:       len(scenes),
    })
    s.taskRepo.Save(ctx, task)
    if task.IsCancelled() return

    // ... 划分场景逻辑 ...

    // 步骤5: 生成场景图片
    for i, scene := range scenes {
        if task.IsCancelled() return

        task.UpdateProgress("生成场景图片", 5, 83, ProgressDetails{
            CharactersExtracted: len(characters),
            CharactersGenerated: len(characters),
            ScenesDivided:       len(scenes),
            ScenesGenerated:     i,
        })
        s.taskRepo.Save(ctx, task)

        // ... 生成场景图片逻辑 ...
    }

    // 步骤6: 完成
    task.MarkCompleted()
    s.taskRepo.Save(ctx, task)
}
```

---

## 6. 优化建议

### 6.1 性能优化

1. **数据库索引**: 在`tasks`表的`status`和`created_at`字段上建立索引
2. **缓存**: 使用Redis缓存任务状态，减少数据库查询
3. **并发控制**: 限制同时执行的任务数量 (使用worker pool)
4. **分批生成**: 大批量场景图片可分批生成，避免超时

### 6.2 用户体验优化

1. **实时通知**: 使用WebSocket推送任务进度，替代轮询
2. **断点续传**: 任务失败后支持从失败步骤重试
3. **预估时间**: 根据历史数据预估任务完成时间
4. **进度动画**: 使用动画展示当前步骤，提升视觉体验

### 6.3 可靠性优化

1. **任务持久化**: 服务重启后能恢复正在执行的任务
2. **错误重试**: AI服务调用失败时自动重试 (指数退避)
3. **超时处理**: 单个步骤超时后标记任务失败
4. **监控告警**: 监控任务失败率，及时告警

---

## 7. 前后端联调检查清单

### 7.1 认证功能检查

- [ ] 前端 Supabase Client 正确初始化
- [ ] 登录功能正常，返回 JWT Token
- [ ] Token 自动添加到 API 请求头
- [ ] Token 过期后自动刷新
- [ ] 未登录用户访问受保护路由重定向到登录页
- [ ] 后端中间件正确验证 JWT Token
- [ ] 用户ID正确提取并存入上下文

### 7.2 接口对接检查

- [ ] 所有接口请求头正确携带 Authorization Token
- [ ] 未授权请求返回 401 状态码和错误信息

- [ ] 创建任务接口返回正确的任务ID
- [ ] 任务状态接口返回完整的进度信息
- [ ] 任务列表接口正确分页和筛选
- [ ] 任务完成后返回完整的结果数据
- [ ] 任务失败后返回可重试标志
- [ ] 取消任务接口正确更新任务状态
- [ ] 用户只能访问自己的任务（数据隔离）

### 7.3 前端功能检查

- [ ] 登录后跳转到首页
- [ ] 文件上传和文本输入模式切换正常
- [ ] 表单验证正确 (标题、内容)
- [ ] 任务创建后跳转到详情页
- [ ] 任务列表页正确显示用户的所有任务
- [ ] 任务列表筛选和分页功能正常
- [ ] 点击任务卡片跳转到任务详情页
- [ ] 轮询正常启动和停止
- [ ] 进度条和步骤指示器实时更新
- [ ] 任务完成后显示结果
- [ ] 任务失败后显示错误信息
- [ ] 取消按钮功能正常
- [ ] 退出登录功能正常

### 7.4 异常场景测试

- [ ] Token过期时自动刷新
- [ ] 网络超时或断开
- [ ] 后端服务重启
- [ ] AI服务调用失败
- [ ] 数据库连接失败
- [ ] 任务执行超过15分钟
- [ ] 用户中途刷新页面
- [ ] 并发创建多个任务
- [ ] 尝试访问其他用户的任务（权限验证）

---

## 8. 总结

本文档描述了简化版漫画生成功能的完整设计方案，包含用户认证、任务管理和数据隔离机制。

**核心变化**:
- 从多步骤手动操作简化为一键生成
- 引入异步任务机制，提升用户体验
- 前端通过轮询实时展示任务进度
- 集成 Supabase 认证，实现用户管理
- 任务数据按用户隔离，保证数据安全

**关键接口** (共4个):
1. `POST /api/v1/manga/generate` - 创建任务
2. `GET /api/v1/manga/task/:task_id` - 查询任务状态
3. `GET /api/v1/manga/tasks` - 获取任务列表
4. `POST /api/v1/manga/task/:task_id/cancel` - 取消任务

**认证机制**:
- 使用 Supabase Authentication 管理用户
- JWT Token 认证，支持自动刷新
- 所有 API 接口需要 Token 认证
- 后端中间件验证 Token 并提取用户 ID

**数据隔离**:
- 每个任务关联一个 `user_id`
- 用户只能访问自己的任务
- 查询时自动添加用户过滤条件
- 防止横向越权访问

**前端页面**:
1. **登录页** - Supabase 邮箱/密码登录
2. **首页** - 输入小说内容，创建任务
3. **任务列表页** - 查看所有任务，支持筛选和分页
4. **任务详情页** - 实时查看任务进度和结果

**前端交互流程**:
1. 用户登录 → 获取 JWT Token → 存储到本地
2. 访问 API 时自动添加 Token 到请求头
3. 用户提交表单 → 创建任务 → 获取任务 ID
4. 跳转到任务详情页 → 启动轮询 → 实时显示进度
5. 任务完成 → 停止轮询 → 展示结果

**技术要点**:
- **后端**: Go + Gin, DDD 架构，Supabase JWT 验证
- **前端**: React + TypeScript, Supabase Client, Axios 拦截器
- **异步任务**: Goroutine 异步执行，任务状态持久化
- **数据库**: 任务表新增 `user_id` 字段和组合索引
- **轮询机制**: 前端每 2 秒轮询，15 分钟超时
- **错误处理**: Token 过期自动刷新，支持任务取消和重试

**安全性**:
- JWT Token 签名验证
- 用户数据严格隔离
- 防止 CSRF 和 XSS 攻击
- API 访问权限控制

---

## 9. 下一步开发任务

### 9.1 后端开发任务

1. **认证中间件**
   - [ ] 实现 Supabase JWT 验证中间件
   - [ ] 配置 JWT Secret 环境变量
   - [ ] 添加用户上下文提取逻辑

2. **数据库迁移**
   - [ ] 在 `tasks` 表添加 `user_id` 字段
   - [ ] 创建 `user_id` 和组合索引
   - [ ] 更新现有 Task Entity 添加 UserID 字段

3. **Handler 更新**
   - [ ] 修改 `GenerateManga` 提取用户 ID
   - [ ] 实现 `GetTaskList` 接口
   - [ ] 所有接口添加用户权限验证

4. **Service 更新**
   - [ ] `CreateTask` 方法添加 userID 参数
   - [ ] `GetTaskStatus` 添加用户权限验证
   - [ ] `GetTaskList` 实现分页和筛选
   - [ ] `CancelTask` 添加用户权限验证

### 9.2 前端开发任务

1. **认证集成**
   - [ ] 安装 `@supabase/supabase-js`
   - [ ] 创建 Supabase Client
   - [ ] 实现登录/注册页面
   - [ ] 实现 `useAuth` Hook
   - [ ] 创建 `ProtectedRoute` 组件

2. **API Client 更新**
   - [ ] 实现 Axios 请求拦截器（自动添加 Token）
   - [ ] 实现响应拦截器（处理 401 错误）
   - [ ] Token 过期时自动刷新

3. **任务列表页**
   - [ ] 创建 `TaskListPage` 组件
   - [ ] 实现任务列表展示
   - [ ] 实现状态筛选功能
   - [ ] 实现分页功能
   - [ ] 点击任务跳转到详情页

4. **路由更新**
   - [ ] 添加 `/login` 路由
   - [ ] 添加 `/tasks` 路由
   - [ ] 所有受保护路由使用 `ProtectedRoute`
   - [ ] 首页添加"查看任务列表"链接

### 9.3 环境配置

1. **后端配置**
   - [ ] 添加 `SUPABASE_JWT_SECRET` 到 `.env`
   - [ ] 安装 `github.com/golang-jwt/jwt/v5`
   - [ ] 更新配置结构体

2. **前端配置**
   - [ ] 添加 `VITE_SUPABASE_URL` 到 `.env`
   - [ ] 添加 `VITE_SUPABASE_ANON_KEY` 到 `.env`
   - [ ] 更新 API Base URL 配置

### 9.4 测试任务

1. **单元测试**
   - [ ] 认证中间件测试
   - [ ] Task Entity 测试
   - [ ] Handler 层测试（带用户权限）

2. **集成测试**
   - [ ] 完整的登录 → 创建任务 → 查看列表流程
   - [ ] Token 过期刷新测试
   - [ ] 数据隔离测试（不同用户）

3. **端到端测试**
   - [ ] 使用 Docker Compose 启动完整环境
   - [ ] 测试所有用户交互流程
   - [ ] 测试异常场景处理

---

**相关文档**:
- [API.md](API.md) - 完整API文档
- [ARCHITECTURE.md](ARCHITECTURE.md) - 系统架构设计
- [CHARACTER_CONSISTENCY.md](CHARACTER_CONSISTENCY.md) - 角色一致性设计
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南

---

*文档版本: v2.0*
*最后更新: 2025-01-26*
*新增内容: 用户认证、任务列表、数据隔离*
