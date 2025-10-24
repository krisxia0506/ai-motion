# Frontend 开发路线图

## 项目概述

AI-Motion 前端开发路线图，基于 React 19 + TypeScript + Vite 7，构建现代化、响应式的用户界面。

**当前版本:** v0.1.0-alpha
**目标版本:** v1.0.0
**最后更新:** 2025-10-24

---

## 图例

- ✅ 已完成
- 🚧 进行中
- ⏳ 计划中
- ❌ 未开始

---

## Phase 1: 项目基础 (Foundation)

### 1.1 项目初始化 ✅

- [x] Vite + React + TypeScript 项目创建
- [x] 项目目录结构搭建
- [x] ESLint 配置
- [x] TypeScript 配置 (tsconfig.json)
- [x] 基础依赖安装 (react-router-dom 等)
- [x] 开发服务器配置

**完成度:** 100%
**备注:** 基础脚手架已完成

---

### 1.2 开发环境配置 ⏳

#### 样式系统

- [ ] CSS 解决方案选择
  - [ ] 选项 A: Tailwind CSS (推荐)
  - [ ] 选项 B: CSS Modules + 全局样式
  - [ ] 选项 C: styled-components
- [ ] 全局样式变量定义
  - [ ] `src/styles/variables.css`
  - [ ] 颜色主题 (primary, secondary, danger, etc.)
  - [ ] 间距系统 (spacing scale)
  - [ ] 字体系统 (font family, sizes)
- [ ] 响应式设计断点
  - [ ] 移动端 (<768px)
  - [ ] 平板 (768px-1024px)
  - [ ] 桌面端 (>1024px)

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 2-3 天

#### UI 组件库集成 (可选)

- [ ] 组件库选择
  - [ ] 选项 A: Headless UI (@headlessui/react)
  - [ ] 选项 B: Radix UI (@radix-ui/react-*)
  - [ ] 选项 C: 自定义组件
- [ ] Toast 通知库
  - [ ] react-hot-toast 或 react-toastify
- [ ] Icon 库
  - [ ] react-icons 或 @heroicons/react

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 1-2 天

---

### 1.3 开发工具配置 ❌

- [ ] Prettier 配置
  - [ ] `.prettierrc` 文件
  - [ ] 与 ESLint 集成
- [ ] Git Hooks
  - [ ] Husky + lint-staged
  - [ ] Pre-commit 代码检查
- [ ] VSCode 配置
  - [ ] `.vscode/settings.json`
  - [ ] 推荐扩展列表

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 1 天

---

## Phase 2: 核心基础设施 (Core Infrastructure)

### 2.1 路由系统 ❌

#### React Router 配置

- [ ] 路由结构设计
  - [ ] `src/App.tsx` 路由配置
  - [ ] 路由懒加载 (React.lazy + Suspense)
- [ ] 页面路由定义
  - [ ] `/` - 首页
  - [ ] `/novels` - 小说列表
  - [ ] `/novels/:id` - 小说详情
  - [ ] `/novels/:id/characters` - 角色管理
  - [ ] `/novels/:id/generate` - 场景生成
  - [ ] `/novels/:id/export` - 导出管理
- [ ] 布局组件
  - [ ] `src/components/Layout.tsx`
  - [ ] 导航栏、侧边栏、页脚

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 2-3 天

---

### 2.2 API 集成层 ❌

#### Axios 配置

- [ ] HTTP 客户端配置
  - [ ] `src/services/api.ts`
  - [ ] Base URL 配置 (从环境变量读取)
  - [ ] 请求超时设置
  - [ ] 请求/响应拦截器
- [ ] 请求拦截器
  - [ ] 添加认证 Token (JWT)
  - [ ] 添加 Request ID
- [ ] 响应拦截器
  - [ ] 统一错误处理 (401, 500, etc.)
  - [ ] Toast 错误提示
  - [ ] 自动重试机制 (可选)

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 2-3 天

#### API Service 层

