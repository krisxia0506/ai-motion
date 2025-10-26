# 开发指南

本文档为开发者提供详细的开发指南，包括环境搭建、代码规范、测试和贡献流程。

---

## 开发环境设置

### 前置要求

- **Go**: 1.24 或更高版本
- **Node.js**: 20 或更高版本
- **Supabase**: PostgreSQL 数据库 (代替 MySQL)
- **Git**: 用于版本控制
- **Docker**: 20+ (可选，推荐)

### 克隆项目

```bash
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion
```

### 后端开发环境

#### 1. 安装 Go 依赖

```bash
cd backend

# 初始化模块 (如果需要)
go mod init github.com/xiajiayi/ai-motion

# 下载依赖
go mod download

# 整理依赖
go mod tidy
```

#### 2. 配置环境变量

```bash
# 在项目根目录创建 .env 文件
cp .env.example .env

# 编辑配置
vim .env
```

详细配置请参考 [配置文档](./CONFIGURATION.md)。

#### 3. 启动后端服务

```bash
# 开发模式运行
go run cmd/main.go

# 或使用 Make 命令
cd ..
make dev-backend
```

后端服务将运行在 http://localhost:8080

### 前端开发环境

#### 1. 安装 Node.js 依赖

```bash
cd frontend

# 使用 npm
npm install

# 或使用 pnpm (推荐)
pnpm install

# 或使用 yarn
yarn install
```

#### 2. 配置前端环境变量

```bash
# 创建 .env 文件
cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=AI-Motion
EOF
```

#### 3. 启动前端开发服务器

```bash
# 使用 npm
npm run dev

# 或使用 Make 命令
cd ..
make dev-frontend
```

前端服务将运行在 http://localhost:3000

### 使用 Make 命令

项目提供了便捷的 Make 命令：

```bash
# 查看所有可用命令
make help

# 安装所有依赖
make install

# 启动开发环境 (后端 + 前端)
make dev

# 构建项目
make build

# 运行测试
make test

# 清理编译产物
make clean
```

---

## 项目架构

### 后端架构 (DDD)

AI-Motion 后端采用 **DDD (领域驱动设计)** 架构：

```
backend/
├── cmd/                    # 应用入口
│   └── main.go            # 主程序
├── internal/              # 内部业务逻辑
│   ├── domain/           # 领域层 - 核心业务逻辑
│   │   ├── novel/        # 小说领域
│   │   ├── character/    # 角色领域
│   │   ├── scene/        # 场景领域
│   │   └── media/        # 媒体领域
│   ├── application/      # 应用层 - 用例编排
│   │   ├── service/      # 应用服务
│   │   └── dto/          # 数据传输对象
│   ├── infrastructure/   # 基础设施层 - 技术实现
│   │   ├── repository/   # 数据仓储
│   │   ├── ai/           # AI 服务客户端
│   │   └── storage/      # 文件存储
│   └── interfaces/       # 接口层 - HTTP/gRPC
│       ├── http/         # HTTP 处理器
│       └── middleware/   # 中间件
└── pkg/                  # 公共包
    └── utils/           # 工具函数
```

#### DDD 分层职责

**1. 领域层 (Domain Layer)**
- 核心业务逻辑
- 领域模型和实体
- 领域服务
- 仓储接口定义

**2. 应用层 (Application Layer)**
- 用例编排
- 事务管理
- DTO 转换
- 应用服务

**3. 基础设施层 (Infrastructure Layer)**
- 外部服务集成 (AI API)
- 数据库实现
- 文件存储实现
- 第三方库集成

**4. 接口层 (Interface Layer)**
- HTTP 路由和处理器
- 请求/响应序列化
- 中间件 (认证、日志等)
- API 文档

### 前端架构

```
frontend/
├── src/
│   ├── components/       # React 组件
│   │   ├── common/      # 通用组件
│   │   └── features/    # 功能组件
│   ├── pages/           # 页面组件
│   ├── services/        # API 服务
│   ├── hooks/           # 自定义 Hooks
│   ├── utils/           # 工具函数
│   ├── types/           # TypeScript 类型定义
│   ├── styles/          # 全局样式
│   ├── App.tsx          # 根组件
│   └── main.tsx         # 入口文件
├── public/              # 静态资源
└── vite.config.ts       # Vite 配置
```

---

## 代码规范

### Go 代码规范

遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)。

