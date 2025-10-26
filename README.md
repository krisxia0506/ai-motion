# AI-Motion - 智能动漫生成系统

基于 AI 技术的小说到动漫自动生成系统，支持角色一致性、智能配音和图文展示。

---

## 项目简介

AI-Motion 是一个创新的自动化动漫生成平台，能够将文字小说转换为视觉化的动漫内容。

### 核心特性

- 📖 **小说解析** - 自动分析小说情节、角色和场景
- 🎨 **角色一致性** - 确保同一角色在整个动漫中保持视觉一致性
- 🖼️ **图文生成** - 生成高质量的场景图片配合文字展示
- 🎙️ **智能配音** - 自动为角色生成语音，支持多角色语音区分
- 🎬 **动漫输出** - 以图配文 + 声音的形式呈现动漫内容

### 技术栈

- **后端**: Go 1.24+ | Gin | DDD 架构 | Supabase (PostgreSQL)
- **前端**: React 19 + TypeScript | Vite 7 | Tailwind CSS
- **AI 服务**: Gemini 2.5 Flash Image (文生图/图生图) | Sora2 (视频生成)
- **DevOps**: Docker + Docker Compose
- **设计系统**: 主题色 #2FB2F1 (青蓝色) | 现代化渐变设计

---

## 项目结构

```
ai-motion/
├── backend/              # Go 后端 (DDD 架构)
│   ├── cmd/             # 应用入口
│   ├── internal/        # 业务逻辑
│   │   ├── domain/     # 领域层
│   │   ├── application/ # 应用层
│   │   ├── infrastructure/ # 基础设施层
│   │   └── interfaces/ # 接口层
│   └── pkg/            # 公共包
├── frontend/            # React 前端
│   └── src/
│       ├── components/ # UI 组件
│       ├── pages/      # 页面
│       └── services/   # API 服务
├── docs/               # 项目文档
├── docker/             # Docker 配置
└── scripts/            # 工具脚本
```

详细架构说明请查看 [架构文档](docs/ARCHITECTURE.md)。

---

## 快速开始

### Docker 部署 (推荐)

```bash
# 1. 克隆项目
git clone https://github.com/xiajiayi/ai-motion.git
cd ai-motion

# 2. 配置环境变量
cp .env.example .env
vim .env  # 填入数据库和 AI API 配置

# 3. 启动服务
docker-compose up -d

# 4. 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8080
```

### 本地开发

```bash
make install  # 安装依赖
make dev      # 启动开发环境
```

详细指南请查看 [快速启动文档](QUICKSTART.md) 和 [部署文档](docs/DEPLOYMENT.md)。

---

## 文档

- [快速启动指南](QUICKSTART.md) - 快速开始使用
- [API 接口文档](docs/API.md) - 完整的 API 参考
- [部署指南](docs/DEPLOYMENT.md) - 生产环境部署
- [配置说明](docs/CONFIGURATION.md) - 详细配置指南
- [开发指南](docs/DEVELOPMENT.md) - 开发者文档
- [架构文档](docs/ARCHITECTURE.md) - 系统架构设计
- [故障排查](docs/TROUBLESHOOTING.md) - 常见问题解决

---

## 常用命令

```bash
make help         # 查看所有命令
make install      # 安装依赖
make dev          # 启动开发环境
make build        # 编译项目
make docker-up    # 启动 Docker 服务
make docker-down  # 停止 Docker 服务
make test         # 运行测试
```

---

## 项目状态

**当前版本**: v0.1.0-alpha (开发中)

基础架构已搭建完成，核心功能正在开发中。

### 开发路线图

- [x] 项目初始化与容器化部署
- [x] 基础架构搭建
- [ ] DDD 架构完整实现
- [ ] 小说解析引擎
- [ ] AI 服务集成 (Gemini / Sora2)
- [ ] 角色一致性方案
- [ ] 前端完整界面
- [ ] 用户系统
- [ ] 生产环境优化

---

## 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

详细贡献指南请查看 [开发文档](docs/DEVELOPMENT.md)。

---

## 许可证

MIT License - 详见 LICENSE 文件

---

## 联系方式

- 提交 [Issue](https://github.com/krisxia0506/ai-motion/issues)
- 查看 [讨论区](https://github.com/krisxia0506/ai-motion/discussions)

---

⭐ 如果觉得项目有帮助，欢迎 Star 支持！