- [ ] Novel API Service
  - [ ] `src/services/novelApi.ts`
  - [ ] uploadNovel(data) -> Novel
  - [ ] getNovel(id) -> Novel
  - [ ] listNovels(page, size) -> Novel[]
  - [ ] parseNovel(id) -> void
  - [ ] deleteNovel(id) -> void
- [ ] Character API Service
  - [ ] `src/services/characterApi.ts`
  - [ ] getCharactersByNovel(novelId) -> Character[]
  - [ ] getCharacter(id) -> Character
  - [ ] updateCharacter(id, data) -> Character
  - [ ] generateReferenceImage(id) -> imageURL
- [ ] Scene API Service
  - [ ] `src/services/sceneApi.ts`
  - [ ] getScenesByNovel(novelId) -> Scene[]
  - [ ] getScene(id) -> Scene
- [ ] Generation API Service
  - [ ] `src/services/generationApi.ts`
  - [ ] generateSceneImage(sceneId) -> Media
  - [ ] generateSceneVideo(sceneId) -> Media
  - [ ] batchGenerate(sceneIds) -> Media[]
  - [ ] getGenerationStatus(taskId) -> Status

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 4-5 天

---

### 2.3 类型定义 ❌

#### TypeScript 类型

- [ ] Novel 类型
  - [ ] `src/types/novel.ts`
  - [ ] Novel, NovelStatus, UploadNovelRequest, NovelResponse
- [ ] Character 类型
  - [ ] `src/types/character.ts`
  - [ ] Character, Appearance, Personality
- [ ] Scene 类型
  - [ ] `src/types/scene.ts`
  - [ ] Scene, Dialogue, Description
- [ ] Media 类型
  - [ ] `src/types/media.ts`
  - [ ] Media, MediaType, MediaStatus, MediaMetadata
- [ ] API 响应类型
  - [ ] `src/types/api.ts`
  - [ ] ApiResponse<T>, PaginatedResponse<T>, ErrorResponse

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 2-3 天

---

### 2.4 状态管理 ❌

#### 状态管理方案选择

- [ ] 选择状态管理库
  - [ ] 选项 A: Zustand (推荐,轻量级)
  - [ ] 选项 B: Context API (简单场景)
  - [ ] 选项 C: Redux Toolkit (复杂场景)

#### Zustand Store 实现

- [ ] Novel Store
  - [ ] `src/store/novelStore.ts`
  - [ ] novels: Novel[]
  - [ ] selectedNovel: Novel | null
  - [ ] loading, error
  - [ ] Actions: setNovels, setSelectedNovel, addNovel, removeNovel
- [ ] Character Store
  - [ ] `src/store/characterStore.ts`
  - [ ] characters: Character[]
  - [ ] Actions: setCharacters, updateCharacter
- [ ] Generation Store
  - [ ] `src/store/generationStore.ts`
  - [ ] generationTasks: GenerationTask[]
  - [ ] Actions: addTask, updateTaskStatus
- [ ] UI Store
  - [ ] `src/store/uiStore.ts`
  - [ ] sidebarOpen, currentTheme
  - [ ] Actions: toggleSidebar, setTheme

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

## Phase 3: 通用组件库 (Common Components)

### 3.1 基础 UI 组件 ❌

- [ ] Button
  - [ ] `src/components/common/Button.tsx`
  - [ ] Variants: primary, secondary, danger, outline
  - [ ] Sizes: small, medium, large
  - [ ] Loading state, disabled state
- [ ] Input
  - [ ] `src/components/common/Input.tsx`
  - [ ] Text, textarea, file input
  - [ ] Error state, helper text
- [ ] Select / Dropdown
  - [ ] `src/components/common/Select.tsx`
  - [ ] 单选、多选支持
- [ ] Modal / Dialog
  - [ ] `src/components/common/Modal.tsx`
  - [ ] 确认框、表单弹窗
- [ ] Card
  - [ ] `src/components/common/Card.tsx`
  - [ ] 内容卡片容器

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 4-5 天

---