**命名规范**:
```go
// 包名: 小写，单数
package novel

// 接口: 名词或形容词，大写开头
type NovelRepository interface {}

// 结构体: 大写开头 (公开) 或小写开头 (私有)
type Novel struct {}
type novelService struct {}

// 函数/方法: 大写开头 (公开) 或小写开头 (私有)
func ParseNovel() {}
func validateChapter() {}

// 常量: 大写开头或全大写
const MaxChapters = 1000
const DEFAULT_PAGE_SIZE = 20
```

**错误处理**:
```go
// 总是检查错误
result, err := someOperation()
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}

// 使用有意义的错误信息
if novel == nil {
    return errors.New("novel not found")
}
```

**注释规范**:
```go
// Package novel provides novel parsing and management functionality.
package novel

// Novel represents a novel document with chapters and metadata.
type Novel struct {
    ID       string
    Title    string
    Chapters []Chapter
}

// ParseNovel parses a novel file and extracts chapters and characters.
// It returns an error if the file format is invalid.
func ParseNovel(file io.Reader) (*Novel, error) {
    // Implementation
}
```

**HTTP 响应规范**:

从 `v0.1.0-alpha` 开始,所有 HTTP Handler 必须使用统一的响应助手函数,位于 `backend/internal/interfaces/http/response` 包。

```go
import "github.com/xiajiayi/ai-motion/internal/interfaces/http/response"

// ✅ 正确: 使用响应助手
func (h *NovelHandler) Upload(c *gin.Context) {
    var req dto.UploadNovelRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.InvalidParams(c, "Invalid request: "+err.Error())
        return
    }
    
    novel, err := h.novelService.Upload(c.Request.Context(), &req)
    if err != nil {
        response.InternalError(c, "Failed to upload: "+err.Error())
        return
    }
    
    response.Success(c, novel)
}

// ❌ 错误: 手动创建响应
func (h *NovelHandler) Upload(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "data": novel,  // 缺少 code 和 message 字段
    })
}
```

**可用的响应助手**:
- `response.Success(c, data)` - 成功响应
- `response.SuccessWithMessage(c, message, data)` - 自定义成功消息
- `response.SuccessList(c, items, page, pageSize, total)` - 分页列表
- `response.InvalidParams(c, message)` - 参数错误 (10001)
- `response.ResourceNotFound(c, message)` - 资源不存在 (10002)
- `response.AIServiceError(c, message)` - AI 服务错误 (40001)
- `response.InternalError(c, message)` - 内部错误 (50002)

所有响应遵循统一格式:
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

