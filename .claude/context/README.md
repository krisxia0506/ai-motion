# Context Files

这个目录包含了可以按需导入的上下文文件，用于减少主 CLAUDE.md 中不必要的内容加载。

## 文件说明

- **commands.md** - 构建、运行、测试和 Docker 命令
- **git-workflow.md** - Git 提交规范和分支策略
- **development-tasks.md** - 常见开发任务的分步指南
- **troubleshooting.md** - 常见问题和解决方案
- **project-status.md** - 当前开发状态和路线图

## 使用方式

在与 Claude Code 对话时，如果需要特定上下文，可以通过以下方式引用:

```
@.claude/context/commands.md
```

或者在 CLAUDE.md 文件中使用 `@import` 语法:

```markdown
@.claude/context/commands.md
```

## 设计原则

这些文件被设计为:
- **按需加载**: 只有在需要时才引用，避免每次请求都加载所有上下文
- **单一职责**: 每个文件只包含一个特定主题的内容
- **易于维护**: 独立的文件更容易更新和管理