### 3.2 反馈组件 ❌

- [ ] Loading Spinner
  - [ ] `src/components/common/LoadingSpinner.tsx`
  - [ ] 全屏加载、局部加载
- [ ] Error Message
  - [ ] `src/components/common/ErrorMessage.tsx`
  - [ ] 错误提示组件
- [ ] Empty State
  - [ ] `src/components/common/EmptyState.tsx`
  - [ ] 空数据提示
- [ ] Progress Bar
  - [ ] `src/components/common/ProgressBar.tsx`
  - [ ] 进度条 (用于生成进度)
- [ ] Toast 通知
  - [ ] 使用 react-hot-toast
  - [ ] 成功、错误、警告提示

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 2-3 天

---

### 3.3 布局组件 ❌

- [ ] Header / Navbar
  - [ ] `src/components/common/Header.tsx`
  - [ ] Logo, 导航菜单, 用户信息
- [ ] Sidebar
  - [ ] `src/components/common/Sidebar.tsx`
  - [ ] 可折叠侧边栏
- [ ] Footer
  - [ ] `src/components/common/Footer.tsx`
  - [ ] 版权信息、链接
- [ ] Container
  - [ ] `src/components/common/Container.tsx`
  - [ ] 内容容器,最大宽度限制

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 2-3 天

---

## Phase 4: 功能页面 - 小说管理 (Novel Management)

### 4.1 首页 ❌

- [ ] 首页设计
  - [ ] `src/pages/HomePage.tsx`
  - [ ] 欢迎信息
  - [ ] 快速操作入口 (上传小说、查看列表)
  - [ ] 统计信息展示 (总小说数、总场景数)

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 1-2 天

---

### 4.2 小说上传 ❌

#### 上传组件

- [ ] NovelUpload 组件
  - [ ] `src/components/features/novel/NovelUpload.tsx`
  - [ ] 表单: 标题、作者、内容文本框
  - [ ] 文件上传支持 (TXT, EPUB)
  - [ ] 拖拽上传功能
  - [ ] 上传进度条
  - [ ] 表单验证 (必填项、文件大小限制)
- [ ] 上传成功处理
  - [ ] Toast 提示
  - [ ] 跳转到小说详情页

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

### 4.3 小说列表 ❌

#### 列表页面

- [ ] NovelListPage
  - [ ] `src/pages/NovelListPage.tsx`
  - [ ] 小说列表展示 (分页)
  - [ ] 搜索/筛选功能
  - [ ] 排序功能 (按时间、标题)
- [ ] NovelCard 组件
  - [ ] `src/components/features/novel/NovelCard.tsx`
  - [ ] 显示: 标题、作者、状态、创建时间
  - [ ] 操作: 查看、删除
  - [ ] 状态徽章 (uploaded, parsing, completed)
- [ ] 分页组件
  - [ ] `src/components/common/Pagination.tsx`
  - [ ] 页码跳转、上一页/下一页

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

### 4.4 小说详情 ❌

#### 详情页面

- [ ] NovelDetailPage
  - [ ] `src/pages/NovelDetailPage.tsx`
  - [ ] 小说基本信息 (标题、作者、状态)
  - [ ] 章节列表 (可折叠)
  - [ ] 操作按钮: 解析、查看角色、生成场景、导出
- [ ] NovelDetail 组件
  - [ ] `src/components/features/novel/NovelDetail.tsx`
  - [ ] 小说元数据展示
  - [ ] 文本预览 (前 500 字)
  - [ ] 编辑/删除操作
- [ ] ChapterList 组件
  - [ ] `src/components/features/novel/ChapterList.tsx`
  - [ ] 章节标题、字数、状态
  - [ ] 展开查看章节内容

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 4-5 天

---

## Phase 5: 功能页面 - 角色管理 (Character Management)

### 5.1 角色列表 ❌

#### 角色页面

- [ ] CharacterPage
  - [ ] `src/pages/CharacterPage.tsx`
  - [ ] 角色列表展示 (卡片或表格)
  - [ ] 添加/编辑角色功能
  - [ ] 生成参考图按钮