详细说明请参考 [backend/CLAUDE.md](../backend/CLAUDE.md#unified-response-helpers) 和 [API 设计规范](./API_DESIGN_GUIDELINES.md)。

### TypeScript/React 代码规范

遵循 [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript) 和 [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)。

**组件规范**:
```typescript
// 使用函数组件和 Hooks
import React, { useState, useEffect } from 'react';

interface NovelListProps {
  onSelect: (id: string) => void;
}

export const NovelList: React.FC<NovelListProps> = ({ onSelect }) => {
  const [novels, setNovels] = useState<Novel[]>([]);

  useEffect(() => {
    // Fetch novels
  }, []);

  return (
    <div className="novel-list">
      {novels.map(novel => (
        <div key={novel.id} onClick={() => onSelect(novel.id)}>
          {novel.title}
        </div>
      ))}
    </div>
  );
};
```

**API 服务规范**:
```typescript
// services/novelApi.ts
import axios from 'axios';

const API_BASE = import.meta.env.VITE_API_BASE_URL;

export const novelApi = {
  async uploadNovel(file: File): Promise<Novel> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await axios.post(`${API_BASE}/api/v1/novel/upload`, formData);
    return response.data;
  },

  async getNovel(id: string): Promise<Novel> {
    const response = await axios.get(`${API_BASE}/api/v1/novel/${id}`);
    return response.data;
  },
};
```

---

## 数据库设计

**注意**: AI-Motion 已迁移至 Supabase (PostgreSQL + PostgREST)。不再使用 MySQL。

### 主要表结构

详细的表结构请参考 [DATABASE_DESIGN_STANDARDS.md](./DATABASE_DESIGN_STANDARDS.md) 和 [`database/schema/`](../database/schema/) 目录。

**aimotion_novel 表** (示例):
```sql
CREATE TABLE aimotion_novel (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100),
    status SMALLINT DEFAULT 0 CHECK (status IN (0, 1, 2)),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_novel IS '小说表';
COMMENT ON COLUMN aimotion_novel.status IS '状态:0-草稿,1-解析中,2-已完成';
```

### 数据库迁移

数据库迁移文件位于 `backend/internal/infrastructure/database/migrations/`。

Supabase 使用 PostgREST 自动生成 REST API,在 Supabase Dashboard 中执行 SQL 迁移即可创建表结构。

---

## 测试

### 后端测试

**单元测试**:
```go
// internal/domain/novel/novel_test.go
package novel_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNovelCreation(t *testing.T) {
    novel := NewNovel("Test Title", "Test Author")

    assert.NotNil(t, novel)
    assert.Equal(t, "Test Title", novel.Title)
    assert.Equal(t, "Test Author", novel.Author)
}
```

**运行测试**:
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/domain/novel

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试

**组件测试**:
```typescript
// components/NovelList.test.tsx
import { render, screen } from '@testing-library/react';
import { NovelList } from './NovelList';

describe('NovelList', () => {
  it('renders novel titles', () => {
    const novels = [
      { id: '1', title: 'Novel 1' },
      { id: '2', title: 'Novel 2' },
    ];

    render(<NovelList novels={novels} onSelect={() => {}} />);

    expect(screen.getByText('Novel 1')).toBeInTheDocument();
    expect(screen.getByText('Novel 2')).toBeInTheDocument();
  });
});
```

**运行测试**:
```bash
cd frontend

# 运行测试
npm test

# 查看覆盖率
npm run test:coverage
```

---

## Git 工作流

### 分支策略

- `main` - 主分支，始终保持可部署状态
- `develop` - 开发分支
- `feature/*` - 功能分支
- `bugfix/*` - 修复分支
- `release/*` - 发布分支

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 格式
<type>(<scope>): <subject>

# 示例
feat(novel): add novel parsing functionality
fix(api): resolve character list endpoint error
docs(readme): update installation guide
refactor(domain): simplify character model
test(novel): add unit tests for parser
chore(deps): update dependencies
```

**类型说明**:
- `feat`: 新功能
- `fix`: 修复
- `docs`: 文档
- `style`: 格式调整
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

### Pull Request 流程

1. 创建功能分支
```bash
git checkout -b feature/novel-parser
```

2. 开发并提交代码
```bash
git add .
git commit -m "feat(novel): implement novel parser"
```

3. 推送分支
```bash
git push origin feature/novel-parser
```

4. 创建 Pull Request
   - 填写清晰的标题和描述
   - 关联相关 Issue
   - 等待代码审查

5. 合并后删除分支
```bash
git branch -d feature/novel-parser
```

---

## 调试技巧

### 后端调试

**使用 Delve 调试器**:
```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/main.go

# 或在 VSCode 中使用 launch.json
```

**日志调试**:
```go
import "log"

log.Printf("Debug: %+v", variable)
log.Println("Processing novel:", novel.ID)
```

### 前端调试

**浏览器开发者工具**:
- F12 打开开发者工具
- Console 查看日志
- Network 查看网络请求
- React DevTools 查看组件状态

**VSCode 调试配置**:
```json
{
  "type": "chrome",
  "request": "launch",
  "name": "Launch Chrome",
  "url": "http://localhost:3000",
  "webRoot": "${workspaceFolder}/frontend/src"
}
```

---

## 性能优化

### 后端优化

- 使用数据库索引
- 实现查询结果缓存
- 优化 N+1 查询问题
- 使用连接池
- 实现 API 响应压缩

### 前端优化

- 代码分割 (Code Splitting)
- 懒加载组件
- 使用 React.memo 避免不必要的重渲染
- 优化图片加载
- 实现虚拟滚动

---

## 贡献指南

### 如何贡献

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some amazing feature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码审查要点

- 代码符合项目规范
- 包含必要的测试
- 文档已更新
- 没有引入新的警告或错误
- 性能没有显著下降

### 提交 Issue

提交 Issue 时请包含：
- 清晰的问题描述
- 重现步骤
- 预期行为
- 实际行为
- 环境信息 (OS, Go/Node 版本等)
- 相关日志或截图

---

## 资源链接

### Go 资源
- [Go 官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### React/TypeScript 资源
- [React 官方文档](https://react.dev/)
- [TypeScript 手册](https://www.typescriptlang.org/docs/)
- [Vite 文档](https://vitejs.dev/)

### DDD 资源
- [Domain-Driven Design Reference](https://www.domainlanguage.com/ddd/)
- [Implementing DDD](https://vaughnvernon.com/)

---

*开发文档版本: v0.1.0-alpha*
