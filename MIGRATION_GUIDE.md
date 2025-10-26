# Media-Novel Association Migration Guide

## 问题描述

之前的后端流程存在以下问题：

1. **Task** 有 `novel_id` 字段
2. **Media** 只有 `scene_id` 字段，但没有 `novel_id`
3. 当查询任务详情时，无法直接通过 Novel 查询到关联的 Media
4. 在漫画生成流程中，我们不使用 Scene，直接生成与 Novel 关联的图片

## 解决方案

添加 `novel_id` 字段到 Media 表，建立以下关系：
- **Task → Novel** (已存在)
- **Novel → Media** (新增)
- **Media → Scene** (可选，保留用于场景生成模式)

## 已完成的代码更改

### 1. 数据库迁移
- ✅ 创建迁移文件: `000008_add_novel_id_to_media.up.sql`
- ✅ 添加 `novel_id` 列到 `aimotion_media` 表
- ✅ 创建索引以优化查询性能
- ✅ 添加外键约束确保数据完整性
- ✅ 将 `scene_id` 改为可选（支持漫画生成模式）

### 2. 领域实体更新
**文件**: `backend/internal/domain/media/entity.go`
- ✅ 添加 `NovelID` 字段到 `Media` 结构体
- ✅ 添加 `NewMediaForNovel()` 构造函数（用于漫画生成）
- ✅ 更新验证逻辑（`NovelID` 或 `SceneID` 至少需要一个）

### 3. 仓储接口和实现
**文件**:
- `backend/internal/domain/media/repository.go`
- `backend/internal/infrastructure/repository/supabase/media_repository.go`
- `backend/internal/infrastructure/repository/mysql/media_repository.go`

更改：
- ✅ 添加 `FindByNovelID(ctx, novelID)` 方法
- ✅ 更新 `Save()` 方法以保存 `novel_id`
- ✅ 更新所有查询方法以包含 `novel_id` 字段
- ✅ 更新映射函数以正确处理 `novel_id`

### 4. 应用服务更新
**文件**: `backend/internal/application/service/manga_workflow_service.go`

更改：
- ✅ 使用 `media.NewMediaForNovel(novelID, ...)` 而不是 `media.NewMedia(sceneID, ...)`
- ✅ 在 `buildTaskResult()` 中使用 `FindByNovelID()` 查询媒体文件
- ✅ 移除了基于硬编码 ID 的查询逻辑

## 需要手动执行的步骤

### ⚠️ 重要：执行数据库迁移

由于我们使用 Supabase 作为数据库，需要手动执行迁移：

1. 打开 Supabase 项目控制台
2. 导航到 **SQL Editor**
3. 打开文件 `MIGRATION_MANUAL.sql` (在项目根目录)
4. 复制所有内容并粘贴到 SQL Editor
5. 点击 **Run** 执行迁移
6. 验证输出，确保迁移成功

## 数据流程

### 之前的流程（有问题）
```
用户请求 → Task (novel_id) → Novel
                                  ↓
                              (无法直接查询)
                                  ↓
                              Media (只有 scene_id)
```

### 现在的流程（已修复）
```
用户请求 → Task (novel_id) → Novel (id)
                                  ↓
                          通过 novel_id 查询
                                  ↓
                              Media (novel_id, scene_id?)
                                  ↓
                              返回所有关联的图片
```

## API 测试

在执行迁移后，测试以下流程：

### 1. 创建漫画生成任务
```bash
curl -X POST http://localhost:8080/api/v1/manga/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试故事",
    "author": "测试作者",
    "content": "这是一个测试故事的内容，需要足够长以通过验证..."
  }'
```

### 2. 查询任务状态（应该能看到生成的媒体文件）
```bash
curl http://localhost:8080/api/v1/manga/task/{task_id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

预期响应应包含：
```json
{
  "task_id": "...",
  "status": "completed",
  "result": {
    "novel_id": "...",
    "title": "测试故事",
    "scene_count": 10,
    "scenes": [
      {
        "id": "...",
        "sequence_num": 1,
        "description": "漫画面板 1",
        "image_url": "https://..."
      },
      ...
    ]
  }
}
```

## 验证检查清单

- [ ] 在 Supabase SQL Editor 中执行了 `MIGRATION_MANUAL.sql`
- [ ] 迁移成功完成（没有错误）
- [ ] 后端服务正常启动（`curl http://localhost:8080/health`）
- [ ] 可以创建新的漫画生成任务
- [ ] 查询任务详情时能看到生成的图片列表
- [ ] 图片 URL 可以正常访问

## 回滚计划

如果需要回滚，在 Supabase SQL Editor 中运行：

```sql
-- 删除外键约束
ALTER TABLE aimotion_media
DROP CONSTRAINT IF EXISTS fk_aimotion_media_novel;

-- 删除索引
DROP INDEX IF EXISTS idx_aimotion_media_novel_id;

-- 删除 novel_id 列
ALTER TABLE aimotion_media
DROP COLUMN IF EXISTS novel_id;

-- 恢复 scene_id 为 NOT NULL（如果需要）
-- 注意：只有在没有 scene_id 为 NULL 的记录时才能执行
-- ALTER TABLE aimotion_media
-- ALTER COLUMN scene_id SET NOT NULL;
```

## 技术细节

### 数据库架构变更

**之前**:
```sql
CREATE TABLE aimotion_media (
    id VARCHAR(36) PRIMARY KEY,
    scene_id VARCHAR(36) NOT NULL,  -- 必填，依赖 scene
    ...
);
```

**之后**:
```sql
CREATE TABLE aimotion_media (
    id VARCHAR(36) PRIMARY KEY,
    novel_id VARCHAR(36),           -- 新增：关联小说
    scene_id VARCHAR(36),           -- 改为可选
    ...
    FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id)
);
```

### 代码变更总结

1. **实体层**: 添加 `NovelID` 字段和新构造函数
2. **仓储层**: 添加 `FindByNovelID()` 方法
3. **应用层**: 修改工作流使用新的关联方式

## 常见问题

### Q: 为什么不直接使用 scene_id？
A: 在简化的漫画生成流程中，我们不需要创建 Scene 实体。直接关联 Novel 更简单高效。

### Q: Scene 模式还能用吗？
A: 可以。`scene_id` 仍然存在且功能完整，只是改为可选。如果未来需要更复杂的场景管理，可以使用 Scene 模式。

### Q: 旧的媒体记录怎么办？
A: 旧记录的 `novel_id` 会是 NULL。新生成的媒体会自动填充 `novel_id`。

## 总结

此次修改建立了 Task、Novel 和 Media 之间的正确关联关系，使得：
- ✅ 查询任务详情时可以直接获取所有关联的媒体文件
- ✅ 数据模型更清晰，符合业务逻辑
- ✅ 支持两种模式：直接小说生成（漫画）和场景生成（动画）
- ✅ 查询性能更好（通过 novel_id 索引）

执行完手动迁移后，整个系统即可正常工作！