- [ ] CharacterList 组件
  - [ ] `src/components/features/character/CharacterList.tsx`
  - [ ] 角色卡片: 名字、外貌、性格、参考图
  - [ ] 筛选: 主要角色 / 配角
- [ ] CharacterCard 组件
  - [ ] `src/components/features/character/CharacterCard.tsx`
  - [ ] 角色头像 (参考图)
  - [ ] 角色信息展示
  - [ ] 操作: 编辑、生成参考图、删除

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 4-5 天

---

### 5.2 角色编辑 ❌

#### 编辑组件

- [ ] CharacterEditor 组件
  - [ ] `src/components/features/character/CharacterEditor.tsx`
  - [ ] 表单: 名字、外貌描述、性格描述
  - [ ] 参考图管理: 上传、删除、预览
  - [ ] 保存/取消按钮
- [ ] 表单验证
  - [ ] 必填项检查
  - [ ] 字数限制

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 2-3 天

---

### 5.3 参考图生成 ❌

#### 生成功能

- [ ] ReferenceImageGenerator 组件
  - [ ] `src/components/features/character/ReferenceImageGenerator.tsx`
  - [ ] 触发生成按钮
  - [ ] 生成进度显示
  - [ ] 生成结果预览
  - [ ] 重新生成/接受/拒绝选项
- [ ] 图片预览组件
  - [ ] `src/components/common/ImagePreview.tsx`
  - [ ] 图片放大查看
  - [ ] 图片下载

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

## Phase 6: 功能页面 - 场景生成 (Scene Generation)

### 6.1 场景列表 ❌

#### 场景管理页面

- [ ] SceneGenerationPage
  - [ ] `src/pages/SceneGenerationPage.tsx`
  - [ ] 场景列表 (按章节分组)
  - [ ] 批量选择场景
  - [ ] 批量生成按钮
- [ ] SceneList 组件
  - [ ] `src/components/features/scene/SceneList.tsx`
  - [ ] 场景卡片: 序号、描述、角色、状态
  - [ ] 场景预览图 (如已生成)
  - [ ] 单个生成按钮
- [ ] SceneCard 组件
  - [ ] `src/components/features/scene/SceneCard.tsx`
  - [ ] 场景描述文本
  - [ ] 涉及角色列表
  - [ ] 生成状态: 未生成、生成中、已完成

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 4-5 天

---

### 6.2 场景生成控制 ❌

#### 生成组件

- [ ] SceneGenerator 组件
  - [ ] `src/components/features/generation/SceneGenerator.tsx`
  - [ ] 生成类型选择: 图片、视频、图片+视频
  - [ ] 生成参数配置 (风格、尺寸、时长)
  - [ ] 开始生成按钮
  - [ ] 批量生成队列管理
- [ ] GenerationProgress 组件
  - [ ] `src/components/features/generation/GenerationProgress.tsx`
  - [ ] 进度条 (当前场景 / 总场景)
  - [ ] 实时状态更新
  - [ ] 取消生成按钮
  - [ ] 错误提示

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 5-6 天

---

### 6.3 生成结果预览 ❌

#### 预览组件

- [ ] GenerationResult 组件
  - [ ] `src/components/features/generation/GenerationResult.tsx`
  - [ ] 图片/视频预览
  - [ ] 画廊视图 (多场景浏览)
  - [ ] 下载按钮
  - [ ] 重新生成按钮
- [ ] MediaViewer 组件
  - [ ] `src/components/common/MediaViewer.tsx`
  - [ ] 图片查看器 (缩放、平移)
  - [ ] 视频播放器 (播放、暂停、进度条)

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 3-4 天

---

## Phase 7: 功能页面 - 导出管理 (Export)

### 7.1 导出页面 ❌

#### 导出配置

