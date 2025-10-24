# AI-Motion 数据库表结构

本目录包含 AI-Motion 项目的所有数据库表结构定义,遵循[阿里巴巴开发手册](https://github.com/alibaba/p3c)的数据库设计规范。

## 目录结构

```
database/schema/
├── README.md                              # 本文件
├── init.sql                               # 初始化脚本(按顺序执行所有表)
├── 01_aimotion_novel.sql                  # 小说表
├── 02_aimotion_novel_content.sql          # 小说内容表
├── 03_aimotion_chapter.sql                # 章节表
├── 04_aimotion_character.sql              # 角色表
├── 05_aimotion_character_image.sql        # 角色图片表
├── 06_aimotion_scene.sql                  # 场景表
├── 07_aimotion_scene_character.sql        # 场景角色关联表
└── 08_aimotion_media.sql                  # 媒体表
```

## 表说明

### 核心表

| 文件名 | 表名 | 说明 | 关键字段 |
|--------|------|------|----------|
| `01_aimotion_novel.sql` | aimotion_novel | 小说基本信息 | title, author, status |
| `02_aimotion_novel_content.sql` | aimotion_novel_content | 小说内容(大字段独立) | novel_id, content |
| `03_aimotion_chapter.sql` | aimotion_chapter | 章节信息 | novel_id, sequence_num |
| `04_aimotion_character.sql` | aimotion_character | 角色信息 | novel_id, name |
| `05_aimotion_character_image.sql` | aimotion_character_image | 角色图片(参考图/场景图) | character_id, image_url |
| `06_aimotion_scene.sql` | aimotion_scene | 场景信息 | chapter_id, sequence_num |
| `07_aimotion_scene_character.sql` | aimotion_scene_character | 场景角色关联 | scene_id, character_id |
| `08_aimotion_media.sql` | aimotion_media | 生成的媒体(图片/视频) | scene_id, url |

## 使用方法

### 方式1: 使用 init.sql 初始化所有表

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS aimotion CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行初始化脚本
mysql -u root -p aimotion < database/schema/init.sql
```

### 方式2: 单独执行每个表的 SQL 文件

```bash
# 按顺序执行
mysql -u root -p aimotion < database/schema/01_aimotion_novel.sql
mysql -u root -p aimotion < database/schema/02_aimotion_novel_content.sql
mysql -u root -p aimotion < database/schema/03_aimotion_chapter.sql
# ... 依次执行其他文件
```

### 方式3: 使用 Docker Compose

```bash
# Docker Compose 会自动执行 init.sql
docker-compose up -d mysql
```

## 设计规范

所有表结构遵循以下规范:

### 1. 命名规范
- 表名: `aimotion_` 前缀 + 业务名称
- 字段名: 小写字母 + 下划线
- 索引名: `idx_` (普通索引) 或 `uk_` (唯一索引) + 字段名

### 2. 必备字段
每张表都包含以下必备字段:
- `id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY` - 主键ID
- `is_deleted TINYINT UNSIGNED DEFAULT 0` - 逻辑删除标识
- `gmt_create DATETIME NOT NULL` - 创建时间
- `gmt_modified DATETIME NOT NULL` - 修改时间

### 3. 数据类型
- 主键: `BIGINT UNSIGNED AUTO_INCREMENT`
- 非负整数: `INT UNSIGNED`, `TINYINT UNSIGNED`
- 小数: `DECIMAL(M,D)` (禁用 FLOAT/DOUBLE)
- 短文本: `VARCHAR(N)` (N ≤ 5000)
- 长文本: `TEXT`, `LONGTEXT` (独立表)

### 4. 索引设计
- 频繁查询的字段建立索引
- 外键关联字段建立索引
- 唯一性约束使用唯一索引
- 避免在频繁更新的字段上建索引

### 5. 约束规范
- **禁用外键约束** - 所有关联关系在应用层维护
- 使用逻辑删除 (`is_deleted`) 而非物理删除

## 表关系图

```
aimotion_novel (小说)
    ├── 1:1 → aimotion_novel_content (内容)
    ├── 1:N → aimotion_chapter (章节)
    ├── 1:N → aimotion_character (角色)
    └── 1:N → aimotion_scene (场景)

aimotion_character (角色)
    └── 1:N → aimotion_character_image (图片)

aimotion_scene (场景)
    ├── N:M → aimotion_character (通过 aimotion_scene_character)
    └── 1:N → aimotion_media (媒体)
```

## 注意事项

1. **大字段拆分**: `aimotion_novel_content` 独立存储小说内容,避免影响主表查询性能
2. **关联表设计**: `aimotion_scene_character` 是场景和角色的多对多关联表
3. **逻辑删除**: 所有删除操作使用 `is_deleted=1` 标记,不物理删除数据
4. **字符集**: 统一使用 `utf8mb4` 支持 emoji 和特殊字符
5. **索引优化**: 根据实际查询需求调整索引,避免过度索引

## 相关文档

- [数据库设计规范](../../docs/DATABASE_DESIGN_STANDARDS.md) - 完整的设计规范文档
- [架构文档](../../docs/ARCHITECTURE.md) - 系统架构和数据模型说明
- [开发文档](../../docs/DEVELOPMENT.md) - 开发环境搭建和规范

## 后续工作

- [ ] 创建数据库迁移脚本 (使用 golang-migrate 或类似工具)
- [ ] 更新 Go entity 定义,匹配新的字段名
- [ ] 实现 Repository 层的 CRUD 操作
- [ ] 编写数据库集成测试
- [ ] 添加数据库性能监控
