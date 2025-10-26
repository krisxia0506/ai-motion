# AI-Motion 数据库设计规范

本文档基于《阿里巴巴Java开发手册》的数据库设计规范,结合 AI-Motion 项目实际情况制定。

**注意**: AI-Motion 现已迁移至 Supabase (PostgreSQL + PostgREST)。本文档同时包含 PostgreSQL 相关规范。

## 目录

- [命名规范](#命名规范)
- [字段规范](#字段规范)
- [索引规范](#索引规范)
- [约束规范](#约束规范)
- [表设计规范](#表设计规范)
- [标准模板](#标准模板)

---

## 命名规范

### 1. 库名规范

**规则**: 库名与应用名称尽量一致

```sql
-- ✅ 正确 (PostgreSQL)
CREATE DATABASE aimotion WITH ENCODING 'UTF8';

-- ✅ 正确 (MySQL - 已废弃)
CREATE DATABASE aimotion CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ❌ 错误
CREATE DATABASE ai_motion_db;
CREATE DATABASE app_database;
```

### 2. 表名规范

**规则**:
- 表名、字段名必须使用小写字母或数字
- 禁止数字开头
- 表名不使用复数名词
- 表的命名最好是 "业务名称_表的作用"

```sql
-- ✅ 正确
aimotion_novel
aimotion_character
aimotion_scene
aimotion_media

-- ❌ 错误
novels           -- 使用了复数
Novel            -- 使用了大写
1_novel          -- 数字开头
novel            -- 缺少业务前缀
```

### 3. 字段名规范

**规则**:
- 必须使用小写字母或数字
- 禁止数字开头
- 表达是与否概念的字段,必须使用 `is_xxx` 命名
- 时间字段使用 `gmt_create`(创建时间) 和 `gmt_modified`(修改时间)

```sql
-- ✅ 正确
is_deleted TINYINT UNSIGNED DEFAULT 0
gmt_create DATETIME NOT NULL
gmt_modified DATETIME NOT NULL

-- ❌ 错误
deleted BOOLEAN              -- 未使用 is_ 前缀
created_at TIMESTAMP         -- 未使用 gmt_create
isDeleted TINYINT            -- 驼峰命名
```

---

## 字段规范

### 1. 必备三字段

**规则**: 每张表必须包含 `id`, `gmt_create`, `gmt_modified`, `is_deleted`

**PostgreSQL 版本**:
```sql
CREATE TABLE aimotion_example (
    id BIGSERIAL PRIMARY KEY,
    -- 业务字段...
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_example IS '示例表';
COMMENT ON COLUMN aimotion_example.id IS '主键ID';
COMMENT ON COLUMN aimotion_example.is_deleted IS '逻辑删除:0-未删除,1-已删除';
COMMENT ON COLUMN aimotion_example.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_example.gmt_modified IS '修改时间';
```

**MySQL 版本 (已废弃)**:
```sql
CREATE TABLE aimotion_example (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    -- 业务字段...
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除:0-未删除,1-已删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='示例表';
```

**说明**:
- `id`: 必为主键,PostgreSQL 使用 `BIGSERIAL` 自增 (MySQL 使用 `BIGINT UNSIGNED AUTO_INCREMENT`)
- 如果使用分库分表,则 `id` 类型为 `VARCHAR`,非自增,使用分布式ID生成器
- `gmt_create`: 创建时间,PostgreSQL 使用 `TIMESTAMP` (MySQL 使用 `DATETIME`)
- `gmt_modified`: 修改时间,PostgreSQL 使用 `TIMESTAMP` (MySQL 使用 `DATETIME`)
- `is_deleted`: 逻辑删除标识,PostgreSQL 使用 `SMALLINT` (MySQL 使用 `TINYINT UNSIGNED`),1表示删除,0表示未删除

### 2. 数据类型选择

#### 整数类型
- **PostgreSQL**: 主键使用 `BIGSERIAL`,其他整数使用 `INTEGER`, `SMALLINT` 等,添加 `CHECK` 约束确保非负
- **MySQL (已废弃)**: 非负数必须使用 `UNSIGNED`,主键使用 `BIGINT UNSIGNED AUTO_INCREMENT`

```sql
-- ✅ 正确 (PostgreSQL)
id BIGSERIAL PRIMARY KEY
sequence_num INTEGER CHECK (sequence_num >= 0)
is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1))

-- ✅ 正确 (MySQL - 已废弃)
id BIGINT UNSIGNED AUTO_INCREMENT
sequence_num INT UNSIGNED
is_deleted TINYINT UNSIGNED

-- ❌ 错误
id INT                    -- 应使用 BIGSERIAL (PostgreSQL) 或 BIGINT UNSIGNED (MySQL)
age INT                   -- 年龄非负,应添加 CHECK 约束或使用 UNSIGNED
```

#### 小数类型
- 小数类型为 `DECIMAL`,禁止使用 `FLOAT` 和 `DOUBLE`
- `FLOAT` 和 `DOUBLE` 存在精度损失问题

```sql
-- ✅ 正确
price DECIMAL(10, 2) COMMENT '价格(元)'
rate DECIMAL(5, 4) COMMENT '比率'

-- ❌ 错误
price FLOAT              -- 禁止使用 FLOAT
rate DOUBLE              -- 禁止使用 DOUBLE
```

#### 字符串类型
- 定长字符串使用 `CHAR`
- 可变长字符串使用 `VARCHAR`,长度不要超过 5000
- 超长文本使用 `TEXT`,独立出来一张表

```sql
-- ✅ 正确 (PostgreSQL)
country_code CHAR(2)
title VARCHAR(255)
name VARCHAR(100)

COMMENT ON COLUMN table_name.country_code IS '国家代码';
COMMENT ON COLUMN table_name.title IS '标题';
COMMENT ON COLUMN table_name.name IS '姓名';

-- 超长文本独立表 (PostgreSQL)
CREATE TABLE aimotion_novel_content (
    id BIGSERIAL PRIMARY KEY,
    novel_id BIGINT NOT NULL,
    content TEXT,
    gmt_create TIMESTAMP NOT NULL,
    gmt_modified TIMESTAMP NOT NULL
);

CREATE INDEX idx_novel_id ON aimotion_novel_content(novel_id);

-- ❌ 错误
content VARCHAR(10000)    -- 超过5000,应使用 TEXT 并独立表
```

#### 布尔类型
- **PostgreSQL**: 可使用 `BOOLEAN` 或 `SMALLINT` 配合 `CHECK` 约束
- **MySQL (已废弃)**: 使用 `TINYINT UNSIGNED` 表示是与否
- 必须使用 `is_xxx` 命名

```sql
-- ✅ 正确 (PostgreSQL - 推荐)
is_deleted BOOLEAN DEFAULT FALSE
is_active BOOLEAN DEFAULT TRUE

-- ✅ 正确 (PostgreSQL - 兼容方式)
is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1))
is_active SMALLINT DEFAULT 1 CHECK (is_active IN (0, 1))

-- ✅ 正确 (MySQL - 已废弃)
is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '0-未删除,1-已删除'
is_active TINYINT UNSIGNED DEFAULT 1 COMMENT '0-未激活,1-已激活'

-- ❌ 错误
deleted BOOLEAN           -- 缺少 is_ 前缀
active SMALLINT           -- 缺少 CHECK 约束
```

---

## 索引规范

### 1. 索引命名规范

**规则**:
- 主键索引: `PRIMARY KEY`
- 唯一索引: `uk_字段名` (uk = unique key)
- 普通索引: `idx_字段名` (idx = index)

```sql
-- ✅ 正确
PRIMARY KEY (id)
UNIQUE KEY uk_novel_id_sequence (novel_id, sequence_num)
INDEX idx_novel_id (novel_id)
INDEX idx_gmt_create (gmt_create)

-- ❌ 错误
INDEX novel_id_index (novel_id)
UNIQUE KEY novel_sequence (novel_id, sequence_num)
```

### 2. 索引设计原则

- 单表索引数量不超过 5 个
- 单个索引字段数不超过 5 个
- 频繁查询的字段建立索引
- 频繁更新的字段不建索引
- 区分度低的字段不建索引 (如性别)

```sql
-- ✅ 正确示例
CREATE TABLE aimotion_scene (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    novel_id BIGINT UNSIGNED NOT NULL,
    chapter_id BIGINT UNSIGNED NOT NULL,
    sequence_num INT UNSIGNED NOT NULL,
    
    INDEX idx_novel_id (novel_id),              -- 频繁按小说查询
    INDEX idx_chapter_id (chapter_id),          -- 频繁按章节查询
    UNIQUE KEY uk_chapter_sequence (chapter_id, sequence_num)  -- 保证章节内序号唯一
);
```

---

## 约束规范

### 1. 禁用外键约束

**规则**: 不得使用外键与级联,一切外键概念必须在应用层解决

**原因**:
- 外键与级联更新适用于单机低并发,不适合分布式、高并发集群
- 级联更新是强阻塞,存在数据库更新风暴的风险
- 外键影响数据库的插入速度

```sql
-- ❌ 错误
CREATE TABLE aimotion_chapter (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    novel_id BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id) ON DELETE CASCADE
);

-- ✅ 正确 - 在应用层维护关联关系
CREATE TABLE aimotion_chapter (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '关联小说ID',
    INDEX idx_novel_id (novel_id)
);
```

### 2. 应用层关联实现

在应用层通过代码维护关联关系:

```go
// 删除小说时,应用层处理关联数据
func (s *NovelService) Delete(ctx context.Context, novelID int64) error {
    // 1. 删除关联章节
    if err := s.chapterRepo.DeleteByNovelID(ctx, novelID); err != nil {
        return err
    }
    
    // 2. 删除关联角色
    if err := s.characterRepo.DeleteByNovelID(ctx, novelID); err != nil {
        return err
    }
    
    // 3. 删除小说
    return s.novelRepo.Delete(ctx, novelID)
}
```

---

## 表设计规范

### 1. 分库分表阈值

**规则**: 单表行数超过 500 万行或单表容量超过 2GB,才推荐进行分库分表

**说明**: 如果预计三年后的数据量达不到这个级别,请不要在创建表时就分库分表

```sql
-- AI-Motion 当前阶段不需要分库分表
-- 预估数据量:
-- - aimotion_novel: < 10万行
-- - aimotion_chapter: < 500万行
-- - aimotion_character: < 100万行
-- - aimotion_scene: < 1000万行 (可能需要分表)
-- - aimotion_media: < 2000万行 (可能需要分表)
```

### 2. 垂直拆分

大字段独立表:

```sql
-- 主表
CREATE TABLE aimotion_novel (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL COMMENT '标题',
    author VARCHAR(100) COMMENT '作者',
    status TINYINT UNSIGNED DEFAULT 0 COMMENT '状态',
    is_deleted TINYINT UNSIGNED DEFAULT 0,
    gmt_create DATETIME NOT NULL,
    gmt_modified DATETIME NOT NULL,
    INDEX idx_gmt_create (gmt_create)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小说表';

-- 内容独立表
CREATE TABLE aimotion_novel_content (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    content LONGTEXT COMMENT '小说内容',
    gmt_create DATETIME NOT NULL,
    gmt_modified DATETIME NOT NULL,
    UNIQUE KEY uk_novel_id (novel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小说内容表';
```

### 3. 水平拆分 (未来考虑)

当 `aimotion_scene` 或 `aimotion_media` 表数据量过大时,可按小说ID哈希分表:

```sql
-- 分表示例 (当前不需要实现)
aimotion_scene_0
aimotion_scene_1
...
aimotion_scene_15  -- 16张表
```

---

## 标准模板

### 1. 基础表模板

**PostgreSQL 版本**:
```sql
CREATE TABLE aimotion_example (
    id BIGSERIAL PRIMARY KEY,
    
    -- 业务字段
    name VARCHAR(100) NOT NULL,
    type SMALLINT DEFAULT 0 CHECK (type IN (0, 1)),
    amount DECIMAL(10, 2) DEFAULT 0.00,
    
    -- 必备字段
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加注释
COMMENT ON TABLE aimotion_example IS '示例表';
COMMENT ON COLUMN aimotion_example.id IS '主键ID';
COMMENT ON COLUMN aimotion_example.name IS '名称';
COMMENT ON COLUMN aimotion_example.type IS '类型:0-类型A,1-类型B';
COMMENT ON COLUMN aimotion_example.amount IS '金额';
COMMENT ON COLUMN aimotion_example.is_deleted IS '逻辑删除:0-未删除,1-已删除';
COMMENT ON COLUMN aimotion_example.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_example.gmt_modified IS '修改时间';

-- 索引
CREATE INDEX idx_type ON aimotion_example(type);
CREATE INDEX idx_gmt_create ON aimotion_example(gmt_create);
```

### 2. 关联表模板

**PostgreSQL 版本**:
```sql
CREATE TABLE aimotion_scene_character (
    id BIGSERIAL PRIMARY KEY,
    scene_id BIGINT NOT NULL,
    character_id BIGINT NOT NULL,
    role_type SMALLINT DEFAULT 0 CHECK (role_type IN (0, 1)),
    
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_scene_character UNIQUE (scene_id, character_id)
);

-- 注释
COMMENT ON TABLE aimotion_scene_character IS '场景角色关联表';
COMMENT ON COLUMN aimotion_scene_character.id IS '主键ID';
COMMENT ON COLUMN aimotion_scene_character.scene_id IS '场景ID';
COMMENT ON COLUMN aimotion_scene_character.character_id IS '角色ID';
COMMENT ON COLUMN aimotion_scene_character.role_type IS '角色类型:0-主角,1-配角';
COMMENT ON COLUMN aimotion_scene_character.is_deleted IS '逻辑删除';
COMMENT ON COLUMN aimotion_scene_character.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_scene_character.gmt_modified IS '修改时间';

-- 索引
CREATE INDEX idx_character_id ON aimotion_scene_character(character_id);
```

### 3. AI-Motion 核心表设计

#### 小说表

```sql
CREATE TABLE aimotion_novel (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    title VARCHAR(255) NOT NULL COMMENT '小说标题',
    author VARCHAR(100) COMMENT '作者',
    status TINYINT UNSIGNED DEFAULT 0 COMMENT '状态:0-草稿,1-解析中,2-已完成',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除:0-未删除,1-已删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_status (status),
    INDEX idx_gmt_create (gmt_create)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小说表';
```

#### 小说内容表

```sql
CREATE TABLE aimotion_novel_content (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    content LONGTEXT COMMENT '小说内容',
    
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_novel_id (novel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小说内容表';
```

#### 章节表

```sql
CREATE TABLE aimotion_chapter (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    title VARCHAR(255) COMMENT '章节标题',
    sequence_num INT UNSIGNED NOT NULL COMMENT '章节序号',
    content TEXT COMMENT '章节内容',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    UNIQUE KEY uk_novel_sequence (novel_id, sequence_num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='章节表';
```

#### 角色表

```sql
CREATE TABLE aimotion_character (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    name VARCHAR(100) NOT NULL COMMENT '角色名称',
    appearance TEXT COMMENT '外貌描述',
    personality TEXT COMMENT '性格描述',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
```

#### 角色参考图表

```sql
CREATE TABLE aimotion_character_image (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    image_url VARCHAR(512) NOT NULL COMMENT '图片URL',
    image_type TINYINT UNSIGNED DEFAULT 0 COMMENT '图片类型:0-参考图,1-场景图',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色图片表';
```

#### 场景表

```sql
CREATE TABLE aimotion_scene (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    chapter_id BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    sequence_num INT UNSIGNED NOT NULL COMMENT '场景序号',
    description TEXT COMMENT '场景描述',
    location VARCHAR(255) COMMENT '地点',
    time_of_day VARCHAR(50) COMMENT '时间段',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    INDEX idx_chapter_id (chapter_id),
    UNIQUE KEY uk_chapter_sequence (chapter_id, sequence_num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场景表';
```

#### 场景角色关联表

```sql
CREATE TABLE aimotion_scene_character (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    scene_id BIGINT UNSIGNED NOT NULL COMMENT '场景ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_scene_character (scene_id, character_id),
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='场景角色关联表';
```

#### 媒体表

```sql
CREATE TABLE aimotion_media (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    type TINYINT UNSIGNED NOT NULL COMMENT '媒体类型:0-图片,1-视频',
    scene_id BIGINT UNSIGNED NOT NULL COMMENT '场景ID',
    url VARCHAR(512) NOT NULL COMMENT '媒体URL',
    metadata JSON COMMENT '元数据',
    status TINYINT UNSIGNED DEFAULT 0 COMMENT '状态:0-生成中,1-已完成,2-失败',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_scene_id (scene_id),
    INDEX idx_type (type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='媒体表';
```

---

## 规范清单

基于《阿里巴巴Java开发手册》,AI-Motion 项目遵循以下规范:

- [x] **规则1**: 库名与应用名称一致 (`aimotion`)
- [x] **规则2**: 表名、字段名使用小写字母或数字,禁止数字开头
- [x] **规则3**: 表名不使用复数名词
- [x] **规则4**: 表名加业务前缀 `aimotion_`
- [x] **规则5**: 必备四字段 `id`, `gmt_create`, `gmt_modified`, `is_deleted`
- [x] **规则6**: 单表行数超过 500 万行或 2GB 才分库分表
- [x] **规则7**: 布尔字段使用 `is_xxx` 命名,类型为 `TINYINT UNSIGNED`
- [x] **规则8**: 小数类型使用 `DECIMAL`,禁用 `FLOAT` 和 `DOUBLE`
- [x] **规则9**: 定长字符串使用 `CHAR`
- [x] **规则10**: `VARCHAR` 不超过 5000,超长文本独立表
- [x] **规则11**: 唯一索引 `uk_`,普通索引 `idx_`
- [x] **规则12**: 禁用外键,应用层解决关联

---

## 后续工作

1. **数据库迁移**: 根据本规范创建数据库迁移脚本
2. **Go Model 适配**: 更新 `backend/internal/domain/**/entity.go` 中的字段名
3. **Repository 实现**: 更新 `backend/internal/infrastructure/repository/mysql/*.go`
4. **集成测试**: 编写数据库集成测试,验证新 schema
5. **文档同步**: 更新 `API.md` 和 `DEVELOPMENT.md` 中的数据库相关说明

---

## 参考资料

- [阿里巴巴Java开发手册 - MySQL数据库规约](https://github.com/alibaba/p3c)
- [MySQL官方文档](https://dev.mysql.com/doc/)