- [ ] ExportPage
  - [ ] `src/pages/ExportPage.tsx`
  - [ ] 导出格式选择 (MP4, MOV, AVI)
  - [ ] 导出质量选择 (720p, 1080p, 4K)
  - [ ] 音频设置 (背景音乐、音效、旁白)
  - [ ] 字幕设置 (是否显示、语言)
  - [ ] 开始导出按钮
- [ ] ExportConfig 组件
  - [ ] `src/components/features/export/ExportConfig.tsx`
  - [ ] 配置表单
  - [ ] 预估文件大小
  - [ ] 预估导出时间
- [ ] ExportProgress 组件
  - [ ] `src/components/features/export/ExportProgress.tsx`
  - [ ] 导出进度条
  - [ ] 当前步骤显示 (视频拼接、音频合成、字幕生成)
  - [ ] 完成后下载链接

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 5-6 天

---

## Phase 8: 高级功能 (Advanced Features)

### 8.1 实时通知 ❌

#### WebSocket 集成

- [ ] WebSocket 客户端
  - [ ] `src/services/websocket.ts`
  - [ ] 连接管理 (重连机制)
  - [ ] 消息订阅/取消订阅
- [ ] useWebSocket Hook
  - [ ] `src/hooks/useWebSocket.ts`
  - [ ] 监听生成进度更新
  - [ ] 监听任务完成通知
- [ ] 实时通知 UI
  - [ ] 右上角通知图标
  - [ ] 通知列表 (下拉框)
  - [ ] 消息标记已读

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 4-5 天

---

### 8.2 主题切换 ❌

- [ ] 深色模式支持
  - [ ] CSS 变量切换
  - [ ] 主题切换按钮
  - [ ] 本地存储记住用户偏好
- [ ] 主题配置
  - [ ] `src/styles/themes/`
  - [ ] light.css, dark.css

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 2-3 天

---

### 8.3 国际化 (i18n) ❌

- [ ] i18n 库集成
  - [ ] react-i18next
  - [ ] `src/locales/` 目录
  - [ ] zh-CN.json, en-US.json
- [ ] 语言切换
  - [ ] 语言选择器组件
  - [ ] 本地存储语言偏好

**完成度:** 0%
**优先级:** P3 (可选)
**预计工期:** 3-4 天

---

## Phase 9: 自定义 Hooks (Custom Hooks)

### 9.1 数据获取 Hooks ❌

- [ ] useNovel Hook
  - [ ] `src/hooks/useNovel.ts`
  - [ ] 获取单个小说,返回 {novel, loading, error, refetch}
- [ ] useNovels Hook
  - [ ] `src/hooks/useNovels.ts`
  - [ ] 获取小说列表,支持分页
- [ ] useCharacters Hook
  - [ ] `src/hooks/useCharacters.ts`
  - [ ] 获取角色列表
- [ ] useGeneration Hook
  - [ ] `src/hooks/useGeneration.ts`
  - [ ] 生成场景,返回任务状态

**完成度:** 0%
**优先级:** P0 (高)
**预计工期:** 3-4 天

---

### 9.2 工具 Hooks ❌

- [ ] useDebounce Hook
  - [ ] `src/hooks/useDebounce.ts`
  - [ ] 防抖搜索输入
- [ ] useLocalStorage Hook
  - [ ] `src/hooks/useLocalStorage.ts`
  - [ ] 本地存储状态管理
- [ ] useIntersectionObserver Hook
  - [ ] `src/hooks/useIntersectionObserver.ts`
  - [ ] 无限滚动加载

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 1-2 天

---

## Phase 10: 测试 (Testing)

### 10.1 单元测试 ❌

#### 测试环境配置

- [ ] 测试库安装
  - [ ] @testing-library/react
  - [ ] @testing-library/jest-dom
  - [ ] vitest (Vite 推荐)
- [ ] 测试配置
  - [ ] `vitest.config.ts`
  - [ ] 测试工具函数

#### 组件测试

- [ ] 通用组件测试
  - [ ] Button.test.tsx
  - [ ] Input.test.tsx
  - [ ] Modal.test.tsx
- [ ] 功能组件测试
  - [ ] NovelCard.test.tsx
  - [ ] CharacterCard.test.tsx
- [ ] Hooks 测试
  - [ ] useNovel.test.ts
  - [ ] useCharacters.test.ts

**目标覆盖率:** 70%+
**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 5-7 天

---

### 10.2 E2E 测试 ❌

- [ ] E2E 测试工具
  - [ ] Playwright 或 Cypress
- [ ] 关键流程测试
  - [ ] 小说上传流程
  - [ ] 角色生成参考图流程
  - [ ] 场景批量生成流程

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 4-5 天

---

## Phase 11: 性能优化 (Performance Optimization)

### 11.1 代码分割 ❌

- [ ] 路由懒加载
  - [ ] 所有页面组件使用 React.lazy
- [ ] 组件懒加载
  - [ ] 大型组件按需加载 (Modal, ImageViewer)
- [ ] Vite 配置优化
  - [ ] `vite.config.ts` 手动分包
  - [ ] vendor, api, ui 分别打包

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 1-2 天

---

### 11.2 性能优化 ❌

- [ ] React 优化
  - [ ] 使用 React.memo 包裹组件
  - [ ] 使用 useMemo 缓存计算值
  - [ ] 使用 useCallback 缓存函数
- [ ] 虚拟滚动
  - [ ] 长列表使用 react-window 或 react-virtual
- [ ] 图片优化
  - [ ] 图片懒加载
  - [ ] 响应式图片 (srcset)
  - [ ] WebP 格式支持

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 3-4 天

---

## Phase 12: 部署与 DevOps (Deployment)

### 12.1 构建优化 ❌

- [ ] 生产构建配置
  - [ ] 环境变量管理 (.env.production)
  - [ ] Source map 配置
  - [ ] 压缩优化
- [ ] 静态资源优化
  - [ ] Gzip / Brotli 压缩
  - [ ] CDN 配置 (可选)

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 1-2 天

---

### 12.2 Docker 部署 🚧

- [ ] Dockerfile 编写
  - [ ] 多阶段构建 (build + nginx)
  - [ ] 优化镜像大小
- [ ] Nginx 配置
  - [ ] `nginx.conf`
  - [ ] SPA 路由支持 (try_files)
  - [ ] API 代理配置

**完成度:** 30%
**优先级:** P0 (高)
**预计工期:** 1-2 天

---

### 12.3 CI/CD ❌

- [ ] GitHub Actions
  - [ ] 自动化测试
  - [ ] 自动化构建
  - [ ] 自动化部署 (可选)
- [ ] 代码质量检查
  - [ ] ESLint 检查
  - [ ] TypeScript 类型检查
  - [ ] 测试覆盖率检查

**完成度:** 0%
**优先级:** P1 (中)
**预计工期:** 2-3 天

---

## Phase 13: 文档与规范 (Documentation)

### 13.1 组件文档 ❌

- [ ] Storybook 集成 (可选)
  - [ ] 组件展示和文档
  - [ ] 交互式示例
- [ ] README 完善
  - [ ] 开发指南
  - [ ] 组件使用说明

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 3-4 天

---

### 13.2 用户指南 ❌

- [ ] 用户手册
  - [ ] 功能介绍
  - [ ] 操作指南 (截图)
- [ ] FAQ 页面
  - [ ] 常见问题解答

**完成度:** 0%
**优先级:** P2 (低)
**预计工期:** 2-3 天

---

## 总体进度

| Phase | 名称 | 完成度 | 状态 |
|-------|------|--------|------|
| Phase 1 | 项目基础 | 40% | 🚧 进行中 |
| Phase 2 | 核心基础设施 | 0% | ❌ 未开始 |
| Phase 3 | 通用组件库 | 0% | ❌ 未开始 |
| Phase 4 | 小说管理 | 0% | ❌ 未开始 |
| Phase 5 | 角色管理 | 0% | ❌ 未开始 |
| Phase 6 | 场景生成 | 0% | ❌ 未开始 |
| Phase 7 | 导出管理 | 0% | ❌ 未开始 |
| Phase 8 | 高级功能 | 0% | ❌ 未开始 |
| Phase 9 | 自定义 Hooks | 0% | ❌ 未开始 |
| Phase 10 | 测试 | 0% | ❌ 未开始 |
| Phase 11 | 性能优化 | 0% | ❌ 未开始 |
| Phase 12 | 部署与 DevOps | 10% | ⏳ 计划中 |
| Phase 13 | 文档与规范 | 0% | ❌ 未开始 |

**总体完成度:** ~10%

---

## 里程碑 (Milestones)

### M1: MVP - v0.2.0 (20% 完成)
- ✅ 项目初始化
- ⏳ 基础路由
- ⏳ API 集成
- ⏳ 小说上传功能
- ⏳ 小说列表展示

**目标:** 完成基本的小说管理功能

---

### M2: Alpha - v0.5.0 (0% 完成)
- 角色管理功能
- 参考图生成 UI
- 场景列表展示
- 基础生成功能

**目标:** 完成角色和场景管理

---

### M3: Beta - v0.8.0 (0% 完成)
- 完整的场景生成流程
- 实时进度更新
- 导出功能
- 性能优化

**目标:** 完整的用户体验

---

### M4: 正式版 - v1.0.0 (0% 完成)
- 所有功能完善
- 测试覆盖率 70%+
- 生产环境部署
- 完善文档

**目标:** 生产环境可用

---

## 优先级说明

- **P0 (高)**: 核心功能,必须完成
- **P1 (中)**: 重要功能,尽快完成
- **P2 (低)**: 优化功能,可延后
- **P3 (可选)**: 扩展功能,按需开发

---

## 技术栈总结

### 核心依赖
- **React** 19 - UI 框架
- **TypeScript** - 类型系统
- **Vite** 7 - 构建工具
- **React Router** - 路由管理

### 推荐依赖
- **Zustand** - 状态管理
- **Axios** - HTTP 客户端
- **Tailwind CSS** - 样式框架
- **react-hot-toast** - Toast 通知
- **react-icons** - 图标库

### 开发工具
- **ESLint** - 代码检查
- **Prettier** - 代码格式化
- **Vitest** - 单元测试
- **Playwright** - E2E 测试

---

## 下一步行动

### 本周计划
1. [ ] 选择并配置样式方案 (Tailwind CSS)
2. [ ] 完成路由配置和页面骨架
3. [ ] 实现 API 服务层
4. [ ] 定义 TypeScript 类型
5. [ ] 创建基础 UI 组件 (Button, Input, Modal)

### 本月计划
1. [ ] 完成 Phase 1-3 (基础设施 + 通用组件)
2. [ ] 完成 Phase 4 (小说管理功能)
3. [ ] 开始 Phase 5 (角色管理功能)
4. [ ] 达到 M1 (MVP) 里程碑

---

## 风险与挑战

### 技术风险

1. **生成任务实时更新**
   - 风险: WebSocket 连接不稳定
   - 缓解: 实现轮询降级方案

2. **大文件上传**
   - 风险: 超时、内存溢出
   - 缓解: 分块上传、进度显示

3. **图片/视频预览性能**
   - 风险: 大量媒体文件卡顿
   - 缓解: 懒加载、虚拟滚动

### UX 风险

1. **生成等待时间长**
   - 风险: 用户流失
   - 缓解: 进度可视化、后台任务

2. **复杂操作流程**
   - 风险: 用户不理解
   - 缓解: 引导提示、帮助文档

---

## 参考文档

- [Frontend CLAUDE.md](./CLAUDE.md) - 前端开发指南
- [API.md](../docs/API.md) - API 接口文档
- [ARCHITECTURE.md](../docs/ARCHITECTURE.md) - 系统架构
- [README.md](../README.md) - 项目概览

---

## 更新日志

- **2025-10-24**: 创建初始路线图,定义 13 个开发阶段
